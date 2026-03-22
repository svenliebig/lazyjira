package config

import (
	"testing"
)

func TestLoad_Empty(t *testing.T) {
	// Clear env vars that might be set
	t.Setenv("JIRA_CLOUD_URL", "")
	t.Setenv("JIRA_API_TOKEN", "")
	t.Setenv("XDG_CONFIG_HOME", t.TempDir()) // use temp dir to avoid reading real config

	cfg, err := Load(Flags{})
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg.IsComplete() {
		t.Error("Expected incomplete config, got complete")
	}
}

func TestLoad_EnvVars(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("JIRA_CLOUD_URL", "https://mycompany.atlassian.net")
	t.Setenv("JIRA_EMAIL", "user@mycompany.com")
	t.Setenv("JIRA_API_TOKEN", "my-secret-token")

	cfg, err := Load(Flags{})
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if !cfg.IsComplete() {
		t.Error("Expected complete config")
	}
	if cfg.JiraCloudURL != "https://mycompany.atlassian.net" {
		t.Errorf("Expected JiraCloudURL %q, got %q", "https://mycompany.atlassian.net", cfg.JiraCloudURL)
	}
	if cfg.JiraEmail != "user@mycompany.com" {
		t.Errorf("Expected JiraEmail %q, got %q", "user@mycompany.com", cfg.JiraEmail)
	}
	if cfg.JiraAPIToken != "my-secret-token" {
		t.Errorf("Expected JiraAPIToken %q, got %q", "my-secret-token", cfg.JiraAPIToken)
	}
}

func TestLoad_FlagOverridesEnv(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("JIRA_CLOUD_URL", "https://env.atlassian.net")
	t.Setenv("JIRA_EMAIL", "env@example.com")
	t.Setenv("JIRA_API_TOKEN", "env-token")

	flags := Flags{
		JiraCloudURL: "https://flag.atlassian.net",
		JiraEmail:    "flag@example.com",
		JiraAPIToken: "flag-token",
	}

	cfg, err := Load(flags)
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg.JiraCloudURL != "https://flag.atlassian.net" {
		t.Errorf("Expected flag URL to override env, got %q", cfg.JiraCloudURL)
	}
	if cfg.JiraAPIToken != "flag-token" {
		t.Errorf("Expected flag token to override env, got %q", cfg.JiraAPIToken)
	}
}

func TestLoad_PartialEnv(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("JIRA_CLOUD_URL", "https://mycompany.atlassian.net")
	t.Setenv("JIRA_API_TOKEN", "")

	cfg, err := Load(Flags{})
	if err != nil {
		t.Fatalf("Load() returned unexpected error: %v", err)
	}
	if cfg.IsComplete() {
		t.Error("Expected incomplete config when only URL is set")
	}
}
