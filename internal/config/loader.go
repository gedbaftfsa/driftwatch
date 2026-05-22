package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ServiceDefinition represents a declared infrastructure definition for a service.
type ServiceDefinition struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Replicas    int               `yaml:"replicas"`
	Environment map[string]string `yaml:"environment"`
	Ports       []int             `yaml:"ports"`
	Resources   ResourceLimits    `yaml:"resources"`
}

// ResourceLimits holds CPU and memory constraints.
type ResourceLimits struct {
	CPU    string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

// DriftWatchConfig is the top-level configuration file structure.
type DriftWatchConfig struct {
	Version  string              `yaml:"version"`
	Services []ServiceDefinition `yaml:"services"`
}

// LoadConfig reads and parses a driftwatch YAML config file from the given path.
func LoadConfig(path string) (*DriftWatchConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg DriftWatchConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// Validate performs basic sanity checks on the loaded configuration.
func (c *DriftWatchConfig) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("config version must be specified")
	}
	seen := make(map[string]struct{})
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service at index %d has no name", i)
		}
		if _, dup := seen[svc.Name]; dup {
			return fmt.Errorf("duplicate service name %q", svc.Name)
		}
		seen[svc.Name] = struct{}{}
	}
	return nil
}
