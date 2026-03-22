package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// IssueDetailModel shows the detail of a single issue.
type IssueDetailModel struct {
	issue       jira.Issue
	viewport    viewport.Model
	width       int
	height      int
	initialized bool
}

func NewIssueDetailModel(issue jira.Issue, width, height int) IssueDetailModel {
	vp := viewport.New(width, height-4) // reserve space for header
	content := buildDetailContent(issue, width)
	vp.SetContent(content)

	return IssueDetailModel{
		issue:       issue,
		viewport:    vp,
		width:       width,
		height:      height,
		initialized: true,
	}
}

func buildDetailContent(issue jira.Issue, width int) string {
	var sb strings.Builder

	// Header section
	sb.WriteString(shared.StyleIssueKey.Render(issue.Key))
	sb.WriteString("\n")
	sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Summary))
	sb.WriteString("\n\n")

	// Status
	sb.WriteString(shared.StyleMuted.Render("Status: "))
	sb.WriteString(shared.StyleIssueStatus.Render(issue.Fields.Status.Name))
	sb.WriteString("\n\n")

	// Assignee
	if issue.Fields.Assignee != nil {
		sb.WriteString(shared.StyleMuted.Render("Assignee: "))
		sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Assignee.DisplayName))
		sb.WriteString("\n")
	}

	// Reporter
	if issue.Fields.Reporter != nil {
		sb.WriteString(shared.StyleMuted.Render("Reporter: "))
		sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Reporter.DisplayName))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(shared.StyleMuted.Render(strings.Repeat("─", min(width-4, 60))))
	sb.WriteString("\n\n")

	// Description
	sb.WriteString(shared.StyleModalTitle.Render("Description"))
	sb.WriteString("\n")
	if issue.Fields.Description != "" {
		sb.WriteString(shared.StyleNormalItem.Render(issue.Fields.Description))
	} else {
		sb.WriteString(shared.StyleMuted.Render("No description provided."))
	}
	sb.WriteString("\n")

	return sb.String()
}

func (m IssueDetailModel) Init() tea.Cmd {
	return nil
}

func (m *IssueDetailModel) SetSize(w, h int) {
	if !m.initialized {
		return
	}
	m.width = w
	m.height = h
	m.viewport.Width = w
	m.viewport.Height = h - 4
	m.viewport.SetContent(buildDetailContent(m.issue, w))
}

func (m IssueDetailModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m IssueDetailModel) View() string {
	scrollPct := fmt.Sprintf("  %3.f%%", m.viewport.ScrollPercent()*100)
	footer := shared.StyleMuted.Render(scrollPct)
	return m.viewport.View() + "\n" + footer
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
