# 2026-03-22 — Release Automation

## Goal

Set up automated releases via GitHub Actions, mirroring the pipeline used in [svenliebig/phopy](https://github.com/svenliebig/phopy).

## Completed work

- Created `.github/workflows/ci.yml` — runs `go test ./...` on pushes to `main` and on all pull requests
- Created `.github/workflows/release-please.yml` — monitors `main` for conventional commits, maintains a release PR with auto-generated changelog, triggers GoReleaser when the PR is merged
- Created `.github/workflows/release.yml` — runs GoReleaser on any `v*` tag push or manual `workflow_dispatch`
- Created `.release-please-config.json` — release type `go`, single package at root
- Created `.release-please-manifest.json` — initial version `0.1.0`
- Created `.goreleaser.yaml` — builds `lazyjira` binary for Linux, macOS, Windows × amd64/arm64/386/ARMv6/ARMv7; packages as `.tar.gz` (`.zip` for Windows); generates `checksums.txt`; groups changelog by Features / Fixes / Chore
- Added `var version = "dev"` to `main.go` so GoReleaser can inject the release version at build time via `-X main.version={{.Version}}`
- Updated `docs/architecture/07_deployment_view.md` to document the release pipeline
- Updated `docs/technology/STACK.md` to list the CI/CD tools

## Special cases

- The `release-please.yml` and `release.yml` workflows overlap intentionally: `release-please` is the primary path (commit-driven), while `release.yml` provides a manual escape hatch via `workflow_dispatch` or direct tag push.
- No Homebrew tap is configured (unlike phopy) — binaries are distributed exclusively via GitHub Releases.
- GoReleaser version is pinned to `v1.26.2` to match the phopy setup and ensure deterministic builds.
- Go version is read from `go.mod` in all workflows (`go-version-file: go.mod`) rather than hard-coded, so it stays in sync automatically.

## Open items

- None
