package modals

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// CopyModal shows copy options for the current issue.
type CopyModal struct{}

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
		case "esc":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "k":
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
	row := func(key, label string) string {
		return "  " + shared.StyleKeyHint.Render(key) + "  " + shared.StyleNormalItem.Render(label) + "\n"
	}

	content := row("k", "Copy issue key") +
		row("u", "Copy issue URL") +
		row("t", "Copy issue title") +
		row("d", "Copy description") +
		"\n" + shared.StyleMuted.Render("  esc: cancel")

	return Wrap("Copy", content)
}
