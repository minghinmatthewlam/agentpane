package dialogs

import "github.com/charmbracelet/lipgloss"

func dialogStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
}
