package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
)

const panelHintH = 2 // lines reserved for hint at bottom of panels

// View implements tea.Model.
func (m *Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return "fixer — Training Labs"
	}

	// Fixed layout dimensions
	headerH := 1
	footerH := 1
	panelGap := 1 // gap between left sidebar and right content

	availableH := m.height - headerH - footerH
	availableW := m.width

	// Sidebar width: fixed 24 chars
	sidebarW := 24
	contentW := availableW - sidebarW - panelGap
	if contentW < 20 {
		contentW = 20
		sidebarW = availableW - contentW - panelGap
	}

	// Tasks bar height: fixed 8 lines
	tasksBarH := 8
	infoPanelH := availableH - tasksBarH - panelGap
	if infoPanelH < 5 {
		infoPanelH = 5
		tasksBarH = availableH - infoPanelH - panelGap
	}

	// Build header
	var sb strings.Builder
	sb.WriteString(m.styles.header.Render(
		centerPad("fixer — Training Labs", availableW),
	))
	sb.WriteString("\n")

	// Build sidebar (Width/Height set inner area; border adds +2 outside)
	sidebar := m.renderToolsSidebar(sidebarW-2, availableH-2)

	// Build right content: tasks bar + info panel
	tasksBar := m.renderTasksBar(contentW-2, tasksBarH-2)
	infoPanel := m.renderInfoPanel(contentW-2, infoPanelH-2)

	gapLine := strings.Repeat(" ", contentW)
	rightCol := lipgloss.JoinVertical(lipgloss.Left, tasksBar, gapLine, infoPanel)

	// Join left and right
	gapCol := strings.Repeat("\n", availableH)
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, sidebar, gapCol, rightCol))
	sb.WriteString("\n")

	// Footer
	sb.WriteString(m.styles.bottomBar.Render(m.buildFooterText(availableW)))

	return sb.String()
}

func (m *Model) renderToolsSidebar(w, h int) string {
	innerW := w - 2 // Width sets inner area (content + padding), padding is 1+1
	innerH := h - 2 // Height sets inner area (content + padding), padding is 1+1
	if innerW < 4 {
		innerW = 4
	}
	if innerH < 1 {
		innerH = 1
	}

	var lines []string

	// Title: " TOOLS " centered
	titleSep := strings.Repeat("─", innerW)
	lines = append(lines, centerPad("Tools", innerW))
	lines = append(lines, titleSep)

	// Tool list
	if len(m.toolGroups) == 0 {
		lines = append(lines, centerPad("No tools", innerW))
	} else {
		for i, group := range m.toolGroups {
			isSelected := i == m.selectedToolIdx
			text := group.Category
			if isSelected {
				text = "▸ " + text
			}
			rendered := centerPad(text, innerW)
			lines = append(lines, rendered)
		}
	}

	// Fill remaining space with empty lines
	for len(lines) < innerH-panelHintH {
		lines = append(lines, "")
	}

	// Hint at bottom
	lines = append(lines, "")
	lines = append(lines, centerPad("j/k navigate", innerW))

	content := strings.Join(lines, "\n")
	return m.styles.sidebar.Width(w).Height(h).Render(content)
}

func (m *Model) renderTasksBar(w, h int) string {
	innerW := w - 2
	innerH := h - 2
	if innerW < 4 {
		innerW = 4
	}
	if innerH < 1 {
		innerH = 1
	}

	var lines []string

	// Title: " Tasks " centered
	titleSep := strings.Repeat("─", innerW)
	lines = append(lines, centerPad("Tasks", innerW))
	lines = append(lines, titleSep)

	// Tasks rendered in a grid (2 rows of tasks if needed)
	if len(m.toolGroups) > 0 {
		currentGroup := m.toolGroups[m.selectedToolIdx]
		if len(currentGroup.Labs) > 0 {
			lines = append(lines, "") // gap after title

			// Build task buttons
			var taskButtons []string
			for j, labWithState := range currentGroup.Labs {
				lab := labWithState.Lab
				isSelected := j == m.selectedLabIdx

				stateChar := "□"
				if labWithState.State == StateActive {
					stateChar = "●"
				} else if labWithState.State == StateIdle {
					stateChar = "○"
				}

				name := truncate(lab.Name, 20)
				btnText := fmt.Sprintf(" %s %s", stateChar, name)

				btnW := lipgloss.Width(btnText) + 4
				if isSelected {
					taskButtons = append(taskButtons, m.styles.selectedTask.Width(btnW).Render(btnText))
				} else {
					taskButtons = append(taskButtons, m.styles.unselectedTask.Width(btnW).Render(btnText))
				}
			}

			// Center the task buttons row
			buttonsRow := strings.Join(taskButtons, " ")
			lines = append(lines, centerPad(buttonsRow, innerW))
			lines = append(lines, "")
		}
	}

	// Fill remaining space
	for len(lines) < innerH-panelHintH {
		lines = append(lines, "")
	}

	// Hint at bottom
	lines = append(lines, "")
	lines = append(lines, centerPad("h/l navigate", innerW))

	content := strings.Join(lines, "\n")
	return m.styles.tasksBar.Width(w).Height(h).Render(content)
}

