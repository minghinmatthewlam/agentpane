package dashboard

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/tui/common"
)

func (m Model) View() string {
	if m.tooNarrow {
		return m.renderTooNarrow()
	}

	if m.dialog != nil {
		return m.renderWithDialog()
	}

	return m.renderDashboard()
}

func (m Model) renderTooNarrow() string {
	msg := fmt.Sprintf(
		"Terminal too narrow.\n\nResize to at least %d columns.\nCurrent: %d x %d",
		minWidth, m.width, m.height,
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, msg)
}

func (m Model) renderDashboard() string {
	header := m.renderHeader()

	var leftPanel string
	var rightPanel string
	if m.tab == TabTemplates {
		leftPanel = m.renderTemplatesList()
		rightPanel = m.renderTemplatePreview()
	} else {
		leftPanel = m.renderTree()
		rightPanel = m.renderPanePreview()
	}

	leftWidth := m.width / 3
	rightWidth := m.width - leftWidth - 3

	leftStyled := common.PanelStyle.Width(leftWidth).Render(leftPanel)
	rightStyled := common.PanelStyle.Width(rightWidth).Render(rightPanel)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftStyled, rightStyled)

	status := m.renderStatus()
	footer := m.renderFooter()
	if status != "" {
		return lipgloss.JoinVertical(lipgloss.Left, header, content, status, footer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func (m Model) renderTree() string {
	var b strings.Builder
	b.WriteString(common.TitleStyle.Render("Sessions"))
	b.WriteString("\n\n")

	// Show filter input if active or has value
	filterValue := strings.TrimSpace(m.filterInput.Value())
	if m.filterActive || filterValue != "" {
		if m.filterActive {
			b.WriteString("Filter: " + m.filterInput.View())
		} else {
			b.WriteString(common.DimSelectedStyle.Render("Filter: " + filterValue))
		}
		b.WriteString("\n\n")
	}

	tree := m.buildTree()
	if len(tree) == 0 {
		if filterValue != "" {
			b.WriteString(common.DimSelectedStyle.Render("No matching sessions"))
		} else {
			b.WriteString(common.DimSelectedStyle.Render("No sessions"))
		}
		b.WriteString("\n")
		return b.String()
	}

	for i, item := range tree {
		cursor := "  "
		style := common.NormalStyle
		if i == m.treeIndex {
			cursor = "→ "
			style = common.SelectedStyle
		}

		if item.Type == ItemSession {
			// Session row
			indicator := "○"
			// Check if any pane in this session is active
			for j := range m.snapshot.Sessions {
				if m.snapshot.Sessions[j].Name == item.Session {
					if sessionHasActive(m.snapshot.Sessions[j]) {
						indicator = "●"
					}
					break
				}
			}

			// Mark current session
			name := item.Session
			if item.Session == m.snapshot.CurrentSession {
				name = "● " + name
			}

			line := fmt.Sprintf("%s%s %s", cursor, indicator, name)
			b.WriteString(style.Render(line))
		} else {
			// Pane row (indented)
			pane := item.Pane
			indicator := "○"
			if pane.Status == domain.StatusActive {
				indicator = "●"
			} else if pane.Status == domain.StatusUnknown {
				indicator = "?"
			}

			typeBadge := fmt.Sprintf("[%s]", pane.Type)
			line := fmt.Sprintf("%s    %s %s %s", cursor, indicator, pane.Title, typeBadge)
			b.WriteString(style.Render(line))
		}
		b.WriteString("\n")
	}

	// Show count if filtered
	if filterValue != "" && len(m.filteredSessions()) != len(m.snapshot.Sessions) {
		b.WriteString("\n")
		b.WriteString(common.DimSelectedStyle.Render(fmt.Sprintf("showing %d of %d sessions", len(m.filteredSessions()), len(m.snapshot.Sessions))))
	}

	return b.String()
}

func (m Model) renderPanePreview() string {
	session := m.selectedSession()
	if session == nil {
		return common.DimSelectedStyle.Render("Select a session to preview panes")
	}

	var b strings.Builder
	b.WriteString(common.TitleStyle.Render(fmt.Sprintf("Preview: %s", session.Name)))
	b.WriteString("\n\n")

	if len(session.Panes) == 0 {
		b.WriteString(common.DimSelectedStyle.Render("No panes in this session"))
		return b.String()
	}

	// Calculate height per pane (divide available height among panes)
	// Reserve some lines for header and spacing
	availableHeight := m.height - 10
	if availableHeight < 5 {
		availableHeight = 5
	}
	linesPerPane := availableHeight / len(session.Panes)
	if linesPerPane < 3 {
		linesPerPane = 3
	}
	if linesPerPane > 15 {
		linesPerPane = 15
	}

	paneStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(0, 1)

	for _, pane := range session.Panes {
		// Pane header
		indicator := "○"
		if pane.Status == domain.StatusActive {
			indicator = "●"
		}
		header := fmt.Sprintf("%s %s [%s]", indicator, pane.Title, pane.Type)
		b.WriteString(common.DimSelectedStyle.Render(header))
		b.WriteString("\n")

		// Pane content
		content := m.capturedContent[pane.ID]
		if content == "" {
			content = "(no content)"
		}

		// Truncate content to fit
		lines := strings.Split(content, "\n")
		maxLines := linesPerPane - 2 // Account for header and border
		if maxLines < 1 {
			maxLines = 1
		}
		if len(lines) > maxLines {
			lines = lines[len(lines)-maxLines:] // Show last N lines
		}

		// Trim trailing empty lines
		for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
			lines = lines[:len(lines)-1]
		}

		previewContent := strings.Join(lines, "\n")
		if previewContent == "" {
			previewContent = "(empty)"
		}

		b.WriteString(paneStyle.Render(previewContent))
		b.WriteString("\n\n")
	}

	return b.String()
}

func (m Model) renderHeader() string {
	// Line 1: Session tabs for quick switching between repos
	sessionBar := m.renderSessionTabs()

	// Line 2: Mode tabs (Sessions/Templates)
	tabs := []string{"Sessions", "Templates"}
	var tabViews []string
	for i, tab := range tabs {
		style := common.TabStyle
		if (i == 0 && m.tab == TabSessions) || (i == 1 && m.tab == TabTemplates) {
			style = common.ActiveTabStyle
		}
		tabViews = append(tabViews, style.Render(tab))
	}
	modeTabs := lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)

	return lipgloss.JoinVertical(lipgloss.Left, sessionBar, modeTabs)
}

