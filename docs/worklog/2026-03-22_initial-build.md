# 2026-03-22 — Initial build, bug fixes, UX overhaul, architecture docs

## Goal

Build the jira-cli project from scratch based on the design documentation in `./docs/`, fix all runtime issues encountered during real usage, implement two UX improvements, and produce full arc42 architecture documentation with ADRs.

---

## Completed work

### Project scaffolding
- Created `go.mod` with module `github.com/svenliebig/jira-cli`, Go 1.24.2
- Dependencies: `charmbracelet/bubbletea v1.3.10`, `charmbracelet/bubbles v1.0.0`, `charmbracelet/lipgloss v1.1.0`, `atotto/clipboard v0.1.4`
- Implemented the full package structure:
  - `internal/config` — credential loading with three-level resolution chain
  - `internal/jira` — REST client, issue listing/fetching, ADF-to-text conversion, transitions
  - `internal/git` — `CommitsForIssue(key)` via `git log --oneline --all`
  - `internal/ollama` — local LLM client (`POST /api/generate`, `stream:false`)
  - `internal/clipboard` — thin wrapper over `atotto/clipboard`
  - `internal/browser` — `open`/`xdg-open`/`start` by `runtime.GOOS`
  - `internal/tui/shared` — shared messages, key constants, styles (circular import break)
  - `internal/tui/views` — `HomeModel`, `IssueListModel` (split-panel), `IssueDetailModel`
  - `internal/tui/modals` — `AuthModal`, `CopyModal`, `AIModal`, `TransitionModal`, `HelpModal`, `ListSelectorModal`
  - `internal/tui/app.go` — root Bubble Tea model, `handleKey()`, state machines

### UX: actions available in list view
All issue actions (`y` copy, `o` open in browser, `a` AI summary, `t` transition) were made available directly in the list view, not only when an issue detail view is open. The `currentIssue` field in the root model tracks the highlighted list item at all times.

### UX: lazygit-style split-panel layout
Replaced the navigate-to-detail-view model with a persistent split-panel:
- Left panel (40%): scrollable issue list
- Right panel (60%): issue detail (description, metadata)
- `Enter` moves focus to the right panel; `Esc` in the right panel moves focus back left
- Border colour changes to indicate which panel is active (`ColorFocus` vs `ColorBorder`)

### Architecture documentation
- arc42 sections 01–12 in `docs/architecture/`
- 10 ADRs in `docs/architecture/decisions/`
- `docs/architecture/README.md` as index

---

## Bug fixes

### 1. Runtime panic on startup — nil pointer in `list.(*Model).updatePagination`

**Root cause:** `updateChildSizes()` in `app.go` called `SetSize` on `IssueListModel` and `IssueDetailModel` during the very first `WindowSizeMsg`, before these models had been constructed. The zero-value `list.Model` inside `IssueListModel` had a nil paginator, causing a nil-pointer dereference.

**Fix:** Added an `initialized bool` field to both `IssueListModel` and `IssueDetailModel`. `SetSize` is a no-op until `initialized = true`, which is set only inside the respective constructors (`NewIssueListModel`, `NewIssueDetailModel`). The `WindowSizeMsg` handler in `app.go` calls `SetSize` freely — the guard inside each model ensures it is safe.

**Why this way:** The `WindowSizeMsg` arrives before `IssueListLoadedMsg` (which triggers construction). Making `SetSize` idempotent on the zero value was cleaner than gating the call in `app.go`, because the calling site doesn't know whether the child has been constructed yet.

**Smoke tests added:** `internal/tui/smoke_test.go` — sends a `WindowSizeMsg` to a fresh model and verifies no panic. These tests catch this class of zero-value initialisation bug.

---

### 2. 403 Forbidden from Jira Cloud API

**Root cause:** Initial implementation used `Authorization: Bearer <api-token>`. Jira Cloud API tokens require HTTP Basic Auth: `Authorization: Basic base64(email:token)`. The email was also not collected anywhere — the config had only `JiraCloudURL` and `JiraAPIToken`.

**Fix:**
- Added `JiraEmail` to `Config`, `Flags`, CLI flags (`--jira-email`), and `AuthCompletedMsg`
- Added an email `textinput` field to `AuthModal`
- Changed `newRequest()` from `req.Header.Set("Authorization", "Bearer "+c.apiToken)` to `req.SetBasicAuth(c.email, c.apiToken)`
- Updated `NewClient` signature: `NewClient(baseURL, email, apiToken string)`
- Updated `issues_test.go` to assert `Basic dGVzdEBleGFtcGxlLmNvbTp0ZXN0LXRva2Vu` in the Authorization header

---

### 3. 410 Gone from Jira API

