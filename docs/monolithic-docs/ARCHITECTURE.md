# Project Architecture

## Overview

Branch-Aware CI is a modular system designed to automatically detect Git branches and make intelligent CI/CD decisions.

```
┌──────────────────────────────────────────────────────┐
│                    User Input                         │
│  (Git Repo, Config File, Output Format)              │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│                  Git Detector                         │
│  • Detect current branch                             │
│  • Parse branch type (feature, hotfix, etc.)         │
│  • Extract metadata (tickets, names)                 │
│  • Identify protected branches                       │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│              Configuration Loader                     │
│  • Load .branchci.yml                                │
│  • Parse YAML configuration                          │
│  • Apply defaults if not found                       │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│                 Policy Engine                         │
│  • Match branch patterns                             │
│  • Select environment                                │
│  • Apply policies                                    │
│  • Make deployment decision                          │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│               Output Formatter                        │
│  • JSON / YAML                                       │
│  • Environment variables                             │
│  • GitHub Actions outputs                            │
│  • Human-readable                                    │
└─────────────────────┬────────────────────────────────┘
                      │
                      ▼
┌──────────────────────────────────────────────────────┐
│                    Output                             │
│  (Deployment Decision + Variables)                   │
└──────────────────────────────────────────────────────┘
```

## Components

### 1. Git Detector (`pkg/git/`)

**Responsibility**: Interact with Git repository and extract branch information

**Key Functions**:
- `DetectBranch()` - Get current branch
- `parseBranchType()` - Determine branch type
- `checkProtected()` - Identify protected branches
- Metadata extraction (ticket numbers, feature names)

**Dependencies**:
- `go-git/go-git` - Git operations

**Data Structures**:
```go
type BranchInfo struct {
    Name        string            // Full branch name
    ShortName   string            // Branch without refs/heads/
    Type        string            // feature, hotfix, main, etc.
    Metadata    map[string]string // Extracted data
    IsProtected bool              // Protected status
}
```

### 2. Configuration Loader (`pkg/config/`)

**Responsibility**: Load and parse configuration files

**Key Functions**:
- `LoadConfig()` - Read YAML configuration
- `DefaultConfig()` - Provide sensible defaults
- `SaveConfig()` - Write configuration
- `findConfigFile()` - Search common locations

**Dependencies**:
- `gopkg.in/yaml.v3` - YAML parsing

**Data Structures**:
```go
type Config struct {
    Environments   map[string]EnvironmentConfig
    BranchMappings []BranchMapping
    Policies       PolicyConfig
}

type EnvironmentConfig struct {
    Name             string
    RequiresApproval bool
    AllowedBranches  []string
    Variables        map[string]string
    NotifyOnDeploy   bool
}

type BranchMapping struct {
    Pattern     string
    Environment string
    Actions     []string
    Priority    int
}
```

### 3. Policy Engine (`pkg/policy/`)

**Responsibility**: Evaluate rules and make decisions

**Key Functions**:
- `Evaluate()` - Main decision-making function
- `findBestMapping()` - Priority-based pattern matching
- `matchesPattern()` - Branch pattern matching
- `shouldDeploy()` - Deployment decision
- `applyPolicies()` - Apply global rules

**Decision Logic**:
1. Find matching branch pattern (highest priority)
2. Determine target environment
3. Check allowed branches
4. Apply policies (tests, approvals, etc.)
5. Generate warnings if needed

**Data Structures**:
```go
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
```

### 4. Output Formatter (`pkg/output/`)

**Responsibility**: Format decisions for different consumers

**Supported Formats**:
- **JSON** - Structured data for scripts
- **YAML** - Human-readable structured data
- **Environment Variables** - For shell scripts
- **GitHub Actions** - `$GITHUB_ENV` and `$GITHUB_OUTPUT`
- **Human-readable** - Debug output

**Key Functions**:
- `Format()` - Main formatting dispatcher
- `formatJSON()` - JSON output
- `formatGitHubOutput()` - GitHub Actions integration
- `formatHuman()` - Pretty-printed output

