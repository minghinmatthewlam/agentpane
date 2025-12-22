package cmd

import (
	"fmt"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/spf13/cobra"
)

func NewRenameCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rename [name]",
		Short: "Rename the current pane",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !a.InTmux() {
				return fmt.Errorf("must be run inside a tmux session")
			}

			title := ""
			if len(args) == 1 {
				title = args[0]
			}

			_, err := a.Rename(app.RenameOptions{Title: title})
			return err
		},
	}

	return cmd
}
