<div align="center">

# git-remote-commits

<a href="https://github.com/marcuwynu23/git-remote-commits/releases"><img src="https://img.shields.io/github/v/release/marcuwynu23/git-remote-commits" alt="Release"></a>
<a href="LICENSE"><img src="https://img.shields.io/github/license/marcuwynu23/git-remote-commits?logo=github" alt="License"></a>
<a href="https://github.com/marcuwynu23/git-remote-commits/stargazers"><img src="https://img.shields.io/github/stars/marcuwynu23/git-remote-commits" alt="Stars"></a>
<img src="https://img.shields.io/github/go-mod/go-version/marcuwynu23/git-remote-commits" alt="Go version">

<strong>Live Git remote monitoring in your terminal.</strong> An `htop`-style dashboard for repository commits with auto-refresh, remote tracking, and commit detail panels.

➡️ **[Read the full user guide →](USER-GUIDE.md)**

</div>

---

## Table of Contents

- [What Is git-remote-commits?](#what-is-git-remote-commits)
- [Use Cases](#use-cases)
- [Benefits](#benefits)
- [Comparison](#comparison)
- [User Guide](USER-GUIDE.md)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Reference](#cli-reference)
- [Configuration](#configuration)
- [Example Output](#example-output)
- [CI/CD Integration](#cicd-integration)
- [Development](#development)
- [Architecture](#architecture)
- [Community Standards](#community-standards)
- [License](#license)

---

## What Is git-remote-commits?

**git-remote-commits** is a terminal UI that monitors Git activity in real time — like an `htop`-style dashboard for repository commits.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), [Lip Gloss](https://github.com/charmbracelet/lipgloss), and [Glamour](https://github.com/charmbracelet/glamour).

### What It Does

- **Displays** a live commit list with date, author, message, and ASCII graph.
- **Tracks** the current branch's remote upstream dynamically — not a fixed branch.
- **Highlights** new commits between auto-refresh polls with a green dot indicator.
- **Shows** upstream sync status: up to date, ahead, behind, or diverged.
- **Renders** commit detail panels with Markdown bodies, file changes, and diff previews.
- **Pulls** remote changes automatically before each refresh cycle.
- **Detects** GitHub connectivity to determine online/offline state.

### Why Use It?

| Problem                                                     | How git-remote-commits Solves It                                              |
| ----------------------------------------------------------- | ----------------------------------------------------------------------------- |
| Need to see team commits in real time                       | **Live dashboard** — auto-refreshes every 3 seconds, highlighting new commits |
| Unsure if your branch is synced with remote                 | **Sync status line** — shows ahead/behind/diverged counts at a glance         |
| Context-switching to `git log` breaks flow                  | **Always-on terminal view** — stays open and updates while you work           |
| CI/CD pipeline produces new commits you miss                | **New-commit notification** — visual green dot + terminal bell                |
| Need to review a commit's diff without leaving the terminal | **Built-in commit panel** — shows diff, file changes, and rendered Markdown   |
| Working in a monorepo with many remotes                     | **Configurable remote** — pass any remote name as an argument                 |

### The Philosophy

1. **Minimal setup, maximum value.** One command to launch, zero config files for basic use.
2. **Your process stays yours.** The tool reads your repo state — it never writes config, changes branches, or modifies Git state beyond pulling.
3. **Terminal-native.** Designed for tmux, SSH sessions, and developers who live in the terminal.

---

## Use Cases

| Scenario                       | How git-remote-commits Helps                                            |
| ------------------------------ | ----------------------------------------------------------------------- |
| **Pair programming / mobbing** | See teammates' commits appear live as they push to the shared branch    |
| **CI/CD monitoring**           | Watch the commit list update after a pipeline produces new commits      |
| **Pre-deploy sanity check**    | Verify the remote branch is exactly where you expect it before merging  |
| **Open source maintenance**    | Monitor incoming PR merges and maintainer commits on the base branch    |
| **Learning Git workflows**     | Visualize commit history, branch topology, and remote sync in real time |

---

## Benefits for Developers

- **Zero configuration** — launch inside any Git repo and it works
- **No TUI lock-in** — standard terminal controls, `q` to quit cleanly
- **Active feedback** — terminal bell rings when new commits arrive
- **Remote-aware** — automatically pulls and compares with `<remote>/<branch>`
- **Commit inspection** — toggle a detail panel with diff, file list, and rendered Markdown
- **Cross-platform** — binaries for Linux, macOS, and Windows (amd64 + arm64)
- **Responsive layout** — status line adapts to terminal width
- **Scrollable history** — navigate thousands of commits with `g`/`G`/`PgUp`/`PgDn`
- **Lightweight** — no heavy JavaScript runtime or system daemon
- **Open source** — MIT licensed, contributions welcome

---

## Comparison

| Aspect                       | git-remote-commits           | `git log --oneline` | `tig`       | `lazygit`       | Manual (`git pull && git log`) |
| ---------------------------- | ---------------------------- | ------------------- | ----------- | --------------- | ------------------------------ |
| **Setup time**               | ~10 seconds                  | Instant             | ~30 seconds | ~30 seconds     | Ongoing effort                 |
| **Auto-refresh**             | ✅ Every 3s                  | ❌                  | ❌          | ❌              | ❌                             |
| **Remote sync status**       | ✅ Ahead/behind/diverged     | ❌                  | ❌          | ⚠️ Partial      | ❌                             |
| **New commit highlighting**  | ✅ Green dot + terminal bell | ❌                  | ❌          | ❌              | ❌                             |
| **Commit detail panel**      | ✅ Diff + files + Markdown   | ❌                  | ✅          | ✅              | ❌                             |
| **Auto-pull before refresh** | ✅                           | ❌                  | ❌          | ❌              | Manual                         |
| **GitHub online detection**  | ✅                           | ❌                  | ❌          | ❌              | ❌                             |
| **Cross-platform**           | ✅ Linux/macOS/Windows       | ✅                  | ⚠️ Limited  | ✅              | ✅                             |
| **Config files**             | None required                | ✅ Flags            | ✅ `.tigrc` | ✅ `config.yml` | N/A                            |
| **Learning curve**           | ~1 minute                    | ~1 minute           | ~10 minutes | ~10 minutes     | Instant                        |
| **Background process**       | ❌ Full-screen TUI           | ❌                  | ❌          | ❌              | ❌                             |

---

## Installation

### Go Install (recommended)

```bash
go install github.com/marcuwynu23/git-remote-commits@latest
```

### Binary Download

Download the pre-built binary for your platform from the [releases page](https://github.com/marcuwynu23/git-remote-commits/releases).

| Platform | Architecture | File                                   |
| -------- | ------------ | -------------------------------------- |
| Linux    | amd64        | `git-remote-commits-linux-amd64`       |
| Linux    | arm64        | `git-remote-commits-linux-arm64`       |
| macOS    | amd64        | `git-remote-commits-darwin-amd64`      |
| macOS    | arm64        | `git-remote-commits-darwin-arm64`      |
| Windows  | amd64        | `git-remote-commits-windows-amd64.exe` |
| Windows  | arm64        | `git-remote-commits-windows-arm64.exe` |

### Verify

```bash
git-remote-commits --help
git-remote-commits --version
```

### Prerequisites

| Requirement | Version                            |
| ----------- | ---------------------------------- |
| Go          | 1.18+ (1.24+ to build from source) |
| Git         | Any modern version in `PATH`       |

---

## Quick Start

```bash
# Launch inside any Git repository
cd my-project
git-remote-commits

# Or specify a custom remote name
git-remote-commits upstream
```

That's it. The TUI opens showing your commit history with live remote tracking.

---

## CLI Reference

### `git-remote-commits [remote]`

Launches the TUI dashboard.

```bash
git-remote-commits
git-remote-commits origin
git-remote-commits upstream
```

| Argument | Default  | Description                                   |
| -------- | -------- | --------------------------------------------- |
| `remote` | `origin` | Remote name to compare against current branch |

| Flag            | Description              |
| --------------- | ------------------------ |
| `-h, --help`    | Show help message        |
| `-v, --version` | Show application version |

### Keyboard Controls

| Key                 | Action                                 |
| ------------------- | -------------------------------------- |
| `up` / `k`          | Move selection up                      |
| `down` / `j`        | Move selection down                    |
| `g` / `Home`        | Jump to latest commit (top)            |
| `G` / `End`         | Jump to initial commit (bottom)        |
| `p` / `P`           | Toggle commit detail panel             |
| `r`                 | Force immediate refresh                |
| `[` / `]`           | Scroll commit panel by 1 line          |
| `u` / `d`           | Scroll commit panel by chunk (5 lines) |
| `Ctrl+U` / `Ctrl+D` | Scroll commit panel by chunk           |
| `PgUp` / `PgDn`     | Scroll commit panel by chunk           |
| `?` / `h` / `H`     | Toggle help overlay                    |
| `q` / `Ctrl+C`      | Quit                                   |

---

## Configuration

git-remote-commits is designed to require **zero configuration** for basic use.

| Setting              | How to configure                                                   |
| -------------------- | ------------------------------------------------------------------ |
| **Remote name**      | Pass as CLI argument: `git-remote-commits <name>`                  |
| **Build version**    | Inject at build time: `-X 'main.Version=<version>'`                |
| **Refresh interval** | Hardcoded in `model/model.go` (`defaultRefresh = 3 * time.Second`) |

There are no YAML, JSON, or TOML config files. The tool reads your repository state directly.

---

## Example Output

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                                                                              │
│  git remote-commits v2.1.1                                                   │
│                                                                              │
│  Repository git-remote-commits  Branch main  Status clean                    │
│  Remote origin/main  Sync online  Refresh 3:04PM                            │
│                                                                              │
│  2026-04-17: 3:04 PM  ● abc123f  (HEAD -> main)  jdoe  Fix refresh race     │
│  2026-04-17: 2:58 PM    def456g  (tag: v2.1.0)    jdoe  Add diff preview    │
│  2026-04-17: 2:30 PM    ghi789h  (origin/main)     jdoe  Refactor model     │
│  ...                                                                         │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐     │
│  │ Commit Hash: abc123f    File Changes                                │     │
│  │ Author: jdoe            ~ model/model.go                             │     │
│  │ When: April 17, 2026    ~ git/git.go                                 │     │
│  │ Title: Fix refresh race                                              │     │
│  │                                                                      │     │
│  │ Fix a race condition where a pending refresh could hang the UI.     │     │
│  └─────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  [up/j | down/k]  [PgUp/u | PgDn/d]  [g/Home | G/End]  [p]  [r]  [?/h]     │
└──────────────────────────────────────────────────────────────────────────────┘
```

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
    tags: ["v*"]
  workflow_dispatch:
jobs:
  release:
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

---

## Development

### Prerequisites

| Tool | Version | Purpose         |
| ---- | ------- | --------------- |
| Go   | 1.24+   | Build and test  |
| Git  | Any     | Source control  |
| Make | Any     | Task automation |

### Commands

| Command         | Description                                       |
| --------------- | ------------------------------------------------- |
| `make test`     | Run `go test ./...`                               |
| `make build`    | Build to `bin/git-remote-commits` (version `dev`) |
| `make release`  | Cross-compile 6 binaries to `dist/`               |
| `make shortcut` | (Windows) Build to `C:/Bin/git-remote-commits/`   |
| `make clean`    | Remove `bin/` and `dist/`                         |
| `go run .`      | Run directly from source                          |

### Project Structure

```
.
├── main.go              # Entry point, CLI args, Bubble Tea bootstrap
├── git/
│   └── git.go           # Git command execution, snapshot collection, commit parsing
├── model/
│   └── model.go         # Bubble Tea model, state management, polling loop, key handling
├── ui/
│   └── ui.go            # Terminal rendering, Lip Gloss styles, layout, commit panel
├── .github/
│   ├── workflows/
│   │   ├── test.yml     # CI test workflow
│   │   └── release.yml  # CD release workflow
│   ├── ISSUE_TEMPLATE/  # Bug report and feature request forms
│   └── FUNDING.yml      # Funding configuration
├── Makefile             # Build automation
├── go.mod / go.sum      # Go module definition
└── CHANGELOG.md         # Auto-generated release history
```

---

## Architecture

1. **`main.go`** parses CLI arguments, validates the working directory is a Git repository, and starts a Bubble Tea program with alternate screen mode.
2. **`model/model.go`** is the core state machine. It runs a polling loop that every 3 seconds calls `git.CollectSnapshot()`, detects new commits by hash comparison, and triggers re-renders.
3. **`git/git.go`** wraps Git CLI commands (`rev-parse`, `status`, `log`, `pull`, `rev-list`, `show`, `remote get-url`) and pings the GitHub API to detect online/offline state.
4. **`ui/ui.go`** renders the header, status line, scrollable commit list, commit detail panel (with parsed file changes and Glamour Markdown rendering), help overlay, and footer shortcuts.

Data flows in one direction: Git commands → `git.Snapshot` → Bubble Tea model (`snapshotMsg`) → `ui.Render()` → terminal output.

---

## Community Standards

- [Code of Conduct](CODE_OF_CONDUCT.md)
- [Contributing Guide](CONTRIBUTING.md)
- [User Guide](USER-GUIDE.md)
- [Security Policy](SECURITY.md)
- Pull request template: `.github/PULL_REQUEST_TEMPLATE.md`
- Issue templates: `.github/ISSUE_TEMPLATE/*.yml`

## Funding

Support this project via PayPal: [paypal.me/wynumarcu23](https://paypal.me/wynumarcu23)

---

## License

This project is licensed under the **MIT License** — see [LICENSE](LICENSE).
