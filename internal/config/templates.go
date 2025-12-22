package config

import (
	"embed"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed templates/*.yml
var templatesFS embed.FS

func LoadBuiltinTemplates() (map[string]Template, error) {
	templates := make(map[string]Template)

	entries, err := templatesFS.ReadDir("templates")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		data, err := templatesFS.ReadFile(filepath.Join("templates", entry.Name()))
		if err != nil {
			return nil, err
		}

		var tmpl Template
		if err := yaml.Unmarshal(data, &tmpl); err != nil {
			return nil, err
		}

		name := strings.TrimSuffix(entry.Name(), ".yml")
		templates[name] = tmpl
	}

	return templates, nil
}
