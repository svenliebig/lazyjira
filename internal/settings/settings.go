package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Settings holds persisted user preferences.
type Settings struct {
	ActiveTheme string `json:"activeTheme"`
}

// Load reads settings from disk. If the file doesn't exist the defaults are returned.
func Load() (*Settings, error) {
	s := &Settings{ActiveTheme: "default"}
	path, err := settingsFilePath()
	if err != nil {
		return s, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return s, nil
	}
	_ = json.Unmarshal(data, s)
	return s, nil
}

// Save persists settings to disk.
func Save(s *Settings) error {
	path, err := settingsFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func settingsFilePath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "lazyjira", "settings.json"), nil
}
