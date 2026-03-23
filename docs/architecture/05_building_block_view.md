# 5. Building Block View

## Level 1 — Top-level Packages

```
┌──────────────────────────────────────────────────────────────────────┐
│                             lazyjira                                  │
│                                                                       │
│  ┌────────┐  ┌────────────┐  ┌──────────┐  ┌────────┐  ┌────────┐  │
│  │ config │  │ exclusions │  │   jira   │  │  tui   │  │ integr │  │
│  └────────┘  └────────────┘  └──────────┘  └────────┘  └────────┘  │
│                                                                       │
│  ┌──────────┐  ┌──────────┐                                          │
│  │  theme   │  │ settings │                                          │
│  └──────────┘  └──────────┘                                          │
│                                                                       │
│  main.go — composition root                                           │
└──────────────────────────────────────────────────────────────────────┘
```

### main.go — Composition Root

Parses CLI flags, loads config, loads settings, loads custom themes, applies the saved theme, constructs a `jira.Client` if credentials are available, and hands everything to `tui.New()` before starting the Bubble Tea event loop.

### internal/config

Responsible for the full credential resolution chain: config file → environment variables → CLI flags (each level overrides the previous). Owns the JSON schema for `~/.config/lazyjira/config.json` and the `IsComplete()` check.

### internal/theme

Owns the `Theme` struct, the list of predefined themes, the `Current` global variable, `SetTheme()`, and `LoadCustom()` which reads user-defined themes from `~/.config/lazyjira/themes.json`. Has no knowledge of the TUI. Used by `shared.RefreshStyles()` to rebuild all lipgloss styles when the theme changes.

### internal/settings

Owns the `Settings` struct (`ActiveTheme string`) and `Load()`/`Save()` backed by `~/.config/lazyjira/settings.json`. Follows the same XDG path convention as `config` and `exclusions`. Has no knowledge of the TUI.

### internal/exclusions

Manages the user's personal list of exclusion rules. Persists `[]Rule` as JSON to `~/.config/lazyjira/exclusions.json`. Exposes `Add`, `Remove`, `Rules`, and `Filter` — the last of which removes matching issues from a `[]jira.Issue` slice in memory. Has no knowledge of the TUI or the Jira API beyond the domain types it filters.

### internal/jira

Stateless HTTP client for Jira Cloud REST API v3. Encapsulates authentication, JSON parsing (including ADF-to-text conversion), and error mapping. Has no knowledge of the TUI.

### internal/tui

The entire terminal UI. Composed of three sub-packages:
- `shared/` — message types, key constants, and styles shared across views and modals
- `views/` — full-screen composable view models
- `modals/` — overlay modal models

### Integrations (four small packages)

| Package | Responsibility |
|---------|---------------|
| `internal/git` | Runs `git log` to find commits for an issue key |
| `internal/ollama` | POSTs to local Ollama HTTP API for text generation |
| `internal/clipboard` | Thin wrapper around `atotto/clipboard` |
| `internal/browser` | OS-agnostic URL opener |

---

## Level 2 — TUI Package

```
internal/tui/
│
├── app.go                  Root model (state machine, message router)
│
├── shared/
│   ├── messages.go         All tea.Msg types (domain events)
│   ├── keys.go             Key constants
│   └── styles.go           Lipgloss styles rebuilt by RefreshStyles() from theme.Current
│
├── views/
│   ├── home.go             Home/landing screen (stateless)
│   ├── issue_list.go       Split-panel: list (left) + detail (right)
│   └── issue_detail.go     Full-screen issue detail (viewport)
│
└── modals/
    ├── modal.go            Wrap() helper — renders any modal in a styled box
    ├── auth_modal.go       First-run credential entry (3 textinput fields)
    ├── help_modal.go       Keyboard shortcut reference overlay
    ├── list_selector.go    "l" sub-menu: choose which list to load
    ├── copy_modal.go       "y" sub-menu: copy key / URL / title / description
    ├── transition_modal.go "t" sub-menu: numbered transition picker
    ├── ai_modal.go         "a→s" AI summary with spinner and viewport
    ├── exclude_modal.go    "x" sub-menu: exclude by key or parent
    └── settings_modal.go   "s" settings screen: theme selector
```

### Root Model (`app.go`)

The root model holds the complete application state and routes every message to the correct child. It is the only component that knows about all other components.

**Responsibilities:**
- Owns `viewState` and `modalState` enums and their transitions
- Tracks `currentIssue` (the issue all actions operate on)
- Tracks `allIssues` (raw API results) separately from the filtered display list
- Dispatches keyboard input: modal-first, then chord resolution, then global keys, then active view
- Fires async commands (`tea.Cmd`) for all I/O operations
- Renders header, content area, status bar, and modal overlay

