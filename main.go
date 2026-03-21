package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/config"
	"github.com/svenliebig/jira-cli/internal/jira"
	"github.com/svenliebig/jira-cli/internal/tui"
)

func main() {
	var flags config.Flags
	flag.StringVar(&flags.JiraCloudURL, "jira-cloud-url", "", "Jira Cloud URL")
	flag.StringVar(&flags.JiraAPIToken, "jira-api-token", "", "Jira API Token")
	flag.Parse()

	cfg, err := config.Load(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	var jiraClient *jira.Client
	if cfg.IsComplete() {
		jiraClient = jira.NewClient(cfg.JiraCloudURL, cfg.JiraAPIToken)
	}

	model := tui.New(cfg, jiraClient)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
