# ghstats

Single-binary Go CLI that renders GitHub profile cards as SVG, wrapped by a
published GitHub Action. `main.go` parses flags → `internal/github` fetches →
`internal/card` renders → `internal/theme` supplies palettes.

## Adding or changing a config knob

A new flag, input, or default is not done when the code compiles. Every knob
has a surface in five places, and a change that lands in some but not all of
them ships a half-wired feature or a README that lies:

1. `main.go` — the CLI flag
2. `action.yml` — the matching Action input and its default
3. `entrypoint.sh` — the `INPUT_*` → flag translation
4. `README.md` — **both** reference tables (Action inputs, CLI flags) **and
   the runnable examples**: the workflow YAML under "Use as a GitHub Action"
   and the `ghstats …` command under "Use as a CLI"
5. `docs/system-architecture.md` — when the knob changes the fetch pipeline
   or its query cost

The examples are the step most often missed, and they are what people copy.
Grep the flag name across the repo before calling the change complete; every
hit that is a reference table or example should mention it.

## Behavior changes are user-visible

The Action is published to the Marketplace, so existing callers get whatever
`@v1` points at. Default-off for anything that moves numbers on already-
rendered cards, and say so in the input description. Releases are tag-driven:
push `vX.Y.Z`, and `release.yml` tests, builds, and force-moves the floating
`v1` tag.

## Cards and the demo gallery

Adding a card means updating `internal/card/card.go`, the card table in
`README.md`, and the layout in `.github/workflows/demo.yml` — the demo
generator embeds a hardcoded list of SVGs, so a new card renders into every
theme directory but stays invisible in the gallery until it is added there.
Per-theme pages mirror the table layout in the author's profile README at
`tiennm99/tiennm99`.

## GitHub API constraints worth remembering

- `commitContributionsByRepository` caps at 100 repos per query and truncates
  silently. Contribution years are queried by quarter, and a saturated quarter
  is re-asked by month. Only widen windows with that ceiling in mind.
- `contributionCalendar` is clipped to the exact `from`/`to`, with no
  week-boundary spillover, so disjoint windows never double-count days.
- Test files named `*_windows_test.go` (or any GOOS/GOARCH suffix) are
  silently excluded from the build on other platforms — `go test` reports
  "ok" while running nothing.
