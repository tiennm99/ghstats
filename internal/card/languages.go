package card

import (
	"fmt"
	"strings"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type languagesCard struct{}

func (languagesCard) Filename() string { return "1-languages.svg" }

func (languagesCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	const (
		width    = 500
		height   = 220
		topN     = 6
		barX     = 25
		barY     = 60
		barW     = 450
		barH     = 10
		legendX0 = 25
	)

	var b strings.Builder
	b.WriteString(header(width, height, t.Background, t.Title, "Top Languages"))

	langs := p.Languages
	if len(langs) > topN {
		langs = langs[:topN]
	}

	if len(langs) == 0 {
		fmt.Fprintf(&b, `
  <text x="25" y="90" font-size="13" fill="%s">No language data available.</text>`, t.Muted)
		b.WriteString(footer)
		return []byte(b.String()), nil
	}

	var total int64
	for _, l := range langs {
		total += l.Bytes
	}

	// Stacked bar.
	fmt.Fprintf(&b, `
  <rect x="%d" y="%d" width="%d" height="%d" rx="5" fill="%s"/>
  <g>`,
		barX, barY, barW, barH, t.Muted)

	offset := float64(barX)
	for _, l := range langs {
		w := float64(barW) * float64(l.Bytes) / float64(total)
		fmt.Fprintf(&b, `
    <rect x="%.2f" y="%d" width="%.2f" height="%d" fill="%s"/>`,
			offset, barY, w, barH, colorOrAccent(l.Color, t.Accent))
		offset += w
	}
	b.WriteString(`
  </g>`)

	// Legend: two columns of up to 3 rows.
	for i, l := range langs {
		col := i % 2
		row := i / 2
		x := legendX0 + col*230
		y := 110 + row*24
		pct := 100 * float64(l.Bytes) / float64(total)
		fmt.Fprintf(&b, `
  <circle cx="%d" cy="%d" r="6" fill="%s"/>
  <text x="%d" y="%d" font-size="13" fill="%s">%s %.2f%%</text>`,
			x+6, y-4, colorOrAccent(l.Color, t.Accent),
			x+20, y, t.Text, escapeXML(l.Name), pct)
	}

	b.WriteString(footer)
	return []byte(b.String()), nil
}

func colorOrAccent(c, fallback string) string {
	if c == "" {
		return fallback
	}
	return c
}
