package tui

import (
	"charm.land/lipgloss/v2"
	"github.com/felipeospina21/tuishell/style"
)

// DefaultTheme returns a Jira/Atlassian blue theme.
func DefaultTheme() style.Theme {
	return style.Theme{
		Primary:       lipgloss.Color("#2684FF"),
		PrimaryBright: lipgloss.Color("#0065FF"),
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
		StatusDev:     lipgloss.Color("#4E8212"),
		StatusAccent1: lipgloss.Color("#0065FF"),
		StatusAccent2: lipgloss.Color("#003884"),
	}
}
