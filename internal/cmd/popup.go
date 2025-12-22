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
		Long:  "Opens the dashboard as a tmux popup overlay (falls back to a tmux window if popups are unsupported).",
		RunE: func(cmd *cobra.Command, args []string) error {
			if os.Getenv("TMUX") == "" {
				fmt.Fprintln(os.Stderr, "Not in tmux, opening dashboard directly")
				return runDashboard(a)
			}

			supported, err := a.SupportsPopup()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to detect tmux popup support, opening dashboard directly:", err)
				return runDashboard(a)
			}

			if supported {
				if err := a.OpenPopup("agentpane", "dashboard"); err == nil {
					return nil
				}
				// Popup failed; fall back to a tmux window.
				fmt.Fprintln(os.Stderr, "Popup failed, opening dashboard in a tmux window")
			}

			if err := a.OpenDashboardWindow(); err != nil {
				fmt.Fprintln(os.Stderr, "Failed to open dashboard window, opening dashboard directly:", err)
				return runDashboard(a)
			}
			return nil
		},
	}
}
