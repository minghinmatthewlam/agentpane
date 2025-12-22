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

func FindRepoConfigPath(startDir string) (string, bool, error) {
	dir := startDir
	for {
		path := RepoConfigPath(dir)
		if _, err := os.Stat(path); err == nil {
			return path, true, nil
		} else if err != nil && !os.IsNotExist(err) {
			return "", false, err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false, nil
		}
		dir = parent
	}
}
