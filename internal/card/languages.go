package card

import (
	"fmt"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type languagesCard struct{}

func (languagesCard) Filename() string { return "1-languages.svg" }

func (languagesCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	// TODO: render language breakdown from p.Languages.
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200">
  <rect width="100%%" height="100%%" fill="%s"/>
  <text x="20" y="40" fill="%s" font-family="sans-serif" font-size="20">Top Languages</text>
  <text x="20" y="80" fill="%s" font-family="sans-serif" font-size="12">%d languages tracked</text>
</svg>`, t.Background, t.Title, t.Text, len(p.Languages))
	return []byte(svg), nil
}
