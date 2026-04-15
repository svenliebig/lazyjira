# 2026-04-15 — Change action modal (`c`), sprint assignment, tab system & Boards tab

## Goal

Introduce a `c` (change) chord key that groups issue property-mutation actions.
Initial options: `c-s` (sprint) and `c-a` (assign, duplicating `a`).

Sprint selection requires knowing which Jira board a project uses. Rather than
auto-resolving or hiding this in settings, a dedicated **Boards tab** is
introduced as the first step of a general tab system. Board-to-project mappings
are persisted in `~/.config/lazyjira/boards.json` (XDG pattern).

Tabs are cycled with `[` (left) and `]` (right).

---

## Completed work

### Tab system

- New tab bar rendered at the top of the main layout; active tab is highlighted
  using the current theme's accent colour
- Initial tabs: **Issues** (index 0) and **Boards** (index 1)
- `[` / `]` cycle left / right through tabs; wraps around
- Added `KeyTabPrev = "["` and `KeyTabNext = "]"` to `internal/tui/shared/keys.go`
- Root model `app.go` gained `activeTab int` field; tab switching clears any
  pending chord key and hides open modals
- Status bar updated to show tab cycling hint `[:prev  ]:next`

### Boards tab (`internal/tui/views/boards_view.go`)

- Lists all project → board mappings currently saved in `boards.json`
- `a` on the boards tab opens an inline form to add a new mapping:
  project key (text field) → board ID (numeric field)
- `d` / `delete` removes the selected mapping
- `e` opens the mapping for editing in the same inline form
- Navigation: `↑`/`k`, `↓`/`j`; `enter` to confirm form; `esc` to cancel

### Board config persistence (`internal/boards/`)

- New package `internal/boards/` following the same pattern as `internal/settings/`
  and `internal/exclusions/`
- `boards.json` schema: `{ "PROJECT-KEY": boardId, ... }` (string → int map)
- `Load() (map[string]int, error)` — reads `~/.config/lazyjira/boards.json`;
  returns empty map if file does not exist (first run)
- `Save(map[string]int) error` — writes the map atomically

### `c` chord: ChangeModal (`internal/tui/modals/change_modal.go`)

- Numbered/keyed action list:
  ```
  Change issue
    [s] Sprint
    [a] Assign
  ```
- `s` / `a` (or their number) emit `ChangeActionSelectedMsg{Action string}`
- `esc` cancels

### Sprint action (`c-s`)

- Requires a board configured for the issue's project in `boards.json`
- If no board is found: modal closes, status message shown:
  `"No board configured for <PROJECT>. Add one in the Boards tab (])."`
- `SprintPickerModal` in `internal/tui/modals/sprint_picker_modal.go`:
  - Fetches active and future sprints via
    `GET /rest/agile/1.0/board/{boardId}/sprint?state=active,future`
  - Live filter on sprint name; `↑`/`k`, `↓`/`j`; `enter` confirms; `esc` cancels
- New Jira client methods in `internal/jira/sprints.go`:
  - `GetSprints(ctx, boardID int) ([]Sprint, error)`
  - `MoveIssueToSprint(ctx, issueKey string, sprintID int) error` —
    `POST /rest/agile/1.0/sprint/{sprintId}/issue` with `{"issues":["<key>"]}`
- New shared types: `Sprint{ID int, Name string, State string}`
- New messages: `SprintsLoadedMsg`, `SprintSelectedMsg`, `SprintMoveDoneMsg`
- On success: updates `Sprint` field in-place on `m.allIssues` and
  `m.currentIssue`; status message `"Moved to sprint <name>"`

### Assign via change modal (`c-a`)

- Dispatches the same `fetchAssignableUsersCmd()` that `a` already uses
- No duplication of logic; existing `AssignModal` reused unchanged

### App wiring (`internal/tui/app.go`)

- `activeTab` field; `KeyTabPrev`/`KeyTabNext` handlers
- `modalChange` state constant; `changeModal modals.ChangeModal` field
- `KeyChange` (`c`) handler: opens `ChangeModal`
- `ChangeActionSelectedMsg` handler: routes to sprint or assign sub-flow
- `boards.BoardConfig` loaded at startup alongside settings and exclusions;
  passed to the Boards view and consulted by the sprint flow

---

## Special cases

### Agile API base path

Sprint endpoints live under `/rest/agile/1.0/`, not `/rest/api/3/`. The client
uses a separate base URL constant for the Agile path.

### `a` key unchanged

`a` continues to open the assign flow directly. `c-a` is an alias for
discoverability only.

### Tab switching clears chord state

Switching tabs while a chord is pending (e.g. mid-`c`) cancels the chord
silently, same as pressing `esc`.

### Boards tab key isolation

The global key switch in `handleKey()` processes issue-action keys (`a`, `e`,
`d`, `t`, etc.) before delegating to the active view. Without a guard, pressing
`a` on the Boards tab would open the assign modal instead of the board-add form.

Fix: immediately after the modal-delegation check, if `activeTab == tabBoards`
the handler returns `updateActiveChild(msg)` directly, bypassing the global
switch entirely. The only keys processed before this point are `ctrl+c` (quit)
and `[`/`]` (tab cycling), which should work from any context.

---

## Open items

- Kanban-only (Business) projects have no sprints. A follow-up should detect
  this and surface a clear message rather than returning an empty sprint list.
- Sprint picker currently fetches on every `c-s` invocation. Caching per board
  ID (invalidated on move) would reduce API calls.
- Future tabs could include: Sprint view (backlog/active sprint on the board),
  Filters, My Work.
