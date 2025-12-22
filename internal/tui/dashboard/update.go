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
		// Capture pane content for preview
		cmds := []tea.Cmd{m.scheduleRefresh()}
		if capCmd := m.capturePaneContent(); capCmd != nil {
			cmds = append(cmds, capCmd)
		}
		return m, tea.Batch(cmds...)
	case capturedContentMsg:
		m.capturedContent = msg.content
		return m, nil
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

type capturedContentMsg struct {
	content map[string]string
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "tab":
		// Tab switches between Sessions and Templates tabs
		if m.tab == TabSessions {
			m.tab = TabTemplates
			m.focus = FocusLeft
		} else {
			m.tab = TabSessions
		}
		return m, nil
	case "left", "h":
		// Jump to previous session in tree
		if m.tab == TabSessions {
			tree := m.buildTree()
			// Find current session and jump to previous one
			for i := m.treeIndex - 1; i >= 0; i-- {
				if tree[i].Type == ItemSession {
					m.treeIndex = i
					break
				}
			}
		}
		return m, nil
	case "right", "l":
		// Jump to next session in tree
		if m.tab == TabSessions {
			tree := m.buildTree()
			// Find next session
			for i := m.treeIndex + 1; i < len(tree); i++ {
				if tree[i].Type == ItemSession {
					m.treeIndex = i
					break
				}
			}
		}
		return m, nil
	case "up":
		if m.tab == TabTemplates {
			if m.focus == FocusLeft {
				if m.templateIndex > 0 {
					m.templateIndex--
				}
			}
		} else {
			// Sessions tab - navigate tree
			if m.treeIndex > 0 {
				m.treeIndex--
			}
		}
		return m, nil
	case "down":
		if m.tab == TabTemplates {
			if m.focus == FocusLeft {
				if m.templateIndex < len(m.templates)-1 {
					m.templateIndex++
				}
			}
		} else {
			// Sessions tab - navigate tree
			tree := m.buildTree()
			if m.treeIndex < len(tree)-1 {
				m.treeIndex++
			}
		}
		return m, nil
	case "k":
		// Kill session (when cursor is on a session)
		if m.tab == TabSessions {
			item := m.selectedTreeItem()
			if item != nil && item.Type == ItemSession {
				m.confirmAction = confirmKillSession
				m.confirmSession = item.Session
				m.dialog = dialogs.NewConfirm(
					"Kill session?",
					fmt.Sprintf("This will kill session '%s' and all its panes.", item.Session),
				)
			}
		}
		return m, nil
	case "enter":
		if m.tab == TabSessions {
			if session := m.selectedSession(); session != nil {
				m.attachSession = session.Name
				return m, tea.Quit
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
		// Rename only works when cursor is on a pane
		if pane := m.selectedPane(); pane != nil {
			session := m.selectedSession()
			if session != nil {
				m.renameSession = session.Name
				m.renamePaneID = pane.ID
				m.dialog = dialogs.NewRename(pane.Title)
			}
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
		// Delete/close pane - only works when cursor is on a pane
		if pane := m.selectedPane(); pane != nil {
			m.confirmAction = confirmClosePane
			m.confirmPaneID = pane.ID
			m.dialog = dialogs.NewConfirm(
				"Close pane?",
				"This will kill the pane and any running processes.",
			)
		}
		return m, nil
	case "o":
		// Open new session
		m.dialog = dialogs.NewOpenSession()
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
			// Reset tree index to stay within bounds
			tree := m.buildTree()
			if m.treeIndex >= len(tree) {
				m.treeIndex = 0
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

// buildTree creates a flattened tree of sessions and their panes
func (m Model) buildTree() []TreeItem {
	filtered := m.filteredSessions()
	var items []TreeItem
	for _, session := range filtered {
		items = append(items, TreeItem{
			Type:    ItemSession,
			Session: session.Name,
		})
		for i := range session.Panes {
			items = append(items, TreeItem{
				Type:    ItemPane,
				Session: session.Name,
				Pane:    &session.Panes[i],
			})
		}
	}
	return items
}

// selectedTreeItem returns the currently selected tree item
func (m Model) selectedTreeItem() *TreeItem {
	tree := m.buildTree()
	if m.treeIndex >= 0 && m.treeIndex < len(tree) {
		return &tree[m.treeIndex]
	}
	return nil
}

// selectedSession returns the session containing the current tree cursor
func (m Model) selectedSession() *domain.Session {
	item := m.selectedTreeItem()
	if item == nil {
		return nil
	}
	// Find the session by name
	for i := range m.snapshot.Sessions {
		if m.snapshot.Sessions[i].Name == item.Session {
			return &m.snapshot.Sessions[i]
		}
	}
	return nil
}

// selectedPane returns the pane if the cursor is on a pane item
func (m Model) selectedPane() *domain.Pane {
	item := m.selectedTreeItem()
	if item == nil || item.Type != ItemPane {
		return nil
	}
	return item.Pane
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
	case dialogs.OpenSessionResult:
		m.dialog = nil
		if msg.Cancelled {
			return m, nil
		}
		if strings.TrimSpace(msg.Path) == "" {
			m.errorMsg = "path cannot be empty"
			return m, nil
		}
		// Store the path and quit - session will be created after dashboard exits
		m.openSessionPath = msg.Path
		return m, tea.Quit
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
			m.attachSession = m.confirmSession
			return m, tea.Quit
		case confirmKillSession:
			if err := m.app.KillSession(m.confirmSession); err != nil {
				m.errorMsg = err.Error()
				m.confirmAction = confirmNone
				return m, nil
			}
			m.statusMsg = fmt.Sprintf("session '%s' killed", m.confirmSession)
			m.confirmAction = confirmNone
			m.treeIndex = 0 // Reset to first item
			return m, m.refreshSnapshot()
		default:
			m.confirmAction = confirmNone
		}
	case dialogs.HelpResult:
		m.dialog = nil
		return m, nil
	}

	return m, cmd
}

// capturePaneContent captures content from panes in the selected session
func (m Model) capturePaneContent() tea.Cmd {
	session := m.selectedSession()
	if session == nil {
		return nil
	}

	return func() tea.Msg {
		content := make(map[string]string)
		for _, pane := range session.Panes {
			captured, err := m.app.CapturePaneContent(pane.ID)
			if err == nil {
				content[pane.ID] = captured
			}
		}
		return capturedContentMsg{content: content}
	}
}
