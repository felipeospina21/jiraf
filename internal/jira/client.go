package jira

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/felipeospina21/jiraf/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	demoMode   bool
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		token:      cfg.APIToken,
		demoMode:   cfg.DemoMode,
		httpClient: &http.Client{
			Timeout:       10 * time.Second,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
	}
}

func (c *Client) post(path string, body any) error {
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.baseURL+path, strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("jira: HTTP %d — %s", resp.StatusCode, b)
	}
	return nil
}

// AddComment posts a comment on an issue.
func (c *Client) AddComment(issueKey, body string) error {
	if c.demoMode {
		time.Sleep(500 * time.Millisecond)
		return nil
	}
	return c.post(fmt.Sprintf("/rest/api/2/issue/%s/comment", issueKey), map[string]string{"body": body})
}

func (c *Client) get(path string, out any) error {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return fmt.Errorf("jira: authentication redirect (HTTP %d) — check your VPN connection and JIRAF_TOKEN", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		return fmt.Errorf("jira: unexpected response (HTTP %d, %s) — check VPN connection and JIRAF_TOKEN", resp.StatusCode, ct)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("jira: HTTP %d — %s", resp.StatusCode, body)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
