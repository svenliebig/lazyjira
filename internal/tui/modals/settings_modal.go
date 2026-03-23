package modals

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/theme"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// SettingsModal lets the user browse and select app settings (theme, etc.).
type SettingsModal struct {
	themes      []theme.Theme
	cursor      int
	activeTheme string
}

func NewSettingsModal(themes []theme.Theme, activeTheme string) SettingsModal {
	// Set initial cursor to the active theme.
	cursor := 0
	for i, t := range themes {
		if t.Name == activeTheme {
			cursor = i
			break
		}
	}
	return SettingsModal{
		themes:      themes,
		cursor:      cursor,
		activeTheme: activeTheme,
	}
}

func (m SettingsModal) Init() tea.Cmd { return nil }

func (m SettingsModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "h":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
			}
		case "enter", "l":
			if len(m.themes) > 0 {
				name := m.themes[m.cursor].Name
				return m, func() tea.Msg { return shared.ThemeSelectedMsg{Name: name} }
			}
		}
	}
	return m, nil
}

func (m SettingsModal) View() string {
	content := shared.StyleMuted.Render("Theme") + "\n\n"

	for i, t := range m.themes {
		active := t.Name == m.activeTheme
		prefix := "  "
		labelStyle := shared.StyleNormalItem
		if i == m.cursor {
			prefix = shared.StyleSelectedItem.Render(">") + " "
			labelStyle = shared.StyleSelectedItem
		}
		label := t.Name
		if active {
			label += "  " + shared.StyleSuccess.Render("(active)")
		}
		content += prefix + labelStyle.Render(label) + "\n"
	}

	content += "\n" + shared.StyleMuted.Render("  ↑/k  ↓/j: navigate   enter/l: apply   h/q/esc: cancel")

	return Wrap("Settings", content)
}
