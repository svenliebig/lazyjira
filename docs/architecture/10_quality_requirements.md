# 10. Quality Requirements

## Quality Tree

```
Quality
├── Responsiveness
│   ├── UI never blocks during network I/O
│   └── Loading state shown immediately on action
│
├── Usability
│   ├── Discoverability — status bar always shows available actions
│   ├── Consistency — vim-style navigation everywhere
│   └── Efficiency — actions reachable in ≤2 keystrokes from any view
│
├── Reliability
│   ├── Errors surface gracefully (no crashes)
│   └── Credentials stored securely (0600 permissions)
│
├── Portability
│   ├── Single binary, no installer
│   └── macOS, Linux, Windows supported
│
└── Privacy
    ├── No telemetry
    ├── AI runs locally (Ollama)
    └── Credentials never logged
```

## Quality Scenarios

### QS-01 — Responsiveness: Network Latency

**Stimulus:** User triggers "list assigned issues" (`l` → `a`) while on a slow VPN connection (500ms RTT to Jira Cloud).

**Response:** The TUI immediately shows "Loading..." text and remains interactive (the user can press `?` or `q`). The issue list appears when the API response arrives. The UI does not freeze.

**Measure:** The render loop continues to respond to input within one frame (≈16ms) regardless of API latency.

**Satisfied by:** `fetchAssignedCmd` is a `tea.Cmd` running in a goroutine. The Bubble Tea event loop is never blocked.

---

### QS-02 — Usability: Discoverability

**Stimulus:** A new user opens the tool for the first time after authenticating.

**Response:** The status bar shows `l:list  ?:help  q:quit`. Pressing `?` reveals a full shortcut reference. Pressing `l` shows `a - Assigned Issues`. No action requires knowledge outside what is displayed on screen.

**Measure:** A user with no prior knowledge of the tool can discover and execute all primary workflows within 2 minutes.

**Satisfied by:** Dynamic status bar rendering in `renderStatusBar()`, which adapts to every application state. `HelpModal` lists all shortcuts.

---

### QS-03 — Reliability: API Error Handling

**Stimulus:** The Jira API returns an unexpected error (e.g., 401 Unauthorized, 500 Internal Server Error, network timeout).

**Response:** The error message is displayed in the content area. The application remains running. The user can retry by pressing `l` → `a` again or navigate to other screens.

**Measure:** No panic, no crash, no silent failure. Error text includes the HTTP status code and response body.

**Satisfied by:** `ErrMsg` pattern — all `tea.Cmd` closures catch errors and return `ErrMsg` instead of panicking. Root model renders `m.err` in the content area.

---

### QS-04 — Reliability: Nil Pointer Safety on Startup

**Stimulus:** The application starts, receives a `WindowSizeMsg`, and `IssueListModel` has not yet been populated with data.

**Response:** No panic. The list and detail views are not resized until they have been properly initialized.

**Measure:** Zero panics in the startup sequence.

**Satisfied by:** `initialized bool` guard in `IssueListModel.SetSize()` and `IssueDetailModel.SetSize()`. Smoke tests in `internal/tui/smoke_test.go` verify this.

---

### QS-05 — Privacy: Credential Storage

**Stimulus:** A user saves their Jira API token via the auth modal.

**Response:** The config file is written to `~/.config/lazyjira/config.json` with mode `0600`. The API token is not logged, printed, or transmitted anywhere other than the Authorization header of Jira API requests.

**Measure:** `ls -la ~/.config/lazyjira/config.json` shows `-rw-------`. The token field uses `EchoPassword` mode in the TUI input.

**Satisfied by:** `config.Save()` uses `os.WriteFile(path, data, 0600)`. `textinput.EchoPassword` in `AuthModal`.

---

### QS-06 — Portability: Single Binary

**Stimulus:** A developer wants to install lazyjira on a new machine.

**Response:** A single binary is copied to a directory in `$PATH`. No additional packages, runtimes, or dependencies are required (beyond OS-level clipboard tools on Linux).

**Measure:** `go build -o lazyjira .` produces a self-contained binary on all target platforms.

**Satisfied by:** Pure Go with no CGo. All dependencies are compiled into the binary by `go build`.
