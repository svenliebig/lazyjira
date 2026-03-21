package modals

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/jira-cli/internal/git"
	"github.com/svenliebig/jira-cli/internal/jira"
	"github.com/svenliebig/jira-cli/internal/ollama"
	"github.com/svenliebig/jira-cli/internal/tui/shared"
)

type aiState int

const (
	aiIdle aiState = iota
	aiLoadingCommits
	aiGenerating
	aiDone
	aiError
)

// AIModal shows the AI assistance panel.
type AIModal struct {
	state      aiState
	spinner    spinner.Model
	viewport   viewport.Model
	issue      *jira.Issue
	commits    []string
	summary    string
	errMsg     string
	ollamaClient *ollama.Client
}

func NewAIModal(jiraClient *jira.Client) AIModal {
	_ = jiraClient // jiraClient kept for future use (e.g., fetching issue details)
	s := spinner.New()
	s.Spinner = spinner.Dot

	vp := viewport.New(56, 10)

	return AIModal{
		state:        aiIdle,
		spinner:      s,
		viewport:     vp,
		ollamaClient: ollama.NewClient(),
	}
}

func (m *AIModal) SetIssue(issue *jira.Issue) {
	m.issue = issue
}

func (m *AIModal) SetCommits(commits []string) {
	m.commits = commits
	m.state = aiGenerating
}

func (m *AIModal) SetSummary(summary string) {
	m.summary = summary
	m.state = aiDone
	m.viewport.SetContent(summary)
}

func (m AIModal) Init() tea.Cmd {
	return nil
}

func (m AIModal) GenerateCmd() tea.Cmd {
	commits := m.commits
	issue := m.issue
	client := m.ollamaClient

	return func() tea.Msg {
		var prompt string
		if issue != nil {
			commitList := strings.Join(commits, "\n")
			if len(commitList) == 0 {
				commitList = "(no commits found for this issue)"
			}
			prompt = fmt.Sprintf(
				"Summarize the following Jira issue and related git commits in 2-3 sentences.\n\nIssue: %s\nSummary: %s\nDescription: %s\n\nCommits:\n%s",
				issue.Key,
				issue.Fields.Summary,
				issue.Fields.Description,
				commitList,
			)
		} else {
			commitList := strings.Join(commits, "\n")
			prompt = fmt.Sprintf("Summarize these git commits:\n%s", commitList)
		}

		result, err := client.Generate(context.Background(), prompt)
		if err != nil {
			return shared.ErrMsg{Err: fmt.Errorf("AI generation failed: %w", err)}
		}
		return shared.AISummaryMsg{Summary: result}
	}
}

func (m AIModal) fetchCommitsCmd() tea.Cmd {
	issue := m.issue
	return func() tea.Msg {
		var issueKey string
		if issue != nil {
			issueKey = issue.Key
		}
		commits, err := git.CommitsForIssue(issueKey)
		if err != nil {
			// Not a git repo — return empty commits rather than error
			return shared.AICommitsLoadedMsg{Commits: []string{}}
		}
		return shared.AICommitsLoadedMsg{Commits: commits}
	}
}

func (m AIModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, func() tea.Msg { return shared.CloseModalMsg{} }
		case "s":
			if m.state == aiIdle {
				m.state = aiLoadingCommits
				return m, tea.Batch(m.spinner.Tick, m.fetchCommitsCmd())
			}
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	if m.state == aiDone {
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m AIModal) View() string {
	var content string

	switch m.state {
	case aiIdle:
		content = shared.StyleMuted.Render("Press s to generate an AI summary for this issue.\n\n") +
			shared.StyleMuted.Render("esc: close")
	case aiLoadingCommits:
		content = m.spinner.View() + " " + shared.StyleMuted.Render("Loading git commits...")
	case aiGenerating:
		content = m.spinner.View() + " " + shared.StyleMuted.Render("Generating AI summary...")
	case aiDone:
		content = m.viewport.View() + "\n\n" + shared.StyleMuted.Render("esc: close")
	case aiError:
		content = shared.StyleError.Render("Error: "+m.errMsg) + "\n\n" + shared.StyleMuted.Render("esc: close")
	}

	return Wrap("AI Assistance", content)
}
