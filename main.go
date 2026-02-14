package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nadeesha_medagama/branch-aware-ci/pkg/config"
	"github.com/nadeesha_medagama/branch-aware-ci/pkg/git"
	"github.com/nadeesha_medagama/branch-aware-ci/pkg/output"
	"github.com/nadeesha_medagama/branch-aware-ci/pkg/policy"
)

var (
	version = "1.0.0"
)

func main() {
	// Command-line flags
	configPath := flag.String("config", "", "Path to config file (default: .branchci.yml)")
	outputFormat := flag.String("format", "human", "Output format (json, yaml, env, github-env, github-output, human)")
	repoPath := flag.String("repo", ".", "Path to Git repository")
	initConfig := flag.Bool("init", false, "Initialize a default config file")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("branch-aware-ci v%s\n", version)
		os.Exit(0)
	}

	// Initialize config
	if *initConfig {
		if err := initializeConfig(*configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ Configuration file created successfully!")
		os.Exit(0)
	}

	// Run the main analysis
	if err := run(*repoPath, *configPath, *outputFormat); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(repoPath, configPath, outputFormat string) error {
	// Detect Git branch
	detector := git.NewDetector(repoPath)
	branchInfo, err := detector.DetectBranch()
	if err != nil {
		return fmt.Errorf("failed to detect branch: %w", err)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Evaluate policy and make decision
	engine := policy.NewEngine(cfg)
	decision, err := engine.Evaluate(branchInfo)
	if err != nil {
		return fmt.Errorf("failed to evaluate policy: %w", err)
	}

	// Format and output result
	formatter := output.NewFormatter(output.Format(outputFormat))
	result, err := formatter.Format(decision)
	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Println(result)
	return nil
}

func initializeConfig(configPath string) error {
	if configPath == "" {
		configPath = ".branchci.yml"
	}

	// Check if file already exists
	if _, err := os.Stat(configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", configPath)
	}

	// Create default config
	cfg := config.DefaultConfig()

	// Save to file
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return err
	}

	return nil
}
