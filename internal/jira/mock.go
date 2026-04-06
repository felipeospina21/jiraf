package jira

import "time"

var now = time.Now()

func jt(t time.Time) JiraTime { return JiraTime{Time: t} }

var MockIssues = []Issue{
	{
		Key: "UCP-1042",
		Fields: IssueFields{
			Summary:     "Fix session timeout on idle users",
			Status:      NameField{Name: "In Progress"},
			Priority:    NameField{Name: "High"},
			IssueType:   NameField{Name: "Bug"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "Sarah Chen"},
			Created:     jt(now.Add(-5 * 24 * time.Hour)),
			Updated:     jt(now.Add(-2 * time.Hour)),
			DueDate:     now.Add(3 * 24 * time.Hour).Format("2006-01-02"),
			Labels:      []string{"backend", "auth"},
			Components:  []NameField{{Name: "API"}},
			FixVersions: []NameField{{Name: "v2.3"}},
			Subtasks:    make([]interface{}, 2),
			Comment:     CommentField{Total: 4},
			Sprint:      &SprintField{Name: "Sprint 42"},
		},
	},
	{
		Key: "UCP-1038",
		Fields: IssueFields{
			Summary:     "Add pagination to dashboard list endpoint",
			Status:      NameField{Name: "To Do"},
			Priority:    NameField{Name: "Medium"},
			IssueType:   NameField{Name: "Story"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "James Park"},
			Created:     jt(now.Add(-8 * 24 * time.Hour)),
			Updated:     jt(now.Add(-1 * 24 * time.Hour)),
			Labels:      []string{"backend"},
			Components:  []NameField{{Name: "API"}, {Name: "Dashboard"}},
			FixVersions: []NameField{{Name: "v2.3"}},
			Comment:     CommentField{Total: 1},
			Sprint:      &SprintField{Name: "Sprint 42"},
		},
	},
	{
		Key: "UCP-1035",
		Fields: IssueFields{
			Summary:     "Migrate user preferences to new schema",
			Status:      NameField{Name: "In Review"},
			Priority:    NameField{Name: "Medium"},
			IssueType:   NameField{Name: "Task"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "Maria Lopez"},
			Created:     jt(now.Add(-12 * 24 * time.Hour)),
			Updated:     jt(now.Add(-6 * time.Hour)),
			DueDate:     now.Add(1 * 24 * time.Hour).Format("2006-01-02"),
			Labels:      []string{"database", "migration"},
			Components:  []NameField{{Name: "Database"}},
			Subtasks:    make([]interface{}, 3),
			Comment:     CommentField{Total: 7},
			Sprint:      &SprintField{Name: "Sprint 42"},
		},
	},
	{
		Key: "UCP-1029",
		Fields: IssueFields{
			Summary:     "Investigate memory spike on report generation",
			Status:      NameField{Name: "In Progress"},
			Priority:    NameField{Name: "Critical"},
			IssueType:   NameField{Name: "Bug"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "Alex Rivera"},
			Created:     jt(now.Add(-3 * 24 * time.Hour)),
			Updated:     jt(now.Add(-30 * time.Minute)),
			Labels:      []string{"performance"},
			Components:  []NameField{{Name: "Reports"}},
			Comment:     CommentField{Total: 2},
			Sprint:      &SprintField{Name: "Sprint 42"},
		},
	},
	{
		Key: "UCP-1021",
		Fields: IssueFields{
			Summary:     "Update API docs for v2.3 endpoints",
			Status:      NameField{Name: "Done"},
			Priority:    NameField{Name: "Low"},
			IssueType:   NameField{Name: "Task"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "Felipe Ospina"},
			Created:     jt(now.Add(-15 * 24 * time.Hour)),
			Updated:     jt(now.Add(-2 * 24 * time.Hour)),
			Labels:      []string{"docs"},
			FixVersions: []NameField{{Name: "v2.3"}},
			Comment:     CommentField{Total: 0},
			Sprint:      &SprintField{Name: "Sprint 41"},
		},
	},
}
