# ADR-001 — Use Bubble Tea as TUI Framework

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

jira-cli is a terminal user interface application. A TUI framework is needed to handle terminal control sequences, keyboard input, layout, and the event loop.

The inspiration source, lazygit, uses [gocui](https://github.com/jesseduffield/gocui) — a lower-level panel management library. Several modern alternatives exist in the Go ecosystem.

## Decision

Use [Bubble Tea](https://github.com/charmbracelet/bubbletea) as the TUI framework, together with [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling and [bubbles](https://github.com/charmbracelet/bubbles) for pre-built components.

## Rationale

**Elm Model-View-Update architecture:**
Bubble Tea enforces unidirectional data flow. All state changes happen through `Update(msg) → (Model, Cmd)`. This makes reasoning about application state deterministic and eliminates a large class of UI bugs caused by mutable shared state.

**Async I/O as first-class citizens:**
`tea.Cmd` provides a structured, framework-managed way to run goroutines and feed results back as messages. There is no need to manage goroutines, channels, or mutexes manually.

**Rich component ecosystem:**
The `bubbles` library provides production-quality `list`, `viewport`, `textinput`, and `spinner` components that would take significant effort to build correctly from scratch (filtering, pagination, scrolling, cursor management).

**Active maintenance and documentation:**
The Charm Bracelet ecosystem (bubbletea + lipgloss + bubbles) is actively maintained, well-documented, and widely adopted in the Go CLI community.

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **gocui** (lazygit's choice) | Lower-level, requires manual event loop management, less structured state management |
| **tview** | Widget-based, imperative API — less suited to the functional/immutable style preferred |
| **tcell** | Raw terminal library; would require building all components from scratch |

## Consequences

- Application structure strictly follows Elm MVU; all state in one root `Model`
- I/O must be wrapped in `tea.Cmd` closures — direct function calls for I/O inside `Update` are prohibited
- Type assertions required when casting `tea.Model` back to concrete types after `Update` calls
- The bubbles components assume specific message types (e.g., `spinner.TickMsg`); these must be routed correctly
