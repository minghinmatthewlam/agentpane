#!/bin/sh
set -euo pipefail

SOCK="agentpane-test-$$"
ROOT="$(pwd)"
TMPROOT="$(mktemp -d)"
HOME_DIR="$TMPROOT/home"
REPO="$TMPROOT/repo"

mkdir -p "$HOME_DIR" "$REPO"

cleanup() {
  tmux -L "$SOCK" kill-server 2>/dev/null || true
  rm -rf "$TMPROOT"
}
trap cleanup EXIT

go build -o "$ROOT/agentpane" "$ROOT/cmd/agentpane"

cd "$REPO"
HOME="$HOME_DIR" AGENTPANE_TMUX_SOCKET="$SOCK" "$ROOT/agentpane" up --detach

tmux -L "$SOCK" has-session -t repo

PANES=$(tmux -L "$SOCK" list-panes -t repo:0 | wc -l | tr -d ' ')
[ "$PANES" -eq 2 ] || (echo "expected 2 panes, got $PANES" && exit 1)

PANE_ID=$(tmux -L "$SOCK" list-panes -t repo:0 -F '#{pane_id}' | head -n1)
TMUX=1 TMUX_PANE="$PANE_ID" HOME="$HOME_DIR" AGENTPANE_TMUX_SOCKET="$SOCK" "$ROOT/agentpane" add shell

PANES=$(tmux -L "$SOCK" list-panes -t repo:0 | wc -l | tr -d ' ')
[ "$PANES" -eq 3 ] || (echo "expected 3 panes, got $PANES" && exit 1)

TMUX=1 TMUX_PANE="$PANE_ID" HOME="$HOME_DIR" AGENTPANE_TMUX_SOCKET="$SOCK" "$ROOT/agentpane" rename "my-title"

TITLE=$(tmux -L "$SOCK" display -p -t "$PANE_ID" '#{pane_title}')
[ "$TITLE" = "my-title" ] || (echo "expected title my-title, got $TITLE" && exit 1)

tmux -L "$SOCK" kill-server

HOME="$HOME_DIR" AGENTPANE_TMUX_SOCKET="$SOCK" "$ROOT/agentpane" up --detach
TITLE=$(tmux -L "$SOCK" display -p -t repo:0.0 '#{pane_title}')
[ "$TITLE" = "codex-1" ] || (echo "expected title codex-1 after restart, got $TITLE" && exit 1)

echo "Integration tests passed"

