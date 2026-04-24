package tui

import (
	"errors"
	"strings"
	"testing"

	"redalf.de/fixer/labs"
)

func TestNewModel(t *testing.T) {
	m := NewModel([]labs.Lab{})
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
		{Name: "Lab A", Image: "alpine:latest", Goal: "Do stuff"},
		{Name: "Lab B", Image: "nginx:latest", Goal: "Fix nginx"},
	}
	m := NewModel(labsList)
	if len(m.labs) != 2 {
		t.Errorf("expected 2 labs, got %d", len(m.labs))
	}
	if m.selectedIdx != 0 {
		t.Errorf("expected selectedIdx=0, got %d", m.selectedIdx)
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
			m := NewModel(nil)
			result := m.normalizeName(tt.input)
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
	m := NewModel(nil)

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
	m := NewModel(nil)

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
	m := NewModel(nil)

	msg := DockerOpComplete{
		Op:   "start",
		Err:  nil,
		ID:   "abc123def456",
		Out:  "",
		Passed: false,
	}
	m.handleDockerComplete(msg)

	if m.containerID != "abc123def456" {
		t.Errorf("expected containerID 'abc123def456', got %q", m.containerID)
	}
	if len(m.logBuffer) != 0 {
		t.Errorf("expected no log entries, got %d", len(m.logBuffer))
	}
}

func TestHandleDockerComplete_StartError(t *testing.T) {
	m := NewModel(nil)

	msg := DockerOpComplete{
		Op:     "start",
		Err:    errors.New("docker not found"),
		ID:     "",
		Out:    "",
		Passed: false,
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
	m := NewModel(nil)
	m.activeLab = labs.Lab{
		Name:   "Test Lab",
		Checks: []string{"check1", "check2"},
	}

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
	m := NewModel(nil)
	m.activeLab = labs.Lab{
		Name:   "Test Lab",
		Checks: []string{"check1", "check2"},
	}

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
	m := NewModel(nil)
	m.activeLab = labs.Lab{
		Name:   "Test Lab",
		Checks: []string{"test -f /tmp/solved", "curl -f http://localhost"},
	}

	for _, check := range m.activeLab.Checks {
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
		if cr.Command != m.activeLab.Checks[i] {
			t.Errorf("check %d: expected command %q, got %q", i, m.activeLab.Checks[i], cr.Command)
		}
	}
}

func TestHandleDockerComplete_CheckNotAllDone(t *testing.T) {
	m := NewModel(nil)
	m.activeLab = labs.Lab{
		Name:   "Test Lab",
		Checks: []string{"a", "b", "c"},
	}

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
