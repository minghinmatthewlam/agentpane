package state

import (
	"os"
	"path/filepath"
)

func DefaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share", "agentpane", "state.yml"), nil
}
