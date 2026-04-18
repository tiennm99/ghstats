package github

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"
)

// contributionYearGQL mirrors contributionYearQuery.
type contributionYearGQL struct {
	User *struct {
		ContributionsCollection struct {
			TotalCommitContributions int `json:"totalCommitContributions"`
			ContributionCalendar     struct {
				Weeks []struct {
					ContributionDays []struct {
						ContributionCount int    `json:"contributionCount"`
						Date              string `json:"date"`
					} `json:"contributionDays"`
				} `json:"weeks"`
			} `json:"contributionCalendar"`
			CommitContributionsByRepository []struct {
				Contributions struct {
					TotalCount int `json:"totalCount"`
				} `json:"contributions"`
				Repository repoNode `json:"repository"`
			} `json:"commitContributionsByRepository"`
		} `json:"contributionsCollection"`
	} `json:"user"`
}

// FetchContributionsAllTime iterates p.ContributionYears and issues one
// contributionsCollection query per year. Each year's payload contributes:
//
//   - Days → p.DailyContributionsAllTime
//   - Commit count → p.TotalCommitsAllTime
//   - Repos the user committed in → p.SeedRepos (deduplicated by owner/name)
//
// Fork and private repos are filtered client-side per opts so the caller can
// run the same pipeline with different visibility policies.
func (c *Client) FetchContributionsAllTime(ctx context.Context, p *Profile, opts FetchOptions) error {
	years := append([]int(nil), p.ContributionYears...)
	sort.Ints(years) // ascending so the concatenated series is oldest→newest

	seen := map[string]int{} // "owner/name" → index in p.SeedRepos

	for _, y := range years {
		from := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(y, 12, 31, 23, 59, 59, 0, time.UTC)
		if now := time.Now().UTC(); to.After(now) {
			to = now
		}

		vars := map[string]any{
			"login": p.Login,
			"from":  from.Format(time.RFC3339),
			"to":    to.Format(time.RFC3339),
		}
		var resp contributionYearGQL
		if err := c.query(ctx, contributionYearQuery, vars, &resp); err != nil {
			return err
		}
		if resp.User == nil {
			// Don't abort the run — other years may still yield data — but
			// make the partial-data case visible instead of rendering an
			// empty all-time card silently.
			fmt.Fprintf(os.Stderr, "warn: contribution year %d returned no user data\n", y)
			continue
		}

		cc := resp.User.ContributionsCollection
		p.TotalCommitsAllTime += cc.TotalCommitContributions
		for _, w := range cc.ContributionCalendar.Weeks {
			for _, d := range w.ContributionDays {
				t, err := time.Parse("2006-01-02", d.Date)
				if err != nil {
					continue
				}
				p.DailyContributionsAllTime = append(p.DailyContributionsAllTime, DailyContribution{
					Date:  t,
					Count: d.ContributionCount,
				})
			}
		}

		for _, cr := range cc.CommitContributionsByRepository {
			r := cr.Repository
			if r.IsFork && !opts.IncludeForks {
				continue
			}
			if r.IsPrivate && !opts.IncludePrivate {
				continue
			}
			info := r.toRepoInfo(p.Login)
			key := info.Owner + "/" + info.Name
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = len(p.SeedRepos)
			p.SeedRepos = append(p.SeedRepos, info)
		}
	}
	return nil
}
