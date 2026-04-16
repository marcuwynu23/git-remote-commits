# commitview

`commitview` is a terminal UI that monitors Git activity in real time, like an `htop`-style dashboard for repository commits.

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
```

Launch from inside a Git repository directory.

## Controls

- `up/down` or `j/k`: move commit selection
- `r`: refresh immediately
- `q` or `Ctrl+C`: quit

## Build and Release

This project includes a `Makefile` with common tasks:

```bash
make build
make release
make clean
```

- `build`: local build to `bin/commitview`
- `release`: cross-platform binaries in `dist/`
  - `commitview-linux-amd64`
  - `commitview-darwin-amd64`
  - `commitview-darwin-arm64`
  - `commitview-windows-amd64.exe`
- `clean`: remove `bin/` and `dist/`

## Project Structure

- `main.go`: app bootstrap
- `git/`: git command execution and parsing
- `model/`: Bubble Tea state management and polling loop
- `ui/`: terminal rendering/layout

## Notes

- Refresh interval is currently set to 3 seconds in `model/model.go`.
- Remote tracking assumes `origin/main`.

## License

This project is licensed under the **MIT License** - see [LICENSE](LICENSE).
