package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/provider"
	"github.com/minghinmatthewlam/agentpane/internal/state"
	"github.com/minghinmatthewlam/agentpane/internal/tmux"
)

type App struct {
	tmux      *tmux.Client
	providers *provider.Registry
	state     *state.StoreFile
	logger    *log.Logger
}

func New() (*App, error) {
	tmuxClient, err := tmux.NewClient()
	if err != nil {
		return nil, err
	}

	statePath, err := state.DefaultPath()
	if err != nil {
		return nil, err
	}

	return &App{
		tmux:      tmuxClient,
		providers: provider.NewRegistry(),
		state:     state.NewStoreFile(statePath),
		logger:    log.New(os.Stderr, "agentpane: ", log.LstdFlags),
	}, nil
}

func (a *App) InTmux() bool { return a.tmux.InTmux() }

func (a *App) Attach(name string) error {
	if a.tmux.InTmux() {
		return a.tmux.SwitchClient(name)
	}
	return a.tmux.AttachSession(name)
}

func (a *App) SupportsPopup() (bool, error) {
	return a.tmux.SupportsPopup()
}

func (a *App) OpenPopup(command string, args ...string) error {
	return a.tmux.DisplayPopup(command, args...)
}

func (a *App) OpenDashboardWindow() error {
	return a.tmux.OpenWindow("agentpane-dashboard", "agentpane dashboard")
}

func (a *App) EnsureDashboardWindow() error {
	if !a.tmux.InTmux() {
		return fmt.Errorf("must be run inside tmux")
	}
	session, err := a.tmux.CurrentSession()
	if err != nil {
		return err
	}
	const win = "agentpane-dashboard"
	ok, err := a.tmux.HasWindow(session, win)
	if err != nil {
		return err
	}
	if ok {
		if err := a.tmux.SelectWindow(session, win); err != nil {
			// In detached/scripted contexts, tmux may report "no current client".
			if strings.Contains(err.Error(), "no current client") {
				return nil
			}
			return err
		}
		return nil
	}
	if err := a.tmux.NewWindow(session, win, "", "agentpane dashboard"); err != nil {
		return err
	}
	if err := a.tmux.SelectWindow(session, win); err != nil {
		if strings.Contains(err.Error(), "no current client") {
			return nil
		}
		return err
	}
	return nil
}

func (a *App) Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("agentpane: %v", err))
	}
}

func (a *App) CapturePaneContent(paneID string) (string, error) {
	return a.tmux.CapturePaneContent(paneID)
}

func (a *App) KillSession(name string) error {
	return a.tmux.KillSession(name)
}
