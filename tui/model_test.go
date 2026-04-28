package tui

import (
	"errors"
	"strings"
	"testing"

	"redalf.de/fixer/labs"
)

func TestNewModel(t *testing.T) {
	m := NewModel([]labs.Lab{}, "")
	if m == nil {
		t.Fatal("expected non-nil model")
	}
	if m.mode != ModeTUI {
		t.Errorf("expected ModeTUI, got %v", m.mode)
	}
	if m.dockerCli == nil {
		t.Fatal("expected non-nil dockerCli")
	}
}

func TestNewModel_WithLabs(t *testing.T) {
	labsList := []labs.Lab{
		{Name: "Lab A", Image: "alpine:latest", Goal: "Do stuff", Category: "grep"},
		{Name: "Lab B", Image: "nginx:latest", Goal: "Fix nginx", Category: "grep"},
	}
	m := NewModel(labsList, "")
	if len(m.toolGroups) != 1 {
		t.Errorf("expected 1 tool group, got %d", len(m.toolGroups))
	}
	if len(m.toolGroups[0].Labs) != 2 {
		t.Errorf("expected 2 labs in first group, got %d", len(m.toolGroups[0].Labs))
	}
}

func TestNormalizeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Fix Nginx", "fix-nginx"},
		{"Hello World 123", "hello-world-123"},
		{"simple", "simple"},
		{"ALL CAPS", "all-caps"},
		{"with spaces", "with-spaces"},
		{"with/slashes", "with-slashes"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeName(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello world", 5, "hell…"},
		{"hello world", 3, "he…"},
		{"", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := truncate(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.n, result, tt.expected)
			}
		})
	}
}

func TestLog(t *testing.T) {
	m := NewModel(nil, "")

	m.log("msg1")
	m.log("msg2")
	m.log("msg3")

	if len(m.logBuffer) != 3 {
		t.Errorf("expected 3 log entries, got %d", len(m.logBuffer))
	}
	if m.logBuffer[0] != "msg1" {
		t.Errorf("expected 'msg1', got %q", m.logBuffer[0])
	}
	if m.logBuffer[2] != "msg3" {
		t.Errorf("expected 'msg3', got %q", m.logBuffer[2])
	}
}

func TestLog_CapsAtMax(t *testing.T) {
	m := NewModel(nil, "")

	for i := 0; i < 300; i++ {
		m.log("log entry")
	}

	if len(m.logBuffer) > 200 {
		t.Errorf("expected at most 200 log entries, got %d", len(m.logBuffer))
	}
	if len(m.logBuffer) < 150 {
		t.Errorf("expected at least 150 log entries, got %d", len(m.logBuffer))
	}
}

func TestHandleDockerComplete_Start(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Test Lab", Image: "alpine:latest", Category: "grep"},
	}, "")

	msg := DockerOpComplete{
		Op:      "start",
		Err:     nil,
		ID:      "abc123def456",
		Out:     "",
		Passed:  false,
		Command: "",
	}
	m.handleDockerComplete(msg)

	if len(m.logBuffer) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(m.logBuffer))
	}
	if !strings.Contains(m.logBuffer[0], "ready") {
		t.Errorf("expected 'ready' in log, got %q", m.logBuffer[0])
	}
}

func TestHandleDockerComplete_StartError(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Test Lab", Image: "alpine:latest", Category: "grep"},
	}, "")

	msg := DockerOpComplete{
		Op:      "start",
		Err:     errors.New("docker not found"),
		ID:      "",
		Out:     "",
		Passed:  false,
		Command: "",
	}
	m.handleDockerComplete(msg)

	if len(m.logBuffer) != 1 {
		t.Errorf("expected 1 log entry, got %d", len(m.logBuffer))
	}
	if !strings.Contains(m.logBuffer[0], "Docker error") {
		t.Errorf("expected 'Docker error' in log, got %q", m.logBuffer[0])
	}
}

func TestHandleDockerComplete_Check(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Test Lab", Image: "alpine:latest", Category: "grep", Validate: []string{"check1", "check2"}},
	}, "")

	// First check passes
	m.handleDockerComplete(DockerOpComplete{
		Op:      "check",
		Command: "check1",
		Passed:  true,
		Out:     "output1",
	})
	if len(m.validationResults) != 1 {
		t.Errorf("expected 1 validation result, got %d", len(m.validationResults))
	}
	if len(m.logBuffer) != 0 {
		t.Errorf("expected no log entries yet, got %d", len(m.logBuffer))
	}
	if m.lastValidation != nil {
		t.Error("expected no lastValidation until all checks done")
	}

	// Second check passes — all done
	m.handleDockerComplete(DockerOpComplete{
		Op:      "check",
		Command: "check2",
		Passed:  true,
		Out:     "output2",
	})
	if len(m.validationResults) != 0 {
		t.Error("expected validationResults to be cleared after all checks")
	}
	if m.lastValidation == nil {
		t.Fatal("expected lastValidation to be set")
	}
	if !m.lastValidation.Passed {
		t.Error("expected all checks to pass")
	}
	if m.lastValidation.Checks[0].Command != "check1" {
		t.Errorf("expected check1 command, got %q", m.lastValidation.Checks[0].Command)
	}
}

