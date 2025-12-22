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
make build      # build only
make install    # build and install to $GOPATH/bin
```

## Quick start

```bash
# Open the interactive dashboard
agentpane

# Or create/attach to a session for the current repo
agentpane up
```

Running `agentpane up` will:
1. Create a tmux session named after your repo (or attach if it exists)
2. Set up panes using the default template (Codex + Claude)
3. Auto-add a tmux keybinding (`prefix + g`) to your `~/.tmux.conf` for quick dashboard access

## Commands

| Command | Description |
|---------|-------------|
| `agentpane` | Open the interactive dashboard (default) |
| `agentpane up` | Create/attach session for current repo |
| `agentpane up --template <name>` | Use a specific template |
| `agentpane add codex\|claude\|shell` | Add a pane to current session |
| `agentpane rename [name]` | Rename current pane |
| `agentpane dashboard` | Open interactive dashboard |
| `agentpane dashboard --tmux-window` | Open dashboard in a dedicated tmux window |
| `agentpane templates` | List available templates |
| `agentpane templates --apply <name>` | Apply template to current session |
| `agentpane templates --apply <name> --force` | Replace existing panes with template |
| `agentpane search [query]` | Search sessions and panes |
| `agentpane init` | Generate `.agentpane.yml` config for repo |

## Dashboard

The dashboard provides a tree view of all sessions and their panes, with live preview of pane content.

### Keyboard shortcuts

| Key | Action |
|-----|--------|
| `↑/↓` | Navigate tree (sessions and panes) |
| `←/→` | Jump between sessions |
| `Tab` | Switch tab (Sessions / Templates) |
| `Enter` | Attach to session / Apply template |
| `/` | Filter sessions by name |
| `o` | Open new session (folder picker) |
| `c` | Quick-add Claude pane |
| `x` | Quick-add Codex pane |
| `s` | Quick-add Shell pane |
| `a` | Add pane (type selection dialog) |
| `r` | Rename pane (when cursor on pane) |
| `d` | Close pane (when cursor on pane) |
| `k` | Kill session (when cursor on session) |
| `?` | Show help |
| `q` | Quit dashboard |

## tmux keybinding

For quick access to the dashboard from anywhere in tmux, add to `~/.tmux.conf`:

```bash
bind-key g run-shell "agentpane dashboard --tmux-window"
```

Then reload:

```bash
tmux source-file ~/.tmux.conf
```

> **Note:** `agentpane up` automatically adds this keybinding if not already present.

Now press `prefix + g` (e.g., `Ctrl-b g`) to open the dashboard in a dedicated tmux window.

## Configuration

### Repo config: `.agentpane.yml`

Place in your repo root to customize session name and default panes:

```yaml
session: my-app   # optional custom session name
layout:
  panes:
    - type: codex
      title: "Codex"      # optional
    - type: claude
      title: "Claude"
    - type: shell
      title: "Dev Server"
```

Generate a starter config:

```bash
agentpane init
```

### Global config: `~/.config/agentpane/config.yml`

```yaml
default_template: duo

providers:
  codex:
    command: codex           # command to run
    args: []                 # optional args
  claude:
    command: claude
  shell:
    command: $SHELL
```

## Templates

Built-in templates:

| Template | Panes |
|----------|-------|
| `simple` | 1 Claude |
| `duo` | Codex + Claude (default) |
| `trio` | Codex + Claude + Shell |
| `quad` | 2x Codex + 2x Claude |
| `full` | 2x Codex + 2x Claude + Shell |

Use a template:

```bash
agentpane up --template quad
```

Apply template to existing session:

```bash
agentpane templates --apply trio --force
```

## Environment variables

| Variable | Description |
|----------|-------------|
| `AGENTPANE_TMUX_SOCKET` | Use a custom tmux socket (useful for testing) |

## Notes

- Works on macOS and Linux
- Requires tmux 3.2+ for popup support (falls back to windows on older versions)
- Pane titles persist across tmux restarts when pane IDs still exist
- If Codex or Claude aren't in PATH, panes fall back to shell

## License

MIT
