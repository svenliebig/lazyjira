# Stack

The technology stack for this project should align be similar to the implementations of [lazygit](https://github.com/jesseduffield/lazygit). It's a Go project that uses a TUI library for the UI.

## CI/CD

| Tool | Purpose |
|------|---------|
| GitHub Actions | CI/CD runner for tests, releases |
| [Release Please](https://github.com/googleapis/release-please-action) | Automated changelog and version bumping based on conventional commits |
| [GoReleaser](https://goreleaser.com/) | Multi-platform binary builds and GitHub Release publishing |

Releases are fully automated: merging a Release Please PR to `main` builds and publishes binaries for Linux, macOS, and Windows (amd64, arm64, 386, ARMv6/v7).
