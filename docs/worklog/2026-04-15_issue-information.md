# 2026-04-15 — Issue overview: estimated time, time remaining & sprint

## Goal

Extend the issue detail panel (right side of the split view) to surface three additional fields from Jira: original estimated time, time remaining, and the sprint the issue belongs to.

---

## Completed work

### Jira types (`internal/jira/types.go`)

- Added `TimeTracking` struct with fields `OriginalEstimateSeconds int` and `RemainingEstimateSeconds int`.
- Added `Sprint` struct with fields `ID int`, `Name string`, `State string`.
- Extended `IssueFields` with:
  - `TimeTracking TimeTracking`
  - `Sprint *Sprint` (pointer; nil when no sprint is set)

### Jira client (`internal/jira/issues.go`)

- Added `timetracking` and `sprint` to the `fields` list requested in the JQL search payload so both fields are returned with every issue fetch — no extra request needed.
- Both fields are deserialized from the JSON response into the extended `IssueFields`.

### Issue detail view (`internal/tui/views/issue_detail.go`)

- `buildDetailContent` extended to render the new fields between the assignee/reporter block and the horizontal divider:
  - **Sprint** — shown only when `issue.Fields.Sprint != nil`; renders the sprint name with its state in parentheses (e.g. `Sprint 42 (active)`).
  - **Original estimate** — shown only when `OriginalEstimateSeconds > 0`; value is formatted as a human-readable duration (e.g. `2h 30m`).
  - **Time remaining** — shown only when `RemainingEstimateSeconds > 0`; same duration format.
- A helper `formatSeconds(s int) string` converts seconds to `Xh Ym` notation, omitting the hours or minutes component when zero.

---

## Special cases

- The Jira REST API v3 returns `sprint` as a custom field (`customfield_10020`) in the `fields` object. The field is an array of sprint objects; only the last element (the active or most recent sprint) is used.
- `timetracking.originalEstimateSeconds` and `timetracking.remainingEstimateSeconds` are absent from the JSON response (not just zero) when time tracking is disabled for the issue. The JSON deserializer leaves the struct fields at their zero values, so the `> 0` guard doubles as a "field present" check.
- Sprint `state` values from the API are `"active"`, `"closed"`, and `"future"`. They are displayed as-is (lower-case) without mapping to avoid coupling to undocumented Jira internals.

---

## Open items

- None.
