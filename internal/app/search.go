package app

import (
	"strings"

	"github.com/minghinmatthewlam/agentpane/internal/domain"
)

type SearchResult struct {
	Session string
	PaneID  string
	Title   string
	Type    domain.PaneType
}

func (a *App) Search(query string) ([]SearchResult, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return nil, nil
	}

	snapshot, err := a.Snapshot()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, session := range snapshot.Sessions {
		if strings.Contains(strings.ToLower(session.Name), query) {
			results = append(results, SearchResult{
				Session: session.Name,
			})
		}
		for _, pane := range session.Panes {
			if strings.Contains(strings.ToLower(pane.Title), query) {
				results = append(results, SearchResult{
					Session: session.Name,
					PaneID:  pane.ID,
					Title:   pane.Title,
					Type:    pane.Type,
				})
			}
		}
	}

	return results, nil
}
