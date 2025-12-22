package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewHelpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "help",
		Short: "Show help and recommended tmux configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(`agentpane - AI Coding Agent Manager

COMMANDS:
  up              Create or attach to session for current repo
  add <type>      Add pane (codex, claude, shell)
  rename [name]   Rename current pane
  dashboard       Open navigation TUI
  popup           Open dashboard as tmux popup
  templates       Browse and apply templates
  search          Search across sessions/panes
  init            Generate .agentpane.yml

QUICK ACCESS (add to ~/.tmux.conf):

    bind-key g run-shell "agentpane popup"

Then press Prefix+g from any pane to open the dashboard popup.
Reload config with: tmux source-file ~/.tmux.conf

DASHBOARD KEYS:
  ←/→, h/l    Switch between sessions (repos)
  ↑/↓, j/k    Navigate panes
  Tab         Switch panels
  Enter       Attach to session
  a           Add pane
  r           Rename pane
  x           Close pane
  t           Switch tabs (Sessions/Templates)
  q           Quit
`)
		},
	}
}
