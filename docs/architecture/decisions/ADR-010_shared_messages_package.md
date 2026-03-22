# ADR-010 — Shared Sub-package to Prevent Circular Imports

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

The TUI layer is split into three Go packages:

| Package | Role |
|---------|------|
| `internal/tui` | Root model (`app.go`), top-level `Update`/`View` |
| `internal/tui/views` | `IssueListModel`, `IssueDetailModel`, `HomeModel` |
| `internal/tui/modals` | `AuthModal`, `CopyModal`, `AIModal`, `TransitionModal`, etc. |

All three packages need to share:
- **`tea.Msg` types** — e.g. `AuthCompletedMsg`, `IssueListLoadedMsg`, `CopyActionMsg` — so that `views` and `modals` can emit messages that `app.go` can match in its `Update` switch
- **Key constants** — e.g. `KeyCopy = "y"`, `KeyOpen = "o"` — so that both the status bar renderer in `app.go` and the chord handler in `views` use identical string literals
- **Lipgloss styles** — e.g. `ColorBorder`, `ColorFocus` — so that `views/issue_list.go` and `app.go` apply consistent colours without duplicating declarations

The naive approach — defining these in `internal/tui` and importing them from `views` and `modals` — creates a circular import: `internal/tui` imports `internal/tui/views`, so `internal/tui/views` cannot import `internal/tui`.

## Decision

Extract all shared definitions into a dedicated sub-package `internal/tui/shared/`:

```
internal/tui/shared/
    keys.go       — key constant strings
    messages.go   — all tea.Msg type definitions
    styles.go     — shared lipgloss styles and exported colour tokens
```

All three packages (`internal/tui`, `internal/tui/views`, `internal/tui/modals`) import `internal/tui/shared`. None of them import each other.

## Rationale

**Breaks the cycle without restructuring the domain:**
The circular dependency exists because Go's import graph must be a DAG. Introducing a leaf package that has no dependencies on the other three packages resolves the cycle with the smallest possible structural change.

**One canonical definition per constant:**
Before the shared package, key strings and style values were duplicated or defined in `app.go` and passed down via constructor arguments. A single source of truth eliminates drift (e.g., a key constant changing in one file but not another).

**Cohesion of the shared package:**
All three categories of shared content (messages, keys, styles) are tightly coupled — a new feature typically requires a new `tea.Msg`, a new key constant, and possibly a new style. Keeping them in one package makes feature additions self-contained.

**No business logic in shared:**
`internal/tui/shared` contains only pure data: constants, struct definitions, and style declarations. It has no behaviour, no `Update` logic, and no dependencies on `bubbletea` beyond the `tea.Msg` interface (an empty interface). This keeps the leaf package stable and prevents the cycle from re-emerging.

## Package structure

```go
// shared/messages.go
type AuthCompletedMsg struct { URL, Email, Token string }
type IssueListLoadedMsg struct { Issues []jira.Issue }
type IssueSelectedMsg  struct { Issue *jira.Issue }
type CopyActionMsg     struct { Action string }
type AICommitsLoadedMsg struct { Commits []string }
type AISummaryMsg      struct { Summary string }
type ErrMsg            struct { Err error }
type CloseModalMsg     struct{}
// ...

// shared/keys.go
const (
    KeyHelp       = "?"
    KeyList       = "l"
    KeyCopy       = "y"   // chord initiator
    KeyOpen       = "o"
    KeyAI         = "a"   // chord initiator
    KeyTransition = "t"
    KeyVimUp      = "k"
    KeyVimDown    = "j"
    // ...
)

// shared/styles.go
var (
    ColorBorder = colorBorder   // exported for views
    ColorFocus  = colorFocus    // exported for views
    // lipgloss style declarations...
)
```

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **Define messages in `internal/tui`, pass via constructor** | Constructors become unwieldy; views/modals cannot reference type names for type assertions |
| **Define messages in `views`, import from `modals` and `app`** | Still creates a cycle if `app` imports both `views` and `modals` |
| **Single flat package for all TUI code** | Eliminates the cycle by collapsing boundaries, but removes the separation of concerns between layout, modals, and orchestration |
| **Interface-based message passing** | Would require every consumer to define its own interface; no shared type names for switch-case matching in `app.go` |

## Consequences

- `internal/tui/shared` is a stable leaf package; adding to it does not affect the import graph
- All new `tea.Msg` types must be defined in `shared/messages.go` — `views` and `modals` may not define their own message types
- New key bindings require a constant in `shared/keys.go` before use; this prevents magic strings across the codebase
- The package name `shared` is intentionally generic — it is an internal coordination package, not a public API
