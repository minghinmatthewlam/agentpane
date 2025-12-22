package dashboard

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
	"github.com/minghinmatthewlam/agentpane/internal/tui/dialogs"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dialog != nil {
		return m.updateDialog(msg)
	}

	// Handle filter input when active
	if m.filterActive {
		return m.handleFilterInput(msg)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tooNarrow = msg.Width < minWidth || msg.Height < minHeight
		return m, nil
	case snapshotMsg:
		m.snapshot = msg.snapshot
		m.errorMsg = ""
		return m, m.scheduleRefresh()
	case errMsg:
		m.errorMsg = msg.err.Error()
		return m, m.scheduleRefresh()
	case templatesMsg:
		m.templates = msg.templates
		return m, nil
	case tickMsg:
		return m, m.refreshSnapshot()
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "t":
		if m.tab == TabSessions {
			m.tab = TabTemplates
			m.focus = FocusLeft
		} else {
			m.tab = TabSessions
			m.focus = FocusLeft
		}
		return m, nil
	case "tab":
		if m.focus == FocusLeft {
			m.focus = FocusRight
		} else {
			m.focus = FocusLeft
		}
		return m, nil
	case "left", "h":
		// Move to previous session tab
		if m.sessionIndex > 0 {
			m.sessionIndex--
			m.paneIndex = 0
		}
		return m, nil
	case "right", "l":
		// Move to next session tab
		filtered := m.filteredSessions()
		if m.sessionIndex < len(filtered)-1 {
			m.sessionIndex++
			m.paneIndex = 0
		}
		return m, nil
	case "up", "k":
		if m.focus == FocusLeft {
			if m.tab == TabTemplates {
				if m.templateIndex > 0 {
					m.templateIndex--
				}
			} else if m.sessionIndex > 0 {
				m.sessionIndex--
				m.paneIndex = 0
			}
		} else if m.paneIndex > 0 {
			m.paneIndex--
		}
		return m, nil
	case "down", "j":
		if m.focus == FocusLeft {
			if m.tab == TabTemplates {
				if m.templateIndex < len(m.templates)-1 {
					m.templateIndex++
				}
			} else {
				filtered := m.filteredSessions()
				if m.sessionIndex < len(filtered)-1 {
					m.sessionIndex++
					m.paneIndex = 0
				}
			}
		} else {
			session := m.selectedSession()
			if session != nil && m.paneIndex < len(session.Panes)-1 {
				m.paneIndex++
			}
		}
		return m, nil
	case "enter":
		if m.tab == TabSessions {
			if m.focus == FocusLeft && m.selectedSession() != nil {
				return m, m.attachToSession(m.selectedSession().Name)
			}
		} else if m.tab == TabTemplates {
			if tmpl := m.selectedTemplate(); tmpl != nil {
				sessionName := m.selectedSessionName()
				if sessionName == "" {
					m.errorMsg = "no session selected"
					return m, nil
				}
				m.confirmAction = confirmApplyTemplate
				m.confirmSession = sessionName
				m.confirmTemplate = tmpl.Name
				m.dialog = dialogs.NewConfirm(
					"Apply template?",
					fmt.Sprintf("Replace panes in session %s with template %s?", sessionName, tmpl.Name),
				)
				return m, nil
			}
		}
		return m, nil
	case "a":
		m.dialog = dialogs.NewAddPane()
		return m, nil
	case "r":
		if pane := m.selectedPane(); pane != nil {
			session := m.selectedSession()
			if session != nil {
				m.renameSession = session.Name
				m.renamePaneID = pane.ID
				m.dialog = dialogs.NewRename(pane.Title)
			}
			return m, nil
		}
		return m, nil
	case "c":
		// Quick-add Claude pane
		m.addPaneType = domain.PaneClaude
		return m, tea.Quit
	case "x":
		// Quick-add Codex pane
		m.addPaneType = domain.PaneCodex
		return m, tea.Quit
	case "s":
		// Quick-add Shell pane
		m.addPaneType = domain.PaneShell
		return m, tea.Quit
	case "d":
		// Delete/close pane
		if pane := m.selectedPane(); pane != nil {
			m.confirmAction = confirmClosePane
			m.confirmPaneID = pane.ID
			m.dialog = dialogs.NewConfirm(
				"Close pane?",
				"This will kill the pane and any running processes.",
			)
			return m, nil
		}
		return m, nil
	case "?":
		m.dialog = dialogs.NewHelp()
		return m, nil
	case "/":
		// Activate session filter
		m.filterActive = true
		m.filterInput.Focus()
		return m, nil
	}

	return m, nil
}

