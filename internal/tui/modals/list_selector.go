package modals

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// ListSelectorModal lets the user choose what list to display.
type ListSelectorModal struct {
	cursor int
}

var listItems = []struct {
	shortcut string
	label    string
	listType string
}{
	{"a", "Assigned Issues", "assigned"},
	{"x", "Excluded Issues", "excluded"},
}

func NewListSelectorModal() ListSelectorModal {
	return ListSelectorModal{}
}

func (m ListSelectorModal) Init() tea.Cmd {
	return nil
}

func (m ListSelectorModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "h":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(listItems)-1 {
				m.cursor++
			}
			return m, nil
		case "enter", "l":
			t := listItems[m.cursor].listType
			return m, func() tea.Msg { return shared.ListSelectedMsg{Type: t} }
		case "a":
			return m, func() tea.Msg {
				return shared.ListSelectedMsg{Type: "assigned"}
			}
		case "x":
			return m, func() tea.Msg {
				return shared.ListSelectedMsg{Type: "excluded"}
			}
		}
	}
	return m, nil
}

func (m ListSelectorModal) View() string {
	content := shared.StyleMuted.Render("Select list type:") + "\n\n"
	for i, item := range listItems {
		prefix := "  "
		labelStyle := shared.StyleNormalItem
		if i == m.cursor {
			prefix = shared.StyleSelectedItem.Render(">") + " "
			labelStyle = shared.StyleSelectedItem
		}
		content += prefix + shared.StyleKeyHint.Render(item.shortcut) + "  " + labelStyle.Render(item.label) + "\n"
	}
	content += "\n" + shared.StyleMuted.Render("  ↑/k  ↓/j: navigate   enter/l: select   h/q/esc: cancel")

	return Wrap("Issue Lists", content)
}
