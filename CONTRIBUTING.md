# Contributing to git-remote-commits

Thank you for your interest in contributing! This guide covers everything you need to know to contribute effectively.

- [Code of Conduct](#code-of-conduct)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Makefile Reference](#makefile-reference)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Conventions](#commit-conventions)
- [PR Process](#pr-process)
- [Release Process](#release-process)
- [Questions](#questions)

---

## Code of Conduct

This project follows the [Contributor Covenant v2.1](CODE_OF_CONDUCT.md). All contributors must abide by its terms.

---

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.24+ | Build, test, and run |
| Git | Any | Version control |
| Make | Any | Task automation |
| A modern terminal | — | Verify TUI rendering |

---

## Project Structure

```
.
├── main.go                   # Entry point: CLI parsing, repo check, Bubble Tea bootstrap
├── git/
│   └── git.go                # Git command execution, snapshot collection, commit parsing
├── model/
│   └── model.go              # Bubble Tea model: state machine, polling loop, key handling
├── ui/
│   └── ui.go                 # Rendering: Lip Gloss styles, layout, commit list, panel, help
├── .github/
│   ├── workflows/
│   │   ├── test.yml          # CI test workflow (push/PR to main)
│   │   └── release.yml       # CD release workflow (tags v*, manual dispatch)
│   ├── ISSUE_TEMPLATE/       # Bug report and feature request YAML forms
│   ├── FUNDING.yml           # PayPal funding link
│   └── pull_request_template.md
├── Makefile                  # Build automation: test, build, release, clean
├── go.mod / go.sum           # Go module dependencies
├── CHANGELOG.md              # Auto-generated release history
├── GUIDE.md                  # TUI components guide (Bubble Tea + Lip Gloss patterns)
├── README.md                 # Project front door
├── USER-GUIDE.md             # Comprehensive user reference
├── CONTRIBUTING.md           # This file
├── CODE_OF_CONDUCT.md        # Community standards
├── SECURITY.md               # Security policy
└── LICENSE                   # MIT License
```

---

## Makefile Reference

| Command | Description |
|---|---|
| `make test` | Run `go test ./...` |
| `make build` | Build to `bin/git-remote-commits` with version `dev` |
| `make release` | Cross-compile 6 binaries to `dist/` |
| `make shortcut` | (Windows only) Build to `C:/Bin/git-remote-commits/git-remote-commits.exe` |
| `make clean` | Remove `bin/` and `dist/` directories |

---

## Development Workflow

1. **Fork** the repository on GitHub.
2. **Clone** your fork:

```bash
git clone https://github.com/YOUR_USERNAME/git-remote-commits.git
cd git-remote-commits
```

3. **Create a feature branch** from `main`:

```bash
git checkout -b feat/my-feature
```

4. **Make your changes.** Keep them focused and small.

5. **Run tests** to verify nothing is broken:

```bash
make test
```

6. **Build and manually verify** the TUI renders correctly:

```bash
make build
./bin/git-remote-commits
```

7. **Commit** your changes using conventional commits (see [below](#commit-conventions)).

8. **Push** your branch:

```bash
git push -u origin feat/my-feature
```

9. **Open a pull request** against `main`. Fill out the PR template with:
   - What changed and why
   - Screenshots or terminal output for UI changes
   - Confirmation that tests pass

---

## Coding Standards

### Naming

- Go conventions: `camelCase` for unexported, `PascalCase` for exported.
- Package names: short, lowercase, no underscores.
- File names: lowercase, underscore-separated (`git.go`, `model.go`).

### Imports

Group standard library, third-party, and internal imports separated by blank lines:

```go
import (
    "fmt"
    "os"

    "github.com/charmbracelet/bubbletea"

    "git-remote-commits/git"
)
```

### Errors

- Use `errors.New()` or `fmt.Errorf()` for static messages.
- Wrap errors with context: `fmt.Errorf("failed to read branch: %w", err)`.
- Prefer explicit error returns over panics.

### Formatting

- Run `go fmt ./...` before committing.
- Use `go vet ./...` to catch common issues.

### Style Rules

- **Exported functions** get a doc comment.
- **No unused parameters** — use `_` for unused variables.
- **Prefer table-driven tests** for logic-heavy functions.

---

## Testing

- **Test directory:** Package-level `_test.go` files alongside production code.
- **Coverage target:** Aim for >70% coverage on new logic.
- **Test pattern:** Table-driven tests using anonymous structs:

```go
func TestIsHexHashToken(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  bool
    }{
        {"valid short hash", "abc123f", true},
        {"valid full hash", strings.Repeat("a", 40), true},
        {"too short", "abc12", false},
        {"invalid chars", "abc123z", false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := isHexHashToken(tt.input); got != tt.want {
                t.Errorf("isHexHashToken(%q) = %v, want %v", tt.input, got, tt.want)
            }
        })
    }
}
```

- **Run tests:**

```bash
make test
# Or:
go test ./... -v -count=1
```

- **TUI tests** are currently manual (launch the binary and verify rendering).

---

## Commit Conventions

We follow [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages. This enables automatic changelog generation.

```
<type>(<scope>): <description>
```

### Types

| Type | Usage |
|---|---|
| `feat` | A new feature visible to the user |
| `fix` | A bug fix |
| `chore` | Maintenance, dependency updates, tooling |
| `docs` | Documentation changes only |
| `refactor` | Code restructuring without behavior change |
| `test` | Adding or updating tests |
| `perf` | Performance improvements |
| `style` | Formatting, linting, whitespace (no logic change) |

### Scope

The scope should be the package or area affected:

- `cli` — argument parsing, `main.go`
- `git` — git command execution, snapshot collection
- `model` — Bubble Tea model, state, key handling
- `ui` — rendering, styles, layout
- `docs` — any documentation file
- `release` — build and release configuration

### Examples

```
feat(git): add support for git-worktree repos
fix(model): handle detached HEAD correctly
docs(ui): update keyboard shortcut reference in README
chore(deps): bump bubbletea to v0.26.0
refactor(git): extract ahead/behind calculation into helper
```

### Breaking Changes

Add `!` after the type/scope and include a `BREAKING CHANGE` trailer:

```
feat(api)!: remove deprecated --verbose flag

BREAKING CHANGE: The --verbose flag has been removed. Use --debug instead.
```

---

## PR Process

### Before Opening a PR

- [ ] Branch is up to date with `main`
- [ ] `make test` passes
- [ ] `make build` succeeds
- [ ] `go vet ./...` is clean
- [ ] UI changes verified manually (binary runs correctly)
- [ ] Docs updated (`README.md`, `USER-GUIDE.md`, etc.) if behavior changed
- [ ] Tests added or updated for logic changes

### During Review

- Keep the PR scope focused. One feature/fix per PR.
- Address review feedback with additional commits (no squashing until merge).
- Maintainers will squash-merge to keep a clean history.

### What Gets Merged

- Features that are scoped and well-tested
- Bug fixes with clear reproduction steps
- Documentation improvements
- Refactoring that improves code quality without changing behavior

### What Doesn't

- Large, unfocused PRs touching many packages
- Changes without tests (for logic changes)
- Breaking changes without prior discussion (open an issue first)

---

## Release Process

1. A maintainer creates a version tag:

```bash
git tag v2.2.0
git push origin v2.2.0
```

2. The [Release workflow](.github/workflows/release.yml) triggers automatically:
   - Runs `make test` to verify
   - Cross-compiles binaries for all 6 platforms
   - Publishes artifacts to GitHub Releases
   - Version is injected via `-X 'main.Version=<tag>'`

3. The release appears at [github.com/marcuwynu23/git-remote-commits/releases](https://github.com/marcuwynu23/git-remote-commits/releases).

### Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **Patch** (`v2.2.1`): bug fixes, small improvements
- **Minor** (`v2.3.0`): new features, non-breaking changes
- **Major** (`v3.0.0`): breaking changes to CLI, API, or behavior

---

## Questions

If you have questions not covered here:

- Open a [Discussion](https://github.com/marcuwynu23/git-remote-commits/discussions)
- File an [Issue](https://github.com/marcuwynu23/git-remote-commits/issues)
- Reach out via [PayPal](https://paypal.me/wynumarcu23)
