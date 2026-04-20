package issues

import (
	"slices"

	"charm.land/bubbles/v2/key"
	"github.com/felipeospina21/jiraf/internal/tui"
)

type IssuesKeyMap struct {
	Details       key.Binding
	OpenInBrowser key.Binding
	Transition    key.Binding
	Refetch       key.Binding
	tui.GlobalKeyMap
}

func (k IssuesKeyMap) ShortHelp() []key.Binding {
	return slices.Concat(
		[]key.Binding{k.Details, k.OpenInBrowser, k.Transition, k.Refetch},
		tui.CommonKeys,
	)
}

func (k IssuesKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		tui.CommonKeys,
		{k.Details, k.OpenInBrowser, k.Transition, k.Refetch},
	}
}

var Keybinds = IssuesKeyMap{
	Details: key.NewBinding(
		key.WithKeys("enter", "l"),
		key.WithHelp("enter", "details"),
	),
	OpenInBrowser: key.NewBinding(
		key.WithKeys("x"),
		key.WithHelp("x", "open in browser"),
	),
	Transition: key.NewBinding(
		key.WithKeys("T"),
		key.WithHelp("T", "transition"),
	),
	Refetch: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refetch"),
	),
	GlobalKeyMap: tui.GlobalKeys(false),
}
