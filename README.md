# ghstats

> Generate SVG cards summarizing a GitHub user's profile — written in Go.

`ghstats` is a single-binary CLI that fetches public data for a GitHub user and writes a themed set of SVGs (profile details, top languages, stats, productive time) you can embed in your README.

## Status

⚠️ Early work-in-progress. Skeleton only — cards render placeholder SVGs. Roadmap below.

## Install

```sh
go install github.com/tiennm99/ghstats@latest
```

Or build from source:

```sh
git clone https://github.com/tiennm99/ghstats
cd ghstats
go build -o ghstats .
```

## Usage

```sh
export GITHUB_TOKEN=ghp_xxx       # PAT with `repo` + `read:user` for private repo stats
ghstats -user tiennm99 -theme dracula -out output
```

Flags:

| Flag     | Default                 | Description                                      |
| -------- | ----------------------- | ------------------------------------------------ |
| `-user`  | (required)              | GitHub username                                  |
| `-token` | `$GITHUB_TOKEN`         | Personal access token                            |
| `-out`   | `output`                | Output directory (cards land at `<out>/<theme>`) |
| `-theme` | `dracula`               | `dracula`, `default`, `github`                   |

## Output

```
output/
  dracula/
    0-profile-details.svg
    1-languages.svg
    2-stats.svg
    3-productive-time.svg
```

Embed in a README:

```md
![profile](./output/dracula/0-profile-details.svg)
![languages](./output/dracula/1-languages.svg)
```

## Roadmap

- [ ] GitHub GraphQL + REST client (`internal/github`)
  - [ ] Profile basics, followers, repos
  - [ ] Commit histogram for productive time
  - [ ] Language bytes aggregation with `linguist-vendored` respect
  - [ ] Private repo support via PAT
- [ ] Card renderers (`internal/card`)
  - [ ] Profile details
  - [ ] Top languages (by bytes + by commit)
  - [ ] Stats (stars, commits, PRs, issues, contributed-to)
  - [ ] Productive time heatmap
- [ ] Themes (`internal/theme`) — pull the full set from github-readme-stats
- [ ] GitHub Action wrapper for use in profile READMEs
- [ ] Tests + examples

## Credits & inspiration

Standing on the shoulders of these projects:

- [**github-profile-summary-cards**](https://github.com/vn7n24fzkq/github-profile-summary-cards) by [@vn7n24fzkq](https://github.com/vn7n24fzkq) — the card layout, theme set, and output structure are directly inspired by this tool.
- [**profile-summary-for-github**](https://github.com/tipsy/profile-summary-for-github) by [@tipsy](https://github.com/tipsy) — the original web-based profile-summary generator; inspired the breakdowns (repos by language, most-commit language, etc.).

## License

Apache-2.0 — see [LICENSE](LICENSE).
