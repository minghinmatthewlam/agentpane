package cmd

import (
	"fmt"

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
			return runDashboard(a)
		},
	}
	return cmd
}

func runDashboard(a *app.App) error {
	model := dashboard.NewModel(a)
	finalModel, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}

	// Check if user requested to add a pane
	if m, ok := finalModel.(dashboard.Model); ok {
		if paneType := m.AddPaneType(); paneType != "" {
			result, err := a.Add(app.AddOptions{Type: paneType})
			if err != nil {
				return fmt.Errorf("failed to add pane: %w", err)
			}
			if result.FellBackToShell {
				fmt.Printf("Warning: %s not found in PATH, created shell pane instead\n", paneType)
			}
			fmt.Printf("Created pane '%s'\n", result.Title)
		}
	}

	return nil
}
