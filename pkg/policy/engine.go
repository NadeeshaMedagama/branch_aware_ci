package policy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/config"
	"github.com/nadeesha_medagama/branch-aware-ci/pkg/git"
)

// Decision represents a CI/CD decision based on branch analysis
type Decision struct {
	BranchName       string            `json:"branch_name" yaml:"branch_name"`
	BranchType       string            `json:"branch_type" yaml:"branch_type"`
	Environment      string            `json:"environment" yaml:"environment"`
	ShouldDeploy     bool              `json:"should_deploy" yaml:"should_deploy"`
	RequiresApproval bool              `json:"requires_approval" yaml:"requires_approval"`
	Actions          []string          `json:"actions" yaml:"actions"`
	Variables        map[string]string `json:"variables" yaml:"variables"`
	Warnings         []string          `json:"warnings,omitempty" yaml:"warnings,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Engine evaluates policies and makes CI/CD decisions
type Engine struct {
	config *config.Config
}

// NewEngine creates a new policy engine
func NewEngine(cfg *config.Config) *Engine {
	return &Engine{config: cfg}
}

// Evaluate evaluates the branch and returns a decision
func (e *Engine) Evaluate(branchInfo *git.BranchInfo) (*Decision, error) {
	decision := &Decision{
		BranchName: branchInfo.ShortName,
		BranchType: branchInfo.Type,
		Actions:    []string{},
		Variables:  make(map[string]string),
		Warnings:   []string{},
		Metadata:   branchInfo.Metadata,
	}

	// Find matching branch mapping
	mapping := e.findBestMapping(branchInfo.ShortName)
	if mapping == nil {
		decision.Environment = "development"
		decision.ShouldDeploy = false
		decision.Warnings = append(decision.Warnings, "No matching branch mapping found, using development environment")
	} else {
		decision.Environment = mapping.Environment
		decision.Actions = mapping.Actions
		decision.ShouldDeploy = e.shouldDeploy(mapping.Actions, branchInfo.ShortName)
	}

	// Apply environment configuration
	if envConfig, exists := e.config.Environments[decision.Environment]; exists {
		decision.RequiresApproval = envConfig.RequiresApproval
		for k, v := range envConfig.Variables {
			decision.Variables[k] = v
		}

		// Check if branch is allowed for this environment
		if !e.isBranchAllowed(branchInfo.ShortName, envConfig.AllowedBranches) {
			decision.Warnings = append(decision.Warnings,
				fmt.Sprintf("Branch %s may not be allowed to deploy to %s",
					branchInfo.ShortName, decision.Environment))
		}
	}

	// Apply policies
	e.applyPolicies(decision, branchInfo)

	return decision, nil
}

// findBestMapping finds the best matching branch mapping based on priority
func (e *Engine) findBestMapping(branchName string) *config.BranchMapping {
	var bestMatch *config.BranchMapping
	highestPriority := -1

	for i, mapping := range e.config.BranchMappings {
		if e.matchesPattern(branchName, mapping.Pattern) {
			if mapping.Priority > highestPriority {
				highestPriority = mapping.Priority
				bestMatch = &e.config.BranchMappings[i]
			}
		}
	}

	return bestMatch
}

// matchesPattern checks if branch name matches a pattern
func (e *Engine) matchesPattern(branchName, pattern string) bool {
	// Exact match
	if branchName == pattern {
		return true
	}

	// Wildcard pattern (e.g., "feature/*")
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(branchName, prefix+"/")
	}

	// Glob pattern
	matched, _ := filepath.Match(pattern, branchName)
	return matched
}

// shouldDeploy determines if deployment should occur
func (e *Engine) shouldDeploy(actions []string, branchName string) bool {
	// Check if "deploy" is in actions
	for _, action := range actions {
		if action == "deploy" {
			return true
		}
	}

	// Check auto-deploy branches
	for _, autoBranch := range e.config.Policies.AutoDeployBranches {
		if branchName == autoBranch {
			return true
		}
	}

	return false
}

// isBranchAllowed checks if branch is allowed for the environment
func (e *Engine) isBranchAllowed(branchName string, allowedPatterns []string) bool {
	if len(allowedPatterns) == 0 {
		return true
	}

	for _, pattern := range allowedPatterns {
		if e.matchesPattern(branchName, pattern) {
			return true
		}
	}

	return false
}

// applyPolicies applies policy rules to the decision
func (e *Engine) applyPolicies(decision *Decision, branchInfo *git.BranchInfo) {
	// Check blocked patterns
	for _, blockedPattern := range e.config.Policies.BlockedBranchPatterns {
		if e.matchesPattern(branchInfo.ShortName, blockedPattern) {
			decision.ShouldDeploy = false
			decision.Warnings = append(decision.Warnings,
				fmt.Sprintf("Branch matches blocked pattern: %s", blockedPattern))
		}
	}

	// Add required actions based on policies
	if e.config.Policies.RequireTests && !contains(decision.Actions, "test") {
		decision.Actions = append(decision.Actions, "test")
	}

	if e.config.Policies.RequireCodeReview && branchInfo.IsProtected {
		decision.RequiresApproval = true
	}
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
