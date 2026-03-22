package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/config"
	"github.com/svenliebig/lazyjira/internal/exclusions"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui"
)

var version = "dev"

func main() {
	var flags config.Flags
	flag.StringVar(&flags.JiraCloudURL, "jira-cloud-url", "", "Jira Cloud URL")
	flag.StringVar(&flags.JiraEmail, "jira-email", "", "Jira account email")
	flag.StringVar(&flags.JiraAPIToken, "jira-api-token", "", "Jira API Token")
	flag.Parse()

	cfg, err := config.Load(flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	var jiraClient *jira.Client
	if cfg.IsComplete() {
		jiraClient = jira.NewClient(cfg.JiraCloudURL, cfg.JiraEmail, cfg.JiraAPIToken)
	}

	store, err := exclusions.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load exclusions: %v\n", err)
	}

	model := tui.New(cfg, jiraClient, store)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
