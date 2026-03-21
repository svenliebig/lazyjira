package modals

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// ListSelectorModal lets the user choose what list to display.
type ListSelectorModal struct{}

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
		case "esc", "q":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "a":
			return m, func() tea.Msg {
				return shared.ListSelectedMsg{Type: "assigned"}
			}
		}
	}
	return m, nil
}

func (m ListSelectorModal) View() string {
	content := shared.StyleMuted.Render("Select list type:") + "\n\n" +
		"  " + shared.StyleKeyHint.Render("a") + "  " + shared.StyleNormalItem.Render("Assigned Issues") + "\n\n" +
		shared.StyleMuted.Render("esc: cancel")

	return Wrap("Issue Lists", content)
}
