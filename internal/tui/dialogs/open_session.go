package dialogs

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OpenSessionResult struct {
	Cancelled bool
	Path      string
}

type OpenSessionModel struct {
	input textinput.Model
}

func NewOpenSession() OpenSessionModel {
	ti := textinput.New()
	ti.Placeholder = "~/path/to/project"
	ti.Width = 40

	// Pre-fill with home directory
	if home, err := os.UserHomeDir(); err == nil {
		ti.SetValue(home + "/")
	}

	ti.Focus()
	return OpenSessionModel{input: ti}
}

func (m OpenSessionModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m OpenSessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return OpenSessionResult{Cancelled: true} }
		case "enter":
			path := expandPath(m.input.Value())
			return m, func() tea.Msg { return OpenSessionResult{Path: path} }
		case "tab":
			// Simple tab completion - expand to first matching directory
			m.input.SetValue(tabComplete(m.input.Value()))
			m.input.CursorEnd()
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m OpenSessionModel) View() string {
	content := "Open Session\n\n" +
		"Path: " + m.input.View() + "\n\n" +
		"[Enter] open  [Tab] complete  [Esc] cancel"
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(content)
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	}
	return path
}

// tabComplete attempts to complete the path to the first matching directory
func tabComplete(path string) string {
	expanded := expandPath(path)

	// If path ends with /, list contents
	if strings.HasSuffix(expanded, "/") {
		entries, err := os.ReadDir(expanded)
		if err != nil {
			return path
		}
		for _, e := range entries {
			if e.IsDir() {
				return path + e.Name() + "/"
			}
		}
		return path
	}

	// Otherwise, complete the last component
	dir := filepath.Dir(expanded)
	base := filepath.Base(expanded)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return path
	}

	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), base) {
			// Reconstruct path preserving ~ if used
			if strings.HasPrefix(path, "~/") {
				return "~/" + filepath.Join(filepath.Dir(path[2:]), e.Name()) + "/"
			}
			return filepath.Join(dir, e.Name()) + "/"
		}
	}

	return path
}
