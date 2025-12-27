package overlay

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

type itemKind int

const (
	itemSession itemKind = iota
	itemPane
)

type item struct {
	kind      itemKind
	session   string
	pane      *domain.Pane
	indicator string
}

type Model struct {
	app      *app.App
	snapshot domain.Snapshot
	items    []item
	index    int
	width    int
	height   int
	errorMsg string

	selectedKind    itemKind
	selectedSession string
	selectedPaneID  string
}

type snapshotMsg struct {
	snapshot domain.Snapshot
}

type errMsg struct {
	err error
}

type tickMsg struct{}

func NewModel(a *app.App) Model {
	return Model{app: a}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.refreshSnapshot(),
	)
}

func (m Model) refreshSnapshot() tea.Cmd {
	return func() tea.Msg {
		snapshot, err := m.app.Snapshot()
		if err != nil {
			return errMsg{err: err}
		}
		return snapshotMsg{snapshot: snapshot}
	}
}

func (m Model) scheduleRefresh() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
