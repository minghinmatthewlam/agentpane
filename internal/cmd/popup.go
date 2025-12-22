package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/tui/dashboard"
	"github.com/spf13/cobra"
)

func NewPopupCmd(a *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "popup",
		Short: "Open dashboard as tmux popup",
		Long:  "Opens the dashboard as a tmux popup overlay. Falls back to regular dashboard if not in tmux.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Getenv("TMUX") == "" {
				fmt.Fprintln(os.Stderr, "Not in tmux, opening dashboard directly")
				return runDashboardTUI(a)
			}

			if err := a.OpenPopup("agentpane dashboard"); err != nil {
				fmt.Fprintln(os.Stderr, "Popup failed, opening dashboard directly:", err)
				return runDashboardTUI(a)
			}
			return nil
		},
	}
}

func runDashboardTUI(a *app.App) error {
	model := dashboard.NewModel(a)
	_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	return err
}
