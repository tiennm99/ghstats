package card

import (
	"fmt"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type productiveCard struct{}

func (productiveCard) Filename() string { return "3-productive-time.svg" }

func (productiveCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	// TODO: render the [7][24]int heatmap from p.Productive.
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200">
  <rect width="100%%" height="100%%" fill="%s"/>
  <text x="20" y="40" fill="%s" font-family="sans-serif" font-size="20">Productive Time</text>
  <text x="20" y="80" fill="%s" font-family="sans-serif" font-size="12">Heatmap placeholder</text>
</svg>`, t.Background, t.Title, t.Text)
	_ = p
	return []byte(svg), nil
}
