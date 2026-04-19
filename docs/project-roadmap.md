# Project Roadmap

## Phase 0 — Skeleton (✅ done)

- Module layout, flag parsing, placeholder SVG renderers.

## Phase 1 — Five core cards (✅ done)

- Profile details, repos-per-language, most-commit-language, stats, productive-time.
- GraphQL profile query + per-repo commit history.
- Docker-based Action wrapper, release workflow, 65-theme palette.

## Phase 2 — Chart quality (✅ done)

- Match github-profile-summary-cards visual style (donuts, 24h bar chart, proper axes).
- Octicon labels on profile + stats cards.
- Smooth area chart for contributions (Catmull-Rom → cubic Bezier).

## Phase 3 — All-time variants (✅ done)

- Unified commit-history fetch splits into last-year and all-time buckets.
- Per-year `contributionsCollection` loop yields `DailyContributionsAllTime` + `TotalCommitsAllTime`.
- Three new cards: most-commit-language-all-time, productive-time-all-time, contributions-all-time.
- Stats card gains a lifetime commits row.

## Phase 4 — Accurate repo sampling (✅ done)

- Seed list built from `commitContributionsByRepository` across every active year.
- `-include-forks` / `-include-private` visibility flags (defaults later flipped on in Phase 7).
- `-top-repos` demoted to an optional cap (default 0 = unlimited).
- Commit-history query takes `$owner` so forks and non-owned repos are probeable.

## Phase 5 — Byte-weighted attribution (✅ done)

- Each commit distributes fractionally across repo's language bytes, not just primary.
- Improves mixed-code repo accuracy; still inaccurate for Markdown-heavy repos (linguist prose-exclusion).

## Phase 6 — Code-review remediation (✅ done)

Follow-up after the full-project review (`plans/reports/code-review-260418-2223-full-project.md`):

