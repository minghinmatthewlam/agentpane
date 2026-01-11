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
	globalLoaded, err := loadYAMLFile(globalPath, func(data []byte) error {
		return yaml.Unmarshal(data, &globalCfg)
	})
	if err != nil {
		return nil, err
	}

	var repoCfg RepoConfig
	repoLoaded := false
	repoPath, found, err := FindRepoConfigPath(cwd)
	if err != nil {
		return nil, err
	}
	if found {
		repoLoaded, err = loadYAMLFile(repoPath, func(data []byte) error {
			return yaml.Unmarshal(data, &repoCfg)
		})
		if err != nil {
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

func loadYAMLFile(path string, decode func([]byte) error) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if err := decode(data); err != nil {
		return false, err
	}
	return true, nil
}
