# 2026-04-15 — UI bug fixes: duplicate "Issues" title and status badge alignment

## Goal

Fix two visual bugs in the split-panel issue list view.

---

## Completed work

### Bug 1: Duplicate "Issues" headline under the tab (`internal/tui/views/issue_list.go`)

The `bubbles/list` component was initialised with `l.Title = "Issues"`, causing it to render its own "Issues" heading directly below the tab bar, which already shows an "Issues" tab label.

**Fix:** replaced `l.Title = "Issues"` with `l.SetShowTitle(false)` so the list component no longer renders its own title.

### Bug 2: Status badge misaligned with its label (`internal/tui/views/issue_list.go`)

`StyleIssueStatus` uses a `lipgloss.NormalBorder()`, making the rendered badge three lines tall (top border, content, bottom border). Appending it to the `"Status:"` label string via `strings.Builder` placed the badge on the same start line as the label, but the border box extended two lines below, misaligning the two elements.

**Fix:** replaced the raw string concatenation with `lipgloss.JoinHorizontal(lipgloss.Center, label, badge)`. Lipgloss vertically centers the single-line label against the three-line badge, so both appear visually aligned on the same row.

---

## Files changed

| File | Change |
|---|---|
| `internal/tui/views/issue_list.go` | `SetShowTitle(false)` instead of setting title text; `JoinHorizontal` for status row |

---

## Open items

- None.
