# 2026-03-22 — Exclude feature

## Goal

Implement an exclude feature that lets users hide issues from the list based on criteria (issue key or parent issue). Add a dedicated view to list and manage active exclusion rules.

---

## Completed work

### New package: `internal/exclusions`
- `Rule{Type, Value}` — a single exclusion criterion; `Type` is `"key"` or `"parent"`, `Value` is the relevant issue key
- `Store` — holds the active rules and persists them to `~/.config/lazyjira/exclusions.json` (same XDG base as config)
- `Add(Rule)` — appends a rule (deduplicates) and saves
- `Remove(Rule)` — removes a rule and saves
- `Rules()` — returns a copy of all active rules
- `Filter([]jira.Issue)` — removes issues matching any active rule from a slice

### Jira: parent issue support
- Added `IssueParent{Key}` type and `IssueFields.Parent *IssueParent`
- Added `parentResponse` to `issueFieldsResponse` and `"parent"` to the JQL fields list in `ListAssigned`
- `convertIssue` maps parent key from API response to domain type

### New modal: `ExcludeModal`
- `p` — Exclude all issues by parent key; shows `(PARENT-KEY)` in the label when available
- `k` — Exclude the current issue by its own key
- When the issue has no parent, the `p` option is rendered with strikethrough and muted colour and the key press is a no-op

### New view: `ExcludedListModel`
- Full-width `list.Model` showing all active exclusion rules
- Each item shows the rule value (styled as issue key) and a description of its type
- `CurrentRule()` returns the highlighted rule for removal
- Empty state renders a muted "No exclusions configured." message

### List selector
- Added `x` → "Excluded Issues" to the `ListSelectorModal`; emits `ListSelectedMsg{Type: "excluded"}`

### App wiring
- New `viewExcludedList` view state and `modalExclude` modal state
- `allIssues []jira.Issue` stored in root model — raw API results kept separate from the filtered display list
- `exclusions *exclusions.Store` field on root model; loaded at startup in `main.go`
- `IssueListLoadedMsg` now filters through `exclusions.Filter` before building `IssueListModel`
- `ExcludeActionMsg` handler: adds rule to store, re-filters `allIssues`, rebuilds `IssueListModel` in place
- `ListSelectedMsg{Type:"excluded"}` handler: builds `ExcludedListModel` from current rules, switches to `viewExcludedList`, clears `currentIssue`
- ESC from `viewExcludedList` returns to home
- `updateChildSizes` updated to call `excludedListView.SetSize`
- Status bar shows `x:exclude` hint when an issue is selected; shows `x:remove` hint in excluded list view

### Key binding: `x`
- In issue list (or detail): opens `ExcludeModal` for the current issue
- In excluded list: removes the highlighted rule immediately (no confirmation modal)

---

## Special cases

### `x` is context-aware — same key, two behaviours
The `KeyExclude` handler in `handleKey()` checks `m.currentView` before deciding what to do:
- `viewExcludedList` → remove the highlighted rule
- anything else → open the exclude modal for the current issue

This mirrors the user's stated intent ("in the list of excluded issues, press x to remove") without needing a separate key binding. The status bar always reflects the correct meaning of `x` for the active view so the user is never surprised.

### `allIssues` separates raw data from filtered display
When issues are fetched from the API, they are stored in `m.allIssues` (unfiltered). The `IssueListModel` is always built from `exclusions.Filter(m.allIssues)`. This means:
- Adding an exclusion does not require a network round-trip — the filter is reapplied locally
- If the user later removes an exclusion from the excluded list and returns to the issue list (via `l → a`), the re-fetch will again apply the current (smaller) rule set

The trade-off: if a previously excluded issue is modified in Jira after exclusion, it won't appear until the user manually refreshes (re-selects the assigned list). This is acceptable for the current use case.

### Strikethrough for unavailable option
The "Exclude by parent" option in `ExcludeModal` uses `lipgloss.NewStyle().Strikethrough(true)` when the issue has no parent. The key press (`p`) is also silently ignored in `Update`. This gives both visual and functional unavailability without adding a separate "disabled" state machine to the modal.

### Immediate persistence
`Store.Add` and `Store.Remove` write to disk synchronously before returning. There is no buffered/batched save. Given the low frequency of exclusion changes (not a hot path), this is simpler than a deferred save and ensures rules survive unexpected exits.

### Empty store is always valid
`exclusions.Load()` returns a non-nil `*Store` even on error (e.g. JSON parse failure). `main.go` logs the warning but continues with an empty store. This means the app always starts, and the user loses stored exclusions in the rare case of a corrupted file — preferable to a startup crash.

---

## Open items

- No confirmation when removing an exclusion — `x` removes immediately. Could be revisited if users accidentally remove rules.
- Excluding an issue by key hides it permanently until manually unexcluded; there is no "exclude until resolved" or time-based rule type.
- The excluded list shows raw rules (key strings), not live issue titles. If a project key changes, the rule becomes stale silently.
