package domain

import "time"

type PaneType string

const (
	PaneCodex   PaneType = "codex"
	PaneClaude  PaneType = "claude"
	PaneShell   PaneType = "shell"
	PaneUnknown PaneType = "unknown"
)

type PaneStatus string

const (
	StatusActive  PaneStatus = "active"
	StatusExited  PaneStatus = "exited"
	StatusUnknown PaneStatus = "unknown"
)

type AgentStatus string

const (
	AgentStatusIdle    AgentStatus = "idle"
	AgentStatusRunning AgentStatus = "running"
)

type Pane struct {
	ID             string
	Index          int
	Title          string
	Type           PaneType
	Status         PaneStatus
	AgentStatus    AgentStatus
	PID            int
	CurrentCommand string
	CurrentPath    string
	LastActive     time.Time
}

type Session struct {
	Name      string
	Path      string
	CreatedAt time.Time
	Attached  bool
	Panes     []Pane
}

type Snapshot struct {
	Sessions       []Session
	CurrentSession string
	CurrentPane    string
}
