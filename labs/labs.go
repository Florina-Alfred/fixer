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
	Name     string   `yaml:"name"`
	Image    string   `yaml:"image"`
	Goal     string   `yaml:"goal"`
	Setup    []string `yaml:"setup"`
	Checks   []string `yaml:"checks"`
	Hints    []string `yaml:"hints"`
	Category string   `yaml:"category"`
	Level    string   `yaml:"level"`
}

// LoadAll loads all lab files (*.yaml, *.yml) from the given directory.
func LoadAll(dir string) ([]Lab, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading lab directory %s: %w", dir, err)
	}

	var labs []Lab
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".yml") && !strings.HasSuffix(name, ".yaml") {
			continue
		}

		path := filepath.Join(dir, name)
		lab, err := Load(path)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", path, err)
		}
		labs = append(labs, lab)
	}

	sort.Slice(labs, func(i, j int) bool {
		return labs[i].Name < labs[j].Name
	})

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
