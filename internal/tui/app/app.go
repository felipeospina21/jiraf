package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/mrjira/internal/config"
	"github.com/felipeospina21/mrjira/internal/jira"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/modal"
	"github.com/felipeospina21/tuishell/statusline"
	"github.com/felipeospina21/tuishell/style"
	"github.com/felipeospina21/tuishell/table"
)

var theme = style.DefaultTheme()

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

	titleStyle = lipgloss.NewStyle().
			Margin(0, 0, 0, 1).
			Foreground(theme.Primary).
			Bold(true)
)

const (
	leftPanelWidth = 30
	tableBorderX   = 2
	tableBorderY   = 2
	headerLines    = 1
	tableOverhead  = 1
)

// Messages
type issuesFetchedMsg struct {
	issues []jira.Issue
	err    error
}

// Board list item
type boardItem struct {
	board config.Board
}

func (i boardItem) Title() string       { return i.board.Name }
func (i boardItem) Description() string { return i.board.Key }
func (i boardItem) FilterValue() string { return i.board.Name }

// Table columns
var cols = []table.Column{
	{Name: "key", Title: "Key", Width: 10},
	{Name: "type", Title: "Type", Width: 6},
	{Name: "priority", Title: "Pri", Width: 6, Centered: true},
	{Name: "summary", Title: "Summary", Width: 40},
	{Name: "status", Title: "Status", Width: 10},
	{Name: "assignee", Title: "Assignee", Width: 10},
	{Name: "reporter", Title: "Reporter", Width: 10},
	{Name: "labels", Title: "Labels", Width: 10},
	{Name: "components", Title: "Comp", Width: 8},
	{Name: "sprint", Title: "Sprint", Width: 10},
	{Name: "comments", Title: "💬", Width: 3, Centered: true},
	{Name: "subtasks", Title: "📋", Width: 3, Centered: true},
	{Name: "fixver", Title: "Fix", Width: 6},
	{Name: "due", Title: "Due", Width: 8},
	{Name: "created", Title: "Created", Width: 4},
	{Name: "updated", Title: "Updated", Width: 4},
}

func issueToRow(i jira.Issue) table.Row {
	f := i.Fields
	assignee, reporter := "", ""
	if f.Assignee != nil {
		assignee = f.Assignee.DisplayName
	}
	if f.Reporter != nil {
		reporter = f.Reporter.DisplayName
	}
	sprint := ""
	if f.Sprint != nil {
		sprint = f.Sprint.Name
	}
	labels := strings.Join(f.Labels, ", ")
	comps := nameList(f.Components)
	fixVer := nameList(f.FixVersions)

	return table.Row{
		i.Key,
		f.IssueType.Name,
		f.Priority.Name,
		f.Summary,
		f.Status.Name,
		assignee,
		reporter,
		labels,
		comps,
		sprint,
		strconv.Itoa(f.Comment.Total),
		strconv.Itoa(len(f.Subtasks)),
		fixVer,
		f.DueDate,
		table.FormatTime(f.Created.Time),
		table.FormatTime(f.Updated.Time),
	}
}

func nameList(items []jira.NameField) string {
	names := make([]string, len(items))
	for i, n := range items {
		names[i] = n.Name
	}
	return strings.Join(names, ", ")
}

func getTableCols(width int) []table.Column {
	w := table.ColWidth
	visibleCount := 0
	for _, c := range cols {
		if c.Width > 0 {
			visibleCount++
		}
	}
	contentWidth := width - visibleCount*2
	used := 0
	summaryIdx := -1
	result := make([]table.Column, len(cols))
	for i, c := range cols {
		c.Width = w(contentWidth, c.Width)
		if c.Name == "summary" {
			summaryIdx = i
		}
		used += c.Width
		result[i] = c
	}
	if summaryIdx >= 0 {
		result[summaryIdx].Width += contentWidth - used
	}
	return result
}

// Model
type Model struct {
	boards        list.Model
	table         table.Model
	statusline    statusline.Model
	modal         modal.Model
	spinner       spinner.Model
	client        *jira.Client
	layout        tuishell.Layout
	ctx           tuishell.AppContext
	cfg           *config.Config
	selectedBoard string
	isLeftOpen    bool
}

