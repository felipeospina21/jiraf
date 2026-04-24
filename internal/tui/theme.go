package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/jiraf/internal/config"
	"github.com/felipeospina21/tuishell/style"
)

func applyOverride(target *color.Color, val *string) {
	if val != nil {
		*target = lipgloss.Color(*val)
	}
}

// BuildTheme returns DefaultTheme with any non-nil overrides applied.
func BuildTheme(overrides config.ThemeOverrides) style.Theme {
	t := DefaultTheme()
	if overrides.Preset != nil {
		if preset, ok := style.Presets[*overrides.Preset]; ok {
			t = preset
		}
	}
	applyOverride(&t.Primary, overrides.Primary)
	applyOverride(&t.PrimaryBright, overrides.PrimaryBright)
	applyOverride(&t.PrimaryFg, overrides.PrimaryFg)
	applyOverride(&t.PrimaryDim, overrides.PrimaryDim)
	applyOverride(&t.Info, overrides.Info)
	applyOverride(&t.InfoBright, overrides.InfoBright)
	applyOverride(&t.Success, overrides.Success)
	applyOverride(&t.SuccessBright, overrides.SuccessBright)
	applyOverride(&t.Danger, overrides.Danger)
	applyOverride(&t.DangerBright, overrides.DangerBright)
	applyOverride(&t.Warning, overrides.Warning)
	applyOverride(&t.WarningBright, overrides.WarningBright)
	applyOverride(&t.Caution, overrides.Caution)
	applyOverride(&t.Text, overrides.Text)
	applyOverride(&t.TextInverse, overrides.TextInverse)
	applyOverride(&t.TextDimmed, overrides.TextDimmed)
	applyOverride(&t.Muted, overrides.Muted)
	applyOverride(&t.Dim, overrides.Dim)
	applyOverride(&t.Border, overrides.Border)
	applyOverride(&t.ModalBorder, overrides.ModalBorder)
	applyOverride(&t.SurfaceDim, overrides.SurfaceDim)
	applyOverride(&t.SelectionBorder, overrides.SelectionBorder)
	applyOverride(&t.StatusText, overrides.StatusText)
	applyOverride(&t.StatusNormal, overrides.StatusNormal)
	applyOverride(&t.StatusLoading, overrides.StatusLoading)
	applyOverride(&t.StatusError, overrides.StatusError)
	applyOverride(&t.StatusDemo, overrides.StatusDemo)
	applyOverride(&t.StatusAccent1, overrides.StatusAccent1)
	applyOverride(&t.StatusAccent2, overrides.StatusAccent2)
	return t
}

// DefaultTheme returns a Jira/Atlassian blue theme.
func DefaultTheme() style.Theme {
	return style.Theme{
		Primary:       lipgloss.Color("#2684FF"),
		PrimaryBright: lipgloss.Color("#8AC4FF"),
		PrimaryFg:     lipgloss.Color("#DEEBFF"),
		PrimaryDim:    lipgloss.Color("#002B6B"),

		Info:          lipgloss.Color("#4C9AFF"),
		InfoBright:    lipgloss.Color("#2684FF"),
		Success:       lipgloss.Color("#6beaaf"),
		SuccessBright: lipgloss.Color("#3ad994"),
		Danger:        lipgloss.Color("#f9a8a8"),
		DangerBright:  lipgloss.Color("#f47575"),
		Warning:       lipgloss.Color("#ffe043"),
		WarningBright: lipgloss.Color("#ffcc14"),
		Caution:       lipgloss.Color("#ff8237"),

		Text:            lipgloss.Color("#C4C4C4"),
		TextInverse:     lipgloss.Color("#111"),
		TextDimmed:      lipgloss.Color("#777777"),
		Muted:           lipgloss.Color("#999999"),
		Dim:             lipgloss.Color("#444444"),
		Border:          lipgloss.Color("#3f4145"),
		ModalBorder:     lipgloss.Color("#666666"),
		SurfaceDim:      lipgloss.Color("#091E42"),
		SelectionBorder: lipgloss.Color("#0052CC"),

		StatusText:    lipgloss.Color("#FFFDF5"),
		StatusNormal:  lipgloss.Color("#0747A6"),
		StatusLoading: lipgloss.Color("#1A7A94"),
		StatusError:   lipgloss.Color("#CE3060"),
		StatusDemo:     lipgloss.Color("#4E8212"),
		StatusAccent1: lipgloss.Color("#0065FF"),
		StatusAccent2: lipgloss.Color("#003884"),
	}
}
