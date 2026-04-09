package issues

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui/icon"
	"github.com/felipeospina21/tuishell"
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
	{Name: "created", Title: icon.Clock, Width: 3},
	{Name: "priority", Title: "Priority", Width: 4, Centered: true},
	{Name: "key", Title: "Key", Width: 8},
	{Name: "summary", Title: "Summary", Width: 35},
	{Name: "status", Title: "Status", Width: 10},
	{Name: "type", Title: "Type", Width: 5},
	{Name: "assignee", Title: "Assignee", Width: 10},
	{Name: "reporter", Title: "Reporter", Width: 0},
	{Name: "labels", Title: "Labels", Width: 0},
	{Name: "components", Title: "Comp", Width: 0},
	{Name: "sprint", Title: "Sprint", Width: 10},
	{Name: "sprint_start", Title: "Start", Width: 5},
	{Name: "sprint_end", Title: "End", Width: 5},
	{Name: "sp", Title: icon.Weight, Width: 3, Centered: true},
	{Name: "comments", Title: icon.Comment, Width: 3, Centered: true},
	{Name: "subtasks", Title: icon.Subtask, Width: 3, Centered: true},
	{Name: "fixver", Title: "Fix", Width: 0},
	{Name: "updated", Title: icon.UserUpdate, Width: 3},
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
		if len(m.Table.Rows()) > 0 {
			h := m.height - headerLines - tableBorderY - tableOverhead
			if h < 3 {
				h = 3
			}
			m.Table.SetColumns(getTableCols(tableW))
			m.Table.SetWidth(tableW)
			m.Table.SetHeight(h)
		}
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
		return m, func() tea.Msg {
			return tuishell.FinishTaskMsg{Err: nil}
		}

	case tea.KeyPressMsg:
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() tea.View {
	header := fmt.Sprintf("%s - Issues", m.SelectedBoard)
	return tea.NewView(table.RenderPanel(theme, &m.Table, m.Loading, m.SpinnerView, header))
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
	return table.Row{
		table.FormatTime(f.Created.Time),
		f.Priority.Name,
		i.Key,
		f.Summary,
		f.Status.Name,
		f.IssueType.Name,
		assignee,
		reporter,
		strings.Join(f.Labels, ", "),
		nameList(f.Components),
		sprintName(f.ActiveSprint()),
		sprintDate(f.ActiveSprint(), true),
		sprintDate(f.ActiveSprint(), false),
		formatSP(f.StoryPoints),
		strconv.Itoa(f.Comment.Total),
		strconv.Itoa(len(f.Subtasks)),
		nameList(f.FixVersions),
		table.FormatTime(f.Updated.Time),
	}
}

func formatSP(sp *float64) string {
	if sp == nil {
		return "-"
	}
	return strconv.Itoa(int(*sp))
}

func sprintName(s *jira.Sprint) string {
	if s == nil {
		return ""
	}
	name := s.Name
	if before, _, ok := strings.Cut(name, " - "); ok {
		name = strings.TrimSpace(before)
	}
	return strings.TrimPrefix(name, "Sprint ")
}

func sprintDate(s *jira.Sprint, start bool) string {
	if s == nil {
		return ""
	}
	d := s.EndDate
	if start {
		d = s.StartDate
	}
	for _, layout := range []string{
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02T15:04:05.000-0700",
		time.RFC3339,
	} {
		if t, err := time.Parse(layout, d); err == nil {
			return t.Format("Jan 02")
		}
	}
	return ""
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
