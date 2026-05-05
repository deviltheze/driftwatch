package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Host represents a remote server to monitor for config drift.
type Host struct {
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	KeyPath  string `yaml:"key_path"`
}

// Check defines a single file or command to watch on remote hosts.
type Check struct {
	Name    string   `yaml:"name"`
	Path    string   `yaml:"path"`
	Hosts   []string `yaml:"hosts"`
}

// Config is the top-level configuration for driftwatch.
type Config struct {
	Interval time.Duration `yaml:"interval"`
	Hosts    []Host        `yaml:"hosts"`
	Checks   []Check       `yaml:"checks"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if cfg.Interval == 0 {
		cfg.Interval = 60 * time.Second
	}

	for i := range cfg.Hosts {
		if cfg.Hosts[i].Port == 0 {
			cfg.Hosts[i].Port = 22
		}
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if len(c.Hosts) == 0 {
		return fmt.Errorf("at least one host must be defined")
	}
	for _, h := range c.Hosts {
		if h.Name == "" {
			return fmt.Errorf("host is missing a name")
		}
		if h.Address == "" {
			return fmt.Errorf("host %q is missing an address", h.Name)
		}
		if h.User == "" {
			return fmt.Errorf("host %q is missing a user", h.Name)
		}
	}
	for _, ch := range c.Checks {
		if ch.Name == "" {
			return fmt.Errorf("a check is missing a name")
		}
		if ch.Path == "" {
			return fmt.Errorf("check %q is missing a path", ch.Name)
		}
	}
	return nil
}
