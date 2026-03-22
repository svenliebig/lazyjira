# ADR-008 — Two-Key Chord System for Sub-actions

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

Several features in lazyjira require a choice between multiple related sub-actions:
- **Copy**: copy the issue key, URL, title, or description
- **AI assistance**: generate a work summary (and potentially other AI actions in the future)

These sub-actions could be exposed as:
1. Nested menus opened by a single key, navigated with arrows
2. Two-key chord sequences (e.g., `y` then `k` to copy the key)
3. Single unique keys per action (e.g., `K` for key, `U` for URL, etc.)

## Decision

Use a two-key chord system for grouping related sub-actions:
- `y` opens the copy menu, then `k`/`u`/`t`/`d` select the specific copy action
- `a` opens the AI menu, then `s` triggers the AI summary

The first key (`y`, `a`) opens a modal overlay that displays the available second keys and their descriptions. `Esc` cancels the chord.

## Rationale

**Mnemonics reduce cognitive load:**
`y` (yank) is the vim convention for copy operations. `k` for key, `u` for URL, `t` for title, `d` for description are intuitive sub-keys. The two-level structure groups related actions without requiring the user to remember many individual bindings.

**Status bar feedback:**
When a first key is pressed, the status bar immediately updates to show the available second keys. The user never has to guess what to press next.

**No mode pollution:**
Chords avoid adding many single-letter bindings at the top level that could conflict with navigation keys. The top-level key space remains clean.

**Lazy menu disclosure:**
The modal only appears when the first key is pressed. Users who never use the copy feature are never exposed to its sub-keys.

## Implementation

`pendingKey string` in the root model tracks the first key of an active chord. The keyboard dispatch in `handleKey()` checks this before the global key switch:

```go
if m.pendingKey == shared.KeyCopy {   // "y"
    m.pendingKey = ""
    switch key {
    case "k": // copy key
    case "u": // copy URL
    case "t": // copy title
    case "d": // copy description
    }
}
```

The copy modal (`CopyModal`) and AI modal (`AIModal`) provide the visual representation — they display the available sub-keys and their labels.

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **Unique single keys** | Would pollute the top-level keyspace; harder to discover |
| **Arrow-navigated menu** | More keypresses to reach an action; less efficient for common operations |
| **Prefix-less modal** | A single key `y` opens a modal you navigate with arrows — less efficient than chord |

## Consequences

- `pendingKey` state in the root model must be cleared on `Esc` and after any chord resolution
- The status bar must always reflect the current `pendingKey` state so the user knows a chord is "open"
- Adding new chord families requires only adding a new `pendingKey` case in `handleKey()` and a matching modal
- First keys (`y`, `a`) cannot be used as standalone actions — they are exclusively chord initiators
