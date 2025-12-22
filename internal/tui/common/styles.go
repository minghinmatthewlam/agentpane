package common

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("33")  // blue
	ColorSecondary = lipgloss.Color("245") // gray
	ColorSuccess   = lipgloss.Color("70")  // green
	ColorWarning   = lipgloss.Color("214")

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorSecondary).
			Padding(1)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	NormalStyle = lipgloss.NewStyle()

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))

	DimSelectedStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary)

	TabStyle = lipgloss.NewStyle().
			Padding(0, 2)

	ActiveTabStyle = TabStyle.Copy().
			Bold(true).
			Foreground(ColorPrimary)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Padding(1, 0)
)
