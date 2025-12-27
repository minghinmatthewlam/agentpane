package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case snapshotMsg:
		m.snapshot = msg.snapshot
		m.items = buildItems(m.snapshot)
		m.errorMsg = ""
		if m.selectedSession == "" && len(m.items) > 0 {
			m.index = 0
			m.syncSelection()
		} else {
			m.restoreSelection()
		}
		return m, m.scheduleRefresh()
	case errMsg:
		m.errorMsg = msg.err.Error()
		return m, m.scheduleRefresh()
	case tickMsg:
		return m, m.refreshSnapshot()
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if m.index > 0 {
			m.index--
		}
		m.syncSelection()
		return m, nil
	case "down", "j":
		if m.index < len(m.items)-1 {
			m.index++
		}
		m.syncSelection()
		return m, nil
	case "enter":
		return m.attachSelected()
	}
	return m, nil
}

func (m Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionRelease {
		return m, nil
	}
	if msg.Button != tea.MouseButtonLeft {
		return m, nil
	}
	idx := msg.Y
	if idx < 0 || idx >= len(m.items) {
		return m, nil
	}
	m.index = idx
	m.syncSelection()
	return m.attachSelected()
}

func (m Model) attachSelected() (tea.Model, tea.Cmd) {
	if len(m.items) == 0 {
		return m, nil
	}
	selected := m.items[m.index]
	if selected.session == "" {
		return m, nil
	}
	if err := m.app.Attach(selected.session); err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}
	return m, nil
}

func (m *Model) syncSelection() {
	if len(m.items) == 0 || m.index < 0 || m.index >= len(m.items) {
		m.selectedKind = itemSession
		m.selectedSession = ""
		m.selectedPaneID = ""
		return
	}
	item := m.items[m.index]
	m.selectedKind = item.kind
	m.selectedSession = item.session
	if item.kind == itemPane && item.pane != nil {
		m.selectedPaneID = item.pane.ID
	} else {
		m.selectedPaneID = ""
	}
}

func (m *Model) restoreSelection() {
	if len(m.items) == 0 {
		m.index = 0
		m.syncSelection()
		return
	}

	if m.selectedKind == itemPane && m.selectedPaneID != "" {
		for i, it := range m.items {
			if it.kind == itemPane && it.pane != nil && it.pane.ID == m.selectedPaneID {
				m.index = i
				m.syncSelection()
				return
			}
		}
	}

	if m.selectedSession != "" {
		for i, it := range m.items {
			if it.kind == itemSession && it.session == m.selectedSession {
				m.index = i
				m.syncSelection()
				return
			}
		}
	}

	if m.index >= len(m.items) {
		m.index = len(m.items) - 1
	}
	m.syncSelection()
}