func NewApp() tea.Model {
	cfg := &config.GlobalConfig
	if err := config.Load(cfg); err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}

	ctx := tuishell.AppContext{DevMode: cfg.DevMode}
	ctx.FocusedPanel = tuishell.LeftPanel

	items := make([]list.Item, len(cfg.Filters.Boards))
	for i, b := range cfg.Filters.Boards {
		items[i] = boardItem{board: b}
	}

	l := list.New(items, list.NewDefaultDelegate(), 30, 10)
	l.Title = "Boards"
	l.SetShowHelp(false)

	sl := statusline.New(theme, ctx.DevMode, tuishell.GlobalKeys(ctx.DevMode))
	sl.ProjectLabel = "🎫 mrjira"

	return Model{
		boards:     l,
		table:      table.New(table.WithFocused(true)),
		statusline: sl,
		modal:      modal.New(&ctx, theme),
		spinner: spinner.New(
			spinner.WithSpinner(spinner.Line),
			spinner.WithStyle(statusline.SpinnerStyle(theme)),
		),
		client:     jira.NewClient(cfg),
		ctx:        ctx,
		cfg:        cfg,
		isLeftOpen: true,
	}
}

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
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.ctx.FocusedPanel == tuishell.LeftPanel {
				if item, ok := m.boards.SelectedItem().(boardItem); ok {
					m.selectedBoard = item.board.Key
					m.statusline.Status = statusline.ModesEnum.Loading
					return m, m.fetchIssues(item.board.Key)
				}
			}
		}
		var cmd tea.Cmd
		if m.ctx.FocusedPanel == tuishell.LeftPanel {
			m.boards, cmd = m.boards.Update(msg)
		} else {
			m.table, cmd = m.table.Update(msg)
		}
		return m, cmd

	case issuesFetchedMsg:
		if msg.err != nil {
			m.statusline.Status = statusline.ModesEnum.Error
			m.statusline.Content = msg.err.Error()
			m.modal.Header = "Error"
			m.modal.IsError = true
			m.modal.Content = msg.err.Error()
			m.modal.SetFocus()
			return m, nil
		}
		rows := make([]table.Row, len(msg.issues))
		for i, issue := range msg.issues {
			rows[i] = issueToRow(issue)
		}
		m.recomputeLayout()
		tableW := m.layout.MainPanel.Width - table.DocStyle(theme).GetHorizontalFrameSize() - tableBorderX
		h := m.layout.ContentH - headerLines - tableBorderY - tableOverhead
		if h < 3 {
			h = 3
		}
		s := table.ThemedStyles(theme)
		m.table = table.InitModel(table.InitModelParams{
			Rows:   rows,
			Colums: getTableCols(tableW),
			Styles: &s,
			Width:  tableW,
			Height: h,
		})
		m.table.W = tableW
		m.table.H = h
		m.ctx.FocusedPanel = tuishell.MainPanel
		m.statusline.Status = statusline.ModesEnum.Normal
		m.statusline.Content = fmt.Sprintf("%d issues", len(msg.issues))
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		if m.statusline.Status == statusline.ModesEnum.Loading {
			m.statusline.Content = m.spinner.View()
		}
		return m, cmd
	}
	return m, nil
}

func (m Model) fetchIssues(projectKey string) tea.Cmd {
	return func() tea.Msg {
		issues, err := m.client.GetMyIssues(projectKey)
		return issuesFetchedMsg{issues: issues, err: err}
	}
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
	m.boards.SetHeight(m.layout.LeftPanel.Height)
	m.statusline.Width = m.layout.Statusline.Width
	m.table.EmptyMessage = "Select a board"
	m.table.W = m.layout.MainPanel.Width - table.DocStyle(theme).GetHorizontalFrameSize() - tableBorderX
	m.table.H = m.layout.MainPanel.Height - tableBorderY
}

func (m Model) View() tea.View {
	left := leftPanelStyle.Render(m.boards.View())

	header := titleStyle.Render(fmt.Sprintf("%s - Issues", m.selectedBoard))
	tbl := table.DocStyle(theme).Render(m.table.View())
	main := lipgloss.JoinVertical(0, header, tbl)

	body := lipgloss.JoinHorizontal(0, left, main)
	sl := m.statusline.View()
	screen := mainFrameStyle.Render(lipgloss.JoinVertical(0, body, sl))

	if m.modal.IsError {
		screen = m.modal.View(screen)
	}

	v := tea.NewView(screen)
	v.AltScreen = true
	return v
}
