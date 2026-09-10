package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFile = "agent-kit.yaml"

// InitLogFile is written at the project root after every successful setup
// run. It records what was installed and how long the run took so setups
// can be benchmarked and compared over time.
const InitLogFile = "agent-kit-init-log.md"

// ProjectConfig stores the project's agent-kit configuration.
// Written by `agent-kit init`, read by `agent-kit add`.
type ProjectConfig struct {
	ProjectName     string   `yaml:"projectName"`
	Prefix          string   `yaml:"prefix"`
	Stack           []string `yaml:"stack"`
	PackageManager  string   `yaml:"packageManager"`
	Monorepo        bool     `yaml:"monorepo"`
	TicketTracker   string   `yaml:"ticketTracker"`
	FeatureFlow     []string `yaml:"featureFlow"`
	IncludeOpenSpec bool     `yaml:"includeOpenSpec"`
	IncludeClaude   bool     `yaml:"includeClaude"`
	IncludeHooks    bool     `yaml:"includeHooks"`
	IncludeCI       bool     `yaml:"includeCI"`
	OptionalSkills  []string `yaml:"optionalSkills"`
	InstalledSkills []string `yaml:"installedSkills"`
}

// Load reads the project config from the current directory.
func Load() (*ProjectConfig, error) {
	return LoadFrom(".")
}

// LoadFrom reads the project config from the given directory.
func LoadFrom(dir string) (*ProjectConfig, error) {
	path := filepath.Join(dir, ConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg ProjectConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save writes the project config to the given directory.
func (c *ProjectConfig) Save(dir string) error {
	path := filepath.Join(dir, ConfigFile)
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// HasSkill checks if a skill is already installed.
func (c *ProjectConfig) HasSkill(name string) bool {
	for _, s := range c.InstalledSkills {
		if s == name {
			return true
		}
	}
	return false
}

// AddSkill records a skill as installed.
func (c *ProjectConfig) AddSkill(name string) {
	if !c.HasSkill(name) {
		c.InstalledSkills = append(c.InstalledSkills, name)
	}
}
