package jira

import (
	"encoding/json"
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
	Key            string         `json:"key"`
	Fields         IssueFields    `json:"fields"`
	RenderedFields RenderedFields `json:"renderedFields"`
}

type RenderedFields struct {
	Description string `json:"description"`
}

type IssueFields struct {
	Summary     string        `json:"summary"`
	Description string        `json:"description"`
	Status      NameField     `json:"status"`
	Priority    NameField     `json:"priority"`
	IssueType   NameField     `json:"issuetype"`
	Assignee    *UserField    `json:"assignee"`
	Reporter    *UserField    `json:"reporter"`
	Created     JiraTime      `json:"created"`
	Updated     JiraTime      `json:"updated"`
	Labels      []string      `json:"labels"`
	Components  []NameField   `json:"components"`
	FixVersions []NameField   `json:"fixVersions"`
	Subtasks    []interface{} `json:"subtasks"`
	Comment     CommentField  `json:"comment"`
	StoryPoints *float64      `json:"customfield_10106"`
	SprintRaw   SprintRawField `json:"customfield_10105"`
}

// SprintRawField handles the sprint custom field which can be null, a string, or an array of strings.
type SprintRawField []string

func (s *SprintRawField) UnmarshalJSON(b []byte) error {
	str := strings.TrimSpace(string(b))
	if str == "null" || str == "" {
		return nil
	}
	// Try array of strings first
	var arr []string
	if err := json.Unmarshal(b, &arr); err == nil {
		*s = arr
		return nil
	}
	// Try single string
	var single string
	if err := json.Unmarshal(b, &single); err == nil {
		*s = []string{single}
		return nil
	}
	return nil
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

// Sprint holds parsed sprint data from the Java toString format.
type Sprint struct {
	Name      string
	StartDate string
	EndDate   string
}

// ParseSprint extracts name, startDate, and endDate from the Java toString sprint string.
func ParseSprint(raw string) Sprint {
	get := func(key string) string {
		prefix := key + "="
		idx := strings.Index(raw, prefix)
		if idx < 0 {
			return ""
		}
		start := idx + len(prefix)
		rest := raw[start:]
		end := strings.IndexAny(rest, ",]")
		if end < 0 {
			return rest
		}
		v := rest[:end]
		if v == "<null>" {
			return ""
		}
		return v
	}
	return Sprint{
		Name:      get("name"),
		StartDate: get("startDate"),
		EndDate:   get("endDate"),
	}
}

// ActiveSprint returns the parsed active sprint from SprintRaw, or nil if none.
func (f IssueFields) ActiveSprint() *Sprint {
	for _, raw := range f.SprintRaw {
		if strings.Contains(raw, "state=ACTIVE") {
			s := ParseSprint(raw)
			return &s
		}
	}
	if len(f.SprintRaw) > 0 {
		s := ParseSprint(f.SprintRaw[len(f.SprintRaw)-1])
		return &s
	}
	return nil
}

type SearchResponse struct {
	Total  int     `json:"total"`
	Issues []Issue `json:"issues"`
}
