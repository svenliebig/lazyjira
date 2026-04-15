# Boards

The Boards tab manages the mapping between Jira projects and their associated
Agile boards. This mapping is required by any feature that needs board-level
data (sprint picker, future sprint views).

## Navigating to the Boards tab

Press `]` from anywhere in the application to move to the next tab, or `[` to
move to the previous tab. The active tab is highlighted in the tab bar at the
top of the screen.

## Managing board mappings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `a` | Add a new project → board mapping |
| `e` | Edit the selected mapping |
| `d` | Delete the selected mapping |

When adding or editing a mapping, an inline form asks for:
- **Project key** — the Jira project key (e.g. `PROJ`)
- **Board ID** — the numeric Jira board ID (visible in the board URL on Jira:
  `.../jira/software/projects/PROJ/boards/123` → board ID is `123`)

Mappings are saved to `~/.config/lazyjira/boards.json`.

## Why a separate tab?

Board IDs are not exposed in Jira's standard REST API (they belong to the Agile
API). Rather than silently picking a board or hiding this in a settings file,
the Boards tab makes the configuration explicit and editable without leaving
the TUI.

This tab is also the foundation for future board-level views (sprint backlogs,
active sprint view, etc.).
