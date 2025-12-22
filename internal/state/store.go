package state

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type StoreFile struct {
	path string
}

func NewStoreFile(path string) *StoreFile {
	return &StoreFile{path: path}
}

func (s *StoreFile) Path() string { return s.path }

func (s *StoreFile) Load() (*Store, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var st Store
	if err := yaml.Unmarshal(data, &st); err != nil {
		return nil, err
	}
	if st.Sessions == nil {
		st.Sessions = make(map[string]*SessionState)
	}
	if st.Version == 0 {
		st.Version = 1
	}
	return &st, nil
}

func (s *StoreFile) Save(st *Store) error {
	if st == nil {
		return errors.New("state store is nil")
	}
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := yaml.Marshal(st)
	if err != nil {
		return err
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	tmp, err := os.CreateTemp(dir, ".agentpane-state-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpName, s.path)
}
