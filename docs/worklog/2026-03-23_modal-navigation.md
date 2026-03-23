# 2026-03-23 — Modal cursor navigation

## Goal

Add arrow-key and vim-style (`h`, `j`, `k`, `l`) navigation to all action-selection modals so users can navigate options without memorising every shortcut. Keys already reserved for navigation (`h`, `j`, `k`, `l`) must not be reused as action shortcuts inside modals.

---

## Completed work

### Cursor navigation in all four action modals

All modals (`CopyModal`, `TransitionModal`, `ListSelectorModal`, `ExcludeModal`) gained a `cursor int` field and a consistent navigation contract:

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` / `l` | Confirm highlighted selection |
| `h` / `esc` | Cancel / close modal |

The currently highlighted row is rendered with a `>` prefix and the `StyleSelectedItem` style (bold purple), matching the existing selected-item style used in list views.

### Key conflict resolution

`k` was previously used as a direct shortcut in two modals, conflicting with vim-up navigation:

- **`CopyModal`**: `k` (copy issue key) replaced by `i`
- **`ExcludeModal`**: `k` (exclude by issue key) replaced by `i`

### Non-conflicting shortcuts preserved

Shortcuts that do not overlap with navigation keys are kept as direct single-key aliases alongside cursor navigation:

- `CopyModal`: `i` (key), `u` (URL), `t` (title), `d` (description)
- `TransitionModal`: `1`–`9` (direct transition select)
- `ListSelectorModal`: `a` (assigned), `x` (excluded)
- `ExcludeModal`: `i` (key), `p` (parent, when available)

### Help modal updated

- Removed the now-stale `y → k` entry; added `y → i` for copy issue key.
- Added `enter/l: Select item` and `h: Cancel / close modal (in action modals)` entries.

### Status bar updated

The copy-chord hint in the status bar was updated to replace `k:key` with `i:key` and to include brief navigation hints (`↑/k ↓/j: navigate`, `enter/l: select`, `h/esc: cancel`).

---

## Special cases

### `h` as cancel in modals

`h` (vim-left) is interpreted as "go back / cancel" inside action modals. This mirrors the spatial metaphor: navigating left exits the current context. Outside modals, `h` continues to work as the vim-left navigation key for the main list views, so there is no behavioural change at the global level.

### `l` as confirm in modals

`l` (vim-right) means "navigate right into the selection" — confirming the highlighted item. This is consistent with how vim-style navigation typically works in nested contexts (left = out, right = in).

### `ExcludeModal` cursor clamps to available items

When the current issue has no parent, `NewExcludeModal` initialises `cursor = 1` (pointing to the key row) and the up-navigation guard prevents the cursor from landing on position 0 (the disabled parent row). The parent row is still rendered with strikethrough and muted colour for discoverability, but it is never focusable.

### Dead-code chord resolution in `app.go`

The chord-resolution block in `handleKey()` (lines ~275–294) that handles `y → k/u/t/d` is effectively dead code: by the time a second key arrives, `activeModal == modalCopy` causes an early return to `updateActiveChild`, so the chord block is never reached. It was left in place to avoid unrelated refactoring during this session; its only user-visible artefact (the `k:key` status bar hint) was removed.

---

## Open items

- The dead chord-resolution block in `app.go` could be removed in a future cleanup session.
