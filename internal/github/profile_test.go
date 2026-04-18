package github

import "testing"

func TestSortLanguages(t *testing.T) {
	bytes := map[string]int64{
		"Go":         500,
		"Python":     300,
		"TypeScript": 500, // tie with Go → alphabetical wins
		"HTML":       100,
	}
	colors := map[string]string{
		"Go":         "#00ADD8",
		"Python":     "#3572A5",
		"TypeScript": "#3178c6",
	}
	got := sortLanguages(bytes, colors)

	wantOrder := []string{"Go", "TypeScript", "Python", "HTML"}
	if len(got) != len(wantOrder) {
		t.Fatalf("len=%d want %d", len(got), len(wantOrder))
	}
	for i, name := range wantOrder {
		if got[i].Name != name {
			t.Errorf("pos %d: %q want %q", i, got[i].Name, name)
		}
	}
	if got[0].Color != "#00ADD8" {
		t.Errorf("Go color=%q want #00ADD8", got[0].Color)
	}
	if got[3].Color != "" {
		t.Errorf("HTML color=%q want empty (missing from colors)", got[3].Color)
	}
}
