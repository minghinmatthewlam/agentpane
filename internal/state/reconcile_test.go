package state

import (
	"testing"
	"time"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

func TestReconcileUpdatesTitlesAndAddsNew(t *testing.T) {
	stateStore := &Store{
		Version: 1,
		Sessions: map[string]*SessionState{
			"repo": {
				Path:      "/tmp/repo",
				CreatedAt: time.Now().Add(-time.Hour),
				Panes: []*PaneState{
					{TmuxID: "%0", Type: "codex", Title: "old-title", CreatedAt: time.Now()},
					{TmuxID: "%9", Type: "shell", Title: "orphan", CreatedAt: time.Now()},
				},
			},
		},
	}

	tmuxSessions := []domain.Session{
		{
			Name: "repo",
			Panes: []domain.Pane{
				{ID: "%0", Title: "new-title", CurrentCommand: "zsh"},
				{ID: "%1", Title: "shell-1", CurrentCommand: "zsh"},
			},
		},
	}

	output := Reconcile(ReconcileInput{
		CurrentState: stateStore,
		TmuxSessions: tmuxSessions,
	})

	if len(output.TitleUpdates) != 1 {
		t.Fatalf("expected 1 title update, got %d", len(output.TitleUpdates))
	}
	if output.TitleUpdates[0].PaneID != "%0" || output.TitleUpdates[0].Title != "old-title" {
		t.Fatalf("unexpected title update: %#v", output.TitleUpdates[0])
	}
	if len(output.OrphanedPanes) != 1 || output.OrphanedPanes[0] != "%9" {
		t.Fatalf("expected orphaned pane %%9, got %#v", output.OrphanedPanes)
	}
	if len(output.NewPanes) != 1 || output.NewPanes[0].PaneID != "%1" {
		t.Fatalf("expected new pane %%1, got %#v", output.NewPanes)
	}
	if len(output.UpdatedState.Sessions["repo"].Panes) != 2 {
		t.Fatalf("expected 2 panes in updated state")
	}
}
