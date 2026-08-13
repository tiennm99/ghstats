package github

import (
	"testing"
	"time"
)

func TestContributionWindowsFullYear(t *testing.T) {
	now := time.Date(2026, 8, 13, 0, 0, 0, 0, time.UTC)
	got := contributionWindows(2024, now)
	if len(got) != 4 {
		t.Fatalf("want 4 quarters for a completed year, got %d", len(got))
	}

	wantStarts := []string{"2024-01-01", "2024-04-01", "2024-07-01", "2024-10-01"}
	wantEnds := []string{"2024-03-31", "2024-06-30", "2024-09-30", "2024-12-31"}
	for i, w := range got {
		if s := w[0].Format("2006-01-02"); s != wantStarts[i] {
			t.Errorf("quarter %d start = %s, want %s", i, s, wantStarts[i])
		}
		if e := w[1].Format("2006-01-02"); e != wantEnds[i] {
			t.Errorf("quarter %d end = %s, want %s", i, e, wantEnds[i])
		}
	}

	// Windows must not overlap, or days would be counted twice in the
	// all-time contribution series.
	for i := 1; i < len(got); i++ {
		if !got[i][0].After(got[i-1][1]) {
			t.Errorf("quarter %d starts at/before the end of quarter %d", i, i-1)
		}
	}
}

func TestContributionWindowsCurrentYearStopsAtNow(t *testing.T) {
	now := time.Date(2026, 8, 13, 10, 30, 0, 0, time.UTC)
	got := contributionWindows(2026, now)
	if len(got) != 3 {
		t.Fatalf("want 3 quarters through August, got %d", len(got))
	}
	if last := got[2][1]; !last.Equal(now) {
		t.Errorf("final window end = %s, want clamped to now (%s)", last, now)
	}
}

func TestContributionWindowsFutureYearIsEmpty(t *testing.T) {
	now := time.Date(2026, 8, 13, 0, 0, 0, 0, time.UTC)
	if got := contributionWindows(2027, now); len(got) != 0 {
		t.Fatalf("want no windows for a future year, got %d", len(got))
	}
}

func TestRepoAffiliationsAndOwnership(t *testing.T) {
	orgAdmin := repoNode{Name: "chambai", ViewerPermission: "ADMIN"}
	orgAdmin.Owner = &struct {
		Login string `json:"login"`
	}{Login: "miti99dev"}

	orgWrite := orgAdmin
	orgWrite.ViewerPermission = "WRITE"

	own := repoNode{Name: "ghstats"}
	own.Owner = &struct {
		Login string `json:"login"`
	}{Login: "tiennm99"}

	off := FetchOptions{}
	on := FetchOptions{IncludeOrgRepos: true}

	if got := repoAffiliations(off); len(got) != 1 || got[0] != "OWNER" {
		t.Errorf("affiliations with org repos off = %v, want [OWNER]", got)
	}
	if got := repoAffiliations(on); len(got) != 2 {
		t.Errorf("affiliations with org repos on = %v, want OWNER + ORGANIZATION_MEMBER", got)
	}

	cases := []struct {
		name string
		node repoNode
		opts FetchOptions
		want bool
	}{
		{"own repo, org off", own, off, true},
		{"own repo, org on", own, on, true},
		{"org admin repo, org off", orgAdmin, off, false},
		{"org admin repo, org on", orgAdmin, on, true},
		{"org write-only repo, org on", orgWrite, on, false},
	}
	for _, tc := range cases {
		if got := ownedByUser(tc.node, "tiennm99", tc.opts); got != tc.want {
			t.Errorf("%s: ownedByUser = %v, want %v", tc.name, got, tc.want)
		}
	}
}
