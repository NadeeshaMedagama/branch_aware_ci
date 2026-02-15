# 🚀 Branch-Aware CI - Quick Reference Card

## 📦 Installation

```bash
# CLI Tool
go install github.com/NadeeshaMedagama/branch_aware_ci@latest

# From Source
git clone https://github.com/NadeeshaMedagama/branch_aware_ci.git
cd branch-aware-ci
make build
```

## 🎯 CLI Usage

```bash
# Basic usage
branch-aware-ci

# Show version
branch-aware-ci -version

# Different output formats
branch-aware-ci -format json
branch-aware-ci -format yaml
branch-aware-ci -format env

# Custom config
branch-aware-ci -config .branchci.yml

# Different repo
branch-aware-ci -repo /path/to/repo

# Initialize config
branch-aware-ci -init
```

## 🔄 GitHub Actions Usage

```yaml
# Minimal
- uses: NadeeshaMedagama/branch_aware_ci@v1
  id: branch

# With options
- uses: NadeeshaMedagama/branch_aware_ci@v1
  id: branch
  with:
    config-path: '.branchci.yml'
    output-format: 'github-output'
```

## 📤 Outputs

```yaml
outputs:
  branch_name:       # e.g., "feature/user-auth"
  branch_type:       # e.g., "feature"
  environment:       # e.g., "development"
  should_deploy:     # e.g., "true"
  requires_approval: # e.g., "false"
  actions:          # e.g., "test,deploy"
```

## 🔧 Using Outputs

```yaml
# Conditional deployment
- name: Deploy
  if: steps.branch.outputs.should_deploy == 'true'
  run: ./deploy.sh ${{ steps.branch.outputs.environment }}

# Environment-specific jobs
- name: Production Deploy
  if: steps.branch.outputs.environment == 'production'
  run: ./deploy-prod.sh

# Action-based execution
- name: Run Tests
  if: contains(steps.branch.outputs.actions, 'test')
  run: npm test
```

## ⚙️ Configuration (.branchci.yml)

```yaml
# Minimal config
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main]

branch_mappings:
  - pattern: main
    environment: production
    actions: [test, deploy]
    priority: 100

  - pattern: feature/*
    environment: development
    actions: [test]
    priority: 50

policies:
  require_tests: true
```

## 🌿 Default Branch Mappings

| Branch Pattern | Environment | Deploy | Approval |
|----------------|-------------|--------|----------|
| main, master   | production  | ✅     | ✅       |
| staging        | staging     | ✅     | ❌       |
| develop        | staging     | ✅     | ❌       |
| release/*      | staging     | ✅     | ❌       |
| feature/*      | development | ❌     | ❌       |
| bugfix/*       | development | ❌     | ❌       |
| hotfix/*       | staging     | ✅     | ❌       |

## 🎨 Common Patterns

### Auto-Deploy

```yaml
- uses: NadeeshaMedagama/branch_aware_ci@v1
  id: branch

- if: steps.branch.outputs.should_deploy == 'true'
  run: ./deploy.sh ${{ steps.branch.outputs.environment }}
```

### Environment Variables

```yaml
- uses: NadeeshaMedagama/branch_aware_ci@v1
  with:
    output-format: github-env

- run: echo "Deploying to $ENVIRONMENT"
```

### Multi-Environment Pipeline

```yaml
deploy-dev:
  if: needs.detect.outputs.environment == 'development'
  # dev deployment

deploy-staging:
  if: needs.detect.outputs.environment == 'staging'
  # staging deployment

deploy-prod:
  if: needs.detect.outputs.environment == 'production'
  # production deployment
```

## 🐳 Docker

```bash
# Build
docker build -t branch-aware-ci .

# Run
docker run -v $(pwd):/repo branch-aware-ci -repo /repo

# With format
docker run -v $(pwd):/repo branch-aware-ci -repo /repo -format json
```

## 🔧 Makefile Targets

```bash
make build         # Build binary
make test          # Run tests
make clean         # Clean artifacts
make install       # Install globally
make run           # Build and run
make docker-build  # Build Docker image
make docker-run    # Run in Docker
make coverage      # Test coverage
```

## 📊 Output Formats

| Format | Use Case | Example |
|--------|----------|---------|
| `json` | Scripts, APIs | `{"branch_name": "main"}` |
| `yaml` | Config files | `branch_name: main` |
| `env` | Shell scripts | `BRANCH_NAME=main` |
| `github-env` | GitHub Actions env | Writes to $GITHUB_ENV |
| `github-output` | GitHub Actions output | Sets step outputs |
| `human` | Debugging | Pretty-printed text |

## 🎯 Branch Types

| Type | Pattern | Example |
|------|---------|---------|
| main | `main`, `master` | `main` |
| staging | `staging` | `staging` |
| develop | `develop`, `development` | `develop` |
| feature | `feature/*` | `feature/user-auth` |
| bugfix | `bugfix/*` | `bugfix/login-fix` |
| hotfix | `hotfix/*` | `hotfix/security-patch` |
| release | `release/*` | `release/v1.0.0` |

## 🛡️ Best Practices

1. **Use high priority for specific branches**
   ```yaml
   - pattern: main
     priority: 100
   ```

2. **Protected branches require approval**
   ```yaml
   production:
     requires_approval: true
   ```

3. **Test all branches**
   ```yaml
   policies:
     require_tests: true
   ```

4. **Block temporary branches**
   ```yaml
   policies:
     blocked_branch_patterns: [temp/*, wip/*]
   ```

## 🔍 Troubleshooting

| Issue | Solution |
|-------|----------|
| Branch not detected | Ensure you're in a Git repository |
| Wrong environment | Check priority in branch_mappings |
| Variables missing | Use `-format github-env` or `github-output` |
| Config not loading | Check file path and YAML syntax |

## 📚 Documentation

- **README.md** - Main documentation
- **QUICKSTART.md** - 5-minute guide
- **docs/CONFIGURATION.md** - Config reference
- **docs/USE-CASES.md** - Examples
- **docs/ARCHITECTURE.md** - Technical details

## 🌐 Resources

- **GitHub**: https://github.com/NadeeshaMedagama/branch_aware_ci
- **Docker Hub**: NadeeshaMedagama/branch_aware_ci
- **Go Module**: github.com/NadeeshaMedagama/branch_aware_ci

## 💡 Quick Examples

### Feature Branch Preview
```yaml
- if: steps.branch.outputs.branch_type == 'feature'
  run: |
    PREVIEW_URL="https://${{ steps.branch.outputs.branch_name }}.preview.com"
    echo "Preview: $PREVIEW_URL"
```

### Hotfix Fast Track
```yaml
- if: steps.branch.outputs.branch_type == 'hotfix'
  run: |
    # Skip some checks for urgent fixes
    ./quick-deploy.sh
```

### Production Gate
```yaml
- if: steps.branch.outputs.environment == 'production'
  run: |
    # Extra security scans
    npm audit --production
    trivy scan .
```

---

**Need help?** Check the full documentation or open an issue!

**Quick Start:** `branch-aware-ci -init` → `branch-aware-ci`

🚀 **Ready to automate your CI/CD!**

