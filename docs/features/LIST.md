# List

The user can display different lists of issues by pressing `l` to open the list selector modal, then choosing a list type.

## List Types

### Assigned Issues (`l → a`)

Shows all issues currently assigned to the authenticated user where the status category is not Done, ordered by last updated. The list is filtered by any active [exclusion rules](#excluded-issues-l--x) before being displayed.

### Excluded Issues (`l → x`)

Shows all active exclusion rules. Each row displays the excluded value (issue key or parent key) and the rule type. Pressing `x` on a highlighted row removes that exclusion immediately.

## Split-panel Layout

When a list is displayed, the screen is split into two panels:

- **Left panel (~40%)** — scrollable list of issues or rules. Navigation with `j`/`k` or arrow keys.
- **Right panel (~60%)** — detail view of the currently highlighted issue (key, summary, status, assignee, reporter, description). Updates live as the cursor moves.

Pressing `Enter` moves focus to the right panel for scrolling long descriptions. Pressing `Esc` returns focus to the left panel.

## Filtering

Issues fetched from the Jira API are filtered against the local exclusion rules before the list is built. The raw API results are kept in memory, so adding or removing an exclusion updates the displayed list instantly without a new network request.
