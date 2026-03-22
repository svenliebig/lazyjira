# 6. Runtime View

## Scenario 1 — Application Startup (First Run, No Config)

```
main()
  │
  ├─ config.Load(flags)
  │   └─ reads ~/.config/jira-cli/config.json  → not found, returns empty Config
  │
  ├─ cfg.IsComplete() → false  (no credentials)
  │
  ├─ jiraClient = nil
  │
  ├─ exclusions.Load()
  │   └─ reads ~/.config/jira-cli/exclusions.json → not found, returns empty Store
  │
  └─ tui.New(cfg, nil, store)
       └─ cfg incomplete → activeModal = modalAuth
       └─ tea.NewProgram(model, tea.WithAltScreen()).Run()
            │
            ├─ WindowSizeMsg → updateChildSizes()
            │
            └─ authModal rendered as overlay
                 └─ user fills URL, email, token, presses Enter
                      └─ AuthCompletedMsg{URL, Email, Token}
                           ├─ config.Save(cfg) → ~/.config/jira-cli/config.json
                           ├─ jiraClient = jira.NewClient(url, email, token)
                           └─ activeModal = modalNone → home view shown
```

---

## Scenario 2 — Listing Assigned Issues

```
User presses "l"
  └─ handleKey("l")
       └─ activeModal = modalListSelector
            └─ listSelectorModal.View() renders "a - Assigned Issues"

User presses "a"
  └─ listSelectorModal.Update("a")
       └─ returns ListSelectedMsg{Type: "assigned"}

Root model receives ListSelectedMsg
  ├─ activeModal = modalNone
  ├─ loading = true
  └─ return fetchAssignedCmd(jiraClient)   ← tea.Cmd, runs in goroutine

fetchAssignedCmd executes:
  └─ jiraClient.ListAssigned(ctx)
       └─ POST /rest/api/3/search/jql
            body: {jql: "assignee = currentUser() AND statusCategory != Done", fields: [...]}
            header: Authorization: Basic base64(email:token)
       └─ Parse JSON → []Issue (with ADF-to-text conversion for description)
       └─ return IssueListLoadedMsg{Issues: [...]}

Root model receives IssueListLoadedMsg
  ├─ allIssues = msg.Issues                              ← raw, unfiltered
  ├─ filtered = exclusions.Filter(allIssues)             ← apply local rules
  ├─ issueListView = NewIssueListModel(filtered, width, height-2)
  ├─ currentView = viewIssueList
  ├─ currentIssue = issueListView.CurrentIssue()         ← first visible issue
  └─ loading = false → split-panel rendered
```

---

## Scenario 3 — Navigating the Split Panel

```
User presses "j" (or ↓)
  └─ handleKey → no chord, no modal, not Esc
       └─ updateActiveChild(KeyMsg{"j"})
            └─ issueListView.Update(KeyMsg{"j"})
                 ├─ focusRight = false → delegate to list.Model
                 ├─ list.Model handles "j" → cursor moves down
                 ├─ list.Index() changed
                 ├─ detail.SetContent(buildIssueDetail(newIssue))
                 └─ detail.GotoTop()
            └─ m.issueListView = updated
            └─ m.currentIssue = issueListView.CurrentIssue()  ← new issue

User presses "Enter"
  └─ issueListView.Update(KeyMsg{"enter"})
       └─ focusRight = false, issues not empty
       └─ focusRight = true
       └─ status bar shows "j/k:scroll  esc:back  o:open  y:copy  t:transition  a:AI"

User presses "j" (right panel focused)
  └─ issueListView.Update(KeyMsg{"j"})
       └─ focusRight = true → delegate to viewport.Model
       └─ viewport scrolls down

User presses "Esc"
  └─ handleKey("esc")
       └─ issueListView.IsFocusRight() = true
       └─ issueListView = issueListView.BlurRight()
       └─ focusRight = false, list regains focus
```

---

## Scenario 4 — Transitioning an Issue

```
User presses "t" (currentIssue != nil)
  └─ handleKey("t")
       └─ loading = true
       └─ return fetchTransitionsCmd()  ← tea.Cmd

fetchTransitionsCmd executes:
  └─ jiraClient.GetTransitions(ctx, "PROJ-42")
       └─ GET /rest/api/3/issue/PROJ-42/transitions
       └─ return TransitionsLoadedMsg{Transitions: [...]}

Root model receives TransitionsLoadedMsg
  ├─ transitionModal = NewTransitionModal(transitions)
  ├─ activeModal = modalTransition
  └─ loading = false

transitionModal renders numbered list:
  "1  In Progress → In Review"
  "2  In Review → Done"

User presses "1"
  └─ transitionModal.Update(KeyMsg{"1"})
       └─ transitions[0].ID = "11"
       └─ return TransitionSelectedMsg{ID: "11"}

Root model receives TransitionSelectedMsg
  ├─ activeModal = modalNone
  ├─ loading = true
  └─ return doTransitionCmd("11")  ← tea.Cmd

doTransitionCmd executes:
  └─ jiraClient.DoTransition(ctx, "PROJ-42", "11")
       └─ POST /rest/api/3/issue/PROJ-42/transitions
            body: {transition: {id: "11"}}
       └─ 204 No Content → success
       └─ return TransitionDoneMsg{}

Root model receives TransitionDoneMsg
  └─ loading = false
  └─ statusMsg = "Transition applied"
```

