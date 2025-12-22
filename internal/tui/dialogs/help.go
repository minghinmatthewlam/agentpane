package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpResult struct{}

type HelpModel struct{}

func NewHelp() HelpModel { return HelpModel{} }

func (m HelpModel) Init() tea.Cmd { return nil }

func (m HelpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, func() tea.Msg { return HelpResult{} }
	}
	return m, nil
}

func (m HelpModel) View() string {
	content := `Keys:
  ←/→         Switch session
  ↑/k, ↓/j    Navigate
  Tab         Focus panel
  t           Switch tab
  Enter       Attach / Apply
  c           Add Claude pane
  x           Add Codex pane
  s           Add Shell pane
  a           Add pane (dialog)
  r           Rename pane
  d           Close pane
  /           Filter sessions
  ?           Help
  q           Quit`
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(content)
}
