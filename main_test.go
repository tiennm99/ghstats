package main

import (
	"strings"
	"testing"
	"time"
)

// TestUTCOffsetLabel checks that half-hour and quarter-hour zones render
// with a decimal, matching github-profile-summary-cards' "UTC+X.NN" style.
func TestUTCOffsetLabel(t *testing.T) {
	cases := []struct {
		zone string
		want string // must appear in the label; exact value varies by DST
	}{
		{"UTC", "UTC+0.00"},
		{"Asia/Saigon", "UTC+7.00"},
		{"Asia/Kolkata", "UTC+5.50"},  // half-hour zone
		{"Asia/Kathmandu", "UTC+5.75"}, // quarter-hour zone
	}
	for _, tc := range cases {
		loc, err := time.LoadLocation(tc.zone)
		if err != nil {
			t.Skipf("%s unavailable: %v", tc.zone, err)
		}
		got := utcOffsetLabel(loc)
		if !strings.Contains(got, tc.want) {
			t.Errorf("utcOffsetLabel(%q) = %q, want prefix %q", tc.zone, got, tc.want)
		}
	}
}
