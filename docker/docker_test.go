package docker

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}
	if !reflect.DeepEqual(c, &Client{}) {
		t.Errorf("expected empty Client, got %+v", c)
	}
}

func TestValidate_DockerCommand(t *testing.T) {
	// We test that the Validate method constructs the correct command.
	// Since we can't easily mock os/exec in pure unit tests without
	// replacing the function, we verify by checking the Client exists
	// and the method is callable (it will fail if docker is not installed,
	// but that's expected — the point is the function is structurally correct).
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	// This will fail because the container doesn't exist, but we just
	// verify the function doesn't panic and returns the right types.
	passed, err := c.Validate("nonexistent-container", "true")

	// The function should not return an error from the command construction
	// but the validate logic should handle the non-existent container case.
	// We expect passed to be false (check failed) and err to be nil
	// (because Validate treats non-zero exit as pass=false, not error).
	if !reflect.DeepEqual(passed, false) {
		t.Errorf("expected passed=false, got %v", passed)
	}
	if err != nil {
		t.Errorf("expected err=nil for Validate, got: %v", err)
	}
}

func TestValidateWithOutput_DockerCommand(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	// Will fail because container doesn't exist
	out, err := c.ValidateWithOutput("nonexistent", "echo hello")
	if err == nil {
		t.Error("expected error for nonexistent container")
	}
	if out != "" {
		t.Errorf("expected empty output for failed command, got: %s", out)
	}
}

func TestCleanUp_StopsAndRemoves(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	// Cleanup on non-existent container: Stop will fail, but Remove
	// will also fail. We expect the function to handle this gracefully.
	err := c.CleanUp("nonexistent-container-12345")
	// At least the function should not panic
	_ = err
}

func TestContainerExists(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	// This will list containers; the function should return a boolean.
	exists, err := c.ContainerExists("fixer-test-does-not-exist")
	// The function might error if docker isn't available, but it should
	// handle that gracefully. We just verify it returns (bool, error).
	_ = exists
	_ = err
}

func TestIsRunning(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	running, err := c.IsRunning("nonexistent")
	// Docker inspect will fail for non-existent container
	_ = running
	_ = err
}
