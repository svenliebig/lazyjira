# 2026-04-15 — Unassign action

## Goal

Add a direct `u` key binding that unassigns the current user from the selected issue and removes it from the assigned issues list without requiring a modal or confirmation step.

---

## Completed work

### Jira client: `UnassignIssue`
- New method `(c *Client) UnassignIssue(ctx context.Context, key string) error`
- Calls `PUT /rest/api/3/issue/{key}/assignee` with body `{"accountId": null}`
- Expects `204 No Content` on success

### Shared messages
- Added `UnassignDoneMsg struct{}` to `internal/tui/shared/messages.go`

### Key binding: `u`
- Added `KeyUnassign = "u"` constant to `internal/tui/shared/keys.go`
- Handled in `handleKey()` in `app.go`: active when `m.currentIssue != nil && m.jiraClient != nil`
- Triggers `doUnassignCmd()` which calls `jiraClient.UnassignIssue` asynchronously
- On `UnassignDoneMsg`: removes the issue from `m.allIssues`, re-filters via `m.exclusions.Filter`, rebuilds `IssueListModel`, sets `m.statusMsg = "Unassigned"`, and returns to `viewIssueList`

### Status bar
- Added `u:unassign` hint to all contexts where an issue is selected

---

## Special cases

### Immediate removal from list
After a successful unassign the issue is removed from `m.allIssues` so it does not reappear until the next API fetch. Because the assigned-issues JQL query filters by `assignee = currentUser()`, the issue would simply not be returned on the next refresh anyway — removing it locally avoids a stale entry and keeps the UX consistent.

### No confirmation modal
Unassigning is directly reversible in Jira, so no confirmation step is warranted. This matches the behaviour of `o` (open in browser) and the immediate remove-exclusion action on the excluded list view.

### Key choice: `u`
`u` is unused at the top-level key map. It is the natural mnemonic for "unassign" and avoids conflicting with the `y-u` chord (copy URL), which is only active while the copy modal is open.

---

## Open items

- No undo within the TUI — the user must reassign in Jira if the action was accidental.
- Unassigning an issue while on the detail view (`viewIssueDetail`) returns the user to `viewIssueList`. There is no case for staying on the detail view of an issue that is no longer relevant to the user.
