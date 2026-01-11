package common

import "github.com/minghinmatthewlam/agentpane/internal/domain"

func AgentStatusIndicator(status domain.AgentStatus) string {
	if status == domain.AgentStatusRunning {
		return "●"
	}
	return "○"
}

func SessionStatusIndicator(session domain.Session) string {
	for _, pane := range session.Panes {
		if pane.AgentStatus == domain.AgentStatusRunning {
			return "●"
		}
	}
	return "○"
}
