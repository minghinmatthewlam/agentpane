package cmd

import (
	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/spf13/cobra"
)

func NewRootCmd(a *app.App) *cobra.Command {
	root := &cobra.Command{
		Use:   "agentpane",
		Short: "tmux-based AI coding agent manager",
		// Default to dashboard when no subcommand given
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDashboard(a)
		},
	}

	root.AddCommand(NewUpCmd(a))
	root.AddCommand(NewInitCmd(a))
	root.AddCommand(NewAddCmd(a))
	root.AddCommand(NewRenameCmd(a))
	root.AddCommand(NewDashboardCmd(a))
	root.AddCommand(NewPopupCmd(a))
	root.AddCommand(NewTemplatesCmd(a))
	root.AddCommand(NewHelpCmd())
	root.AddCommand(NewSearchCmd(a))
	root.AddCommand(NewOverlayCmd(a))
	root.AddCommand(NewOverlayViewCmd(a))
	return root
}
