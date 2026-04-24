package docker

import (
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps Docker CLI operations used by the fixer tool.
type Client struct{}

// New creates a new Docker client.
func New() *Client {
	return &Client{}
}

// ContainerExists checks if a container with the given name exists (stopped or running).
func (c *Client) ContainerExists(name string) (bool, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("checking containers: %w", err)
	}

	names := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, n := range names {
		if n == name {
			return true, nil
		}
	}
	return false, nil
}

// Start launches a container from the given image with a generated name.
func (c *Client) Start(image, name string) (string, error) {
	cmd := exec.Command("docker", "run", "-d", "--name", name, image, "tail", "-f", "/dev/null")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("starting container: %w", err)
	}

	return strings.TrimSpace(string(out)), nil
}

// Setup runs setup commands inside a container after it starts.
func (c *Client) Setup(containerName string, commands []string) error {
	for _, cmd := range commands {
		if err := c.execCommand(containerName, cmd); err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
	}
	return nil
}

func (c *Client) execCommand(containerName, command string) error {
	cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)
	return cmd.Run()
}

// Stop stops a running container.
func (c *Client) Stop(name string) error {
	cmd := exec.Command("docker", "stop", name, "--time", "5")
	return cmd.Run()
}

// Remove removes a container completely.
func (c *Client) Remove(name string) error {
	cmd := exec.Command("docker", "rm", "-f", name)
	return cmd.Run()
}

// CleanUp removes a container regardless of its state.
func (c *Client) CleanUp(name string) error {
	_ = c.Stop(name)
	return c.Remove(name)
}

// Validate runs the given command inside the container and returns its exit status.
// A return code of 0 means the check passed.
func (c *Client) Validate(containerName, command string) (bool, error) {
	cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		// Non-zero exit means the check failed.
		return false, nil
	}
	return true, nil
}

// ValidateWithOutput runs the given command inside the container and returns stdout.
func (c *Client) ValidateWithOutput(containerName, command string) (string, error) {
	cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("validation failed: %s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// IsRunning checks if a container is currently running.
func (c *Client) IsRunning(name string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", name)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "true", nil
}
