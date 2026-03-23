# 2026-03-23 — Theme system & settings modal

## Goal

Make the color scheme configurable through a theme file in the configuration folder. Add a settings screen (key `s`) accessible from the home screen where the user can select a theme from a list of predefined themes or from custom themes defined in a themes file.

---

## Completed work

### New packages

**`internal/theme`**

- `Theme` struct with fields: `name`, `primary`, `success`, `error`, `muted`, `bg`, `surface`, `text`, `subtext`, `border`, `focus`.
- `Predefined []Theme` — seven built-in themes: `default`, `dracula`, `nord`, `catppuccin-mocha`, `catppuccin-macchiato`, `catppuccin-frappe`, `catppuccin-latte`.
- `Current Theme` — the active theme; defaults to `default`.
- `SetTheme(t Theme)` — sets `Current`.
- `FindByName(name string, custom []Theme) (Theme, bool)` — searches predefined then custom by name.
- `All(custom []Theme) []Theme` — returns predefined + custom combined.
- `LoadCustom() ([]Theme, error)` — reads `~/.config/lazyjira/themes.json` (missing file is not an error).

**`internal/settings`**

- `Settings{ActiveTheme string}` — persisted user preferences.
- `Load() (*Settings, error)` — reads `~/.config/lazyjira/settings.json`; returns `{ActiveTheme: "default"}` if absent.
- `Save(s *Settings) error` — writes with `0600` permissions; creates directory if needed.

### Refactored styles (`internal/tui/shared/styles.go`)

All `lipgloss.Style` variables are now package-level `var`s initialised (and re-initialised) by `RefreshStyles()`. An `init()` call ensures the default theme is applied on startup. Callers invoke `theme.SetTheme(t); shared.RefreshStyles()` to switch themes at runtime — no other code needs to change.

`colorPrimary` etc. are no longer package-level `const`-like vars; they are computed inside `RefreshStyles()` from `theme.Current`.

### New modal (`internal/tui/modals/settings_modal.go`)

`SettingsModal` shows all available themes (predefined + custom) as a navigable list. The currently active theme is marked `(active)`. Navigation and confirmation use the standard modal key contract.

| Key | Action |
|-----|--------|
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `enter` / `l` | Apply selected theme |
| `h` / `q` / `esc` | Cancel |

On confirm, emits `ThemeSelectedMsg{Name}`.

### Root model changes (`internal/tui/app.go`)

- Added `modalSettings` state.
- `s` key (globally, no view restriction) opens `SettingsModal` with all themes and the current active theme.
- `ThemeSelectedMsg` handler: calls `theme.SetTheme` + `shared.RefreshStyles()`, updates `appSettings.ActiveTheme`, calls `settings.Save()`.
- `tui.New()` signature extended: `New(cfg, jiraClient, store, appSettings, customThemes)`.
- Status bar now shows `s:settings` alongside `l:list`.

### Startup wiring (`main.go`)

1. Load settings (`settings.Load()`).
2. Load custom themes (`theme.LoadCustom()`).
3. If saved `ActiveTheme` is found, apply it via `theme.SetTheme` + `shared.RefreshStyles()` before the TUI starts.

---

## Custom theme file format

`~/.config/lazyjira/themes.json`:

```json
[
  {
    "name": "my-theme",
    "primary": "#FF6B6B",
    "success": "#6BCB77",
    "error":   "#FF4D4D",
    "muted":   "#888888",
    "bg":      "#1A1A2E",
    "surface": "#16213E",
    "text":    "#FFFFFF",
    "subtext": "#CCCCCC",
    "border":  "#0F3460",
    "focus":   "#FF6B6B"
  }
]
```

Custom themes appear after the predefined themes in the settings modal. The file is optional; if absent, only the predefined themes are available.

---

## Predefined themes

| Name | Character |
|------|-----------|
| `default` | Dark purple on dark grey |
| `dracula` | Classic Dracula |
| `nord` | Arctic Nord |
| `catppuccin-mocha` | Catppuccin Mocha (darkest) |
| `catppuccin-macchiato` | Catppuccin Macchiato |
| `catppuccin-frappe` | Catppuccin Frappé |
| `catppuccin-latte` | Catppuccin Latte (light) |
