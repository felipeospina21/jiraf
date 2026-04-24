package app

import (
	tea "charm.land/bubbletea/v2"
	"github.com/felipeospina21/jiraf/internal/tui/issues"
)

func (m Model) fetchIssues(projectKey string) tea.Cmd {
	filters := m.activeFilters
	return func() tea.Msg {
		iss, err := m.client.GetMyIssues(projectKey, filters)
		return issues.FetchedMsg{Issues: iss, Err: err}
	}
}

func (m Model) fetchTransitions(issueKey string) tea.Cmd {
	return func() tea.Msg {
		t, err := m.client.GetTransitions(issueKey)
		return TransitionsFetchedMsg{Transitions: t, Err: err}
	}
}

func (m Model) doTransition(issueKey, transitionID string) tea.Cmd {
	return func() tea.Msg {
		err := m.client.DoTransition(issueKey, transitionID)
		return TransitionDoneMsg{Err: err}
	}
}
