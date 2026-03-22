# ADR-006 — Three-Level Config Resolution Chain

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

The application needs Jira credentials (URL, email, API token) to function. These need to be configurable in multiple ways to support different usage scenarios: interactive first-run setup, scripting via environment variables, and short-lived overrides via CLI flags.

## Decision

Resolve credentials using a three-level chain where each level overrides the previous:

```
1. Config file:    ~/.config/jira-cli/config.json         (lowest priority)
2. Env variables:  JIRA_CLOUD_URL / JIRA_EMAIL / JIRA_API_TOKEN
3. CLI flags:      --jira-cloud-url / --jira-email / --jira-api-token  (highest priority)
```

If none of the three levels provides complete credentials, the TUI shows the interactive authentication modal on startup.

## Rationale

**Config file for persistence:**
Most users set up credentials once and reuse them indefinitely. A persistent config file avoids re-entering credentials on every run.

**Environment variables for automation:**
CI/CD pipelines, scripts, and container environments cannot use interactive prompts. Environment variables are the standard way to inject secrets in these contexts.

**CLI flags for testing and overrides:**
Short-lived overrides (e.g., pointing at a test Jira instance) are most conveniently done via flags without modifying the config file or environment.

**Each level overrides the previous:**
This is the conventional precedence order in most Unix CLI tools (e.g., git, aws-cli, kubectl). Users familiar with one tool will intuitively understand the order.

**XDG compliance:**
The config file path follows the XDG Base Directory Specification (`$XDG_CONFIG_HOME/jira-cli/config.json`, defaulting to `~/.config/jira-cli/config.json`). This respects user and system conventions on Linux and macOS.

## Implementation

```go
// config.Load(flags Flags) *Config:
// 1. Read file
cfg := readFromFile()
// 2. Override with env vars
if v := os.Getenv("JIRA_CLOUD_URL"); v != "" { cfg.JiraCloudURL = v }
if v := os.Getenv("JIRA_EMAIL"); v != "" { cfg.JiraEmail = v }
if v := os.Getenv("JIRA_API_TOKEN"); v != "" { cfg.JiraAPIToken = v }
// 3. Override with flags
if flags.JiraCloudURL != "" { cfg.JiraCloudURL = flags.JiraCloudURL }
if flags.JiraEmail != "" { cfg.JiraEmail = flags.JiraEmail }
if flags.JiraAPIToken != "" { cfg.JiraAPIToken = flags.JiraAPIToken }
```

## Consequences

- Users have three ways to authenticate; this needs to be documented clearly
- The `IsComplete()` check requires all three fields to be present — partial config triggers the auth modal even if some values are set
- CLI flags expose the API token in the process list (`ps aux`) — this is a known trade-off for the flag approach; env vars are preferred for secrets in production use
- Config is saved with `0600` permissions; written only when the user submits the auth modal
