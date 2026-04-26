package labs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Lab represents a declarative training lab configuration.
type Lab struct {
	Name        string   `yaml:"name"`
	Image       string   `yaml:"image"`
	Goal        string   `yaml:"goal"`
	Setup       []string `yaml:"setup"`
	Validate    []string `yaml:"validate"`
	Hints       []string `yaml:"hints"`
	Category    string   `yaml:"category"`
	Level       string   `yaml:"level"`
	Description string   `yaml:"description"`
}

// LoadAll loads all labs from the given directory.
// Supports both flat YAML files and folder-based structure (category/lab/lab.yml).
func LoadAll(dir string) ([]Lab, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading lab directory %s: %w", dir, err)
	}

	var labs []Lab

	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(dir, name)

		if entry.IsDir() {
			// Check for folder-based structure: category/lab/lab.yml
			categoryLabs, err := loadCategoryDir(path)
			if err != nil {
				continue // Skip invalid category directories
			}
			labs = append(labs, categoryLabs...)
		} else if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			// Support flat YAML files for backward compatibility
			lab, err := Load(path)
			if err != nil {
				return nil, fmt.Errorf("loading %s: %w", path, err)
			}
			labs = append(labs, lab)
		}
	}

	sort.Slice(labs, func(i, j int) bool {
		return labs[i].Name < labs[j].Name
	})

	return labs, nil
}

// loadCategoryDir loads labs from a category directory (e.g., grep/, find/).
func loadCategoryDir(categoryDir string) ([]Lab, error) {
	entries, err := os.ReadDir(categoryDir)
	if err != nil {
		return nil, err
	}

	var labs []Lab
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Look for lab.yml or lab.yaml in the subdirectory
		labDir := filepath.Join(categoryDir, entry.Name())
		ymlPath := filepath.Join(labDir, "lab.yml")
		yamlPath := filepath.Join(labDir, "lab.yaml")

		var labPath string
		if _, err := os.Stat(ymlPath); err == nil {
			labPath = ymlPath
		} else if _, err := os.Stat(yamlPath); err == nil {
			labPath = yamlPath
		} else {
			continue // No lab file found
		}

		lab, err := Load(labPath)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", labPath, err)
		}
		labs = append(labs, lab)
	}

	return labs, nil
}

// Load reads a single lab file from the given path.
func Load(path string) (Lab, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Lab{}, fmt.Errorf("reading %s: %w", path, err)
	}

	var lab Lab
	if err := yaml.Unmarshal(data, &lab); err != nil {
		return Lab{}, fmt.Errorf("parsing %s: %w", path, err)
	}

	if lab.Name == "" {
		return Lab{}, fmt.Errorf("lab %s: missing name", path)
	}
	if lab.Image == "" {
		return Lab{}, fmt.Errorf("lab %s: missing image", path)
	}

	return lab, nil
}
