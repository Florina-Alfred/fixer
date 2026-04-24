package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.
func (m *Model) View() string {
	var sb strings.Builder

	sb.WriteString(m.styles.header.Render("fixer — Training Labs\n"))
	sb.WriteString(m.renderMainView())
	sb.WriteString("\n")
	sb.WriteString(m.renderBottomBar())

	return sb.String()
}

func (m *Model) renderMainView() string {
	var sb strings.Builder

	left := m.renderLeftPanel()
	right := m.renderRightPanel()

	// Render panels side by side
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, right))

	return sb.String()
}

func (m *Model) renderLeftPanel() string {
	var sb strings.Builder

	if len(m.toolGroups) == 0 {
		return m.styles.leftPanel.Render("No labs available")
	}

	// Render each tool group
	for i, group := range m.toolGroups {
		isSelected := i == m.selectedToolIdx

		// Tool name
		if isSelected {
			sb.WriteString(m.styles.selectedTool.Render(" " + group.Category + " "))
		} else {
			sb.WriteString(m.styles.unselectedTool.Render(" " + group.Category + " "))
		}
		sb.WriteString("\n")

		// Labs in this tool (shown horizontally)
		for j, labWithState := range group.Labs {
			lab := labWithState.Lab
			isLabSelected := isSelected && j == m.selectedLabIdx

			// Determine the prefix based on state
			prefix := "  "
			if isLabSelected {
				switch labWithState.State {
				case StateActive:
					prefix = "  [●] " // Green active
				case StateIdle:
					prefix = "  [○] " // Yellow idle
				default:
					prefix = "  [ ] " // No state
				}
			}

			// Truncate lab name to fit
			labName := truncate(lab.Name, 20)

			if isLabSelected {
				sb.WriteString(m.styles.selectedLab.Render(prefix + labName))
			} else {
				sb.WriteString(m.styles.unselectedLab.Render(prefix + labName))
			}

			// Add state indicator
			if labWithState.State == StateActive {
				sb.WriteString(m.styles.activeLab.Render(" ●"))
			} else if labWithState.State == StateIdle {
				sb.WriteString(m.styles.idleLab.Render(" ○"))
			}

			sb.WriteString("\n")
		}

		if i < len(m.toolGroups)-1 {
			sb.WriteString("\n")
		}
	}

	return m.styles.leftPanel.Render(sb.String())
}

func (m *Model) renderRightPanel() string {
	var sb strings.Builder

	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		return m.styles.rightPanel.Render("Select a lab to view details")
	}

	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx]

	// Title
	sb.WriteString(m.styles.title.Render(lab.Lab.Name) + "\n\n")

	// Status indicator
	switch lab.State {
	case StateActive:
		sb.WriteString(m.styles.activeLab.Render("● Container: Active (in shell)") + "\n\n")
	case StateIdle:
		sb.WriteString(m.styles.idleLab.Render("○ Container: Running (idle)") + "\n\n")
	default:
		sb.WriteString("□ Container: Stopped\n\n")
	}

	// Goal
	if lab.Lab.Goal != "" {
		sb.WriteString(m.styles.goal.Render("Goal: " + lab.Lab.Goal) + "\n\n")
	}

	// Description
	if lab.Lab.Description != "" {
		sb.WriteString("Description:\n")
		for _, line := range strings.Split(lab.Lab.Description, "\n") {
			sb.WriteString("  " + line + "\n")
		}
		sb.WriteString("\n")
	}

	// Hints
	if len(lab.Lab.Hints) > 0 {
		sb.WriteString(m.styles.hint.Render("Hints:\n"))
		for _, hint := range lab.Lab.Hints {
			sb.WriteString("  • " + hint + "\n")
		}
		sb.WriteString("\n")
	}

	// Validation commands
	if len(lab.Lab.Validate) > 0 {
		sb.WriteString("Validation:\n")
		for _, check := range lab.Lab.Validate {
			sb.WriteString("  $ " + check + "\n")
		}
		sb.WriteString("\n")
	}

	// Last validation result
	if m.lastValidation != nil {
		sb.WriteString("\n")
		if m.lastValidation.Passed {
			sb.WriteString(m.styles.checkPassed.Render("✓ All checks passed!"))
		} else {
			sb.WriteString(m.styles.checkFailed.Render("✗ Some checks failed"))
		}
	}

	// Log section
	if m.showLog && len(m.logBuffer) > 0 {
		sb.WriteString("\n" + lipgloss.NewStyle().Bold(true).Render("Log:") + "\n")
		for _, line := range m.logBuffer {
			sb.WriteString("  " + m.styles.logStyle.Render(line) + "\n")
		}
	} else if !m.showLog {
		sb.WriteString("\n" + m.styles.hint.Render("(press l to show log)") + "\n")
	}

	// Quick task reminder
	sb.WriteString("\n" + m.styles.hint.Render("(press t to show task details)") + "\n")

	return m.styles.rightPanel.Render(sb.String())
}

func (m *Model) renderBottomBar() string {
	var parts []string

	parts = append(parts, m.styles.modeTUI.Render(" TUI MODE "))

	hints := []string{
		"↑/↓: tools",
		"←/→: labs",
		"Enter: start",
		"e: shell",
		"s: stop",
		"r: reset",
		"v: validate",
		"t: task",
		"l: log",
	}
	parts = append(parts, strings.Join(hints, "  |  "))

	barText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2C3136")).
		Padding(0, 2).
		Width(m.width).
		Render(strings.Join(parts, "    "))

	return barText
}
