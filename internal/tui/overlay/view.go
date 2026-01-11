package overlay

import (
	"fmt"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/tui/common"
)

func (m Model) View() string {
	if len(m.items) == 0 {
		if m.errorMsg != "" {
			return common.ErrorStyle.Render("Error: " + m.errorMsg)
		}
		return common.DimSelectedStyle.Render("No sessions")
	}

	var b strings.Builder
	for i, it := range m.items {
		cursor := "  "
		style := common.NormalStyle
		if i == m.index {
			cursor = "> "
			style = common.SelectedStyle
		}
		line := renderItem(it)
		if m.width > 0 {
			available := m.width - len(cursor)
			line = truncate(line, available)
		}
		b.WriteString(style.Render(cursor + line))
		b.WriteString("\n")
	}
	if m.errorMsg != "" {
		b.WriteString(common.ErrorStyle.Render("Error: " + m.errorMsg))
	}
	return b.String()
}

func renderItem(it item) string {
	if it.kind == itemSession {
		return fmt.Sprintf("%s %s", it.indicator, it.session)
	}
	label := it.pane.Title
	typeBadge := fmt.Sprintf("[%s]", it.pane.Type)
	return fmt.Sprintf("  %s %s %s", it.indicator, label, typeBadge)
}

func buildItems(snapshot domain.Snapshot) []item {
	var items []item
	for i := range snapshot.Sessions {
		session := &snapshot.Sessions[i]
		var visiblePanes []domain.Pane
		for j := range session.Panes {
			pane := &session.Panes[j]
			if strings.TrimSpace(pane.Title) == "agentpane-status" {
				continue
			}
			visiblePanes = append(visiblePanes, *pane)
		}
		visibleSession := domain.Session{
			Name:  session.Name,
			Panes: visiblePanes,
		}
		sessionIndicator := common.SessionStatusIndicator(visibleSession)
		sessionItem := item{
			kind:      itemSession,
			session:   session.Name,
			indicator: sessionIndicator,
		}
		items = append(items, sessionItem)
		for j := range visiblePanes {
			pane := &visiblePanes[j]
			items = append(items, item{
				kind:      itemPane,
				session:   session.Name,
				pane:      pane,
				indicator: common.AgentStatusIndicator(pane.AgentStatus),
			})
		}
	}
	return items
}

func truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}

