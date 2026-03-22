package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	assignee := "unassigned"
	if i.issue.Fields.Assignee != nil {
		assignee = i.issue.Fields.Assignee.DisplayName
	}
	return fmt.Sprintf("%s · %s", i.issue.Fields.Status.Name, assignee)
}

func (i issueItem) FilterValue() string {
	return i.issue.Key + " " + i.issue.Fields.Summary
}

// IssueListModel is a split-panel view: issue list on the left, detail on the right.
type IssueListModel struct {
	list        list.Model
	detail      viewport.Model
	issues      []jira.Issue
	initialized bool
	focusRight  bool
	width       int
	height      int
}

// splitWidths returns the content widths for the left and right panels.
// The left panel uses a right border (1 char), so total = leftContent + 1 + rightContent.
func splitWidths(total int) (leftContent, rightContent int) {
	leftContent = total * 2 / 5
	rightContent = total - leftContent - 1
	return
}

func NewIssueListModel(issues []jira.Issue, width, height int) IssueListModel {
	leftW, rightW := splitWidths(width)

	items := make([]list.Item, len(issues))
	for i, issue := range issues {
		items[i] = issueItem{issue: issue}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, leftW, height)
	l.Title = "Issues"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(false)

	vp := viewport.New(rightW-2, height)

	m := IssueListModel{
		list:        l,
		detail:      vp,
		issues:      issues,
		initialized: true,
		width:       width,
		height:      height,
	}
	if len(issues) > 0 {
		m.detail.SetContent(buildIssueDetail(issues[0], rightW-4))
	}
	return m
}

// CurrentIssue returns the currently highlighted issue, or nil if the list is empty.
func (m IssueListModel) CurrentIssue() *jira.Issue {
	if len(m.issues) == 0 {
		return nil
	}
	idx := m.list.Index()
	if idx < 0 || idx >= len(m.issues) {
		return nil
	}
	issue := m.issues[idx]
	return &issue
}

// IsFocusRight reports whether focus is in the detail panel.
func (m IssueListModel) IsFocusRight() bool { return m.focusRight }

// BlurRight returns the model with focus moved back to the list panel.
func (m IssueListModel) BlurRight() IssueListModel {
	m.focusRight = false
	return m
}

func (m *IssueListModel) SetSize(w, h int) {
	if !m.initialized {
		return
	}
	m.width = w
	m.height = h
	leftW, rightW := splitWidths(w)
	m.list.SetSize(leftW, h)
	m.detail.Width = rightW - 2
	m.detail.Height = h
	if issue := m.CurrentIssue(); issue != nil {
		m.detail.SetContent(buildIssueDetail(*issue, rightW-4))
	}
}

func (m IssueListModel) Init() tea.Cmd { return nil }

func (m IssueListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if m.focusRight {
			switch keyMsg.String() {
			case "esc":
				m.focusRight = false
				return m, nil
			}
			var cmd tea.Cmd
			m.detail, cmd = m.detail.Update(msg)
			return m, cmd
		}
		// Left panel: Enter moves focus to detail
		if keyMsg.String() == shared.KeyEnter && len(m.issues) > 0 {
			m.focusRight = true
			return m, nil
		}
	}

	prevIdx := m.list.Index()
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	// Refresh detail panel when cursor moves
	if m.list.Index() != prevIdx {
		if issue := m.CurrentIssue(); issue != nil {
			_, rightW := splitWidths(m.width)
			m.detail.SetContent(buildIssueDetail(*issue, rightW-4))
			m.detail.GotoTop()
		}
	}

	return m, cmd
}

func (m IssueListModel) View() string {
	leftW, rightW := splitWidths(m.width)

	listBorderColor := shared.ColorBorder
	if !m.focusRight {
		listBorderColor = shared.ColorFocus
	}
	leftPane := lipgloss.NewStyle().
		Width(leftW).
		Height(m.height).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(listBorderColor).
		Render(m.list.View())

	detailBorderColor := shared.ColorBorder
	if m.focusRight {
		detailBorderColor = shared.ColorFocus
	}
	_ = detailBorderColor

	var detailContent string
	if m.CurrentIssue() != nil {
		detailContent = m.detail.View()
	} else {
		detailContent = shared.StyleMuted.Render("No issue selected")
	}
	rightPane := lipgloss.NewStyle().
		Width(rightW - 2).
		Height(m.height).
		Padding(0, 1).
		Render(detailContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)
}

func buildIssueDetail(issue jira.Issue, width int) string {
	var sb strings.Builder

	sb.WriteString(shared.StyleIssueKey.Render(issue.Key))
	sb.WriteString("\n")
	sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Summary))
	sb.WriteString("\n\n")

	sb.WriteString(shared.StyleMuted.Render("Status:   "))
	sb.WriteString(shared.StyleIssueStatus.Render(issue.Fields.Status.Name))
	sb.WriteString("\n\n")

	if issue.Fields.Assignee != nil {
		sb.WriteString(shared.StyleMuted.Render("Assignee: "))
		sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Assignee.DisplayName))
		sb.WriteString("\n")
	}
	if issue.Fields.Reporter != nil {
		sb.WriteString(shared.StyleMuted.Render("Reporter: "))
		sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Reporter.DisplayName))
		sb.WriteString("\n")
	}

	lineWidth := width
	if lineWidth > 60 {
		lineWidth = 60
	}
	if lineWidth > 0 {
		sb.WriteString("\n")
		sb.WriteString(shared.StyleMuted.Render(strings.Repeat("─", lineWidth)))
		sb.WriteString("\n\n")
	}

	sb.WriteString(shared.StyleModalTitle.Render("Description"))
	sb.WriteString("\n")
	if issue.Fields.Description != "" {
		sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Description))
	} else {
		sb.WriteString(shared.StyleMuted.Render("No description provided."))
	}

	return sb.String()
}

