package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the branch-aware CI configuration
type Config struct {
	Environments   map[string]EnvironmentConfig `yaml:"environments"`
	BranchMappings []BranchMapping              `yaml:"branch_mappings"`
	Policies       PolicyConfig                 `yaml:"policies"`
}

// EnvironmentConfig defines settings for a specific environment
type EnvironmentConfig struct {
	Name             string            `yaml:"name"`
	RequiresApproval bool              `yaml:"requires_approval"`
	AllowedBranches  []string          `yaml:"allowed_branches"`
	Variables        map[string]string `yaml:"variables"`
	NotifyOnDeploy   bool              `yaml:"notify_on_deploy"`
}

// BranchMapping maps branch patterns to environments
type BranchMapping struct {
	Pattern     string   `yaml:"pattern"`
	Environment string   `yaml:"environment"`
	Actions     []string `yaml:"actions"`
	Priority    int      `yaml:"priority"`
}

// PolicyConfig defines CI/CD policies
type PolicyConfig struct {
	RequireTests          bool     `yaml:"require_tests"`
	RequireCodeReview     bool     `yaml:"require_code_review"`
	BlockedBranchPatterns []string `yaml:"blocked_branch_patterns"`
	AutoDeployBranches    []string `yaml:"auto_deploy_branches"`
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Environments: map[string]EnvironmentConfig{
			"production": {
				Name:             "production",
				RequiresApproval: true,
				AllowedBranches:  []string{"main", "master"},
				Variables: map[string]string{
					"ENV": "production",
				},
				NotifyOnDeploy: true,
			},
			"staging": {
				Name:             "staging",
				RequiresApproval: false,
				AllowedBranches:  []string{"staging", "develop"},
				Variables: map[string]string{
					"ENV": "staging",
				},
				NotifyOnDeploy: true,
			},
			"development": {
				Name:             "development",
				RequiresApproval: false,
				AllowedBranches:  []string{"feature/*", "bugfix/*", "hotfix/*"},
				Variables: map[string]string{
					"ENV": "development",
				},
				NotifyOnDeploy: false,
			},
		},
		BranchMappings: []BranchMapping{
			{Pattern: "main", Environment: "production", Actions: []string{"deploy", "notify"}, Priority: 100},
			{Pattern: "master", Environment: "production", Actions: []string{"deploy", "notify"}, Priority: 100},
			{Pattern: "staging", Environment: "staging", Actions: []string{"deploy", "notify"}, Priority: 90},
			{Pattern: "develop", Environment: "staging", Actions: []string{"deploy"}, Priority: 80},
			{Pattern: "release/*", Environment: "staging", Actions: []string{"deploy", "test"}, Priority: 85},
			{Pattern: "feature/*", Environment: "development", Actions: []string{"test"}, Priority: 50},
			{Pattern: "bugfix/*", Environment: "development", Actions: []string{"test"}, Priority: 50},
			{Pattern: "hotfix/*", Environment: "staging", Actions: []string{"test", "deploy"}, Priority: 70},
		},
		Policies: PolicyConfig{
			RequireTests:          true,
			RequireCodeReview:     true,
			BlockedBranchPatterns: []string{},
			AutoDeployBranches:    []string{"main", "master", "staging"},
		},
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	// If config path is empty, try default locations
	if configPath == "" {
		configPath = findConfigFile()
	}

	// If no config file found, use defaults
	if configPath == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// findConfigFile searches for config file in common locations
func findConfigFile() string {
	possiblePaths := []string{
		".branchci.yml",
		".branchci.yaml",
		".github/branchci.yml",
		".github/branchci.yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
