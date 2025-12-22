package cmd

import (
	"fmt"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/spf13/cobra"
)

func NewTemplatesCmd(a *app.App) *cobra.Command {
	var (
		applyName string
		force     bool
		session   string
	)

	cmd := &cobra.Command{
		Use:   "templates",
		Short: "List or apply templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			if applyName != "" {
				_, err := a.ApplyTemplate(app.ApplyTemplateOptions{
					Session:  session,
					Template: applyName,
					Force:    force,
				})
				return err
			}

			templates, err := a.ListTemplates()
			if err != nil {
				return err
			}
			if len(templates) == 0 {
				fmt.Println("No templates available")
				return nil
			}
			for _, t := range templates {
				desc := ""
				if t.Description != "" {
					desc = " - " + t.Description
				}
				fmt.Printf("%s (%d panes)%s\n", t.Name, len(t.Panes), desc)
			}
			fmt.Printf("\nApply with: agentpane templates --apply <name>%s\n", templateSessionHint(session))
			return nil
		},
	}

	cmd.Flags().StringVar(&applyName, "apply", "", "Apply a template to a session")
	cmd.Flags().BoolVar(&force, "force", false, "Apply template without confirmation")
	cmd.Flags().StringVar(&session, "session", "", "Session name to apply the template to")
	return cmd
}

func templateSessionHint(session string) string {
	if strings.TrimSpace(session) != "" {
		return " --session " + session
	}
	return ""
}
