package card

import (
	"fmt"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type profileCard struct{}

func (profileCard) Filename() string { return "0-profile-details.svg" }

func (profileCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	// TODO: render a real profile details card.
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="200">
  <rect width="100%%" height="100%%" fill="%s"/>
  <text x="20" y="40" fill="%s" font-family="sans-serif" font-size="24">%s</text>
  <text x="20" y="72" fill="%s" font-family="sans-serif" font-size="14">%s</text>
</svg>`, t.Background, t.Title, p.Login, t.Muted, p.Bio)
	return []byte(svg), nil
}
