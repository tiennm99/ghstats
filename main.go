// ghstats generates SVG cards summarizing a GitHub user's profile.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tiennm99/ghstats/internal/card"
	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

func main() {
	var (
		user    = flag.String("user", "", "GitHub username (required)")
		token   = flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub token (or env GITHUB_TOKEN)")
		out     = flag.String("out", "output", "output directory for SVG cards")
		themeID = flag.String("theme", "dracula", "theme id (dracula, default, github)")
	)
	flag.Parse()

	if *user == "" {
		fmt.Fprintln(os.Stderr, "error: -user is required")
		flag.Usage()
		os.Exit(2)
	}

	th, ok := theme.Lookup(*themeID)
	if !ok {
		fmt.Fprintf(os.Stderr, "error: unknown theme %q\n", *themeID)
		os.Exit(2)
	}

	client := github.NewClient(*token)
	profile, err := client.Profile(*user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: fetch profile: %v\n", err)
		os.Exit(1)
	}

	if err := card.RenderAll(profile, th, *out); err != nil {
		fmt.Fprintf(os.Stderr, "error: render cards: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote cards to %s/%s/\n", *out, th.ID)
}
