package app

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/modal"
	"github.com/felipeospina21/tuishell/statusline"
	"github.com/felipeospina21/tuishell/style"
	"github.com/felipeospina21/tuishell/table"
)

var theme = style.DefaultTheme()

// Panel styles matching mrglab's look.
var (
	leftPanelStyle = lipgloss.NewStyle().
			PaddingRight(4).
			MarginBottom(2).
			Foreground(theme.Primary).
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(theme.Border).
			Width(30)

	rightPanelStyle = lipgloss.NewStyle().
			MarginTop(1).
			Border(lipgloss.NormalBorder(), true, false, true, true).
			BorderForeground(theme.Border)

	mainFrameStyle = style.MainFrameStyle(theme)
)

const leftPanelWidth = 30

type Model struct {
	projects   list.Model
	table      table.Model
	statusline statusline.Model
	modal      modal.Model
	spinner    spinner.Model
	layout     tuishell.Layout
	ctx        tuishell.AppContext
	isLeftOpen bool
}

func NewApp() tea.Model {
	ctx := tuishell.AppContext{DevMode: true}
	ctx.FocusedPanel = tuishell.LeftPanel

	items := []list.Item{
		listItem{name: "My Board"},
		listItem{name: "Sprint 42"},
		listItem{name: "Backlog"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 30, 10)
	l.Title = "Projects"
	l.SetShowHelp(false)

	sl := statusline.New(theme, ctx.DevMode, tuishell.GlobalKeys(true))
	sl.ProjectLabel = "🎫 mrjira"

	return Model{
		projects:   l,
		table:      table.New(table.WithFocused(true)),
		statusline: sl,
		modal:      modal.New(&ctx, theme),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(statusline.SpinnerStyle(theme)),
		),
		ctx:        ctx,
		isLeftOpen: true,
	}
}

type listItem struct{ name string }

func (i listItem) Title() string       { return i.name }
func (i listItem) Description() string { return "" }
func (i listItem) FilterValue() string { return i.name }

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.statusline.Init(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ctx.Window = msg
		m.recomputeLayout()
		return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) recomputeLayout() {
	cfg := tuishell.LayoutConfig{
		MainFrameStyle:  mainFrameStyle,
		StatusBarStyle:  statusline.StatusBarStyle(),
		LeftPanelStyle:  leftPanelStyle,
		RightPanelStyle: rightPanelStyle,
		LeftPanelWidth:  leftPanelWidth,
		StatuslineLines: 1,
	}
	m.layout = tuishell.ComputeLayout(m.ctx.Window, cfg, m.isLeftOpen, false, false)
	m.projects.SetHeight(m.layout.LeftPanel.Height)
	m.statusline.Width = m.layout.Statusline.Width
	m.table.EmptyMessage = "Select a project"
	m.table.W = m.layout.MainPanel.Width - table.DocStyle(theme).GetHorizontalFrameSize() - 2
	m.table.H = m.layout.MainPanel.Height - 2
}

func (m Model) View() tea.View {
	left := leftPanelStyle.Render(m.projects.View())
	tbl := table.DocStyle(theme).Render(m.table.View())
	body := lipgloss.JoinHorizontal(0, left, tbl)
	sl := m.statusline.View()
	screen := mainFrameStyle.Render(lipgloss.JoinVertical(0, body, sl))

	v := tea.NewView(screen)
	v.AltScreen = true
	return v
}
