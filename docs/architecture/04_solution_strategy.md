# 4. Solution Strategy

## Fundamental Decisions

### 1. Elm Architecture via Bubble Tea

The application is built on the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework, which implements the Elm Model-View-Update (MVU) pattern in Go. This is the same foundation that lazygit uses via gocui, adapted for a functional, message-driven style.

**Consequences:**
- All UI state lives in a single immutable `Model` struct
- UI changes happen only through `Update(msg) → (Model, Cmd)` — pure functions
- All I/O happens in `Cmd` closures (goroutines managed by the runtime), keeping `Update` side-effect-free
- Reasoning about application state is straightforward: given a model and a message, the next state is deterministic

### 2. State Machine for Navigation

Navigation between screens and modals is modelled as two explicit finite state machines in the root model:

- `viewState`: `Home → IssueList → IssueDetail`
- `modalState`: `None | Auth | Help | ListSelector | Copy | AI | Transition`

Modals are always rendered as overlays on top of the current view. Only one modal can be active at a time. This approach avoids a stack-based navigation model and keeps state transitions predictable.

### 3. Split-Panel Layout for Issue Browsing

The issue list uses a lazygit-style split panel (40% list / 60% detail) rather than separate screens. This allows the user to see issue details while navigating the list, matching the mental model developers already have from tools like lazygit and k9s.

### 4. No External Jira SDK

Jira Cloud REST API v3 calls are made directly using Go's `net/http`. The four API calls needed (search, get issue, get transitions, do transition) do not justify the complexity and dependency overhead of a full SDK. See [ADR-002](./decisions/ADR-002_no_jira_sdk.md).

### 5. Local AI — No Cloud Dependency

AI-assisted work summaries use a locally running Ollama instance. This keeps the feature available without API keys, cloud costs, or sending issue content to third-party services. See [ADR-005](./decisions/ADR-005_local_ollama.md).

## Technology Choices

| Concern | Choice | Reason |
|---------|--------|--------|
| TUI framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) | Elm architecture, active community, Charm ecosystem |
| Styling | [Lipgloss](https://github.com/charmbracelet/lipgloss) | Declarative, composable, pairs with Bubble Tea |
| List component | [bubbles/list](https://github.com/charmbracelet/bubbles) | Free filtering, pagination, keyboard nav |
| Scrolling | [bubbles/viewport](https://github.com/charmbracelet/bubbles) | Efficient large-text scrolling |
| Text input | [bubbles/textinput](https://github.com/charmbracelet/bubbles) | Password mode, cursor, validation |
| Spinner | [bubbles/spinner](https://github.com/charmbracelet/bubbles) | Async loading feedback |
| Clipboard | [atotto/clipboard](https://github.com/atotto/clipboard) | Cross-platform (macOS, Linux, Windows) |
| HTTP | Go standard `net/http` | No dependencies, sufficient for 4 endpoints |
| Config | Go standard `encoding/json` | Config is simple JSON; no YAML library needed |
| AI | Ollama HTTP API | Local-first, model-agnostic |

## How Quality Goals are Addressed

| Quality Goal | Approach |
|-------------|----------|
| Responsiveness | All I/O in `tea.Cmd` goroutines; spinner shown during loading |
| Simplicity | Direct HTTP over SDK; no router; flat state machine |
| Discoverability | Dynamic status bar shows available shortcuts for current state |
| Portability | Standard Go binary; platform-specific code isolated to `browser/open.go` |
| Privacy | Ollama runs locally; credentials stored in XDG config with `0600` permissions |
