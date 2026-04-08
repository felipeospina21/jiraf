package issues

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/mrjira/internal/jira"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/loader"
	"github.com/felipeospina21/tuishell/style"
	"github.com/felipeospina21/tuishell/table"
)

var theme = style.DefaultTheme()

const (
	tableBorderX  = 2
	tableBorderY  = 2
	headerLines   = 1
	tableOverhead = 1
)

var titleStyle = lipgloss.NewStyle().
	Margin(0, 0, 0, 1).
	Foreground(theme.Primary).
	Bold(true)

// FetchedMsg carries the result of an issues fetch.
type FetchedMsg struct {
	Issues []jira.Issue
	Err    error
}

var cols = []table.Column{
	{Name: "key", Title: "Key", Width: 8},
	{Name: "type", Title: "Type", Width: 5},
	{Name: "priority", Title: "Pri", Width: 4, Centered: true},
	{Name: "summary", Title: "Summary", Width: 35},
	{Name: "status", Title: "Status", Width: 10},
	{Name: "assignee", Title: "Assignee", Width: 10},
	{Name: "reporter", Title: "Reporter", Width: 0},
	{Name: "labels", Title: "Labels", Width: 0},
	{Name: "components", Title: "Comp", Width: 0},
	{Name: "sprint", Title: "Sprint", Width: 10},
	{Name: "comments", Title: "💬", Width: 3, Centered: true},
	{Name: "subtasks", Title: "📋", Width: 3, Centered: true},
	{Name: "fixver", Title: "Fix", Width: 0},
	{Name: "due", Title: "Due", Width: 6},
	{Name: "created", Title: "Created", Width: 3},
	{Name: "updated", Title: "Updated", Width: 3},
}

// Model is the main-panel issues table.
type Model struct {
	Table         table.Model
	SelectedBoard string
	Loading       bool
	SpinnerView   string
	width         int
	height        int
}

func New() Model {
	return Model{
		Table: table.New(table.WithFocused(true)),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.Table.EmptyMessage = "Select a board"
		tableW := m.tableWidth()
		m.Table.W = tableW
		m.Table.H = m.height - tableBorderY
		return m, nil

	case FetchedMsg:
		if msg.Err != nil {
			return m, func() tea.Msg {
				return tuishell.FinishTaskMsg{Err: msg.Err}
			}
		}
		rows := make([]table.Row, len(msg.Issues))
		for i, issue := range msg.Issues {
			rows[i] = issueToRow(issue)
		}
		tableW := m.tableWidth()
		h := m.height - headerLines - tableBorderY - tableOverhead
		if h < 3 {
			h = 3
		}
		s := table.ThemedStyles(theme)
		m.Table = table.InitModel(table.InitModelParams{
			Rows:   rows,
			Colums: getTableCols(tableW),
			Styles: &s,
			Width:  tableW,
			Height: h,
		})
		m.Table.W = tableW
		m.Table.H = h
		count := len(msg.Issues)
		return m, tea.Batch(
			func() tea.Msg { return tuishell.FinishTaskMsg{Keybinds: tuishell.GlobalKeys(false)} },
			func() tea.Msg { return tuishell.SetStatusMsg{Content: fmt.Sprintf("%d issues", count)} },
		)

	case tea.KeyPressMsg:
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() tea.View {
	header := titleStyle.Render(fmt.Sprintf("%s - Issues", m.SelectedBoard))
	if m.Loading {
		loaderView := lipgloss.NewStyle().
			Width(m.width).
			Height(m.height - lipgloss.Height(header)).
			Render(loader.View(theme, m.SpinnerView))
		return tea.NewView(lipgloss.JoinVertical(0, header, loaderView))
	}
	tbl := table.RenderTable(theme, m.Table.View())
	return tea.NewView(lipgloss.JoinVertical(0, header, tbl))
}

func (m Model) tableWidth() int {
	return m.width - table.DocStyle(theme).GetHorizontalFrameSize() - tableBorderX
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
	return table.Row{
		i.Key,
		f.IssueType.Name,
		f.Priority.Name,
		f.Summary,
		f.Status.Name,
		assignee,
		reporter,
		strings.Join(f.Labels, ", "),
		nameList(f.Components),
		sprint,
		strconv.Itoa(f.Comment.Total),
		strconv.Itoa(len(f.Subtasks)),
		nameList(f.FixVersions),
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
