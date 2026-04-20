package app

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/jiraf/internal/config"
	jirafExec "github.com/felipeospina21/jiraf/internal/exec"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui/boards"
	"github.com/felipeospina21/jiraf/internal/tui/details"
	"github.com/felipeospina21/jiraf/internal/tui/icon"
	"github.com/felipeospina21/jiraf/internal/tui/issues"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/popover"
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

	// Transition state
	pendingTransition      string // issue key awaiting transition
	pendingTransitionName  string // selected transition name awaiting confirmation
	confirmPopover         popover.ConfirmModel
	transitionPicker       popover.ListModel
	transitionIDByName     map[string]string
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

	return Model{
		Shell:            s,
		Details:          &det,
		client:           client,
		confirmPopover:   popover.NewConfirm(theme),
		transitionPicker: popover.NewList(theme),
	}
}

func (m Model) Init() tea.Cmd {
	return m.Shell.Init()
}

// TransitionsFetchedMsg carries the result of fetching available transitions.
type TransitionsFetchedMsg struct {
	Transitions []jira.Transition
	Err         error
}

// TransitionDoneMsg carries the result of executing a transition.
type TransitionDoneMsg struct {
	Err error
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

	case issues.OpenInBrowserMsg:
		url := config.GlobalConfig.BaseURL + "/browse/" + msg.IssueKey
		jirafExec.OpenBrowser(url)
		return m, nil

	case issues.RefetchMsg:
		if main, ok := m.Shell.Main.(issues.Model); ok && main.SelectedBoard != "" {
			main.Loading = true
			main.SpinnerView = m.Shell.Spinner.View()
			m.Shell.Main = main
			return m, func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: m.fetchIssues(main.SelectedBoard)}
			}
		}

	case issues.TransitionMsg:
		m.pendingTransition = msg.Issue.Key
		return m, func() tea.Msg {
			return tuishell.StartTaskMsg{Cmd: m.fetchTransitions(msg.Issue.Key)}
		}

	case TransitionsFetchedMsg:
		if msg.Err != nil {
			m.pendingTransition = ""
			return m, func() tea.Msg {
				return tuishell.FinishTaskMsg{Err: msg.Err}
			}
		}
		items := make([]tuishell.ListPopoverItem, len(msg.Transitions))
		m.transitionIDByName = make(map[string]string, len(msg.Transitions))
		for i, t := range msg.Transitions {
			items[i] = tuishell.ListPopoverItem{Label: t.Name, Value: t.Name}
			m.transitionIDByName[t.Name] = t.ID
		}
		m.transitionPicker.Open(fmt.Sprintf("Transition %s", m.pendingTransition), items)
		return m, func() tea.Msg { return tuishell.FinishTaskMsg{Err: nil} }

	case tuishell.SelectListPopoverMsg:
		if m.pendingTransition != "" {
			m.pendingTransitionName = msg.Value
			m.transitionPicker.Close()
			m.confirmPopover.Open(
				"Transition Issue",
				fmt.Sprintf("Transition %s to %s?", m.pendingTransition, msg.Value),
				"Confirm",
				"Cancel",
			)
		}

	case tuishell.ConfirmPopoverYesMsg:
		m.confirmPopover.Close()
		issueKey := m.pendingTransition
		transitionID := m.transitionIDByName[m.pendingTransitionName]
		m.pendingTransition = ""
		m.pendingTransitionName = ""
		m.transitionIDByName = nil
		return m, func() tea.Msg {
			return tuishell.StartTaskMsg{Cmd: m.doTransition(issueKey, transitionID)}
		}

	case tuishell.ConfirmPopoverNoMsg:
		m.confirmPopover.Close()
		m.pendingTransition = ""
		m.pendingTransitionName = ""
		m.transitionIDByName = nil

	case tuishell.CloseListPopoverMsg:
		m.pendingTransition = ""
		m.transitionIDByName = nil
		m.transitionPicker.Close()

	case TransitionDoneMsg:
		if msg.Err != nil {
			return m, func() tea.Msg {
				return tuishell.FinishTaskMsg{Err: msg.Err}
			}
		}
		var cmds []tea.Cmd
		cmds = append(cmds, func() tea.Msg {
			return tuishell.FinishTaskMsg{Err: nil, Keybinds: issues.Keybinds}
		})
		cmds = append(cmds, func() tea.Msg {
			return tuishell.SetStatusMsg{Content: "✓ Transition complete"}
		})
		// Refetch issues if we have a selected board
		if main, ok := m.Shell.Main.(issues.Model); ok && main.SelectedBoard != "" {
			main.Loading = true
			m.Shell.Main = main
			cmds = append(cmds, func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: m.fetchIssues(main.SelectedBoard)}
			})
		}
		return m, tea.Batch(cmds...)

	case tuishell.FinishTaskMsg:
		if main, ok := m.Shell.Main.(issues.Model); ok {
			main.Loading = false
			m.Shell.Main = main
		}
	}

	// Route keys to confirm popover when open
	if m.confirmPopover.IsOpen() {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			var cmd tea.Cmd
			m.confirmPopover, cmd = m.confirmPopover.Update(msg)
			return m, cmd
		}
	}

	// Route keys to transition picker when open
	if m.transitionPicker.IsOpen() {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			var cmd tea.Cmd
			m.transitionPicker, cmd = m.transitionPicker.Update(msg)
			return m, cmd
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
	v := m.Shell.RenderView()
	w := m.Shell.Ctx.Window.Width
	h := m.Shell.Ctx.Window.Height
	if m.confirmPopover.IsOpen() {
		screen := m.confirmPopover.View(v.Content, w, h)
		return tea.View{Content: screen, AltScreen: v.AltScreen}
	}
	if m.transitionPicker.IsOpen() {
		screen := m.transitionPicker.View(v.Content, w, h)
		return tea.View{Content: screen, AltScreen: v.AltScreen}
	}
	return v
}

func (m *Model) syncKeybinds() {
	switch m.Shell.Ctx.FocusedPanel {
	case tuishell.LeftPanel:
		m.Shell.Statusline.Keybinds = boards.Keybinds
	case tuishell.MainPanel:
		m.Shell.Statusline.Keybinds = issues.Keybinds
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

func (m Model) fetchTransitions(issueKey string) tea.Cmd {
	return func() tea.Msg {
		t, err := m.client.GetTransitions(issueKey)
		return TransitionsFetchedMsg{Transitions: t, Err: err}
	}
}

func (m Model) doTransition(issueKey, transitionID string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.DoTransition(issueKey, transitionID)
		return TransitionDoneMsg{Err: err}
	}
}
