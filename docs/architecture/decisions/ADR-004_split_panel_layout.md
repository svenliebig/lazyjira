# ADR-004 — Split-Panel Layout as Primary Issue View

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

When viewing a list of issues, the user needs to be able to quickly scan the list and read the details of the currently highlighted issue. Two layout approaches were considered:

1. **Sequential navigation** — list view and detail view are separate full-screen views; Enter navigates from list to detail, Esc returns.
2. **Split panel** — list and detail are shown side by side; navigating the list instantly updates the detail panel on the right.

## Decision

Use a split-panel layout for the issue list view: the left panel (~40% width) shows the scrollable issue list, and the right panel (~60% width) shows the details of the currently highlighted issue.

Keyboard focus can be moved to the right panel with Enter (for scrolling the description), and returned to the left panel with Esc. All issue actions (copy, open, transition, AI) are available regardless of which panel has focus.

## Rationale

**Faster information consumption:**
The user can browse the list while simultaneously reading each issue's description, without navigating back and forth between screens. This matches the mental model of tools like lazygit, k9s, and `ranger`.

**Reduced navigation depth:**
The user needs at most one keypress (Enter) to shift focus to the detail panel, versus two keypresses (Enter to open, Esc to return) in the sequential model.

**Actions available from the list:**
Because `currentIssue` is always set to the highlighted issue in the list, all action keys (y, o, t, a) work immediately without needing to "enter" the issue first. This was the primary motivation for the split design.

**Familiar pattern:**
lazygit, which served as the UX inspiration for jira-cli, uses an identical split-panel approach for its commit and branch lists.

## Implementation

The split is implemented inside `IssueListModel` as a composite of two Bubble Tea components:
- **Left:** `bubbles/list.Model` showing issue summaries
- **Right:** `bubbles/viewport.Model` showing the full issue detail

Width ratio: `leftContent = total * 2/5`, `rightContent = total - left - 1` (1 char for the divider border).

`focusRight bool` tracks which panel receives keyboard input. The active panel is indicated by a highlighted border colour (purple = focused, grey = inactive).

When the list cursor moves, the right panel content is rebuilt via `buildIssueDetail()` and the viewport is scrolled to the top.

## Consequences

- `IssueListModel` is more complex than a simple list wrapper — it manages two child components and a focus state
- The initial implementation's `viewIssueDetail` full-screen view is now unreachable (technical debt TD-01)
- The right panel width must be recalculated and rerendered whenever the terminal is resized
- The split ratio (40/60) is hardcoded — not user-configurable at this stage
