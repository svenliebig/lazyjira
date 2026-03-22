package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	JiraCloudURL string `json:"jiraCloudUrl"`
	JiraEmail    string `json:"jiraEmail"`
	JiraAPIToken string `json:"jiraApiToken"`
}

type Flags struct {
	JiraCloudURL string
	JiraEmail    string
	JiraAPIToken string
}

func Load(flags Flags) (*Config, error) {
	cfg := &Config{}

	// 1. Config file
	if path, err := configFilePath(); err == nil {
		if data, err := os.ReadFile(path); err == nil {
			_ = json.Unmarshal(data, cfg)
		}
	}

	// 2. Environment variables (override file)
	if v := os.Getenv("JIRA_CLOUD_URL"); v != "" {
		cfg.JiraCloudURL = v
	}
	if v := os.Getenv("JIRA_EMAIL"); v != "" {
		cfg.JiraEmail = v
	}
	if v := os.Getenv("JIRA_API_TOKEN"); v != "" {
		cfg.JiraAPIToken = v
	}

	// 3. CLI flags (highest priority)
	if flags.JiraCloudURL != "" {
		cfg.JiraCloudURL = flags.JiraCloudURL
	}
	if flags.JiraEmail != "" {
		cfg.JiraEmail = flags.JiraEmail
	}
	if flags.JiraAPIToken != "" {
		cfg.JiraAPIToken = flags.JiraAPIToken
	}

	return cfg, nil
}

func (c *Config) IsComplete() bool {
	return c.JiraCloudURL != "" && c.JiraEmail != "" && c.JiraAPIToken != ""
}

func Save(cfg *Config) error {
	path, err := configFilePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func configFilePath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "lazyjira", "config.json"), nil
}
