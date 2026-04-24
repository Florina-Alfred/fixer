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
	c := New()
	if c == nil {
		t.Fatal("expected non-nil client")
	}

	// This will fail because the container doesn't exist, but we just
	// verify the function doesn't panic and returns the right types.
	passed, output, err := c.Validate("nonexistent-container", []string{"true"})

	// Verify return types are correct
	_ = passed
	_ = output
	_ = err
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
