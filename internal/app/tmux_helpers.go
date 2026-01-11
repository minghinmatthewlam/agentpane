package app

func (a *App) enablePaneTitles(session string) error {
	if err := a.tmux.SetOption(session, "pane-border-status", "top"); err != nil {
		return err
	}
	if err := a.tmux.SetOption(session, "pane-border-format", " #{pane_title} "); err != nil {
		return err
	}
	return nil
}
