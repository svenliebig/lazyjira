package modals

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// SprintPickerModal shows a filterable list of sprints for selection.
type SprintPickerModal struct {
	sprints  []jira.Sprint
	filter   string
	cursor   int
	filtered []jira.Sprint
}

func NewSprintPickerModal(sprints []jira.Sprint) SprintPickerModal {
	m := SprintPickerModal{sprints: sprints}
	m.filtered = sprints
	return m
}

func (m SprintPickerModal) Init() tea.Cmd { return nil }

func (m SprintPickerModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "backspace":
			if len(m.filter) > 0 {
				m.filter = m.filter[:len(m.filter)-1]
				m.refilter()
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.filtered)-1 {
				m.cursor++
			}
		case "enter":
			if m.cursor < len(m.filtered) {
				s := m.filtered[m.cursor]
				return m, func() tea.Msg { return shared.SprintSelectedMsg{Sprint: s} }
			}
		default:
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.refilter()
			}
		}
	}
	return m, nil
}

func (m *SprintPickerModal) refilter() {
	m.cursor = 0
	if m.filter == "" {
		m.filtered = m.sprints
		return
	}
	q := strings.ToLower(m.filter)
	result := m.filtered[:0:0]
	for _, s := range m.sprints {
		if strings.Contains(strings.ToLower(s.Name), q) {
			result = append(result, s)
		}
	}
	m.filtered = result
}

func (m SprintPickerModal) View() string {
	var sb strings.Builder

	filterLine := "  Filter: " + m.filter + "_"
	sb.WriteString(shared.StyleMuted.Render(filterLine) + "\n\n")

	if len(m.filtered) == 0 {
		sb.WriteString(shared.StyleMuted.Render("  No matching sprints.\n"))
	} else {
		for i, s := range m.filtered {
			prefix := "  "
			nameStyle := shared.StyleNormalItem
			if i == m.cursor {
				prefix = shared.StyleSelectedItem.Render(">") + " "
				nameStyle = shared.StyleSelectedItem
			}
			label := s.Name
			if s.State != "" {
				label += "  " + shared.StyleMuted.Render("("+s.State+")")
			}
			sb.WriteString(prefix + nameStyle.Render(label) + "\n")
		}
	}

	sb.WriteString("\n" + shared.StyleMuted.Render("  type: filter   ↑/k  ↓/j: navigate   enter: select   esc: cancel"))

	return Wrap("Select Sprint", sb.String())
}
