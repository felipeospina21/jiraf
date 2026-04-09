package jira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/felipeospina21/jiraf/internal/config"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
	devMode    bool
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		token:      cfg.APIToken,
		devMode:    cfg.DevMode,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jira: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}
