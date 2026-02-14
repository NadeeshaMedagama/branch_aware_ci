package git

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

// BranchInfo contains information about the current Git branch
type BranchInfo struct {
	Name        string            // Full branch name (e.g., "feature/user-auth")
	ShortName   string            // Branch name without refs/heads/ prefix
	Type        string            // Branch type (e.g., "feature", "hotfix", "release", "main")
	Metadata    map[string]string // Extracted metadata from branch name
	IsProtected bool              // Whether this is a protected branch
}

// Detector handles Git branch detection
type Detector struct {
	repoPath string
}

// NewDetector creates a new Git detector
func NewDetector(repoPath string) *Detector {
	if repoPath == "" {
		repoPath = "."
	}
	return &Detector{repoPath: repoPath}
}

// DetectBranch detects the current Git branch
func (d *Detector) DetectBranch() (*BranchInfo, error) {
	repo, err := git.PlainOpen(d.repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	branchName := head.Name().Short()
	info := &BranchInfo{
		Name:      branchName,
		ShortName: branchName,
		Metadata:  make(map[string]string),
	}

	// Determine branch type and extract metadata
	d.parseBranchType(info)
	d.checkProtected(info)

	return info, nil
}

// parseBranchType determines the branch type and extracts metadata
func (d *Detector) parseBranchType(info *BranchInfo) {
	name := info.ShortName

	// Common branch patterns
	patterns := map[string]*regexp.Regexp{
		"feature": regexp.MustCompile(`^feature/(.+)$`),
		"hotfix":  regexp.MustCompile(`^hotfix/(.+)$`),
		"bugfix":  regexp.MustCompile(`^bugfix/(.+)$`),
		"release": regexp.MustCompile(`^release/(.+)$`),
		"develop": regexp.MustCompile(`^(develop|development)$`),
		"staging": regexp.MustCompile(`^staging$`),
		"main":    regexp.MustCompile(`^(main|master)$`),
	}

	for branchType, pattern := range patterns {
		if matches := pattern.FindStringSubmatch(name); matches != nil {
			info.Type = branchType
			if len(matches) > 1 {
				// Extract the suffix (e.g., "user-auth" from "feature/user-auth")
				info.Metadata["suffix"] = matches[1]

				// Try to extract ticket number (e.g., JIRA-123)
				ticketPattern := regexp.MustCompile(`([A-Z]+-\d+)`)
				if ticketMatches := ticketPattern.FindStringSubmatch(matches[1]); ticketMatches != nil {
					info.Metadata["ticket"] = ticketMatches[1]
				}
			}
			return
		}
	}

	// Default type for unknown patterns
	info.Type = "unknown"
}

// checkProtected determines if the branch is a protected branch
func (d *Detector) checkProtected(info *BranchInfo) {
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
func (d *Detector) GetRepositoryRoot() (string, error) {
	repo, err := git.PlainOpen(d.repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	return filepath.Clean(worktree.Filesystem.Root()), nil
}
