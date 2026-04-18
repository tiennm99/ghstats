// Package theme defines SVG color palettes used by card renderers.
package theme

// Theme describes the colors applied to a rendered card.
type Theme struct {
	ID         string
	Background string
	Text       string
	Title      string
	Accent     string
	Muted      string
}

// Built-in palettes. Add new themes by appending to the map.
var themes = map[string]Theme{
	"dracula": {
		ID:         "dracula",
		Background: "#282a36",
		Text:       "#f8f8f2",
		Title:      "#ff79c6",
		Accent:     "#bd93f9",
		Muted:      "#6272a4",
	},
	"default": {
		ID:         "default",
		Background: "#ffffff",
		Text:       "#24292f",
		Title:      "#0969da",
		Accent:     "#2188ff",
		Muted:      "#57606a",
	},
	"github": {
		ID:         "github",
		Background: "#0d1117",
		Text:       "#c9d1d9",
		Title:      "#58a6ff",
		Accent:     "#3fb950",
		Muted:      "#8b949e",
	},
}

// Lookup returns the theme with the given id.
func Lookup(id string) (Theme, bool) {
	t, ok := themes[id]
	return t, ok
}
