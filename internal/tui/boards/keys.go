package boards

import (
	"slices"

	"charm.land/bubbles/v2/key"
	"github.com/felipeospina21/jiraf/internal/tui"
)

type BoardsKeyMap struct {
	SelectBoard key.Binding
	CursorUp    key.Binding
	CursorDown  key.Binding
	Filter      key.Binding
	tui.GlobalKeyMap
}

func (k BoardsKeyMap) ShortHelp() []key.Binding {
	return slices.Concat(
		[]key.Binding{k.SelectBoard, k.CursorUp, k.CursorDown, k.Filter},
		tui.CommonKeys,
	)
}

func (k BoardsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		tui.CommonKeys,
		{k.SelectBoard, k.CursorUp, k.CursorDown, k.Filter},
	}
}

var Keybinds = BoardsKeyMap{
	SelectBoard: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select board"),
	),
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	GlobalKeyMap: tui.GlobalKeys(false),
}
