# 🦒 jiraf

# Jiraf

![screenshot](docs/assets/new_jiraf_logo.png)

[![Go Version](https://img.shields.io/github/go-mod/go-version/felipeospina21/jiraf)](https://github.com/felipeospina21/jiraf)
[![License](https://img.shields.io/github/license/felipeospina21/jiraf)](https://github.com/felipeospina21/jiraf/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/felipeospina21/jiraf)](https://goreportcard.com/report/github.com/felipeospina21/jiraf)

jiraf is a TUI to manage Jira issues from the command line. Built on [tuishell](https://github.com/felipeospina21/tuishell).

### Why "jiraf"?

**ji**ra + **f**etch = **jiraf** — and it sounds like _giraffe_ 🦒, because sometimes you need a long neck to see across all your Jira boards.

## Requirements

- Nerd Font (Symbols) v3.2.1 or higher [ download ](https://github.com/ryanoasis/nerd-fonts/releases/download/v3.2.1/NerdFontsSymbolsOnly.zip)

## Install

```bash
go install github.com/felipeospina21/jiraf@latest
```

### From source

```bash
git clone https://github.com/felipeospina21/jiraf.git
cd jiraf
go build -o jiraf .
```

To make it available system-wide:

```bash
sudo mv jiraf /usr/local/bin/
```

## Usage

Launch the TUI:

```bash
jiraf
```

### Dev mode

Run with mocked data (no API calls):

```bash
jiraf -dev
```

## Features

- Browse Jira issues for your configured boards
- View issues by board with sprint filtering
- Filter issues server-side by status, priority, type, and sprint number
- View issue details with markdown rendering and comments
- Transition issue status with confirmation
- Open issues in browser
- Custom theme colors via config with preset support
- Navigate between boards and issues panels

## Keybindings

### Global

| Key      | Action            |
| -------- | ----------------- |
| `?`      | Toggle help       |
| `ctrl+c` | Quit              |
| `ctrl+o` | Toggle side panel |

### Boards panel

| Key     | Action      |
| ------- | ----------- |
| `enter` | View issues |

### Issues panel

| Key   | Action          |
| ----- | --------------- |
| `↑/k` | Move up         |
| `↓/j` | Move down       |
| `enter/l` | View details |
| `x`   | Open in browser |
| `T`   | Transition status |
| `r`   | Refetch issues  |
| `/`   | Open filter     |
| `F`   | Clear filters   |

## Config

Config file is read from `~/.config/jiraf/jiraf.toml`.

**You need to set a Jira API token as an environment variable. Add it to your shell config (e.g. `.zshrc`) to persist it:**

```bash
export JIRAF_TOKEN="YOUR_JIRA_API_TOKEN"
```

### Config properties

| Option           | Description           | Default | Required |
| ---------------- | --------------------- | ------- | -------- |
| `base_url`       | Jira instance URL     | —       | Yes      |
| `filters.boards` | List of board objects  | —       | Yes      |
| `theme`          | Color overrides (see [Theme](#theme)) | Jira blue | No |

### Board object

Each board in `filters.boards` has the following fields:

| Field  | Type     | Description                          |
| ------ | -------- | ------------------------------------ |
| `name` | `string` | Display name shown in the board list |
| `id`   | `string` | Jira board ID                        |
| `key`  | `string` | Jira project key (e.g. `PROJ`)       |

### Theme

You can customize the color theme by adding a `[theme]` section to your config file.

**Presets:** Use a preconfigured palette by name. Individual overrides can be combined with a preset.

| Preset                 | Description                  |
| ---------------------- | ---------------------------- |
| `catppuccin-mocha`     | Catppuccin Mocha (dark)      |
| `catppuccin-macchiato` | Catppuccin Macchiato         |
| `catppuccin-frappe`    | Catppuccin Frappé            |
| `catppuccin-latte`     | Catppuccin Latte (light)     |
| `rose-pine`            | Rosé Pine                    |
| `tokyo-night`          | Tokyo Night                  |
| `dracula`              | Dracula                      |

```toml
[theme]
preset = "catppuccin-mocha"
```

**Individual overrides:** Each property accepts a hex color string. Only the tokens you specify are overridden — everything else falls back to the preset (or the default Jira blue palette if no preset is set).

| Token             | Description                        | Default   |
| ----------------- | ---------------------------------- | --------- |
| `primary`         | Primary accent color               | `#2684FF` |
| `primary_bright`  | Brighter primary variant           | `#8AC4FF` |
| `primary_fg`      | Foreground on primary backgrounds  | `#DEEBFF` |
| `primary_dim`     | Dimmed primary for subtle accents  | `#002B6B` |
| `info`            | Informational color                | `#4C9AFF` |
| `info_bright`     | Brighter info variant              | `#2684FF` |
| `success`         | Success / positive state           | `#6beaaf` |
| `success_bright`  | Brighter success variant           | `#3ad994` |
| `danger`          | Error / destructive state          | `#f9a8a8` |
| `danger_bright`   | Brighter danger variant            | `#f47575` |
| `warning`         | Warning state                      | `#ffe043` |
| `warning_bright`  | Brighter warning variant           | `#ffcc14` |
| `caution`         | Caution / attention state          | `#ff8237` |
| `text`            | Default text color                 | `#C4C4C4` |
| `text_inverse`    | Text on light backgrounds          | `#111`    |
| `text_dimmed`     | Dimmed / secondary text            | `#777777` |
| `muted`           | Muted text                         | `#999999` |
| `dim`             | Dimmed UI elements                 | `#444444` |
| `border`          | Panel borders                      | `#3f4145` |
| `modal_border`    | Modal overlay border               | `#666666` |
| `surface_dim`     | Dark surface background            | `#091E42` |
| `selection_border`| Selected item border               | `#0052CC` |
| `status_text`     | Status bar text                    | `#FFFDF5` |
| `status_normal`   | Status bar normal mode             | `#0747A6` |
| `status_loading`  | Status bar loading state           | `#1A7A94` |
| `status_error`    | Status bar error state             | `#CE3060` |
| `status_dev`      | Status bar dev mode indicator      | `#4E8212` |
| `status_accent1`  | Status bar accent 1                | `#0065FF` |
| `status_accent2`  | Status bar accent 2                | `#003884` |

### Config example

```toml
base_url = "https://your-org.atlassian.net"

[filters]
boards = [
    { name = "Some Board", id = "3588", key = "SMB" },
    { name = "Platform Sprint", id = "42", key = "PLAT" },
]

# Optional: use a preset theme
[theme]
preset = "catppuccin-mocha"

# Or override specific colors (works with or without a preset)
# primary = "#0052CC"
# success = "#22C55E"
# danger = "#EF4444"
```

## Related

- [mrglab](https://github.com/felipeospina21/mrglab) — GitLab merge requests TUI
- [tuishell](https://github.com/felipeospina21/tuishell) — Shared TUI framework

## License

This project is licensed under the [MIT License](LICENSE).
