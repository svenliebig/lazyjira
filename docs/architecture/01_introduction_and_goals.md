# 1. Introduction and Goals

## What is lazyjira?

lazyjira is a terminal user interface (TUI) tool for minimal, focused interaction with Jira Cloud. It aims to give developers fast access to their Jira workflow without leaving the terminal or navigating a web browser.

The tool is inspired by [lazygit](https://github.com/jesseduffield/lazygit) in look, feel, and interaction model: a full-screen TUI with keyboard-driven navigation, context-aware action menus, and a split-panel layout.

## Requirements Overview

| ID | Requirement |
|----|-------------|
| R-01 | Authenticate against Jira Cloud using a URL and API token |
| R-02 | Credentials can be provided via CLI flags, environment variables, or a config file |
| R-03 | List issues assigned to the current user |
| R-04 | Show issue details (summary, description, status, assignee, reporter) |
| R-05 | Transition an issue to a new status |
| R-06 | Copy issue key, URL, title, or description to the clipboard |
| R-07 | Open an issue in the browser |
| R-08 | Generate an AI-assisted work summary from git commits linked to an issue |
| R-09 | All interaction through keyboard shortcuts — no mouse required |
| R-10 | Navigation via arrow keys and vim-style keys (h/j/k/l) |

## Quality Goals

| Priority | Quality Goal | Motivation |
|----------|-------------|------------|
| 1 | **Responsiveness** | The TUI must never block on I/O. All network calls and subprocess executions are asynchronous. |
| 2 | **Simplicity** | Minimal surface area. Only the operations a developer needs daily. No feature bloat. |
| 3 | **Discoverability** | Every available action is visible in the status bar. Users should not need to memorize shortcuts. |
| 4 | **Portability** | Runs on macOS, Linux, and Windows without modification. |
| 5 | **Privacy** | AI features use a local LLM (Ollama). No issue data leaves the user's machine for AI processing. |

## Stakeholders

| Role | Interest |
|------|---------|
| Developer (primary user) | Fast Jira interaction during development without switching contexts |
| Team Lead | Issue status transitions from the terminal during standups or code review |
| Security-conscious team | Credentials stored locally, AI runs locally — no third-party cloud services beyond Jira itself |
