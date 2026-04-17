# Contributing to git-remote-commits

Thank you for contributing.

## Getting started

1. Fork the repository.
2. Clone your fork.
3. Create a feature branch from `main`.
4. Run local checks before opening a PR:

```bash
make test
make build
```

## Development guidelines

- Keep changes focused and small.
- Follow existing Go style and naming.
- Update docs (`README.md`, this guide) when behavior changes.
- Prefer adding or updating tests when introducing logic changes.

## Pull request process

1. Make sure your branch is up to date with `main`.
2. Ensure `make test` passes.
3. Fill out the pull request template.
4. Describe what changed and why.
5. Include screenshots or terminal output when UI behavior changes.

## Commit messages

- Use concise, clear messages.
- Prefer intent-first wording, e.g. `fix refresh queue race`.

## Reporting issues

Use the issue templates for bug reports or feature requests and include reproducible steps.
