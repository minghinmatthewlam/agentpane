package dialogs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
	input           textinput.Model
	completions     []string // available completions
	completionIndex int      // current index in completions
	completionBase  string   // original input before completing
}

func NewOpenSession() OpenSessionModel {
	ti := textinput.New()
	ti.Placeholder = "~/path/to/project"
	ti.Width = 40

	// Pre-fill with ~/ for cleaner display
	ti.SetValue("~/")

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
			return m.handleTab(), nil
		case "shift+tab", "up":
			// Navigate completions backward
			if len(m.completions) > 1 {
				return m.handleShiftTab(), nil
			}
			return m, nil
		case "down":
			// Navigate completions forward
			if len(m.completions) > 1 {
				return m.handleTabCycle(), nil
			}
			return m, nil
		default:
			// Any other key resets completions
			m.completions = nil
			m.completionIndex = 0
			m.completionBase = ""
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m OpenSessionModel) handleTab() OpenSessionModel {
	currentValue := m.input.Value()

	// If we have active completions and user is cycling through them
	if len(m.completions) > 0 && m.completionBase != "" {
		// Cycle to next completion
		m.completionIndex = (m.completionIndex + 1) % len(m.completions)
		newPath := buildCompletionPath(m.completionBase, m.completions[m.completionIndex])
		m.input.SetValue(newPath)
		m.input.CursorEnd()
		return m
	}

	// Calculate new completions
	m.completionBase = currentValue
	m.completions = getCompletions(currentValue)
	m.completionIndex = 0

	if len(m.completions) == 0 {
		// No completions - do nothing
		return m
	}

	// Apply the first completion (whether single or multiple matches)
	// For single match: this completes it
	// For multiple matches: this shows the list with first item selected
	newPath := buildCompletionPath(m.completionBase, m.completions[0])
	m.input.SetValue(newPath)
	m.input.CursorEnd()

	// If single match, clear completions for next Tab cycle
	if len(m.completions) == 1 {
		m.completions = nil
		m.completionBase = ""
	}

	return m
}

func (m OpenSessionModel) handleShiftTab() OpenSessionModel {
	if len(m.completions) <= 1 {
		return m
	}

	// Cycle to previous completion
	m.completionIndex--
	if m.completionIndex < 0 {
		m.completionIndex = len(m.completions) - 1
	}
	newPath := buildCompletionPath(m.completionBase, m.completions[m.completionIndex])
	m.input.SetValue(newPath)
	m.input.CursorEnd()

	return m
}

func (m OpenSessionModel) handleTabCycle() OpenSessionModel {
	if len(m.completions) <= 1 {
		return m
	}

	// Cycle to next completion
	m.completionIndex = (m.completionIndex + 1) % len(m.completions)
	newPath := buildCompletionPath(m.completionBase, m.completions[m.completionIndex])
	m.input.SetValue(newPath)
	m.input.CursorEnd()

	return m
}

func (m OpenSessionModel) View() string {
	selectedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	var lines []string
	lines = append(lines, "Open Session")
	lines = append(lines, "")
	lines = append(lines, "Path: "+m.input.View())

	// Show completions if multiple options
	if len(m.completions) > 1 {
		lines = append(lines, "")
		maxShow := 8
		for i, c := range m.completions {
			if i >= maxShow {
				remaining := len(m.completions) - maxShow
				lines = append(lines, dimStyle.Render(fmt.Sprintf("  ... and %d more", remaining)))
				break
			}
			if i == m.completionIndex {
				lines = append(lines, selectedStyle.Render("→ "+c+"/"))
			} else {
				lines = append(lines, dimStyle.Render("  "+c+"/"))
			}
		}
	}

	lines = append(lines, "")
	lines = append(lines, "[Enter] open  [↑/↓] navigate  [Tab] complete  [Esc] cancel")

	content := strings.Join(lines, "\n")

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2)
	return style.Render(content)
}

// expandPath expands ~ to home directory, preserving trailing slash
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			expanded := filepath.Join(home, path[2:])
			// Preserve trailing slash if original had it
			if strings.HasSuffix(path, "/") && !strings.HasSuffix(expanded, "/") {
				expanded += "/"
			}
			return expanded
		}
	}
	if path == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	}
	return path
}

// getCompletions returns all matching directory names for the given path
func getCompletions(path string) []string {
	expanded := expandPath(path)

	var dir, prefix string
	if strings.HasSuffix(expanded, "/") {
		// Path ends with / - list all dirs in this directory
		dir = expanded
		prefix = ""
	} else {
		// Complete the last component
		dir = filepath.Dir(expanded)
		prefix = filepath.Base(expanded)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var completions []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		// Skip hidden directories unless prefix starts with .
		if strings.HasPrefix(name, ".") && !strings.HasPrefix(prefix, ".") {
			continue
		}
		if prefix == "" || strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			completions = append(completions, name)
		}
	}

	sort.Strings(completions)
	return completions
}

// buildCompletionPath constructs the full path from base and completion
func buildCompletionPath(base, completion string) string {
	expanded := expandPath(base)

	var dir string
	if strings.HasSuffix(expanded, "/") {
		dir = expanded
	} else {
		dir = filepath.Dir(expanded)
	}

	fullPath := filepath.Join(dir, completion) + "/"

	// Preserve ~ prefix if original used it
	if strings.HasPrefix(base, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			if strings.HasPrefix(fullPath, home) {
				return "~" + fullPath[len(home):]
			}
		}
	}

	return fullPath
}
