package app

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/mrjira/internal/config"
	"github.com/felipeospina21/mrjira/internal/jira"
	"github.com/felipeospina21/mrjira/internal/tui/boards"
	"github.com/felipeospina21/mrjira/internal/tui/issues"
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

// Model wraps shell.Model with mrjira-specific domain logic.
type Model struct {
	Shell  shell.Model
	client *jira.Client
}

func NewApp() tea.Model {
	cfg := &config.GlobalConfig
	if err := config.Load(cfg); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	client := jira.NewClient(cfg)
	left := boards.New(cfg.Filters.Boards)
	main := issues.New()

	s := shell.New(shell.Config{
		Theme:          theme,
		LeftPanel:      left,
		MainPanel:      main,
		StatusLabel:    "🎫 mrjira",
		Keybinds:       tuishell.GlobalKeys(cfg.DevMode),
		DevMode:        cfg.DevMode,
		LeftPanelWidth: 30,
		LeftPanelStyle: leftPanelStyle,
	})

	return Model{Shell: s, client: client}
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
		return m, func() tea.Msg { return tuishell.StartTaskMsg{Cmd: m.fetchIssues(msg.Key)} }

	case issues.FetchedMsg:
		// Clear loading state and route to main panel
		if main, ok := m.Shell.Main.(issues.Model); ok {
			main.Loading = false
			m.Shell.Main = main
		}
		var cmd tea.Cmd
		m.Shell.Main, cmd = m.Shell.Main.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.Shell, cmd = m.Shell.Update(msg)

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

func (m Model) fetchIssues(projectKey string) tea.Cmd {
	return func() tea.Msg {
		iss, err := m.client.GetMyIssues(projectKey)
		return issues.FetchedMsg{Issues: iss, Err: err}
	}
}
