package views

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/svenliebig/lazyjira/internal/tui/shared"
)

type boardEntry struct {
	projectKey string
	boardID    int
}

// BoardsModel renders the Boards tab for managing project → board mappings.
type BoardsModel struct {
	entries   []boardEntry
	cursor    int
	width     int
	height    int
	forming   bool
	editIdx   int // -1 = new entry
	formKey   string
	formID    string
	formFocus int // 0 = project key field, 1 = board ID field
	formErr   string
}

func NewBoardsModel(boards map[string]int, width, height int) BoardsModel {
	entries := make([]boardEntry, 0, len(boards))
	for k, v := range boards {
		entries = append(entries, boardEntry{projectKey: k, boardID: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].projectKey < entries[j].projectKey
	})
	return BoardsModel{
		entries: entries,
		width:   width,
		height:  height,
		editIdx: -1,
	}
}

func (m *BoardsModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m BoardsModel) Init() tea.Cmd { return nil }

func (m BoardsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.forming {
			return m.handleFormKey(msg)
		}
		return m.handleListKey(msg)
	}
	return m, nil
}

func (m BoardsModel) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.entries)-1 {
			m.cursor++
		}
	case "a":
		m.forming = true
		m.editIdx = -1
		m.formKey = ""
		m.formID = ""
		m.formFocus = 0
		m.formErr = ""
	case "e":
		if m.cursor < len(m.entries) {
			e := m.entries[m.cursor]
			m.forming = true
			m.editIdx = m.cursor
			m.formKey = e.projectKey
			m.formID = fmt.Sprintf("%d", e.boardID)
			m.formFocus = 0
			m.formErr = ""
		}
	case "d":
		if m.cursor < len(m.entries) {
			key := m.entries[m.cursor].projectKey
			return m, func() tea.Msg { return shared.BoardDeletedMsg{ProjectKey: key} }
		}
	}
	return m, nil
}

func (m BoardsModel) handleFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.forming = false
		m.formErr = ""
	case "tab":
		m.formFocus = 1 - m.formFocus
	case "enter":
		projectKey := strings.TrimSpace(strings.ToUpper(m.formKey))
		if projectKey == "" {
			m.formErr = "Project key is required"
			return m, nil
		}
		boardID, err := strconv.Atoi(strings.TrimSpace(m.formID))
		if err != nil || boardID <= 0 {
			m.formErr = "Board ID must be a positive number"
			return m, nil
		}
		m.forming = false
		m.formErr = ""
		return m, func() tea.Msg {
			return shared.BoardSavedMsg{ProjectKey: projectKey, BoardID: boardID}
		}
	case "backspace":
		if m.formFocus == 0 && len(m.formKey) > 0 {
			m.formKey = m.formKey[:len(m.formKey)-1]
		} else if m.formFocus == 1 && len(m.formID) > 0 {
			m.formID = m.formID[:len(m.formID)-1]
		}
	default:
		if len(msg.String()) == 1 {
			if m.formFocus == 0 {
				m.formKey += msg.String()
			} else {
				m.formID += msg.String()
			}
		}
	}
	return m, nil
}

func (m BoardsModel) View() string {
	if m.forming {
		return m.renderForm()
	}
	return m.renderList()
}

func (m BoardsModel) renderList() string {
	var sb strings.Builder
	if len(m.entries) == 0 {
		sb.WriteString(shared.StyleMuted.Render("  No boards configured. Press a to add one.\n"))
	} else {
		for i, e := range m.entries {
			prefix := "  "
			keyStyle := shared.StyleNormalItem
			if i == m.cursor {
				prefix = shared.StyleSelectedItem.Render(">") + " "
				keyStyle = shared.StyleSelectedItem
			}
			line := keyStyle.Render(e.projectKey) + "  " + shared.StyleMuted.Render(fmt.Sprintf("Board %d", e.boardID))
			sb.WriteString(prefix + line + "\n")
		}
	}
	sb.WriteString("\n" + shared.StyleMuted.Render("  a:add  e:edit  d:delete  j/k:navigate"))
	return shared.StyleContentArea.Render(sb.String())
}

func (m BoardsModel) renderForm() string {
	title := "Add Board"
	if m.editIdx >= 0 {
		title = "Edit Board"
	}

	var sb strings.Builder
	sb.WriteString(shared.StyleModalTitle.Render(title) + "\n\n")

	projLabel := "  Project Key: "
	if m.formFocus == 0 {
		sb.WriteString(shared.StyleSelectedItem.Render(projLabel) + shared.StyleNormalItem.Render(m.formKey+"_") + "\n")
	} else {
		sb.WriteString(shared.StyleMuted.Render(projLabel) + shared.StyleNormalItem.Render(m.formKey) + "\n")
	}

	boardLabel := "  Board ID:    "
	if m.formFocus == 1 {
		sb.WriteString(shared.StyleSelectedItem.Render(boardLabel) + shared.StyleNormalItem.Render(m.formID+"_") + "\n")
	} else {
		sb.WriteString(shared.StyleMuted.Render(boardLabel) + shared.StyleNormalItem.Render(m.formID) + "\n")
	}

	if m.formErr != "" {
		sb.WriteString("\n" + shared.StyleError.Render("  "+m.formErr) + "\n")
	}

	sb.WriteString("\n" + shared.StyleMuted.Render("  tab: next field   enter: save   esc: cancel"))

	return shared.StyleContentArea.Render(sb.String())
}