func (m Model) handleFilterInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "esc":
			// Close filter, keep the filter text (enter) or clear it (esc)
			if msg.String() == "esc" {
				m.filterInput.SetValue("")
			}
			m.filterActive = false
			m.filterInput.Blur()
			// Reset session index to stay within filtered bounds
			filtered := m.filteredSessions()
			if m.sessionIndex >= len(filtered) {
				m.sessionIndex = 0
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.tooNarrow = msg.Width < minWidth || msg.Height < minHeight
		return m, nil
	}

	// Forward to text input
	var cmd tea.Cmd
	m.filterInput, cmd = m.filterInput.Update(msg)
	return m, cmd
}

// filteredSessions returns sessions matching the current filter
func (m Model) filteredSessions() []domain.Session {
	filter := strings.ToLower(strings.TrimSpace(m.filterInput.Value()))
	if filter == "" {
		return m.snapshot.Sessions
	}

	var filtered []domain.Session
	for _, s := range m.snapshot.Sessions {
		if strings.Contains(strings.ToLower(s.Name), filter) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

func (m Model) selectedSession() *domain.Session {
	filtered := m.filteredSessions()
	if m.sessionIndex < len(filtered) {
		return &filtered[m.sessionIndex]
	}
	return nil
}

func (m Model) selectedPane() *domain.Pane {
	session := m.selectedSession()
	if session == nil {
		return nil
	}
	if m.paneIndex < len(session.Panes) {
		return &session.Panes[m.paneIndex]
	}
	return nil
}

func (m Model) selectedTemplate() *app.TemplateSummary {
	if len(m.templates) == 0 {
		return nil
	}
	if m.templateIndex < 0 || m.templateIndex >= len(m.templates) {
		return nil
	}
	return &m.templates[m.templateIndex]
}

func (m Model) selectedSessionName() string {
	if m.snapshot.CurrentSession != "" {
		return m.snapshot.CurrentSession
	}
	if s := m.selectedSession(); s != nil {
		return s.Name
	}
	return ""
}

func (m Model) attachToSession(name string) tea.Cmd {
	return tea.Sequence(
		tea.ExitAltScreen,
		func() tea.Msg {
			_ = m.app.Attach(name)
			return tea.Quit()
		},
	)
}

func (m Model) scheduleRefresh() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) updateDialog(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dialog == nil {
		return m, nil
	}

	var cmd tea.Cmd
	m.dialog, cmd = m.dialog.Update(msg)

	switch msg := msg.(type) {
	case dialogs.AddPaneResult:
		m.dialog = nil
		if msg.Cancelled {
			return m, nil
		}
		result, err := m.app.Add(app.AddOptions{Type: domain.PaneType(msg.Type)})
		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		if result.FellBackToShell {
			m.statusMsg = "provider unavailable; created shell pane"
		} else {
			m.statusMsg = fmt.Sprintf("created pane %s", result.Title)
		}
		return m, m.refreshSnapshot()
	case dialogs.RenameResult:
		m.dialog = nil
		if msg.Cancelled {
			m.renameSession = ""
			m.renamePaneID = ""
			return m, nil
		}
		if strings.TrimSpace(msg.Title) == "" {
			m.errorMsg = "title cannot be empty"
			return m, nil
		}
		if _, err := m.app.Rename(app.RenameOptions{
			Title:   msg.Title,
			Session: m.renameSession,
			PaneID:  m.renamePaneID,
		}); err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}
		m.renameSession = ""
		m.renamePaneID = ""
		m.statusMsg = "pane renamed"
		return m, m.refreshSnapshot()
	case dialogs.ConfirmResult:
		m.dialog = nil
		if !msg.Accepted {
			m.confirmAction = confirmNone
			return m, nil
		}
		switch m.confirmAction {
		case confirmClosePane:
			if err := m.app.ClosePane(m.confirmPaneID); err != nil {
				m.errorMsg = err.Error()
				return m, nil
			}
			m.statusMsg = "pane closed"
			m.confirmAction = confirmNone
			return m, m.refreshSnapshot()
		case confirmApplyTemplate:
			_, err := m.app.ApplyTemplate(app.ApplyTemplateOptions{
				Session:  m.confirmSession,
				Template: m.confirmTemplate,
				Force:    true,
			})
			if err != nil {
				m.errorMsg = err.Error()
				m.confirmAction = confirmNone
				return m, nil
			}
			// Auto-attach to the session after applying template
			m.confirmAction = confirmNone
			return m, m.attachToSession(m.confirmSession)
		default:
			m.confirmAction = confirmNone
		}
	case dialogs.HelpResult:
		m.dialog = nil
		return m, nil
	}

	return m, cmd
}
