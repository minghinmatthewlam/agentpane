package cmd

import (
	"fmt"
	"os"

	"github.com/minghinmatthewlam/agentpane/internal/app"
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
				return runDashboard(a)
			}

			if err := a.OpenPopup("agentpane dashboard"); err != nil {
				fmt.Fprintln(os.Stderr, "Popup failed, opening dashboard directly:", err)
				return runDashboard(a)
			}
			return nil
		},
	}
}
