package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

// CopySetup copies lab setup files into the container.
func (c *Client) CopySetup(containerName, labDir string) error {
	// Copy setup directory
	setupSrc := filepath.Join(labDir, "setup")
	setupDst := "/setup"
	if _, err := os.Stat(setupSrc); err == nil {
		// First ensure the destination directory exists
		mkdirCmd := exec.Command("docker", "exec", containerName, "mkdir", "-p", "/setup")
		if err := mkdirCmd.Run(); err != nil {
			return fmt.Errorf("creating setup dir: %w", err)
		}
		// Then copy
		cmd := exec.Command("docker", "cp", setupSrc+"/.", containerName+":"+setupDst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("copying setup: %w, output: %s", err, string(out))
		}
	}

	// Copy validate directory
	valSrc := filepath.Join(labDir, "validate")
	valDst := "/validate"
	if _, err := os.Stat(valSrc); err == nil {
		mkdirCmd := exec.Command("docker", "exec", containerName, "mkdir", "-p", "/validate")
		if err := mkdirCmd.Run(); err != nil {
			return fmt.Errorf("creating validate dir: %w", err)
		}
		cmd := exec.Command("docker", "cp", valSrc+"/.", containerName+":"+valDst)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("copying validate: %w, output: %s", err, string(out))
		}
	}

	return nil
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

// Validate runs validation commands and returns true if all pass.
func (c *Client) Validate(containerName string, commands []string) (bool, string, error) {
	var output strings.Builder
	allPassed := true

	for _, cmd := range commands {
		out, err := c.execCommandWithOutput(containerName, cmd)
		output.WriteString(out)
		if err != nil {
			allPassed = false
			output.WriteString(fmt.Sprintf("\nError running validation: %v\n", err))
		}
	}

	return allPassed, output.String(), nil
}

func (c *Client) execCommand(containerName, command string) error {
	cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)
	return cmd.Run()
}

func (c *Client) execCommandWithOutput(containerName, command string) (string, error) {
	cmd := exec.Command("docker", "exec", containerName, "sh", "-c", command)
	out, err := cmd.CombinedOutput()
	return string(out), err
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

// IsRunning checks if a container is currently running.
func (c *Client) IsRunning(name string) (bool, error) {
	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", name)
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) == "true", nil
}
