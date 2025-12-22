package domain

import "testing"

func TestParsePaneType(t *testing.T) {
	got, err := ParsePaneType("codex")
	if err != nil || got != PaneCodex {
		t.Fatalf("expected codex, got %v err=%v", got, err)
	}
	if _, err := ParsePaneType("bad"); err == nil {
		t.Fatalf("expected error for invalid type")
	}
}

func TestInferPaneType(t *testing.T) {
	if got := InferPaneType("codex", ""); got != PaneCodex {
		t.Fatalf("expected codex, got %s", got)
	}
	if got := InferPaneType("zsh", "claude-1"); got != PaneClaude {
		t.Fatalf("expected claude, got %s", got)
	}
	if got := InferPaneType("zsh", "random"); got != PaneUnknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}
