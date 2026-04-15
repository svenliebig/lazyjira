package modals

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// ChangeModal presents a keyed list of issue property actions.
type ChangeModal struct{}

func NewChangeModal() ChangeModal { return ChangeModal{} }

func (m ChangeModal) Init() tea.Cmd { return nil }

func (m ChangeModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "s", "1":
			return m, func() tea.Msg { return shared.ChangeActionSelectedMsg{Action: "sprint"} }
		case "a", "2":
			return m, func() tea.Msg { return shared.ChangeActionSelectedMsg{Action: "assign"} }
		}
	}
	return m, nil
}

func (m ChangeModal) View() string {
	var sb strings.Builder
	sb.WriteString("  " + shared.StyleKeyHint.Render("1") + "  " + shared.StyleNormalItem.Render("[s] Sprint") + "\n")
	sb.WriteString("  " + shared.StyleKeyHint.Render("2") + "  " + shared.StyleNormalItem.Render("[a] Assign") + "\n")
	sb.WriteString("\n" + shared.StyleMuted.Render("  key/number: select   esc: cancel"))
	return Wrap("Change Issue", sb.String())
}
