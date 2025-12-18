package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"os"
)

// SelectorItem represents an item in the selector
type SelectorItem struct {
	Title       string
	Description string
	Value       string
}

// SelectorModel is the bubbletea model for interactive selection
type SelectorModel struct {
	items    []SelectorItem
	cursor   int
	selected int
	title    string
	done     bool
	canceled bool
	width    int
	height   int
}

// NewSelector creates a new selector model
func NewSelector(title string, items []SelectorItem) SelectorModel {
	return SelectorModel{
		items:    items,
		cursor:   0,
		selected: -1,
		title:    title,
	}
}

func (m SelectorModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, tea.EnterAltScreen)
}

func (m SelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.canceled = true
			m.done = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			} else {
				m.cursor = len(m.items) - 1 // Wrap to bottom
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			} else {
				m.cursor = 0 // Wrap to top
			}

		case "pgup", "ctrl+u":
			// Jump up 5 items
			m.cursor -= 5
			if m.cursor < 0 {
				m.cursor = 0
			}

		case "pgdown", "ctrl+d":
			// Jump down 5 items
			m.cursor += 5
			if m.cursor >= len(m.items) {
				m.cursor = len(m.items) - 1
			}

		case "home", "g":
			m.cursor = 0

		case "end", "G":
			m.cursor = len(m.items) - 1

		case "enter", " ", "l":
			m.selected = m.cursor
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}


func (m SelectorModel) View() string {
	if m.done {
		return ""
	}

	// Use stored dimensions or get from terminal
	width := m.width
	height := m.height
	if width == 0 {
		if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && w > 0 {
			width = w
		} else {
			width = 80
		}
	}
	if height == 0 {
		if _, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil && h > 0 {
			height = h
		} else {
			height = 24
		}
	}

	// Calculate content width (max 70, min 40)
	contentWidth := min(width-8, 70)
	if contentWidth < 40 {
		contentWidth = width - 4
	}

	// Determine if descriptions should be shown
	showDesc := height > 20
	
	// Calculate how many items we can show
	// Overhead: Title(2) + Progress(2) + Help(2) + Border(4) + Padding(2) + Scroll indicators(2) = ~14 lines
	linesPerItem := 1
	if showDesc {
		linesPerItem = 2 // title + description
	}
	
	availableLines := height - 14
	if availableLines < 3 {
		availableLines = 3
	}
	maxVisibleItems := availableLines / linesPerItem
	if maxVisibleItems < 3 {
		maxVisibleItems = 3 // Minimum 3 items visible
	}
	if maxVisibleItems > len(m.items) {
		maxVisibleItems = len(m.items)
	}

	// Calculate viewport start/end to keep cursor visible
	viewportStart := 0
	viewportEnd := len(m.items)
	
	if len(m.items) > maxVisibleItems {
		// Need scrolling - keep cursor visible within viewport
		if m.cursor < viewportStart {
			viewportStart = m.cursor
		} else if m.cursor >= viewportStart+maxVisibleItems {
			viewportStart = m.cursor - maxVisibleItems + 1
		}
		
		// Recalculate based on cursor position
		viewportStart = m.cursor - maxVisibleItems/2
		if viewportStart < 0 {
			viewportStart = 0
		}
		
		viewportEnd = viewportStart + maxVisibleItems
		if viewportEnd > len(m.items) {
			viewportEnd = len(m.items)
			viewportStart = viewportEnd - maxVisibleItems
			if viewportStart < 0 {
				viewportStart = 0
			}
		}
	}


	var b strings.Builder

	// Title with fancy styling
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF79C6")).
		Background(lipgloss.Color("#44475a")).
		Padding(0, 2).
		Width(contentWidth).
		Align(lipgloss.Center)

	b.WriteString(titleStyle.Render(m.title))
	b.WriteString("\n\n")

	// Show scroll up indicator if needed
	if viewportStart > 0 {
		scrollUpStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffb86c")).
			Align(lipgloss.Center).
			Width(contentWidth)
		b.WriteString(scrollUpStyle.Render(fmt.Sprintf("▲ %d more above", viewportStart)))
		b.WriteString("\n")
	}

	// Items with improved visibility (only visible ones)
	for i := viewportStart; i < viewportEnd; i++ {
		item := m.items[i]
		var line string
		var style lipgloss.Style

		if i == m.cursor {
			// Selected item - highlighted background
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#282a36")).
				Background(lipgloss.Color("#50fa7b")).
				Bold(true).
				Padding(0, 1).
				Width(contentWidth - 4)
			line = fmt.Sprintf("▸ %s", item.Title)
		} else {
			// Normal item
			style = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f8f8f2")).
				Padding(0, 1).
				Width(contentWidth - 4)
			line = fmt.Sprintf("  %s", item.Title)
		}

		b.WriteString(style.Render(line))
		b.WriteString("\n")

		// Description with brighter color (only if space allows)
		if item.Description != "" && showDesc {
			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#8be9fd")).
				PaddingLeft(4)
			b.WriteString(descStyle.Render(item.Description))
			b.WriteString("\n")
		}
	}

	// Show scroll down indicator if needed
	if viewportEnd < len(m.items) {
		scrollDownStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffb86c")).
			Align(lipgloss.Center).
			Width(contentWidth)
		b.WriteString(scrollDownStyle.Render(fmt.Sprintf("▼ %d more below", len(m.items)-viewportEnd)))
		b.WriteString("\n")
	}

	// Visual progress bar
	progress := float64(m.cursor+1) / float64(len(m.items))
	barWidth := contentWidth - 10
	if barWidth < 10 {
		barWidth = 10
	}
	filled := int(progress * float64(barWidth))
	if filled < 1 {
		filled = 1
	}
	
	progressBar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	posStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#50fa7b")).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(posStyle.Render(progressBar))
	b.WriteString("\n")
	
	// Position text
	posTextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6272a4")).
		Align(lipgloss.Center).
		Width(contentWidth)
	b.WriteString(posTextStyle.Render(fmt.Sprintf("%d of %d", m.cursor+1, len(m.items))))
	b.WriteString("\n")

	// Help with better visibility
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#bd93f9")).
		Align(lipgloss.Center).
		Width(contentWidth).
		Italic(true)

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑↓/jk nav • PgUp/Dn jump • enter select • q quit"))


	// Wrap everything in a bordered box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#ff79c6")).
		Padding(1, 2).
		Width(contentWidth)

	boxContent := boxStyle.Render(b.String())

	// Center only horizontally, keep at top
	centeredStyle := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center)

	return centeredStyle.Render(boxContent)

}



