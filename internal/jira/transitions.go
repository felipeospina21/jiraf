package jira

import (
	"fmt"
	"time"
)

// Transition represents an available status transition for an issue.
type Transition struct {
	ID   string    `json:"id"`
	Name string    `json:"name"`
	To   NameField `json:"to"`
}

type transitionsResponse struct {
	Transitions []Transition `json:"transitions"`
}

// GetTransitions fetches available transitions for an issue.
func (c *Client) GetTransitions(issueKey string) ([]Transition, error) {
	if c.devMode {
		time.Sleep(500 * time.Millisecond)
		return MockTransitions, nil
	}
	var resp transitionsResponse
	if err := c.get(fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueKey), &resp); err != nil {
		return nil, err
	}
	return resp.Transitions, nil
}

// DoTransition executes a status transition on an issue.
func (c *Client) DoTransition(issueKey, transitionID string) error {
	if c.devMode {
		time.Sleep(500 * time.Millisecond)
		return nil
	}
	body := map[string]any{"transition": map[string]string{"id": transitionID}}
	return c.post(fmt.Sprintf("/rest/api/2/issue/%s/transitions", issueKey), body)
}
