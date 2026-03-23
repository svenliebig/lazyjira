package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/x/ansi"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/svenliebig/lazyjira/internal/browser"
	"github.com/svenliebig/lazyjira/internal/clipboard"
	"github.com/svenliebig/lazyjira/internal/config"
	"github.com/svenliebig/lazyjira/internal/exclusions"
	"github.com/svenliebig/lazyjira/internal/jira"
	"github.com/svenliebig/lazyjira/internal/tui/modals"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
	"github.com/svenliebig/lazyjira/internal/tui/views"
)

type viewState int

const (
	viewHome viewState = iota
	viewIssueList
	viewIssueDetail
	viewExcludedList
)

type modalState int

const (
	modalNone modalState = iota
	modalAuth
	modalHelp
	modalListSelector
	modalCopy
	modalAI
	modalTransition
	modalExclude
)

// Model is the root bubbletea model.
type Model struct {
	width  int
	height int
	cfg    *config.Config

	jiraClient  *jira.Client
	currentView viewState
	activeModal modalState
	pendingKey  string

	// View models
	homeView         views.HomeModel
	issueListView    views.IssueListModel
	issueDetailView  views.IssueDetailModel
	excludedListView views.ExcludedListModel

	// Modal models
	authModal       modals.AuthModal
	helpModal       modals.HelpModal
	listModal       modals.ListSelectorModal
	copyModal       modals.CopyModal
	aiModal         modals.AIModal
	transitionModal modals.TransitionModal
	excludeModal    modals.ExcludeModal

	// State
	currentIssue *jira.Issue
	allIssues    []jira.Issue
	exclusions   *exclusions.Store
	loading      bool
	err          string
	statusMsg    string
}

func New(cfg *config.Config, jiraClient *jira.Client, store *exclusions.Store) Model {
	m := Model{
		cfg:        cfg,
		jiraClient: jiraClient,
		exclusions: store,
	}
	if !cfg.IsComplete() {
		m.activeModal = modalAuth
		m.authModal = modals.NewAuthModal()
	}
	m.helpModal = modals.NewHelpModal()
	m.listModal = modals.NewListSelectorModal()
	return m
}

