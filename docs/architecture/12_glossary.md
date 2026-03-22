# 12. Glossary

| Term | Definition |
|------|-----------|
| **ADF** | Atlassian Document Format — the structured JSON format Jira uses to represent rich-text content such as issue descriptions. jira-cli converts ADF to plain text for terminal display. |
| **API Token** | A personal access token generated at `id.atlassian.net`. Used in place of a password for Jira Cloud REST API authentication. Combined with the user's email address and Base64-encoded as a Basic Auth credential. |
| **Basic Auth** | HTTP authentication scheme where credentials are sent as `Authorization: Basic base64(username:password)`. For Jira Cloud API tokens, `username` is the account email and `password` is the API token. |
| **Bubble Tea** | A Go framework for building terminal user interfaces, based on the Elm Model-View-Update architecture. Used as the application framework for jira-cli. |
| **bubbles** | A companion library to Bubble Tea providing pre-built TUI components: `list`, `viewport`, `textinput`, `spinner`, etc. |
| **Chord** | A two-key keyboard sequence used to access a sub-menu of actions. For example, `y` then `k` copies the issue key, and `a` then `s` generates an AI summary. |
| **Command (`tea.Cmd`)** | In Bubble Tea, a function that runs asynchronously (in a goroutine) and returns a `tea.Msg`. Used for all I/O: HTTP calls, subprocess execution, clipboard writes. |
| **Elm Architecture** | A unidirectional data-flow pattern (Model–View–Update) originating from the Elm language. State is immutable; changes happen only through messages. Bubble Tea implements this pattern in Go. |
| **Focus** | Which UI element receives keyboard input. In the split-panel view, focus is either on the left list or the right detail panel, tracked by `focusRight bool`. |
| **JQL** | Jira Query Language — a SQL-like syntax for searching issues in Jira. Example: `assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC`. |
| **Lipgloss** | A Go library for declaring terminal styles (colours, borders, padding, alignment) using a CSS-inspired API. Used for all visual styling in jira-cli. |
| **Modal** | An overlay rendered on top of the current view. In jira-cli, modals are used for authentication, help, list selection, copy actions, transitions, and AI assistance. Only one modal can be active at a time. |
| **Model** | In Bubble Tea, the immutable struct holding all application state at a given point in time. In jira-cli, `tui.Model` is the root model. |
| **Message (`tea.Msg`)** | An event in the Bubble Tea event loop. Can be a key press, window resize, or the result of an async command. The `Update` function receives messages and returns a new model and optional command. |
| **Ollama** | An open-source runtime for running large language models locally. jira-cli uses Ollama (default model: `llama3`) to generate AI work summaries without sending data to cloud services. |
| **Split-panel** | A layout where the screen is divided vertically into two panels — a list on the left and details on the right — allowing the user to browse and read simultaneously. Inspired by lazygit's panel layout. |
| **Status bar** | The bottom row of the TUI showing context-aware keyboard shortcut hints. Updated by `renderStatusBar()` based on the current view state, modal state, and pending chord key. |
| **Transition** | In Jira, a workflow action that moves an issue from one status to another (e.g., "In Progress" → "In Review"). Transitions are fetched from the Jira API and presented in the transition modal. |
| **TEA** | The Elm Architecture as implemented in Go by Bubble Tea: **T**he **E**lm **A**rchitecture. |
| **Viewport** | A `bubbles` component that renders a portion of a long text string with scroll support. Used for issue descriptions and AI summaries. |
| **viewState** | An enum in the root model tracking which full-screen view is currently displayed: `viewHome`, `viewIssueList`, or `viewIssueDetail`. |
| **modalState** | An enum in the root model tracking which modal overlay is currently active: `modalNone`, `modalAuth`, `modalHelp`, `modalListSelector`, `modalCopy`, `modalAI`, or `modalTransition`. |
| **XDG** | The [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html), a Linux standard for where user configuration, data, and cache files should be stored. jira-cli follows XDG for its config file location. |
