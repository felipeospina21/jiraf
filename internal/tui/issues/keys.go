package issues

import (
	"slices"

	"charm.land/bubbles/v2/key"
	"github.com/felipeospina21/jiraf/internal/tui"
)

type IssuesKeyMap struct {
	Details    key.Binding
	Transition key.Binding
	tui.GlobalKeyMap
}

func (k IssuesKeyMap) ShortHelp() []key.Binding {
	return slices.Concat(
		[]key.Binding{k.Details, k.Transition},
		tui.CommonKeys,
	)
}

func (k IssuesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		tui.CommonKeys,
		{k.Details, k.Transition},
	}
}

var Keybinds = IssuesKeyMap{
	Details: key.NewBinding(
		key.WithKeys("enter", "l"),
		key.WithHelp("enter", "details"),
	),
	Transition: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "transition"),
	),
	GlobalKeyMap: tui.GlobalKeys(false),
}
