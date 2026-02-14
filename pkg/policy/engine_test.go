package policy

import (
	"testing"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/config"
	"github.com/nadeesha_medagama/branch-aware-ci/pkg/git"
)

func TestEvaluate(t *testing.T) {
	cfg := config.DefaultConfig()
	engine := NewEngine(cfg)

	tests := []struct {
		name             string
		branchName       string
		branchType       string
		expectedEnv      string
		expectedDeploy   bool
		expectedApproval bool
	}{
		{
			name:             "main branch",
			branchName:       "main",
			branchType:       "main",
			expectedEnv:      "production",
			expectedDeploy:   true,
			expectedApproval: true,
		},
		{
			name:             "staging branch",
			branchName:       "staging",
			branchType:       "staging",
			expectedEnv:      "staging",
			expectedDeploy:   true,
			expectedApproval: false,
		},
		{
			name:             "feature branch",
			branchName:       "feature/user-auth",
			branchType:       "feature",
			expectedEnv:      "development",
			expectedDeploy:   false,
			expectedApproval: false,
		},
		{
			name:             "hotfix branch",
			branchName:       "hotfix/critical-bug",
			branchType:       "hotfix",
			expectedEnv:      "staging",
			expectedDeploy:   true,
			expectedApproval: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branchInfo := &git.BranchInfo{
				ShortName:   tt.branchName,
				Type:        tt.branchType,
				Metadata:    make(map[string]string),
				IsProtected: tt.branchType == "main",
			}

			decision, err := engine.Evaluate(branchInfo)
			if err != nil {
				t.Fatalf("Evaluate failed: %v", err)
			}

			if decision.Environment != tt.expectedEnv {
				t.Errorf("Expected environment %s, got %s", tt.expectedEnv, decision.Environment)
			}

			if decision.ShouldDeploy != tt.expectedDeploy {
				t.Errorf("Expected ShouldDeploy %v, got %v", tt.expectedDeploy, decision.ShouldDeploy)
			}

			if decision.RequiresApproval != tt.expectedApproval {
				t.Errorf("Expected RequiresApproval %v, got %v", tt.expectedApproval, decision.RequiresApproval)
			}
		})
	}
}

func TestMatchesPattern(t *testing.T) {
	engine := NewEngine(config.DefaultConfig())

	tests := []struct {
		name       string
		branchName string
		pattern    string
		expected   bool
	}{
		{"exact match", "main", "main", true},
		{"wildcard match", "feature/auth", "feature/*", true},
		{"wildcard no match", "bugfix/auth", "feature/*", false},
		{"glob match", "release-v1.0", "release-*", true},
		{"no match", "develop", "main", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.matchesPattern(tt.branchName, tt.pattern)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for pattern %s on branch %s",
					tt.expected, result, tt.pattern, tt.branchName)
			}
		})
	}
}

func TestFindBestMapping(t *testing.T) {
	cfg := config.DefaultConfig()
	engine := NewEngine(cfg)

	tests := []struct {
		name        string
		branchName  string
		expectedEnv string
	}{
		{"main branch", "main", "production"},
		{"master branch", "master", "production"},
		{"staging branch", "staging", "staging"},
		{"feature branch", "feature/test", "development"},
		{"release branch", "release/v1.0", "staging"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := engine.findBestMapping(tt.branchName)
			if mapping == nil {
				t.Fatalf("No mapping found for %s", tt.branchName)
			}

			if mapping.Environment != tt.expectedEnv {
				t.Errorf("Expected environment %s, got %s", tt.expectedEnv, mapping.Environment)
			}
		})
	}
}
