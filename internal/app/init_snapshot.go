package app

import (
	"fmt"

	"github.com/minghinmatthewlam/agentpane/internal/config"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

func (a *App) SnapshotCurrentLayout() (config.Layout, string, error) {
	session, err := a.tmux.CurrentSession()
	if err != nil {
		return config.Layout{}, "", err
	}

	panes, err := a.tmux.ListPanes(session)
	if err != nil {
		return config.Layout{}, "", err
	}
	if len(panes) == 0 {
		return config.Layout{}, "", fmt.Errorf("no panes found in session %q", session)
	}

	specs := make([]config.PaneSpec, 0, len(panes))
	for _, p := range panes {
		specs = append(specs, config.PaneSpec{
			Type:  string(domain.InferPaneType("", p.Title)),
			Title: p.Title,
		})
	}

	return config.Layout{Panes: specs}, session, nil
}
