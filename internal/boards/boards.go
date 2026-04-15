package boards

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Store holds project key → board ID mappings and persists them to disk.
type Store struct {
	data map[string]int
	path string
}

// Load reads the board store from disk. Returns an empty store if the file does
// not exist.
func Load() (*Store, error) {
	p := storePath()
	s := &Store{data: make(map[string]int), path: p}

	raw, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return s, err
	}

	_ = json.Unmarshal(raw, &s.data)
	return s, nil
}

// Get returns the board ID for a project key and whether it was found.
func (s *Store) Get(projectKey string) (int, bool) {
	id, ok := s.data[projectKey]
	return id, ok
}

// Set adds or updates a project → board mapping and persists it.
func (s *Store) Set(projectKey string, boardID int) error {
	s.data[projectKey] = boardID
	return s.save()
}

// Delete removes a mapping and persists the result.
func (s *Store) Delete(projectKey string) error {
	delete(s.data, projectKey)
	return s.save()
}

// All returns a copy of all mappings.
func (s *Store) All() map[string]int {
	out := make(map[string]int, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

func storePath() string {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "lazyjira", "boards.json")
}
