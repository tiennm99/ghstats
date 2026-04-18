package card

import (
	"fmt"
	"strings"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type productiveCard struct{}

func (productiveCard) Filename() string { return "4-productive-time.svg" }

var weekdayLabels = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

func (productiveCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	const (
		width    = 650
		height   = 240
		cellSize = 18
		cellGap  = 3
		gridX    = 55
		gridY    = 60
	)

	var b strings.Builder
	b.WriteString(header(width, height, t.Background, t.Title, "Productive Time (last year, by hour)"))

	max := 0
	for _, row := range p.Productive {
		for _, v := range row {
			if v > max {
				max = v
			}
		}
	}

	// Weekday labels.
	for i, d := range weekdayLabels {
		y := gridY + i*(cellSize+cellGap) + cellSize - 4
		fmt.Fprintf(&b, `
  <text x="25" y="%d" font-size="11" fill="%s">%s</text>`,
			y, t.Muted, d)
	}

	// Hour labels along top (every 3 hours).
	for h := 0; h < 24; h += 3 {
		x := gridX + h*(cellSize+cellGap)
		fmt.Fprintf(&b, `
  <text x="%d" y="55" font-size="10" fill="%s">%02dh</text>`,
			x, t.Muted, h)
	}

	// Cells.
	for d := 0; d < 7; d++ {
		for h := 0; h < 24; h++ {
			count := p.Productive[d][h]
			opacity := heatOpacity(count, max)
			x := gridX + h*(cellSize+cellGap)
			y := gridY + d*(cellSize+cellGap)
			fmt.Fprintf(&b, `
  <rect x="%d" y="%d" width="%d" height="%d" rx="3" fill="%s" fill-opacity="%.2f"><title>%s %02d:00 — %d commits</title></rect>`,
				x, y, cellSize, cellSize, t.Accent, opacity,
				weekdayLabels[d], h, count)
		}
	}

	b.WriteString(footer)
	return []byte(b.String()), nil
}

// heatOpacity returns the fill-opacity for a cell. Zero is almost transparent
// so the grid is still visible; max-count maps to fully opaque.
func heatOpacity(count, max int) float64 {
	if max == 0 {
		return 0.08
	}
	const floor = 0.10
	ratio := float64(count) / float64(max)
	return floor + (1.0-floor)*ratio
}
