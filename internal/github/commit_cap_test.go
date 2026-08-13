package github

import "testing"

func TestReachedCommitCap(t *testing.T) {
	cases := []struct {
		name       string
		seen       int
		maxPerRepo int
		want       bool
	}{
		{"under the cap", 100, 500, false},
		{"at the cap", 500, 500, true},
		{"past the cap", 700, 500, true},
		// Zero means unlimited, so even a fresh repo keeps paginating. The
		// old behavior stopped here and rendered empty commit-derived cards.
		{"zero cap, nothing seen yet", 0, 0, false},
		{"zero cap, deep into history", 100_000, 0, false},
		{"negative cap treated as unlimited", 10, -1, false},
	}
	for _, tc := range cases {
		if got := reachedCommitCap(tc.seen, tc.maxPerRepo); got != tc.want {
			t.Errorf("%s: reachedCommitCap(%d, %d) = %v, want %v",
				tc.name, tc.seen, tc.maxPerRepo, got, tc.want)
		}
	}
}
