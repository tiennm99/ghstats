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

// maxRepositoriesPerWindow mirrors the maxRepositories argument in
// contributionYearQuery. GitHub caps it at 100, and a window that returns
// exactly that many repos is almost certainly truncated.
const maxRepositoriesPerWindow = 100

// contributionWindows splits one calendar year into quarters, dropping any
// quarter that starts after now and clamping the last one to now.
//
// A year-wide window silently loses repos once the user commits in more than
// maxRepositoriesPerWindow of them in that year — the query returns the top
// 100 and says nothing about the rest. Quarters bring most years under the
// ceiling; a quarter that still saturates gets split again by month. The API
// clips contributionCalendar to the exact window (no week-boundary
// spillover), so concatenating quarters yields each day once.
func contributionWindows(year int, now time.Time) [][2]time.Time {
	var out [][2]time.Time
	for q := 0; q < 4; q++ {
		from := time.Date(year, time.Month(q*3+1), 1, 0, 0, 0, 0, time.UTC)
		if from.After(now) {
			break
		}
		to := from.AddDate(0, 3, 0).Add(-time.Second)
		if to.After(now) {
			to = now
		}
		out = append(out, [2]time.Time{from, to})
	}
	return out
}

// FetchContributionsAllTime iterates p.ContributionYears and issues one
// contributionsCollection query per quarter. Each window's payload
// contributes:
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
	now := time.Now().UTC()

	for _, y := range years {
		for _, w := range contributionWindows(y, now) {
			saturated, err := c.fetchContributionWindow(ctx, p, opts, w[0], w[1], seen, true)
			if err != nil {
				return err
			}
			if !saturated {
				continue
			}
			// The quarter came back at the ceiling, so its repo list is
			// truncated. Re-ask month by month to recover the tail. Only the
			// repo lists are merged — the quarter already contributed its
			// days and commit totals, and folding them again would
			// double-count.
			for _, m := range monthWindows(w[0], w[1]) {
				if _, err := c.fetchContributionWindow(ctx, p, opts, m[0], m[1], seen, false); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// monthWindows splits an arbitrary window into calendar months, clamped to the
// window's own bounds so the pieces cover exactly the same span.
func monthWindows(from, to time.Time) [][2]time.Time {
	var out [][2]time.Time
	start := from
	for start.Before(to) || start.Equal(to) {
		next := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0)
		end := next.Add(-time.Second)
		if end.After(to) {
			end = to
		}
		out = append(out, [2]time.Time{start, end})
		start = next
	}
	return out
}

// fetchContributionWindow folds a single contributionsCollection window into
// the profile and reports whether its repo list came back at the ceiling (and
// so is truncated). foldCalendar is false for the month re-queries that follow
// a saturated quarter: those exist only to recover repos, and re-folding their
// days would count the same contributions twice.
func (c *Client) fetchContributionWindow(
	ctx context.Context,
	p *Profile,
	opts FetchOptions,
	from, to time.Time,
	seen map[string]int,
	foldCalendar bool,
) (bool, error) {
	vars := map[string]any{
		"login": p.Login,
		"from":  from.Format(time.RFC3339),
		"to":    to.Format(time.RFC3339),
	}
	var resp contributionYearGQL
	if err := c.query(ctx, contributionYearQuery, vars, &resp); err != nil {
		return false, err
	}
	if resp.User == nil {
		// Don't abort the run — other windows may still yield data — but
		// make the partial-data case visible instead of rendering an
		// empty all-time card silently.
		fmt.Fprintf(os.Stderr, "warn: contributions %s..%s returned no user data\n",
			from.Format("2006-01-02"), to.Format("2006-01-02"))
		return false, nil
	}

	cc := resp.User.ContributionsCollection
	if foldCalendar {
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
	}

	saturated := len(cc.CommitContributionsByRepository) >= maxRepositoriesPerWindow
	if saturated && !foldCalendar {
		// A month at the ceiling has nowhere finer to go, so this is as
		// complete as the API will answer. Say so instead of implying the
		// list covers everything.
		fmt.Fprintf(os.Stderr,
			"warn: contributions %s..%s hit the %d-repo ceiling even at month granularity; some repos are missing from commit probing\n",
			from.Format("2006-01-02"), to.Format("2006-01-02"), maxRepositoriesPerWindow)
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
	return saturated, nil
}
