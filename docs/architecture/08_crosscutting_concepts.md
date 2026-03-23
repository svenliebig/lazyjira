# 8. Cross-cutting Concepts

## Error Handling

All errors are propagated to the root model via the `shared.ErrMsg` message type. Components never crash the application on error; they return an `ErrMsg` from their `tea.Cmd` closure.

```go
// Pattern used consistently across all async commands:
return func() tea.Msg {
    result, err := someOperation()
    if err != nil {
        return shared.ErrMsg{Err: fmt.Errorf("context: %w", err)}
    }
    return SomeSuccessMsg{Data: result}
}
```

The root model's `Update` function catches `ErrMsg` and sets `m.err`, which is rendered in the content area. The user can then dismiss or navigate away.

HTTP errors include the response status code and body to aid debugging:
```go
if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
}
```

---

## Async I/O Pattern

All operations that involve I/O (HTTP, subprocess, local file system) are executed as `tea.Cmd` functions — closures that run in a goroutine managed by the Bubble Tea runtime. The UI remains interactive during any pending operation.

Loading states are indicated via `m.loading = true`, which renders "Loading..." in the content area. Where multiple messages are expected (commits + AI generation), `tea.Batch(cmd1, cmd2)` is used.

The rule: **`Update()` is always synchronous and side-effect-free. I/O belongs in commands.**

---

## Context and Cancellation

All I/O operations accept a `context.Context`:

```go
func (c *Client) ListAssigned(ctx context.Context) ([]Issue, error)
func (c *Client) GetTransitions(ctx context.Context, key string) ([]Transition, error)
func (c *ollama) Generate(ctx context.Context, prompt string) (string, error)
```

Commands create a `context.Background()` at the call site. This provides a hook for future cancellation (e.g., abort on navigation away) without changing function signatures.

---

## Styling and Theming

All visual styles are defined centrally in `internal/tui/shared/styles.go` as exported `lipgloss.Style` variables. No inline style literals exist outside this file.

Styles are **dynamic**: `RefreshStyles()` rebuilds all variables from `theme.Current`. The function is called once at startup (via `init()`) and again whenever the user selects a new theme. This means all call sites (`shared.StyleHeader.Render(…)` etc.) remain unchanged.

**Theme struct fields** (defined in `internal/theme`):

| Field | Default (`default` theme) | Use |
|-------|---------------------------|-----|
| `primary` | `#7C3AED` | Headings, key hints, focused borders, issue keys |
| `success` | `#10B981` | Success messages, status badges |
| `error` | `#EF4444` | Error messages |
| `muted` | `#9CA3AF` | Labels, separators, secondary text |
| `bg` | `#1F2937` | Modal backgrounds |
| `surface` | `#374151` | Status bar background |
| `text` | `#F9FAFB` | Primary text |
| `subtext` | `#D1D5DB` | Secondary text (status bar) |
| `border` | `#4B5563` | Inactive panel borders |
| `focus` | `#7C3AED` | Active panel borders |

`ColorBorder` and `ColorFocus` are exported from `shared` for use in view packages that need to apply border colours programmatically; they are also updated by `RefreshStyles()`.

**Predefined themes:** `default`, `dracula`, `nord`, `catppuccin-mocha`, `catppuccin-macchiato`, `catppuccin-frappe`, `catppuccin-latte`.

**Custom themes** can be added to `~/.config/lazyjira/themes.json` as a JSON array of `Theme` objects. They appear after the predefined themes in the settings modal.

---

## Keyboard Handling

Key constants are defined in `internal/tui/shared/keys.go` and referenced by name throughout the codebase to avoid magic strings.

**Dispatch order in `handleKey()`:**

1. `ctrl+c` — unconditional quit
2. If any modal is active — delegate to modal (`updateActiveChild`)
3. `esc` — contextual back/close (checks pending key, focus state, view stack)
4. Chord resolution (e.g., `pendingKey == "y"` + new key)
5. Global key switch (l, ?, q, y, o, a, t)
6. Delegate to active view

This ordering ensures that:
- Modals always intercept input first
- `esc` always works as a cancel/back
- Global shortcuts don't fire inside a modal

**Chord system:**

Two-key sequences are implemented via `pendingKey string` in the root model:

```
Press "y" → pendingKey = "y", open copy modal
Press "k" → if pendingKey == "y": copy key; pendingKey = ""
Press "esc" → pendingKey = ""; close modal
```

The copy modal and AI modal both provide visual feedback for which second key is expected.

---

## Status Bar

The status bar (`renderStatusBar()`) is the primary discoverability mechanism. Its content adapts to the current application state:

| State | Shown hints |
|-------|-------------|
| Home / List view | `l:list  s:settings  ?:help  q:quit` |
| List view, left focus | `l:list  s:settings  ?:help  q:quit  enter:focus detail  [action keys if issue selected]` |
| List view, right focus | `j/k:scroll  esc:back  [action keys]` |
| Excluded list view | `j/k:navigate  x:remove  esc:back  ?:help` |
| Copy chord active | `↑/k ↓/j:navigate  enter/l:select  u:url  t:title  d:desc  h/esc:cancel` |
| AI chord active | `s:summary  esc:cancel` |

Action keys (`o:open  y:copy  t:transition  a:AI  x:exclude`) appear whenever `currentIssue != nil`.

---

## Config Persistence

All persistent files live in `$XDG_CONFIG_HOME/lazyjira/` (default `~/.config/lazyjira/`). Every writer follows the same pattern: `os.MkdirAll(dir, 0700)` then `os.WriteFile(path, data, 0600)`. This ensures:
- Directory is created if absent
- Only the owning user can read or write the files

| File | Package | Contents |
|------|---------|---------|
| `config.json` | `internal/config` | Jira credentials (URL, email, token) |
| `exclusions.json` | `internal/exclusions` | Array of exclusion rules |
| `settings.json` | `internal/settings` | User preferences (`activeTheme`) |
| `themes.json` | `internal/theme` | Optional array of custom `Theme` objects |

All files use `encoding/json` with `json.MarshalIndent`. No migration mechanism is needed at this scale; unknown fields are ignored on read.

---

## Testing Strategy

The project uses Go's standard `testing` package. Key test patterns:

- **`internal/config`**: Table-driven tests using `t.Setenv()` to isolate environment variables, with `t.TempDir()` to avoid touching the real config file.
- **`internal/jira`**: `httptest.NewServer` creates an in-process HTTP server. Tests verify both the correct request format (method, auth header) and correct response parsing.
- **`internal/tui`**: Smoke tests that instantiate the root model and simulate `WindowSizeMsg` and key presses, catching panics and verifying non-empty output without a real terminal.

No mocking framework is used; test doubles are achieved via `httptest` or by passing `nil` for optional dependencies (e.g., `jiraClient = nil` to test the auth-required state).
