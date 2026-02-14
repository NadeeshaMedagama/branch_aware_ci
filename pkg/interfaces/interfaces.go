package interfaces

import (
	"context"
)

// BranchInfo represents information about a Git branch
type BranchInfo struct {
	Name        string
	ShortName   string
	Type        string
	Metadata    map[string]string
	IsProtected bool
}

// Decision represents a CI/CD decision
type Decision struct {
	BranchName       string
	BranchType       string
	Environment      string
	ShouldDeploy     bool
	RequiresApproval bool
	Actions          []string
	Variables        map[string]string
	Warnings         []string
	Metadata         map[string]string
}

// Config represents the application configuration
type Config struct {
	Environments   map[string]EnvironmentConfig
	BranchMappings []BranchMapping
	Policies       PolicyConfig
}

// EnvironmentConfig defines settings for a specific environment
type EnvironmentConfig struct {
	Name             string
	RequiresApproval bool
	AllowedBranches  []string
	Variables        map[string]string
	NotifyOnDeploy   bool
}

// BranchMapping maps branch patterns to environments
type BranchMapping struct {
	Pattern     string
	Environment string
	Actions     []string
	Priority    int
}

// PolicyConfig defines CI/CD policies
type PolicyConfig struct {
	RequireTests          bool
	RequireCodeReview     bool
	BlockedBranchPatterns []string
	AutoDeployBranches    []string
}

// IBranchDetector defines the interface for branch detection
type IBranchDetector interface {
	DetectBranch(ctx context.Context, repoPath string) (*BranchInfo, error)
	GetBranchInfo(ctx context.Context, repoPath string, branchName string) (*BranchInfo, error)
}

// IPolicyEngine defines the interface for policy evaluation
type IPolicyEngine interface {
	Evaluate(ctx context.Context, branchInfo *BranchInfo, config *Config) (*Decision, error)
	ValidatePolicy(ctx context.Context, config *Config) (bool, []string, error)
}

// IConfigManager defines the interface for configuration management
type IConfigManager interface {
	GetConfig(ctx context.Context, configPath string) (*Config, error)
	UpdateConfig(ctx context.Context, config *Config, configPath string) error
	ValidateConfig(ctx context.Context, config *Config) (bool, []string, error)
}

// INotifier defines the interface for notifications
type INotifier interface {
	Notify(ctx context.Context, decision *Decision) error
}

// IHealthChecker defines the interface for health checks
type IHealthChecker interface {
	HealthCheck(ctx context.Context) error
	ReadinessCheck(ctx context.Context) error
}