func (m Model) Init() tea.Cmd {
	if m.activeModal == modalAuth {
		return m.authModal.Init()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateChildSizes()
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)

	case shared.IssueListLoadedMsg:
		m.loading = false
		m.allIssues = msg.Issues
		filtered := m.exclusions.Filter(m.allIssues)
		m.issueListView = views.NewIssueListModel(filtered, m.width, m.height-2)
		m.currentView = viewIssueList
		m.currentIssue = m.issueListView.CurrentIssue()
		return m, nil

	case shared.IssueSelectedMsg:
		m.currentIssue = &msg.Issue
		m.issueDetailView = views.NewIssueDetailModel(msg.Issue, m.width, m.height-4)
		m.currentView = viewIssueDetail
		return m, nil

	case shared.TransitionsLoadedMsg:
		m.transitionModal = modals.NewTransitionModal(msg.Transitions)
		m.activeModal = modalTransition
		m.loading = false
		return m, nil

	case shared.TransitionSelectedMsg:
		if m.currentIssue != nil && m.jiraClient != nil {
			m.activeModal = modalNone
			m.loading = true
			return m, m.doTransitionCmd(msg.ID)
		}
		return m, nil

	case shared.TransitionDoneMsg:
		m.activeModal = modalNone
		m.loading = false
		m.statusMsg = "Transition applied"
		return m, nil

	case shared.ListSelectedMsg:
		m.activeModal = modalNone
		switch msg.Type {
		case "excluded":
			m.excludedListView = views.NewExcludedListModel(m.exclusions.Rules(), m.width, m.height-2)
			m.currentView = viewExcludedList
			m.currentIssue = nil
		default: // "assigned"
			if m.jiraClient != nil {
				m.loading = true
				return m, fetchAssignedCmd(m.jiraClient)
			}
			m.err = "Jira client not configured"
		}
		return m, nil

	case shared.AuthCompletedMsg:
		m.cfg.JiraCloudURL = msg.URL
		m.cfg.JiraEmail = msg.Email
		m.cfg.JiraAPIToken = msg.Token
		_ = config.Save(m.cfg)
		m.jiraClient = jira.NewClient(msg.URL, msg.Email, msg.Token)
		m.activeModal = modalNone
		m.statusMsg = "Authenticated!"
		return m, nil

	case shared.AICommitsLoadedMsg:
		m.aiModal.SetCommits(msg.Commits)
		return m, m.aiModal.GenerateCmd()

	case shared.AISummaryMsg:
		m.aiModal.SetSummary(msg.Summary)
		return m, nil

	case shared.CopyActionMsg:
		if m.currentIssue != nil {
			var text string
			switch msg.Action {
			case "key":
				text = m.currentIssue.Key
			case "url":
				text = m.cfg.JiraCloudURL + "/browse/" + m.currentIssue.Key
			case "title":
				text = m.currentIssue.Fields.Summary
			case "desc":
				text = m.currentIssue.Fields.Description
			}
			if text != "" {
				_ = clipboard.Write(text)
				m.activeModal = modalNone
				m.pendingKey = ""
				m.statusMsg = "Copied!"
			}
		}
		return m, nil

	case shared.CopyMsg:
		_ = clipboard.Write(msg.Text)
		m.activeModal = modalNone
		m.pendingKey = ""
		m.statusMsg = "Copied!"
		return m, nil

	case shared.CloseModalMsg:
		m.activeModal = modalNone
		m.pendingKey = ""
		return m, nil

	case shared.ExcludeActionMsg:
		if m.exclusions != nil {
			_ = m.exclusions.Add(exclusions.Rule{Type: msg.Type, Value: msg.Value})
			filtered := m.exclusions.Filter(m.allIssues)
			m.issueListView = views.NewIssueListModel(filtered, m.width, m.height-2)
			m.currentIssue = m.issueListView.CurrentIssue()
		}
		m.activeModal = modalNone
		m.statusMsg = "Issue excluded"
		return m, nil

	case shared.ErrMsg:
		m.loading = false
		m.err = msg.Err.Error()
		return m, nil
	}

	return m.updateActiveChild(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Always allow quit
	if key == "ctrl+c" {
		return m, tea.Quit
	}

	// If modal active, delegate to modal
	if m.activeModal != modalNone {
		return m.updateActiveChild(msg)
	}

	// ESC clears pending key or goes back
	if key == shared.KeyEsc {
		if m.pendingKey != "" {
			m.pendingKey = ""
			return m, nil
		}
		if m.currentView == viewIssueList && m.issueListView.IsFocusRight() {
			m.issueListView = m.issueListView.BlurRight()
			return m, nil
		}
		if m.currentView == viewIssueDetail {
			m.currentView = viewIssueList
			m.currentIssue = m.issueListView.CurrentIssue()
			return m, nil
		}
		if m.currentView == viewIssueList {
			m.currentView = viewHome
			m.currentIssue = nil
			return m, nil
		}
		if m.currentView == viewExcludedList {
			m.currentView = viewHome
			return m, nil
		}
		return m, nil
	}

	// Chord resolution for copy
	if m.pendingKey == shared.KeyCopy {
		m.pendingKey = ""
		if m.currentIssue != nil {
			switch key {
			case "k":
				text := m.currentIssue.Key
				return m, func() tea.Msg { return shared.CopyMsg{Text: text} }
			case "u":
				url := m.cfg.JiraCloudURL + "/browse/" + m.currentIssue.Key
				return m, func() tea.Msg { return shared.CopyMsg{Text: url} }
			case "t":
				text := m.currentIssue.Fields.Summary
				return m, func() tea.Msg { return shared.CopyMsg{Text: text} }
			case "d":
				text := m.currentIssue.Fields.Description
				return m, func() tea.Msg { return shared.CopyMsg{Text: text} }
			}
		}
		return m, nil
	}

	// Chord resolution for AI
	if m.pendingKey == shared.KeyAI {
		m.pendingKey = ""
		if key == "s" && m.currentIssue != nil {
			m.aiModal = modals.NewAIModal(m.jiraClient)
			m.aiModal.SetIssue(m.currentIssue)
			m.activeModal = modalAI
			return m, m.aiModal.Init()
		}
		return m, nil
	}

	// Global keys
	switch key {
	case shared.KeyHelp:
		m.activeModal = modalHelp
		return m, nil

	case shared.KeyQuit:
		if m.currentView == viewHome {
			return m, tea.Quit
		}

	case shared.KeyList:
		m.activeModal = modalListSelector
		return m, m.listModal.Init()

	case shared.KeyCopy:
		if m.currentIssue != nil {
			m.pendingKey = shared.KeyCopy
			m.copyModal = modals.NewCopyModal()
			m.activeModal = modalCopy
			return m, nil
		}

	case shared.KeyOpen:
		if m.currentIssue != nil {
			url := m.cfg.JiraCloudURL + "/browse/" + m.currentIssue.Key
			_ = browser.OpenURL(url)
			m.statusMsg = "Opening in browser..."
		}

	case shared.KeyAI:
		if m.currentIssue != nil {
			m.pendingKey = shared.KeyAI
			m.aiModal = modals.NewAIModal(m.jiraClient)
			m.aiModal.SetIssue(m.currentIssue)
			m.activeModal = modalAI
			return m, nil
		}

	case shared.KeyTransition:
		if m.currentIssue != nil && m.jiraClient != nil {
			m.loading = true
			return m, m.fetchTransitionsCmd()
		}

	case shared.KeyExclude:
		if m.currentView == viewExcludedList {
			if rule := m.excludedListView.CurrentRule(); rule != nil {
				_ = m.exclusions.Remove(*rule)
				m.excludedListView = views.NewExcludedListModel(m.exclusions.Rules(), m.width, m.height-2)
				m.statusMsg = "Exclusion removed"
			}
			return m, nil
		}
		if m.currentIssue != nil {
			m.excludeModal = modals.NewExcludeModal(m.currentIssue)
			m.activeModal = modalExclude
			return m, nil
		}
	}

	// Delegate navigation to active view
	return m.updateActiveChild(msg)
}

