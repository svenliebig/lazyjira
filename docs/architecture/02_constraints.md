# 2. Architecture Constraints

## Technical Constraints

| ID | Constraint | Rationale |
|----|-----------|-----------|
| TC-01 | Implemented in **Go** | Chosen for static binaries, cross-platform support, and similarity to the lazygit tech stack |
| TC-02 | TUI framework must be **Bubble Tea** | Consistent with the lazygit-inspired approach; Bubble Tea is the standard for modern Go TUI apps |
| TC-03 | No CGo dependencies | Ensures cross-compilation without a C toolchain |
| TC-04 | Single compiled binary | No runtime dependencies; install by copying the binary |
| TC-05 | Jira **Cloud** REST API v3 only | On-premise Jira (Data Center) is out of scope |
| TC-06 | AI features require a running **local Ollama** instance | No cloud AI provider; Ollama must be running at `localhost:11434` |

## Organisational Constraints

| ID | Constraint | Rationale |
|----|-----------|-----------|
| OC-01 | Config stored in `$XDG_CONFIG_HOME/jira-cli/config.json` | Follows the XDG Base Directory specification for portability and user expectations on Linux/macOS |
| OC-02 | Credentials must never be logged or printed | API tokens are sensitive; echo mode is used for the token input field |
| OC-03 | No telemetry or analytics | Privacy-first; no data collection |

## Conventions

| ID | Convention |
|----|-----------|
| CV-01 | All packages placed under `internal/` to prevent external imports |
| CV-02 | Message types defined in `internal/tui/shared/` to avoid circular imports between TUI packages |
| CV-03 | Asynchronous operations always wrapped as `tea.Cmd`; no goroutines launched directly |
| CV-04 | Go standard `context.Context` threaded through all I/O calls for cancellation support |
