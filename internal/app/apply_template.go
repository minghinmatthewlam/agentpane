package app

import (
	"fmt"
	"sort"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/config"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/state"
)

type ApplyTemplateOptions struct {
	Session  string
	Template string
	Force    bool
}

type ApplyTemplateResult struct {
	Session  string
	Template string
	Panes    int
}

func (a *App) ApplyTemplate(opts ApplyTemplateOptions) (ApplyTemplateResult, error) {
	session := opts.Session
	if session == "" {
		if !a.tmux.InTmux() {
			return ApplyTemplateResult{}, fmt.Errorf("session is required when not inside tmux")
		}
		var err error
		session, err = a.tmux.CurrentSession()
		if err != nil {
			return ApplyTemplateResult{}, err
		}
	}

	sessionPath, err := a.tmux.SessionPath(session)
	if err != nil {
		sessionPath = ""
	}

	loaded, err := config.LoadAll(sessionPath)
	if err != nil {
		return ApplyTemplateResult{}, err
	}
	a.applyProviderOverrides(loaded.Merged)

	tmpl, ok := loaded.Merged.Templates[opts.Template]
	if !ok {
		return ApplyTemplateResult{}, fmt.Errorf("unknown template %q", opts.Template)
	}

	panes, err := a.tmux.ListPanes(session)
	if err != nil {
		return ApplyTemplateResult{}, err
	}
	if len(panes) > 0 && !opts.Force {
		return ApplyTemplateResult{}, fmt.Errorf("session has %d panes; use --force to apply template", len(panes))
	}

	if err := a.tmux.SetOption(session, "pane-border-status", "top"); err != nil {
		return ApplyTemplateResult{}, err
	}
	if err := a.tmux.SetOption(session, "pane-border-format", " #{pane_title} "); err != nil {
		return ApplyTemplateResult{}, err
	}

	if len(panes) == 0 {
		return ApplyTemplateResult{}, fmt.Errorf("session has no panes")
	}

	// Keep first pane, kill the rest.
	firstPaneID := panes[0].ID
	for i := 1; i < len(panes); i++ {
		if err := a.tmux.KillPane(panes[i].ID); err != nil {
			a.logger.Printf("failed to kill pane %s: %v", panes[i].ID, err)
		}
	}

	typeCounts := map[domain.PaneType]int{}
	var paneStates []*state.PaneState
	paneIDs := []string{firstPaneID}

	firstRes, err := a.configurePaneSpec(firstPaneID, tmpl.Panes[0], typeCounts)
	if err != nil {
		return ApplyTemplateResult{}, err
	}
	paneStates = append(paneStates, &state.PaneState{
		TmuxID:    firstPaneID,
		Type:      string(firstRes.Type),
		Title:     firstRes.Title,
		CreatedAt: time.Now(),
	})

	for i := 1; i < len(tmpl.Panes); i++ {
		newPaneID, err := a.tmux.SplitPane(session, sessionPath)
		if err != nil {
			return ApplyTemplateResult{}, err
		}
		paneIDs = append(paneIDs, newPaneID)
		res, err := a.configurePaneSpec(newPaneID, tmpl.Panes[i], typeCounts)
		if err != nil {
			return ApplyTemplateResult{}, err
		}
		paneStates = append(paneStates, &state.PaneState{
			TmuxID:    newPaneID,
			Type:      string(res.Type),
			Title:     res.Title,
			CreatedAt: time.Now(),
		})
	}

	layout := "tiled"
	if len(tmpl.Panes) == 2 {
		layout = "even-horizontal"
	}
	_ = a.tmux.SelectLayout(session, layout)

	if err := a.replaceSessionState(session, sessionPath, paneStates); err != nil {
		return ApplyTemplateResult{}, err
	}

	return ApplyTemplateResult{
		Session:  session,
		Template: opts.Template,
		Panes:    len(paneStates),
	}, nil
}

func (a *App) replaceSessionState(session, path string, panes []*state.PaneState) error {
	st := a.loadStateOrNew()
	if err := a.attachServerID(st); err != nil {
		return err
	}

	st.Sessions[session] = &state.SessionState{
		Path:      path,
		CreatedAt: time.Now(),
		Panes:     panes,
	}

	// Ensure deterministic ordering for stability
	sort.Slice(st.Sessions[session].Panes, func(i, j int) bool {
		return st.Sessions[session].Panes[i].TmuxID < st.Sessions[session].Panes[j].TmuxID
	})

	return a.state.Save(st)
}
