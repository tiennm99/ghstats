package card

import (
	"fmt"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type statsCard struct{}

func (statsCard) Filename() string { return "2-stats.svg" }

func (statsCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	// TODO: totals for stars, commits, PRs, issues, contributed-to repos.
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200">
  <rect width="100%%" height="100%%" fill="%s"/>
  <text x="20" y="40" fill="%s" font-family="sans-serif" font-size="20">Stats</text>
  <text x="20" y="80" fill="%s" font-family="sans-serif" font-size="12">%d public repos · %d followers · %d following</text>
</svg>`, t.Background, t.Title, t.Text, p.PublicRepos, p.Followers, p.Following)
	return []byte(svg), nil
}
