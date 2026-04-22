package boards

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/felipeospina21/jiraf/internal/tui"
)

const ellipsis = "…"

var theme = tui.DefaultTheme()

var TitleStyle = lipgloss.NewStyle().
	Foreground(theme.Info)

type DefaultItemStyles struct {
	NormalTitle    lipgloss.Style
	NormalDesc     lipgloss.Style
	SelectedTitle  lipgloss.Style
	SelectedDesc   lipgloss.Style
	DimmedTitle    lipgloss.Style
	DimmedDesc     lipgloss.Style
	FilterMatch    lipgloss.Style
}

func NewDefaultItemStyles() (s DefaultItemStyles) {
	s.NormalTitle = lipgloss.NewStyle().
		Foreground(theme.PrimaryFg).
		Padding(0, 0, 0, 2)

	s.NormalDesc = lipgloss.NewStyle().
		Foreground(theme.TextDimmed).
		Padding(0, 0, 0, 2)

	s.SelectedTitle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.SelectionBorder).
		Foreground(theme.PrimaryBright).
		Padding(0, 0, 0, 1)

	s.SelectedDesc = s.SelectedTitle.
		Foreground(theme.PrimaryBright)

	s.DimmedTitle = lipgloss.NewStyle().
		Foreground(theme.TextDimmed).
		Padding(0, 0, 0, 2)

	s.DimmedDesc = lipgloss.NewStyle().
		Foreground(theme.Dim)

	s.FilterMatch = lipgloss.NewStyle().Underline(true)

	return s
}

type itemDelegate struct {
	ShowDescription bool
	Styles          DefaultItemStyles
}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	var (
		title, desc  string
		matchedRunes []int
		s            = &d.Styles
	)

	if i, ok := listItem.(item); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	if m.Width() <= 0 {
		return
	}

	textwidth := m.Width() - s.NormalTitle.GetPaddingLeft() - s.NormalTitle.GetPaddingRight()
	title = ansi.Truncate(title, textwidth, ellipsis)
	if d.ShowDescription {
		var lines []string
		for i, line := range strings.Split(desc, "\n") {
			if i >= d.Height()-1 {
				break
			}
			lines = append(lines, ansi.Truncate(line, textwidth, ellipsis))
		}
		desc = strings.Join(lines, "\n")
	}

	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.Render(desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Render(desc)
	} else {
		if isFiltered {
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		fmt.Fprintf(w, "%s\n%s", title, desc) //nolint: errcheck
		return
	}
	fmt.Fprintf(w, "%s", title) //nolint: errcheck
}
