package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.
func (m *Model) View() string {
	var sb strings.Builder

	sb.WriteString(m.styles.header.Render("fixer — Training Labs"))
	sb.WriteString("\n")

	// Main layout: sidebar (left) + content (right)
	sidebar := m.renderToolsSidebar()
	content := m.renderContentArea()

	// Make sidebar span the same height as content
	contentHeight := lipgloss.Height(content)
	sidebarStyled := lipgloss.NewStyle().
		Height(contentHeight).
		Render(sidebar)

	// Join horizontally
	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, sidebarStyled, content)

	sb.WriteString(mainLayout)
	sb.WriteString("\n")
	sb.WriteString(m.renderBottomBar())

	return sb.String()
}

func (m *Model) renderToolsSidebar() string {
	var sb strings.Builder

	if len(m.toolGroups) == 0 {
		return m.styles.sidebar.Render("No tools")
	}

	for i, group := range m.toolGroups {
		isSelected := i == m.selectedToolIdx

		if isSelected {
			sb.WriteString(m.styles.selectedTool.Render("▸ " + group.Category) + "\n")
		} else {
			sb.WriteString(m.styles.unselectedTool.Render("   " + group.Category) + "\n")
		}
	}

	return m.styles.sidebar.Render(sb.String())
}

func (m *Model) renderContentArea() string {
	tasksBar := m.renderTasksBar()
	infoPanel := m.renderInfoPanel()

	// Tasks bar (top) + Info panel (bottom)
	return lipgloss.JoinVertical(lipgloss.Left, tasksBar, infoPanel)
}

func (m *Model) renderTasksBar() string {
	var sb strings.Builder

	if len(m.toolGroups) == 0 {
		return m.styles.tasksBar.Render("No tools")
	}

	currentGroup := m.toolGroups[m.selectedToolIdx]

	if len(currentGroup.Labs) == 0 {
		return m.styles.tasksBar.Render("  No tasks")
	}

	// Render tasks horizontally as a bar
	for j, labWithState := range currentGroup.Labs {
		lab := labWithState.Lab
		isSelected := j == m.selectedLabIdx

		// State indicator
		stateChar := "□"
		if labWithState.State == StateActive {
			stateChar = "●"
		} else if labWithState.State == StateIdle {
			stateChar = "○"
		}

		name := truncate(lab.Name, 15)

		if isSelected {
			sb.WriteString(m.styles.selectedTask.Render(" " + stateChar + " " + name + " "))
		} else {
			sb.WriteString(m.styles.unselectedTask.Render(" " + stateChar + " " + name + " "))
		}

		// Add separator between tasks
		if j < len(currentGroup.Labs)-1 {
			sb.WriteString(m.styles.tasksBarSep.Render("│"))
		}
	}

	return m.styles.tasksBar.Render(sb.String())
}

func (m *Model) renderInfoPanel() string {
	var sb strings.Builder

	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		return m.styles.infoPanel.Render("Select a task")
	}

	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx]

	// Title
	sb.WriteString(m.styles.taskTitle.Render(" " + lab.Lab.Name + "\n"))

	// Status indicator
	switch lab.State {
	case StateActive:
		sb.WriteString("  " + m.styles.activeLab.Render("● Active (in shell)") + "\n\n")
	case StateIdle:
		sb.WriteString("  " + m.styles.idleLab.Render("○ Running (idle)") + "\n\n")
	default:
		sb.WriteString("  " + m.styles.stoppedLab.Render("□ Stopped") + "\n\n")
	}

	// Goal
	if lab.Lab.Goal != "" {
		sb.WriteString(m.styles.infoLabel.Render("Goal:") + "\n")
		sb.WriteString("  " + lab.Lab.Goal + "\n\n")
	}

	// Description
	if lab.Lab.Description != "" {
		sb.WriteString(m.styles.infoLabel.Render("Description:") + "\n")
		for _, line := range strings.Split(lab.Lab.Description, "\n") {
			sb.WriteString("  " + line + "\n")
		}
		sb.WriteString("\n")
	}

	// Hints
	if len(lab.Lab.Hints) > 0 {
		sb.WriteString(m.styles.infoLabel.Render("Hints:") + "\n")
		for _, hint := range lab.Lab.Hints {
			sb.WriteString("  • " + hint + "\n")
		}
		sb.WriteString("\n")
	}

	// Validation commands
	if len(lab.Lab.Validate) > 0 {
		sb.WriteString(m.styles.infoLabel.Render("Validation:") + "\n")
		for _, check := range lab.Lab.Validate {
			sb.WriteString("  $ " + check + "\n")
		}
		sb.WriteString("\n")
	}

	// Last validation result
	if m.lastValidation != nil {
		sb.WriteString("\n")
		if m.lastValidation.Passed {
			sb.WriteString(m.styles.checkPassed.Render("  ✓ All checks passed!"))
		} else {
			sb.WriteString(m.styles.checkFailed.Render("  ✗ Some checks failed"))
		}
		sb.WriteString("\n")
	}

	// Log section
	if m.showLog && len(m.logBuffer) > 0 {
		sb.WriteString(m.styles.infoLabel.Render("Log:") + "\n")
		for _, line := range m.logBuffer {
			sb.WriteString("  " + m.styles.logStyle.Render(line) + "\n")
		}
	} else if !m.showLog {
		sb.WriteString("\n  " + m.styles.dimText.Render("(press o to show log)"))
	}

	return m.styles.infoPanel.Render(sb.String())
}

func (m *Model) renderBottomBar() string {
	var parts []string

	parts = append(parts, m.styles.bottomBarMode.Render("TUI MODE"))
	parts = append(parts, m.styles.dimText.Render("↑/↓: tools"))
	parts = append(parts, m.styles.dimText.Render("h/l: tasks"))
	parts = append(parts, m.styles.dimText.Render("Enter: start"))
	parts = append(parts, m.styles.dimText.Render("e: shell"))
	parts = append(parts, m.styles.dimText.Render("s: stop"))
	parts = append(parts, m.styles.dimText.Render("r: reset"))
	parts = append(parts, m.styles.dimText.Render("v: validate"))
	parts = append(parts, m.styles.dimText.Render("t: task"))
	parts = append(parts, m.styles.dimText.Render("o: log"))
	parts = append(parts, m.styles.dimText.Render("q: quit"))

	return m.styles.bottomBar.Render(strings.Join(parts, "  │  "))
}
