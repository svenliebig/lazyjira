package shared

import "github.com/svenliebig/lazyjira/internal/jira"

// Auth
type AuthCompletedMsg struct {
	URL   string
	Email string
	Token string
}

// Issues
type IssueListLoadedMsg struct{ Issues []jira.Issue }
type IssueSelectedMsg struct{ Issue jira.Issue }
type TransitionsLoadedMsg struct{ Transitions []jira.Transition }
type TransitionDoneMsg struct{}
type TransitionSelectedMsg struct{ ID string }

// List selection
type ListSelectedMsg struct{ Type string }

// AI
type AICommitsLoadedMsg struct{ Commits []string }
type AISummaryMsg struct{ Summary string }

// Errors
type ErrMsg struct{ Err error }

func (e ErrMsg) Error() string { return e.Err.Error() }

// Copy actions
type CopyMsg struct{ Text string }
type CopyActionMsg struct{ Action string }

// Close modal
type CloseModalMsg struct{}

// Exclusions
type ExcludeActionMsg struct {
	Type  string // "key" or "parent"
	Value string // issue key or parent key
}

// Unassign
type UnassignDoneMsg struct{}

// Assign
type UsersLoadedMsg struct{ Users []jira.User }
type UserSelectedMsg struct{ User jira.User }
type AssignDoneMsg struct{ User jira.User }

// Theme
type ThemeSelectedMsg struct{ Name string }
