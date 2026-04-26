package pty

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"
)

// Session represents an active PTY session for a Docker exec shell.
type Session struct {
	f *os.File
}

// StartDockerExec starts a docker exec sh process inside a container and wraps it in a PTY.
func StartDockerExec(containerName string) (*Session, error) {
	cmd := exec.Command("docker", "exec", "-it", containerName, "sh")
	f, err := ptyStart(cmd)
	if err != nil {
		return nil, fmt.Errorf("starting pty: %w", err)
	}
	return &Session{f: f}, nil
}

// ExecuteShell starts an interactive docker exec shell session, forwarding
// the terminal directly to the PTY. It restores the terminal to normal mode
// when the shell exits (either via exit command or Ctrl+C).
func ExecuteShell(containerName string) error {
	// Get actual terminal size
	cols, rows, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		cols, rows = 80, 24
	}

	cmd := exec.Command("docker", "exec", "-it", containerName, "sh")
	f, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
	if err != nil {
		return fmt.Errorf("starting shell: %w", err)
	}

	// Put the real terminal in raw mode
	origState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		f.Close()
		return fmt.Errorf("setting raw mode: %w", err)
	}

	// Channels to signal when copies are done
	outputDone := make(chan struct{})
	inputDone := make(chan struct{})

	// Forward PTY output to stdout in goroutine
	go func() {
		io.Copy(os.Stdout, f)
		close(outputDone)
	}()

	// Forward stdin to the PTY in another goroutine
	go func() {
		io.Copy(f, os.Stdin)
		close(inputDone)
	}()

	// Wait for the command to finish (shell exits)
	_ = cmd.Wait()

	// Close the PTY to signal both goroutines to finish
	f.Close()

	// Wait for both goroutines to finish with timeout
	select {
	case <-outputDone:
	case <-time.After(200 * time.Millisecond):
	}
	select {
	case <-inputDone:
	case <-time.After(200 * time.Millisecond):
	}

	// Restore the terminal
	term.Restore(int(os.Stdin.Fd()), origState)
	fmt.Println()
	return nil
}

// Read reads all available output from the PTY.
func (s *Session) Read(p []byte) (int, error) {
	return s.f.Read(p)
}

// Write sends input bytes to the PTY.
func (s *Session) Write(p []byte) (int, error) {
	return s.f.Write(p)
}

// Close terminates the PTY session.
func (s *Session) Close() error {
	return s.f.Close()
}

// Fd returns the file descriptor for the PTY.
func (s *Session) Fd() uintptr {
	return s.f.Fd()
}

// File returns the underlying file.
func (s *Session) File() *os.File {
	return s.f
}

// ptyStart is an interfaceable wrapper around pty.Start.
func ptyStart(cmd *exec.Cmd) (*os.File, error) {
	return pty.StartWithSize(cmd, &pty.Winsize{
		Rows: 40,
		Cols: 120,
	})
}
