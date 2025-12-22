package config

import (
	"fmt"
)

var validPaneTypes = map[string]bool{
	"codex":  true,
	"claude": true,
	"shell":  true,
}

func ValidateGlobal(cfg *Config) error {
	if cfg == nil {
		return nil
	}
	if cfg.DefaultPaneType != "" && !validPaneTypes[cfg.DefaultPaneType] {
		return fmt.Errorf("invalid default_pane_type %q", cfg.DefaultPaneType)
	}
	if cfg.Providers != nil {
		for k := range cfg.Providers {
			if !validPaneTypes[k] {
				return fmt.Errorf("invalid providers key %q (expected codex, claude, shell)", k)
			}
		}
	}
	for name, tmpl := range cfg.Templates {
		if err := validateTemplate(name, tmpl); err != nil {
			return err
		}
	}
	return nil
}

func ValidateRepo(rc *RepoConfig) error {
	if rc == nil {
		return nil
	}
	if len(rc.Layout.Panes) == 0 {
		return fmt.Errorf("repo config layout.panes must not be empty")
	}
	for i, p := range rc.Layout.Panes {
		if !validPaneTypes[p.Type] {
			return fmt.Errorf("repo config layout.panes[%d].type invalid: %q", i, p.Type)
		}
	}
	return nil
}

func validateTemplate(name string, tmpl Template) error {
	if len(tmpl.Panes) == 0 {
		return fmt.Errorf("template %q panes must not be empty", name)
	}
	for i, p := range tmpl.Panes {
		if !validPaneTypes[p.Type] {
			return fmt.Errorf("template %q panes[%d].type invalid: %q", name, i, p.Type)
		}
	}
	return nil
}