func (m Model) updateActiveChild(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.activeModal {
	case modalAuth:
		var updated tea.Model
		updated, cmd = m.authModal.Update(msg)
		m.authModal = updated.(modals.AuthModal)
	case modalHelp:
		var updated tea.Model
		updated, cmd = m.helpModal.Update(msg)
		m.helpModal = updated.(modals.HelpModal)
	case modalListSelector:
		var updated tea.Model
		updated, cmd = m.listModal.Update(msg)
		m.listModal = updated.(modals.ListSelectorModal)
	case modalCopy:
		var updated tea.Model
		updated, cmd = m.copyModal.Update(msg)
		m.copyModal = updated.(modals.CopyModal)
	case modalAI:
		var updated tea.Model
		updated, cmd = m.aiModal.Update(msg)
		m.aiModal = updated.(modals.AIModal)
	case modalTransition:
		var updated tea.Model
		updated, cmd = m.transitionModal.Update(msg)
		m.transitionModal = updated.(modals.TransitionModal)
	case modalExclude:
		var updated tea.Model
		updated, cmd = m.excludeModal.Update(msg)
		m.excludeModal = updated.(modals.ExcludeModal)
	default:
		switch m.currentView {
		case viewIssueList:
			var updated tea.Model
			updated, cmd = m.issueListView.Update(msg)
			m.issueListView = updated.(views.IssueListModel)
			m.currentIssue = m.issueListView.CurrentIssue()
		case viewIssueDetail:
			var updated tea.Model
			updated, cmd = m.issueDetailView.Update(msg)
			m.issueDetailView = updated.(views.IssueDetailModel)
		case viewExcludedList:
			var updated tea.Model
			updated, cmd = m.excludedListView.Update(msg)
			m.excludedListView = updated.(views.ExcludedListModel)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	header := m.renderHeader()
	content := m.renderContent()
	statusBar := m.renderStatusBar()

	view := lipgloss.JoinVertical(lipgloss.Left, header, content, statusBar)

	// Overlay modal if active
	if m.activeModal != modalNone {
		overlay := m.renderModal()
		if overlay != "" {
			view = overlayModal(view, overlay, m.width, m.height)
		}
	}

	return view
}

func (m Model) renderHeader() string {
	title := "lazyjira"
	if m.currentIssue != nil {
		title += "  " + shared.StyleIssueKey.Render(m.currentIssue.Key)
	}
	return shared.StyleHeader.Width(m.width).Render(title)
}

func (m Model) renderContent() string {
	contentHeight := m.height - 2 // header + statusbar
	if contentHeight < 0 {
		contentHeight = 1
	}
	style := lipgloss.NewStyle().Width(m.width).Height(contentHeight)

	if m.loading {
		return style.Render(shared.StyleMuted.Render("\n  Loading..."))
	}
	if m.err != "" {
		return style.Render(shared.StyleError.Render("\n  Error: " + m.err))
	}

	switch m.currentView {
	case viewHome:
		return style.Render(m.homeView.View())
	case viewIssueList:
		// The split view manages its own layout — no extra wrapping
		return m.issueListView.View()
	case viewIssueDetail:
		return style.Render(m.issueDetailView.View())
	case viewExcludedList:
		return style.Render(m.excludedListView.View())
	}
	return style.Render("")
}

func (m Model) renderStatusBar() string {
	var hints []string
	if m.pendingKey == shared.KeyCopy {
		hints = []string{"↑/k↓/j:navigate", "enter/l:select", "u:url", "t:title", "d:desc", "h/esc:cancel"}
	} else if m.pendingKey == shared.KeyAI {
		hints = []string{"s:summary", "esc:cancel"}
	} else if m.currentView == viewExcludedList {
		hints = []string{"j/k:navigate", "x:remove", "esc:back", "?:help"}
	} else if m.currentView == viewIssueList && m.issueListView.IsFocusRight() {
		hints = []string{"j/k:scroll", "esc:back"}
		if m.currentIssue != nil {
			hints = append(hints, "o:open", "y:copy", "t:transition", "a:AI", "x:exclude")
		}
	} else {
		hints = []string{"l:list", "?:help", "q:quit"}
		if m.currentView == viewIssueList {
			hints = append(hints, "enter:focus detail")
		}
		if m.currentIssue != nil {
			hints = append(hints, "o:open", "y:copy", "t:transition", "a:AI", "x:exclude")
		}
	}

	parts := make([]string, 0, len(hints)*2)
	for i, h := range hints {
		parts = append(parts, shared.StyleKeyHint.Render(h))
		if i < len(hints)-1 {
			parts = append(parts, shared.StyleKeyHintSep.Render("  "))
		}
	}
	if m.statusMsg != "" {
		parts = append(parts, "  "+shared.StyleSuccess.Render(m.statusMsg))
	}

	content := strings.Join(parts, "")
	// Truncate to content area width (m.width minus padding) to prevent wrapping to a second line,
	// which would push the total rendered height over the terminal height and scroll the header off.
	if m.width > 2 {
		content = ansi.Truncate(content, m.width-2, "")
	}
	return shared.StyleStatusBar.Width(m.width).Render(content)
}

func (m Model) renderModal() string {
	switch m.activeModal {
	case modalAuth:
		return m.authModal.View()
	case modalHelp:
		return m.helpModal.View()
	case modalListSelector:
		return m.listModal.View()
	case modalCopy:
		return m.copyModal.View()
	case modalAI:
		return m.aiModal.View()
	case modalTransition:
		return m.transitionModal.View()
	case modalExclude:
		return m.excludeModal.View()
	}
	return ""
}

func overlayModal(background, modal string, w, h int) string {
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, modal,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#00000080")))
}

func (m *Model) updateChildSizes() {
	contentH := m.height - 2
	m.issueListView.SetSize(m.width, contentH)
	m.issueDetailView.SetSize(m.width, contentH)
	m.excludedListView.SetSize(m.width, contentH)
}

func (m Model) fetchTransitionsCmd() tea.Cmd {
	client := m.jiraClient
	key := m.currentIssue.Key
	return func() tea.Msg {
		transitions, err := client.GetTransitions(context.Background(), key)
		if err != nil {
			return shared.ErrMsg{Err: err}
		}
		return shared.TransitionsLoadedMsg{Transitions: transitions}
	}
}

func (m Model) doTransitionCmd(transitionID string) tea.Cmd {
	client := m.jiraClient
	key := m.currentIssue.Key
	return func() tea.Msg {
		if err := client.DoTransition(context.Background(), key, transitionID); err != nil {
			return shared.ErrMsg{Err: fmt.Errorf("transition failed: %w", err)}
		}
		return shared.TransitionDoneMsg{}
	}
}

func fetchAssignedCmd(client *jira.Client) tea.Cmd {
	return func() tea.Msg {
		issues, err := client.ListAssigned(context.Background())
		if err != nil {
			return shared.ErrMsg{Err: fmt.Errorf("failed to load issues: %w", err)}
		}
		return shared.IssueListLoadedMsg{Issues: issues}
	}
}
