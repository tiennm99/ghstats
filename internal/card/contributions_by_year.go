package card

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type contributionsByYearCard struct{}

func (contributionsByYearCard) Filename() string { return "contributions-by-year.svg" }

func (contributionsByYearCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	const (
		width    = 340
		height   = 200
		leftAxis = 35
		rightPad = 15
		topPad   = 45
		chartH   = 110
		barGap   = 2
	)

	buckets := aggregateByYear(p.DailyContributionsAllTime)

	var b strings.Builder
	b.WriteString(header(width, height, t.Background, t.Stroke, t.StrokeOpacity, t.Title, "Contributions by Year"))

	if len(buckets) == 0 {
		fmt.Fprintf(&b, `
  <text x="25" y="100" font-size="13" fill="%s">No contribution data available.</text>`, t.Muted)
		b.WriteString(footer)
		return []byte(b.String()), nil
	}

	chartW := width - leftAxis - rightPad
	barW := float64(chartW-barGap*(len(buckets)-1)) / float64(len(buckets))

	var maxVal int
	peakIdx := 0
	for i, bk := range buckets {
		if bk.Count > maxVal {
			maxVal = bk.Count
			peakIdx = i
		}
	}
	yMax := float64(maxVal)
	if yMax == 0 {
		yMax = 1
	}
	ticks := niceTicks(yMax, 5)
	if len(ticks) > 0 {
		yMax = ticks[len(ticks)-1]
	}

	// Y axis + ticks.
	fmt.Fprintf(&b, `
  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s"/>`,
		leftAxis, topPad, leftAxis, topPad+chartH, t.Muted)
	for _, v := range ticks {
		y := topPad + chartH - int(float64(chartH)*v/yMax)
		fmt.Fprintf(&b, `
  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s"/>
  <text x="%d" y="%d" font-size="10" fill="%s" text-anchor="end">%s</text>`,
			leftAxis-4, y, leftAxis, y, t.Muted,
			leftAxis-6, y+3, t.Muted, escapeXML(formatTick(v)))
	}

	// X axis baseline.
	fmt.Fprintf(&b, `
  <line x1="%d" y1="%d" x2="%d" y2="%d" stroke="%s"/>`,
		leftAxis, topPad+chartH, leftAxis+chartW, topPad+chartH, t.Muted)

	// Bars + year labels. Peak year uses Accent; others use a muted Accent
	// mix so the eye snaps to the year that matters most.
	dim := mixHex(t.Background, t.Accent, 0.55)
	labelStride := yearLabelStride(len(buckets))
	for i, bk := range buckets {
		barH := float64(chartH) * float64(bk.Count) / yMax
		x := float64(leftAxis) + (barW+float64(barGap))*float64(i)
		y := float64(topPad+chartH) - barH
		fill := dim
		if i == peakIdx {
			fill = t.Accent
		}
		fmt.Fprintf(&b, `
  <rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" rx="2" fill="%s"><title>%d — %d commits</title></rect>`,
			x, y, barW, barH, fill, bk.Year, bk.Count)

		if i%labelStride == 0 || i == len(buckets)-1 {
			cx := x + barW/2
			fmt.Fprintf(&b, `
  <text x="%.2f" y="%d" font-size="10" fill="%s" text-anchor="middle">%d</text>`,
				cx, topPad+chartH+14, t.Muted, bk.Year)
		}
	}

	b.WriteString(footer)
	return []byte(b.String()), nil
}

// yearBucket holds a calendar-year aggregate count.
type yearBucket struct {
	Year  int
	Count int
}

// aggregateByYear bins the daily series into year totals, ascending. Missing
// years between first and last (rare but possible) become zero rows so the
// x-axis stays chronologically continuous.
func aggregateByYear(days []github.DailyContribution) []yearBucket {
	if len(days) == 0 {
		return nil
	}
	counts := map[int]int{}
	var minY, maxY int
	minY, maxY = days[0].Date.Year(), days[0].Date.Year()
	for _, d := range days {
		y := d.Date.Year()
		counts[y] += d.Count
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	out := make([]yearBucket, 0, maxY-minY+1)
	for y := minY; y <= maxY; y++ {
		out = append(out, yearBucket{Year: y, Count: counts[y]})
	}
	// Defensive sort in case caller ever passes an unordered slice.
	sort.Slice(out, func(i, j int) bool { return out[i].Year < out[j].Year })
	return out
}

// yearLabelStride picks how many years between printed x-axis labels so the
// axis stays legible when the user has a long GitHub history.
func yearLabelStride(n int) int {
	switch {
	case n <= 8:
		return 1
	case n <= 16:
		return 2
	default:
		return 3
	}
}
