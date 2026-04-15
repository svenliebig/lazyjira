# ADR-013 — General Tab System

**Date:** 2026-04-15
**Status:** Accepted

## Context

The application started as a single-view TUI focused on the issue list.
Introducing board configuration (needed for sprint assignment) required a place
in the UI to manage project → board mappings without hiding them in a plain JSON
file. At the same time, future features (sprint views, filter management) will
need first-class screen real estate.

A tab system provides a natural, discoverable container for these views.

## Decision

Implement a general tab bar rendered at the top of the main layout.

- Tabs are cycled with `[` (previous) and `]` (next); wrapping around
- The active tab's title is highlighted with the theme accent colour
- Initial tabs: **Issues** and **Boards**
- Tab state (`activeTab int`) lives on the root `Model` in `app.go`
- Switching tabs cancels any pending chord or open modal

Tab-level views live in `internal/tui/views/` following the existing split-panel
pattern. Each view implements the standard `tea.Model` interface.

## Consequences

**Positive**
- Clear home for board configuration and future board-level views
- `[`/`]` is low-conflict (not used by any existing chord)
- Extensible without changing the root model significantly

**Negative**
- Tab bar consumes one line of vertical space
- Switching tabs resets view-local state (scroll position, filters); acceptable
  at this stage since tabs are independent contexts

## Alternatives considered

- **Settings-only board config** — board IDs buried in `settings.json`; not
  discoverable and does not scale to future board views
- **Modal-based board picker** — triggered at first `c-s` use; one-off UX,
  hard to review or edit later
