package engine

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/NadeeshaMedagama/branch_aware_ci/pkg/interfaces"
)

// PolicyEngine implements the IPolicyEngine interface
// Following Single Responsibility Principle
type PolicyEngine struct{}

// NewPolicyEngine creates a new instance of PolicyEngine
func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{}
}

// Evaluate implements IPolicyEngine.Evaluate
func (e *PolicyEngine) Evaluate(ctx context.Context, branchInfo *interfaces.BranchInfo, config *interfaces.Config) (*interfaces.Decision, error) {
	decision := &interfaces.Decision{
		BranchName: branchInfo.ShortName,
		BranchType: branchInfo.Type,
		Actions:    []string{},
		Variables:  make(map[string]string),
		Warnings:   []string{},
		Metadata:   branchInfo.Metadata,
	}

	// Find matching branch mapping
	mapping := e.findBestMapping(branchInfo.ShortName, config)
	if mapping == nil {
		decision.Environment = "development"
		decision.ShouldDeploy = false
		decision.Warnings = append(decision.Warnings, "No matching branch mapping found, using development environment")
	} else {
		decision.Environment = mapping.Environment
		decision.Actions = mapping.Actions
		decision.ShouldDeploy = e.shouldDeploy(mapping.Actions, branchInfo.ShortName, config)
	}

	// Apply environment configuration
	if envConfig, exists := config.Environments[decision.Environment]; exists {
		decision.RequiresApproval = envConfig.RequiresApproval
		for k, v := range envConfig.Variables {
			decision.Variables[k] = v
		}

		if !e.isBranchAllowed(branchInfo.ShortName, envConfig.AllowedBranches) {
			decision.Warnings = append(decision.Warnings,
				fmt.Sprintf("Branch %s may not be allowed to deploy to %s",
					branchInfo.ShortName, decision.Environment))
		}
	}

	// Apply policies
	e.applyPolicies(decision, branchInfo, config)

	return decision, nil
}

// ValidatePolicy implements IPolicyEngine.ValidatePolicy
func (e *PolicyEngine) ValidatePolicy(ctx context.Context, config *interfaces.Config) (bool, []string, error) {
	var errors []string

	// Validate environments
	if len(config.Environments) == 0 {
		errors = append(errors, "No environments defined")
	}

	// Validate branch mappings
	if len(config.BranchMappings) == 0 {
		errors = append(errors, "No branch mappings defined")
	}

	// Check for duplicate priorities
	priorities := make(map[int]bool)
	for _, mapping := range config.BranchMappings {
		if priorities[mapping.Priority] {
			errors = append(errors, fmt.Sprintf("Duplicate priority found: %d", mapping.Priority))
		}
		priorities[mapping.Priority] = true
	}

	return len(errors) == 0, errors, nil
}

// findBestMapping finds the best matching branch mapping based on priority
func (e *PolicyEngine) findBestMapping(branchName string, config *interfaces.Config) *interfaces.BranchMapping {
	var bestMatch *interfaces.BranchMapping
	highestPriority := -1

	for i := range config.BranchMappings {
		mapping := &config.BranchMappings[i]
		if e.matchesPattern(branchName, mapping.Pattern) {
			if mapping.Priority > highestPriority {
				highestPriority = mapping.Priority
				bestMatch = mapping
			}
		}
	}

	return bestMatch
}

// matchesPattern checks if branch name matches a pattern
func (e *PolicyEngine) matchesPattern(branchName, pattern string) bool {
	if branchName == pattern {
		return true
	}

	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(branchName, prefix+"/")
	}

	matched, _ := filepath.Match(pattern, branchName)
	return matched
}

// shouldDeploy determines if deployment should occur
func (e *PolicyEngine) shouldDeploy(actions []string, branchName string, config *interfaces.Config) bool {
	for _, action := range actions {
		if action == "deploy" {
			return true
		}
	}

	for _, autoBranch := range config.Policies.AutoDeployBranches {
		if branchName == autoBranch {
			return true
		}
	}

	return false
}

// isBranchAllowed checks if branch is allowed for the environment
func (e *PolicyEngine) isBranchAllowed(branchName string, allowedPatterns []string) bool {
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
func (e *PolicyEngine) applyPolicies(decision *interfaces.Decision, branchInfo *interfaces.BranchInfo, config *interfaces.Config) {
	// Check blocked patterns
	for _, blockedPattern := range config.Policies.BlockedBranchPatterns {
		if e.matchesPattern(branchInfo.ShortName, blockedPattern) {
			decision.ShouldDeploy = false
			decision.Warnings = append(decision.Warnings,
				fmt.Sprintf("Branch matches blocked pattern: %s", blockedPattern))
		}
	}

	// Add required actions based on policies
	if config.Policies.RequireTests && !contains(decision.Actions, "test") {
		decision.Actions = append(decision.Actions, "test")
	}

	if config.Policies.RequireCodeReview && branchInfo.IsProtected {
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
