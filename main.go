// ghstats generates SVG cards summarizing a GitHub user's profile.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tiennm99/ghstats/internal/card"
	"github.com/tiennm99/ghstats/internal/github"
	"github.com/tiennm99/ghstats/internal/theme"
)

func main() {
	var (
		user       = flag.String("user", "", "GitHub username (required)")
		token      = flag.String("token", os.Getenv("GITHUB_TOKEN"), "GitHub token (or env GITHUB_TOKEN)")
		out        = flag.String("out", "output", "output directory")
		themesFlag = flag.String("themes", "dracula", "comma-separated theme ids, or 'all'")
		tzName     = flag.String("tz", "Local", "timezone for productive-time card (IANA name, e.g. Asia/Saigon)")
		topRepos   = flag.Int("top-repos", 10, "owned repos to sample for productive-time heatmap (0 to skip)")
		perRepo    = flag.Int("commits-per-repo", 100, "max commits sampled per repo")
		listThemes = flag.Bool("list-themes", false, "print available theme ids and exit")
	)
	flag.Parse()

	if *listThemes {
		for _, id := range theme.IDs() {
			fmt.Println(id)
		}
		return
	}

	if *user == "" {
		fmt.Fprintln(os.Stderr, "error: -user is required")
		flag.Usage()
		os.Exit(2)
	}

	selected, err := resolveThemes(*themesFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	loc, err := time.LoadLocation(*tzName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: unknown timezone %q, falling back to UTC\n", *tzName)
		loc = time.UTC
	}

	client := github.NewClient(*token)
	profile, err := client.FetchProfile(*user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: fetch profile: %v\n", err)
		os.Exit(1)
	}

	if *topRepos > 0 && profile.ID != "" {
		repos := profile.TopRepos
		if len(repos) > *topRepos {
			repos = repos[:*topRepos]
		}
		if err := client.FetchProductive(profile, repos, loc, *perRepo); err != nil {
			fmt.Fprintf(os.Stderr, "warn: productive-time + commits-per-language fetch: %v\n", err)
		}
	}

	for _, t := range selected {
		if err := card.RenderAll(profile, t, *out); err != nil {
			fmt.Fprintf(os.Stderr, "error: render %s: %v\n", t.ID, err)
			os.Exit(1)
		}
		fmt.Printf("wrote %s/%s/\n", *out, t.ID)
	}
}

func resolveThemes(spec string) ([]theme.Theme, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, fmt.Errorf("no themes specified")
	}
	if spec == "all" {
		ids := theme.IDs()
		out := make([]theme.Theme, 0, len(ids))
		for _, id := range ids {
			if t, ok := theme.Lookup(id); ok {
				out = append(out, t)
			}
		}
		return out, nil
	}
	var out []theme.Theme
	for _, id := range strings.Split(spec, ",") {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}
		t, ok := theme.Lookup(id)
		if !ok {
			return nil, fmt.Errorf("unknown theme %q (use -list-themes)", id)
		}
		out = append(out, t)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no valid themes")
	}
	return out, nil
}
