package provider

import (
	"testing"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

func TestGetWithFallbackUsesShellWhenUnavailable(t *testing.T) {
	r := NewRegistry()
	r.providers[domain.PaneCodex].Executable = "definitely-not-a-binary"

	p, actual, ok := r.GetWithFallback(domain.PaneCodex)
	if !ok {
		t.Fatalf("expected provider, got none")
	}
	if actual != domain.PaneShell {
		t.Fatalf("expected fallback to shell, got %s", actual)
	}
	if p.Type != domain.PaneShell {
		t.Fatalf("expected shell provider, got %s", p.Type)
	}
}

func TestOverrideCommand(t *testing.T) {
	r := NewRegistry()
	r.Override(domain.PaneCodex, "codex --foo")

	p, ok := r.Get(domain.PaneCodex)
	if !ok {
		t.Fatalf("expected codex provider")
	}
	if p.Command != "codex --foo" {
		t.Fatalf("expected override command, got %s", p.Command)
	}
}
