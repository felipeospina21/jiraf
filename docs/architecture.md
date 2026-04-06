# Architecture Overview

## Module Structure

```
projects/go/
├── tuishell/        # Shared TUI library
├── mrglab/          # GitLab merge requests TUI
├── mrjira/          # Jira issues TUI (this repo)
└── tuishell-hub/    # Launcher that hosts both apps
```

## tuishell (shared library)

Reusable 3-panel TUI layout extracted from mrglab. No `main` package — library only.

### Packages

| Package | Purpose |
|---------|---------|
| `tuishell` (root) | `AppContext`, `Layout`, `ComputeLayout()`, `GlobalKeyMap`, `KeyMatcher()`, `Max/Min/Clamp/Truncate` |
| `style` | `Theme` struct (30 semantic color tokens), `DefaultTheme()`, `MainFrameStyle()`, color palettes (`Blue`, `Red`, `Green`, `Yellow`, `Violet`, `Orange`) |
| `table` | Table widget with scrolling, themed styles, `InitModel()`, `FormatTime()`, `ColWidth()`, `GetColIndex()` |
| `modal` | Overlay modal with dim background, keybindings, copy/submit/close messages |
| `statusline` | Status bar with mode colors (normal/loading/error/dev), `ProjectLabel` slot, spinner style |
| `loader` | Themed loading spinner view |
| `hub` | App picker launcher — list of `AppEntry` items, switches between child `tea.Model`s |

### Theme

Components receive a `style.Theme` at construction. `DefaultTheme()` uses mrglab's violet palette. mrjira can define its own:

```go
theme := style.Theme{
    Primary:       lipgloss.Color("#0052CC"),  // Jira blue
    PrimaryBright: lipgloss.Color("#2684FF"),
    // ... 28 more tokens
}
```

### Layout

`ComputeLayout()` takes a `LayoutConfig` (panel styles, widths) and window size, returns computed dimensions for left/main/right panels and statusline. Supports left panel toggle, right panel toggle, and right panel fullscreen.

## tuishell-hub (launcher)

Single binary that imports both apps and presents a picker.

### How it works

1. User runs `tuishell-hub` → sees app picker list (mrglab, mrjira)
2. `enter` launches selected app — child `tea.Model` receives all messages
3. `ctrl+h` from any app returns to the picker (app state preserved)
4. `ctrl+c` quits

### Import pattern

Each app exposes an `export/` package with `NewApp() tea.Model`:

```go
// tuishell-hub/main.go
apps := []hub.AppEntry{
    {Name: "mrglab", Desc: "GitLab merge requests", NewModel: func() tea.Model { return mrglabExport.NewApp() }},
    {Name: "mrjira", Desc: "Jira issues",           NewModel: func() tea.Model { return mrjiraExport.NewApp() }},
}
m := hub.New(apps)
```

### Dependency graph

```
tuishell-hub
├── imports tuishell/hub
├── imports mrglab/export  → mrglab/internal/... → tuishell/*
└── imports mrjira/export  → mrjira/internal/... → tuishell/*
```

No circular dependencies. `tuishell` never imports `mrglab` or `mrjira`.

## mrjira (this repo)

### Current state

Placeholder app rendering an empty tuishell layout:
- Left panel with a `bubbles/list` (hardcoded project names)
- Center panel with an empty table ("Select a project")
- Statusline at bottom

### File structure

```
mrjira/
├── main.go                      # Standalone entry point
├── export/export.go             # NewApp() for tuishell-hub
├── internal/tui/app/app.go      # Model, Update, View using tuishell layout
└── go.mod                       # replace ../tuishell for local dev
```

### Local development

All modules use `replace` directives for local development:

```
# mrjira/go.mod
replace github.com/felipeospina21/tuishell => ../tuishell

# tuishell-hub/go.mod
replace (
    github.com/felipeospina21/mrglab   => ../mrglab
    github.com/felipeospina21/mrjira   => ../mrjira
    github.com/felipeospina21/tuishell => ../tuishell
)
```

Run standalone: `cd mrjira && go run .`
Run via hub: `cd tuishell-hub && go run .`