- Donut chart's single-slice (100%) rendering no longer produces an empty arc.
- `FetchContributionsAllTime` warns on stderr when a year returns nil user data.
- `attributeCommit` receives a precomputed per-repo byte total instead of re-summing every commit.
- `Profile.TotalContributions` → `TotalContributionsLastYear` (accurate semantics).
- `context.Context` threaded through all fetchers; `-timeout` flag (default 30m); Ctrl-C cancels in-flight requests.
- Rate-limit awareness: on 429 or exhausted primary limit, honor `Retry-After` / `X-RateLimit-Reset` up to 5 min and retry once.
- Release workflow gates docker + binaries on a test job; no more shipping broken tags.
- Docker base images and third-party GitHub Actions pinned to SHA with version comments.
- Stats card label "Contributed to (non-fork)" corrected to "Contributed to" (the query doesn't filter forks).
- Tests: fixed stale XML-escape assertion, added `TestDonutSingleSlice`, added `TestUTCOffsetLabel` for half-hour zones.

## Phase 7 — Release polish & Marketplace publish (✅ done)

- Visibility defaults flipped on: `-include-forks`, `-include-private` now default `true` (private silently no-ops if token lacks scope).
- Output filenames dropped the numeric prefix: `0-profile-details.svg` → `profile-details.svg` etc. Embedders reference by name.
- Card dimensions shrunk `500×220` → `340×200` to match github-profile-summary-cards so two cards fit per row in a README.
- Action `action.yml` name set to `ghstats-cards` for Marketplace (the bare `ghstats` is taken); repo stays `tiennm99/ghstats`.
- `v1.0.0`, `v1.1.0`, `v1.1.1` tagged and released. Prebuilt binaries (linux/darwin/windows × amd64/arm64) ship with each; Docker image pushed to `ghcr.io/tiennm99/ghstats`.
- Floating `v1` major tag created; `release.yml` has an `update-major-tag` job that force-moves `v1` to the latest patch after test+docker+binaries pass, so consumers pinned to `tiennm99/ghstats@v1` auto-pick new releases.
- README badges (Marketplace / Release / License) + direct Marketplace link for cross-navigation.
- Repo topics expanded for Marketplace discoverability (`ghstats-cards`, `profile-readme`, `stats-cards`, etc.).
- An attempted repo rename to `tiennm99/ghstats-cards` was committed and reverted (commits `399a3dc` + `8bd2128` on record) — GHCR path immutability and the cost of breaking pinned consumers outweighed the Marketplace-name cosmetic benefit.

## Phase 7.6 — S-tier breadth cards (✅ done)

Five new cards that ride on data already fetched — zero extra API calls:

- `contributions-heatmap` — canonical 7×53 calendar grid with a theme-derived 5-bucket intensity ramp.
- `contributions-by-year` — one bar per active year, peak year highlighted.
- `productive-weekday` + `productive-weekday-all-time` — mirror the hour-of-day pair; `FetchProductive` now also fills `Weekday` / `WeekdayAllTime` histograms.
- `top-starred-repos` — top 5 owned non-fork repos by ⭐; required threading `Stars` through `RepoInfo`.
- `streak` — current + longest streak + active days/total. Pure post-processing of `DailyContributionsAllTime`.

Card count: 9 → 14. `FetchProductive` still pays for commit-history pagination once; the new cards are pure renderers.

## Phase 7.5 — Demo gallery for theme discovery (✅ done)

- New `.github/workflows/demo.yml` renders every card for every theme against the repo owner's profile on each push to `main`.
- Output lands in `demo/<theme>/`, with an auto-generated `demo/README.md` TOC so reviewers can browse palettes side-by-side with real data instead of cloning and running the CLI.
- Loop prevention: workflow skips pushes that only touch `demo/**`, `output/**`, `**.md`, or `LICENSE`; `GITHUB_TOKEN`-driven pushes don't retrigger workflows by design.
- Consumer impact: none — this is a repo-internal discovery aid, not a shipped feature.

---

## Phase 8 — Per-commit file classification (planned)

**Goal**: fix the Markdown-blog misattribution case (and any repo where linguist's byte view disagrees with what files user actually edited).

**Approach**: `GET /repos/{owner}/{repo}/commits/{sha}` per commit → classify each file with `go-enry`. Weight by `additions + deletions`.

**Cost**: ~1 REST call per commit. At current defaults (30 seed repos × 500 commits = 15,000 commits worst case) this is heavy — needs `-accurate-languages` opt-in flag, schedule weekly not daily.

**Research**: see `plans/reports/researcher-260418-2001-accurate-language-stats.md`.

**Status**: designed, not implemented.

## Phase 9 — Partial bare clone for lifetime all-repo stats (planned)

**Goal**: lifetime language stats across **every** repo a user has committed in, without the 500-commits-per-repo cap.

**Approach**: `git clone --filter=blob:none --bare` per seed repo + `git log --author --numstat` → go-enry.

**Cost**: ~5% of full-clone disk (trees only, no blobs); 3–5 minutes runtime for 100 repos; zero REST calls.

**Trade-off**: needs disk + git binary on runner. Lowlighter/metrics' indepth mode does similar but clones full blobs; we'd skip those.

**Status**: researched only; behind `-deep` flag when landed.

## Phase 10 — User-configurable repo exclusion (planned)

**Goal**: let users drop throwaway repos (experiments, forks they stashed) from stats without disabling forks globally.

**Approach**: `-exclude-repo owner1/name1,owner2/name2` flag. Filter seed list before probing.

**Cost**: negligible (client-side filter).

**Status**: pending user demand.

## Phase 11 — Expand ownerAffiliations (planned)

**Goal**: catch work done in org repos where user is a collaborator, not owner (e.g., company monorepos).

**Approach**: expose `-affiliations OWNER,COLLABORATOR,ORGANIZATION_MEMBER` flag. Requires thinking about whether to *display* private org work on a public profile card.

**Status**: blocked on deciding the privacy default.

---

## Known limitations (not roadmap items — by design)

| Limitation | Reason |
| --- | --- |
| Markdown/prose excluded from byte counts | Linguist's default; we defer to linguist |
| No real-time API | Scope: scheduled batch renderer, not a server |
| No WakaTime integration | Out of scope — WakaTime cards already exist (athul/waka-readme, anmol098/waka-readme-stats) |
| No heatmap (7×24) variant of productive time | Simplified to 24-hour bar chart to match reference project |
| Hard width of 340 px per card | Matches github-profile-summary-cards; customising would cascade through every chart's geometry. |

## Tracked research reports

All in `plans/reports/`:
- `researcher-260418-2001-accurate-language-stats.md` — metrics vs GRS vs go-enry feasibility
- `researcher-260418-2012-profile-stats-survey.md` — follow-up survey across 6 more tools
- `analysis-260418-2140-most-commit-language-all-time.md` — hand-reconstruction of tiennm99's card output, showing exactly why each language lands where
- `code-review-260418-2223-full-project.md` — adversarial review of the whole codebase; findings all closed in Phase 6