**Does NOT:**
- Know the internal structure of any modal or view
- Perform I/O directly
- Contain rendering logic for views or modals

### Views

Each view implements `Init() tea.Cmd`, `Update(tea.Msg) (tea.Model, tea.Cmd)`, and `View() string`.

| View | Description |
|------|-------------|
| `HomeModel` | Static welcome screen. No state. |
| `IssueListModel` | Split-panel view. Manages a `bubbles/list` on the left and a `bubbles/viewport` on the right. Tracks `focusRight bool` to route keyboard input. Exposes `CurrentIssue()`, `IsFocusRight()`, `BlurRight()` for root model coordination. |
| `IssueDetailModel` | Full-screen scrollable viewport for a single issue. Retained for potential direct navigation; currently reached only programmatically. |
| `ExcludedListModel` | Full-width list of active exclusion rules. Each item shows the rule type (`key` or `parent`) and value. Exposes `CurrentRule()` so the root model can pass the highlighted rule to `exclusions.Store.Remove()`. |

### Modals

Each modal emits one of the message types in `shared/messages.go` when closed or when an action is confirmed. Modals are stateless where possible; state is introduced only when required (auth inputs, AI generation state).

| Modal | Emits |
|-------|-------|
| `AuthModal` | `AuthCompletedMsg` or `CloseModalMsg` |
| `HelpModal` | `CloseModalMsg` |
| `ListSelectorModal` | `ListSelectedMsg{Type}` or `CloseModalMsg` |
| `CopyModal` | `CopyActionMsg{Action}` or `CloseModalMsg` |
| `TransitionModal` | `TransitionSelectedMsg{ID}` or `CloseModalMsg` |
| `AIModal` | Internally handles `AICommitsLoadedMsg`, `AISummaryMsg`; emits `CloseModalMsg` on close |
| `ExcludeModal` | `ExcludeActionMsg{Type, Value}` or `CloseModalMsg`. The `p` (parent) option is visually strikethrough and non-functional when the current issue has no parent. |
| `SettingsModal` | `ThemeSelectedMsg{Name}` or `CloseModalMsg`. Receives the full theme list (predefined + custom) and the name of the currently active theme from the root model. |

---

## Level 3 — IssueListModel (Split Panel)

```
IssueListModel
│
├── list.Model (bubbles/list)        ← left panel, ~40% width
│   └── []issueItem                  ← implements list.Item
│       ├── Title()   → "[KEY]  Summary text"
│       ├── Description() → "Status · Assignee"
│       └── FilterValue() → "KEY Summary"
│
├── viewport.Model (bubbles/viewport) ← right panel, ~60% width
│   └── content = buildIssueDetail()
│       ├── Key + Summary
│       ├── Status (bordered badge)
│       ├── Assignee / Reporter
│       ├── Horizontal divider
│       └── Description (ADF→text)
│
└── focusRight bool
    ├── false → key events go to list.Model (j/k navigate, Enter toggles)
    └── true  → key events go to viewport.Model (j/k scroll, Esc returns)
```

**Width calculation:**
```
total width
├── left content: total × 2/5       (≈40%)
├── border:       1 char
└── right content: total - left - 1 (≈59%)
```

---

## Level 3 — AIModal State Machine

```
AIModal
│
├── state: aiIdle
│   └── key "s" → fetchCommitsCmd() + spinner.Tick
│       ↓
├── state: aiLoadingCommits
│   └── AICommitsLoadedMsg → SetCommits() → GenerateCmd()
│       ↓
├── state: aiGenerating
│   └── AISummaryMsg → SetSummary()
│       ↓
└── state: aiDone
    └── viewport scrolling active
```

`fetchCommitsCmd()` calls `git.CommitsForIssue(issueKey)`.
`GenerateCmd()` builds a prompt combining issue metadata and commit messages, then calls `ollama.Client.Generate()`.

---

## Package Dependency Graph

```
main
 └─► config
 └─► exclusions
 │    └─► jira (types only)
 └─► jira
 └─► settings
 └─► theme
 └─► tui/shared  (RefreshStyles called at startup)
 └─► tui/app
      └─► tui/shared
      │    └─► theme
      └─► tui/views
      │    └─► tui/shared
      │    └─► jira (types only)
      │    └─► exclusions
      └─► tui/modals
      │    └─► tui/shared
      │    └─► jira (types only)
      │    └─► theme
      │    └─► git
      │    └─► ollama
      └─► config
      └─► exclusions
      └─► settings
      └─► theme
      └─► jira
      └─► clipboard
      └─► browser
```

No circular dependencies. `tui/shared` is the only package imported by both `tui/views` and `tui/modals`, and it imports nothing from the `tui` tree itself. `theme` and `settings` are leaf packages with no dependencies on the TUI or on each other.
