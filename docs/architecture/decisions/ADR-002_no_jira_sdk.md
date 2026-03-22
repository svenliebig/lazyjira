# ADR-002 — Direct HTTP Instead of Jira SDK

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

The application needs to communicate with the Jira Cloud REST API. There are Go libraries that wrap the Jira API (e.g., `go-jira`, `andygrunwald/go-jira`). An alternative is to call the API directly using Go's `net/http`.

## Decision

Call the Jira Cloud REST API v3 directly using Go's standard `net/http` package. No third-party Jira SDK is used.

## Rationale

**Minimal API surface:**
jira-cli uses exactly four API endpoints:
1. `POST /rest/api/3/search/jql` — list assigned issues
2. `GET /rest/api/3/issue/{key}` — get single issue
3. `GET /rest/api/3/issue/{key}/transitions` — get available transitions
4. `POST /rest/api/3/issue/{key}/transitions` — apply transition

A full SDK would map hundreds of endpoints and data types. For four simple endpoints, the SDK's overhead in terms of dependency weight, API learning curve, and generated code noise exceeds the benefit.

**Control over the request format:**
The `POST /rest/api/3/search/jql` endpoint requires a specific JSON body with explicit field selection (`"fields": [...]`). SDK abstractions often hide this level of control or use their own field mapping logic.

**Smaller dependency graph:**
`net/http` is part of the Go standard library. Adding a Jira SDK adds transitive dependencies, increases binary size, and introduces an external version-management burden.

**ADF parsing:**
Jira issue descriptions use Atlassian Document Format (ADF), a nested JSON structure. No existing Go Jira SDK handles ADF-to-text conversion for terminal display correctly. A custom `adfToText()` function was necessary regardless.

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **andygrunwald/go-jira** | Wraps deprecated `/rest/api/2`, unclear v3 support, heavy transitive dependencies |
| **ctreminiom/go-atlassian** | Full Atlassian Cloud SDK; significant API surface for 4 endpoints; adds many transitive deps |

## Consequences

- Four API call implementations in `internal/jira/issues.go` (~250 lines including response structs)
- ADF-to-text conversion implemented locally (`adfToText()` recursive function)
- Auth header (`SetBasicAuth`) and JSON marshaling managed manually per request
- Any new Jira feature requires a new hand-written API call
- Full control over request format and error handling
