package tmux

import "testing"

func TestParsePanesEscapedDelimiter(t *testing.T) {
	output := "%0\\0370\\037codex-1\\037zsh\\037/Users/test\\037123\n"
	panes, err := ParsePanes(output)
	if err != nil {
		t.Fatalf("ParsePanes error: %v", err)
	}
	if len(panes) != 1 {
		t.Fatalf("expected 1 pane, got %d", len(panes))
	}
	if panes[0].ID != "%0" || panes[0].Title != "codex-1" {
		t.Fatalf("unexpected pane: %#v", panes[0])
	}
}

func TestParseSessionsEscapedDelimiter(t *testing.T) {
	output := "repo\\037/Users/test\\0371700000000\\0371\n"
	sessions, err := ParseSessions(output)
	if err != nil {
		t.Fatalf("ParseSessions error: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(sessions))
	}
	if sessions[0].Name != "repo" || sessions[0].Path != "/Users/test" {
		t.Fatalf("unexpected session: %#v", sessions[0])
	}
}
