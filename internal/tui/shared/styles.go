package shared

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary   = lipgloss.Color("#7C3AED")
	colorSecondary = lipgloss.Color("#6B7280") //nolint:unused
	colorSuccess   = lipgloss.Color("#10B981")
	colorError     = lipgloss.Color("#EF4444")
	colorMuted     = lipgloss.Color("#9CA3AF")
	colorBg        = lipgloss.Color("#1F2937")
	colorSurface   = lipgloss.Color("#374151")
	colorText      = lipgloss.Color("#F9FAFB")
	colorSubtext   = lipgloss.Color("#D1D5DB")
	colorBorder    = lipgloss.Color("#4B5563") //nolint:unused
	colorFocus     = lipgloss.Color("#7C3AED")

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
)
