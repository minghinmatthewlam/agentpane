package app

func (a *App) ClosePane(paneID string) error {
	if err := a.tmux.KillPane(paneID); err != nil {
		return err
	}

	st := a.loadStateOrNew()

	for _, session := range st.Sessions {
		filtered := session.Panes[:0]
		for _, p := range session.Panes {
			if p.TmuxID != paneID {
				filtered = append(filtered, p)
			}
		}
		session.Panes = filtered
	}

	if err := a.attachServerID(st); err != nil {
		return err
	}
	return a.state.Save(st)
}
