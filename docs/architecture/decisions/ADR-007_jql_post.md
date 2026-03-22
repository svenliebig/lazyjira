# ADR-007 — POST /rest/api/3/search/jql for Issue Search

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

The original implementation used `GET /rest/api/3/search?jql=...` to list assigned issues. This endpoint was removed by Atlassian (HTTP 410 Gone) and replaced with `/rest/api/3/search/jql`, documented in [CHANGE-2046](https://developer.atlassian.com/changelog/#CHANGE-2046).

The new endpoint supports both `GET /rest/api/3/search/jql?jql=...` and `POST /rest/api/3/search/jql` with a JSON body.

Additionally, the initial implementation showed 50 "unassigned" results — the JQL filter was not working and the `assignee` field was missing from responses.

## Decision

Use `POST /rest/api/3/search/jql` with a JSON request body:

```json
{
  "jql": "assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC",
  "maxResults": 50,
  "fields": ["summary", "description", "status", "assignee", "reporter"]
}
```

## Rationale

**POST over GET for complex queries:**
The `currentUser()` JQL function contains parentheses (`()`). While these are technically valid in URL query strings per RFC 3986, some proxies and API gateways may reject or misinterpret them. A POST body avoids URL-encoding ambiguities entirely.

**Explicit field selection is required:**
The new `/search/jql` endpoint does not guarantee returning all fields by default. Without explicitly requesting `["summary", "description", "status", "assignee", "reporter"]`, the `assignee` field was absent from responses — causing all 50 results to show as "unassigned". Explicit field selection is both correct and more efficient (smaller response payloads).

**Filter to active issues:**
Adding `AND statusCategory != Done` excludes completed issues. This gives the user a relevant view of their active workload rather than a historical archive of all assignments.

**`currentUser()` with Basic Auth:**
The `currentUser()` JQL function resolves to the authenticated user based on the email in the Basic Auth header. Combined with the decision in ADR-003, this correctly returns only issues assigned to the account whose email is in the config.

## Root Cause of the Original Bug

The initial failure (50 unassigned results) had two causes:
1. The `GET` form with URL-encoded JQL may not have resolved `currentUser()` correctly, returning all issues up to the max
2. The `assignee` field was not returned without explicit field selection in the new endpoint

## Consequences

- Issue search uses POST instead of GET — unusual for a read operation, but correct per the Atlassian API specification
- Field selection must be kept in sync if new fields are needed (e.g., labels, priority)
- The `statusCategory != Done` filter hides resolved issues — this is a deliberate UX choice, not a technical limitation
- The `maxResults: 50` limit remains hardcoded (see TD-05 for pagination debt)
