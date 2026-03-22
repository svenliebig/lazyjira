package modals

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// AuthModal collects Jira credentials from the user.
type AuthModal struct {
	inputs  []textinput.Model
	focused int
	err     string
}

func NewAuthModal() AuthModal {
	urlInput := textinput.New()
	urlInput.Placeholder = "https://yourcompany.atlassian.net"
	urlInput.Focus()
	urlInput.CharLimit = 256
	urlInput.Width = 50

	emailInput := textinput.New()
	emailInput.Placeholder = "you@example.com"
	emailInput.CharLimit = 256
	emailInput.Width = 50

	tokenInput := textinput.New()
	tokenInput.Placeholder = "your-api-token"
	tokenInput.CharLimit = 512
	tokenInput.Width = 50
	tokenInput.EchoMode = textinput.EchoPassword
	tokenInput.EchoCharacter = '•'

	return AuthModal{
		inputs:  []textinput.Model{urlInput, emailInput, tokenInput},
		focused: 0,
	}
}

func (m AuthModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m AuthModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "tab", "shift+tab":
			if msg.String() == "tab" {
				m.focused = (m.focused + 1) % len(m.inputs)
			} else {
				m.focused = (m.focused - 1 + len(m.inputs)) % len(m.inputs)
			}
			for i := range m.inputs {
				if i == m.focused {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, nil
		case "enter":
			if m.focused < len(m.inputs)-1 {
				// Move to next field
				m.focused++
				for i := range m.inputs {
					if i == m.focused {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
				return m, nil
			}
			// Last field — submit
			url := m.inputs[0].Value()
			email := m.inputs[1].Value()
			token := m.inputs[2].Value()
			if url == "" || email == "" || token == "" {
				m.err = "All fields are required"
				return m, nil
			}
			return m, func() tea.Msg {
				return shared.AuthCompletedMsg{URL: url, Email: email, Token: token}
			}
		}
	}

	var cmds []tea.Cmd
	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m AuthModal) View() string {
	content := shared.StyleMuted.Render("Jira Cloud URL:") + "\n" +
		m.inputs[0].View() + "\n\n" +
		shared.StyleMuted.Render("Email:") + "\n" +
		m.inputs[1].View() + "\n\n" +
		shared.StyleMuted.Render("API Token:") + "\n" +
		m.inputs[2].View() + "\n\n" +
		shared.StyleMuted.Render("tab: next field  enter: confirm  esc: cancel")

	if m.err != "" {
		content += "\n" + shared.StyleError.Render(m.err)
	}

	return Wrap("Authentication Required", content)
}
