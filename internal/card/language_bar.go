package card

import (
	"fmt"
	"strings"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

// renderLanguageCard draws a horizontal stacked bar + legend from a list of
// LangStats. Shared by the repos-per-language and most-commit-language cards.
//
// title is the card heading; empty is rendered as the "no data" fallback.
func renderLanguageCard(title string, stats []github.LangStat, t theme.Theme) []byte {
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
	b.WriteString(header(width, height, t.Background, t.Title, title))

	if len(stats) > topN {
		stats = stats[:topN]
	}

	if len(stats) == 0 {
		fmt.Fprintf(&b, `
  <text x="25" y="90" font-size="13" fill="%s">No data available.</text>`, t.Muted)
		b.WriteString(footer)
		return []byte(b.String())
	}

	var total int64
	for _, s := range stats {
		total += s.Value
	}

	fmt.Fprintf(&b, `
  <rect x="%d" y="%d" width="%d" height="%d" rx="5" fill="%s"/>
  <g>`,
		barX, barY, barW, barH, t.Muted)

	offset := float64(barX)
	for _, s := range stats {
		w := float64(barW) * float64(s.Value) / float64(total)
		fmt.Fprintf(&b, `
    <rect x="%.2f" y="%d" width="%.2f" height="%d" fill="%s"/>`,
			offset, barY, w, barH, colorOrAccent(s.Color, t.Accent))
		offset += w
	}
	b.WriteString(`
  </g>`)

	for i, s := range stats {
		col := i % 2
		row := i / 2
		x := legendX0 + col*230
		y := 110 + row*24
		pct := 100 * float64(s.Value) / float64(total)
		fmt.Fprintf(&b, `
  <circle cx="%d" cy="%d" r="6" fill="%s"/>
  <text x="%d" y="%d" font-size="13" fill="%s">%s %.2f%%</text>`,
			x+6, y-4, colorOrAccent(s.Color, t.Accent),
			x+20, y, t.Text, escapeXML(s.Name), pct)
	}

	b.WriteString(footer)
	return []byte(b.String())
}

func colorOrAccent(c, fallback string) string {
	if c == "" {
		return fallback
	}
	return c
}
