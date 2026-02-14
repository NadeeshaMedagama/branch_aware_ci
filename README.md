# 🌿 Branch-Aware CI/CD

> Automatically detect Git branches and make intelligent CI/CD decisions for your GitHub Actions workflows

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub Actions](https://img.shields.io/badge/GitHub-Actions-2088FF?style=flat&logo=github-actions)](https://github.com/features/actions)

## 🚀 The Problem This Solves

Ever forget to update branch names in your CI/CD workflows? Tired of maintaining multiple workflow files for different environments? **Branch-Aware CI** automatically:

✅ **Detects your Git branch** without manual configuration  
✅ **Maps branches to environments** (main → production, feature/* → development)  
✅ **Makes deployment decisions** based on customizable policies  
✅ **Outputs structured data** for use in subsequent workflow steps  
✅ **Prevents mistakes** by enforcing branch-based rules  

---

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
  - [As a GitHub Action](#as-a-github-action)
  - [As a CLI Tool](#as-a-cli-tool)
- [Configuration](#-configuration)
- [Examples](#-examples)
- [Output Formats](#-output-formats)
- [Architecture](#-architecture)
- [Contributing](#-contributing)

---

## ⚡ Quick Start

### 1. Add to your GitHub Actions workflow:

```yaml
name: Deploy

on:
  push:
    branches:
      - '**'

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Detect branch and environment
        id: branch
        uses: nadeesha_medagama/branch-aware-ci@v1
      
      - name: Deploy
        if: steps.branch.outputs.should_deploy == 'true'
        run: |
          echo "Deploying to ${{ steps.branch.outputs.environment }}"
          # Your deployment commands here
```

### 2. That's it! 🎉

The action automatically:
- Detects `main` → deploys to **production**
- Detects `staging` → deploys to **staging**
- Detects `feature/*` → runs tests only (no deploy)

---

## 🎯 Features

### 🔍 Smart Branch Detection
- Automatically identifies branch type (feature, hotfix, release, main, etc.)
- Extracts metadata (ticket numbers, feature names)
- Recognizes protected branches

### 🌍 Environment Mapping
- Pre-configured mappings for common workflows
- Customizable via `.branchci.yml`
- Priority-based pattern matching

### 🛡️ Policy Enforcement
- Require approvals for production
- Enforce testing requirements
- Block specific branch patterns
- Auto-deploy rules

### 📤 Multiple Output Formats
- **JSON/YAML** - for scripts and tools
- **GitHub Actions outputs** - for workflow steps
- **Environment variables** - for shell scripts
- **Human-readable** - for debugging

### 🔧 Flexible Configuration
- YAML-based configuration
- Override defaults per project
- Template configurations included

---

## 💻 Installation

### As a GitHub Action (Recommended)

Add to your workflow file (`.github/workflows/*.yml`):

```yaml
- uses: nadeesha_medagama/branch-aware-ci@v1
  with:
    config-path: '.branchci.yml'  # optional
    output-format: 'github-output' # optional
```

### As a CLI Tool

#### Using Go:
```bash
go install github.com/nadeesha_medagama/branch-aware-ci@latest
```

#### From source:
```bash
git clone https://github.com/nadeesha_medagama/branch-aware-ci.git
cd branch-aware-ci
go build -o branch-aware-ci
```

#### Using Docker:
```bash
docker build -t branch-aware-ci .
docker run -v $(pwd):/repo branch-aware-ci -repo /repo
```

---

## 🔨 Usage

### As a GitHub Action

```yaml
jobs:
  detect:
    runs-on: ubuntu-latest
    outputs:
      environment: ${{ steps.detect.outputs.environment }}
      should_deploy: ${{ steps.detect.outputs.should_deploy }}
    
    steps:
      - uses: actions/checkout@v4
      - name: Detect branch
        id: detect
        uses: nadeesha_medagama/branch-aware-ci@v1
      
      - name: Show results
        run: |
          echo "Branch: ${{ steps.detect.outputs.branch_name }}"
          echo "Environment: ${{ steps.detect.outputs.environment }}"
          echo "Deploy: ${{ steps.detect.outputs.should_deploy }}"
```

### As a CLI Tool

```bash
# Basic usage (auto-detects current branch)
branch-aware-ci

# With custom config
branch-aware-ci -config .branchci.yml

# Different output formats
branch-aware-ci -format json
branch-aware-ci -format yaml
branch-aware-ci -format env

# For different repository
branch-aware-ci -repo /path/to/repo

# Initialize default config
branch-aware-ci -init

# Show version
branch-aware-ci -version
```

---

## ⚙️ Configuration

Create `.branchci.yml` in your repository root:

```yaml
# Environment definitions
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches:
      - main
      - master
    variables:
      ENV: production
    notify_on_deploy: true

  staging:
    name: staging
    requires_approval: false
    allowed_branches:
      - staging
      - develop
    variables:
      ENV: staging
    notify_on_deploy: true

  development:
    name: development
    requires_approval: false
    allowed_branches:
      - feature/*
      - bugfix/*
    variables:
      ENV: development
    notify_on_deploy: false

# Branch pattern mappings
branch_mappings:
  - pattern: main
    environment: production
    actions: [test, deploy, notify]
    priority: 100

  - pattern: feature/*
    environment: development
    actions: [test]
    priority: 50

# Global policies
policies:
  require_tests: true
  require_code_review: true
  auto_deploy_branches:
    - main
    - staging
```

See [`.branchci.example.yml`](.branchci.example.yml) for a complete example.

---

## 📚 Examples

### Example 1: Simple Auto-Deploy

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - id: branch
        uses: nadeesha_medagama/branch-aware-ci@v1
      
      - name: Deploy
        if: steps.branch.outputs.should_deploy == 'true'
        run: |
          ./deploy.sh ${{ steps.branch.outputs.environment }}
```

### Example 2: Environment-Specific Jobs

```yaml
jobs:
  detect:
    runs-on: ubuntu-latest
    outputs:
      environment: ${{ steps.detect.outputs.environment }}
    steps:
      - uses: actions/checkout@v4
      - id: detect
        uses: nadeesha_medagama/branch-aware-ci@v1

  deploy-prod:
    needs: detect
    if: needs.detect.outputs.environment == 'production'
    runs-on: ubuntu-latest
    steps:
      - run: echo "Deploying to production"

  deploy-staging:
    needs: detect
    if: needs.detect.outputs.environment == 'staging'
    runs-on: ubuntu-latest
    steps:
      - run: echo "Deploying to staging"
```

### Example 3: Using Environment Variables

```yaml
steps:
  - uses: nadeesha_medagama/branch-aware-ci@v1
    with:
      output-format: github-env
  
  - name: Use variables
    run: |
      echo "Environment: $ENVIRONMENT"
      echo "Branch Type: $BRANCH_TYPE"
      echo "Should Deploy: $SHOULD_DEPLOY"
```

See [`.github/workflows/`](.github/workflows/) for more complete examples.

---

## 📤 Output Formats

### GitHub Actions Outputs

```yaml
outputs:
  branch_name: "feature/user-auth"
  branch_type: "feature"
  environment: "development"
  should_deploy: "false"
  requires_approval: "false"
  actions: "test"
```

### JSON Output

```json
{
  "branch_name": "main",
  "branch_type": "main",
  "environment": "production",
  "should_deploy": true,
  "requires_approval": true,
  "actions": ["test", "deploy", "notify"],
  "variables": {
    "ENV": "production"
  }
}
```

### Environment Variables

```bash
BRANCH_NAME=main
BRANCH_TYPE=main
ENVIRONMENT=production
SHOULD_DEPLOY=true
REQUIRES_APPROVAL=true
ACTIONS=test,deploy,notify
ENV=production
```

---

## 🏗️ Architecture

```
┌─────────────────┐
│   Git Branch    │
│  Detection      │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Configuration  │
│  Loading        │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Policy Engine  │
│  Evaluation     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Decision       │
│  Formatting     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Output         │
└─────────────────┘
```

### Components

- **`pkg/git/`** - Git repository interaction and branch detection
- **`pkg/config/`** - Configuration loading and parsing
- **`pkg/policy/`** - Decision engine and rule evaluation
- **`pkg/output/`** - Output formatting (JSON, YAML, env, etc.)

---

## 🎨 Use Cases

### ✅ Automatic Environment Selection
- `main` branch → production deployment
- `staging` branch → staging deployment
- `feature/*` → development/preview environment

### ✅ Conditional CI/CD Steps
- Run security scans only for production deployments
- Skip deployments for feature branches
- Require approvals for protected branches

### ✅ Multi-Environment Pipelines
- Single workflow handles all environments
- Different rules per branch type
- Customizable per project

### ✅ Prevent Configuration Mistakes
- No hardcoded branch names
- Automatic branch detection
- Policy enforcement

---

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](docs/monolithic-docs/CONTRIBUTING.md) for guidelines.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/nadeesha_medagama/branch-aware-ci.git
cd branch-aware-ci

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o branch-aware-ci

# Run locally
./branch-aware-ci
```

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🌟 Show Your Support

If this project helped you, please give it a ⭐️!

---

## 📞 Contact

- **Author**: Nadeesha Medagama
- **GitHub**: [@nadeesha_medagama](https://github.com/nadeesha_medagama)

---

## 🗺️ Roadmap

- [ ] Support for more version control platforms (GitLab, Bitbucket)
- [ ] Web UI for configuration visualization
- [ ] Integration with popular deployment tools
- [ ] Monorepo path filtering
- [ ] Parallel environment deployments
- [ ] Rollback strategies

---

**Built with ❤️ for the DevOps community**

