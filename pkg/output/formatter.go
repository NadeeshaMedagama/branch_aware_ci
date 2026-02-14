package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/policy"
)

// Format represents the output format
type Format string

const (
	FormatJSON         Format = "json"
	FormatYAML         Format = "yaml"
	FormatEnv          Format = "env"
	FormatGitHubEnv    Format = "github-env"
	FormatGitHubOutput Format = "github-output"
	FormatHuman        Format = "human"
)

// Formatter handles output formatting
type Formatter struct {
	format Format
}

// NewFormatter creates a new formatter
func NewFormatter(format Format) *Formatter {
	return &Formatter{format: format}
}

// Format formats the decision according to the specified format
func (f *Formatter) Format(decision *policy.Decision) (string, error) {
	switch f.format {
	case FormatJSON:
		return f.formatJSON(decision)
	case FormatYAML:
		return f.formatYAML(decision)
	case FormatEnv:
		return f.formatEnv(decision), nil
	case FormatGitHubEnv:
		return f.formatGitHubEnv(decision)
	case FormatGitHubOutput:
		return f.formatGitHubOutput(decision)
	case FormatHuman:
		return f.formatHuman(decision), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", f.format)
	}
}

// formatJSON formats as JSON
func (f *Formatter) formatJSON(decision *policy.Decision) (string, error) {
	data, err := json.MarshalIndent(decision, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(data), nil
}

// formatYAML formats as YAML
func (f *Formatter) formatYAML(decision *policy.Decision) (string, error) {
	data, err := yaml.Marshal(decision)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return string(data), nil
}

// formatEnv formats as environment variables
func (f *Formatter) formatEnv(decision *policy.Decision) string {
	var lines []string

	lines = append(lines, fmt.Sprintf("BRANCH_NAME=%s", decision.BranchName))
	lines = append(lines, fmt.Sprintf("BRANCH_TYPE=%s", decision.BranchType))
	lines = append(lines, fmt.Sprintf("ENVIRONMENT=%s", decision.Environment))
	lines = append(lines, fmt.Sprintf("SHOULD_DEPLOY=%t", decision.ShouldDeploy))
	lines = append(lines, fmt.Sprintf("REQUIRES_APPROVAL=%t", decision.RequiresApproval))

	if len(decision.Actions) > 0 {
		lines = append(lines, fmt.Sprintf("ACTIONS=%s", strings.Join(decision.Actions, ",")))
	}

	for k, v := range decision.Variables {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	return strings.Join(lines, "\n")
}

// formatGitHubEnv formats for GitHub Actions environment file
func (f *Formatter) formatGitHubEnv(decision *policy.Decision) (string, error) {
	envFile := os.Getenv("GITHUB_ENV")
	if envFile == "" {
		return "", fmt.Errorf("GITHUB_ENV not set")
	}

	file, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open GITHUB_ENV file: %w", err)
	}
	defer file.Close()

	envContent := f.formatEnv(decision)
	if _, err := file.WriteString(envContent + "\n"); err != nil {
		return "", fmt.Errorf("failed to write to GITHUB_ENV: %w", err)
	}

	return "Environment variables written to $GITHUB_ENV", nil
}

// formatGitHubOutput formats for GitHub Actions output
func (f *Formatter) formatGitHubOutput(decision *policy.Decision) (string, error) {
	outputFile := os.Getenv("GITHUB_OUTPUT")
	if outputFile == "" {
		return "", fmt.Errorf("GITHUB_OUTPUT not set")
	}

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to open GITHUB_OUTPUT file: %w", err)
	}
	defer file.Close()

	var lines []string
	lines = append(lines, fmt.Sprintf("branch_name=%s", decision.BranchName))
	lines = append(lines, fmt.Sprintf("branch_type=%s", decision.BranchType))
	lines = append(lines, fmt.Sprintf("environment=%s", decision.Environment))
	lines = append(lines, fmt.Sprintf("should_deploy=%t", decision.ShouldDeploy))
	lines = append(lines, fmt.Sprintf("requires_approval=%t", decision.RequiresApproval))
	lines = append(lines, fmt.Sprintf("actions=%s", strings.Join(decision.Actions, ",")))

	output := strings.Join(lines, "\n") + "\n"
	if _, err := file.WriteString(output); err != nil {
		return "", fmt.Errorf("failed to write to GITHUB_OUTPUT: %w", err)
	}

	return "Output variables written to $GITHUB_OUTPUT", nil
}

// formatHuman formats for human-readable output
func (f *Formatter) formatHuman(decision *policy.Decision) string {
	var lines []string

	lines = append(lines, "🌿 Branch Analysis")
	lines = append(lines, "==================")
	lines = append(lines, fmt.Sprintf("Branch:      %s", decision.BranchName))
	lines = append(lines, fmt.Sprintf("Type:        %s", decision.BranchType))
	lines = append(lines, fmt.Sprintf("Environment: %s", decision.Environment))
	lines = append(lines, "")

	lines = append(lines, "📋 CI/CD Decision")
	lines = append(lines, "==================")

	if decision.ShouldDeploy {
		lines = append(lines, "✅ Should Deploy: Yes")
	} else {
		lines = append(lines, "❌ Should Deploy: No")
	}

	if decision.RequiresApproval {
		lines = append(lines, "⚠️  Requires Approval: Yes")
	} else {
		lines = append(lines, "✓  Requires Approval: No")
	}

	if len(decision.Actions) > 0 {
		lines = append(lines, fmt.Sprintf("Actions:     %s", strings.Join(decision.Actions, ", ")))
	}

	if len(decision.Variables) > 0 {
		lines = append(lines, "")
		lines = append(lines, "🔧 Variables")
		lines = append(lines, "============")
		for k, v := range decision.Variables {
			lines = append(lines, fmt.Sprintf("%s=%s", k, v))
		}
	}

	if len(decision.Warnings) > 0 {
		lines = append(lines, "")
		lines = append(lines, "⚠️  Warnings")
		lines = append(lines, "===========")
		for _, warning := range decision.Warnings {
			lines = append(lines, fmt.Sprintf("- %s", warning))
		}
	}

	if len(decision.Metadata) > 0 {
		lines = append(lines, "")
		lines = append(lines, "📊 Metadata")
		lines = append(lines, "===========")
		for k, v := range decision.Metadata {
			lines = append(lines, fmt.Sprintf("%s: %s", k, v))
		}
	}

	return strings.Join(lines, "\n")
}
