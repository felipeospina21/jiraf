package jira

import (
	"fmt"
	"net/url"
	"time"
)

const issueFields = "summary,description,status,priority,issuetype,assignee,reporter,created,updated,labels,components,fixVersions,subtasks,comment,customfield_10106,customfield_10105"

// GetMyIssues fetches issues assigned to the current user for a given project key.
func (c *Client) GetMyIssues(projectKey string) ([]Issue, error) {
	if c.devMode {
		time.Sleep(1 * time.Second)
		return MockIssues, nil
	}

	jql := fmt.Sprintf("assignee=currentUser() AND project=%s AND status NOT IN (Done, Withdrawn) ORDER BY priority DESC, updated DESC", projectKey)
	path := fmt.Sprintf("/rest/api/2/search?jql=%s&fields=%s&expand=renderedFields&maxResults=50", url.QueryEscape(jql), issueFields)

	var resp SearchResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return resp.Issues, nil
}