## Data Flow

### Example: Feature Branch

```
Input: feature/user-auth

1. Git Detector
   └─> BranchInfo {
         Name: "feature/user-auth"
         Type: "feature"
         Metadata: {suffix: "user-auth"}
       }

2. Configuration Loader
   └─> Loads .branchci.yml or defaults

3. Policy Engine
   └─> Matches pattern "feature/*"
   └─> Environment: "development"
   └─> Actions: ["test"]
   └─> ShouldDeploy: false

4. Output Formatter
   └─> JSON:
       {
         "branch_name": "feature/user-auth",
         "environment": "development",
         "should_deploy": false
       }
```

### Example: Main Branch

```
Input: main

1. Git Detector
   └─> BranchInfo {
         Name: "main"
         Type: "main"
         IsProtected: true
       }

2. Configuration Loader
   └─> Loads production environment config

3. Policy Engine
   └─> Matches pattern "main" (priority: 100)
   └─> Environment: "production"
   └─> Actions: ["test", "deploy", "notify"]
   └─> ShouldDeploy: true
   └─> RequiresApproval: true

4. Output Formatter
   └─> GitHub Output:
       branch_name=main
       environment=production
       should_deploy=true
```

## Directory Structure

```
branch-aware-ci/
├── main.go                 # CLI entry point
├── go.mod                  # Go dependencies
├── Dockerfile             # Container image
├── action.yml             # GitHub Action definition
├── Makefile              # Build automation
│
├── pkg/                   # Core packages
│   ├── git/              # Git operations
│   │   ├── detector.go
│   │   └── detector_test.go
│   ├── config/           # Configuration
│   │   └── config.go
│   ├── policy/           # Decision engine
│   │   ├── engine.go
│   │   └── engine_test.go
│   └── output/           # Formatters
│       └── formatter.go
│
├── .github/              # GitHub integration
│   └── workflows/        # Example workflows
│       ├── branch-aware-ci.yml
│       ├── simple-example.yml
│       └── advanced-example.yml
│
├── docs/                 # Documentation
│   ├── CONFIGURATION.md
│   └── USE-CASES.md
│
└── bin/                  # Build output (gitignored)
    └── branch-aware-ci
```

## Extension Points

### Adding New Branch Types

Edit `pkg/git/detector.go`:

```go
patterns := map[string]*regexp.Regexp{
    "feature":  regexp.MustCompile(`^feature/(.+)$`),
    "new-type": regexp.MustCompile(`^new-type/(.+)$`),
}
```

### Adding New Output Formats

Edit `pkg/output/formatter.go`:

```go
const (
    FormatJSON   Format = "json"
    FormatCustom Format = "custom"
)

func (f *Formatter) formatCustom(decision *policy.Decision) (string, error) {
    // Implementation
}
```

### Adding New Policies

Edit `pkg/policy/engine.go`:

```go
func (e *Engine) applyPolicies(decision *Decision, branchInfo *git.BranchInfo) {
    // Add new policy logic
}
```

## Design Principles

1. **Modularity**: Each package has a single responsibility
2. **Testability**: Core logic is decoupled from I/O
3. **Extensibility**: Easy to add new formats, policies, etc.
4. **Default-friendly**: Works without configuration
5. **Type-safe**: Go's type system prevents errors
6. **Idiomatic Go**: Follows Go best practices

## Performance Considerations

- **Git operations**: Cached repository handles
- **Configuration**: Loaded once per execution
- **Pattern matching**: Priority-sorted for early exit
- **Memory**: Minimal allocations, reuse structures

## Security

- **No credential storage**: Relies on existing Git auth
- **Read-only operations**: Never modifies repository
- **Input validation**: All patterns validated
- **Container isolation**: Docker runs with minimal privileges

## Future Enhancements

1. **Caching**: Cache decisions for performance
2. **Remote config**: Load config from URL
3. **Plugins**: External policy engines
4. **Web UI**: Visual configuration editor
5. **Metrics**: Prometheus integration

