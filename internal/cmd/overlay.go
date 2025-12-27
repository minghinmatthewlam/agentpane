package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/tui/overlay"
)

func NewOverlayCmd(a *app.App) *cobra.Command {
	var position string
	var size int

	cmd := &cobra.Command{
		Use:   "overlay",
		Short: "Toggle the session status overlay",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.ToggleOverlay(app.OverlayOptions{
				Position: position,
				Size:     size,
			})
		},
	}

	cmd.Flags().StringVar(&position, "position", "left", "Overlay position: left, right, top, bottom")
	cmd.Flags().IntVar(&size, "size", 30, "Overlay size in columns/rows")

	return cmd
}

func NewOverlayViewCmd(a *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "overlay-view",
		Short:  "Internal: overlay view",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runOverlayView(a)
		},
	}
	return cmd
}

func runOverlayView(a *app.App) error {
	model := overlay.NewModel(a)
	_, err := tea.NewProgram(model, tea.WithMouseCellMotion()).Run()
	return err
}
