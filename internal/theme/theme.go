// Package theme defines SVG color palettes used by card renderers.
package theme

import "sort"

// Theme describes the colors applied to a rendered card.
type Theme struct {
	ID         string
	Background string
	Text       string
	Title      string
	Accent     string
	Muted      string
}

// Built-in palettes. Port of a curated subset from github-readme-stats.
var themes = map[string]Theme{
	"default":         {ID: "default", Background: "#fffefe", Text: "#434d58", Title: "#2f80ed", Accent: "#4c71f2", Muted: "#6a737d"},
	"dark":            {ID: "dark", Background: "#151515", Text: "#9f9f9f", Title: "#fff", Accent: "#79ff97", Muted: "#666"},
	"radical":         {ID: "radical", Background: "#141321", Text: "#a9fef7", Title: "#fe428e", Accent: "#f8d847", Muted: "#a9fef7"},
	"merko":           {ID: "merko", Background: "#0b1p08", Text: "#68b684", Title: "#abd200", Accent: "#b7d364", Muted: "#68b684"},
	"gruvbox":         {ID: "gruvbox", Background: "#282828", Text: "#fbf1c7", Title: "#fabd2f", Accent: "#8ec07c", Muted: "#a89984"},
	"tokyonight":      {ID: "tokyonight", Background: "#1a1b27", Text: "#a9b1d6", Title: "#70a5fd", Accent: "#bf91f3", Muted: "#565f89"},
	"onedark":         {ID: "onedark", Background: "#282c34", Text: "#aaa", Title: "#e4bf7a", Accent: "#8eb573", Muted: "#5c6370"},
	"cobalt":          {ID: "cobalt", Background: "#193549", Text: "#dbe4ee", Title: "#e683d9", Accent: "#0088ff", Muted: "#6fc3df"},
	"synthwave":       {ID: "synthwave", Background: "#2b213a", Text: "#e5289e", Title: "#e2e9ec", Accent: "#ef8539", Muted: "#e5289e"},
	"highcontrast":    {ID: "highcontrast", Background: "#000000", Text: "#ffffff", Title: "#e7f216", Accent: "#00ffff", Muted: "#ffffff"},
	"dracula":         {ID: "dracula", Background: "#282a36", Text: "#f8f8f2", Title: "#ff79c6", Accent: "#bd93f9", Muted: "#6272a4"},
	"prussian":        {ID: "prussian", Background: "#172f45", Text: "#c8c9db", Title: "#bddfff", Accent: "#38b2ac", Muted: "#6c95b8"},
	"monokai":         {ID: "monokai", Background: "#272822", Text: "#d6ebbf", Title: "#eb1f6a", Accent: "#e28905", Muted: "#75715e"},
	"vue":             {ID: "vue", Background: "#fffefe", Text: "#476582", Title: "#41b883", Accent: "#35495e", Muted: "#476582"},
	"vue-dark":        {ID: "vue-dark", Background: "#1d1f21", Text: "#bbb", Title: "#41b883", Accent: "#41b883", Muted: "#888"},
	"shades-of-purple":{ID: "shades-of-purple", Background: "#2d2b55", Text: "#a599e9", Title: "#fad000", Accent: "#b362ff", Muted: "#a599e9"},
	"nightowl":        {ID: "nightowl", Background: "#011627", Text: "#acb4c2", Title: "#7fdbca", Accent: "#82aaff", Muted: "#637777"},
	"buefy":           {ID: "buefy", Background: "#ffffff", Text: "#363636", Title: "#7957d5", Accent: "#ff3860", Muted: "#7a7a7a"},
	"blue-green":      {ID: "blue-green", Background: "#040f0f", Text: "#2dd4bf", Title: "#afebcd", Accent: "#26a69a", Muted: "#5e8b7e"},
	"algolia":         {ID: "algolia", Background: "#050f2c", Text: "#ffffff", Title: "#00aeff", Accent: "#2dde98", Muted: "#8c8c8c"},
	"great-gatsby":    {ID: "great-gatsby", Background: "#000000", Text: "#ffd700", Title: "#ffa726", Accent: "#ffb74d", Muted: "#9e9e9e"},
	"darcula":         {ID: "darcula", Background: "#242424", Text: "#ba5f17", Title: "#ba5f17", Accent: "#2f81f7", Muted: "#8b949e"},
	"bear":            {ID: "bear", Background: "#1f2023", Text: "#8f9396", Title: "#e03c8a", Accent: "#00aeff", Muted: "#8f9396"},
	"solarized-dark":  {ID: "solarized-dark", Background: "#002b36", Text: "#859900", Title: "#268bd2", Accent: "#d33682", Muted: "#586e75"},
	"solarized-light": {ID: "solarized-light", Background: "#fdf6e3", Text: "#657b83", Title: "#268bd2", Accent: "#d33682", Muted: "#93a1a1"},
	"chartreuse-dark": {ID: "chartreuse-dark", Background: "#000000", Text: "#ffffff", Title: "#7fff00", Accent: "#7fff00", Muted: "#5fcf00"},
	"nord":            {ID: "nord", Background: "#2e3440", Text: "#d8dee9", Title: "#88c0d0", Accent: "#81a1c1", Muted: "#4c566a"},
	"github":          {ID: "github", Background: "#ffffff", Text: "#24292f", Title: "#0969da", Accent: "#2188ff", Muted: "#57606a"},
	"github-dark":     {ID: "github-dark", Background: "#0d1117", Text: "#c9d1d9", Title: "#58a6ff", Accent: "#3fb950", Muted: "#8b949e"},
	"transparent":     {ID: "transparent", Background: "#00000000", Text: "#434d58", Title: "#2f80ed", Accent: "#4c71f2", Muted: "#6a737d"},
}

// merko had a typo fixed at init.
func init() {
	m := themes["merko"]
	m.Background = "#0b1708"
	themes["merko"] = m
}

// Lookup returns the theme with the given id.
func Lookup(id string) (Theme, bool) {
	t, ok := themes[id]
	return t, ok
}

// IDs returns every registered theme id sorted alphabetically.
func IDs() []string {
	out := make([]string, 0, len(themes))
	for id := range themes {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}
