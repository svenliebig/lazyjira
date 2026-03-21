package views

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/jira"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// issueItem wraps a jira.Issue to implement list.Item.
type issueItem struct {
	issue jira.Issue
}

func (i issueItem) Title() string {
	return fmt.Sprintf("%s  %s", shared.StyleIssueKey.Render(i.issue.Key), i.issue.Fields.Summary)
}

func (i issueItem) Description() string {
	status := i.issue.Fields.Status.Name
	assignee := "unassigned"
	if i.issue.Fields.Assignee != nil {
		assignee = i.issue.Fields.Assignee.DisplayName
	}
	return fmt.Sprintf("%s · %s", status, assignee)
}

func (i issueItem) FilterValue() string {
	return i.issue.Key + " " + i.issue.Fields.Summary
}

// IssueListModel is the view for listing issues.
type IssueListModel struct {
	list        list.Model
	issues      []jira.Issue
	initialized bool
}

func NewIssueListModel(issues []jira.Issue, width, height int) IssueListModel {
	items := make([]list.Item, len(issues))
	for i, issue := range issues {
		items[i] = issueItem{issue: issue}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, width, height)
	l.Title = "Assigned Issues"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	return IssueListModel{
		list:        l,
		issues:      issues,
		initialized: true,
	}
}

func (m IssueListModel) Init() tea.Cmd {
	return nil
}

func (m *IssueListModel) SetSize(w, h int) {
	if !m.initialized {
		return
	}
	m.list.SetSize(w, h)
}

func (m IssueListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == shared.KeyEnter {
			if item, ok := m.list.SelectedItem().(issueItem); ok {
				return m, func() tea.Msg {
					return shared.IssueSelectedMsg{Issue: item.issue}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m IssueListModel) View() string {
	return m.list.View()
}
