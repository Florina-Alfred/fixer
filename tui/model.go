package tui

import (
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"redalf.de/fixer/docker"
	"redalf.de/fixer/labs"
)

const (
	containerNamePrefix = "fixer-lab-"
)

// Mode represents which interface mode the app is in.
type Mode int

const (
	ModeTUI Mode = iota
)

// ContainerState tracks the state of a lab's container.
type ContainerState int

const (
	StateStopped ContainerState = iota
	StateIdle    // Running but not in shell
	StateActive  // Currently in shell
)

// DockerOpComplete signals a Docker operation finished.
type DockerOpComplete struct {
	Op      string
	Err     error
	ID      string
	Out     string
	Passed  bool
	Command string
}

// ValidationResult holds the outcome of lab checks.
type ValidationResult struct {
	Passed bool
	Checks []CheckResult
}

// CheckResult holds the outcome of a single check.
type CheckResult struct {
	Command string
	Passed  bool
	Output  string
}

// LabWithState holds a lab and its container state.
type LabWithState struct {
	Lab       labs.Lab
	State     ContainerState
	Container string
}

// ToolGroup groups labs by category.
type ToolGroup struct {
	Category string
	Labs     []LabWithState
}

// Model is the central state for the Bubble Tea application.
type Model struct {
	mode              Mode
	toolGroups        []ToolGroup
	selectedToolIdx   int
	selectedLabIdx    int
	showLog           bool
	logBuffer         []string
	dockerCli         *docker.Client
	width             int
	height            int
	validationResults []CheckResult
	lastValidation    *ValidationResult
	styles            *styles
	leftWidth         int
	rightWidth        int
	bottomHeight      int
	shellRequested    bool
}

type styles struct {
	header           lipgloss.Style
	leftPanel        lipgloss.Style
	rightPanel       lipgloss.Style
	bottomBar        lipgloss.Style
	selectedTool     lipgloss.Style
	unselectedTool   lipgloss.Style
	selectedLab      lipgloss.Style
	unselectedLab    lipgloss.Style
	activeLab        lipgloss.Style
	idleLab          lipgloss.Style
	title            lipgloss.Style
	modeTUI          lipgloss.Style
	hint             lipgloss.Style
	goal             lipgloss.Style
	checkPassed      lipgloss.Style
	checkFailed      lipgloss.Style
	logStyle         lipgloss.Style
	border           lipgloss.Style
}

func defaultStyles() *styles {
	s := &styles{}

	blue := lipgloss.Color("#006BB6")
	green := lipgloss.Color("#50C878")
	yellow := lipgloss.Color("#FFB347")
	red := lipgloss.Color("#FF6B6B")
	gray := lipgloss.Color("#6C757D")
	white := lipgloss.Color("#FFFFFF")
	darkGray := lipgloss.Color("#343A40")

	s.header = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B35")).
		PaddingLeft(2)

	s.title = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Width(14)

	s.leftPanel = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(gray).
		Padding(1, 2)

	s.rightPanel = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(gray).
		Padding(1, 2)

	s.bottomBar = lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(darkGray).
		Padding(0, 2).
		Width(0)

	s.selectedTool = lipgloss.NewStyle().
		Foreground(blue).
		Bold(true).
		Background(lipgloss.Color("#1a1a2e"))

	s.unselectedTool = lipgloss.NewStyle().
		Foreground(gray)

	s.selectedLab = lipgloss.NewStyle().
		Foreground(white).
		Bold(true).
		Background(lipgloss.Color("#1a1a2e"))

	s.unselectedLab = lipgloss.NewStyle().
		Foreground(gray)

	s.activeLab = lipgloss.NewStyle().
		Foreground(green).
		Bold(true)

	s.idleLab = lipgloss.NewStyle().
		Foreground(yellow).
		Bold(true)

	s.modeTUI = lipgloss.NewStyle().
		Bold(true).
		Foreground(green)

	s.hint = lipgloss.NewStyle().
		Foreground(gray)

	s.goal = lipgloss.NewStyle().
		Foreground(green).
		Italic(true)

	s.checkPassed = lipgloss.NewStyle().
		Foreground(green)

	s.checkFailed = lipgloss.NewStyle().
		Foreground(red)

	s.logStyle = lipgloss.NewStyle().
		Foreground(gray)

	s.border = lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(gray)

	return s
}

