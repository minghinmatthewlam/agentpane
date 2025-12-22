package dashboard

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

type Tab int

const (
	TabSessions Tab = iota
	TabTemplates
)

type Focus int

const (
	FocusLeft Focus = iota
	FocusRight
)

type confirmAction int

const (
	confirmNone confirmAction = iota
	confirmClosePane
	confirmApplyTemplate
)

const (
	minWidth  = 80
	minHeight = 20
)

type Model struct {
	app *app.App

	snapshot domain.Snapshot

	tab           Tab
	focus         Focus
	sessionIndex  int
	paneIndex     int
	templateIndex int
	width         int
	height        int

	dialog tea.Model

	statusMsg string
	errorMsg  string

	tooNarrow bool

	templates []app.TemplateSummary

	confirmAction   confirmAction
	confirmSession  string
	confirmPaneID   string
	confirmTemplate string
}

func NewModel(a *app.App) Model {
	return Model{
		app:   a,
		tab:   TabSessions,
		focus: FocusLeft,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.refreshSnapshot(),
		m.refreshTemplates(),
		tea.EnterAltScreen,
	)
}

func (m Model) refreshSnapshot() tea.Cmd {
	return func() tea.Msg {
		snapshot, err := m.app.Snapshot()
		if err != nil {
			return errMsg{err}
		}
		return snapshotMsg{snapshot}
	}
}

func (m Model) refreshTemplates() tea.Cmd {
	return func() tea.Msg {
		templates, err := m.app.ListTemplates()
		if err != nil {
			return errMsg{err}
		}
		return templatesMsg{templates: templates}
	}
}

type snapshotMsg struct {
	snapshot domain.Snapshot
}

type errMsg struct {
	err error
}

type tickMsg struct{}

type templatesMsg struct {
	templates []app.TemplateSummary
}
