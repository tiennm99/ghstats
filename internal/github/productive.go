package github

import (
	"time"
)

// productiveGQL is the response shape for commitHistoryQuery.
type productiveGQL struct {
	Repository *struct {
		DefaultBranchRef *struct {
			Target *struct {
				History struct {
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
					Nodes []struct {
						CommittedDate string `json:"committedDate"`
					} `json:"nodes"`
				} `json:"history"`
			} `json:"target"`
		} `json:"defaultBranchRef"`
	} `json:"repository"`
}

// FetchProductive fills p.Productive with a [7][24] commit histogram over the
// last year, gathered from the user's top-starred owned repos. Each repo is
// sampled up to maxPerRepo commits to keep the cost bounded.
//
// The timezone loc is applied to CommittedDate so the heatmap reflects when the
// user actually commits, not UTC.
func (c *Client) FetchProductive(p *Profile, repos []string, loc *time.Location, maxPerRepo int) error {
	if loc == nil {
		loc = time.UTC
	}
	since := time.Now().AddDate(-1, 0, 0).UTC().Format(time.RFC3339)

	for _, repo := range repos {
		var cursor *string
		seen := 0
		for {
			if seen >= maxPerRepo {
				break
			}
			vars := map[string]any{
				"login":  p.Login,
				"repo":   repo,
				"userId": p.ID,
				"since":  since,
			}
			if cursor != nil {
				vars["after"] = *cursor
			}

			var resp productiveGQL
			if err := c.query(commitHistoryQuery, vars, &resp); err != nil {
				return err
			}
			if resp.Repository == nil || resp.Repository.DefaultBranchRef == nil ||
				resp.Repository.DefaultBranchRef.Target == nil {
				break
			}
			h := resp.Repository.DefaultBranchRef.Target.History
			for _, n := range h.Nodes {
				t, err := time.Parse(time.RFC3339, n.CommittedDate)
				if err != nil {
					continue
				}
				tl := t.In(loc)
				p.Productive[int(tl.Weekday())][tl.Hour()]++
				seen++
			}
			if !h.PageInfo.HasNextPage {
				break
			}
			end := h.PageInfo.EndCursor
			cursor = &end
		}
	}
	return nil
}