func (m *Model) renderInfoPanel(w, h int) string {
	innerW := w - 2
	innerH := h - 2
	if innerW < 4 {
		innerW = 4
	}
	if innerH < 1 {
		innerH = 1
	}

	var lines []string

	// Title row: "Information" centered
	titleSep := strings.Repeat("─", innerW)
	lines = append(lines, centerPad("Information", innerW))
	lines = append(lines, titleSep)

	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		lines = append(lines, "")
		lines = append(lines, centerPad("Select a task to begin", innerW))
		for len(lines) < innerH {
			lines = append(lines, "")
		}
		content := strings.Join(lines, "\n")
		return m.styles.infoPanel.Width(w).Height(h).Render(content)
	}

	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx]

	// State badge in top-right corner only
	stateText := ""
	switch lab.State {
	case StateActive:
		stateText = m.styles.activeLab.Render("● Active")
	case StateIdle:
		stateText = m.styles.idleLab.Render("○ Idle")
	default:
		stateText = m.styles.stoppedLab.Render("□ Stopped")
	}
	lines = append(lines, lipgloss.NewStyle().Align(lipgloss.Right).Width(innerW).Render(stateText))
	lines = append(lines, "")

	// Lab name centered
	lines = append(lines, centerPad(lab.Lab.Name, innerW))
	lines = append(lines, "")

	// Goal
	if lab.Lab.Goal != "" {
		lines = append(lines, m.styles.infoLabel.Render(" Goal:"))
		wrapped := wrapText(lab.Lab.Goal, innerW-2)
		for _, line := range wrapped {
			lines = append(lines, "  "+line)
		}
		lines = append(lines, "")
	}

	// Description
	if lab.Lab.Description != "" {
		lines = append(lines, m.styles.infoLabel.Render(" Description:"))
		wrapped := wrapText(lab.Lab.Description, innerW-2)
		for _, line := range wrapped {
			lines = append(lines, "  "+line)
		}
		lines = append(lines, "")
	}

	// Hints
	if len(lab.Lab.Hints) > 0 {
		lines = append(lines, m.styles.infoLabel.Render(" Hints:"))
		for _, hint := range lab.Lab.Hints {
			wrapped := wrapText("  • "+hint, innerW-2)
			lines = append(lines, wrapped...)
		}
		lines = append(lines, "")
	}

	// Last validation result
	if m.lastValidation != nil {
		if m.lastValidation.Passed {
			lines = append(lines, centerPad(m.styles.checkPassed.Render("✓ All checks passed!"), innerW))
		} else {
			lines = append(lines, centerPad(m.styles.checkFailed.Render("✗ Some checks failed"), innerW))
		}
		lines = append(lines, "")
	}

	// Log section
	if m.showLog && len(m.logBuffer) > 0 {
		lines = append(lines, m.styles.infoLabel.Render(" Log:"))
		for _, line := range m.logBuffer {
			lines = append(lines, "  "+m.styles.logStyle.Render(line))
		}
	} else if !m.showLog {
		lines = append(lines, centerPad(m.styles.dimText.Render("(press o to show log)"), innerW))
	}

	// Fill remaining with empty lines
	for len(lines) < innerH {
		lines = append(lines, "")
	}

	content := strings.Join(lines, "\n")
	return m.styles.infoPanel.Width(w).Height(h).Render(content)
}

func (m *Model) buildFooterText(width int) string {
	parts := []string{
		m.styles.bottomBarMode.Render("TUI"),
		"↑↓: tools",
		"h/l: tasks",
		"Enter: start",
		"e: shell",
		"s: stop",
		"r: reset",
		"v: validate",
		"o: log",
		"q: quit",
	}
	return centerPad(strings.Join(parts, "  │  "), width-4)
}

// wrapText wraps text to the given width.
func wrapText(text string, maxWidth int) []string {
	var lines []string
	words := strings.Fields(text)
	var current string
	for _, word := range words {
		if current == "" {
			current = word
		} else if len(current)+1+len(word) <= maxWidth {
			current += " " + word
		} else {
			lines = append(lines, current)
			current = word
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	if len(lines) == 0 {
		lines = append(lines, "")
	}
	return lines
}

// rightPad right-aligns text within the given width.
func rightPad(s string, width int) string {
	textLen := utf8.RuneCountInString(s)
	if textLen >= width {
		return s
	}
	return strings.Repeat(" ", width-textLen) + s
}
