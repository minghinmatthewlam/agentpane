package state

import (
	"strings"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

type ReconcileInput struct {
	CurrentState *Store
	TmuxSessions []domain.Session
}

type ReconcileOutput struct {
	UpdatedState  *Store
	TitleUpdates  []TitleUpdate
	OrphanedPanes []string
	NewPanes      []NewPaneInfo
}

type TitleUpdate struct {
	PaneID string
	Title  string
}

type NewPaneInfo struct {
	SessionName  string
	PaneID       string
	InferredType domain.PaneType
}

func Reconcile(input ReconcileInput) ReconcileOutput {
	output := ReconcileOutput{
		UpdatedState: &Store{
			Version:  1,
			Sessions: make(map[string]*SessionState),
		},
	}

	if input.CurrentState != nil {
		output.UpdatedState.Version = input.CurrentState.Version
		output.UpdatedState.ServerID = input.CurrentState.ServerID
	}

	tmuxSessionMap := make(map[string]domain.Session)
	for _, s := range input.TmuxSessions {
		tmuxSessionMap[s.Name] = s
	}

	for name, stateSession := range input.CurrentState.Sessions {
		tmuxSession, exists := tmuxSessionMap[name]
		if !exists {
			continue
		}
		reconciled := reconcileSession(stateSession, tmuxSession, &output)
		output.UpdatedState.Sessions[name] = reconciled
	}

	for name, tmuxSession := range tmuxSessionMap {
		if _, exists := input.CurrentState.Sessions[name]; exists {
			continue
		}
		newSession := createSessionState(tmuxSession, &output)
		output.UpdatedState.Sessions[name] = newSession
	}

	return output
}

func reconcileSession(stateSession *SessionState, tmuxSession domain.Session, output *ReconcileOutput) *SessionState {
	if !stateSession.CreatedAt.IsZero() && tmuxSession.CreatedAt.After(stateSession.CreatedAt.Add(2*time.Second)) {
		// tmux session was recreated (e.g., server restart). Drop stale pane IDs.
		return createSessionState(tmuxSession, output)
	}

	result := &SessionState{
		Path:      stateSession.Path,
		CreatedAt: stateSession.CreatedAt,
		Panes:     make([]*PaneState, 0),
	}

	tmuxPaneMap := make(map[string]domain.Pane)
	for _, p := range tmuxSession.Panes {
		tmuxPaneMap[p.ID] = p
	}

	statePaneMap := make(map[string]*PaneState)
	for _, p := range stateSession.Panes {
		statePaneMap[p.TmuxID] = p
	}

	for _, statePane := range stateSession.Panes {
		tmuxPane, exists := tmuxPaneMap[statePane.TmuxID]
		if !exists {
			output.OrphanedPanes = append(output.OrphanedPanes, statePane.TmuxID)
			continue
		}

		result.Panes = append(result.Panes, statePane)

		if tmuxPane.Title != statePane.Title {
			output.TitleUpdates = append(output.TitleUpdates, TitleUpdate{
				PaneID: statePane.TmuxID,
				Title:  statePane.Title,
			})
		}
	}

	for _, tmuxPane := range tmuxSession.Panes {
		if _, exists := statePaneMap[tmuxPane.ID]; exists {
			continue
		}
		newPane := createPaneState(tmuxPane)
		result.Panes = append(result.Panes, newPane)
		output.NewPanes = append(output.NewPanes, NewPaneInfo{
			SessionName:  tmuxSession.Name,
			PaneID:       tmuxPane.ID,
			InferredType: inferPaneType(tmuxPane),
		})
	}

	return result
}

func createSessionState(tmuxSession domain.Session, output *ReconcileOutput) *SessionState {
	ss := &SessionState{
		Path:      tmuxSession.Path,
		CreatedAt: tmuxSession.CreatedAt,
		Panes:     make([]*PaneState, 0, len(tmuxSession.Panes)),
	}

	for _, pane := range tmuxSession.Panes {
		ss.Panes = append(ss.Panes, createPaneState(pane))
		output.NewPanes = append(output.NewPanes, NewPaneInfo{
			SessionName:  tmuxSession.Name,
			PaneID:       pane.ID,
			InferredType: inferPaneType(pane),
		})
	}

	return ss
}

func createPaneState(pane domain.Pane) *PaneState {
	return &PaneState{
		TmuxID:    pane.ID,
		Type:      string(inferPaneType(pane)),
		Title:     pane.Title,
		CreatedAt: time.Now(),
	}
}

func inferPaneType(pane domain.Pane) domain.PaneType {
	cmd := strings.ToLower(pane.CurrentCommand)
	if strings.Contains(cmd, "codex") {
		return domain.PaneCodex
	}
	if strings.Contains(cmd, "claude") {
		return domain.PaneClaude
	}
	title := strings.ToLower(strings.TrimSpace(pane.Title))
	switch {
	case strings.HasPrefix(title, "codex"):
		return domain.PaneCodex
	case strings.HasPrefix(title, "claude"):
		return domain.PaneClaude
	case strings.HasPrefix(title, "shell"):
		return domain.PaneShell
	}
	return domain.PaneUnknown
}
