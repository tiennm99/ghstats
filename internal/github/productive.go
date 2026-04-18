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

// scaleFactor is the fixed-point multiplier used when distributing a single
// commit across several languages by byte share. Stored in LangStat.Value
// (int64) so the existing sort + percentage math keeps working; the absolute
// magnitude is irrelevant because the card renders percentages.
const scaleFactor = 10_000

// FetchProductive fills p.Productive with a 24-hour commit histogram over the
// last year and p.CommitsByLanguage with commit counts distributed across each
// repo's language byte breakdown. Commits are gathered from the given repos
// (usually p.TopRepos[:N]); each repo is sampled up to maxPerRepo commits to
// keep the cost bounded.
//
// Attribution model: each commit contributes a whole scaleFactor unit,
// partitioned across the repo's languages proportional to linguist byte
// counts. A repo that is 60% Go / 40% Python credits 0.6 to Go and 0.4 to
// Python per commit — a strict upgrade over the previous primary-language-
// only model. Prose languages (Markdown, AsciiDoc, …) remain excluded by
// linguist itself, so blog-style repos still skew toward their detected
// code fraction; fixing that requires per-commit file classification.
//
// The timezone loc is applied to CommittedDate so the heatmap reflects when
// the user actually commits, not UTC.
func (c *Client) FetchProductive(p *Profile, repos []RepoInfo, loc *time.Location, maxPerRepo int) error {
	if loc == nil {
		loc = time.UTC
	}
	since := time.Now().AddDate(-1, 0, 0).UTC().Format(time.RFC3339)

	commitsByLang := map[string]int64{}
	langColor := map[string]string{}

	for _, repo := range repos {
		var cursor *string
		seen := 0
		for {
			if seen >= maxPerRepo {
				break
			}
			vars := map[string]any{
				"login":  p.Login,
				"repo":   repo.Name,
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
				p.Productive[t.In(loc).Hour()]++
				attributeCommit(repo, commitsByLang, langColor)
				seen++
			}
			if !h.PageInfo.HasNextPage {
				break
			}
			end := h.PageInfo.EndCursor
			cursor = &end
		}
	}

	p.CommitsByLanguage = sortLangStats(commitsByLang, langColor)
	return nil
}

// attributeCommit distributes a single commit across the repo's languages
// proportional to byte share. Falls back to the primary language when no
// byte breakdown is available (empty repo or linguist-free repo).
func attributeCommit(repo RepoInfo, commitsByLang map[string]int64, langColor map[string]string) {
	var total int64
	for _, l := range repo.Languages {
		total += l.Bytes
	}
	if total == 0 {
		if repo.PrimaryLanguage != "" {
			commitsByLang[repo.PrimaryLanguage] += scaleFactor
			if _, ok := langColor[repo.PrimaryLanguage]; !ok {
				langColor[repo.PrimaryLanguage] = repo.PrimaryColor
			}
		}
		return
	}
	for _, l := range repo.Languages {
		share := int64(scaleFactor) * l.Bytes / total
		if share == 0 {
			continue
		}
		commitsByLang[l.Name] += share
		if _, ok := langColor[l.Name]; !ok {
			langColor[l.Name] = l.Color
		}
	}
}