// Selected returns the selected index (-1 if canceled)
func (m SelectorModel) Selected() int {
	if m.canceled {
		return -1
	}
	return m.selected
}

// RunSelector runs the interactive selector and returns the selected index
func RunSelector(title string, items []SelectorItem) (int, error) {
	model := NewSelector(title, items)
	p := tea.NewProgram(model)

	finalModel, err := p.Run()
	if err != nil {
		return -1, err
	}

	return finalModel.(SelectorModel).Selected(), nil
}

// SelectFromStrings is a convenience function to select from string slice
func SelectFromStrings(title string, options []string) (int, string, error) {
	items := make([]SelectorItem, len(options))
	for i, opt := range options {
		items[i] = SelectorItem{Title: opt, Value: opt}
	}

	idx, err := RunSelector(title, items)
	if err != nil {
		return -1, "", err
	}

	if idx < 0 || idx >= len(options) {
		return -1, "", nil
	}

	return idx, options[idx], nil
}

// SelectSSHKey shows SSH key selector with auto-suggestions
func SelectSSHKey(keys []string, currentKey string) (int, string, error) {
	items := make([]SelectorItem, len(keys))
	for i, key := range keys {
		desc := ""
		if key == currentKey {
			desc = "Currently selected"
		}
		items[i] = SelectorItem{
			Title:       key,
			Description: desc,
			Value:       key,
		}
	}

	idx, err := RunSelector("Select SSH Key", items)
	if err != nil {
		return -1, "", err
	}

	if idx < 0 || idx >= len(keys) {
		return -1, "", nil
	}

	return idx, keys[idx], nil
}

// SelectAccount shows account selector
func SelectAccountInteractive(accounts []string, activeAccount string) (int, string, error) {
	items := make([]SelectorItem, len(accounts))
	for i, acc := range accounts {
		desc := ""
		if acc == activeAccount {
			desc = "● Active"
		}
		items[i] = SelectorItem{
			Title:       acc,
			Description: desc,
			Value:       acc,
		}
	}

	idx, err := RunSelector("Select Account", items)
	if err != nil {
		return -1, "", err
	}

	if idx < 0 || idx >= len(accounts) {
		return -1, "", nil
	}

	return idx, accounts[idx], nil
}

// SelectMethodInteractive shows method selector (SSH/Token)
func SelectMethodInteractive(hasSSH, hasToken bool) (string, error) {
	var items []SelectorItem

	if hasSSH {
		items = append(items, SelectorItem{
			Title:       "🔑 SSH",
			Description: "Use SSH key authentication",
			Value:       "ssh",
		})
	}
	if hasToken {
		items = append(items, SelectorItem{
			Title:       "🔐 Token (HTTPS)",
			Description: "Use Personal Access Token",
			Value:       "token",
		})
	}

	if len(items) == 0 {
		return "", fmt.Errorf("no authentication methods available")
	}

	if len(items) == 1 {
		return items[0].Value, nil
	}

	idx, err := RunSelector("Select Authentication Method", items)
	if err != nil {
		return "", err
	}

	if idx < 0 || idx >= len(items) {
		return "", nil
	}

	return items[idx].Value, nil
}

// SelectPlatformInteractive shows platform selector
func SelectPlatformInteractive() (string, error) {
	items := []SelectorItem{
		{Title: "🐙 GitHub", Description: "github.com", Value: "github"},
		{Title: "🦊 GitLab", Description: "gitlab.com", Value: "gitlab"},
		{Title: "🪣 Bitbucket", Description: "bitbucket.org", Value: "bitbucket"},
		{Title: "🍵 Gitea", Description: "Self-hosted Gitea", Value: "gitea"},
		{Title: "🏔️ Codeberg", Description: "codeberg.org", Value: "codeberg"},
		{Title: "🌐 Other", Description: "Custom Git server", Value: "other"},
	}

	idx, err := RunSelector("Select Platform", items)
	if err != nil {
		return "", err
	}

	if idx < 0 || idx >= len(items) {
		return "", nil
	}

	return items[idx].Value, nil
}
