# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build          # Build the binary
go test ./...     # Run all tests
go test ./internal/jira/...  # Run a single package's tests
```

## Architecture

**lazyjira** is a terminal UI (TUI) for Jira Cloud, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) (event-driven state machine) and styled with Lipgloss.

### Entry Point & Composition Root

`main.go` loads config, settings, custom themes, and exclusions, then instantiates the Jira HTTP client and starts the Bubble Tea event loop.

### Package Structure

| Package | Role |
|---|---|
| `internal/tui/` | Root Bubble Tea model (`app.go`), views, modals, shared types |
| `internal/tui/views/` | Full-screen panels: issue list (split), issue detail |
| `internal/tui/modals/` | Overlays: auth, help, settings, AI summary, transition, exclusions, etc. |
| `internal/tui/shared/` | Message types (`messages.go`) and Lipgloss styles (`styles.go`) — shared to prevent circular imports (ADR-010) |
| `internal/jira/` | Stateless HTTP client for Jira Cloud REST API v3 (Basic Auth, ADR-002/003) |
| `internal/config/` | Three-level credential resolution: config file → env vars → CLI flags |
| `internal/theme/` | Predefined themes (default, Dracula, Nord, 4 Catppuccin flavors) + custom theme loading |
| `internal/settings/` | User preferences persistence (currently active theme) |
| `internal/exclusions/` | Client-side issue filtering rules |
| `internal/git/`, `internal/ollama/`, `internal/clipboard/`, `internal/browser/` | Thin integration wrappers |

### Message Bus Pattern

All cross-component communication uses typed message structs defined in `internal/tui/shared/messages.go`. The root model in `app.go` routes messages to the active view or modal. Async operations (API calls, file I/O) return `tea.Cmd` values.

### Theme System

Themes are defined in `internal/theme/theme.go`. `SetTheme()` updates the global theme; styles in `internal/tui/shared/styles.go` are rebuilt on theme change. Custom themes are loaded from `~/.config/lazyjira/themes.json`.

### Configuration Files

All stored under `~/.config/lazyjira/` (XDG):
- `config.json` — Jira URL, email, API token
- `settings.json` — Active theme
- `exclusions.json` — Filtered issue rules
- `themes.json` — User-defined custom themes

## Key Design Decisions

- **No Jira SDK** — direct HTTP calls to avoid heavy dependencies (ADR-002)
- **Two-key chord system** for sub-actions (e.g. `e` then `s` for exclude by status) (ADR-008)
- **ADF to plain text** conversion for Jira's Atlassian Document Format description bodies (ADR-009)
- **Local Ollama** for AI summaries — no external AI service dependency (ADR-005)
- Architecture decisions documented in `docs/architecture/decisions/`
