package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.
func (m *Model) View() string {
	return m.renderTUIView()
}

func (m *Model) renderTUIView() string {
	left := m.renderLeftPanel()
	right := m.renderRightPanel()

	mainArea := lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	bottom := m.renderBottomBar()

	return lipgloss.JoinVertical(lipgloss.Left, mainArea, bottom)
}

func (m *Model) renderLeftPanel() string {
	var sb strings.Builder

	if len(m.labs) == 0 {
		sb.WriteString("  No labs found.\n")
		sb.WriteString("  Put *.yml in the labs/ directory.\n")
		return m.styles.leftPanel.Render(" LABS " + "\n" + sb.String())
	}

	for i, lab := range m.labs {
		var line string
		if i == m.selectedIdx {
			line = m.styles.selectedItem.Render("▸ " + lab.Name) + "\n"
		} else {
			line = "  " + m.styles.unselectedItem.Render(lab.Name) + "\n"
		}
		sb.WriteString(line)
	}

	return m.styles.leftPanel.Render(" LABS " + "\n" + sb.String())
}

func (m *Model) renderRightPanel() string {
	var sb strings.Builder

	if m.activeLab.Name == "" {
		sb.WriteString("Select a lab and press Enter to start.\n\n")
		sb.WriteString("Keybindings:\n")
		sb.WriteString("  ↑/↓  Navigate labs\n")
		sb.WriteString("  Enter  Start lab\n")
		sb.WriteString("  e  Enter shell\n")
		sb.WriteString("  s  Stop container\n")
		sb.WriteString("  r  Reset lab\n")
		sb.WriteString("  v  Validate\n")
		sb.WriteString("  l  Toggle log\n")
		return m.styles.rightPanel.Render(sb.String())
	}

	sb.WriteString(m.styles.title.Render("🔧 " + m.activeLab.Name) + "\n")

	if m.activeLab.Goal != "" {
		sb.WriteString("\n  " + m.styles.goal.Render("Goal: " + m.activeLab.Goal) + "\n")
	}

	status := "● Running"
	if m.containerName == "" {
		status = "— None —"
	} else if m.containerID == "" {
		status = "○ Stopped"
	}
	sb.WriteString("\n  Container: " + m.styles.checkPassed.Render(status) + " " + m.containerName + "\n")

	if m.containerID != "" {
		sb.WriteString("  ID: " + truncate(m.containerID, 12) + "\n")
	}

	if len(m.activeLab.Hints) > 0 {
		sb.WriteString("\n  " + lipgloss.NewStyle().Bold(true).Render("Hints:") + "\n")
		for _, hint := range m.activeLab.Hints {
			sb.WriteString("    • " + hint + "\n")
		}
	}

	if len(m.activeLab.Checks) > 0 {
		sb.WriteString("\n  " + lipgloss.NewStyle().Bold(true).Render("Checks:") + "\n")
		for _, check := range m.activeLab.Checks {
			sb.WriteString("    $ " + check + "\n")
		}
	}

	if m.lastValidation != nil {
		sb.WriteString("\n")
		if m.lastValidation.Passed {
			sb.WriteString(m.styles.checkPassed.Render("✓ All checks passed!"))
		} else {
			sb.WriteString(m.styles.checkFailed.Render("✗ Some checks failed"))
		}
	}

	// Log section — shown when showLog is true
	if m.showLog && len(m.logBuffer) > 0 {
		sb.WriteString("\n  " + lipgloss.NewStyle().Bold(true).Render("Log:") + "\n")
		for _, line := range m.logBuffer {
			sb.WriteString("    " + m.styles.logStyle.Render(line) + "\n")
		}
	} else if !m.showLog {
		sb.WriteString("\n  " + m.styles.hint.Render("(press l to show log)") + "\n")
	}

	return m.styles.rightPanel.Render(sb.String())
}

func (m *Model) renderBottomBar() string {
	var parts []string

	parts = append(parts, m.styles.modeTUI.Render(" TUI MODE "))

	hints := []string{"Enter: start", "e: shell", "s: stop", "r: reset", "v: validate", "l: log"}
	parts = append(parts, strings.Join(hints, "  |  "))

	barText := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#2C3136")).
		Padding(0, 2).
		Render(strings.Join(parts, "    "))

	padding := (m.width - lipgloss.Width(barText)) / 2
	if padding < 0 {
		padding = 0
	}
	bar := strings.Repeat(" ", padding) + barText

	return bar
}


