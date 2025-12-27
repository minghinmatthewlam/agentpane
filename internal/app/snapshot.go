package app

import (
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/agentstate"
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
	idleThreshold := agentstate.IdleThreshold()
	now := time.Now()

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
			pane.AgentStatus = a.detectAgentStatus(*pane, now, idleThreshold)
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

func (a *App) detectAgentStatus(pane domain.Pane, now time.Time, idleThreshold time.Duration) domain.AgentStatus {
	// 1. Check fresh state file first
	state, ok, err := agentstate.Read(pane.ID)
	if err == nil && ok {
		if state.PaneID != "" && state.PaneID != pane.ID {
			ok = false
		}
		if ok && !agentstate.MatchesTool(pane.Type, state.Tool) {
			ok = false
		}
		if ok && agentstate.IsFresh(state, now) {
			if mapped, ok := agentstate.MapStatus(state.State); ok {
				return mapped
			}
		}
	}

	// 2. Output-based classification (Codex/Claude)
	if pane.Type == domain.PaneCodex || pane.Type == domain.PaneClaude {
		if output, err := a.tmux.CapturePaneLines(pane.ID, agentstate.OutputLines()); err == nil {
			// Check prompt FIRST - visible prompt is the strongest idle signal
			if status, ok := agentstate.MatchPrompt(output); ok {
				return status
			}
			// Then check recent output for running keywords
			if status, ok := agentstate.MatchOutput(output); ok {
				return status
			}
		}
	}

	// 3. Activity-based detection
	if !pane.LastActive.IsZero() {
		if now.Sub(pane.LastActive) > idleThreshold {
			return domain.AgentStatusIdle
		}
		return domain.AgentStatusRunning
	}

	// 4. Default to idle (not actively working)
	return domain.AgentStatusIdle
}
