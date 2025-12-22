package app

import (
	"fmt"
	"log"
	"os"

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

func (a *App) Must(err error) {
	if err != nil {
		panic(fmt.Sprintf("agentpane: %v", err))
	}
}
