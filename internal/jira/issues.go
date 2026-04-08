package jira

import (
	"fmt"
	"net/url"
	"time"
)

const issueFields = "summary,status,priority,issuetype,assignee,reporter,created,updated,duedate,labels,components,fixVersions,subtasks,comment,sprint"

// GetMyIssues fetches issues assigned to the current user for a given project key.
func (c *Client) GetMyIssues(projectKey string) ([]Issue, error) {
	if c.devMode {
		time.Sleep(1 * time.Second)
		return MockIssues, nil
	}

	jql := fmt.Sprintf("assignee=currentUser() AND project=%s ORDER BY priority ASC, updated DESC", projectKey)
	path := fmt.Sprintf("/rest/api/2/search?jql=%s&fields=%s&maxResults=50", url.QueryEscape(jql), issueFields)

	var resp SearchResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return resp.Issues, nil
}