func (m *Model) initStyles() {
	m.styles = defaultStyles()
}

// groupLabsByCategory groups labs into ToolGroups by category.
func groupLabsByCategory(labsList []labs.Lab) []ToolGroup {
	categoryMap := make(map[string][]labs.Lab)
	for _, lab := range labsList {
		categoryMap[lab.Category] = append(categoryMap[lab.Category], lab)
	}

	var groups []ToolGroup
	for category, categoryLabs := range categoryMap {
		var labsWithState []LabWithState
		for _, lab := range categoryLabs {
			labsWithState = append(labsWithState, LabWithState{
				Lab:   lab,
				State: StateStopped,
			})
		}
		groups = append(groups, ToolGroup{
			Category: category,
			Labs:     labsWithState,
		})
	}

	// Sort by category name
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Category < groups[j].Category
	})

	return groups
}

// NewModel creates a new application model.
func NewModel(labsList []labs.Lab) *Model {
	m := &Model{
		mode:          ModeTUI,
		toolGroups:    groupLabsByCategory(labsList),
		selectedToolIdx: 0,
		selectedLabIdx:  0,
		dockerCli:       docker.New(),
	}
	m.initStyles()
	return m
}

// GetState returns the current model state for restoration.
func (m *Model) GetState() (toolGroups []ToolGroup, selectedToolIdx, selectedLabIdx int, containerName, containerID string) {
	return m.toolGroups, m.selectedToolIdx, m.selectedLabIdx, m.getCurrentContainerName(), ""
}

// getCurrentContainerName returns the container name for the currently selected lab.
func (m *Model) getCurrentContainerName() string {
	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		return ""
	}
	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx]
	return containerNamePrefix + normalizeName(lab.Lab.Name)
}

// normalizeName converts a lab name to a valid container name.
func normalizeName(name string) string {
	result := ""
	for _, c := range name {
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			result += string(c)
		} else if c >= 'A' && c <= 'Z' {
			result += string(unicode.ToLower(c))
		} else {
			result += "-"
		}
	}
	return result
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.leftWidth = m.width / 3
		if m.leftWidth < 25 {
			m.leftWidth = 25
		}
		if m.leftWidth > 45 {
			m.leftWidth = 45
		}
		m.rightWidth = m.width - m.leftWidth - 2
		m.bottomHeight = 3

	case tea.KeyMsg:
		return m.handleKey(msg)

	case DockerOpComplete:
		m.handleDockerComplete(msg)
		return m, nil

	case tea.FocusMsg:
		return m, nil

	case tea.BlurMsg:
		return m, nil
	}

	return m, nil
}

func (m *Model) handleKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	return m.handleTUIKey(k)
}

