package jira

import (
	"fmt"
	"strings"
	"time"
)

// JiraTime handles Jira's non-standard time format.
type JiraTime struct {
	time.Time
}

func (jt *JiraTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}
	// Try standard format first, then Jira's format without colon in tz
	for _, layout := range []string{
		"2006-01-02T15:04:05.000-0700",
		"2006-01-02T15:04:05.000Z0700",
		time.RFC3339,
	} {
		t, err := time.Parse(layout, s)
		if err == nil {
			jt.Time = t
			return nil
		}
	}
	return fmt.Errorf("cannot parse jira time: %s", s)
}

// Issue represents a Jira issue from the search API.
type Issue struct {
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

type IssueFields struct {
	Summary     string        `json:"summary"`
	Status      NameField     `json:"status"`
	Priority    NameField     `json:"priority"`
	IssueType   NameField     `json:"issuetype"`
	Assignee    *UserField    `json:"assignee"`
	Reporter    *UserField    `json:"reporter"`
	Created     JiraTime      `json:"created"`
	Updated     JiraTime      `json:"updated"`
	DueDate     string        `json:"duedate"`
	Labels      []string      `json:"labels"`
	Components  []NameField   `json:"components"`
	FixVersions []NameField   `json:"fixVersions"`
	Subtasks    []interface{} `json:"subtasks"`
	Comment     CommentField  `json:"comment"`
	Sprint      *SprintField  `json:"sprint"`
}

type NameField struct {
	Name string `json:"name"`
}

type UserField struct {
	DisplayName string `json:"displayName"`
}

type CommentField struct {
	Total int `json:"total"`
}

type SprintField struct {
	Name string `json:"name"`
}

type SearchResponse struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}
