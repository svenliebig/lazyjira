package modals

import (
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// ExcludeModal lets the user choose how to exclude the current issue.
// cursor positions: 0 = parent row, 1 = key row.
// When parent is unavailable the cursor starts on the key row and cannot move to parent.
type ExcludeModal struct {
	issue  *jira.Issue
	cursor int
}

func NewExcludeModal(issue *jira.Issue) ExcludeModal {
	cursor := 0
	if issue == nil || issue.Fields.Parent == nil {
		cursor = 1 // start on key row when parent is unavailable
	}
	return ExcludeModal{issue: issue, cursor: cursor}
}

func (m ExcludeModal) Init() tea.Cmd { return nil }

func (m ExcludeModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	hasParent := m.issue != nil && m.issue.Fields.Parent != nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "h":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				next := m.cursor - 1
				// skip disabled parent row
				if next == 0 && !hasParent {
					break
				}
				m.cursor = next
			}
			return m, nil
		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			}
			return m, nil
		case "enter", "l":
			return m, m.selectCurrent(hasParent)
		case "i":
			val := m.issue.Key
			return m, func() tea.Msg {
				return shared.ExcludeActionMsg{Type: "key", Value: val}
			}
		case "p":
			if hasParent {
				val := m.issue.Fields.Parent.Key
				return m, func() tea.Msg {
					return shared.ExcludeActionMsg{Type: "parent", Value: val}
				}
			}
		}
	}
	return m, nil
}

func (m ExcludeModal) selectCurrent(hasParent bool) tea.Cmd {
	if m.cursor == 0 && hasParent {
		val := m.issue.Fields.Parent.Key
		return func() tea.Msg { return shared.ExcludeActionMsg{Type: "parent", Value: val} }
	}
	if m.cursor == 1 {
		val := m.issue.Key
		return func() tea.Msg { return shared.ExcludeActionMsg{Type: "key", Value: val} }
	}
	return nil
}

func (m ExcludeModal) View() string {
	hasParent := m.issue != nil && m.issue.Fields.Parent != nil

	// Parent row
	var parentRow string
	if hasParent {
		prefix := "  "
		labelStyle := shared.StyleNormalItem
		if m.cursor == 0 {
			prefix = shared.StyleSelectedItem.Render(">") + " "
			labelStyle = shared.StyleSelectedItem
		}
		label := "Exclude all by parent issue (" + m.issue.Fields.Parent.Key + ")"
		parentRow = prefix + shared.StyleKeyHint.Render("p") + "  " + labelStyle.Render(label) + "\n"
	} else {
		label := lipgloss.NewStyle().
			Foreground(shared.ColorBorder).
			Strikethrough(true).
			Render("Exclude all by parent issue (no parent)")
		parentRow = "  " + shared.StyleMuted.Render("p") + "  " + label + "\n"
	}

	// Key row
	prefix := "  "
	labelStyle := shared.StyleNormalItem
	if m.cursor == 1 {
		prefix = shared.StyleSelectedItem.Render(">") + " "
		labelStyle = shared.StyleSelectedItem
	}
	keyRow := prefix + shared.StyleKeyHint.Render("i") + "  " + labelStyle.Render("Exclude by issue key ("+m.issue.Key+")") + "\n"

	content := parentRow + keyRow + "\n" + shared.StyleMuted.Render("  ↑/k  ↓/j: navigate   enter/l: select   h/esc: cancel")

	return Wrap("Exclude Issue", content)
}