func (m *Model) handleTUIKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit

	// Navigate tools (vertical)
	case "up", "k":
		if m.selectedToolIdx > 0 {
			m.selectedToolIdx--
			m.selectedLabIdx = 0
		}
		return m, nil

	case "down", "j":
		if m.selectedToolIdx < len(m.toolGroups)-1 {
			m.selectedToolIdx++
			m.selectedLabIdx = 0
		}
		return m, nil

	// Navigate labs within tool (horizontal)
	case "h":
		if len(m.toolGroups) > 0 && len(m.toolGroups[m.selectedToolIdx].Labs) > 0 {
			if m.selectedLabIdx > 0 {
				m.selectedLabIdx--
			}
		}
		return m, nil

	case "l":
		if len(m.toolGroups) > 0 && len(m.toolGroups[m.selectedToolIdx].Labs) > 0 {
			if m.selectedLabIdx < len(m.toolGroups[m.selectedToolIdx].Labs)-1 {
				m.selectedLabIdx++
			}
		}
		return m, nil

	case "enter":
		return m.startLab()

	case "r":
		return m.resetLab()

	case "s":
		return m.stopLab()

	case "v":
		return m.validateLab()

	case "e":
		if !m.hasRunningContainer() {
			m.log("No container running. Start a lab first.")
			return m, nil
		}
		running, err := m.dockerCli.IsRunning(m.getCurrentContainerName())
		if err != nil || !running {
			if err != nil {
				m.log("Container check failed: " + err.Error())
			}
			m.log("Container not running — press Enter to start it first")
			return m, nil
		}
		m.log("Exiting TUI to open shell...")
		m.setContainerActive(true)
		m.shellRequested = true
		return m, tea.Quit

	case "t":
		m.showTask()
		return m, nil

	case "o":
		m.showLog = !m.showLog
		return m, nil
	}

	return m, nil
}

func (m *Model) hasRunningContainer() bool {
	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		return false
	}
	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx]
	return lab.State == StateIdle || lab.State == StateActive
}

func (m *Model) setContainerActive(active bool) {
	if len(m.toolGroups) > 0 && len(m.toolGroups[m.selectedToolIdx].Labs) > 0 {
		if active {
			m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].State = StateActive
		} else {
			m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].State = StateIdle
		}
	}
}

func (m *Model) showTask() {
	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		return
	}
	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].Lab
	m.log("Task: " + lab.Goal)
	if lab.Description != "" {
		m.log("Description: " + lab.Description)
	}
	if len(lab.Hints) > 0 {
		m.log("Hints:")
		for _, hint := range lab.Hints {
			m.log("  - " + hint)
		}
	}
}

func (m *Model) handleDockerComplete(msg DockerOpComplete) {
	if msg.Err != nil && msg.Op != "check" {
		m.log("Docker error (" + msg.Op + "): " + msg.Err.Error())
	}
	if msg.ID != "" {
		// Store container ID in the lab state
		if len(m.toolGroups) > 0 && len(m.toolGroups[m.selectedToolIdx].Labs) > 0 {
			// Container ID is not stored, but we track state
		}
	}

	if msg.Op == "check" {
		m.validationResults = append(m.validationResults, CheckResult{
			Command: msg.Command,
			Passed:  msg.Passed,
			Output:  msg.Out,
		})

		if len(m.toolGroups) > 0 && len(m.toolGroups[m.selectedToolIdx].Labs) > 0 {
			lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].Lab
			if len(m.validationResults) >= len(lab.Validate) {
				allPassed := true
				for _, cr := range m.validationResults {
					if !cr.Passed {
						allPassed = false
						break
					}
				}
				m.lastValidation = &ValidationResult{
					Passed: allPassed,
					Checks: m.validationResults,
				}
				if allPassed {
					m.log("All validation checks passed!")
				} else {
					m.log("Validation failed:")
					for _, cr := range m.validationResults {
						status := "FAIL"
						if cr.Passed {
							status = "PASS"
						}
						m.log("  [" + status + "] " + cr.Command)
					}
				}
				m.validationResults = nil
			}
		}
	} else if msg.Op == "start" && msg.Err == nil {
		m.setContainerActive(false)
		m.log("Container started and ready")
	}
}

