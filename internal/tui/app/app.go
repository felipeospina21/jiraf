package app

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/jiraf/internal/config"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui"
	"github.com/felipeospina21/jiraf/internal/tui/boards"
	"github.com/felipeospina21/jiraf/internal/tui/details"
	"github.com/felipeospina21/jiraf/internal/tui/icon"
	"github.com/felipeospina21/jiraf/internal/tui/issues"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/popover"
	"github.com/felipeospina21/tuishell/shell"
)

func NewApp() tea.Model {
	cfg := &config.GlobalConfig
	if err := config.Load(cfg); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	theme := tui.BuildTheme(cfg.Theme)

	leftPanelStyle := lipgloss.NewStyle().
		PaddingRight(4).
		MarginBottom(2).
		Foreground(theme.Primary).
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(theme.Border).
		Width(30)

	rightPanelStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), true, false, true, true).
		BorderForeground(theme.Border)

	client := jira.NewClient(cfg)
	b := boards.New(cfg.Filters.Boards, theme)
	left := BoardsPanel{Model: &b}
	main := issues.New(theme)
	det := details.New(theme)

	s := shell.New(shell.Config{
		Theme:           theme,
		LeftPanel:       left,
		MainPanel:       main,
		RightPanel:      DetailsPanel{Model: &det},
		AppIcon:         icon.Jira,
		Keybinds:        tuishell.GlobalKeys(cfg.DemoMode),
		DemoMode:         cfg.DemoMode,
		LeftPanelWidth:  30,
		LeftPanelStyle:  leftPanelStyle,
		RightPanelStyle: rightPanelStyle,
	})

	ti := textarea.New()
	ti.Placeholder = "Write your comment..."
	ti.CharLimit = 0

	return Model{
		Shell:            s,
		Details:          &det,
		Input:            ti,
		client:           client,
		filterPopover:    popover.NewFilter(theme),
		confirmPopover:   popover.NewConfirm(theme),
		transitionPicker: popover.NewList(theme),
	}
}
