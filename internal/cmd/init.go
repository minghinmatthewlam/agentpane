package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func NewInitCmd(a *app.App) *cobra.Command {
	var fromCurrent bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Generate .agentpane.yml in current repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			targetDir := cwd
			if root, ok, err := findRepoRoot(cwd); err == nil && ok {
				targetDir = root
			} else if err != nil {
				return err
			}

			path := config.RepoConfigPath(targetDir)
			if _, err := os.Stat(path); err == nil {
				return fmt.Errorf("%s already exists", path)
			} else if err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}

			var rc config.RepoConfig
			if fromCurrent {
				if !a.InTmux() {
					return fmt.Errorf("init --from-current must be run inside tmux")
				}

				layout, session, err := a.SnapshotCurrentLayout()
				if err != nil {
					return err
				}
				rc.Session = session
				rc.Layout = layout
			} else {
				rc.Layout = config.Layout{
					Panes: []config.PaneSpec{{Type: "codex"}, {Type: "claude"}},
				}
			}

			data, err := yaml.Marshal(&rc)
			if err != nil {
				return err
			}

			// Ensure file ends with newline for nicer diffs.
			if !strings.HasSuffix(string(data), "\n") {
				data = append(data, '\n')
			}
			return os.WriteFile(path, data, 0o644)
		},
	}

	cmd.Flags().BoolVar(&fromCurrent, "from-current", false, "Generate layout from current tmux session")
	return cmd
}

func findRepoRoot(start string) (string, bool, error) {
	dir := start
	for {
		gitPath := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return dir, true, nil
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", false, err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}
