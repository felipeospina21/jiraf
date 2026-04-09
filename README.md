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

| Key   | Action    |
| ----- | --------- |
| `↑/k` | Move up   |
| `↓/j` | Move down |

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
| `filters.boards` | List of board objects | —       | Yes      |

### Board object

Each board in `filters.boards` has the following fields:

| Field  | Type     | Description                          |
| ------ | -------- | ------------------------------------ |
| `name` | `string` | Display name shown in the board list |
| `id`   | `string` | Jira board ID                        |
| `key`  | `string` | Jira project key (e.g. `PROJ`)       |

### Config example

```toml
base_url = "https://your-org.atlassian.net"

[filters]
boards = [
    { name = "Some Board", id = "3588", key = "SMB" },
    { name = "Platform Sprint", id = "42", key = "PLAT" },
]
```

## Related

- [mrglab](https://github.com/felipeospina21/mrglab) — GitLab merge requests TUI
- [tuishell](https://github.com/felipeospina21/tuishell) — Shared TUI framework

## License

This project is licensed under the [MIT License](LICENSE).
