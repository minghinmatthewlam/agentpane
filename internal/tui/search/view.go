package search

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("Search\n\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")

	if m.errMsg != "" {
		b.WriteString("Error: " + m.errMsg + "\n")
	} else if len(m.results) == 0 {
		b.WriteString("No results\n")
	} else {
		for _, r := range m.results {
			if r.PaneID == "" {
				b.WriteString(fmt.Sprintf("Session: %s\n", r.Session))
			} else {
				b.WriteString(fmt.Sprintf("Pane: %s (%s) in %s\n", r.Title, r.Type, r.Session))
			}
		}
	}

	b.WriteString("\n[Esc] close")
	style := lipgloss.NewStyle().
		Padding(1, 2)
	return style.Render(b.String())
}
