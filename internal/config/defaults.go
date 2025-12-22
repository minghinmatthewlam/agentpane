package config

func DefaultConfig() *Config {
	return &Config{
		DefaultPaneType: "codex",
		DefaultTemplate: "duo",
		Providers:       map[string]ProviderConfig{},
		Templates:       map[string]Template{},
	}
}
