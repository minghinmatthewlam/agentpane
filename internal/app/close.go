package app

import (
	"errors"
	"os"
)

func (a *App) ClosePane(paneID string) error {
	if err := a.tmux.KillPane(paneID); err != nil {
		return err
	}

	st, err := a.state.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return nil
	}

	for _, session := range st.Sessions {
		filtered := session.Panes[:0]
		for _, p := range session.Panes {
			if p.TmuxID != paneID {
				filtered = append(filtered, p)
			}
		}
		session.Panes = filtered
	}

	return a.state.Save(st)
}
