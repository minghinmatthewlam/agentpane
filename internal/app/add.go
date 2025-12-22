package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/config"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/provider"
	"github.com/minghinmatthewlam/agentpane/internal/state"
)

type AddOptions struct {
	Type          domain.PaneType
	ExplicitTitle string
}

type AddResult struct {
	Title           string
	FellBackToShell bool
	PaneID          string
}

func (a *App) Add(opts AddOptions) (AddResult, error) {
	if !a.tmux.InTmux() {
		return AddResult{}, fmt.Errorf("must be run inside tmux")
	}

	session, err := a.tmux.CurrentSession()
	if err != nil {
		return AddResult{}, err
	}

	cwd, err := a.tmux.SessionPath(session)
	if err != nil {
		cwd = ""
	}

	if err := a.applyConfigOverrides(cwd); err != nil {
		return AddResult{}, err
	}

	prov, actualType, ok := a.providers.GetWithFallback(opts.Type)
	if !ok {
		return AddResult{}, fmt.Errorf("unknown pane type: %s", opts.Type)
	}

	paneID, err := a.tmux.SplitPane(session, cwd)
	if err != nil {
		return AddResult{}, err
	}

	title := strings.TrimSpace(opts.ExplicitTitle)
	if title == "" {
		title, err = a.nextAutoTitle(session, actualType, prov)
		if err != nil {
			return AddResult{}, err
		}
	}

	if err := a.tmux.SetPaneTitle(paneID, title); err != nil {
		return AddResult{}, err
	}
	if err := a.launchProvider(paneID, prov); err != nil {
		return AddResult{}, err
	}

	if err := a.updateStateForNewPane(session, paneID, actualType, title); err != nil {
		return AddResult{}, err
	}

	return AddResult{
		Title:           title,
		FellBackToShell: actualType != opts.Type,
		PaneID:          paneID,
	}, nil
}

func (a *App) applyConfigOverrides(cwd string) error {
	if cwd == "" {
		cwd, _ = os.Getwd()
	}
	loaded, err := config.LoadAll(cwd)
	if err != nil {
		return err
	}
	a.applyProviderOverrides(loaded.Merged)
	return nil
}

func (a *App) nextAutoTitle(session string, t domain.PaneType, prov *provider.Provider) (string, error) {
	count := 0

	st, err := a.state.Load()
	if err == nil {
		if ss, ok := st.Sessions[session]; ok {
			for _, p := range ss.Panes {
				if p.Type == string(t) {
					count++
				}
			}
		}
	}

	if count == 0 {
		// fallback to tmux titles if state is missing
		if panes, err := a.tmux.ListPanes(session); err == nil {
			prefix := strings.ToLower(prov.TitlePrefix)
			for _, p := range panes {
				if strings.HasPrefix(strings.ToLower(p.Title), prefix) {
					count++
				}
			}
		}
	}

	return fmt.Sprintf("%s%d", prov.TitlePrefix, count+1), nil
}

func (a *App) updateStateForNewPane(session, paneID string, t domain.PaneType, title string) error {
	st, err := a.state.Load()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			st = state.NewStore()
		} else {
			st = state.NewStore()
		}
	}

	ss, ok := st.Sessions[session]
	if !ok {
		path, _ := a.tmux.SessionPath(session)
		ss = &state.SessionState{
			Path:      path,
			CreatedAt: time.Now(),
			Panes:     []*state.PaneState{},
		}
		st.Sessions[session] = ss
	}

	ss.Panes = append(ss.Panes, &state.PaneState{
		TmuxID:    paneID,
		Type:      string(t),
		Title:     title,
		CreatedAt: time.Now(),
	})

	return a.state.Save(st)
}
