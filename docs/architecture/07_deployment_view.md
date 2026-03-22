# 7. Deployment View

## Infrastructure

lazyjira is a single compiled binary with no runtime dependencies. It runs entirely on the developer's local machine.

```
Developer's Workstation
│
├── lazyjira binary          (anywhere in PATH, e.g. /usr/local/bin/lazyjira)
│
├── ~/.config/lazyjira/
│   └── config.json          (mode 0600 — user read/write only)
│       {
│         "jiraCloudUrl":  "https://company.atlassian.net",
│         "jiraEmail":     "user@company.com",
│         "jiraApiToken":  "ATATT3xFfGF0..."
│       }
│
├── Ollama (optional)        (localhost:11434 — required only for AI features)
│
└── Git                      (any version — required only for AI features)
```

## Release Pipeline

Releases are automated via GitHub Actions using a two-stage pipeline:

1. **Release Please** (`release-please.yml`) — monitors `main` for conventional commits (`feat:`, `fix:`, etc.), maintains a release PR with an auto-generated `CHANGELOG.md`, and bumps the version in `.release-please-manifest.json` when merged.
2. **GoReleaser** — triggered automatically when Release Please merges a release PR. Builds multi-platform binaries and attaches them to the GitHub release.

Versioning follows [Semantic Versioning](https://semver.org/). The version is injected into the binary at build time via `ldflags` (`-X main.version=<tag>`).

## Installation

Pre-built binaries are available on the [GitHub Releases page](https://github.com/svenliebig/lazyjira/releases). Download the archive for your platform and place the binary in `$PATH`.

Alternatively, install from source:

```
go install github.com/svenliebig/lazyjira@latest
# — or —
go build -o lazyjira . && mv lazyjira /usr/local/bin/
```

## Configuration File Location

The config path follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html):

```
$XDG_CONFIG_HOME/lazyjira/config.json
```

If `XDG_CONFIG_HOME` is not set, defaults to:
```
~/.config/lazyjira/config.json
```

The directory is created with `0700` and the file is written with `0600` permissions. Credentials are never stored in a world-readable location.

## Environment Variables

| Variable | Purpose |
|----------|---------|
| `JIRA_CLOUD_URL` | Jira instance base URL (overrides config file) |
| `JIRA_EMAIL` | Jira account email (overrides config file) |
| `JIRA_API_TOKEN` | API token from Atlassian (overrides config file) |
| `XDG_CONFIG_HOME` | Base directory for config file location |

## Runtime Requirements

| Requirement | When Required | Notes |
|-------------|---------------|-------|
| Network access to Jira Cloud | Always | HTTPS on port 443 |
| `git` in PATH | AI summary feature only | Any recent version |
| Ollama at `localhost:11434` | AI summary feature only | Model: `llama3` by default |
| Terminal with ANSI colour support | Always | Any modern terminal emulator |

## Supported Platforms

| Platform | Clipboard | Browser open | Notes |
|----------|-----------|--------------|-------|
| macOS | `pbcopy`/`pbpaste` via `atotto/clipboard` | `open` | Primary development platform |
| Linux | `xclip` or `xsel` required | `xdg-open` | `xclip` must be installed |
| Windows | Windows API via `atotto/clipboard` | `start` | Supported but not primary |
