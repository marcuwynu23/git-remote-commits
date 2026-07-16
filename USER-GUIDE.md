# git-remote-commits User Guide

Comprehensive reference for git-remote-commits — the live Git remote monitoring TUI.

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Keyboard Controls](#keyboard-controls)
- [Status Meanings](#status-meanings)
- [Commit Detail Panel](#commit-detail-panel)
- [Help Overlay](#help-overlay)
- [Concepts](#concepts)
- [CI/CD Integration](#cicd-integration)
- [Workflows](#workflows)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Installation

### Prerequisites

| Requirement | Version | Notes |
|---|---|---|
| Go | 1.18+ | 1.24+ required to build from source |
| Git | Any | Must be available in `PATH` |
| Terminal | Any modern terminal | Supports 256-color ANSI output |

### Install via `go install`

```bash
go install github.com/marcuwynu23/git-remote-commits@latest
```

This builds the binary to `$GOPATH/bin` (or `$HOME/go/bin` by default). Ensure that directory is in your `PATH`.

### Download a Pre-Built Binary

Visit the [releases page](https://github.com/marcuwynu23/git-remote-commits/releases) and download the archive for your platform:

| Platform | Architecture | Binary Name |
|---|---|---|
| Linux | amd64 | `git-remote-commits-linux-amd64` |
| Linux | arm64 | `git-remote-commits-linux-arm64` |
| macOS (Intel) | amd64 | `git-remote-commits-darwin-amd64` |
| macOS (Apple Silicon) | arm64 | `git-remote-commits-darwin-arm64` |
| Windows | amd64 | `git-remote-commits-windows-amd64.exe` |
| Windows | arm64 | `git-remote-commits-windows-arm64.exe` |

Make the binary executable (`chmod +x`) and place it in a directory on your `PATH`.

### Build from Source

```bash
git clone https://github.com/marcuwynu23/git-remote-commits.git
cd git-remote-commits
make build
```

The binary is written to `bin/git-remote-commits`.

### Verify Installation

```bash
git-remote-commits --version
git-remote-commits --help
```

---

## Quick Start

1. Open a terminal and navigate to any Git repository:

```bash
cd my-project
```

2. Launch the TUI:

```bash
git-remote-commits
```

3. The dashboard appears. After a brief loading animation, you see:
   - **Header**: app name and version
   - **Status line**: repository, branch, working tree status, remote tracking, sync state, last refresh time
   - **Commit list**: scrollable list of commits with date, graph, hash, refs, author, and message
   - **Footer**: keyboard shortcut bar

4. Use `up`/`down` or `j`/`k` to navigate commits. Press `p` to toggle the commit detail panel. Press `q` to quit.

To monitor a different remote:

```bash
git-remote-commits upstream
```

---

## Command Reference

### `git-remote-commits [remote]`

Launches the TUI dashboard inside the current directory. The directory must be a Git repository.

```bash
git-remote-commits
git-remote-commits origin
git-remote-commits my-fork
```

| Argument | Default | Description |
|---|---|---|
| `remote` | `origin` | Remote name used for upstream comparison |

| Flag | Default | Description |
|---|---|---|
| `-h, --help` | — | Print usage information and exit |
| `-v, --version` | — | Print version and exit |

#### Exit Codes

| Code | Condition |
|---|---|
| `0` | Normal exit (user pressed `q` / `Ctrl+C`) |
| `1` | Error: not a Git repo, too many args, unknown flag, or runtime error |

---

## Configuration

git-remote-commits has **no configuration files**. All settings are configured through CLI arguments or build-time injection.

| Setting | Method | Default |
|---|---|---|
| Remote name | CLI argument | `origin` |
| Build version | `-X 'main.Version=<version>'` linker flag | `dev` |
| Refresh interval | Hardcoded in source (`model/model.go`) | `3s` |
| Commit display limit | Hardcoded in source (`model/model.go`) | No limit (all commits) |

To change the refresh interval, edit `defaultRefresh` in `model/model.go` and rebuild:

```go
const defaultRefresh = 5 * time.Second // change from 3s to 5s
```

### Version Injection

When building for release, inject the version via a linker flag:

```bash
go build -ldflags "-X 'main.Version=v2.1.1'" -o git-remote-commits .
```

The `make release` target uses `VERSION` variable and cross-compiles with version injection automatically.

---

## Keyboard Controls

### Navigation

| Key | Action |
|---|---|
| `up` / `k` | Move selection up one commit |
| `down` / `j` | Move selection down one commit |
| `g` / `Home` | Jump to the latest (top) commit |
| `G` / `End` | Jump to the initial (bottom) commit |

### Commit Panel

| Key | Action |
|---|---|
| `p` / `P` | Toggle the commit detail panel on/off |
| `[` | Scroll panel up by 1 line |
| `]` | Scroll panel down by 1 line |
| `u` / `Ctrl+U` | Scroll panel up by 5 lines |
| `d` / `Ctrl+D` | Scroll panel down by 5 lines |
| `PgUp` | Scroll panel up by 5 lines |
| `PgDn` | Scroll panel down by 5 lines |

### General

| Key | Action |
|---|---|
| `r` | Force an immediate refresh (pull + reload) |
| `?` / `h` / `H` | Toggle help overlay |
| `q` / `Ctrl+C` | Quit the application |

---

## Status Meanings

### Working Tree Status

The status line shows `Status: clean` or `Status: dirty`.

| Status | Meaning |
|---|---|
| **clean** | No local uncommitted changes |
| **dirty** | Modified, staged, or untracked files exist |

### Sync Status

The sync status shows the relationship between your local branch and its remote tracking branch.

| Status | Meaning |
|---|---|
| **Sync: online** | GitHub is reachable and remote comparison is active |
| **Sync: offline** | GitHub is unreachable (network issue or remote is not GitHub) |
| **X commits behind remote** | Remote has X commits you don't have locally — run `git pull` |
| **X commits ahead of remote** | You have X commits not yet pushed — run `git push` |
| **A behind / B ahead** | Branches have diverged — both local and remote have unique commits |
| **up to date** | Local and remote are aligned |

### Remote Status Indicator

The sync chip uses color coding:

- **Green** (`online` / `up to date`) — connected and synced
- **Yellow** (`offline` / `behind` / `ahead`) — network issue or out of sync
- **Blue** (`pending`) — refresh in progress

---

## Commit Detail Panel

Press `p` to toggle the commit detail panel. It appears at the bottom of the screen and shows information about the currently selected commit.

### Left Column

| Field | Description |
|---|---|
| **Commit Hash** | Abbreviated hash (`abc123f`) |
| **Author** | Author name and email username |
| **When** | Local date and time of the commit |
| **Refs** | Branch heads and tags pointing to this commit |
| **Title** | The first line of the commit message |
| **Body** | Full commit message body rendered as Markdown (via Glamour) |

### Right Column

**File Changes** lists every file affected by the commit, with an icon indicating the change type:

| Icon | Meaning |
|---|---|
| `+` | File added |
| `-` | File deleted |
| `~` | File modified |
| `→` | File renamed |

### Scrolling

When the panel content exceeds the available height, scroll indicators (`↑` / `↓`) appear in the top-left corner. Use `[`/`]` or `u`/`d` to scroll.

---

## Help Overlay

Press `?`, `h`, or `H` to toggle the full help overlay. It shows all keyboard controls organized by category:

- **Navigation** — cursor movement and jumping
- **Commit Panel** — panel toggle and scrolling
- **General** — refresh, help, and quit

The help overlay closes when you press `?`, `h`, `H`, or `q`.

---

## Concepts

### The Refresh Cycle

1. Every **3 seconds**, the tool checks GitHub connectivity.
2. If online, it runs `git pull <remote> <branch>` to fetch latest changes.
3. It collects a snapshot: branch name, working tree status, commit log, ahead/behind counts.
4. The new snapshot is compared to the previous one. New commits (hashes not seen before) are marked with a green `●`.
5. If new commits are detected, a terminal bell (`\a`) rings.
6. The display re-renders with updated data.

### Dynamic Remote Tracking

Unlike tools that track a hardcoded `origin/main` ref, git-remote-commits resolves the upstream as `<remote>/<current-branch>`. This means:

- If you switch branches while the TUI is running, the next refresh tracks the new branch's remote counterpart.
- There is no manual configuration needed per branch.

### Author Highlighting

If the current commit's author matches the `user.name` configured in your local Git config, the author name appears in **bold green**. Other authors appear in normal green.

---

## CI/CD Integration

### GitHub Actions

**Test workflow** (`.github/workflows/test.yml`):

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: make test
```

**Release workflow** (`.github/workflows/release.yml`):

```yaml
name: Release
on:
  push:
    tags: ['v*']
  workflow_dispatch:
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: make test
  release:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: stable
      - run: make release
      - uses: softprops/action-gh-release@v2
        with:
          files: dist/*
```

### Makefile Targets

| Target | Description |
|---|---|
| `make test` | Run all tests |
| `make build` | Build for current platform (version: `dev`) |
| `make release` | Cross-compile for all 6 supported platforms |
| `make shortcut` | (Windows only) Build to `C:/Bin/git-remote-commits/` |
| `make clean` | Remove `bin/` and `dist/` output directories |

---

## Workflows

### Monorepo Monitoring

If your monorepo has multiple remotes (e.g., `origin` for your fork, `upstream` for the main repo):

```bash
# Monitor the upstream remote
cd my-monorepo
git-remote-commits upstream

# Monitor your fork
git-remote-commits origin
```

### Using with tmux

Split your tmux pane and run the TUI alongside your editor:

```bash
tmux split-window -h
# In the new pane:
cd my-project
git-remote-commits
```

### SSH Session Monitoring

The TUI works over SSH. Connect to a remote dev server and launch:

```bash
ssh dev-server
cd /var/www/my-project
git-remote-commits
```

The terminal bell (`\a`) may not propagate through all SSH configurations.

---

## Troubleshooting

### "current directory is not a git repository"

**Cause:** You ran the tool outside a Git working tree.
**Fix:** `cd` into a Git repository before launching.

```bash
cd my-project
git-remote-commits
```

### "Error: too many arguments"

**Cause:** You passed more than one positional argument.
**Fix:** Provide only the remote name.

```bash
git-remote-commits origin  # correct
git-remote-commits origin main  # incorrect
```

### "Error: unknown flag"

**Cause:** You passed an unrecognized flag.
**Fix:** Use only `-h`/`--help` or `-v`/`--version`.

```bash
git-remote-commits -h     # correct
git-remote-commits --foo  # incorrect
```

### Commit list is empty

**Cause:** The repository has no commits yet, or the `git log` command failed.
**Fix:** Ensure the repository has at least one commit. Check that Git is installed and in `PATH`.

### Sync always shows "offline"

**Cause:** GitHub is unreachable, or the remote is not hosted on GitHub.
**Fix:** Check your network connection. The tool probes `https://api.github.com` — if the remote is on a different host (GitLab, self-hosted), the sync will show "offline" but upstream comparison still works.

### Commit panel shows "No file changes."

**Cause:** The commit selected is a merge commit with no diff, or the diff output is empty.
**Fix:** Select a non-merge commit to see file changes.

### TUI doesn't render correctly

**Cause:** The terminal doesn't support the required ANSI escape codes, or the window is too narrow.
**Fix:** Use a modern terminal emulator (iTerm2, Windows Terminal, GNOME Terminal, Konsole, Alacritty). Resize the terminal to at least 60 columns wide.

---

## FAQ

**Q: Can I change the refresh interval?**
A: Yes — edit `defaultRefresh` in `model/model.go` (line 17) and rebuild the binary.

**Q: Does the tool push changes?**
A: No. It only reads repository state and runs `git pull` to fetch updates. It never pushes, merges, or modifies Git configuration.

**Q: Can I use it with GitLab or Bitbucket?**
A: Yes. Remote tracking and syncing work with any Git host. The "offline" indicator only applies to GitHub connectivity detection.

**Q: Does it work on Windows?**
A: Yes. Windows binaries are provided. Use Windows Terminal or any modern terminal emulator.

**Q: How many commits are displayed?**
A: By default, all commits. The `commitLimit` constant in `model/model.go` can be set to a positive integer to limit the log output.

**Q: Why doesn't the terminal bell ring on new commits?**
A: Some terminals disable the bell by default. Check your terminal's notification/bell settings. SSH sessions may not forward the bell.

**Q: Can I run this in the background?**
A: No — it's a full-screen TUI application. Use tmux or a dedicated terminal window.

**Q: What does the `●` marker mean?**
A: It indicates commits that appeared since the last refresh cycle — new commits from `git pull`.

**Q: How do I report a bug?**
A: Use the [bug report template](https://github.com/marcuwynu23/git-remote-commits/issues/new?template=bug_report.yml) on GitHub.

**Q: How do I request a feature?**
A: Use the [feature request template](https://github.com/marcuwynu23/git-remote-commits/issues/new?template=feature_request.yml) on GitHub.

**Q: Is this project open source?**
A: Yes, MIT licensed. See [LICENSE](LICENSE).

**Q: How is this different from `tig`?**
A: git-remote-commits auto-refreshes, tracks remote sync status, highlights new commits, and works with any remote — all with zero configuration.

---

For additional help, open an issue at [github.com/marcuwynu23/git-remote-commits/issues](https://github.com/marcuwynu23/git-remote-commits/issues).
