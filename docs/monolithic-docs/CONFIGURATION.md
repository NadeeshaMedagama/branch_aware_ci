# Configuration Guide

This guide explains all configuration options for Branch-Aware CI.

## Table of Contents

- [Quick Start](#quick-start)
- [Configuration File](#configuration-file)
- [Environments](#environments)
- [Branch Mappings](#branch-mappings)
- [Policies](#policies)
- [Examples](#examples)

## Quick Start

Initialize a default configuration:

```bash
branch-aware-ci -init
```

This creates `.branchci.yml` with sensible defaults.

## Configuration File

Branch-Aware CI looks for configuration in these locations (in order):

1. Path specified via `-config` flag
2. `.branchci.yml` in repository root
3. `.branchci.yaml` in repository root
4. `.github/branchci.yml`
5. `.github/branchci.yaml`

If no configuration is found, default settings are used.

## Environments

Define deployment environments and their settings:

```yaml
environments:
  production:
    name: production              # Environment name
    requires_approval: true       # Require manual approval
    allowed_branches:             # Branches allowed to deploy
      - main
      - master
    variables:                    # Environment variables
      ENV: production
      LOG_LEVEL: info
      DATABASE_URL: prod.db.example.com
    notify_on_deploy: true        # Send notifications
```

### Environment Properties

| Property | Type | Required | Description |
|----------|------|----------|-------------|
| `name` | string | Yes | Environment identifier |
| `requires_approval` | boolean | No | Whether manual approval is needed |
| `allowed_branches` | array | No | Branch patterns allowed for this environment |
| `variables` | map | No | Environment-specific variables |
| `notify_on_deploy` | boolean | No | Whether to send deployment notifications |

## Branch Mappings

Map branch patterns to environments and actions:

```yaml
branch_mappings:
  - pattern: main                 # Branch pattern (exact or glob)
    environment: production       # Target environment
    actions:                      # Actions to perform
      - test
      - deploy
      - notify
    priority: 100                 # Priority (higher = checked first)
```

### Pattern Syntax

- **Exact match**: `main`, `staging`, `develop`
- **Wildcard**: `feature/*`, `release/*`, `hotfix/*`
- **Glob**: `feature-*`, `release-v*`

### Common Actions

- `test` - Run test suite
- `deploy` - Deploy to environment
- `notify` - Send notifications
- `build` - Build artifacts
- `scan` - Security scanning

### Priority

Higher priority mappings are evaluated first. This allows you to override general patterns with specific rules.

Example:
```yaml
branch_mappings:
  # High priority - specific branch
  - pattern: main
    environment: production
    priority: 100

  # Lower priority - pattern match
  - pattern: release/*
    environment: staging
    priority: 80

  # Lowest priority - catch-all
  - pattern: feature/*
    environment: development
    priority: 50
```

## Policies

Global rules that apply to all branches:

```yaml
policies:
  require_tests: true                    # Always run tests
  require_code_review: true              # Require reviews for protected branches
  blocked_branch_patterns:               # Patterns that cannot deploy
    - temp/*
    - wip/*
  auto_deploy_branches:                  # Auto-deploy these branches
    - main
    - staging
```

### Policy Properties

| Property | Type | Description |
|----------|------|-------------|
| `require_tests` | boolean | Add 'test' action to all decisions |
| `require_code_review` | boolean | Require approval for protected branches |
| `blocked_branch_patterns` | array | Branch patterns that cannot deploy |
| `auto_deploy_branches` | array | Branches that auto-deploy |

## Examples

### Example 1: Simple Configuration

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main]

  staging:
    name: staging
    allowed_branches: [develop]

branch_mappings:
  - pattern: main
    environment: production
    actions: [test, deploy]
    priority: 100

  - pattern: develop
    environment: staging
    actions: [test, deploy]
    priority: 90

  - pattern: feature/*
    environment: staging
    actions: [test]
    priority: 50

policies:
  require_tests: true
  auto_deploy_branches: [main]
```

### Example 2: Multi-Environment Setup

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main, master]
    variables:
      ENV: production
      API_URL: https://api.example.com
    notify_on_deploy: true

  staging:
    name: staging
    allowed_branches: [staging, develop, release/*]
    variables:
      ENV: staging
      API_URL: https://staging-api.example.com
    notify_on_deploy: true

  development:
    name: development
    allowed_branches: [feature/*, bugfix/*]
    variables:
      ENV: development
      API_URL: https://dev-api.example.com

branch_mappings:
  - pattern: main
    environment: production
    actions: [test, security-scan, deploy, notify]
    priority: 100

  - pattern: staging
    environment: staging
    actions: [test, deploy, notify]
    priority: 90

  - pattern: release/*
    environment: staging
    actions: [test, security-scan, deploy]
    priority: 85

  - pattern: feature/*
    environment: development
    actions: [test, deploy]
    priority: 50

policies:
  require_tests: true
  require_code_review: true
  blocked_branch_patterns: [temp/*, wip/*]
  auto_deploy_branches: [main, staging]
```

### Example 3: Hotfix-Focused Configuration

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main, hotfix/*]

  staging:
    name: staging
    allowed_branches: [staging, hotfix/*]

branch_mappings:
  # Hotfixes get special treatment
  - pattern: hotfix/*
    environment: staging
    actions: [test, security-scan, deploy]
    priority: 95

  # Production
  - pattern: main
    environment: production
    actions: [test, deploy, notify]
    priority: 100

  # Regular development
  - pattern: feature/*
    environment: development
    actions: [test]
    priority: 50

policies:
  require_tests: true
  require_code_review: true
  # Hotfixes can auto-deploy to staging
  auto_deploy_branches: [main, hotfix/*]
```

## Variables in GitHub Actions

Environment variables defined in configuration are exported:

```yaml
environments:
  production:
    variables:
      ENV: production
      DATABASE_URL: prod.example.com
```

Access in workflow:

```yaml
- name: Use variables
  run: |
    echo "Environment: $ENV"
    echo "Database: $DATABASE_URL"
```

## Override Defaults

You can override specific defaults while keeping others:

```yaml
# Only specify what you want to change
branch_mappings:
  # Override the main branch mapping
  - pattern: main
    environment: production
    actions: [test, security-scan, deploy, notify, create-release]
    priority: 100

# Other mappings use defaults
```

## Validation

Branch-Aware CI validates your configuration:

- Ensures required fields are present
- Checks for duplicate priorities
- Validates branch patterns
- Warns about unreachable mappings

## Best Practices

1. **Use high priorities for specific branches**
   ```yaml
   - pattern: main
     priority: 100
   ```

2. **Use lower priorities for wildcards**
   ```yaml
   - pattern: feature/*
     priority: 50
   ```

3. **Define clear environments**
   ```yaml
   environments:
     production:
       requires_approval: true
     development:
       requires_approval: false
   ```

4. **Use descriptive action names**
   ```yaml
   actions: [test, security-scan, deploy, notify]
   ```

5. **Block problematic patterns**
   ```yaml
   policies:
     blocked_branch_patterns:
       - temp/*
       - experimental/*
   ```

## Troubleshooting

### Branch not matching expected environment

Check priority values - higher priority mappings are evaluated first.

### Variables not appearing in workflow

Ensure you're using the correct output format:
```bash
-format github-env  # or github-output
```

### Environment requires approval but not showing

Check the environment configuration:
```yaml
environments:
  production:
    requires_approval: true
```

## Reference

See [`.branchci.example.yml`](../../.branchci.example.yml) for a complete working example.

