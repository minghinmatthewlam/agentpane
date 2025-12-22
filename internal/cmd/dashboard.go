package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/tui/dashboard"
	"github.com/spf13/cobra"
)

func NewDashboardCmd(a *app.App) *cobra.Command {
	var tmuxWindow bool
	cmd := &cobra.Command{
		Use:   "dashboard",
		Short: "Open the interactive dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			if tmuxWindow {
				if os.Getenv("TMUX") == "" {
					fmt.Fprintln(os.Stderr, "Not in tmux, opening dashboard directly")
					return runDashboard(a)
				}
				return a.EnsureDashboardWindow()
			}
			return runDashboard(a)
		},
	}
	cmd.Flags().BoolVar(&tmuxWindow, "tmux-window", false, "Open dashboard in a tmux window (recommended over popup)")
	return cmd
}

func runDashboard(a *app.App) error {
	model := dashboard.NewModel(a)
	finalModel, err := tea.NewProgram(model, tea.WithAltScreen()).Run()
	if err != nil {
		return err
	}

	m, ok := finalModel.(dashboard.Model)
	if !ok {
		return nil
	}

	// Check if user requested to attach to a session
	if sessionName := m.AttachSession(); sessionName != "" {
		return a.Attach(sessionName)
	}

	// Check if user requested to open a new session
	if path := m.OpenSessionPath(); path != "" {
		result, err := a.Up(app.UpOptions{Cwd: path})
		if err != nil {
			return fmt.Errorf("failed to open session: %w", err)
		}
		fmt.Printf("Session '%s' %s\n", result.SessionName, result.Action)
		return nil
	}

	// Check if user requested to add a pane
	if paneType := m.AddPaneType(); paneType != "" {
		result, err := a.Add(app.AddOptions{Type: paneType})
		if err != nil {
			return fmt.Errorf("failed to add pane: %w", err)
		}
		if result.FellBackToShell {
			fmt.Printf("Warning: %s not found in PATH, created shell pane instead\n", paneType)
		}
		fmt.Printf("Created pane '%s'\n", result.Title)
	}

	return nil
}
