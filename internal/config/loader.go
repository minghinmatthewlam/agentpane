package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Loaded struct {
	Global     *Config
	Repo       *RepoConfig
	Builtins   map[string]Template
	Merged     *Config
	RepoPath   string
	GlobalPath string
}

func LoadAll(cwd string) (*Loaded, error) {
	builtins, err := LoadBuiltinTemplates()
	if err != nil {
		return nil, err
	}

	globalPath, err := GlobalConfigPath()
	if err != nil {
		return nil, err
	}

	var globalCfg Config
	globalLoaded := false
	if data, err := os.ReadFile(globalPath); err == nil {
		if err := yaml.Unmarshal(data, &globalCfg); err != nil {
			return nil, err
		}
		globalLoaded = true
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var repoCfg RepoConfig
	repoLoaded := false
	repoPath, found, err := FindRepoConfigPath(cwd)
	if err != nil {
		return nil, err
	}
	if found {
		if data, err := os.ReadFile(repoPath); err == nil {
			if err := yaml.Unmarshal(data, &repoCfg); err != nil {
				return nil, err
			}
			repoLoaded = true
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	base := DefaultConfig()
	// Builtins are the default template set; global config can override by name.
	base.Templates = builtins

	var globalPtr *Config
	if globalLoaded {
		globalPtr = &globalCfg
	}
	if err := ValidateGlobal(globalPtr); err != nil {
		return nil, err
	}

	var repoPtr *RepoConfig
	if repoLoaded {
		repoPtr = &repoCfg
	}
	if err := ValidateRepo(repoPtr); err != nil {
		return nil, err
	}

	merged := Merge(base, globalPtr)

	return &Loaded{
		Global:     globalPtr,
		Repo:       repoPtr,
		Builtins:   builtins,
		Merged:     merged,
		RepoPath:   repoPath,
		GlobalPath: globalPath,
	}, nil
}
