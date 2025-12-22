package search

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/minghinmatthewlam/agentpane/internal/app"
)

type Model struct {
	app     *app.App
	input   textinput.Model
	results []app.SearchResult
	errMsg  string
}

func NewModel(a *app.App) Model {
	ti := textinput.New()
	ti.Placeholder = "Search sessions or pane titles"
	ti.Focus()
	return Model{
		app:   a,
		input: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

type resultsMsg struct {
	results []app.SearchResult
	err     error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		}
	case resultsMsg:
		if msg.err != nil {
			m.errMsg = msg.err.Error()
		} else {
			m.errMsg = ""
			m.results = msg.results
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, tea.Batch(cmd, m.searchCmd())
}

func (m Model) searchCmd() tea.Cmd {
	query := strings.TrimSpace(m.input.Value())
	return func() tea.Msg {
		results, err := m.app.Search(query)
		return resultsMsg{results: results, err: err}
	}
}
