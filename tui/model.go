package tui

import (
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

// Model is the central state for the Bubble Tea application.
type Model struct {
	mode              Mode
	labs              []labs.Lab
	selectedIdx       int
	containerID       string
	containerName     string
	activeLab         labs.Lab
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
	leftPanel      lipgloss.Style
	rightPanel     lipgloss.Style
	bottomBar      lipgloss.Style
	selectedItem   lipgloss.Style
	unselectedItem lipgloss.Style
	title          lipgloss.Style
	modeTUI        lipgloss.Style
	hint           lipgloss.Style
	goal           lipgloss.Style
	checkPassed    lipgloss.Style
	checkFailed    lipgloss.Style
	logStyle       lipgloss.Style
	border         lipgloss.Style
}

func defaultStyles() *styles {
	s := &styles{}

	blue := lipgloss.Color("#006BB6")
	green := lipgloss.Color("#50C878")
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

	s.selectedItem = lipgloss.NewStyle().
		Foreground(blue).
		Bold(true)

	s.unselectedItem = lipgloss.NewStyle().
		Foreground(gray)

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

// NewModel creates a new application model.
func NewModel(labsList []labs.Lab) *Model {
	m := &Model{
		mode:        ModeTUI,
		labs:        labsList,
		selectedIdx: 0,
		dockerCli:   docker.New(),
	}
	m.initStyles()
	return m
}

// NewModelWithState creates a new application model with restored state.
func NewModelWithState(labsList []labs.Lab, containerName, containerID string, activeLab labs.Lab, selectedIdx int) *Model {
	m := &Model{
		mode:          ModeTUI,
		labs:          labsList,
		selectedIdx:   selectedIdx,
		containerName: containerName,
		containerID:   containerID,
		activeLab:     activeLab,
		dockerCli:     docker.New(),
	}
	m.initStyles()
	return m
}

// GetState returns the current model state for restoration.
func (m *Model) GetState() (containerName, containerID string, activeLab labs.Lab, selectedIdx int) {
	return m.containerName, m.containerID, m.activeLab, m.selectedIdx
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

	case "up", "k":
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
		return m, nil

	case "down", "j":
		if m.selectedIdx < len(m.labs)-1 {
			m.selectedIdx++
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
		if m.containerName == "" {
			m.log("No container running. Start a lab first.")
			return m, nil
		}
		running, err := m.dockerCli.IsRunning(m.containerName)
		if err != nil || !running {
			if err != nil {
				m.log("Container check failed: " + err.Error())
			}
			m.log("Container not running — press Enter to start it first")
			return m, nil
		}
		m.log("Exiting TUI to open shell...")
		m.shellRequested = true
		return m, tea.Quit

	case "l":
		m.showLog = !m.showLog
		return m, nil
	}

	return m, nil
}

func (m *Model) handleDockerComplete(msg DockerOpComplete) {
	if msg.Err != nil && msg.Op != "check" {
		m.log("Docker error (" + msg.Op + "): " + msg.Err.Error())
	}
	if msg.ID != "" {
		m.containerID = msg.ID
	}

	if msg.Op == "check" {
		m.validationResults = append(m.validationResults, CheckResult{
			Command: msg.Command,
			Passed:  msg.Passed,
			Output:  msg.Out,
		})

		if len(m.validationResults) >= len(m.activeLab.Checks) {
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
}

func (m *Model) startLab() (tea.Model, tea.Cmd) {
	if len(m.labs) == 0 {
		m.log("No labs available")
		return m, nil
	}

	lab := m.labs[m.selectedIdx]
	m.activeLab = lab

	containerName := containerNamePrefix + m.normalizeName(lab.Name)
	m.containerName = containerName

	// Check if container already exists and is running
	exists, _ := m.dockerCli.ContainerExists(containerName)
	if exists {
		running, _ := m.dockerCli.IsRunning(containerName)
		if running {
			m.log("Container already running: " + containerName)
			m.containerName = containerName
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
	if m.containerName == "" {
		m.log("No lab is currently running")
		return m, nil
	}

	m.log("Resetting container: " + m.containerName)
	return m, func() tea.Msg {
		err := m.dockerCli.CleanUp(m.containerName)
		if err != nil {
			m.log("Failed to stop/remove container: " + err.Error())
			return DockerOpComplete{Op: "reset", Err: err}
		}

		lab := m.activeLab
		containerName := m.containerName
		m.log("Starting fresh container for: " + lab.Name)
		id, err := m.dockerCli.Start(lab.Image, containerName)
		if err != nil {
			m.log("Failed to start container: " + err.Error())
			return DockerOpComplete{Op: "reset", Err: err}
		}
		m.log("Container started: " + id[:12])
		return DockerOpComplete{Op: "reset", Err: err, ID: id}
	}
}

func (m *Model) stopLab() (tea.Model, tea.Cmd) {
	if m.containerName == "" {
		m.log("No lab is currently running")
		return m, nil
	}

	m.log("Stopping container: " + m.containerName)
	return m, func() tea.Msg {
		err := m.dockerCli.Stop(m.containerName)
		if err != nil {
			m.log("Failed to stop container: " + err.Error())
			return DockerOpComplete{Op: "stop", Err: err}
		}
		m.containerID = ""
		m.log("Container stopped")
		return DockerOpComplete{Op: "stop", Err: nil}
	}
}

func (m *Model) validateLab() (tea.Model, tea.Cmd) {
	if m.activeLab.Name == "" {
		m.log("No lab selected")
		return m, nil
	}

	var cmds []tea.Cmd
	for _, check := range m.activeLab.Checks {
		cmds = append(cmds, m.runCheck(check))
	}
	return m, tea.Batch(cmds...)
}

func (m *Model) runCheck(command string) tea.Cmd {
	return func() tea.Msg {
		passed, err := m.dockerCli.Validate(m.containerName, command)
		output := ""
		if err != nil {
			output = err.Error()
		}
		return DockerOpComplete{
			Op:      "check",
			Err:     err,
			Out:     output,
			Passed:  passed,
			Command: command,
		}
	}
}

// startPTYReader starts a goroutine that reads from the PTY and sends
// ptyMsg values to m.ptyCh. It runs until m.ptyDone is closed.
// This spawns a real goroutine — NOT a blocking tea.Cmd.
func (m *Model) log(msg string) {
	m.logBuffer = append(m.logBuffer, msg)
	if len(m.logBuffer) > 200 {
		m.logBuffer = m.logBuffer[len(m.logBuffer)-150:]
	}
}

func (m *Model) normalizeName(name string) string {
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

// RequestShell returns true if the TUI should exit and start a shell session.
func (m *Model) RequestShell() bool {
	return m.shellRequested
}

// ContainerName returns the name of the currently running container.
func (m *Model) ContainerName() string {
	return m.containerName
}

// truncate renders text to fit within n runes.
func truncate(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	runes := []rune(s)
	return string(runes[:n-1]) + "…"
}
