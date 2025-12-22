package config

import (
	"os"
	"path/filepath"
)

func GlobalConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "agentpane", "config.yml"), nil
}

func RepoConfigPath(cwd string) string {
	return filepath.Join(cwd, ".agentpane.yml")
}

