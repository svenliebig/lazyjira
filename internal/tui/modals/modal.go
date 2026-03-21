package modals

import (
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

// Wrap wraps content in a styled modal box.
func Wrap(title, content string) string {
	body := shared.StyleModalTitle.Render(title) + "\n" + content
	return shared.StyleModal.Render(body)
}
