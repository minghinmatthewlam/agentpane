package dialogs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RenameResult struct {
	Cancelled bool
	Title     string
}

type RenameModel struct {
	input textinput.Model
}

func NewRename(initial string) RenameModel {
	ti := textinput.New()
	ti.Placeholder = "New title"
	ti.SetValue(initial)
	ti.Focus()
	return RenameModel{input: ti}
}

func (m RenameModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RenameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return RenameResult{Cancelled: true} }
		case "enter":
			return m, func() tea.Msg { return RenameResult{Title: m.input.Value()} }
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m RenameModel) View() string {
	content := "Rename pane:\n\n" + m.input.View() + "\n\n[Enter] save  [Esc] cancel"
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(content)
}
