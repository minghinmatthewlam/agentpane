package app

import (
	"fmt"

	"github.com/minghinmatthewlam/agentpane/internal/provider"
	"github.com/minghinmatthewlam/agentpane/internal/tmux"
)

type App struct {
	tmux      *tmux.Client
	providers *provider.Registry
}

func New() (*App, error) {
	tmuxClient, err := tmux.NewClient()
	if err != nil {
		return nil, err
	}

	return &App{
		tmux:      tmuxClient,
		providers: provider.NewRegistry(),
	}, nil
}

func (a *App) InTmux() bool { return a.tmux.InTmux() }

func (a *App) Attach(name string) error {
	if a.tmux.InTmux() {
		return a.tmux.SwitchClient(name)
	}
	return a.tmux.AttachSession(name)
}

func (a *App) Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("agentpane: %v", err))
	}
}

