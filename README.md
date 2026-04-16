# git-remote-commits

`git-remote-commits` is a terminal UI that monitors Git activity in real time, like an `htop`-style dashboard for repository commits.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).

## Features

- Live commit list with date, author, and commit message.
- Auto-refresh polling (default every 3 seconds).
- New commit highlighting between refreshes.
- Live remote tracking for the current branch upstream (for example `origin/feature-x`), not a fixed branch.
- Upstream status shows `up to date`, ahead, behind, or diverged counts.
- Friendly error handling for non-repo paths and git command failures.

## Requirements

- Go 1.18+ (1.21+ recommended)
- Git installed and available in `PATH`

## Run

```bash
go run .
go run . origin
go run . --help
go run . --version
```

Launch from inside a Git repository directory. The optional `remote` argument defaults to `origin`.

## Controls

- `up/down` or `j/k`: move commit selection
- `p`: toggle commit panel
- `r`: refresh immediately
- `q` or `Ctrl+C`: quit

## Status Meanings

- `Status: clean` means there are no local uncommitted changes.
- `Status: dirty` means there are local changes (modified/staged/untracked files).
- `Sync: up to date` means local branch and remote branch are aligned.
- `Sync: X commits behind remote` means remote has commits you do not have locally yet.
- `Sync: X commits ahead of remote` means your local branch has commits not on remote yet.
- `Sync: A behind / B ahead` means local and remote have diverged.

## Build and Release

This project includes a `Makefile` with common tasks:

```bash
make test
make build
make release
make clean
```

- `test`: run unit/package tests with `go test ./...`
- `build`: local build to `bin/git-remote-commits` with version `dev`
- `release`: cross-platform binaries in `dist/`
  - `git-remote-commits-linux-amd64`
  - `git-remote-commits-linux-arm64`
  - `git-remote-commits-darwin-amd64`
  - `git-remote-commits-darwin-arm64`
  - `git-remote-commits-windows-amd64.exe`
  - `git-remote-commits-windows-arm64.exe`
- `clean`: remove `bin/` and `dist/`

## Project Structure

- `main.go`: app bootstrap
- `git/`: git command execution and parsing
- `model/`: Bubble Tea state management and polling loop
- `ui/`: terminal rendering/layout

## Community Standards

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Contributing Guide](CONTRIBUTING.md)
- Pull request template: `.github/pull_request_template.md`
- Issue templates:
  - `.github/ISSUE_TEMPLATE/bug_report.yml`
  - `.github/ISSUE_TEMPLATE/feature_request.yml`

## CI and Release Automation

- Test workflow: `.github/workflows/test.yml`
  - Runs on pushes and pull requests to `main`
  - Executes `make test`
- Release workflow: `.github/workflows/release.yml`
  - Runs on version tags (`v*`) or manual trigger
  - Release job depends on successful test job
  - Builds artifacts for linux/darwin/windows on amd64+arm64
  - Injects app version from tag name into binaries
  - Publishes `dist/*` files to GitHub Releases

## Funding

- GitHub Sponsors custom link is configured in `.github/FUNDING.yml`
- Support via PayPal: [paypal.me/wynumarcu23](https://paypal.me/wynumarcu23)

## Notes

- Refresh interval is currently set to 3 seconds in `model/model.go`.
- Remote target is `<remote>/<current-branch>` (remote defaults to `origin` unless passed as an argument).
- Polling performs `git pull <remote> <branch>`, then refreshes commit/sync status.
- Header shows app title and build version.

## License

This project is licensed under the **MIT License** - see [LICENSE](LICENSE).
