package cmd

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/tui/dashboard"
	"github.com/spf13/cobra"
)

func NewDashboardCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Open the interactive dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := dashboard.NewModel(a)
			_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
			return err
		},
	}
	return cmd
}