func (m *Model) startLab() (tea.Model, tea.Cmd) {
	if len(m.toolGroups) == 0 {
		m.log("No labs available")
		return m, nil
	}

	if len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		m.log("No labs in this category")
		return m, nil
	}

	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].Lab
	containerName := containerNamePrefix + normalizeName(lab.Name)

	// Check if container already exists and is running
	exists, _ := m.dockerCli.ContainerExists(containerName)
	if exists {
		running, _ := m.dockerCli.IsRunning(containerName)
		if running {
			m.log("Container already running: " + containerName)
			m.setContainerActive(false)
			return m, nil
		}
	}

	_ = m.dockerCli.CleanUp(containerName)

	m.log("Starting container for lab: " + lab.Name)
	return m, func() tea.Msg {
		id, err := m.dockerCli.Start(lab.Image, containerName)
		if err != nil {
			return DockerOpComplete{Op: "start", Err: err, ID: id}
		}

		// Run setup commands if any
		if len(lab.Setup) > 0 {
			setupErr := m.dockerCli.Setup(containerName, lab.Setup)
			if setupErr != nil {
				return DockerOpComplete{Op: "setup", Err: setupErr, ID: id}
			}
		}

		return DockerOpComplete{Op: "start", Err: nil, ID: id}
	}
}

func (m *Model) resetLab() (tea.Model, tea.Cmd) {
	containerName := m.getCurrentContainerName()
	if containerName == "" {
		m.log("No lab selected")
		return m, nil
	}

	m.log("Resetting container: " + containerName)
	return m, func() tea.Msg {
		err := m.dockerCli.CleanUp(containerName)
		if err != nil {
			m.log("Failed to stop/remove container: " + err.Error())
			return DockerOpComplete{Op: "reset", Err: err}
		}
		m.setContainerStateStopped()
		m.log("Container reset complete")
		return DockerOpComplete{Op: "reset", Err: nil}
	}
}

func (m *Model) stopLab() (tea.Model, tea.Cmd) {
	containerName := m.getCurrentContainerName()
	if containerName == "" {
		m.log("No lab selected")
		return m, nil
	}

	m.log("Stopping container: " + containerName)
	return m, func() tea.Msg {
		err := m.dockerCli.CleanUp(containerName)
		if err != nil {
			m.log("Failed to stop container: " + err.Error())
			return DockerOpComplete{Op: "stop", Err: err}
		}
		m.setContainerStateStopped()
		m.log("Container stopped and removed")
		return DockerOpComplete{Op: "stop", Err: nil}
	}
}

func (m *Model) setContainerStateStopped() {
	if len(m.toolGroups) > 0 && len(m.toolGroups[m.selectedToolIdx].Labs) > 0 {
		m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].State = StateStopped
	}
}

func (m *Model) validateLab() (tea.Model, tea.Cmd) {
	if len(m.toolGroups) == 0 || len(m.toolGroups[m.selectedToolIdx].Labs) == 0 {
		m.log("No lab selected")
		return m, nil
	}

	lab := m.toolGroups[m.selectedToolIdx].Labs[m.selectedLabIdx].Lab
	if lab.Name == "" {
		m.log("No lab selected")
		return m, nil
	}

	var cmds []tea.Cmd
	for _, check := range lab.Validate {
		cmds = append(cmds, m.runCheck(check))
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) runCheck(command string) tea.Cmd {
	return func() tea.Msg {
		passed, output, err := m.dockerCli.Validate(m.getCurrentContainerName(), []string{command})
		return DockerOpComplete{
			Op:      "check",
			Err:     err,
			Out:     output,
			Passed:  passed,
			Command: command,
		}
	}
}

func (m *Model) log(msg string) {
	m.logBuffer = append(m.logBuffer, msg)
	if len(m.logBuffer) > 200 {
		m.logBuffer = m.logBuffer[len(m.logBuffer)-150:]
	}
}

// RequestShell returns true if the TUI should exit and start a shell session.
func (m *Model) RequestShell() bool {
	return m.shellRequested
}

// ContainerName returns the name of the currently running container.
func (m *Model) ContainerName() string {
	return m.getCurrentContainerName()
}

// truncate renders text to fit within n runes.
func truncate(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	runes := []rune(s)
	return string(runes[:n-1]) + "…"
}
