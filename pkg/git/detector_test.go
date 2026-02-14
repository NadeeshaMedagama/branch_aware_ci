package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func TestDetectBranch(t *testing.T) {
	// Create a temporary directory for test repository
	tmpDir, err := os.MkdirTemp("", "branch-aware-ci-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize a git repository
	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Add and commit
	worktree, _ := repo.Worktree()
	worktree.Add("test.txt")
	worktree.Commit("Initial commit", &git.CommitOptions{})

	// Test detection
	detector := NewDetector(tmpDir)
	branchInfo, err := detector.DetectBranch()
	if err != nil {
		t.Fatalf("Failed to detect branch: %v", err)
	}

	// Should detect main or master
	if branchInfo.ShortName != "master" && branchInfo.ShortName != "main" {
		t.Errorf("Expected main or master, got %s", branchInfo.ShortName)
	}

	if branchInfo.Type != "main" {
		t.Errorf("Expected type main, got %s", branchInfo.Type)
	}
}

func TestParseBranchType(t *testing.T) {
	tests := []struct {
		name         string
		branchName   string
		expectedType string
	}{
		{"main branch", "main", "main"},
		{"master branch", "master", "main"},
		{"develop branch", "develop", "develop"},
		{"feature branch", "feature/user-auth", "feature"},
		{"hotfix branch", "hotfix/fix-bug", "hotfix"},
		{"bugfix branch", "bugfix/login-issue", "bugfix"},
		{"release branch", "release/v1.0.0", "release"},
		{"staging branch", "staging", "staging"},
		{"unknown branch", "random-branch", "unknown"},
	}

	detector := NewDetector(".")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &BranchInfo{
				ShortName: tt.branchName,
				Metadata:  make(map[string]string),
			}
			detector.parseBranchType(info)

			if info.Type != tt.expectedType {
				t.Errorf("Expected type %s, got %s", tt.expectedType, info.Type)
			}
		})
	}
}

func TestCheckProtected(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		expected   bool
	}{
		{"main is protected", "main", true},
		{"master is protected", "master", true},
		{"develop is protected", "develop", true},
		{"staging is protected", "staging", true},
		{"feature is not protected", "feature/test", false},
		{"hotfix is not protected", "hotfix/test", false},
	}

	detector := NewDetector(".")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &BranchInfo{
				ShortName: tt.branchName,
			}
			detector.checkProtected(info)

			if info.IsProtected != tt.expected {
				t.Errorf("Expected IsProtected=%v, got %v", tt.expected, info.IsProtected)
			}
		})
	}
}

func TestMetadataExtraction(t *testing.T) {
	detector := NewDetector(".")

	info := &BranchInfo{
		ShortName: "feature/JIRA-123-user-authentication",
		Metadata:  make(map[string]string),
	}
	detector.parseBranchType(info)

	// Should extract ticket number
	if ticket, exists := info.Metadata["ticket"]; !exists || ticket != "JIRA-123" {
		t.Errorf("Expected ticket JIRA-123, got %s", ticket)
	}

	// Should extract suffix
	if suffix, exists := info.Metadata["suffix"]; !exists || suffix != "JIRA-123-user-authentication" {
		t.Errorf("Expected suffix, got %s", suffix)
	}
}
