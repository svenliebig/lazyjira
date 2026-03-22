package modals

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// TransitionModal shows available issue transitions.
type TransitionModal struct {
	transitions []jira.Transition
}

func NewTransitionModal(transitions []jira.Transition) TransitionModal {
	return TransitionModal{transitions: transitions}
}

func (m TransitionModal) Init() tea.Cmd {
	return nil
}

func (m TransitionModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "esc", "q":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		}
		// Number keys 1-9 select a transition
		if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
			idx := int(key[0]-'1')
			if idx < len(m.transitions) {
				id := m.transitions[idx].ID
				return m, func() tea.Msg {
					return shared.TransitionSelectedMsg{ID: id}
				}
			}
		}
	}
	return m, nil
}

func (m TransitionModal) View() string {
	if len(m.transitions) == 0 {
		return Wrap("Transitions", shared.StyleMuted.Render("No transitions available.\n\nesc: close"))
	}

	var sb strings.Builder
	for i, t := range m.transitions {
		if i >= 9 {
			break
		}
		num := fmt.Sprintf("%d", i+1)
		sb.WriteString("  " + shared.StyleKeyHint.Render(num) + "  " +
			shared.StyleNormalItem.Render(t.Name) + " → " +
			shared.StyleMuted.Render(t.To.Name) + "\n")
	}
	sb.WriteString("\n" + shared.StyleMuted.Render("  esc: cancel"))

	return Wrap("Select Transition", sb.String())
}
