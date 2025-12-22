package config

func Merge(base *Config, overlay *Config) *Config {
	if base == nil {
		base = DefaultConfig()
	}
	out := &Config{
		DefaultPaneType: base.DefaultPaneType,
		DefaultTemplate: base.DefaultTemplate,
		Providers:       map[string]ProviderConfig{},
		Templates:       map[string]Template{},
	}

	for k, v := range base.Providers {
		out.Providers[k] = v
	}
	for k, v := range base.Templates {
		out.Templates[k] = v
	}

	if overlay == nil {
		return out
	}
	if overlay.DefaultPaneType != "" {
		out.DefaultPaneType = overlay.DefaultPaneType
	}
	if overlay.DefaultTemplate != "" {
		out.DefaultTemplate = overlay.DefaultTemplate
	}
	for k, v := range overlay.Providers {
		out.Providers[k] = v
	}
	for k, v := range overlay.Templates {
		out.Templates[k] = v
	}
	return out
}

