package app

import (
	"fmt"
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/config"
	jirafExec "github.com/felipeospina21/jiraf/internal/exec"
	"github.com/felipeospina21/jiraf/internal/jira"
	"github.com/felipeospina21/jiraf/internal/tui/boards"
	"github.com/felipeospina21/jiraf/internal/tui/issues"
	"github.com/felipeospina21/tuishell"
)

// TransitionsFetchedMsg carries the result of fetching available transitions.
type TransitionsFetchedMsg struct {
	Transitions []jira.Transition
	Err         error
}

// TransitionDoneMsg carries the result of executing a transition.
type TransitionDoneMsg struct {
	Err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case boards.SelectBoardMsg:
		// Reset filters and known values when switching boards
		m.activeFilters = jira.IssueFilters{}
		m.knownStatuses = nil
		m.knownPriorities = nil
		m.knownTypes = nil
		if main, ok := m.Shell.Main.(issues.Model); ok {
			main.SelectedBoard = msg.Key
			main.Loading = true
			main.SpinnerView = m.Shell.Spinner.View()
			m.Shell.Main = main
		}
		return m, tea.Batch(
			func() tea.Msg { return tuishell.CloseLeftPanelMsg{} },
			func() tea.Msg { return tuishell.StartTaskMsg{Cmd: m.fetchIssues(msg.Key)} },
		)

	case issues.FetchedMsg:
		// Collect known filter values from fetched issues
		if msg.Err == nil {
			if m.knownStatuses == nil {
				m.knownStatuses = map[string]bool{}
			}
			if m.knownPriorities == nil {
				m.knownPriorities = map[string]bool{}
			}
			if m.knownTypes == nil {
				m.knownTypes = map[string]bool{}
			}
			for _, i := range msg.Issues {
				m.knownStatuses[i.Fields.Status.Name] = true
				m.knownPriorities[i.Fields.Priority.Name] = true
				m.knownTypes[i.Fields.IssueType.Name] = true
			}
		}

	case issues.ViewDetailsMsg:
		m.Details.SetContent(msg.Issue)
		var cmds []tea.Cmd
		if !m.Shell.IsRightOpen() {
			cmds = append(cmds, func() tea.Msg { return tuishell.OpenRightPanelMsg{} })
		}
		return m, tea.Batch(cmds...)

	case issues.OpenInBrowserMsg:
		url := config.GlobalConfig.BaseURL + "/browse/" + msg.IssueKey
		jirafExec.OpenBrowser(url)
		return m, nil

	case issues.RefetchMsg:
		if main, ok := m.Shell.Main.(issues.Model); ok && main.SelectedBoard != "" {
			main.Loading = true
			main.SpinnerView = m.Shell.Spinner.View()
			m.Shell.Main = main
			return m, func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: m.fetchIssues(main.SelectedBoard)}
			}
		}

	case issues.TransitionMsg:
		m.pendingTransition = msg.Issue.Key
		return m, func() tea.Msg {
			return tuishell.StartTaskMsg{Cmd: m.fetchTransitions(msg.Issue.Key)}
		}

	case TransitionsFetchedMsg:
		if msg.Err != nil {
			m.pendingTransition = ""
			return m, func() tea.Msg {
				return tuishell.FinishTaskMsg{Err: msg.Err}
			}
		}
		items := make([]tuishell.ListPopoverItem, len(msg.Transitions))
		m.transitionIDByName = make(map[string]string, len(msg.Transitions))
		for i, t := range msg.Transitions {
			items[i] = tuishell.ListPopoverItem{Label: t.Name, Value: t.Name}
			m.transitionIDByName[t.Name] = t.ID
		}
		m.transitionPicker.Open(fmt.Sprintf("Transition %s", m.pendingTransition), items)
		return m, func() tea.Msg { return tuishell.FinishTaskMsg{Err: nil} }

	case tuishell.SelectListPopoverMsg:
		if m.pendingTransition != "" {
			m.pendingTransitionName = msg.Value
			m.transitionPicker.Close()
			m.confirmPopover.Open(
				"Transition Issue",
				fmt.Sprintf("Transition %s to %s?", m.pendingTransition, msg.Value),
				"Confirm",
				"Cancel",
			)
		}

	case tuishell.ConfirmPopoverYesMsg:
		m.confirmPopover.Close()
		issueKey := m.pendingTransition
		transitionID := m.transitionIDByName[m.pendingTransitionName]
		m.pendingTransition = ""
		m.pendingTransitionName = ""
		m.transitionIDByName = nil
		return m, func() tea.Msg {
			return tuishell.StartTaskMsg{Cmd: m.doTransition(issueKey, transitionID)}
		}

	case tuishell.ConfirmPopoverNoMsg:
		m.confirmPopover.Close()
		m.pendingTransition = ""
		m.pendingTransitionName = ""
		m.transitionIDByName = nil

	case tuishell.CloseListPopoverMsg:
		m.pendingTransition = ""
		m.transitionIDByName = nil
		m.transitionPicker.Close()

	case issues.OpenFilterMsg:
		// Include active filter values in known sets so they always appear
		for _, v := range m.activeFilters.Statuses {
			if m.knownStatuses == nil {
				m.knownStatuses = map[string]bool{}
			}
			m.knownStatuses[v] = true
		}
		for _, v := range m.activeFilters.Priorities {
			if m.knownPriorities == nil {
				m.knownPriorities = map[string]bool{}
			}
			m.knownPriorities[v] = true
		}
		for _, v := range m.activeFilters.Types {
			if m.knownTypes == nil {
				m.knownTypes = map[string]bool{}
			}
			m.knownTypes[v] = true
		}
		sections := buildFilterSections(m.knownStatuses, m.knownPriorities, m.knownTypes, m.activeFilters)
		m.filterPopover.Open(sections)
		return m, nil

	case tuishell.ApplyFilterPopoverMsg:
		m.filterPopover.Close()
		m.activeFilters = jira.IssueFilters{
			Statuses:   msg.Selections["Status"],
			Priorities: msg.Selections["Priority"],
			Types:      msg.Selections["Type"],
		}
		if main, ok := m.Shell.Main.(issues.Model); ok && main.SelectedBoard != "" {
			main.Loading = true
			main.SpinnerView = m.Shell.Spinner.View()
			m.Shell.Main = main
			var cmds []tea.Cmd
			cmds = append(cmds, func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: m.fetchIssues(main.SelectedBoard)}
			})
			if m.activeFilters.HasActive() {
				cmds = append(cmds, func() tea.Msg {
					return tuishell.SetStatusMsg{Content: filterStatusText(m.activeFilters)}
				})
			} else {
				cmds = append(cmds, func() tea.Msg {
					return tuishell.SetStatusMsg{Content: ""}
				})
			}
			return m, tea.Batch(cmds...)
		}

	case tuishell.CloseFilterPopoverMsg:
		m.filterPopover.Close()

	case issues.ClearFilterMsg:
		m.activeFilters = jira.IssueFilters{}
		if main, ok := m.Shell.Main.(issues.Model); ok && main.SelectedBoard != "" {
			main.Loading = true
			main.SpinnerView = m.Shell.Spinner.View()
			m.Shell.Main = main
			return m, tea.Batch(
				func() tea.Msg { return tuishell.SetStatusMsg{Content: ""} },
				func() tea.Msg { return tuishell.StartTaskMsg{Cmd: m.fetchIssues(main.SelectedBoard)} },
			)
		}

	case TransitionDoneMsg:
		if msg.Err != nil {
			return m, func() tea.Msg {
				return tuishell.FinishTaskMsg{Err: msg.Err}
			}
		}
		var cmds []tea.Cmd
		cmds = append(cmds, func() tea.Msg {
			return tuishell.FinishTaskMsg{Err: nil, Keybinds: issues.Keybinds}
		})
		cmds = append(cmds, func() tea.Msg {
			return tuishell.SetStatusMsg{Content: "✓ Transition complete"}
		})
		// Refetch issues if we have a selected board
		if main, ok := m.Shell.Main.(issues.Model); ok && main.SelectedBoard != "" {
			main.Loading = true
			m.Shell.Main = main
			cmds = append(cmds, func() tea.Msg {
				return tuishell.StartTaskMsg{Cmd: m.fetchIssues(main.SelectedBoard)}
			})
		}
		return m, tea.Batch(cmds...)

	case tuishell.FinishTaskMsg:
		if main, ok := m.Shell.Main.(issues.Model); ok {
			main.Loading = false
			m.Shell.Main = main
		}
	}

	// Route keys to confirm popover when open
	if m.confirmPopover.IsOpen() {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			var cmd tea.Cmd
			m.confirmPopover, cmd = m.confirmPopover.Update(msg)
			return m, cmd
		}
	}

	// Route keys to transition picker when open
	if m.transitionPicker.IsOpen() {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			var cmd tea.Cmd
			m.transitionPicker, cmd = m.transitionPicker.Update(msg)
			return m, cmd
		}
	}

	// Route keys to filter popover when open
	if m.filterPopover.IsOpen() {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			var cmd tea.Cmd
			m.filterPopover, cmd = m.filterPopover.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.Shell, cmd = m.Shell.Update(msg)
	m.syncKeybinds()

	// Sync panel pointer after shell update
	if d, ok := m.Shell.Right.(DetailsPanel); ok {
		m.Details = d.Model
	}

	// Update spinner view on issues panel when loading
	if main, ok := m.Shell.Main.(issues.Model); ok && main.Loading {
		main.SpinnerView = m.Shell.Spinner.View()
		m.Shell.Main = main
	}

	return m, cmd
}

