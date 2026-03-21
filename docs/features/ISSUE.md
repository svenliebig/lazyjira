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
- `a` opens the modal for AI assistance
  - `a-s` creates a summary of the work on the issue, for that we use a [Local LLM](./LOCAL_LLM.md) to generate a summary of the work done. The work is done is resulting from the commits that are linked to the issue, for that the tool can query the git log that is associated with the issue. This action has to be performed while in the directory of the repository. If there is no git repository in the directory, the tool will prompt the user to navigate to a git repository when calling the tool.
- `t` opens the modal for transition the issue to a new status, the available transitions are listed in the modal, mostly fetched beforehand somehow cached or fetched on demand. The user can then select the new status by pressing the number of the status.
