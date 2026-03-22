package modals

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// HelpModal shows keyboard shortcuts.
type HelpModal struct{}

func NewHelpModal() HelpModal {
	return HelpModal{}
}

func (m HelpModal) Init() tea.Cmd {
	return nil
}

func (m HelpModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "?", "q":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		}
	}
	return m, nil
}

func (m HelpModal) View() string {
	shortcuts := [][]string{
		{"l", "Open issue list selector"},
		{"?", "Toggle help"},
		{"q", "Quit (from home screen)"},
		{"esc", "Go back / close modal"},
		{"o", "Open issue in browser"},
		{"y", "Copy sub-menu (when issue selected)"},
		{"y → k", "Copy issue key"},
		{"y → u", "Copy issue URL"},
		{"y → t", "Copy issue title"},
		{"y → d", "Copy issue description"},
		{"t", "Show transitions for issue"},
		{"x", "Exclude issue (when issue selected)"},
		{"x", "Remove exclusion (in excluded list)"},
		{"a", "AI assistance sub-menu"},
		{"a → s", "Generate AI summary"},
		{"↑/k", "Move up"},
		{"↓/j", "Move down"},
		{"enter", "Select item"},
	}

	var sb strings.Builder
	for _, row := range shortcuts {
		key := shared.StyleKeyHint.Render(row[0])
		desc := shared.StyleNormalItem.Render(row[1])
		sb.WriteString("  " + key + strings.Repeat(" ", 14-len(row[0])) + desc + "\n")
	}

	return Wrap("Keyboard Shortcuts", sb.String())
}
