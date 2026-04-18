// Package github fetches profile data from the GitHub API.
package github

import "errors"

// Profile is the aggregate of data other packages render into cards.
// Fields are stubs; flesh them out as cards are implemented.
type Profile struct {
	Login       string
	Name        string
	Bio         string
	Followers   int
	Following   int
	PublicRepos int

	// Top languages aggregated across repos (name → bytes).
	Languages map[string]int64

	// Commit-count histogram indexed by [day-of-week][hour-of-day], local tz.
	Productive [7][24]int
}

// Client wraps GitHub REST + GraphQL access.
type Client struct {
	token string
	// TODO: http.Client, rate-limit handling
}

// NewClient returns a client that authenticates with the given PAT.
// Empty token uses unauthenticated access (low rate limit).
func NewClient(token string) *Client {
	return &Client{token: token}
}

// Profile loads the profile summary for a user.
func (c *Client) Profile(user string) (*Profile, error) {
	if user == "" {
		return nil, errors.New("empty user")
	}
	// TODO: fetch via GraphQL: viewer, user.repositories, contributionsCollection
	return &Profile{Login: user}, nil
}
