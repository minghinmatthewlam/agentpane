package dialogs

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmResult struct {
	Accepted bool
}

type ConfirmModel struct {
	Title string
	Body  string
}

func NewConfirm(title, body string) ConfirmModel {
	return ConfirmModel{Title: title, Body: body}
}

func (m ConfirmModel) Init() tea.Cmd { return nil }

func (m ConfirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			return m, func() tea.Msg { return ConfirmResult{Accepted: true} }
		case "n", "esc", "q":
			return m, func() tea.Msg { return ConfirmResult{Accepted: false} }
		}
	}
	return m, nil
}

func (m ConfirmModel) View() string {
	content := fmt.Sprintf("%s\n\n%s\n\n[y] yes  [n] no", m.Title, m.Body)
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(content)
}
