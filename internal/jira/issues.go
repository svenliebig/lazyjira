package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// adfToText converts an Atlassian Document Format JSON object to plain text.
func adfToText(raw json.RawMessage) string {
	if raw == nil {
		return ""
	}
	var node map[string]json.RawMessage
	if err := json.Unmarshal(raw, &node); err != nil {
		// Try as string
		var s string
		if err2 := json.Unmarshal(raw, &s); err2 == nil {
			return s
		}
		return ""
	}

	var sb strings.Builder

	// Check for text field
	if textRaw, ok := node["text"]; ok {
		var text string
		if err := json.Unmarshal(textRaw, &text); err == nil {
			sb.WriteString(text)
		}
	}

	// Recurse into content array
	if contentRaw, ok := node["content"]; ok {
		var items []json.RawMessage
		if err := json.Unmarshal(contentRaw, &items); err == nil {
			for _, item := range items {
				sb.WriteString(adfToText(item))
			}
		}
	}

	// Add newline after block nodes
	if typeRaw, ok := node["type"]; ok {
		var nodeType string
		if err := json.Unmarshal(typeRaw, &nodeType); err == nil {
			switch nodeType {
			case "paragraph", "heading", "bulletList", "orderedList", "listItem", "blockquote", "codeBlock", "rule":
				if sb.Len() > 0 {
					sb.WriteString("\n")
				}
			}
		}
	}

	return sb.String()
}

// Response types for JSON unmarshaling

type searchResponse struct {
	Issues []issueResponse `json:"issues"`
}

type issueResponse struct {
	ID     string              `json:"id"`
	Key    string              `json:"key"`
	Fields issueFieldsResponse `json:"fields"`
}

type issueFieldsResponse struct {
	Summary      string               `json:"summary"`
	Description  json.RawMessage      `json:"description"`
	Status       statusResponse       `json:"status"`
	Assignee     *userResponse        `json:"assignee"`
	Reporter     *userResponse        `json:"reporter"`
	Parent       *parentResponse      `json:"parent"`
	Sprints      []sprintResponse     `json:"customfield_10020"`
	TimeTracking timeTrackingResponse `json:"timetracking"`
}

type sprintResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

type timeTrackingResponse struct {
	OriginalEstimateSeconds  int `json:"originalEstimateSeconds"`
	RemainingEstimateSeconds int `json:"remainingEstimateSeconds"`
}

type parentResponse struct {
	Key string `json:"key"`
}

type statusResponse struct {
	Name string `json:"name"`
}

type userResponse struct {
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
}

type transitionsResponse struct {
	Transitions []transitionResponse `json:"transitions"`
}

type transitionResponse struct {
	ID   string             `json:"id"`
	Name string             `json:"name"`
	To   transitionToResponse `json:"to"`
}

type transitionToResponse struct {
	Name string `json:"name"`
}

func convertIssue(r issueResponse) Issue {
	issue := Issue{
		ID:  r.ID,
		Key: r.Key,
		Fields: IssueFields{
			Summary:     r.Fields.Summary,
			Description: adfToText(r.Fields.Description),
			Status:      IssueStatus{Name: r.Fields.Status.Name},
		},
	}
	if r.Fields.Assignee != nil {
		issue.Fields.Assignee = &User{
			DisplayName:  r.Fields.Assignee.DisplayName,
			EmailAddress: r.Fields.Assignee.EmailAddress,
		}
	}
	if r.Fields.Reporter != nil {
		issue.Fields.Reporter = &User{
			DisplayName:  r.Fields.Reporter.DisplayName,
			EmailAddress: r.Fields.Reporter.EmailAddress,
		}
	}
	if r.Fields.Parent != nil {
		issue.Fields.Parent = &IssueParent{Key: r.Fields.Parent.Key}
	}
	if len(r.Fields.Sprints) > 0 {
		s := r.Fields.Sprints[len(r.Fields.Sprints)-1]
		issue.Fields.Sprint = &Sprint{ID: s.ID, Name: s.Name, State: s.State}
	}
	issue.Fields.TimeTracking = TimeTracking{
		OriginalEstimateSeconds:  r.Fields.TimeTracking.OriginalEstimateSeconds,
		RemainingEstimateSeconds: r.Fields.TimeTracking.RemainingEstimateSeconds,
	}
	return issue
}

// ListAssigned returns issues assigned to the current user.
func (c *Client) ListAssigned(ctx context.Context) ([]Issue, error) {
	payload := map[string]interface{}{
		"jql":        "assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC",
		"maxResults": 50,
		"fields":     []string{"summary", "description", "status", "assignee", "reporter", "parent", "customfield_10020", "timetracking"},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(http.MethodPost, "/rest/api/3/search/jql", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	issues := make([]Issue, len(result.Issues))
	for i, r := range result.Issues {
		issues[i] = convertIssue(r)
	}
	return issues, nil
}

// GetIssue returns a single issue by key.
func (c *Client) GetIssue(ctx context.Context, key string) (*Issue, error) {
	req, err := c.newRequest(http.MethodGet, "/rest/api/3/issue/"+key, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var r issueResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	issue := convertIssue(r)
	return &issue, nil
}

// GetTransitions returns available transitions for an issue.
func (c *Client) GetTransitions(ctx context.Context, key string) ([]Transition, error) {
	req, err := c.newRequest(http.MethodGet, "/rest/api/3/issue/"+key+"/transitions", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result transitionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	transitions := make([]Transition, len(result.Transitions))
	for i, t := range result.Transitions {
		transitions[i] = Transition{
			ID:   t.ID,
			Name: t.Name,
			To:   TransitionTo{Name: t.To.Name},
		}
	}
	return transitions, nil
}

// DoTransition performs a transition on an issue.
func (c *Client) DoTransition(ctx context.Context, key, transitionID string) error {
	payload := map[string]interface{}{
		"transition": map[string]string{
			"id": transitionID,
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest(http.MethodPost, "/rest/api/3/issue/"+key+"/transitions", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
