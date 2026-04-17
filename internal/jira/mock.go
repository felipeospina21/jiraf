package jira

import "time"

var now = time.Now()

func jt(t time.Time) JiraTime { return JiraTime{Time: t} }

func sp(v float64) *float64 { return &v }

var mockSprint = "com.atlassian.greenhopper.service.sprint.Sprint@abc[name=Sprint 42,state=ACTIVE,startDate=2026-04-01T08:00:00.000-05:00,endDate=2026-04-15T08:00:00.000-05:00]"

var MockTransitions = []Transition{
	{ID: "11", Name: "To Do", To: NameField{Name: "To Do"}},
	{ID: "21", Name: "In Progress", To: NameField{Name: "In Progress"}},
	{ID: "31", Name: "In Review", To: NameField{Name: "In Review"}},
	{ID: "41", Name: "Done", To: NameField{Name: "Done"}},
}

var MockIssues = []Issue{
	{
		Key: "PROJ-1042",
		Fields: IssueFields{
			Summary:     "Fix session timeout on idle users",
			Description: "h3. Problem\nUsers are being logged out after 5 minutes of inactivity.\n\nh3. Acceptance Criteria\n# Session TTL should be 30 minutes per the auth spec\n## Verify with integration tests\n# Idle detection must use server-side timestamps",
			Status:      NameField{Name: "In Progress"},
			Priority:    NameField{Name: "High"},
			IssueType:   NameField{Name: "Bug"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "Sarah Chen"},
			Created:     jt(now.Add(-5 * 24 * time.Hour)),
			Updated:     jt(now.Add(-2 * time.Hour)),
			Labels:      []string{"backend", "auth"},
			Components:  []NameField{{Name: "API"}},
			FixVersions: []NameField{{Name: "v2.3"}},
			Subtasks:    make([]interface{}, 2),
			Comment:     CommentField{Total: 4},
			StoryPoints: sp(5),
			SprintRaw:   []string{mockSprint},
		},
		RenderedFields: RenderedFields{
			Description: "<h3>Problem</h3><p>Users are being logged out after 5 minutes of inactivity.</p><h3>Acceptance Criteria</h3><ol><li>Session TTL should be 30 minutes per the auth spec<ol><li>Verify with integration tests</li></ol></li><li>Idle detection must use server-side timestamps</li></ol>",
		},
	},
	{
		Key: "PROJ-1038",
		Fields: IssueFields{
			Summary:     "Add pagination to dashboard list endpoint",
			Description: "The /api/v2/dashboards endpoint returns all results. Add cursor-based pagination with a default page size of 25.",
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
			StoryPoints: sp(3),
			SprintRaw:   []string{mockSprint},
		},
		RenderedFields: RenderedFields{
			Description: "<p>The /api/v2/dashboards endpoint returns all results. Add cursor-based pagination with a default page size of 25.</p>",
		},
	},
	{
		Key: "PROJ-1035",
		Fields: IssueFields{
			Summary:     "Migrate user preferences to new schema",
			Description: "Move user preferences from the legacy key-value table to the new typed preferences schema. Must be backward compatible.",
			Status:      NameField{Name: "In Review"},
			Priority:    NameField{Name: "Medium"},
			IssueType:   NameField{Name: "Task"},
			Assignee:    &UserField{DisplayName: "Felipe Ospina"},
			Reporter:    &UserField{DisplayName: "Maria Lopez"},
			Created:     jt(now.Add(-12 * 24 * time.Hour)),
			Updated:     jt(now.Add(-6 * time.Hour)),
			Labels:      []string{"database", "migration"},
			Components:  []NameField{{Name: "Database"}},
			Subtasks:    make([]interface{}, 3),
			Comment:     CommentField{Total: 7},
			StoryPoints: sp(8),
			SprintRaw:   []string{mockSprint},
		},
		RenderedFields: RenderedFields{
			Description: "<p>Move user preferences from the legacy key-value table to the new typed preferences schema. Must be backward compatible.</p>",
		},
	},
	{
		Key: "PROJ-1029",
		Fields: IssueFields{
			Summary:     "Investigate memory spike on report generation",
			Description: "Heap usage jumps to 4GB when generating the monthly activity report. Profile and identify the allocation hotspot.",
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
			StoryPoints: sp(2),
			SprintRaw:   []string{mockSprint},
		},
		RenderedFields: RenderedFields{
			Description: "<p>Heap usage jumps to 4GB when generating the monthly activity report. Profile and identify the allocation hotspot.</p>",
		},
	},
	{
		Key: "PROJ-1021",
		Fields: IssueFields{
			Summary:     "Update API docs for v2.3 endpoints",
			Description: "Document the new endpoints added in v2.3: batch operations, webhook management, and audit log queries.",
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
		},
		RenderedFields: RenderedFields{
			Description: "<p>Document the new endpoints added in v2.3: batch operations, webhook management, and audit log queries.</p>",
		},
	},
}
