package modals

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/jira"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// ExcludeModal lets the user choose how to exclude the current issue.
type ExcludeModal struct {
	issue *jira.Issue
}

func NewExcludeModal(issue *jira.Issue) ExcludeModal {
	return ExcludeModal{issue: issue}
}

func (m ExcludeModal) Init() tea.Cmd { return nil }

func (m ExcludeModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "k":
			val := m.issue.Key
			return m, func() tea.Msg {
				return shared.ExcludeActionMsg{Type: "key", Value: val}
			}
		case "p":
			if m.issue.Fields.Parent != nil {
				val := m.issue.Fields.Parent.Key
				return m, func() tea.Msg {
					return shared.ExcludeActionMsg{Type: "parent", Value: val}
				}
			}
		}
	}
	return m, nil
}

func (m ExcludeModal) View() string {
	hasParent := m.issue != nil && m.issue.Fields.Parent != nil

	// Parent row — strikethrough and muted when unavailable
	var parentRow string
	if hasParent {
		label := "Exclude all by parent issue (" + m.issue.Fields.Parent.Key + ")"
		parentRow = "  " + shared.StyleKeyHint.Render("p") + "  " + shared.StyleNormalItem.Render(label) + "\n"
	} else {
		label := lipgloss.NewStyle().
			Foreground(shared.ColorBorder).
			Strikethrough(true).
			Render("Exclude all by parent issue (no parent)")
		parentRow = "  " + shared.StyleMuted.Render("p") + "  " + label + "\n"
	}

	keyRow := "  " + shared.StyleKeyHint.Render("k") + "  " +
		shared.StyleNormalItem.Render("Exclude by issue key ("+m.issue.Key+")") + "\n"

	content := parentRow + keyRow + "\n" + shared.StyleMuted.Render("  esc: cancel")

	return Wrap("Exclude Issue", content)
}
