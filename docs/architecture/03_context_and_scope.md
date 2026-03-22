# 3. System Scope and Context

## Business Context

lazyjira sits between the developer's terminal workflow and external systems. It acts as a thin client for Jira Cloud, with optional integrations to the local git repository and a local AI model.

```
┌─────────────────────────────────────────────────────────────────────┐
│                          Developer's Machine                        │
│                                                                     │
│   ┌──────────┐    keyboard/     ┌─────────────┐                    │
│   │Developer │◄────display─────►│  lazyjira   │                    │
│   └──────────┘                  └──────┬──────┘                    │
│                                        │                            │
│             ┌──────────────────────────┼────────────────┐          │
│             │                          │                │          │
│             ▼                          ▼                ▼          │
│   ┌──────────────────┐    ┌────────────────────┐  ┌──────────┐    │
│   │  Git Repository  │    │  Ollama (local AI) │  │Clipboard │    │
│   │  (cwd)           │    │  localhost:11434    │  │& Browser │    │
│   └──────────────────┘    └────────────────────┘  └──────────┘    │
└───────────────────────────────────┬─────────────────────────────────┘
                                    │ HTTPS
                                    ▼
                         ┌─────────────────────┐
                         │   Jira Cloud        │
                         │   REST API v3       │
                         │   (Atlassian)       │
                         └─────────────────────┘
```

## External Systems

| System | Interface | Direction | Purpose |
|--------|-----------|-----------|---------|
| **Jira Cloud** | HTTPS REST API v3 | Outbound | Read issues, fetch transitions, apply transitions |
| **Git** (local) | `git log` subprocess | Outbound | Find commits linked to an issue by key |
| **Ollama** | HTTP REST (`localhost:11434`) | Outbound | Generate AI work summaries |
| **System Clipboard** | OS API via `atotto/clipboard` | Outbound | Copy issue data |
| **System Browser** | OS `open`/`xdg-open`/`start` | Outbound | Open issue URL |
| **Filesystem** | `os.ReadFile`/`os.WriteFile` | Both | Read and write credentials config |

## Technical Context

```
lazyjira binary
│
├── stdin/stdout/stderr  ──► Terminal emulator (user interaction)
│
├── HTTPS :443           ──► *.atlassian.net (Jira Cloud)
│
├── HTTP :11434          ──► localhost (Ollama)
│
├── exec.Command("git")  ──► Git binary in PATH
│
├── exec.Command("open") ──► macOS / Linux / Windows browser opener
│
└── ~/.config/lazyjira/config.json  ──► Local filesystem
```

## Scope Boundaries

**In scope:**
- Viewing and filtering assigned issues
- Transitioning issue status
- Copying issue metadata to clipboard
- Opening issues in the browser
- AI-assisted work summary generation from linked git commits

**Explicitly out of scope:**
- Creating or editing issues
- Comment creation
- Attachment handling
- Sprint and board management
- Jira Data Center / Server
- Multi-account or multi-workspace support
