package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/tui/boards"
	"github.com/felipeospina21/jiraf/internal/tui/details"
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

// DetailsPanel wraps details.Model to implement tea.Model.
type DetailsPanel struct {
	*details.Model
}

func (p DetailsPanel) Init() tea.Cmd { return nil }

func (p DetailsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, cmd := p.Model.Update(msg)
	p.Model = &m
	return p, cmd
}

func (p DetailsPanel) View() tea.View { return p.Model.View() }
