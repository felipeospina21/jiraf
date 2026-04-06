// Package export exposes the mrjira TUI app for embedding in launchers.
package export

import (
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/mrjira/internal/tui/app"
)

// NewApp returns an initialized mrjira tea.Model ready to run.
func NewApp() tea.Model {
	return app.NewApp()
}
