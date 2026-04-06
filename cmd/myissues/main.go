package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type issue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary  string `json:"summary"`
		Status   struct{ Name string } `json:"status"`
		Priority struct{ Name string } `json:"priority"`
		IssueType struct{ Name string } `json:"issuetype"`
	} `json:"fields"`
}

func main() {
	_ = godotenv.Load()

	baseURL := os.Getenv("MRJIRA_BASE_URL")
	token := os.Getenv("MRJIRA_TOKEN")

	if baseURL == "" || token == "" {
		fmt.Fprintln(os.Stderr, "Set MRJIRA_BASE_URL and MRJIRA_TOKEN in .env")
		os.Exit(1)
	}

	// Default to UCP project key, or pass as arg
	projectKey := "UCP"
	if len(os.Args) > 1 {
		projectKey = os.Args[1]
	}

	jql := fmt.Sprintf("assignee=currentUser() AND project=%s AND sprint in openSprints() ORDER BY priority DESC", projectKey)

	u := fmt.Sprintf("%s/rest/api/2/search?jql=%s&fields=summary,status,priority,issuetype&maxResults=50",
		baseURL, url.QueryEscape(jql))

	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, string(body[:min(len(body), 300)]))
		os.Exit(1)
	}

	var result struct {
		Total  int     `json:"total"`
		Issues []issue `json:"issues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Found %d issues assigned to you\n\n", result.Total)
	fmt.Printf("%-12s %-10s %-12s %-10s %s\n", "Key", "Type", "Priority", "Status", "Summary")
	fmt.Println(string(make([]byte, 80)))
	for _, i := range result.Issues {
		fmt.Printf("%-12s %-10s %-12s %-10s %s\n",
			i.Key,
			i.Fields.IssueType.Name,
			i.Fields.Priority.Name,
			i.Fields.Status.Name,
			i.Fields.Summary,
		)
	}
}
