# ADR-009 — Plain-Text ADF Conversion

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

Jira Cloud REST API v3 returns issue descriptions in [Atlassian Document Format (ADF)](https://developer.atlassian.com/cloud/jira/platform/apis/document/structure/) — a nested JSON tree structure where each node has a `type` and optional `content` (children) and `text` (leaf) fields.

Example ADF:
```json
{
  "type": "doc",
  "content": [{
    "type": "paragraph",
    "content": [{
      "type": "text",
      "text": "This is a description."
    }]
  }]
}
```

The description must be rendered in the terminal. Options:
1. **Render as styled terminal text** (bold, italic, bullet points using ANSI codes)
2. **Convert to plain text** with structural whitespace

## Decision

Convert ADF to plain text using `adfToText()`, a recursive function that:
- Extracts all `"text"` leaf values
- Appends newlines after block nodes (`paragraph`, `heading`, `bulletList`, `orderedList`, `listItem`, `blockquote`, `codeBlock`, `rule`)
- Falls back to empty string for unknown node types

The converted plain text is stored in `Issue.Fields.Description` and displayed in a `bubbles/viewport`.

## Rationale

**Terminal rendering limitations:**
Rich text (bold, italic, colours) in a terminal requires ANSI escape codes. Rendering ADF faithfully would require mapping each ADF mark type to the correct escape sequence, handling nested marks, and ensuring correct reset behaviour — significant implementation effort.

**Viewport handles wrapping and scrolling:**
The `bubbles/viewport` component already handles word wrapping and scroll position. Plain text integrates naturally without any additional processing.

**Descriptions are primarily prose:**
In practice, Jira issue descriptions are mostly prose with occasional bullet lists and code blocks. The structural information (headings, lists) is preserved via newlines, which is sufficient for reading in a terminal.

**Clipboard copy uses plain text:**
When the user copies the description (`y → d`), they expect plain text on the clipboard, not a JSON blob or ANSI escape codes. Plain-text conversion serves both display and clipboard use.

## Limitations

- Inline formatting (bold, italic, code spans) is stripped; only the text content is preserved
- Code blocks are rendered as plain text without syntax highlighting
- Table support is not implemented; table content is partially rendered (text extracted, structure lost)
- Links are rendered as link text only; the URL is dropped

These are accepted trade-offs for a read-only terminal display. The user can always open the issue in the browser (`o`) for the full formatted view.

## Consequences

- `adfToText()` in `internal/jira/issues.go` (~50 lines)
- Description stored as `string` in `Issue.Fields.Description` — no ADF model in the domain layer
- Rich formatting information (bold, italic, colours) is lost at the point of parsing
- The function handles both ADF objects and plain strings (for compatibility with older API responses)
