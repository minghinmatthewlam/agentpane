package domain

import (
	"fmt"
	"strings"
)

func ParsePaneType(s string) (PaneType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "codex":
		return PaneCodex, nil
	case "claude":
		return PaneClaude, nil
	case "shell":
		return PaneShell, nil
	default:
		return "", fmt.Errorf("unknown pane type: %s (use codex, claude, or shell)", s)
	}
}

func InferPaneType(command, title string) PaneType {
	cmd := strings.ToLower(command)
	switch {
	case strings.Contains(cmd, "codex"):
		return PaneCodex
	case strings.Contains(cmd, "claude"):
		return PaneClaude
	}

	t := strings.ToLower(strings.TrimSpace(title))
	switch {
	case strings.HasPrefix(t, "codex"):
		return PaneCodex
	case strings.HasPrefix(t, "claude"):
		return PaneClaude
	case strings.HasPrefix(t, "shell"):
		return PaneShell
	}
	return PaneUnknown
}

func InferPaneTypeFromPane(p Pane) PaneType {
	return InferPaneType(p.CurrentCommand, p.Title)
}
