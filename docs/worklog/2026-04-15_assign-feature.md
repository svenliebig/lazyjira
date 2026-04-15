# 2026-04-15 — Assign feature & AI key remapping

## Goal

Add an `a` key binding that opens a fuzzy-searchable user picker to assign the current issue to any assignable user. Relocate the existing AI chord from `a` to `m` to free up `a` for the new action.

---

## Completed work

### Key remapping: AI `a` → `m`
- Changed `KeyAI = "a"` to `KeyAI = "m"` in `internal/tui/shared/keys.go`
- Added `KeyAssign = "a"` to the same file
- Status bar hints updated from `a:AI` to `m:AI` and `a:assign`

### Jira types: `AccountID` on `User`
- Added `AccountID string` field to `User` in `internal/jira/types.go`
- Added `accountId` JSON tag to `userResponse` in `internal/jira/issues.go`
- `convertIssue` now maps `AccountID` for both `Assignee` and `Reporter`

### Jira client: `SearchAssignableUsers` and `AssignIssue`
- `SearchAssignableUsers(ctx, issueKey) ([]User, error)` — calls `GET /rest/api/3/user/assignable/search?issueKey=<key>&maxResults=50`; returns users scoped to the specific issue, which honours project-level permission schemes
- `AssignIssue(ctx, key, accountID string) error` — calls `PUT /rest/api/3/issue/{key}/assignee` with body `{"accountId": "<id>"}`; expects `204 No Content`

### Shared messages
- Added `UsersLoadedMsg{Users []jira.User}` — carries the user list from the async fetch
- Added `UserSelectedMsg{User jira.User}` — emitted by the modal when the user confirms a selection
- Added `AssignDoneMsg{User jira.User}` — emitted by the Jira call on success; carries the assigned user for the status message

### AssignModal (`internal/tui/modals/assign_modal.go`)
- New modal with a live filter text field (no external text-input component needed)
- Printable key presses append to the filter string; `backspace` removes the last character
- The filtered list is recomputed on every filter change using case-insensitive substring matching on `DisplayName` and `EmailAddress`
- Navigation: `↑`/`k`, `↓`/`j`; confirmation: `enter`; cancel: `esc`

### App wiring (`internal/tui/app.go`)
- Added `modalAssign` state constant
- Added `assignModal modals.AssignModal` field on `Model`
- `KeyAssign` handler: sets `m.loading = true`, dispatches `fetchAssignableUsersCmd()`
- `UsersLoadedMsg` handler: constructs `AssignModal`, switches to `modalAssign`
- `UserSelectedMsg` handler: closes modal, sets loading, dispatches `doAssignCmd(user)`
- `AssignDoneMsg` handler: updates `Assignee` field in-place on `m.allIssues` and `m.currentIssue`; sets status message "Assigned to <name>"
- `fetchAssignableUsersCmd()` and `doAssignCmd(user)` added as `Model` methods

---

## Special cases

### User list scoped to the issue
`/rest/api/3/user/assignable/search?issueKey=<key>` is used instead of the global user search endpoint. This ensures the list only contains users with the permission to be assigned on that specific issue, which avoids presenting users who would be rejected by the API on assignment.

### In-place update instead of list removal
After a successful assign the `Assignee` field is updated in-place on `m.allIssues`. The issue stays visible in the list. This differs from unassign, which removes the issue, because:
- Assigning to someone else is a deliberate handoff — the issue is still worth seeing in context
- Assigning to oneself (if for some reason the user re-assigns to themselves) would cause a jarring disappear/reappear cycle if the list were rebuilt from scratch

### No chord for assign
`a` triggers the assign flow directly (fetch users, then open modal) rather than entering a pending-key chord. A chord would add no value since there is only one assign action.

---

## Open items

- `maxResults=50` on the user search is a hard limit. Projects with many contributors may not show all candidates. A follow-up could add pagination or a server-side search (query the API per keystroke) for very large teams.
