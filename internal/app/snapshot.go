package app

import (
	"errors"
	"os"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/provider"
	"github.com/minghinmatthewlam/agentpane/internal/state"
)

func (a *App) Snapshot() (domain.Snapshot, error) {
	if err := a.applyConfigOverrides(""); err != nil {
		return domain.Snapshot{}, err
	}

	rawSessions, err := a.tmux.ListSessions()
	if err != nil {
		return domain.Snapshot{}, err
	}
	sessions, err := a.convertSessions(rawSessions)
	if err != nil {
		return domain.Snapshot{}, err
	}

	st, err := a.state.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			st = state.NewStore()
		} else {
			st = state.NewStore()
		}
	}

	detector := provider.NewStatusDetector(a.providers)

	for si := range sessions {
		session := &sessions[si]
		stateSession := st.Sessions[session.Name]
		statePaneMap := map[string]*state.PaneState{}
		if stateSession != nil {
			for _, p := range stateSession.Panes {
				statePaneMap[p.TmuxID] = p
			}
		}
		for pi := range session.Panes {
			pane := &session.Panes[pi]
			if sp, ok := statePaneMap[pane.ID]; ok {
				pane.Title = sp.Title
				pane.Type = domain.PaneType(sp.Type)
			} else {
				pane.Type = inferPaneType(pane.CurrentCommand, pane.Title)
			}
			pane.Status = detector.DetectStatus(pane.PID, pane.Type)
		}
	}

	snapshot := domain.Snapshot{
		Sessions: sessions,
	}

	if a.tmux.InTmux() {
		if name, err := a.tmux.CurrentSession(); err == nil {
			snapshot.CurrentSession = name
		}
		if paneID, err := a.tmux.CurrentPane(); err == nil {
			snapshot.CurrentPane = paneID
		}
	}

	return snapshot, nil
}

func inferPaneType(command, title string) domain.PaneType {
	cmd := strings.ToLower(command)
	switch {
	case strings.Contains(cmd, "codex"):
		return domain.PaneCodex
	case strings.Contains(cmd, "claude"):
		return domain.PaneClaude
	}
	t := strings.ToLower(strings.TrimSpace(title))
	switch {
	case strings.HasPrefix(t, "codex"):
		return domain.PaneCodex
	case strings.HasPrefix(t, "claude"):
		return domain.PaneClaude
	case strings.HasPrefix(t, "shell"):
		return domain.PaneShell
	}
	return domain.PaneUnknown
}
