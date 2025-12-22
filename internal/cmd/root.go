package cmd

import (
	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/spf13/cobra"
)

func NewRootCmd(a *app.App) *cobra.Command {
	root := &cobra.Command{
		Use:   "agentpane",
		Short: "tmux-based AI coding agent manager",
	}

	root.AddCommand(NewUpCmd(a))
	root.AddCommand(NewInitCmd(a))
	root.AddCommand(NewAddCmd(a))
	root.AddCommand(NewRenameCmd(a))
	return root
}
