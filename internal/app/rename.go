package app

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/state"
)

type RenameOptions struct {
	Title string
}

type RenameResult struct {
	OldTitle string
	NewTitle string
}

func (a *App) Rename(opts RenameOptions) (RenameResult, error) {
	if !a.tmux.InTmux() {
		return RenameResult{}, fmt.Errorf("must be run inside tmux")
	}

	session, err := a.tmux.CurrentSession()
	if err != nil {
		return RenameResult{}, err
	}

	paneID, err := a.tmux.CurrentPane()
	if err != nil {
		return RenameResult{}, err
	}

	currentTitle, err := a.lookupPaneTitle(session, paneID)
	if err != nil {
		return RenameResult{}, err
	}

	newTitle := strings.TrimSpace(opts.Title)
	if newTitle == "" {
		if !stdinIsTTY() {
			return RenameResult{}, fmt.Errorf("rename requires a title when stdin is not a TTY")
		}
		newTitle, err = promptForTitle(currentTitle)
		if err != nil {
			return RenameResult{}, err
		}
	}
	if newTitle == "" {
		return RenameResult{}, fmt.Errorf("title cannot be empty")
	}

	if err := a.tmux.SetPaneTitle(paneID, newTitle); err != nil {
		return RenameResult{}, err
	}

	if err := a.updateStateForRename(session, paneID, newTitle); err != nil {
		return RenameResult{}, err
	}

	return RenameResult{OldTitle: currentTitle, NewTitle: newTitle}, nil
}

func (a *App) lookupPaneTitle(session, paneID string) (string, error) {
	panes, err := a.tmux.ListPanes(session)
	if err != nil {
		return "", err
	}
	for _, p := range panes {
		if p.ID == paneID {
			return p.Title, nil
		}
	}
	return "", fmt.Errorf("pane %s not found", paneID)
}

func promptForTitle(current string) (string, error) {
	fmt.Fprintf(os.Stderr, "Rename pane (current: %s): ", current)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	title := strings.TrimSpace(line)
	if title == "" {
		return current, nil
	}
	return title, nil
}

func (a *App) updateStateForRename(session, paneID, title string) error {
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

	now := time.Now()
	for _, p := range ss.Panes {
		if p.TmuxID == paneID {
			p.Title = title
			p.RenamedAt = &now
			return a.state.Save(st)
		}
	}

	ss.Panes = append(ss.Panes, &state.PaneState{
		TmuxID:    paneID,
		Type:      "unknown",
		Title:     title,
		CreatedAt: now,
		RenamedAt: &now,
	})

	return a.state.Save(st)
}
