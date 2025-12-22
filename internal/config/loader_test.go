package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllMergesGlobalAndRepo(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	globalPath, err := GlobalConfigPath()
	if err != nil {
		t.Fatalf("global path: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(globalPath), 0o755); err != nil {
		t.Fatalf("mkdir global: %v", err)
	}
	if err := os.WriteFile(globalPath, []byte("default_template: quad\n"), 0o644); err != nil {
		t.Fatalf("write global: %v", err)
	}

	repo := filepath.Join(tmp, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}
	repoCfg := []byte("layout:\n  panes:\n    - type: codex\n    - type: claude\n")
	if err := os.WriteFile(filepath.Join(repo, ".agentpane.yml"), repoCfg, 0o644); err != nil {
		t.Fatalf("write repo: %v", err)
	}

	loaded, err := LoadAll(repo)
	if err != nil {
		t.Fatalf("LoadAll: %v", err)
	}
	if loaded.Repo == nil {
		t.Fatalf("expected repo config")
	}
	if loaded.Merged.DefaultTemplate != "quad" {
		t.Fatalf("expected default_template quad, got %s", loaded.Merged.DefaultTemplate)
	}
	if _, ok := loaded.Merged.Templates["duo"]; !ok {
		t.Fatalf("expected builtin template duo")
	}
	if len(loaded.Repo.Layout.Panes) != 2 {
		t.Fatalf("expected 2 panes in repo layout")
	}
}
