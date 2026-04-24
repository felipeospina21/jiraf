package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui/boards"
	"github.com/felipeospina21/jiraf/internal/tui/details"
	"github.com/felipeospina21/jiraf/internal/tui/issues"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/popover"
	"github.com/felipeospina21/tuishell/shell"
)

// Model wraps shell.Model with jiraf-specific domain logic.
type Model struct {
	Shell   shell.Model
	Details *details.Model
	client  *jira.Client

	// Filter state
	filterPopover    popover.FilterModel
	activeFilters    jira.IssueFilters
	knownStatuses    map[string]bool
	knownPriorities  map[string]bool
	knownTypes       map[string]bool

	// Transition state
	pendingTransition     string // issue key awaiting transition
	pendingTransitionName string // selected transition name awaiting confirmation
	confirmPopover        popover.ConfirmModel
	transitionPicker      popover.ListModel
	transitionIDByName    map[string]string
}

func (m Model) Init() tea.Cmd {
	return m.Shell.Init()
}

func (m *Model) syncKeybinds() {
	switch m.Shell.Ctx.FocusedPanel {
	case tuishell.LeftPanel:
		m.Shell.Statusline.Keybinds = boards.Keybinds
	case tuishell.RightPanel:
		m.Shell.Statusline.Keybinds = details.Keybinds
	case tuishell.MainPanel:
		m.Shell.Statusline.Keybinds = issues.Keybinds
	}
}
