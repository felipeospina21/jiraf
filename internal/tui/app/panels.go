package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/mrjira/internal/tui/boards"
)

// BoardsPanel wraps boards.Model to implement tea.Model and SelectionProvider.
type BoardsPanel struct {
	*boards.Model
}

func (p BoardsPanel) Init() tea.Cmd { return nil }

func (p BoardsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := p.Model.Update(msg)
	if updated, ok := m.(boards.Model); ok {
		p.Model = &updated
	}
	return p, cmd
}

func (p BoardsPanel) View() tea.View { return p.Model.View() }

// SelectedLabel implements tuishell.SelectionProvider.
func (p BoardsPanel) SelectedLabel() string { return p.Model.SelectedLabel() }
