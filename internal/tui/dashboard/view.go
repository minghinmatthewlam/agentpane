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
		leftPanel = m.renderSessionsList()
		rightPanel = m.renderPanesList()
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

func (m Model) renderSessionsList() string {
	var b strings.Builder
	b.WriteString(common.TitleStyle.Render("Sessions"))
	b.WriteString("\n\n")

	for i, session := range m.snapshot.Sessions {
		cursor := "  "
		style := common.NormalStyle
		if i == m.sessionIndex {
			cursor = "→ "
			if m.focus == FocusLeft {
				style = common.SelectedStyle
			} else {
				style = common.DimSelectedStyle
			}
		}

		indicator := "○"
		if sessionHasActive(session) {
			indicator = "●"
		}

		line := fmt.Sprintf("%s%s %s (%d panes)",
			cursor, indicator, session.Name, len(session.Panes))
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderPanesList() string {
	var b strings.Builder

	session := m.selectedSession()
	if session == nil {
		return "No session selected"
	}

	b.WriteString(common.TitleStyle.Render(fmt.Sprintf("Panes in %s", session.Name)))
	b.WriteString("\n\n")

	for i, pane := range session.Panes {
		cursor := "  "
		style := common.NormalStyle
		if i == m.paneIndex {
			cursor = "→ "
			if m.focus == FocusRight {
				style = common.SelectedStyle
			} else {
				style = common.DimSelectedStyle
			}
		}

		indicator := "○"
		if pane.Status == domain.StatusActive {
			indicator = "●"
		} else if pane.Status == domain.StatusUnknown {
			indicator = "?"
		}

		typeBadge := fmt.Sprintf("[%s]", pane.Type)
		line := fmt.Sprintf("%s%-20s %s %s", cursor, pane.Title, typeBadge, indicator)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderHeader() string {
	tabs := []string{"Sessions", "Templates"}
	var tabViews []string
	for i, tab := range tabs {
		style := common.TabStyle
		if (i == 0 && m.tab == TabSessions) || (i == 1 && m.tab == TabTemplates) {
			style = common.ActiveTabStyle
		}
		tabViews = append(tabViews, style.Render(tab))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabViews...)
}

func (m Model) renderFooter() string {
	keys := []string{"[Tab] focus", "[t] tabs", "[q] quit"}
	if m.tab == TabSessions {
		keys = append([]string{"[Enter] attach", "[a] add", "[r] rename", "[x] close", "[?] help"}, keys...)
	} else {
		keys = append([]string{"[Enter] apply", "[?] help"}, keys...)
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
