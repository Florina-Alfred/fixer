package labs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_All(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid lab files
	valid1 := `name: Lab One
image: alpine:latest
goal: Do something
checks:
  - "test -f /tmp/solved"
hints:
  - "Check /tmp"
`
	valid2 := `name: Lab Two
image: nginx:latest
goal: Fix nginx
checks:
  - "curl -f http://localhost"
`
	if err := os.WriteFile(filepath.Join(tmpDir, "one.yml"), []byte(valid1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "two.yaml"), []byte(valid2), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a non-YAML file that should be skipped
	if err := os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("ignore"), 0644); err != nil {
		t.Fatal(err)
	}

	labs, err := LoadAll(tmpDir)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(labs) != 2 {
		t.Fatalf("expected 2 labs, got %d", len(labs))
	}

	if labs[0].Name != "Lab One" {
		t.Errorf("expected 'Lab One', got '%s'", labs[0].Name)
	}
	if labs[1].Name != "Lab Two" {
		t.Errorf("expected 'Lab Two', got '%s'", labs[1].Name)
	}
}

func TestLoad_MissingName(t *testing.T) {
	tmpDir := t.TempDir()
	data := `image: alpine:latest
`
	if err := os.WriteFile(filepath.Join(tmpDir, "bad.yml"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	labs, err := LoadAll(tmpDir)
	if err == nil {
		t.Fatal("expected error for missing name, got nil")
	}
	if len(labs) != 0 {
		t.Errorf("expected 0 labs, got %d", len(labs))
	}
}

func TestLoad_MissingImage(t *testing.T) {
	tmpDir := t.TempDir()
	data := `name: Broken Lab
`
	if err := os.WriteFile(filepath.Join(tmpDir, "bad.yml"), []byte(data), 0644); err != nil {
		t.Fatal(err)
	}

	labs, err := LoadAll(tmpDir)
	if err == nil {
		t.Fatal("expected error for missing image, got nil")
	}
	if len(labs) != 0 {
		t.Errorf("expected 0 labs, got %d", len(labs))
	}
}

func TestLoad_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	labs, err := LoadAll(tmpDir)
	if err != nil {
		t.Fatalf("expected no error for empty dir, got: %v", err)
	}
	if len(labs) != 0 {
		t.Errorf("expected 0 labs, got %d", len(labs))
	}
}

func TestLoad_NonexistentDir(t *testing.T) {
	_, err := LoadAll("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Fatal("expected error for nonexistent dir, got nil")
	}
}

func TestLoad_Single(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "lab-*.yml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	data := `name: Single Lab
image: redis:latest
goal: Fix redis config
validate:
  - "redis-cli ping"
hints:
  - "Check redis config"
  - "Look at error logs"
`
	if _, err := tmpFile.WriteString(data); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	lab, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if lab.Name != "Single Lab" {
		t.Errorf("expected 'Single Lab', got '%s'", lab.Name)
	}
	if lab.Image != "redis:latest" {
		t.Errorf("expected 'redis:latest', got '%s'", lab.Image)
	}
	if lab.Goal != "Fix redis config" {
		t.Errorf("expected 'Fix redis config', got '%s'", lab.Goal)
	}
	if len(lab.Validate) != 1 {
		t.Errorf("expected 1 validation, got %d", len(lab.Validate))
	}
	if len(lab.Hints) != 2 {
		t.Errorf("expected 2 hints, got %d", len(lab.Hints))
	}
}
