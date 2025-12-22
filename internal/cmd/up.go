package cmd

import (
	"fmt"
	"os"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/spf13/cobra"
)

func NewUpCmd(a *app.App) *cobra.Command {
	var (
		sessionName string
		template    string
		detach      bool
	)

	cmd := &cobra.Command{
		Use:   "up",
		Short: "Create or attach to session for current repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			result, err := a.Up(app.UpOptions{
				Cwd:          cwd,
				ExplicitName: sessionName,
				Template:     template,
				Detach:       detach,
			})
			if err != nil {
				return err
			}

			for _, w := range result.Warnings {
				fmt.Fprintln(os.Stderr, "Warning:", w)
			}

			switch result.Action {
			case app.ActionAlreadyIn:
				fmt.Printf("Already in session %s\n", result.SessionName)
			case app.ActionDetached:
				fmt.Printf("Session %s ready\n", result.SessionName)
			case app.ActionCreated:
				// Outside tmux, attach-session will take over the terminal; keep output minimal.
			case app.ActionAttached:
				// Same as created.
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&sessionName, "name", "n", "", "Explicit session name")
	cmd.Flags().StringVarP(&template, "template", "t", "", "Use specific template")
	cmd.Flags().BoolVarP(&detach, "detach", "d", false, "Create but don't attach")

	return cmd
}
