# ADR-011 — Client-side Issue Exclusions

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

Users sometimes have issues assigned to them that they do not want to see in their daily working list — blocked issues, issues awaiting external input, or parent epics that appear alongside their children. They need a way to hide these from the list permanently without resolving or re-assigning them in Jira.

Options for where to store and apply exclusion rules:

1. **Client-side** — a local file (`~/.config/jira-cli/exclusions.json`) that filters results after fetching from the API
2. **Jira labels** — add a label (e.g. `jira-cli-excluded`) to the issue in Jira, then exclude it from the JQL query
3. **Jira saved filters** — manage a saved JQL filter in Jira that the user maintains themselves
4. **Jira custom field** — write a custom field value via the API to mark issues as excluded

## Decision

Store exclusion rules client-side in `~/.config/jira-cli/exclusions.json`, managed entirely by the CLI. The Jira API is not involved in storing or applying exclusions.

Two rule types are supported:
- `{"type": "key", "value": "PROJ-123"}` — hides a single issue by its key
- `{"type": "parent", "value": "PROJ-100"}` — hides all issues whose parent key matches

Rules are applied by filtering the API response in memory before building the issue list.

## Rationale

**No Jira write permissions required:**
The CLI only needs read access and transition write access. Writing labels or custom fields would require broader permissions that many users — especially read-only consumers of a managed Jira instance — may not have. Client-side exclusions work with any permission level.

**No Jira instance side-effects:**
Adding labels or custom fields to issues changes data visible to other team members in Jira. A personal "I don't want to see this in my terminal" preference should not alter the canonical issue state that teammates rely on.

**Exclusions are personal, not shared:**
Whether a developer wants to hide a blocked epic from their terminal view is a local UI preference, not a workflow state. Storing it locally keeps the concern where it belongs.

**Simpler rule types — especially parent exclusion:**
"Exclude all issues with this parent" cannot be expressed as a single-issue label without labelling every child. Client-side filtering can match `issue.Fields.Parent.Key` against a set in O(n) over the already-fetched result set — no additional JQL complexity or API calls required.

**No re-fetch needed when rules change:**
The root model stores `allIssues` (the raw API result) separately from the displayed filtered list. Adding or removing an exclusion re-filters the in-memory slice immediately, with no network round-trip. This is fast and works offline.

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **Jira labels** | Requires write permission; mutates shared issue data visible to other users |
| **Jira saved filters** | Requires manual JQL management in the Jira UI; the CLI cannot write saved filters without the Jira Admin scope |
| **Jira custom field** | Requires admin to create the field; requires write permission; persists in Jira, not the user's local preferences |

## Consequences

- Exclusions are machine-local: switching to a new machine or deleting `~/.config/jira-cli/` loses them
- Exclusions are not synced across devices or team members
- If an excluded issue's key is renamed in Jira (rare but possible), the rule becomes stale silently — the issue will reappear in the list
- Parent exclusion is by parent key at fetch time; if an issue is re-parented in Jira, it may reappear or disappear after the next fetch
- `internal/exclusions` is a new persistence package alongside `internal/config`; both write to `~/.config/jira-cli/` with `0600` permissions

## Implementation

```
internal/exclusions/store.go
  Rule{Type, Value}           — serialised as JSON
  Store.Add(Rule)             — appends + saves immediately
  Store.Remove(Rule)          — removes + saves immediately
  Store.Filter([]Issue)       — O(n) filter, builds key/parent lookup maps
  Load() (*Store, error)      — reads file; returns empty Store on file-not-found
```

The `allIssues` / `issueListView` split in `app.go`:
```
IssueListLoadedMsg → m.allIssues = msg.Issues
                   → m.issueListView = NewIssueListModel(exclusions.Filter(allIssues))

ExcludeActionMsg   → exclusions.Add(rule)
                   → m.issueListView = NewIssueListModel(exclusions.Filter(allIssues))
```

`allIssues` is never modified after being set; the filter produces a new slice each time.
