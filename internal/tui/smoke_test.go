package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/config"
)

func TestSmoke_AuthModal(t *testing.T) {
	cfg := &config.Config{}
	m := New(cfg, nil)
	_ = m.Init()

	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(Model)
	v := m.View()
	if len(v) == 0 {
		t.Fatal("expected non-empty view")
	}
}

func TestSmoke_HomeView(t *testing.T) {
	cfg := &config.Config{JiraCloudURL: "https://test.atlassian.net", JiraAPIToken: "tok"}
	m := New(cfg, nil)

	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(Model)
	v := m.View()
	if len(v) == 0 {
		t.Fatal("expected non-empty view")
	}
}

func TestSmoke_HelpModal(t *testing.T) {
	cfg := &config.Config{JiraCloudURL: "https://test.atlassian.net", JiraAPIToken: "tok"}
	m := New(cfg, nil)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(Model)

	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = model.(Model)
	v := m.View()
	if len(v) == 0 {
		t.Fatal("expected non-empty view after help")
	}
}

func TestSmoke_ListModal(t *testing.T) {
	cfg := &config.Config{JiraCloudURL: "https://test.atlassian.net", JiraAPIToken: "tok"}
	m := New(cfg, nil)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(Model)

	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")})
	m = model.(Model)
	v := m.View()
	if len(v) == 0 {
		t.Fatal("expected non-empty view after list")
	}
}

func TestSmoke_EscKey(t *testing.T) {
	cfg := &config.Config{JiraCloudURL: "https://test.atlassian.net", JiraAPIToken: "tok"}
	m := New(cfg, nil)
	model, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	m = model.(Model)

	// Open help, then close
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("?")})
	m = model.(Model)
	model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = model.(Model)
	v := m.View()
	if len(v) == 0 {
		t.Fatal("expected non-empty view after esc")
	}
}
