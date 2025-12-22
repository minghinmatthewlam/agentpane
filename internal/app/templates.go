package app

import (
	"os"
	"sort"

	"github.com/minghinmatthewlam/agentpane/internal/config"
)

type TemplateSummary struct {
	Name        string
	Description string
	Panes       []config.PaneSpec
}

func (a *App) ListTemplates() ([]TemplateSummary, error) {
	cwd, _ := os.Getwd()
	loaded, err := config.LoadAll(cwd)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(loaded.Merged.Templates))
	for name := range loaded.Merged.Templates {
		names = append(names, name)
	}
	sort.Strings(names)

	out := make([]TemplateSummary, 0, len(names))
	for _, name := range names {
		tmpl := loaded.Merged.Templates[name]
		out = append(out, TemplateSummary{
			Name:        name,
			Description: tmpl.Description,
			Panes:       tmpl.Panes,
		})
	}
	return out, nil
}
