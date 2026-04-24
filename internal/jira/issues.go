package jira

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

const issueFields = "summary,description,status,priority,issuetype,assignee,reporter,created,updated,labels,components,fixVersions,subtasks,comment,customfield_10106,customfield_10105"

// IssueFilters holds optional JQL filter criteria.
type IssueFilters struct {
	Statuses   []string
	Priorities []string
	Types      []string
	Sprint     string
}

// HasActive reports whether any filter is set.
func (f IssueFilters) HasActive() bool {
	return len(f.Statuses) > 0 || len(f.Priorities) > 0 || len(f.Types) > 0 || f.Sprint != ""
}

// GetMyIssues fetches issues assigned to the current user for a given project key.
func (c *Client) GetMyIssues(projectKey string, filters IssueFilters) ([]Issue, error) {
	if c.demoMode {
		time.Sleep(1 * time.Second)
		return MockIssues, nil
	}

	jql := fmt.Sprintf("assignee=currentUser() AND project=%s", projectKey)
	if len(filters.Statuses) > 0 {
		jql += fmt.Sprintf(` AND status IN ("%s")`, strings.Join(filters.Statuses, `", "`))
	} else {
		jql += ` AND status NOT IN (Done, Withdrawn)`
	}
	if len(filters.Priorities) > 0 {
		jql += fmt.Sprintf(` AND priority IN ("%s")`, strings.Join(filters.Priorities, `", "`))
	}
	if len(filters.Types) > 0 {
		jql += fmt.Sprintf(` AND issuetype IN ("%s")`, strings.Join(filters.Types, `", "`))
	}
	if filters.Sprint != "" {
		jql += fmt.Sprintf(` AND sprint = "Sprint %s"`, filters.Sprint)
	}

	jql += " ORDER BY priority DESC, updated DESC"

	path := fmt.Sprintf("/rest/api/2/search?jql=%s&fields=%s&expand=renderedFields&maxResults=50", url.QueryEscape(jql), issueFields)

	var resp SearchResponse
	if err := c.get(path, &resp); err != nil {
		return nil, err
	}
	return resp.Issues, nil
}
