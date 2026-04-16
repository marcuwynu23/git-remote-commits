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
```

Launch from inside a Git repository directory. The optional `remote` argument defaults to `origin`.

## Controls

- `up/down` or `j/k`: move commit selection
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
make build
make release
make clean
```

- `build`: local build to `bin/git-remote-commits`
- `release`: cross-platform binaries in `dist/`
  - `git-remote-commits-linux-amd64`
  - `git-remote-commits-darwin-amd64`
  - `git-remote-commits-darwin-arm64`
  - `git-remote-commits-windows-amd64.exe`
- `clean`: remove `bin/` and `dist/`

## Project Structure

- `main.go`: app bootstrap
- `git/`: git command execution and parsing
- `model/`: Bubble Tea state management and polling loop
- `ui/`: terminal rendering/layout

## Notes

- Refresh interval is currently set to 3 seconds in `model/model.go`.
- Remote target is `<remote>/<current-branch>` (remote defaults to `origin` unless passed as an argument).
- Polling performs `git pull <remote> <branch>`, then refreshes commit/sync status.

## License

This project is licensed under the **MIT License** - see [LICENSE](LICENSE).
