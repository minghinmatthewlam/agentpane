package cmd

import (
	"fmt"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/spf13/cobra"
)

func NewAddCmd(a *app.App) *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "add <codex|claude|shell>",
		Short: "Add a pane to the current tmux session",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !a.InTmux() {
				return fmt.Errorf("must be run inside a tmux session")
			}

			paneType, err := parsePaneType(args[0])
			if err != nil {
				return err
			}

			result, err := a.Add(app.AddOptions{
				Type:          paneType,
				ExplicitTitle: title,
			})
			if err != nil {
				return err
			}

			if result.FellBackToShell {
				fmt.Printf("Warning: %s not found in PATH, created shell pane instead\n", args[0])
			}
			fmt.Printf("Created pane '%s'\n", result.Title)
			return nil
		},
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Custom pane title")
	return cmd
}

func parsePaneType(s string) (domain.PaneType, error) {
	switch s {
	case "codex":
		return domain.PaneCodex, nil
	case "claude":
		return domain.PaneClaude, nil
	case "shell":
		return domain.PaneShell, nil
	default:
		return "", fmt.Errorf("unknown pane type: %s (use codex, claude, or shell)", s)
	}
}
