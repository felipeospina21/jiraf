package details

import (
	"slices"

	"charm.land/bubbles/v2/key"
	"github.com/felipeospina21/jiraf/internal/tui"
)

type DetailsKeyMap struct {
	Fullscreen key.Binding
	Comment    key.Binding
	tui.GlobalKeyMap
}

func (k DetailsKeyMap) ShortHelp() []key.Binding {
	return slices.Concat(
		[]key.Binding{k.Fullscreen, k.Comment},
		tui.CommonKeys,
	)
}

func (k DetailsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		tui.CommonKeys,
		{k.Fullscreen, k.Comment},
	}
}

var Keybinds = DetailsKeyMap{
	Fullscreen: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "fullscreen"),
	),
	Comment: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "comment"),
	),
	GlobalKeyMap: tui.GlobalKeys(false),
}
