package app

import tea "charm.land/bubbletea/v2"

func (m Model) View() tea.View {
	v := m.Shell.RenderView()
	w := m.Shell.Ctx.Window.Width
	h := m.Shell.Ctx.Window.Height
	if m.confirmPopover.IsOpen() {
		screen := m.confirmPopover.View(v.Content, w, h)
		return tea.View{Content: screen, AltScreen: v.AltScreen}
	}
	if m.transitionPicker.IsOpen() {
		screen := m.transitionPicker.View(v.Content, w, h)
		return tea.View{Content: screen, AltScreen: v.AltScreen}
	}
	return v
}
