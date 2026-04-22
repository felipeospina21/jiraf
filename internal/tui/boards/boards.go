package boards

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/config"
	"github.com/felipeospina21/tuishell"
)

// SelectBoardMsg is sent when the user picks a board.
type SelectBoardMsg struct {
	Key string
}

type item struct {
	board config.Board
}

func (i item) Title() string       { return i.board.Name }
func (i item) Description() string { return i.board.Key }
func (i item) FilterValue() string { return i.board.Name }

// Model is the left-panel boards list.
type Model struct {
	List list.Model
}

// New creates a boards panel from config.
func New(boards []config.Board) Model {
	items := make([]list.Item, len(boards))
	for i, b := range boards {
		items[i] = item{board: b}
	}
	l := list.New(items, itemDelegate{ShowDescription: true, Styles: NewDefaultItemStyles()}, 30, 10)
	tuishell.ConfigureList(&l)
	l.Title = "Boards"
	l.Styles.Title = TitleStyle
	l.SetShowHelp(false)
	return Model{List: l}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.SetWidth(msg.Width)
		m.List.SetHeight(msg.Height)
		return m, nil
	case tea.KeyPressMsg:
		if msg.String() == "enter" {
			if it, ok := m.List.SelectedItem().(item); ok {
				return m, func() tea.Msg { return SelectBoardMsg{Key: it.board.Key} }
			}
		}
		var cmd tea.Cmd
		m.List, cmd = m.List.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m Model) View() tea.View {
	return tea.NewView(m.List.View())
}

// SelectedLabel implements tuishell.SelectionProvider.
func (m Model) SelectedLabel() string {
	if it, ok := m.List.SelectedItem().(item); ok {
		return it.board.Key
	}
	return ""
}
