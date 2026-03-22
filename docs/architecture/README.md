# Architecture Documentation — jira-cli

This documentation follows the [arc42](https://arc42.org) template for software architecture.

## Table of Contents

| # | Section |
|---|---------|
| 1 | [Introduction and Goals](./01_introduction_and_goals.md) |
| 2 | [Architecture Constraints](./02_constraints.md) |
| 3 | [System Scope and Context](./03_context_and_scope.md) |
| 4 | [Solution Strategy](./04_solution_strategy.md) |
| 5 | [Building Block View](./05_building_block_view.md) |
| 6 | [Runtime View](./06_runtime_view.md) |
| 7 | [Deployment View](./07_deployment_view.md) |
| 8 | [Cross-cutting Concepts](./08_crosscutting_concepts.md) |
| 9 | [Architecture Decisions](./09_architecture_decisions.md) |
| 10 | [Quality Requirements](./10_quality_requirements.md) |
| 11 | [Risks and Technical Debt](./11_risks_and_technical_debt.md) |
| 12 | [Glossary](./12_glossary.md) |

## Architecture Decision Records

| ID | Title | Status |
|----|-------|--------|
| [ADR-001](./decisions/ADR-001_bubbletea_tui_framework.md) | Use Bubble Tea as TUI framework | Accepted |
| [ADR-002](./decisions/ADR-002_no_jira_sdk.md) | Direct HTTP instead of Jira SDK | Accepted |
| [ADR-003](./decisions/ADR-003_basic_auth.md) | Jira Cloud Basic Auth with email and API token | Accepted |
| [ADR-004](./decisions/ADR-004_split_panel_layout.md) | Split-panel layout as primary issue view | Accepted |
| [ADR-005](./decisions/ADR-005_local_ollama.md) | Local Ollama for AI assistance | Accepted |
| [ADR-006](./decisions/ADR-006_config_resolution_chain.md) | Three-level config resolution chain | Accepted |
| [ADR-007](./decisions/ADR-007_jql_post.md) | POST /rest/api/3/search/jql for issue search | Accepted |
| [ADR-008](./decisions/ADR-008_chord_key_system.md) | Two-key chord system for sub-actions | Accepted |
| [ADR-009](./decisions/ADR-009_adf_to_text.md) | Plain-text ADF conversion | Accepted |
| [ADR-010](./decisions/ADR-010_shared_messages_package.md) | Shared messages package to prevent circular imports | Accepted |
