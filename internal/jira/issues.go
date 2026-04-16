package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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

	var nodeType string
	if typeRaw, ok := node["type"]; ok {
		json.Unmarshal(typeRaw, &nodeType) //nolint:errcheck
	}

	// Special handling for table: collect rows/cells and render with alignment
	if nodeType == "table" {
		return adfTableToText(node)
	}

	// Special handling for date: render timestamp as YYYY-MM-DD
	if nodeType == "date" {
		return adfDateToText(node)
	}

	var sb strings.Builder

	// taskItem: prefix with checkbox indicator
	if nodeType == "taskItem" {
		state := "TODO"
		if attrsRaw, ok := node["attrs"]; ok {
			var attrs map[string]json.RawMessage
			if json.Unmarshal(attrsRaw, &attrs) == nil {
				if stateRaw, ok2 := attrs["state"]; ok2 {
					json.Unmarshal(stateRaw, &state) //nolint:errcheck
				}
			}
		}
		if state == "DONE" {
			sb.WriteString("[x] ")
		} else {
			sb.WriteString("[ ] ")
		}
	}

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
	switch nodeType {
	case "paragraph", "heading", "bulletList", "orderedList", "listItem", "blockquote", "codeBlock", "rule", "taskItem":
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// adfTableToText renders an ADF table node as aligned plain text.
func adfTableToText(node map[string]json.RawMessage) string {
	contentRaw, ok := node["content"]
	if !ok {
		return ""
	}
	var rowNodes []json.RawMessage
	if err := json.Unmarshal(contentRaw, &rowNodes); err != nil {
		return ""
	}

	type tableRow struct {
		cells    []string
		isHeader bool
	}
	var rows []tableRow

	for _, rowRaw := range rowNodes {
		var rowNode map[string]json.RawMessage
		if err := json.Unmarshal(rowRaw, &rowNode); err != nil {
			continue
		}
		cellsRaw, ok := rowNode["content"]
		if !ok {
			continue
		}
		var cellNodes []json.RawMessage
		if err := json.Unmarshal(cellsRaw, &cellNodes); err != nil {
			continue
		}

		var row tableRow
		for i, cellRaw := range cellNodes {
			var cellNode map[string]json.RawMessage
			if err := json.Unmarshal(cellRaw, &cellNode); err != nil {
				continue
			}
			if i == 0 {
				var cellType string
				if typeRaw, ok := cellNode["type"]; ok {
					json.Unmarshal(typeRaw, &cellType) //nolint:errcheck
				}
				row.isHeader = cellType == "tableHeader"
			}
			row.cells = append(row.cells, strings.TrimSpace(adfToText(cellRaw)))
		}
		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return ""
	}

	// Compute column widths
	numCols := 0
	for _, row := range rows {
		if len(row.cells) > numCols {
			numCols = len(row.cells)
		}
	}
	colWidths := make([]int, numCols)
	for _, row := range rows {
		for i, cell := range row.cells {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	var sb strings.Builder
	for _, row := range rows {
		for i := 0; i < numCols; i++ {
			if i > 0 {
				sb.WriteString(" | ")
			}
			cell := ""
			if i < len(row.cells) {
				cell = row.cells[i]
			}
			sb.WriteString(cell)
			// Pad all but the last column
			if i < numCols-1 {
				for j := len(cell); j < colWidths[i]; j++ {
					sb.WriteByte(' ')
				}
			}
		}
		sb.WriteString("\n")
		if row.isHeader {
			for i := 0; i < numCols; i++ {
				if i > 0 {
					sb.WriteString("-+-")
				}
				for j := 0; j < colWidths[i]; j++ {
					sb.WriteByte('-')
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// adfDateToText converts an ADF date node (attrs.timestamp in Unix ms) to YYYY-MM-DD.
func adfDateToText(node map[string]json.RawMessage) string {
	attrsRaw, ok := node["attrs"]
	if !ok {
		return ""
	}
	var attrs map[string]json.RawMessage
	if err := json.Unmarshal(attrsRaw, &attrs); err != nil {
		return ""
	}
	tsRaw, ok := attrs["timestamp"]
	if !ok {
		return ""
	}
	// Timestamp may be a JSON string or number
	var tsMS int64
	var tsStr string
	if err := json.Unmarshal(tsRaw, &tsStr); err == nil {
		tsMS, _ = strconv.ParseInt(tsStr, 10, 64)
	} else {
		json.Unmarshal(tsRaw, &tsMS) //nolint:errcheck
	}
	if tsMS == 0 {
		return ""
	}
	return time.Unix(tsMS/1000, 0).UTC().Format("2006-01-02")
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
	AccountID    string `json:"accountId"`
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
			AccountID:    r.Fields.Assignee.AccountID,
			DisplayName:  r.Fields.Assignee.DisplayName,
			EmailAddress: r.Fields.Assignee.EmailAddress,
		}
	}
	if r.Fields.Reporter != nil {
		issue.Fields.Reporter = &User{
			AccountID:    r.Fields.Reporter.AccountID,
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

// UnassignIssue removes the assignee from an issue.
func (c *Client) UnassignIssue(ctx context.Context, key string) error {
	payload := map[string]interface{}{
		"accountId": nil,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest(http.MethodPut, "/rest/api/3/issue/"+key+"/assignee", bytes.NewReader(data))
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

// SearchAssignableUsers returns users that can be assigned to the given issue.
func (c *Client) SearchAssignableUsers(ctx context.Context, issueKey string) ([]User, error) {
	req, err := c.newRequest(http.MethodGet, "/rest/api/3/user/assignable/search?issueKey="+issueKey+"&maxResults=50", nil)
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

	var results []userResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	users := make([]User, len(results))
	for i, u := range results {
		users[i] = User{
			AccountID:    u.AccountID,
			DisplayName:  u.DisplayName,
			EmailAddress: u.EmailAddress,
		}
	}
	return users, nil
}

// AssignIssue sets the assignee of an issue to the given account ID.
func (c *Client) AssignIssue(ctx context.Context, key, accountID string) error {
	payload := map[string]interface{}{
		"accountId": accountID,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest(http.MethodPut, "/rest/api/3/issue/"+key+"/assignee", bytes.NewReader(data))
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