---

## Scenario 5 — AI Work Summary

```
User presses "a" (currentIssue = PROJ-42)
  └─ handleKey("a") → pendingKey = "a", aiModal created, activeModal = modalAI

User presses "s"
  └─ aiModal.Update(KeyMsg{"s"})
       └─ state = aiLoadingCommits
       └─ return tea.Batch(spinner.Tick, fetchCommitsCmd())

fetchCommitsCmd executes:
  └─ git.CommitsForIssue("PROJ-42")
       └─ exec "git log --oneline --all"
       └─ filter lines containing "PROJ-42"
       └─ return AICommitsLoadedMsg{Commits: ["abc1234 PROJ-42 implement feature"]}

Root model receives AICommitsLoadedMsg
  └─ aiModal.SetCommits(commits)  → state = aiGenerating
  └─ return aiModal.GenerateCmd()

GenerateCmd executes:
  └─ builds prompt:
       "Summarize the following Jira issue and related git commits in 2-3 sentences.
        Issue: PROJ-42
        Summary: Implement login feature
        Description: ...
        Commits:
        abc1234 PROJ-42 implement feature"
  └─ ollama.Client.Generate(ctx, prompt)
       └─ POST http://localhost:11434/api/generate
            body: {model: "llama3", prompt: "...", stream: false}
       └─ return AISummaryMsg{Summary: "The team implemented..."}

Root model receives AISummaryMsg
  └─ aiModal.SetSummary(summary)  → state = aiDone, viewport activated
```

---

## Scenario 6 — Copy Issue Key

```
User presses "y" (currentIssue = PROJ-42)
  └─ handleKey("y")
       └─ pendingKey = "y"
       └─ copyModal = NewCopyModal()
       └─ activeModal = modalCopy

copyModal renders:
  "k  Copy issue key"
  "u  Copy issue URL"
  "t  Copy issue title"
  "d  Copy description"

User presses "k"
  └─ copyModal.Update(KeyMsg{"k"})
       └─ return CopyActionMsg{Action: "key"}

Root model receives CopyActionMsg{Action: "key"}
  └─ text = currentIssue.Key  → "PROJ-42"
  └─ clipboard.Write("PROJ-42")
  └─ activeModal = modalNone
  └─ pendingKey = ""
  └─ statusMsg = "Copied!"
```

---

## Scenario 7 — Excluding an Issue

```
User presses "x" (currentIssue = PROJ-42, viewIssueList)
  └─ handleKey("x")
       └─ currentView != viewExcludedList, currentIssue != nil
       └─ excludeModal = NewExcludeModal(currentIssue)
       └─ activeModal = modalExclude

excludeModal renders:
  "p  Exclude all by parent issue (PROJ-10)"   ← if parent exists
  "k  Exclude by issue key (PROJ-42)"

  — or, if no parent —

  "p  ~~Exclude all by parent issue (no parent)~~"  ← strikethrough, non-interactive
  "k  Exclude by issue key (PROJ-42)"

User presses "k"
  └─ excludeModal.Update(KeyMsg{"k"})
       └─ return ExcludeActionMsg{Type: "key", Value: "PROJ-42"}

Root model receives ExcludeActionMsg{Type: "key", Value: "PROJ-42"}
  ├─ exclusions.Add(Rule{Type: "key", Value: "PROJ-42"})
  │   └─ appends to rules, writes ~/.config/jira-cli/exclusions.json
  ├─ filtered = exclusions.Filter(allIssues)    ← PROJ-42 now absent
  ├─ issueListView = NewIssueListModel(filtered, ...)
  ├─ currentIssue = issueListView.CurrentIssue()
  ├─ activeModal = modalNone
  └─ statusMsg = "Issue excluded"
```

Note: `allIssues` is not modified. The filtered slice is recomputed from the unchanged raw data on every exclusion change. No network request is needed.

---

## Scenario 8 — Viewing and Removing Exclusions

```
User presses "l"
  └─ activeModal = modalListSelector

listSelectorModal renders:
  "a  Assigned Issues"
  "x  Excluded Issues"

User presses "x"
  └─ listSelectorModal.Update(KeyMsg{"x"})
       └─ return ListSelectedMsg{Type: "excluded"}

Root model receives ListSelectedMsg{Type: "excluded"}
  ├─ excludedListView = NewExcludedListModel(exclusions.Rules(), width, height-2)
  ├─ currentView = viewExcludedList
  └─ currentIssue = nil

excludedListView renders:
  "PROJ-42          Excluded by issue key"
  "Parent: PROJ-10  All issues with this parent are excluded"

User navigates to "PROJ-42" row and presses "x"
  └─ handleKey("x")
       └─ currentView == viewExcludedList
       └─ rule = excludedListView.CurrentRule() → Rule{Type:"key", Value:"PROJ-42"}
       └─ exclusions.Remove(rule)
       │   └─ removes from rules, writes ~/.config/jira-cli/exclusions.json
       └─ excludedListView = NewExcludedListModel(exclusions.Rules(), ...)
       └─ statusMsg = "Exclusion removed"
```
