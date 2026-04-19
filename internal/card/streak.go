package card

import (
	"fmt"
	"strings"
	"time"

	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

type streakCard struct{}

func (streakCard) Filename() string { return "streak.svg" }

func (streakCard) SVG(p *github.Profile, t theme.Theme) ([]byte, error) {
	const (
		width  = 340
		height = 200
	)

	stats := computeStreak(p.DailyContributionsAllTime)

	var b strings.Builder
	b.WriteString(header(width, height, t.Background, t.Stroke, t.StrokeOpacity, t.Title, "Streak"))

	// Three large stat columns (current / longest / active-days) side by side,
	// each with a big number on top and a smaller label underneath. Mirrors
	// the classic "streak" card layout so embedders recognise it instantly.
	cols := []struct {
		value string
		label string
		end   string // optional date annotation beneath the label
	}{
		{formatInt(stats.Current), "Current streak", streakRange(stats.CurrentStart, stats.CurrentEnd)},
		{formatInt(stats.Longest), "Longest streak", streakRange(stats.LongestStart, stats.LongestEnd)},
		{fmt.Sprintf("%d / %d", stats.Active, stats.Total), "Active days", ""},
	}
	colW := width / len(cols)
	for i, c := range cols {
		cx := colW*i + colW/2
		fmt.Fprintf(&b, `
  <text x="%d" y="%d" font-size="28" font-weight="700" fill="%s" text-anchor="middle">%s</text>
  <text x="%d" y="%d" font-size="12" fill="%s" text-anchor="middle">%s</text>`,
			cx, 95, t.Accent, escapeXML(c.value),
			cx, 120, t.Text, escapeXML(c.label))
		if c.end != "" {
			fmt.Fprintf(&b, `
  <text x="%d" y="%d" font-size="10" fill="%s" text-anchor="middle">%s</text>`,
				cx, 140, t.Muted, escapeXML(c.end))
		}
	}

	b.WriteString(footer)
	return []byte(b.String()), nil
}

// streakStats is the post-processed daily series summarised for the card.
type streakStats struct {
	Current              int
	CurrentStart, CurrentEnd time.Time
	Longest              int
	LongestStart, LongestEnd time.Time
	Active               int // days with ≥1 contribution
	Total                int // total days observed
}

// computeStreak walks the daily series once. The "current streak" runs
// backwards from the most recent day; if today has 0 contributions we still
// count yesterday as current (a single-day grace) so the card doesn't reset
// the moment a user hasn't pushed yet today.
func computeStreak(days []github.DailyContribution) streakStats {
	var s streakStats
	if len(days) == 0 {
		return s
	}
	s.Total = len(days)

	// Longest streak + active day count: single forward pass.
	var run int
	var runStart time.Time
	for _, d := range days {
		if d.Count > 0 {
			s.Active++
			if run == 0 {
				runStart = d.Date
			}
			run++
			if run > s.Longest {
				s.Longest = run
				s.LongestStart = runStart
				s.LongestEnd = d.Date
			}
		} else {
			run = 0
		}
	}

	// Current streak: walk backwards from the end. Skip at most one trailing
	// zero-day (today-not-pushed-yet) before aborting.
	tail := len(days) - 1
	if days[tail].Count == 0 && tail > 0 {
		tail--
	}
	for i := tail; i >= 0; i-- {
		if days[i].Count == 0 {
			break
		}
		s.Current++
		s.CurrentEnd = days[tail].Date
		s.CurrentStart = days[i].Date
	}
	return s
}

// streakRange formats the open/close dates of a streak as "Mon 2 — Wed 11"
// when both are present. Returns "" when the streak is zero-length so the
// card renders cleanly.
func streakRange(start, end time.Time) string {
	if start.IsZero() || end.IsZero() {
		return ""
	}
	if start.Equal(end) {
		return start.Format("Jan 2, 2006")
	}
	if start.Year() == end.Year() {
		return start.Format("Jan 2") + " — " + end.Format("Jan 2, 2006")
	}
	return start.Format("Jan 2006") + " — " + end.Format("Jan 2006")
}
