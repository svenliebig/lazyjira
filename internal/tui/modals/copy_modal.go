package modals

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// CopyModal shows copy options for the current issue.
type CopyModal struct {
	cursor int
}

var copyItems = []struct {
	shortcut string
	label    string
	action   string
}{
	{"i", "Copy issue key", "key"},
	{"u", "Copy issue URL", "url"},
	{"t", "Copy issue title", "title"},
	{"d", "Copy description", "desc"},
}

func NewCopyModal() CopyModal {
	return CopyModal{}
}

func (m CopyModal) Init() tea.Cmd {
	return nil
}

func (m CopyModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "h":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(copyItems)-1 {
				m.cursor++
			}
			return m, nil
		case "enter", "l":
			action := copyItems[m.cursor].action
			return m, func() tea.Msg { return shared.CopyActionMsg{Action: action} }
		case "i":
			return m, func() tea.Msg { return shared.CopyActionMsg{Action: "key"} }
		case "u":
			return m, func() tea.Msg { return shared.CopyActionMsg{Action: "url"} }
		case "t":
			return m, func() tea.Msg { return shared.CopyActionMsg{Action: "title"} }
		case "d":
			return m, func() tea.Msg { return shared.CopyActionMsg{Action: "desc"} }
		}
	}
	return m, nil
}

func (m CopyModal) View() string {
	var content string
	for i, item := range copyItems {
		prefix := "  "
		labelStyle := shared.StyleNormalItem
		if i == m.cursor {
			prefix = shared.StyleSelectedItem.Render(">") + " "
			labelStyle = shared.StyleSelectedItem
		}
		keyPart := shared.StyleKeyHint.Render(item.shortcut) + "  "
		content += prefix + keyPart + labelStyle.Render(item.label) + "\n"
	}
	content += "\n" + shared.StyleMuted.Render("  ↑/k  ↓/j: navigate   enter/l: select   h/esc: cancel")

	return Wrap("Copy", content)
}
