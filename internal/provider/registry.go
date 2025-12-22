package provider

import (
	"os"
	"os/exec"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

type Provider struct {
	Type        domain.PaneType
	Command     string
	TitlePrefix string
	Executable  string
}

type Registry struct {
	providers map[domain.PaneType]*Provider
	overrides map[domain.PaneType]string
}

func NewRegistry() *Registry {
	return &Registry{
		providers: map[domain.PaneType]*Provider{
			domain.PaneCodex: {
				Type:        domain.PaneCodex,
				Command:     "codex",
				TitlePrefix: "codex-",
				Executable:  "codex",
			},
			domain.PaneClaude: {
				Type:        domain.PaneClaude,
				Command:     "claude",
				TitlePrefix: "claude-",
				Executable:  "claude",
			},
			domain.PaneShell: {
				Type:        domain.PaneShell,
				Command:     "",
				TitlePrefix: "shell-",
				Executable:  "",
			},
		},
		overrides: make(map[domain.PaneType]string),
	}
}

func (r *Registry) Get(t domain.PaneType) (*Provider, bool) {
	p, ok := r.providers[t]
	if !ok {
		return nil, false
	}
	if override, has := r.overrides[t]; has {
		copied := *p
		copied.Command = override
		return &copied, true
	}
	return p, true
}

func (r *Registry) IsAvailable(t domain.PaneType) bool {
	p, ok := r.providers[t]
	if !ok {
		return false
	}
	if p.Executable == "" {
		return true
	}
	_, err := exec.LookPath(p.Executable)
	return err == nil
}

func (r *Registry) GetWithFallback(t domain.PaneType) (*Provider, domain.PaneType, bool) {
	if r.IsAvailable(t) {
		p, ok := r.Get(t)
		return p, t, ok
	}
	p, ok := r.Get(domain.PaneShell)
	return p, domain.PaneShell, ok
}

func (r *Registry) Override(t domain.PaneType, command string) {
	r.overrides[t] = command
}

func DefaultShell() string {
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh
	}
	return "sh"
}

