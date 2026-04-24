// Package details implements the Jira issue details side panel.
package details

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui/icon"
	"github.com/felipeospina21/tuishell"
	"github.com/felipeospina21/tuishell/style"
	"github.com/felipeospina21/tuishell/table"
)

// Model holds the state for the details side panel.
type Model struct {
	Viewport viewport.Model
	Ready    bool
	issue    *jira.Issue
	width    int
	height   int
	// theme-derived styles
	keyStyle     lipgloss.Style
	summaryStyle lipgloss.Style
	labelStyle   lipgloss.Style
	valueStyle   lipgloss.Style
	sectionStyle lipgloss.Style
}

// New creates a new details panel model.
func New(t style.Theme) Model {
	return Model{
		Viewport:     viewport.New(viewport.WithWidth(10), viewport.WithHeight(10)),
		keyStyle:     lipgloss.NewStyle().Foreground(t.Primary).Bold(true).MarginLeft(1),
		summaryStyle: lipgloss.NewStyle().Foreground(t.Text).Bold(true).MarginLeft(1),
		labelStyle:   lipgloss.NewStyle().Foreground(t.TextDimmed).MarginLeft(1).Width(14),
		valueStyle:   lipgloss.NewStyle().Foreground(t.Text),
		sectionStyle: lipgloss.NewStyle().Foreground(t.Primary).Bold(true).MarginLeft(1),
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		match := tuishell.KeyMatcher(msg)
		switch {
		case match(Keybinds.Fullscreen):
			return m, func() tea.Msg { return tuishell.ToggleFullscreenMsg{} }
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerH := lipgloss.Height(m.headerView())
		footerH := lipgloss.Height(m.footerView())
		vpH := msg.Height - headerH - footerH
		if vpH < 1 {
			vpH = 1
		}
		if !m.Ready {
			m.Viewport = viewport.New(viewport.WithWidth(msg.Width), viewport.WithHeight(vpH))
			m.Ready = true
		} else {
			m.Viewport.SetWidth(msg.Width)
			m.Viewport.SetHeight(vpH)
		}
		if m.issue != nil {
			m.renderContent()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.Viewport, cmd = m.Viewport.Update(msg)
	return m, cmd
}

func (m Model) View() tea.View {
	return tea.NewView(fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView()))
}

func (m Model) headerView() string {
	title := titleStyle.Render(" Issue Details ")
	line := strings.Repeat("─", max(0, m.Viewport.Width()-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf(" %3.f%% ", m.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.Viewport.Width()-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

// SetContent formats the issue and sets it as the viewport content.
func (m *Model) SetContent(issue jira.Issue) {
	m.issue = &issue
	m.renderContent()
}

func (m *Model) renderContent() {
	f := m.issue.Fields
	var b strings.Builder

	// Key + Summary
	b.WriteString(m.keyStyle.Render(m.issue.Key))
	b.WriteString("\n")
	b.WriteString(m.summaryStyle.Render(f.Summary))
	b.WriteString("\n\n")

	// Status / Priority / Type
	m.writeField(&b, "Status", f.Status.Name)
	m.writeField(&b, "Priority", f.Priority.Name)
	m.writeField(&b, "Type", f.IssueType.Name)
	b.WriteString("\n")

	// People
	if f.Assignee != nil {
		m.writeField(&b, "Assignee", f.Assignee.DisplayName)
	}
	if f.Reporter != nil {
		m.writeField(&b, "Reporter", f.Reporter.DisplayName)
	}
	b.WriteString("\n")

	// Sprint
	if s := f.ActiveSprint(); s != nil {
		m.writeField(&b, "Sprint", s.Name)
	}

	// Story Points
	if f.StoryPoints != nil {
		m.writeField(&b, "Story Points", strconv.Itoa(int(*f.StoryPoints)))
	}

	// Labels
	if len(f.Labels) > 0 {
		m.writeField(&b, "Labels", strings.Join(f.Labels, ", "))
	}

	// Components
	if len(f.Components) > 0 {
		names := make([]string, len(f.Components))
		for i, c := range f.Components {
			names[i] = c.Name
		}
		m.writeField(&b, "Components", strings.Join(names, ", "))
	}

	// Dates
	m.writeField(&b, "Created", f.Created.Time.Format("2006-01-02 15:04"))
	m.writeField(&b, "Updated", f.Updated.Time.Format("2006-01-02 15:04"))

	// Description
	if f.Description != "" || m.issue.RenderedFields.Description != "" {
		b.WriteString("\n")
		b.WriteString(m.sectionStyle.Render("Description"))
		b.WriteString("\n")
		desc := m.issue.RenderedFields.Description
		if desc == "" {
			desc = f.Description
		}
		md := jira.HTMLToMarkdown(desc)
		b.WriteString(lipgloss.NewStyle().PaddingLeft(1).Render(renderMarkdown(md, m.width-1)))
		b.WriteString("\n")
	}

	// Comments
	if len(f.Comment.Comments) > 0 {
		m.writeHeader(&b, icon.Comment, fmt.Sprintf("Comments (%d)", f.Comment.Total))
		for i, c := range f.Comment.Comments {
			b.WriteString(m.keyStyle.Render(c.Author.DisplayName))
			b.WriteString("  ")
			b.WriteString(m.labelStyle.Render(timeAgo(table.FormatTime(c.Created.Time))))
			b.WriteString("\n")
			b.WriteString(lipgloss.NewStyle().PaddingLeft(1).Render(renderMarkdown(c.Body, m.width-1)))
			b.WriteString("\n")
			if i < len(f.Comment.Comments)-1 {
				b.WriteString("\n")
			}
		}
	} else {
		m.writeHeader(&b, icon.Comment, "Comments")
		b.WriteString(m.labelStyle.Italic(true).Width(0).Render("There are no comments yet in this issue"))
	}

	m.Viewport.SetContent(b.String())
}

func (m *Model) writeHeader(b *strings.Builder, icon, label string) {
	b.WriteString("\n")
	b.WriteString(m.sectionStyle.Render(fmt.Sprintf("%s %s", icon, label)))
	b.WriteString("\n\n")
}

func (m *Model) writeField(b *strings.Builder, label, value string) {
	b.WriteString(m.labelStyle.Render(label+":") + " " + m.valueStyle.Render(value) + "\n")
}

func renderMarkdown(md string, width int) string {
	s := styles.DarkStyleConfig
	zero := uint(0)
	s.Document.Margin = &zero
	r, err := glamour.NewTermRenderer(
		glamour.WithStyles(s),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return md
	}
	out, err := r.Render(md)
	if err != nil {
		return md
	}
	return strings.TrimSpace(out)
}

// Styles that don't depend on theme
var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0)
	}()
)

func timeAgo(time string) string {
	return fmt.Sprintf("%s ago", time)
}
