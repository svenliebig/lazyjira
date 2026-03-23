# User Experience

The overall user experience of the application is that the tool is a command that open a TUI to interact with the user, similar to lazygit for example.

When in the TUI the user has an overview of possible actions and their shortcuts, like when the user presses `l` (list) the user will see a small modal that asks the user to select the type of list they want to see. So following with `a` (assigned) will list the assigned issues.

The `esc` key will always return the user to the previous screen or close the current modal etc.

When pressing `?` the user will see an overview of the possible actions and their shortcuts.

Navigation in the tool is through the arrow keys or the vim-like navigation keys `h`, `j`, `k` and `l`.

## Action modal navigation

All action-selection modals (copy, transition, list selector, exclude) support cursor navigation:

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` / `l` | Confirm highlighted selection |
| `h` / `esc` | Cancel / close modal |

The currently highlighted option is shown with a `>` prefix in bold purple. Direct shortcut keys (e.g., `u` for copy URL, `1`–`9` for transitions) remain available as single-keystroke aliases where they do not conflict with the navigation keys. No action shortcut inside a modal may use `h`, `j`, `k`, or `l`.
