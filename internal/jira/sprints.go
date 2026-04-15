package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type agileSprintListResponse struct {
	Values []agileSprintItem `json:"values"`
}

type agileSprintItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// GetSprints returns active and future sprints for the given board.
func (c *Client) GetSprints(ctx context.Context, boardID int) ([]Sprint, error) {
	path := fmt.Sprintf("/rest/agile/1.0/board/%d/sprint?state=active,future", boardID)
	req, err := c.newRequest(http.MethodGet, path, nil)
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

	var result agileSprintListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	sprints := make([]Sprint, len(result.Values))
	for i, v := range result.Values {
		sprints[i] = Sprint{ID: v.ID, Name: v.Name, State: v.State}
	}
	return sprints, nil
}

// MoveIssueToSprint moves an issue into the given sprint.
func (c *Client) MoveIssueToSprint(ctx context.Context, issueKey string, sprintID int) error {
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue", sprintID)
	body := strings.NewReader(fmt.Sprintf(`{"issues":["%s"]}`, issueKey))

	req, err := c.newRequest(http.MethodPost, path, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
