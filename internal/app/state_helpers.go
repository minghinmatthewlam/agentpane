package app

import (
	"errors"
	"os"

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
