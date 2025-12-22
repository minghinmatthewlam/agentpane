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
  ↑/k, ↓/j    Navigate
  Tab         Focus panel
  t           Switch tab
  Enter       Attach / Apply
  a           Add pane
  r           Rename pane
  x           Close pane
  ?           Help
  q           Quit`
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(content)
}
