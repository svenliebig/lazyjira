# Issue

Issues represent the Jira issues and are the main entity in the tool.

## Actions

If a user is on an issue, may it be on the details page of the issue or highlighting it in one of the lists of issues, the user has the following possibilities of actions.

- `y` opens the modal for copying
  - `y-k` copies the key of the issue
  - `y-u` copies the URL of the issue
  - `y-t` copies the title of the issue
  - `y-d` copies the description of the issue
- `o` opens the issue in the browser
- `a` opens the assign modal — a fuzzy-searchable list of users who can be assigned to the issue. The list is fetched from Jira for the specific issue. Typing filters by display name or email address; `↑`/`k` and `↓`/`j` navigate; `enter` confirms; `esc` cancels. On success a status message "Assigned to <name>" is shown.
- `m` opens the modal for AI assistance
  - `m-s` creates a summary of the work on the issue, for that we use a [Local LLM](./LOCAL_LLM.md) to generate a summary of the work done. The work is done is resulting from the commits that are linked to the issue, for that the tool can query the git log that is associated with the issue. This action has to be performed while in the directory of the repository. If there is no git repository in the directory, the tool will prompt the user to navigate to a git repository when calling the tool.
- `t` opens the modal for transition the issue to a new status, the available transitions are listed in the modal, mostly fetched beforehand somehow cached or fetched on demand. The user can then select the new status by pressing the number of the status.
- `u` unassigns the current issue immediately. The issue is removed from the assigned list without confirmation since it is directly reversible in Jira. A status message "Unassigned" is shown after success.
- `x` opens the modal for excluding the issue from the list
  - `x-p` excludes all issues that share the same parent issue. This option is only available if the issue has a parent; it is shown as strikethrough and non-interactive otherwise.
  - `x-k` excludes the specific issue by its key. The issue will no longer appear in the assigned issues list until the exclusion is removed via `l → x`.
