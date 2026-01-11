package app

import (
	"errors"
	"os"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/state"
)

func (a *App) loadStateOrNew() *state.Store {
	st, err := a.state.Load()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			a.logger.Printf("failed to load state: %v", err)
		}
		return state.NewStore()
	}
	return st
}

func (a *App) attachServerID(st *state.Store) error {
	if st == nil {
		return nil
	}
	serverID, err := a.ensureServerID()
	if err != nil {
		return err
	}
	st.ServerID = serverID
	return nil
}

func (a *App) ensureSessionState(st *state.Store, session string) *state.SessionState {
	if st == nil {
		return nil
	}
	if ss, ok := st.Sessions[session]; ok {
		return ss
	}
	path, _ := a.tmux.SessionPath(session)
	ss := &state.SessionState{
		Path:      path,
		CreatedAt: time.Now(),
		Panes:     []*state.PaneState{},
	}
	st.Sessions[session] = ss
	return ss
}
