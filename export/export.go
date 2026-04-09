// Package export exposes the jiraf TUI app for embedding in launchers.
package export

import (
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/tui/app"
)

// NewApp returns an initialized jiraf tea.Model ready to run.
func NewApp() tea.Model {
	return app.NewApp()
}
