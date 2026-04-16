package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListAssigned(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Basic dGVzdEBleGFtcGxlLmNvbTp0ZXN0LXRva2Vu" {
			t.Errorf("Expected Basic auth, got %s", r.Header.Get("Authorization"))
		}

		resp := searchResponse{
			Issues: []issueResponse{
				{
					ID:  "10001",
					Key: "PROJ-1",
					Fields: issueFieldsResponse{
						Summary:     "Test issue",
						Description: json.RawMessage(`{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Test description"}]}]}`),
						Status:      statusResponse{Name: "In Progress"},
						Assignee: &userResponse{
							DisplayName:  "John Doe",
							EmailAddress: "john@example.com",
						},
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "test-token")
	issues, err := client.ListAssigned(context.Background())
	if err != nil {
		t.Fatalf("ListAssigned() returned error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("Expected 1 issue, got %d", len(issues))
	}

	issue := issues[0]
	if issue.ID != "10001" {
		t.Errorf("Expected ID %q, got %q", "10001", issue.ID)
	}
	if issue.Key != "PROJ-1" {
		t.Errorf("Expected Key %q, got %q", "PROJ-1", issue.Key)
	}
	if issue.Fields.Summary != "Test issue" {
		t.Errorf("Expected Summary %q, got %q", "Test issue", issue.Fields.Summary)
	}
	if issue.Fields.Status.Name != "In Progress" {
		t.Errorf("Expected Status %q, got %q", "In Progress", issue.Fields.Status.Name)
	}
	if issue.Fields.Assignee == nil {
		t.Fatal("Expected non-nil Assignee")
	}
	if issue.Fields.Assignee.DisplayName != "John Doe" {
		t.Errorf("Expected Assignee %q, got %q", "John Doe", issue.Fields.Assignee.DisplayName)
	}
	if issue.Fields.Description == "" {
		t.Error("Expected non-empty description")
	}
}

func TestGetTransitions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		resp := transitionsResponse{
			Transitions: []transitionResponse{
				{
					ID:   "11",
					Name: "To Do",
					To:   transitionToResponse{Name: "To Do"},
				},
				{
					ID:   "21",
					Name: "In Progress",
					To:   transitionToResponse{Name: "In Progress"},
				},
				{
					ID:   "31",
					Name: "Done",
					To:   transitionToResponse{Name: "Done"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "test-token")
	transitions, err := client.GetTransitions(context.Background(), "PROJ-1")
	if err != nil {
		t.Fatalf("GetTransitions() returned error: %v", err)
	}
	if len(transitions) != 3 {
		t.Fatalf("Expected 3 transitions, got %d", len(transitions))
	}

	if transitions[0].ID != "11" {
		t.Errorf("Expected first transition ID %q, got %q", "11", transitions[0].ID)
	}
	if transitions[0].Name != "To Do" {
		t.Errorf("Expected first transition Name %q, got %q", "To Do", transitions[0].Name)
	}
	if transitions[2].Name != "Done" {
		t.Errorf("Expected last transition Name %q, got %q", "Done", transitions[2].Name)
	}
}

func TestListAssigned_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := searchResponse{Issues: []issueResponse{}}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "test-token")
	issues, err := client.ListAssigned(context.Background())
	if err != nil {
		t.Fatalf("ListAssigned() returned error: %v", err)
	}
	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}

func TestAdfToText_Table(t *testing.T) {
	raw := json.RawMessage(`{
		"type": "doc",
		"version": 1,
		"content": [
			{
				"type": "paragraph",
				"content": [{"type": "text", "text": "Intro text"}]
			},
			{
				"type": "table",
				"content": [
					{
						"type": "tableRow",
						"content": [
							{"type": "tableHeader", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Image"}]}]},
							{"type": "tableHeader", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Tag"}]}]},
							{"type": "tableHeader", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "Comment"}]}]}
						]
					},
					{
						"type": "tableRow",
						"content": [
							{"type": "tableCell", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "backend-core"}]}]},
							{"type": "tableCell", "content": [{"type": "paragraph", "content": [{"type": "text", "text": "main-abc123"}]}]},
							{"type": "tableCell", "content": [
								{"type": "taskList", "content": [
									{"type": "taskItem", "attrs": {"state": "TODO"}, "content": [
										{"type": "text", "text": "Deployed at "},
										{"type": "date", "attrs": {"timestamp": "1776297600000"}}
									]}
								]}
							]}
						]
					}
				]
			}
		]
	}`)

	got := adfToText(raw)

	if !strings.Contains(got, "Image") {
		t.Error("Expected table header 'Image' in output")
	}
	if !strings.Contains(got, "Tag") {
		t.Error("Expected table header 'Tag' in output")
	}
	if !strings.Contains(got, "backend-core") {
		t.Error("Expected 'backend-core' in output")
	}
	if !strings.Contains(got, " | ") {
		t.Error("Expected column separator ' | ' in output")
	}
	if !strings.Contains(got, "---") {
		t.Error("Expected header separator '---' in output")
	}
	if !strings.Contains(got, "[ ]") {
		t.Error("Expected task item '[ ]' in output")
	}
	if !strings.Contains(got, "2026-04-16") {
		t.Error("Expected date '2026-04-16' in output")
	}
	if !strings.Contains(got, "Intro text") {
		t.Error("Expected 'Intro text' in output")
	}
}

func TestGetIssue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := issueResponse{
			ID:  "10002",
			Key: "PROJ-2",
			Fields: issueFieldsResponse{
				Summary:     "Another issue",
				Description: json.RawMessage(`{"type":"doc","content":[]}`),
				Status:      statusResponse{Name: "Done"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test@example.com", "test-token")
	issue, err := client.GetIssue(context.Background(), "PROJ-2")
	if err != nil {
		t.Fatalf("GetIssue() returned error: %v", err)
	}
	if issue.Key != "PROJ-2" {
		t.Errorf("Expected Key %q, got %q", "PROJ-2", issue.Key)
	}
	if issue.Fields.Status.Name != "Done" {
		t.Errorf("Expected Status %q, got %q", "Done", issue.Fields.Status.Name)
	}
}
