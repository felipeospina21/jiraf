package issues

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui"
	"github.com/felipeospina21/jiraf/internal/tui/icon"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
	"github.com/felipeospina21/tuishell/table"
)

const (
	tableBorderX  = 2
	tableBorderY  = 2
	headerLines   = 1
	tableOverhead = 1
)

// FetchedMsg carries the result of an issues fetch.
type FetchedMsg struct {
	Issues []jira.Issue
	Err    error
}

// ViewDetailsMsg is sent when the user wants to view issue details.
type ViewDetailsMsg struct {
	Issue jira.Issue
}

// TransitionMsg is sent when the user wants to transition an issue's status.
type TransitionMsg struct {
	Issue jira.Issue
}

// OpenInBrowserMsg is sent when the user wants to open the issue in the browser.
type OpenInBrowserMsg struct {
	IssueKey string
}

// RefetchMsg is sent when the user wants to refetch the issues list.
type RefetchMsg struct{}

// OpenFilterMsg is sent when the user wants to open the filter popover.
type OpenFilterMsg struct{}

// ClearFilterMsg is sent when the user wants to clear all active filters.
type ClearFilterMsg struct{}

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
	Issues        []jira.Issue
	SelectedBoard string
	Loading       bool
	SpinnerView   string
	theme         style.Theme
	width         int
	height        int
}

func New(t style.Theme) Model {
	return Model{
		Table: table.New(table.WithFocused(true)),
		theme: t,
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
		m.Issues = msg.Issues
		rows := make([]table.Row, len(msg.Issues))
		for i, issue := range msg.Issues {
			rows[i] = issueToRow(issue)
		}
		tableW := m.tableWidth()
		h := m.height - headerLines - tableBorderY - tableOverhead
		if h < 3 {
			h = 3
		}
		s := table.ThemedStyles(m.theme)
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
		match := tui.KeyMatcher(msg)
		switch {
		case match(Keybinds.Details):
			idx := m.Table.Cursor()
			if idx >= 0 && idx < len(m.Issues) {
				issue := m.Issues[idx]
				return m, func() tea.Msg { return ViewDetailsMsg{Issue: issue} }
			}
		case match(Keybinds.Transition):
			idx := m.Table.Cursor()
			if idx >= 0 && idx < len(m.Issues) {
				issue := m.Issues[idx]
				return m, func() tea.Msg { return TransitionMsg{Issue: issue} }
			}
		case match(Keybinds.OpenInBrowser):
			idx := m.Table.Cursor()
			if idx >= 0 && idx < len(m.Issues) {
				issueKey := m.Issues[idx].Key
				return m, func() tea.Msg { return OpenInBrowserMsg{IssueKey: issueKey} }
			}
		case match(Keybinds.Refetch):
			return m, func() tea.Msg { return RefetchMsg{} }
		case match(Keybinds.Filter):
			return m, func() tea.Msg { return OpenFilterMsg{} }
		case match(Keybinds.ClearFilter):
			return m, func() tea.Msg { return ClearFilterMsg{} }
		}
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() tea.View {
	header := fmt.Sprintf("%s - Issues", m.SelectedBoard)
	return tea.NewView(table.RenderPanel(m.theme, &m.Table, m.Loading, m.SpinnerView, header))
}

func (m Model) tableWidth() int {
	return m.width - table.DocStyle(m.theme).GetHorizontalFrameSize() - tableBorderX
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
