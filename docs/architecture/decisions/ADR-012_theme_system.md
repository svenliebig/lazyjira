# ADR-012 — Configurable Theme System

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03 |
| Deciders | Project team |

## Context

The application had a single hard-coded color palette defined as package-level `lipgloss.Style` variables in `internal/tui/shared/styles.go`. Users had no way to change the visual appearance without modifying source code. There was also no mechanism for general user preferences to be persisted.

## Decision

Introduce a two-layer theme system:

1. **`internal/theme`** — owns the `Theme` struct, a list of predefined themes, the `Current` global, and `LoadCustom()` which reads user-defined themes from `~/.config/lazyjira/themes.json`.
2. **`internal/settings`** — owns a general `Settings` struct (currently `ActiveTheme string`) persisted to `~/.config/lazyjira/settings.json`.

Styles in `shared/styles.go` are rebuilt from `theme.Current` by a `RefreshStyles()` function. At runtime, switching a theme is: `theme.SetTheme(t); shared.RefreshStyles()`.

A settings modal (key `s`) lets the user pick a theme interactively. The choice is saved immediately to `settings.json` so it is applied on the next launch.

## Rationale

**Separate `theme` and `settings` packages:**
Themes are a standalone concept that could theoretically be loaded without the TUI. `settings` is also a standalone concept for future additional preferences. Keeping them separate avoids a bloated config package and gives each a single responsibility.

**`RefreshStyles()` over style functions:**
Changing all call sites from `shared.StyleHeader.Render(…)` to `shared.StyleHeader().Render(…)` would touch every view and modal. The `RefreshStyles()` approach — rebuilding the same exported vars — keeps all call sites unchanged and is consistent with the existing style.

**Predefined themes in code, custom themes in a file:**
Predefined themes ship with the binary and require no file system access. Custom themes are an opt-in extension point. Merging them in `theme.All(custom)` keeps the modal simple and uniform.

**XDG compliance:**
`settings.json` and `themes.json` follow the same `$XDG_CONFIG_HOME/lazyjira/` path convention as `config.json` and `exclusions.json`, keeping all user data in one directory.

## Implementation

```
startup:
  settings.Load()       → ActiveTheme = "default" (or saved value)
  theme.LoadCustom()    → []Theme from themes.json (or nil)
  theme.FindByName(activeName, custom) → Theme
  theme.SetTheme(t); shared.RefreshStyles()
  tui.New(cfg, client, store, settings, customThemes)

on s key:
  theme.All(customThemes) → combined list
  open SettingsModal with list + current name

on ThemeSelectedMsg{Name}:
  theme.FindByName(name, customThemes) → Theme
  theme.SetTheme(t); shared.RefreshStyles()
  settings.ActiveTheme = name
  settings.Save()
  close modal
```

## Consequences

- Adding a new predefined theme requires only a new entry in `theme.Predefined` — no other changes needed.
- `settings.json` currently stores only `activeTheme`; additional preferences (e.g., keybinding profiles) can be added to the `Settings` struct without breaking existing files.
- `RefreshStyles()` rebuilds all styles synchronously on the main goroutine. This is negligible cost for the number of styles involved.
- Custom theme authors must know all ten field names. No validation is performed beyond JSON parsing; unknown fields are ignored by `encoding/json`.
