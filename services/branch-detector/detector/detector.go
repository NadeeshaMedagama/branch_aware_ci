package detector

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/NadeeshaMedagama/branch_aware_ci/pkg/interfaces"
	"github.com/go-git/go-git/v5"
)

// BranchDetector implements the IBranchDetector interface
// Following Single Responsibility Principle: only handles branch detection
type BranchDetector struct {
	patterns map[string]*regexp.Regexp
}

// NewBranchDetector creates a new instance of BranchDetector
func NewBranchDetector() *BranchDetector {
	return &BranchDetector{
		patterns: compileBranchPatterns(),
	}
}

// DetectBranch implements IBranchDetector.DetectBranch
func (d *BranchDetector) DetectBranch(ctx context.Context, repoPath string) (*interfaces.BranchInfo, error) {
	if repoPath == "" {
		repoPath = "."
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	branchName := head.Name().Short()
	info := &interfaces.BranchInfo{
		Name:      branchName,
		ShortName: branchName,
		Metadata:  make(map[string]string),
	}

	d.parseBranchType(info)
	d.checkProtected(info)

	return info, nil
}

// GetBranchInfo implements IBranchDetector.GetBranchInfo
func (d *BranchDetector) GetBranchInfo(ctx context.Context, repoPath string, branchName string) (*interfaces.BranchInfo, error) {
	// For now, return basic info based on branch name
	info := &interfaces.BranchInfo{
		Name:      branchName,
		ShortName: branchName,
		Metadata:  make(map[string]string),
	}

	d.parseBranchType(info)
	d.checkProtected(info)

	return info, nil
}

// compileBranchPatterns creates and compiles all branch patterns
// Following Open/Closed Principle: easy to add new patterns without modifying existing code
func compileBranchPatterns() map[string]*regexp.Regexp {
	return map[string]*regexp.Regexp{
		"feature": regexp.MustCompile(`^feature/(.+)$`),
		"hotfix":  regexp.MustCompile(`^hotfix/(.+)$`),
		"bugfix":  regexp.MustCompile(`^bugfix/(.+)$`),
		"release": regexp.MustCompile(`^release/(.+)$`),
		"develop": regexp.MustCompile(`^(develop|development)$`),
		"staging": regexp.MustCompile(`^staging$`),
		"main":    regexp.MustCompile(`^(main|master)$`),
	}
}

// parseBranchType determines the branch type and extracts metadata
func (d *BranchDetector) parseBranchType(info *interfaces.BranchInfo) {
	name := info.ShortName

	for branchType, pattern := range d.patterns {
		if matches := pattern.FindStringSubmatch(name); matches != nil {
			info.Type = branchType
			if len(matches) > 1 {
				info.Metadata["suffix"] = matches[1]

				// Extract ticket number (e.g., JIRA-123)
				ticketPattern := regexp.MustCompile(`([A-Z]+-\d+)`)
				if ticketMatches := ticketPattern.FindStringSubmatch(matches[1]); ticketMatches != nil {
					info.Metadata["ticket"] = ticketMatches[1]
				}
			}
			return
		}
	}

	info.Type = "unknown"
}

// checkProtected determines if the branch is a protected branch
func (d *BranchDetector) checkProtected(info *interfaces.BranchInfo) {
	protectedBranches := []string{"main", "master", "develop", "staging", "production"}
	for _, protected := range protectedBranches {
		if strings.EqualFold(info.ShortName, protected) {
			info.IsProtected = true
			return
		}
	}
	info.IsProtected = false
}

// GetRepositoryRoot returns the root path of the Git repository
func (d *BranchDetector) GetRepositoryRoot(repoPath string) (string, error) {
	if repoPath == "" {
		repoPath = "."
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	return filepath.Clean(worktree.Filesystem.Root()), nil
}