func (m Model) renderSessionTabs() string {
	filtered := m.filteredSessions()
	if len(filtered) == 0 {
		return common.DimSelectedStyle.Render("⚡ No sessions")
	}

	// Get currently selected session from tree
	selectedSession := m.selectedSession()
	var selectedName string
	if selectedSession != nil {
		selectedName = selectedSession.Name
	}

	var tabs []string
	tabs = append(tabs, "⚡")

	for _, session := range filtered {
		name := session.Name

		// Add marker if this is the currently attached session
		if session.Name == m.snapshot.CurrentSession {
			name = "● " + name
		}

		style := common.SessionTabStyle
		if session.Name == selectedName {
			style = common.ActiveSessionTabStyle
		}
		tabs = append(tabs, style.Render(name))
	}

	return strings.Join(tabs, " ")
}

func (m Model) renderFooter() string {
	var keys []string
	if m.tab == TabSessions {
		keys = []string{"[Enter] attach", "[o] open", "[c] claude", "[x] codex", "[s] shell", "[k] kill", "[/] filter", "[Tab] templates", "[q] quit"}
	} else {
		keys = []string{"[Enter] apply", "[Tab] sessions", "[q] quit"}
	}
	return common.FooterStyle.Render(strings.Join(keys, "  "))
}

func (m Model) renderWithDialog() string {
	base := m.renderDashboard()
	if m.dialog == nil {
		return base
	}
	dialogView := m.dialog.View()
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, dialogView)
}

func (m Model) renderStatus() string {
	if m.errorMsg != "" {
		return common.ErrorStyle.Render("Error: " + m.errorMsg)
	}
	if m.statusMsg != "" {
		return common.StatusStyle.Render(m.statusMsg)
	}
	return ""
}

func (m Model) renderTemplatesList() string {
	var b strings.Builder
	b.WriteString(common.TitleStyle.Render("Templates"))
	b.WriteString("\n\n")

	for i, tmpl := range m.templates {
		cursor := "  "
		style := common.NormalStyle
		if i == m.templateIndex {
			cursor = "→ "
			if m.focus == FocusLeft {
				style = common.SelectedStyle
			} else {
				style = common.DimSelectedStyle
			}
		}
		line := fmt.Sprintf("%s%s", cursor, tmpl.Name)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	if len(m.templates) == 0 {
		b.WriteString("No templates available\n")
	}
	return b.String()
}

func (m Model) renderTemplatePreview() string {
	tmpl := m.selectedTemplate()
	if tmpl == nil {
		return "No template selected"
	}

	var b strings.Builder
	b.WriteString(common.TitleStyle.Render(fmt.Sprintf("Template: %s", tmpl.Name)))
	b.WriteString("\n\n")
	if tmpl.Description != "" {
		b.WriteString(tmpl.Description)
		b.WriteString("\n\n")
	}
	b.WriteString("Panes:\n")
	for _, p := range tmpl.Panes {
		title := p.Type
		if p.Title != "" {
			title = fmt.Sprintf("%s (%s)", p.Type, p.Title)
		}
		b.WriteString(fmt.Sprintf("  - %s\n", title))
	}
	return b.String()
}

func sessionHasActive(s domain.Session) bool {
	for _, p := range s.Panes {
		if p.Status == domain.StatusActive {
			return true
		}
	}
	return false
}
