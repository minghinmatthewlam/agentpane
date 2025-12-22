package app

import (
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

	st := a.loadStateOrNew()

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
				pane.Type = domain.InferPaneType(pane.CurrentCommand, pane.Title)
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
