package modals

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// AssignModal shows a fuzzy-searchable list of assignable users.
type AssignModal struct {
	users    []jira.User
	filter   string
	cursor   int
	filtered []jira.User
}

func NewAssignModal(users []jira.User) AssignModal {
	m := AssignModal{users: users}
	m.filtered = users
	return m
}

func (m AssignModal) Init() tea.Cmd { return nil }

func (m AssignModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				u := m.filtered[m.cursor]
				return m, func() tea.Msg { return shared.UserSelectedMsg{User: u} }
			}
		default:
			// Append printable single characters to the filter
			if len(msg.String()) == 1 {
				m.filter += msg.String()
				m.refilter()
			}
		}
	}
	return m, nil
}

func (m *AssignModal) refilter() {
	m.cursor = 0
	if m.filter == "" {
		m.filtered = m.users
		return
	}
	q := strings.ToLower(m.filter)
	result := m.filtered[:0:0]
	for _, u := range m.users {
		if strings.Contains(strings.ToLower(u.DisplayName), q) ||
			strings.Contains(strings.ToLower(u.EmailAddress), q) {
			result = append(result, u)
		}
	}
	m.filtered = result
}

func (m AssignModal) View() string {
	var sb strings.Builder

	filterLine := "  Filter: " + m.filter + "_"
	sb.WriteString(shared.StyleMuted.Render(filterLine) + "\n\n")

	if len(m.filtered) == 0 {
		sb.WriteString(shared.StyleMuted.Render("  No matching users.\n"))
	} else {
		for i, u := range m.filtered {
			prefix := "  "
			nameStyle := shared.StyleNormalItem
			if i == m.cursor {
				prefix = shared.StyleSelectedItem.Render(">") + " "
				nameStyle = shared.StyleSelectedItem
			}
			label := u.DisplayName
			if u.EmailAddress != "" {
				label += "  " + shared.StyleMuted.Render("("+u.EmailAddress+")")
			}
			sb.WriteString(prefix + nameStyle.Render(label) + "\n")
		}
	}

	sb.WriteString("\n" + shared.StyleMuted.Render("  type: filter   ↑/k  ↓/j: navigate   enter: assign   esc: cancel"))

	return Wrap("Assign Issue", sb.String())
}
