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
  templates       Browse and apply templates
  search          Search across sessions/panes
  init            Generate .agentpane.yml

RECOMMENDED TMUX CONFIGURATION:

Add this to your ~/.tmux.conf to open dashboard with Prefix+g:

    bind-key g run-shell "agentpane dashboard"

Then reload tmux config:

    tmux source-file ~/.tmux.conf

DASHBOARD KEYS:
  ↑/k, ↓/j    Navigate
  Tab         Switch panels
  t           Switch tabs
  Enter       Attach / Apply
  a           Add pane
  r           Rename pane
  x           Close pane
  ?           Help
  q           Quit
`)
		},
	}
}
