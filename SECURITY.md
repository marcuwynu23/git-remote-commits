# Security Policy

## Supported Versions

| Version | Supported |
|---|---|
| >= 2.x | ✅ Active development and security patches |
| < 2.0 | ❌ No longer supported |

## Reporting a Vulnerability

If you discover a security vulnerability in git-remote-commits, please report it privately.

**Do not** open a public issue. Instead, please:

1. Open a **private security report** on GitHub:
   [github.com/marcuwynu23/git-remote-commits/security/advisories/new](https://github.com/marcuwynu23/git-remote-commits/security/advisories/new)

2. Alternatively, contact the maintainer directly at the GitHub account [`marcuwynu23`](https://github.com/marcuwynu23).

### What to Include

- A clear description of the vulnerability
- Steps to reproduce (if applicable)
- Potential impact
- Any suggested fix (optional)

### Response Timeline

- **Acknowledgment:** within 48 hours
- **Assessment:** within 5 business days
- **Fix:** timeline depends on severity, communicated during assessment

## Security Considerations for This Project

git-remote-commits is a read-only TUI tool that:

- **Executes Git commands** (`git pull`, `git log`, `git show`, etc.) in the current repository
- **Makes outbound HTTP requests** to `https://api.github.com` for connectivity detection
- **Does not** push, merge, or modify remote repositories
- **Does not** store credentials or sensitive data
- **Does not** write configuration files

The binary does not require any network permissions beyond:
- Access to your local Git repository
- Outbound HTTPS to `api.github.com` (for online/offline detection; gracefully degrades if blocked)
