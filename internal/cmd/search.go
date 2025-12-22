package cmd

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/tui/search"
	"github.com/spf13/cobra"
)

func NewSearchCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search across sessions and panes",
		RunE: func(cmd *cobra.Command, args []string) error {
			model := search.NewModel(a)
			_, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
			return err
		},
	}
	return cmd
}