func buildFilterSections(statuses, priorities, types map[string]bool, active jira.IssueFilters) []tuishell.FilterSection {
	makeSection := func(title string, vals map[string]bool, selected []string) tuishell.FilterSection {
		sel := map[string]bool{}
		for _, v := range selected {
			sel[v] = true
		}
		keys := make([]string, 0, len(vals))
		for k := range vals {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		opts := make([]tuishell.FilterOption, len(keys))
		for i, k := range keys {
			opts[i] = tuishell.FilterOption{Label: k, Value: k, Selected: sel[k]}
		}
		return tuishell.FilterSection{Title: title, Options: opts}
	}

	return []tuishell.FilterSection{
		makeSection("Status", statuses, active.Statuses),
		makeSection("Priority", priorities, active.Priorities),
		makeSection("Type", types, active.Types),
	}
}

func filterStatusText(f jira.IssueFilters) string {
	var parts []string
	if len(f.Statuses) > 0 {
		parts = append(parts, fmt.Sprintf("Status(%d)", len(f.Statuses)))
	}
	if len(f.Priorities) > 0 {
		parts = append(parts, fmt.Sprintf("Priority(%d)", len(f.Priorities)))
	}
	if len(f.Types) > 0 {
		parts = append(parts, fmt.Sprintf("Type(%d)", len(f.Types)))
	}
	return "Filtered: " + strings.Join(parts, ", ")
}