**Root cause:** `GET /rest/api/3/search?jql=...` was deprecated by Atlassian (see [CHANGE-2046](https://developer.atlassian.com/changelog/#CHANGE-2046)) and returns HTTP 410 Gone.

**Fix:** Changed to `POST /rest/api/3/search/jql` with a JSON request body. Updated the test server to expect `POST` method.

---

### 4. Issue list showing 50 "unassigned" items

**Root cause:** Two compounding issues:
1. The GET form of the JQL query may not have correctly resolved `currentUser()`, returning all issues up to the `maxResults` limit.
2. The new `/search/jql` endpoint does not return all fields by default. Without explicitly requesting `assignee`, the field was absent from all responses — causing every rendered row to show "unassigned".

**Fix:** Switched to `POST /rest/api/3/search/jql` with a JSON body:
```json
{
  "jql": "assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC",
  "maxResults": 50,
  "fields": ["summary", "description", "status", "assignee", "reporter"]
}
```
The `AND statusCategory != Done` filter also removes completed issues from the default view.

---

### 5. `min` redeclared — compilation error

**Root cause:** Both `views/issue_list.go` and `views/issue_detail.go` defined a local `min(a, b int) int` helper function. Go does not allow two functions with the same name in the same package.

**Fix:** Removed the duplicate definition from `issue_list.go`; kept the one in `issue_detail.go`.

---

## Special cases

### Circular import prevention via `internal/tui/shared`

`internal/tui` (app.go) imports `internal/tui/views` and `internal/tui/modals`. If `views` or `modals` imported `internal/tui` to access shared `tea.Msg` types or key constants, the import graph would be cyclic (a compile error in Go).

**Solution:** `internal/tui/shared` is a leaf package with no dependencies on the other TUI packages. It contains:
- All `tea.Msg` struct definitions (`AuthCompletedMsg`, `IssueListLoadedMsg`, etc.)
- Key constant strings (`KeyCopy = "y"`, `KeyAI = "a"`, etc.)
- Shared lipgloss styles and exported colour tokens (`ColorBorder`, `ColorFocus`)

All three packages import `shared`; none import each other. Any new feature that requires a new message type or key binding must add it to `shared` first.

---

### `initialized` guard on `SetSize`

Both `IssueListModel` and `IssueDetailModel` have an `initialized bool` field. `SetSize` silently returns when `initialized == false`. This guard exists because `WindowSizeMsg` is emitted by the Bubble Tea runtime before the application has loaded its first data (and before the constructors for child models have been called). Without the guard, calling `SetSize` on a zero-value `list.Model` panics in `list.(*Model).updatePagination`.

The pattern to follow: any Bubble Tea child model that is not constructed at program start must include this guard.

---

### `currentIssue` kept in sync on every list update

The root model holds `currentIssue *jira.Issue` so that action keys (`y`, `o`, `a`, `t`) work in the list view without navigating to a detail page. This pointer must always reflect the currently highlighted list row.

Two places update it:
1. `IssueListLoadedMsg` handler — sets `currentIssue` to the first issue after the list loads
2. `updateActiveChild()` — called after every `Update` that routes to `issueListView`; reads `m.issueListView.CurrentIssue()` and writes it back to `m.currentIssue`

If a new code path updates the list model without going through `updateActiveChild`, `currentIssue` will drift. The test for this is: press `j`/`k` in the list and then `y→k` — the copied key should match the highlighted row, not the previously selected one.

---

### ESC priority: blur right panel before navigating back

The root `handleKey()` handles `Esc` to navigate back (e.g., from issue list to home). In the split-panel layout, `Esc` inside the right detail panel should blur the panel (return focus to the left list), not navigate away from the issue list entirely.

The ESC branch checks `m.issueListView.IsFocusRight()` first:
```go
case "esc":
    if m.issueListView.IsFocusRight() {
        m.issueListView.BlurRight()
        return m, nil
    }
    // ... normal back-navigation
```

Without this check, pressing `Esc` while reading the detail panel would jump the user all the way back to the home screen.

---

### ADF-to-text handles both ADF objects and plain strings

`adfToText()` in `internal/jira/issues.go` accepts a `json.RawMessage`. Jira Cloud API v3 always returns ADF, but some older self-hosted instances or test fixtures may return a plain string. The function checks the first byte: if `"` it unmarshals as a plain string; if `{` it recursively traverses the ADF tree. This prevents a parse error when running against non-cloud Jira instances or mocked responses.

---

### Split-panel widths: 40/60, not equal halves

The left panel (issue list) is narrower than the right panel (issue detail) deliberately. Issue keys and titles are short; descriptions are long. An equal 50/50 split wastes space on the list and truncates descriptions unnecessarily. The 40/60 ratio (`splitWidths(total int) (left, right int)` in `issue_list.go`) was chosen to balance readability of both panels on a standard 120-column terminal.

---

## Open items

- Pagination (`maxResults` is hardcoded at 50) — TD-05
- Ollama model name is hardcoded as `"llama3"` — TD-03
- Error messages when Ollama is not running give no guidance — R-03
- No automated integration tests against a real Jira instance — TD-04
- `IssueDetailModel` is no longer used as a standalone view (replaced by the right panel in `IssueListModel`) — could be removed or repurposed (TD-06)
