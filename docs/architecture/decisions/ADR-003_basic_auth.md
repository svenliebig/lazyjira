# ADR-003 — Jira Cloud Basic Auth with Email and API Token

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

Jira Cloud supports multiple authentication methods for REST API access:
- **OAuth 2.0 (3LO)** — requires app registration and a browser-based authorization flow
- **Personal Access Tokens** — supported on Jira Data Center only, not Jira Cloud
- **Basic Auth with API Token** — email address + API token from `id.atlassian.net`, encoded as `Authorization: Basic base64(email:token)`

The initial implementation incorrectly used `Authorization: Bearer <token>`, which returned a 403 Forbidden response from Jira Cloud.

## Decision

Use HTTP Basic Authentication with the user's email address and Jira Cloud API token:

```
Authorization: Basic base64("<email>:<api-token>")
```

Implemented via Go's `req.SetBasicAuth(email, apiToken)`.

Both `email` and `apiToken` are stored in the config file and collected during the auth modal flow.

## Rationale

**No app registration required:**
OAuth 2.0 (3LO) requires registering an Atlassian Connect app or an OAuth 2.0 app in the Atlassian Developer Console. This is impractical for a personal CLI tool.

**No browser dependency:**
The OAuth device flow would introduce browser-based authorization, breaking the terminal-only UX goal.

**Widely used convention:**
Basic auth with API tokens is the standard approach for Jira Cloud REST API access in CLI tools, CI/CD pipelines, and scripts. Users are accustomed to generating API tokens at `id.atlassian.net`.

**JQL `currentUser()` works with Basic Auth:**
The JQL function `currentUser()` resolves to the authenticated user based on the email in the Basic Auth header, making the "assigned to me" query work correctly without additional parameters.

## Why Bearer Failed

`Authorization: Bearer <token>` is the correct format for **OAuth 2.0 access tokens**. Jira Cloud API tokens are not OAuth access tokens — they are personal API tokens intended for Basic Auth. Sending them as Bearer tokens results in 403 Forbidden.

## Alternatives Considered

| Alternative | Reason not chosen |
|-------------|------------------|
| **OAuth 2.0 (3LO)** | Requires app registration, browser callback; inappropriate for a personal CLI |
| **Bearer token (API token)** | Incorrect for Jira Cloud API tokens; returns 403 |

## Consequences

- Users must store their email address in addition to the API token
- The auth modal requires a third input field (added after the initial implementation)
- The config file schema has three fields: `jiraCloudUrl`, `jiraEmail`, `jiraApiToken`
- API tokens from Atlassian should be treated as passwords and stored only in the `0600` config file
- Users should revoke and regenerate tokens if accidentally exposed (e.g., shared in logs or messages)
