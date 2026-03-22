package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/exclusions"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

// exclusionItem wraps an exclusions.Rule to implement list.Item.
type exclusionItem struct {
	rule exclusions.Rule
}

func (i exclusionItem) Title() string {
	switch i.rule.Type {
	case "key":
		return shared.StyleIssueKey.Render(i.rule.Value)
	case "parent":
		return "Parent: " + shared.StyleIssueKey.Render(i.rule.Value)
	}
	return i.rule.Value
}

func (i exclusionItem) Description() string {
	switch i.rule.Type {
	case "key":
		return "Excluded by issue key"
	case "parent":
		return "All issues with this parent are excluded"
	}
	return ""
}

func (i exclusionItem) FilterValue() string { return i.rule.Value }

// ExcludedListModel shows the active exclusion rules.
type ExcludedListModel struct {
	list        list.Model
	rules       []exclusions.Rule
	initialized bool
	width       int
	height      int
}

func NewExcludedListModel(rules []exclusions.Rule, width, height int) ExcludedListModel {
	items := make([]list.Item, len(rules))
	for i, r := range rules {
		items[i] = exclusionItem{rule: r}
	}

	delegate := list.NewDefaultDelegate()
	l := list.New(items, delegate, width, height)
	l.Title = "Excluded Issues"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return ExcludedListModel{
		list:        l,
		rules:       rules,
		initialized: true,
		width:       width,
		height:      height,
	}
}

// CurrentRule returns the highlighted rule, or nil if the list is empty.
func (m ExcludedListModel) CurrentRule() *exclusions.Rule {
	if len(m.rules) == 0 {
		return nil
	}
	idx := m.list.Index()
	if idx < 0 || idx >= len(m.rules) {
		return nil
	}
	r := m.rules[idx]
	return &r
}

func (m *ExcludedListModel) SetSize(w, h int) {
	if !m.initialized {
		return
	}
	m.width = w
	m.height = h
	m.list.SetSize(w, h)
}

func (m ExcludedListModel) Init() tea.Cmd { return nil }

func (m ExcludedListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ExcludedListModel) View() string {
	if len(m.rules) == 0 {
		return shared.StyleMuted.Render("\n  No exclusions configured.")
	}
	return m.list.View()
}
