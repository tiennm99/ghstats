// Package github fetches profile data from the GitHub GraphQL API.
package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const endpoint = "https://api.github.com/graphql"

// Client issues authenticated GraphQL requests.
type Client struct {
	token string
	http  *http.Client
}

// NewClient returns a client authenticated with the given PAT. An empty token
// falls back to unauthenticated requests (60/h rate limit, no private data).
func NewClient(token string) *Client {
	return &Client{
		token: token,
		http:  &http.Client{Timeout: 30 * time.Second},
	}
}

type gqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type gqlError struct {
	Message string   `json:"message"`
	Type    string   `json:"type,omitempty"`
	Path    []string `json:"path,omitempty"`
}

type gqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []gqlError      `json:"errors,omitempty"`
}

// query runs a GraphQL query and unmarshals the `data` field into out.
func (c *Client) query(q string, vars map[string]any, out any) error {
	body, err := json.Marshal(gqlRequest{Query: q, Variables: vars})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ghstats")
	if c.token != "" {
		req.Header.Set("Authorization", "bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("http %d: %s", resp.StatusCode, truncate(raw, 500))
	}

	var r gqlResponse
	if err := json.Unmarshal(raw, &r); err != nil {
		return fmt.Errorf("decode body: %w", err)
	}
	if len(r.Errors) > 0 {
		msgs := make([]string, 0, len(r.Errors))
		for _, e := range r.Errors {
			msgs = append(msgs, e.Message)
		}
		return fmt.Errorf("graphql: %s", joinErrs(msgs))
	}
	if out != nil {
		if err := json.Unmarshal(r.Data, out); err != nil {
			return fmt.Errorf("decode data: %w", err)
		}
	}
	return nil
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "…"
}

func joinErrs(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	out := ss[0]
	for _, s := range ss[1:] {
		out += "; " + s
	}
	return out
}

