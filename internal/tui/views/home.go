package views

import "github.com/svenliebig/lazyjira/internal/tui/shared"

type HomeModel struct{}

func (m HomeModel) View() string {
	return shared.StyleContentArea.Render(
		shared.StyleModalTitle.Render("Welcome to lazyjira") + "\n\n" +
			shared.StyleMuted.Render("Press l to list issues, ? for help, q to quit"),
	)
}
