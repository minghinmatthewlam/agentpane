package app

import (
	"fmt"
	"os"
	"strings"
)

const overlayPaneTitle = "agentpane-status"

type OverlayOptions struct {
	Position string
	Size     int
}

func (a *App) ToggleOverlay(opts OverlayOptions) error {
	if !a.tmux.InTmux() {
		return fmt.Errorf("must be run inside tmux")
	}

	session, err := a.tmux.CurrentSession()
	if err != nil {
		return err
	}
	window, err := a.tmux.CurrentWindowIndex()
	if err != nil {
		return err
	}

	panes, err := a.tmux.ListPanesInWindow(session, window)
	if err != nil {
		return err
	}
	for _, pane := range panes {
		if pane.Title == overlayPaneTitle {
			return a.tmux.KillPane(pane.ID)
		}
	}

	size := opts.Size
	if size <= 0 {
		size = 30
	}
	position := strings.TrimSpace(opts.Position)
	if position == "" {
		position = "left"
	}

	exe, err := os.Executable()
	if err != nil || strings.TrimSpace(exe) == "" {
		return fmt.Errorf("failed to resolve agentpane path: %w", err)
	}
	command := shellQuote(exe) + " overlay-view"
	paneID, err := a.tmux.SplitPaneWithCommand(position, size, command)
	if err != nil {
		return err
	}
	return a.tmux.SetPaneTitle(paneID, overlayPaneTitle)
}

func shellQuote(input string) string {
	if input == "" {
		return "''"
	}
	if !strings.ContainsAny(input, " \t\n'\"\\$&;|<>()[]{}*?") {
		return input
	}
	// Single-quote and escape embedded single quotes.
	return "'" + strings.ReplaceAll(input, "'", `'"'"'`) + "'"
}
