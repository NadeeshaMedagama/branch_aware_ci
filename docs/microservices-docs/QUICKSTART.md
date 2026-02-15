# Quick Start Guide

Get started with Branch-Aware CI in 5 minutes!

## Prerequisites

- Git repository
- GitHub Actions (for GitHub integration)
- OR Go 1.23+ (for local usage)

## Option 1: GitHub Action (Recommended)

### Step 1: Add to your workflow

Create `.github/workflows/deploy.yml`:

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
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Detect branch and environment
        id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Show results
        run: |
          echo "Branch: ${{ steps.branch.outputs.branch_name }}"
          echo "Environment: ${{ steps.branch.outputs.environment }}"
          echo "Should Deploy: ${{ steps.branch.outputs.should_deploy }}"
      
      - name: Deploy
        if: steps.branch.outputs.should_deploy == 'true'
        run: |
          echo "Deploying to ${{ steps.branch.outputs.environment }}"
          # Add your deployment commands here
```

### Step 2: Push to any branch

```bash
git checkout -b feature/test-ci
git add .
git commit -m "Test branch-aware CI"
git push
```

### Step 3: Check the workflow

Go to Actions tab in GitHub and see the results!

## Option 2: CLI Tool

### Step 1: Install

```bash
# Using Go
go install github.com/NadeeshaMedagama/branch_aware_ci@latest

# Or build from source
git clone https://github.com/NadeeshaMedagama/branch_aware_ci.git
cd branch-aware-ci
make build
```

### Step 2: Initialize configuration

```bash
branch-aware-ci -init
```

This creates `.branchci.yml` with default settings.

### Step 3: Run

```bash
# Detect current branch
branch-aware-ci

# Output in JSON
branch-aware-ci -format json

# Use custom config
branch-aware-ci -config my-config.yml
```

## Option 3: Docker

### Step 1: Pull image

```bash
docker pull NadeeshaMedagama/branch_aware_ci:latest
```

### Step 2: Run

```bash
docker run -v $(pwd):/repo NadeeshaMedagama/branch_aware_ci:latest -repo /repo
```

## Customization

### Create custom configuration

`.branchci.yml`:

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main]
    variables:
      ENV: production

  staging:
    name: staging
    allowed_branches: [staging, develop]
    variables:
      ENV: staging

branch_mappings:
  - pattern: main
    environment: production
    actions: [test, deploy, notify]
    priority: 100

  - pattern: staging
    environment: staging
    actions: [test, deploy]
    priority: 90

  - pattern: feature/*
    environment: development
    actions: [test]
    priority: 50

policies:
  require_tests: true
  auto_deploy_branches: [main, staging]
```

## Common Workflows

### Deploy only on main

```yaml
- uses: NadeeshaMedagama/branch_aware_ci@v1
  id: branch

- name: Deploy
  if: steps.branch.outputs.environment == 'production'
  run: ./deploy.sh
```

### Run tests on all branches

```yaml
- uses: NadeeshaMedagama/branch_aware_ci@v1
  id: branch

- name: Test
  if: contains(steps.branch.outputs.actions, 'test')
  run: npm test
```

### Different commands per environment

```yaml
- uses: NadeeshaMedagama/branch_aware_ci@v1
  id: branch

- name: Deploy
  run: |
    case "${{ steps.branch.outputs.environment }}" in
      production)
        ./deploy-prod.sh
        ;;
      staging)
        ./deploy-staging.sh
        ;;
      development)
        ./deploy-dev.sh
        ;;
    esac
```

## Next Steps

- Read the [Configuration Guide](../monolithic-docs/CONFIGURATION.md)
- Check out [Use Cases](../monolithic-docs/USE-CASES.md)
- See [Example Workflows](.github/workflows/)
- Contribute on [GitHub](https://github.com/NadeeshaMedagama/branch_aware_ci)

## Troubleshooting

### "No such file or directory" error

Make sure you're running from a Git repository.

### Branch not detected correctly

Check your `.branchci.yml` configuration and pattern matching.

### Variables not showing up

Use `-format github-env` or `-format github-output` in GitHub Actions.

## Support

- 📖 [Full Documentation](../../README.md)
- 💬 [GitHub Discussions](https://github.com/NadeeshaMedagama/branch_aware_ci/discussions)
- 🐛 [Report Issues](https://github.com/NadeeshaMedagama/branch_aware_ci/issues)

