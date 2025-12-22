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

	SessionTabStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(ColorSecondary)

	ActiveSessionTabStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Bold(true).
				Foreground(lipgloss.Color("229")).
				Background(ColorPrimary)

	CurrentSessionMarker = lipgloss.NewStyle().
				Foreground(ColorSuccess)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Padding(1, 0)

	StatusStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Padding(0, 0)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")).
			Padding(0, 0)
)
