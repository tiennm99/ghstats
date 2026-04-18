package github

import "time"

// Profile is the aggregate of data other packages render into cards.
type Profile struct {
	ID        string
	Login     string
	Name      string
	Bio       string
	AvatarURL string
	Company   string
	Location  string
	Website   string
	CreatedAt time.Time

	Followers   int
	Following   int
	PublicRepos int

	// Totals for the stats card.
	TotalStars         int
	TotalForks         int
	TotalCommits       int
	TotalPRs           int
	TotalIssues        int
	TotalReviews       int
	TotalContributedTo int
	TotalContributions int // lifetime contributions from calendar + restricted

	// Sorted desc by bytes. Color is GitHub's linguist color or "" if absent.
	Languages []LangStat

	// Commit-count histogram indexed by [day-of-week 0=Sunday][hour-of-day 0-23].
	Productive [7][24]int

	// TopRepos is the list of owned repo names sorted by stargazer count desc,
	// populated by FetchProfile. Used as the seed set for FetchProductive.
	TopRepos []string
}

// LangStat is one row in the top-languages card.
type LangStat struct {
	Name  string
	Color string
	Bytes int64
}

// repoNode is the GraphQL shape of one repository node; kept here because
// it's shared by the profile fetcher and the productive-time fetcher.
type repoNode struct {
	Name            string `json:"name"`
	StargazerCount  int    `json:"stargazerCount"`
	ForkCount       int    `json:"forkCount"`
	PrimaryLanguage *struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"primaryLanguage"`
	Languages struct {
		Edges []struct {
			Size int64 `json:"size"`
			Node struct {
				Name  string `json:"name"`
				Color string `json:"color"`
			} `json:"node"`
		} `json:"edges"`
	} `json:"languages"`
}