func TestHandleDockerComplete_CheckFail(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Test Lab", Image: "alpine:latest", Category: "grep", Validate: []string{"check1", "check2"}},
	}, "")

	// Both fail
	m.handleDockerComplete(DockerOpComplete{
		Op:      "check",
		Command: "check1",
		Passed:  false,
		Out:     "failed output",
	})
	m.handleDockerComplete(DockerOpComplete{
		Op:      "check",
		Command: "check2",
		Passed:  true,
	})

	if m.lastValidation == nil {
		t.Fatal("expected lastValidation")
	}
	if m.lastValidation.Passed {
		t.Error("expected validation to fail")
	}
	if m.lastValidation.Checks[0].Command != "check1" {
		t.Errorf("expected check1 command, got %q", m.lastValidation.Checks[0].Command)
	}
	if m.lastValidation.Checks[0].Passed {
		t.Error("expected first check to have Passed=false")
	}
}

func TestHandleDockerComplete_CheckCommandAssociation(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Test Lab", Image: "alpine:latest", Category: "grep", Validate: []string{"test -f /tmp/solved", "curl -f http://localhost"}},
	}, "")

	for _, check := range m.toolGroups[0].Labs[0].Lab.Validate {
		m.handleDockerComplete(DockerOpComplete{
			Op:      "check",
			Command: check,
			Passed:  false,
		})
	}

	if m.lastValidation == nil {
		t.Fatal("expected lastValidation")
	}
	for i, cr := range m.lastValidation.Checks {
		expectedCheck := m.toolGroups[0].Labs[0].Lab.Validate[i]
		if cr.Command != expectedCheck {
			t.Errorf("check %d: expected command %q, got %q", i, expectedCheck, cr.Command)
		}
	}
}

func TestHandleDockerComplete_CheckNotAllDone(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Test Lab", Image: "alpine:latest", Category: "grep", Validate: []string{"a", "b", "c"}},
	}, "")

	// Only 2 of 3 checks done
	m.handleDockerComplete(DockerOpComplete{
		Op:      "check",
		Command: "a",
		Passed:  true,
	})
	m.handleDockerComplete(DockerOpComplete{
		Op:      "check",
		Command: "b",
		Passed:  false,
	})

	if len(m.validationResults) != 2 {
		t.Errorf("expected 2 pending results, got %d", len(m.validationResults))
	}
	if m.lastValidation != nil {
		t.Error("expected lastValidation to be nil until all checks complete")
	}
}

func TestContainerStates(t *testing.T) {
	m := NewModel([]labs.Lab{
		{Name: "Lab A", Image: "alpine:latest", Category: "grep"},
	}, "")

	// Initial state should be stopped
	if m.toolGroups[0].Labs[0].State != StateStopped {
		t.Errorf("expected StateStopped, got %v", m.toolGroups[0].Labs[0].State)
	}

	// Set to active
	m.setContainerActive(true)
	if m.toolGroups[0].Labs[0].State != StateActive {
		t.Errorf("expected StateActive, got %v", m.toolGroups[0].Labs[0].State)
	}

	// Set to idle
	m.setContainerActive(false)
	if m.toolGroups[0].Labs[0].State != StateIdle {
		t.Errorf("expected StateIdle, got %v", m.toolGroups[0].Labs[0].State)
	}

	// Set to stopped
	m.setContainerStateStopped()
	if m.toolGroups[0].Labs[0].State != StateStopped {
		t.Errorf("expected StateStopped, got %v", m.toolGroups[0].Labs[0].State)
	}
}

func TestGroupLabsByCategory(t *testing.T) {
	labsList := []labs.Lab{
		{Name: "Lab A", Image: "alpine:latest", Category: "grep"},
		{Name: "Lab B", Image: "nginx:latest", Category: "find"},
		{Name: "Lab C", Image: "python:latest", Category: "grep"},
	}

	groups := groupLabsByCategory(labsList)

	if len(groups) != 2 {
		t.Errorf("expected 2 groups, got %d", len(groups))
	}

	// Find grep group
	var grepGroup *ToolGroup
	for i := range groups {
		if groups[i].Category == "grep" {
			grepGroup = &groups[i]
			break
		}
	}
	if grepGroup == nil {
		t.Fatal("grep group not found")
	}
	if len(grepGroup.Labs) != 2 {
		t.Errorf("expected 2 labs in grep group, got %d", len(grepGroup.Labs))
	}
}
