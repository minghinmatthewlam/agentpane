# agentpane

`agentpane` is a tmux-based environment for managing AI coding agent panes (Codex CLI, Claude Code, and shells) across repos.

## Installation

### Homebrew (recommended)

```bash
brew install minghinmatthewlam/tap/agentpane
```

### curl (no Go required)

```bash
curl -fsSL https://raw.githubusercontent.com/minghinmatthewlam/agentpane/main/install.sh | bash
```

Optional overrides:

```bash
VERSION=0.1.0 INSTALL_DIR=~/.local/bin \
  curl -fsSL https://raw.githubusercontent.com/minghinmatthewlam/agentpane/main/install.sh | bash
```

### Go install

```bash
go install github.com/minghinmatthewlam/agentpane/cmd/agentpane@latest
```

### From source

```bash
git clone https://github.com/minghinmatthewlam/agentpane.git
cd agentpane
make build
```

## Quick start

```bash
agentpane up
```

This creates a tmux session for the current repo using the default template (Codex + Claude).

## Commands

- `agentpane up` — create/attach session for current repo
- `agentpane add codex|claude|shell` — add a pane
- `agentpane rename [name]` — rename current pane
- `agentpane dashboard` — interactive dashboard
- `agentpane templates` — list templates
- `agentpane templates --apply <name> [--force]` — apply template
- `agentpane search` — search sessions/panes
- `agentpane init` — generate `.agentpane.yml`

### Dashboard keys

- `←/→` or `h/l`: switch session
- `↑/↓` or `j/k`: navigate
- `Tab`: switch focus
- `/`: filter sessions
- `Enter`: attach / apply template
- `c`: quick-add Claude pane
- `x`: quick-add Codex pane
- `s`: quick-add Shell pane
- `a`: add pane (dialog)
- `r`: rename pane
- `d`: close pane
- `q`: quit

## Config

### Repo config: `.agentpane.yml`

```yaml
session: my-app   # optional
layout:
  panes:
    - type: codex
    - type: claude
```

### Global config: `~/.config/agentpane/config.yml`

```yaml
default_template: duo
providers:
  codex:
    command: codex
  claude:
    command: claude
```

## Templates

Built-ins: `simple`, `duo`, `trio`, `quad`, `full`.

```bash
agentpane up --template quad
```

## tmux keybinding (recommended)

Add to `~/.tmux.conf`:

```bash
bind-key g run-shell "agentpane popup"
```

Reload:

```bash
tmux source-file ~/.tmux.conf
```

## Notes

- Works on macOS and Linux.
- Pane titles persist across tmux restarts when pane IDs still exist.
- For integration tests, use `AGENTPANE_TMUX_SOCKET` to isolate tmux.

## License

MIT
