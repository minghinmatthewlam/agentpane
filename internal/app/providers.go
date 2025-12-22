package app

import (
	"github.com/minghinmatthewlam/agentpane/internal/config"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/provider"
)

func (a *App) applyProviderOverrides(cfg *config.Config) {
	a.providers = provider.NewRegistry()
	if cfg == nil {
		return
	}
	for k, v := range cfg.Providers {
		if v.Command == "" {
			continue
		}
		switch k {
		case "codex":
			a.providers.Override(domain.PaneCodex, v.Command)
		case "claude":
			a.providers.Override(domain.PaneClaude, v.Command)
		case "shell":
			a.providers.Override(domain.PaneShell, v.Command)
		}
	}
}
