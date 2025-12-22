package dashboard

import (
	"github.com/charmbracelet/bubbles/textinput"
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
	confirmKillSession
)

const (
	minWidth  = 80
	minHeight = 20
)

// TreeItemType distinguishes between session and pane items in the tree
type TreeItemType int

const (
	ItemSession TreeItemType = iota
	ItemPane
)

// TreeItem represents an item in the session/pane tree
type TreeItem struct {
	Type    TreeItemType
	Session string       // session name
	Pane    *domain.Pane // nil for session items
}

type Model struct {
	app *app.App

	snapshot domain.Snapshot

	tab           Tab
	focus         Focus // used for Templates tab
	treeIndex     int   // cursor position in flattened tree (Sessions tab)
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

	// renameSession/renamePaneID track which pane to rename
	renameSession string
	renamePaneID  string

	// addPaneType is set when user presses c/x/s to quick-add a pane
	addPaneType domain.PaneType

	// attachSession is set when user wants to attach to a session after exit
	attachSession string

	// openSessionPath is set when user wants to open a new session at a path
	openSessionPath string

	// Session filtering
	filterInput  textinput.Model
	filterActive bool

	// Captured pane content for preview
	capturedContent map[string]string // paneID -> content
}

// AddPaneType returns the pane type to add after dashboard exits (empty if none)
func (m Model) AddPaneType() domain.PaneType {
	return m.addPaneType
}

// AttachSession returns the session name to attach to after dashboard exits (empty if none)
func (m Model) AttachSession() string {
	return m.attachSession
}

// OpenSessionPath returns the path to open a new session at after dashboard exits (empty if none)
func (m Model) OpenSessionPath() string {
	return m.openSessionPath
}

func NewModel(a *app.App) Model {
	ti := textinput.New()
	ti.Placeholder = "filter sessions..."
	ti.CharLimit = 50
	ti.Width = 20

	return Model{
		app:             a,
		tab:             TabSessions,
		focus:           FocusLeft,
		filterInput:     ti,
		capturedContent: make(map[string]string),
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
