# jiraf Roadmap

## 1. Critical — Shell Integration Issues

These issues violate the tuishell message contract or cause rendering bugs.

### 1.1 ~~Table does not resize on WindowSizeMsg when rows exist~~ ✅

**File:** `internal/tui/issues/issues.go` — `Update` → `tea.WindowSizeMsg`

The `WindowSizeMsg` handler only sets `m.Table.W` and `m.Table.H` hint fields but never calls `SetColumns`/`SetWidth`/`SetHeight` on the underlying table when rows already exist. When the shell recomputes layout (e.g. left panel toggle, window resize after data loads), the table keeps its old column widths and internal viewport height, causing overflow that hides the statusline and wraps columns.

**Fix:** When `len(m.Table.Rows()) > 0`, call `m.Table.SetColumns(getTableCols(tw))`, `m.Table.SetWidth(tw)`, and `m.Table.SetHeight(h)` — same as the `FetchedMsg` handler does.

### 1.2 ~~FetchedMsg handling bypasses shell.Update~~ ✅

**File:** `internal/tui/app/app.go` — `Update` → `issues.FetchedMsg`

The app intercepts `FetchedMsg`, manually sets `Loading = false`, routes to the panel via `m.Shell.Main.Update(msg)`, and returns early. This skips `m.Shell.Update(msg)` entirely for that tick, so the shell never processes the `spinner.TickMsg` or any other pending state on that cycle. The `FinishTaskMsg` and `SetStatusMsg` returned by the panel come back as commands on the next tick, which works, but the manual `Loading = false` + direct panel routing duplicates what the message flow should handle.

**Fix:** Remove the `FetchedMsg` case from `app.go`. Let the shell's `default` case route it to the main panel via `routeToPanel`. The issues panel already handles `FetchedMsg` and returns `FinishTaskMsg`. The app only needs to sync `SpinnerView` after `m.Shell.Update`, which it already does.

### 1.3 ~~SetStatusMsg sent without Mode resets the mode label~~ ✅

**File:** `internal/tui/issues/issues.go` — `Update` → `FetchedMsg`

```go
func() tea.Msg { return tuishell.SetStatusMsg{Content: fmt.Sprintf("%d issues", count)} }
```

`Mode` is zero-value `""`, so the shell sets `m.Statusline.Status = ""`, rendering an empty mode label in the statusline. The mode color falls through to `StatusNormal` but the label text is blank.

**Fix:** Either drop the `SetStatusMsg` entirely (the `FinishTaskMsg` already resets to normal mode), or set `Mode: statusline.ModesEnum.Normal` explicitly.

### 1.4 ~~FinishTaskMsg hardcodes GlobalKeys(false), ignoring DevMode~~ ✅

**File:** `internal/tui/issues/issues.go` — `Update` → `FetchedMsg`

```go
func() tea.Msg { return tuishell.FinishTaskMsg{Keybinds: tuishell.GlobalKeys(false)} }
```

`DevMode` is always `false` regardless of the app's config. If the app runs in dev mode, the dev keybinds disappear after the first fetch.

**Fix:** The issues panel doesn't have access to `DevMode`. Either pass it at construction, or have the app send `SetKeybindsMsg` with the correct keybinds after `FinishTaskMsg` (like mrglab does), or send `FinishTaskMsg` without `Keybinds` and let the app set them.

### 1.5 ~~Boards panel does not set list width on WindowSizeMsg~~ ✅

**File:** `internal/tui/boards/boards.go` — `Update` → `tea.WindowSizeMsg`

Only calls `m.List.SetHeight(msg.Height)` — never sets width. The `bubbles/list` component needs both dimensions to render correctly. If the shell sends a different width (e.g. after a window resize), the list won't adapt.

**Fix:** Call `m.List.SetWidth(msg.Width)` alongside `SetHeight`.
