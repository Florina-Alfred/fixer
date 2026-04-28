package main

import (
	"fmt"
	"os"
	"path/filepath"

	"redalf.de/fixer/docker"
	"redalf.de/fixer/labs"
	"redalf.de/fixer/pty"
	"redalf.de/fixer/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// savedState holds TUI state to restore after shell exit.
type savedState struct {
	toolGroups      []tui.ToolGroup
	selectedToolIdx int
	selectedLabIdx  int
}

func main() {
	labDir := findLabDir()

	labList, err := labs.LoadAll(labDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not load labs: %v\n", err)
		fmt.Fprintf(os.Stderr, "Make sure lab YAML files are in %s\n", labDir)
	}

	if len(labList) == 0 {
		fmt.Fprintf(os.Stderr, "No labs found. Expected YAML files in %s\n", labDir)
		fmt.Fprintf(os.Stderr, "Create lab files with: name, image, goal, checks, hints\n")
		os.Exit(1)
	}

	var prev *savedState
	for {
		m := tui.NewModel(labList, labDir)

		// Restore state from previous session if available
		if prev != nil {
			m.RestoreState(prev.toolGroups, prev.selectedToolIdx, prev.selectedLabIdx)
		}

		p := tea.NewProgram(m, tea.WithAltScreen())

		_, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if !m.RequestShell() {
			// Clean up all running containers before exit
			m.CleanupContainers()
			return
		}

		if m.ContainerName() == "" {
			fmt.Fprintf(os.Stderr, "No container running. Start a lab first.\n")
			return
		}

		// Save state before shell
		containerName := m.ContainerName()
		prevGroups, selectedToolIdx, selectedLabIdx, _, _ := m.GetState()
		prev = &savedState{
			toolGroups:      prevGroups,
			selectedToolIdx: selectedToolIdx,
			selectedLabIdx:  selectedLabIdx,
		}

		if err := pty.ExecuteShell(containerName); err != nil {
			fmt.Fprintf(os.Stderr, "Shell error: %v\n", err)
		}

		// Normalize terminal state before restarting TUI
		fmt.Print("\033[c\033[?1004l\033[0m")

		// Check if container is still running
		running, _ := docker.New().IsRunning(containerName)
		if !running {
			// Container was stopped, exit
			return
		}

		// Mark container as Idle (not Active) so it can be returned to
		m.SetContainerIdle()

		// Continue loop to restore TUI with existing state
	}
}

func findLabDir() string {
	// Check command-line args first
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--labs" || os.Args[i] == "-d" {
			if i+1 < len(os.Args) {
				return os.Args[i+1]
			}
		}
		if os.Args[i] == "-d" || os.Args[i] == "--dir" {
			if i+1 < len(os.Args) {
				return os.Args[i+1]
			}
		}
	}

	// Check environment variable
	if dir := os.Getenv("FIXER_LABS_DIR"); dir != "" {
		return dir
	}

	// Default: labs/ relative to executable
	exe, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exe)
		defaultDir := filepath.Join(exeDir, "labs")
		if _, err := os.Stat(defaultDir); err == nil {
			return defaultDir
		}
	}

	// Fallback to working directory
	cwd, err := os.Getwd()
	if err == nil {
		cwdDir := filepath.Join(cwd, "labs")
		if _, err := os.Stat(cwdDir); err == nil {
			return cwdDir
		}
	}

	return "labs"
}
