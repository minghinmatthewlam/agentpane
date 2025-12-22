package dialogs

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AddPaneResult struct {
	Cancelled bool
	Type      string
}

type AddPaneModel struct {
	options []string
	index   int
}

func NewAddPane() AddPaneModel {
	return AddPaneModel{
		options: []string{"codex", "claude", "shell"},
		index:   0,
	}
}

func (m AddPaneModel) Init() tea.Cmd { return nil }

func (m AddPaneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return AddPaneResult{Cancelled: true} }
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
			return m, nil
		case "down", "j":
			if m.index < len(m.options)-1 {
				m.index++
			}
			return m, nil
		case "enter":
			return m, func() tea.Msg { return AddPaneResult{Type: m.options[m.index]} }
		}
	}
	return m, nil
}

func (m AddPaneModel) View() string {
	var b strings.Builder
	b.WriteString("Add pane:\n\n")
	for i, opt := range m.options {
		cursor := "  "
		if i == m.index {
			cursor = "→ "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, opt))
	}
	b.WriteString("\n[Enter] select  [Esc] cancel")

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(b.String())
}
