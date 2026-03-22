# 11. Risks and Technical Debt

## Risks

### R-01 — Jira Cloud API Deprecations

**Description:** The Jira Cloud REST API v3 evolves. The `GET /rest/api/3/search` endpoint was already removed (410 Gone) and replaced with `POST /rest/api/3/search/jql`. Further endpoint changes may break the tool without warning.

**Probability:** Medium — Atlassian has a history of versioning changes.

**Impact:** High — the tool becomes non-functional if a used endpoint is removed.

**Mitigation:**
- Subscribe to the [Atlassian Developer Changelog](https://developer.atlassian.com/changelog/)
- Error messages include the full HTTP response body (status + error text), making deprecation notices visible to users
- The small number of API calls (4 endpoints) limits exposure

---

### R-02 — No Issue List Refresh

**Description:** Once the issue list is loaded, it is not refreshed unless the user explicitly re-opens the list (`l` → `a`). If issues are modified externally (by teammates or workflow automations), the displayed list is stale.

**Probability:** High — this will happen in active teams.

**Impact:** Low — the user sees an outdated list but can refresh manually.

**Mitigation:** Add a periodic refresh or an explicit refresh key (`r`) in a future iteration.

---

### R-03 — Local Ollama Dependency

**Description:** The AI summary feature silently fails if Ollama is not running at `localhost:11434`. The error is surfaced as an `ErrMsg`, but the user may not understand why.

**Probability:** High for users who haven't set up Ollama.

**Impact:** Low — the feature degrades gracefully; all other functionality is unaffected.

**Mitigation:** Improve the error message to include setup instructions. Optionally pre-check Ollama availability before entering the AI modal.

---

### R-04 — Git Log Filtering by Issue Key

**Description:** `CommitsForIssue` runs `git log --oneline --all` and filters lines containing the issue key string (e.g., "PROJ-42"). This will produce false positives if the key string appears in unrelated commit messages, and misses commits where the key is referenced differently (e.g., in the body, or with different casing).

**Probability:** Low — commit conventions usually place the issue key consistently.

**Impact:** Low — the AI summary may include irrelevant commits or miss some. The output is informational only.

**Mitigation:** Accept as a limitation at current scope. A future improvement could parse the full commit message including body, or use conventional commit formats.

---

### R-05 — No Authentication Token Validation

**Description:** There is no validation of the Jira URL or API token during the auth modal. The user can save invalid credentials, which will result in 401/403 errors on every API call until credentials are corrected.

**Probability:** Medium — typos in URL or token are common.

**Impact:** Medium — the tool is unusable until credentials are fixed, which requires manual config file editing or re-running the auth flow.

**Mitigation:** Add a lightweight test request (e.g., `GET /rest/api/3/myself`) after the user submits credentials, and show an inline error in the modal if it fails.

---

## Technical Debt

### TD-01 — `viewIssueDetail` is Unreachable

`IssueDetailModel` and `viewIssueDetail` exist in the codebase but are no longer reachable through normal navigation after the split-panel view was introduced. The `IssueSelectedMsg` handler still transitions to `viewIssueDetail`, but no code path emits that message anymore.

**Impact:** Dead code increases maintenance surface.

**Recommendation:** Either remove `viewIssueDetail` and `IssueDetailModel`, or add a gesture to enter it (e.g., double Enter or a dedicated key) for full-screen issue reading.

---

### TD-02 — Transition List is Not Cached

Every time the user opens the transition modal (`t`), a fresh API call is made. For teams with stable transition lists, this is wasteful.

**Impact:** Slightly slower response; extra API calls.

**Recommendation:** Cache transitions per issue key in the root model. Invalidate on `TransitionDoneMsg`.

---

### TD-03 — Ollama Model is Hardcoded

The Ollama model name (`llama3`) is a constant in `internal/ollama/client.go`. Users running different models (Mistral, Llama 3.1, Gemma, etc.) cannot configure this.

**Impact:** Forces users to pull the `llama3` model specifically.

**Recommendation:** Add an optional `ollamaModel` field to `config.json`, defaulting to `"llama3"`.

---

### TD-04 — Status Message Not Auto-cleared

`m.statusMsg` (e.g., "Copied!", "Transition applied", "Authenticated!") is set but never cleared. It remains visible in the status bar until something overwrites it.

**Impact:** Minor visual noise.

**Recommendation:** Clear `statusMsg` after a timeout using a `tea.Tick` command, or clear it on the next key press.

---

### TD-05 — Single Page Issue List (No Pagination)

`ListAssigned` hardcodes `maxResults: 50` and does not implement pagination. Users with more than 50 active assigned issues will not see them all.

**Impact:** Incomplete issue list for busy assignees.

**Recommendation:** Add `startAt` offset support and a "load more" action in the list view (or infinite scroll triggered when the cursor reaches the last item).

---

### TD-06 — `KeyList` and `KeyVimRight` Both Equal `"l"`

In `shared/keys.go`, `KeyList = "l"` and `KeyVimRight = "l"` are defined with the same value. The `l` key triggers "open list selector" in global key handling, which makes vim-right navigation unavailable when in a view that uses it. Currently harmless because `KeyVimRight` is not used in any switch case, but it is a latent naming inconsistency.

**Impact:** Confusing constant naming; potential future bug if `KeyVimRight` is used in a switch.

**Recommendation:** Remove `KeyVimRight` since the list selector intentionally captures `l`, or document why both exist with the same value.
