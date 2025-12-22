package config

type Config struct {
	DefaultPaneType string                    `yaml:"default_pane_type"`
	DefaultTemplate string                    `yaml:"default_template"`
	Providers       map[string]ProviderConfig `yaml:"providers"`
	Templates       map[string]Template       `yaml:"templates"`
}

type ProviderConfig struct {
	Command string `yaml:"command"`
}

type Template struct {
	Description string     `yaml:"description"`
	Panes       []PaneSpec `yaml:"panes"`
}

type PaneSpec struct {
	Type  string `yaml:"type"`
	Title string `yaml:"title,omitempty"`
}

type RepoConfig struct {
	Session string `yaml:"session,omitempty"`
	Layout  Layout `yaml:"layout"`
}

type Layout struct {
	Panes []PaneSpec `yaml:"panes"`
}

