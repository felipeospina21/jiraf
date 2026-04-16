package app

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/jiraf/internal/config"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui/boards"
	"github.com/felipeospina21/jiraf/internal/tui/details"
	"github.com/felipeospina21/jiraf/internal/tui/icon"
	"github.com/felipeospina21/jiraf/internal/tui/issues"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/shell"
	"github.com/felipeospina21/tuishell/style"
)

var theme = style.DefaultTheme()

var leftPanelStyle = lipgloss.NewStyle().
	PaddingRight(4).
	MarginBottom(2).
	Foreground(theme.Primary).
	Border(lipgloss.NormalBorder(), false, true, false, false).
	BorderForeground(theme.Border).
	Width(30)

var rightPanelStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), true, false, true, true).
	BorderForeground(theme.Border)

// Model wraps shell.Model with jiraf-specific domain logic.
type Model struct {
	Shell   shell.Model
	Details *details.Model
	client  *jira.Client
}

func NewApp() tea.Model {
	cfg := &config.GlobalConfig
	if err := config.Load(cfg); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	client := jira.NewClient(cfg)
	b := boards.New(cfg.Filters.Boards)
	left := BoardsPanel{Model: &b}
	main := issues.New()
	det := details.New()

	s := shell.New(shell.Config{
		Theme:           theme,
		LeftPanel:       left,
		MainPanel:       main,
		RightPanel:      DetailsPanel{Model: &det},
		AppIcon:         icon.Jira,
		Keybinds:        tuishell.GlobalKeys(cfg.DevMode),
		DevMode:         cfg.DevMode,
		LeftPanelWidth:  30,
		LeftPanelStyle:  leftPanelStyle,
		RightPanelStyle: rightPanelStyle,
	})

	return Model{Shell: s, Details: &det, client: client}
}

func (m Model) Init() tea.Cmd {
	return m.Shell.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case boards.SelectBoardMsg:
		// Update the board name and set loading state on the issues panel
		if main, ok := m.Shell.Main.(issues.Model); ok {
			main.SelectedBoard = msg.Key
			main.Loading = true
			main.SpinnerView = m.Shell.Spinner.View()
			m.Shell.Main = main
		}
		return m, tea.Batch(
			func() tea.Msg { return tuishell.CloseLeftPanelMsg{} },
			func() tea.Msg { return tuishell.StartTaskMsg{Cmd: m.fetchIssues(msg.Key)} },
		)

	case issues.ViewDetailsMsg:
		m.Details.SetContent(msg.Issue)
		var cmds []tea.Cmd
		if !m.Shell.IsRightOpen() {
			cmds = append(cmds, func() tea.Msg { return tuishell.OpenRightPanelMsg{} })
		}
		return m, tea.Batch(cmds...)

	case details.ClosePanelMsg:
		return m, func() tea.Msg { return tuishell.CloseRightPanelMsg{} }

	case tuishell.FinishTaskMsg:
		if main, ok := m.Shell.Main.(issues.Model); ok {
			main.Loading = false
			m.Shell.Main = main
		}
	}

	var cmd tea.Cmd
	m.Shell, cmd = m.Shell.Update(msg)
	m.syncKeybinds()

	// Sync panel pointer after shell update
	if d, ok := m.Shell.Right.(DetailsPanel); ok {
		m.Details = d.Model
	}

	// Update spinner view on issues panel when loading
	if main, ok := m.Shell.Main.(issues.Model); ok && main.Loading {
		main.SpinnerView = m.Shell.Spinner.View()
		m.Shell.Main = main
	}

	return m, cmd
}

func (m Model) View() tea.View {
	return m.Shell.RenderView()
}

func (m *Model) syncKeybinds() {
	switch m.Shell.Ctx.FocusedPanel {
	case tuishell.LeftPanel:
		m.Shell.Statusline.Keybinds = boards.Keybinds
	default:
		m.Shell.Statusline.Keybinds = tuishell.GlobalKeys(m.Shell.Ctx.DevMode)
	}
}

func (m Model) fetchIssues(projectKey string) tea.Cmd {
	return func() tea.Msg {
		iss, err := m.client.GetMyIssues(projectKey)
		return issues.FetchedMsg{Issues: iss, Err: err}
	}
}
