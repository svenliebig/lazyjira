package shared

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/svenliebig/lazyjira/internal/theme"
)

var (
	// Exported for use in child packages
	ColorBorder lipgloss.Color
	ColorFocus  lipgloss.Color

	StyleHeader      lipgloss.Style
	StyleStatusBar   lipgloss.Style
	StyleKeyHint     lipgloss.Style
	StyleKeyHintSep  lipgloss.Style
	StyleModal       lipgloss.Style
	StyleModalTitle  lipgloss.Style
	StyleSelectedItem lipgloss.Style
	StyleNormalItem  lipgloss.Style
	StyleMuted       lipgloss.Style
	StyleError       lipgloss.Style
	StyleSuccess     lipgloss.Style
	StyleIssueKey    lipgloss.Style
	StyleIssueStatus lipgloss.Style
	StyleContentArea lipgloss.Style
)

func init() {
	RefreshStyles()
}

// RefreshStyles rebuilds all styles from the current theme. Call this after
// calling theme.SetTheme to apply a new color scheme.
func RefreshStyles() {
	t := theme.Current

	colorPrimary := lipgloss.Color(t.Primary)
	colorSuccess := lipgloss.Color(t.Success)
	colorError := lipgloss.Color(t.Error)
	colorMuted := lipgloss.Color(t.Muted)
	colorBg := lipgloss.Color(t.Bg)
	colorSurface := lipgloss.Color(t.Surface)
	colorText := lipgloss.Color(t.Text)
	colorSubtext := lipgloss.Color(t.Subtext)
	colorBorder := lipgloss.Color(t.Border)
	colorFocus := lipgloss.Color(t.Focus)

	ColorBorder = colorBorder
	ColorFocus = colorFocus

	StyleHeader = lipgloss.NewStyle().
		Background(colorPrimary).
		Foreground(colorText).
		Bold(true).
		Padding(0, 2)

	StyleStatusBar = lipgloss.NewStyle().
		Background(colorSurface).
		Foreground(colorSubtext).
		Padding(0, 1)

	StyleKeyHint = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)

	StyleKeyHintSep = lipgloss.NewStyle().
		Foreground(colorMuted)

	StyleModal = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorFocus).
		Background(colorBg).
		Padding(1, 2)

	StyleModalTitle = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true).
		MarginBottom(1)

	StyleSelectedItem = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)

	StyleNormalItem = lipgloss.NewStyle().
		Foreground(colorText)

	StyleMuted = lipgloss.NewStyle().
		Foreground(colorMuted)

	StyleError = lipgloss.NewStyle().
		Foreground(colorError).
		Bold(true)

	StyleSuccess = lipgloss.NewStyle().
		Foreground(colorSuccess)

	StyleIssueKey = lipgloss.NewStyle().
		Foreground(colorPrimary).
		Bold(true)

	StyleIssueStatus = lipgloss.NewStyle().
		Foreground(colorSuccess).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(colorSuccess)

	StyleContentArea = lipgloss.NewStyle().
		Padding(1, 2)
}
