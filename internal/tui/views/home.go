package views

import "github.com/svenliebig/jira-cli/internal/tui/shared"

type HomeModel struct{}

func (m HomeModel) View() string {
	return shared.StyleContentArea.Render(
		shared.StyleModalTitle.Render("Welcome to jira-cli") + "\n\n" +
			shared.StyleMuted.Render("Press l to list issues, ? for help, q to quit"),
	)
}
