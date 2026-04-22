package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

var presets = map[string]style.Theme{
	"catppuccin-mocha": {
		Primary:       lipgloss.Color("#89b4fa"),  // Blue
		PrimaryBright: lipgloss.Color("#b4befe"),  // Lavender
		PrimaryFg:     lipgloss.Color("#cdd6f4"),  // Text
		PrimaryDim:    lipgloss.Color("#45475a"),  // Surface1

		Info:          lipgloss.Color("#74c7ec"),  // Sapphire
		InfoBright:    lipgloss.Color("#89dceb"),  // Sky
		Success:       lipgloss.Color("#a6e3a1"),  // Green
		SuccessBright: lipgloss.Color("#94e2d5"),  // Teal
		Danger:        lipgloss.Color("#f38ba8"),  // Red
		DangerBright:  lipgloss.Color("#eba0ac"),  // Maroon
		Warning:       lipgloss.Color("#f9e2af"),  // Yellow
		WarningBright: lipgloss.Color("#fab387"),  // Peach
		Caution:       lipgloss.Color("#fab387"),  // Peach

		Text:            lipgloss.Color("#cdd6f4"),  // Text
		TextInverse:     lipgloss.Color("#11111b"),  // Crust
		TextDimmed:      lipgloss.Color("#a6adc8"),  // Subtext0
		Muted:           lipgloss.Color("#9399b2"),  // Overlay2
		Dim:             lipgloss.Color("#6c7086"),  // Overlay0
		Border:          lipgloss.Color("#585b70"),  // Surface2
		ModalBorder:     lipgloss.Color("#7f849c"),  // Overlay1
		SurfaceDim:      lipgloss.Color("#1e1e2e"),  // Base
		SelectionBorder: lipgloss.Color("#cba6f7"),  // Mauve

		StatusText:    lipgloss.Color("#cdd6f4"),  // Text
		StatusNormal:  lipgloss.Color("#45475a"),  // Surface1
		StatusLoading: lipgloss.Color("#94e2d5"),  // Teal
		StatusError:   lipgloss.Color("#f38ba8"),  // Red
		StatusDev:     lipgloss.Color("#a6e3a1"),  // Green
		StatusAccent1: lipgloss.Color("#cba6f7"),  // Mauve
		StatusAccent2: lipgloss.Color("#b4befe"),  // Lavender
	},
}
