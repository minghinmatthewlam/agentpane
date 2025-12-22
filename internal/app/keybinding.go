package app

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

const (
	keybindingLine   = `bind-key g run-shell "agentpane dashboard --tmux-window"`
	keybindingMarker = "agentpane dashboard --tmux-window"
	legacyMarker     = "agentpane popup"
)

// EnsureKeybinding checks if the dashboard keybinding exists in tmux.conf
// and adds it if not present. Returns true if it was added.
func (a *App) EnsureKeybinding() (bool, error) {
	confPath := tmuxConfPath()

	// Check if binding already exists
	exists, err := keybindingExists(confPath)
	if err != nil {
		return false, err
	}
	if exists {
		return false, nil
	}

	// Add the keybinding
	if err := appendKeybinding(confPath); err != nil {
		return false, err
	}

	return true, nil
}

func tmuxConfPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// Check XDG location first
	xdgPath := filepath.Join(home, ".config", "tmux", "tmux.conf")
	if _, err := os.Stat(xdgPath); err == nil {
		return xdgPath
	}

	// Default to ~/.tmux.conf
	return filepath.Join(home, ".tmux.conf")
}

func keybindingExists(confPath string) (bool, error) {
	if confPath == "" {
		return false, nil
	}

	file, err := os.Open(confPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Check if line contains our keybinding marker
		if strings.Contains(line, keybindingMarker) || strings.Contains(line, legacyMarker) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

func appendKeybinding(confPath string) error {
	if confPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		confPath = filepath.Join(home, ".tmux.conf")
	}

	// Open file for appending, create if doesn't exist
	file, err := os.OpenFile(confPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Add newline before our content if file is not empty
	info, err := file.Stat()
	if err != nil {
		return err
	}

	content := "\n# agentpane: quick dashboard access (Prefix+g)\n" + keybindingLine + "\n"
	if info.Size() == 0 {
		content = "# agentpane: quick dashboard access (Prefix+g)\n" + keybindingLine + "\n"
	}

	_, err = file.WriteString(content)
	return err
}
