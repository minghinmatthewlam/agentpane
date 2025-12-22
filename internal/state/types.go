package state

import "time"

type Store struct {
	Version  int                      `yaml:"version"`
	ServerID string                   `yaml:"server_id,omitempty"`
	Sessions map[string]*SessionState `yaml:"sessions"`
}

type SessionState struct {
	Path      string       `yaml:"path"`
	CreatedAt time.Time    `yaml:"created_at"`
	Panes     []*PaneState `yaml:"panes"`
}

type PaneState struct {
	TmuxID    string     `yaml:"tmux_id"`
	Type      string     `yaml:"type"`
	Title     string     `yaml:"title"`
	CreatedAt time.Time  `yaml:"created_at"`
	RenamedAt *time.Time `yaml:"renamed_at,omitempty"`
}

func NewStore() *Store {
	return &Store{
		Version:  1,
		Sessions: make(map[string]*SessionState),
	}
}
