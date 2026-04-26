package tui

import (
	"sort"
	"strings"
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
	StateIdle                   // Running but not in shell
	StateActive                 // Currently in shell
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
	header         lipgloss.Style
	headerTitle    lipgloss.Style
	sidebar        lipgloss.Style
	sidebarTitle   lipgloss.Style
	sidebarHint    lipgloss.Style
	tasksBar       lipgloss.Style
	tasksBarTitle  lipgloss.Style
	tasksBarHint   lipgloss.Style
	tasksBarSep    lipgloss.Style
	infoPanel      lipgloss.Style
	infoPanelTitle lipgloss.Style
	infoLegend     lipgloss.Style
	bottomBar      lipgloss.Style
	panelTitle     lipgloss.Style
	taskTitle      lipgloss.Style
	infoLabel      lipgloss.Style
	selectedTool   lipgloss.Style
	unselectedTool lipgloss.Style
	selectedTask   lipgloss.Style
	unselectedTask lipgloss.Style
	activeLab      lipgloss.Style
	idleLab        lipgloss.Style
	stoppedLab     lipgloss.Style
	title          lipgloss.Style
	bottomBarMode  lipgloss.Style
	dimText        lipgloss.Style
	hint           lipgloss.Style
	goal           lipgloss.Style
	checkPassed    lipgloss.Style
	checkFailed    lipgloss.Style
	logStyle       lipgloss.Style
	levelBadge     lipgloss.Style
	goalLabel      lipgloss.Style
	goalText       lipgloss.Style
	hintLabel      lipgloss.Style
	border         lipgloss.Style
}

func defaultStyles() *styles {
	s := &styles{}

	// Nord Dark palette
	nord0 := lipgloss.Color("#2E3440")
	nord1 := lipgloss.Color("#3B4252")
	nord3 := lipgloss.Color("#4C566A")
	nord5 := lipgloss.Color("#E5E9F0")
	nord6 := lipgloss.Color("#ECEFF4")
	nord8 := lipgloss.Color("#88C0D0")
	nord9 := lipgloss.Color("#81A1C1")
	nord12 := lipgloss.Color("#BF616A")
	nord14 := lipgloss.Color("#EBCB8B")
	nord15 := lipgloss.Color("#A3BE8C")

	selBg := nord1
	selFg := nord6

	// Header
	s.header = lipgloss.NewStyle().
		Foreground(nord8).
		Bold(true)

	s.headerTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord8).
		Align(lipgloss.Center)

	// Panels — equal top/bottom/left/right padding, vertically centered content
	panelBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(nord3).
		PaddingLeft(1).
		PaddingRight(1).
		PaddingTop(1).
		PaddingBottom(1).
		AlignVertical(lipgloss.Center)

	s.sidebar = panelBorder
	s.tasksBar = panelBorder
	s.infoPanel = panelBorder

	s.sidebarTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord9).
		Align(lipgloss.Center)

	s.sidebarHint = lipgloss.NewStyle().
		Foreground(nord3).
		Align(lipgloss.Center)

	s.tasksBarTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord9).
		Align(lipgloss.Center)

	s.tasksBarHint = lipgloss.NewStyle().
		Foreground(nord3).
		Align(lipgloss.Center)

	s.tasksBarSep = lipgloss.NewStyle().
		Foreground(nord3)

	s.infoPanelTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord9).
		Align(lipgloss.Center)

	// Bottom bar
	s.bottomBar = lipgloss.NewStyle().
		Foreground(nord3).
		Background(nord0).
		Padding(0, 2).
		Align(lipgloss.Center)

	s.bottomBarMode = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord15)

	// Navigation highlights
	s.selectedTool = lipgloss.NewStyle().
		Foreground(selFg).
		Bold(true).
		Background(selBg).
		Padding(0, 1)

	s.unselectedTool = lipgloss.NewStyle().
		Foreground(nord3).
		Padding(0, 1)

	s.selectedTask = lipgloss.NewStyle().
		Foreground(selFg).
		Bold(true).
		Background(selBg).
		Padding(0, 1)

	s.unselectedTask = lipgloss.NewStyle().
		Foreground(nord3).
		Padding(0, 1)

	// States
	s.activeLab = lipgloss.NewStyle().
		Foreground(nord15).
		Bold(true)

	s.idleLab = lipgloss.NewStyle().
		Foreground(nord14).
		Bold(true)

	s.stoppedLab = lipgloss.NewStyle().
		Foreground(nord3)

	// Info panel content
	s.taskTitle = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord5).
		Align(lipgloss.Center)

	s.infoLabel = lipgloss.NewStyle().
		Bold(true).
		Foreground(nord8)

	s.goal = lipgloss.NewStyle().
		Foreground(nord15).
		Italic(true)

	s.hint = lipgloss.NewStyle().
		Foreground(nord3)

	s.checkPassed = lipgloss.NewStyle().
		Foreground(nord15).
		Bold(true)

	s.checkFailed = lipgloss.NewStyle().
		Foreground(nord12).
		Bold(true)

	s.logStyle = lipgloss.NewStyle().
		Foreground(nord3)

	s.dimText = lipgloss.NewStyle().
		Foreground(nord3)

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
		// Sort labs by name for consistent ordering
		sort.Slice(categoryLabs, func(i, j int) bool {
			return categoryLabs[i].Name < categoryLabs[j].Name
		})
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
		mode:            ModeTUI,
		toolGroups:      groupLabsByCategory(labsList),
		selectedToolIdx: 0,
		selectedLabIdx:  0,
		dockerCli:       docker.New(),
	}
	m.initStyles()
	return m
}

// RestoreState restores container states and selection from a previous session.
func (m *Model) RestoreState(prevGroups []ToolGroup, selectedToolIdx, selectedLabIdx int) {
	// Restore container states by matching lab names
	for i := range prevGroups {
		for j := range m.toolGroups {
			if m.toolGroups[j].Category == prevGroups[i].Category {
				prevNameMap := make(map[string]ContainerState)
				for _, pl := range prevGroups[i].Labs {
					prevNameMap[pl.Lab.Name] = pl.State
				}
				for k := range m.toolGroups[j].Labs {
					if state, ok := prevNameMap[m.toolGroups[j].Labs[k].Lab.Name]; ok {
						m.toolGroups[j].Labs[k].State = state
					}
				}
			}
		}
	}
	// Restore selection
	m.selectedToolIdx = selectedToolIdx
	m.selectedLabIdx = selectedLabIdx
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

// centerPad centers text within the given width by padding with spaces.
func centerPad(s string, width int) string {
	textLen := utf8.RuneCountInString(s)
	if textLen >= width {
		return s
	}
	totalPad := width - textLen
	left := totalPad / 2
	right := totalPad - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
