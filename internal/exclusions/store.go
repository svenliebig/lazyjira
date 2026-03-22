package exclusions

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/svenliebig/lazyjira/internal/jira"
)

// Rule represents a single exclusion criterion.
type Rule struct {
	Type  string `json:"type"`  // "key" or "parent"
	Value string `json:"value"` // issue key or parent key
}

// Store holds the set of active exclusion rules and persists them to disk.
type Store struct {
	rules []Rule
	path  string
}

// Load reads the exclusion store from disk. If the file does not exist, an
// empty store is returned. The returned *Store is always non-nil.
func Load() (*Store, error) {
	p := storePath()
	s := &Store{path: p}

	data, err := os.ReadFile(p)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return s, err
	}

	if err := json.Unmarshal(data, &s.rules); err != nil {
		return s, err
	}
	return s, nil
}

// Add adds a rule to the store and persists it. Duplicate rules are ignored.
func (s *Store) Add(r Rule) error {
	for _, existing := range s.rules {
		if existing == r {
			return nil
		}
	}
	s.rules = append(s.rules, r)
	return s.save()
}

// Remove removes a rule from the store and persists the result.
func (s *Store) Remove(r Rule) error {
	filtered := make([]Rule, 0, len(s.rules))
	for _, existing := range s.rules {
		if existing != r {
			filtered = append(filtered, existing)
		}
	}
	s.rules = filtered
	return s.save()
}

// Rules returns a copy of all active exclusion rules.
func (s *Store) Rules() []Rule {
	out := make([]Rule, len(s.rules))
	copy(out, s.rules)
	return out
}

// Filter removes issues that match any active exclusion rule.
func (s *Store) Filter(issues []jira.Issue) []jira.Issue {
	if len(s.rules) == 0 {
		return issues
	}

	excludedKeys := make(map[string]bool, len(s.rules))
	excludedParents := make(map[string]bool, len(s.rules))
	for _, r := range s.rules {
		switch r.Type {
		case "key":
			excludedKeys[r.Value] = true
		case "parent":
			excludedParents[r.Value] = true
		}
	}

	result := make([]jira.Issue, 0, len(issues))
	for _, issue := range issues {
		if excludedKeys[issue.Key] {
			continue
		}
		if issue.Fields.Parent != nil && excludedParents[issue.Fields.Parent.Key] {
			continue
		}
		result = append(result, issue)
	}
	return result
}

func (s *Store) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s.rules, "", "  ")
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
	return filepath.Join(base, "lazyjira", "exclusions.json")
}
