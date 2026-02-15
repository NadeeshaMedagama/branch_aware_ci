# 🚀 Quick Command Reference

Copy and paste these commands to complete your setup!

---

## 📦 STEP 1: Push to GitHub (2 minutes)

```bash
# Navigate to project
cd /Users/nadeesha_medagama/GolandProjects/branch-aware-ci

# Set new remote
git remote set-url origin https://github.com/NadeeshaMedagama/branch_aware_ci.git

# Stage all changes
git add .

# Commit
git commit -m "chore: migrate to NadeeshaMedagama/branch_aware_ci

- Updated all Go import paths
- Updated module path in go.mod
- Updated all documentation
- Updated GitHub Actions workflows
- Ready for production deployment"

# Push to GitHub
git push -u origin main
```

---

## 🏷️ STEP 2: Create First Release (30 seconds)

```bash
# Create and push v1.0.0 tag
git tag -a v1.0.0 -m "🎉 Initial Release

Features:
- Smart branch detection and type classification
- Environment mapping with priority-based rules
- Policy engine for deployment decisions
- Multiple output formats (JSON, YAML, GitHub Actions)
- Microservices architecture (Gateway, Detector, Policy Engine)
- Docker Compose support
- CodeQL security scanning
- Copilot AI code review
- Automated releases
- Comprehensive documentation

This is the first production-ready release of Branch Aware CI."

# Push tag to trigger release workflow
git push origin v1.0.0
```

---

## 🐳 STEP 3: Test Docker Locally (1 minute)

```bash
# Build all services
docker-compose build

# Start services
docker-compose up -d

# Check health
curl http://localhost:8080/health

# View logs
docker-compose logs -f gateway

# Stop services
docker-compose down
```

---

## 🧪 STEP 4: Test CLI Locally (30 seconds)

```bash
# Build
make build

# Run in current repo
./bin/branch-aware-ci

# Test with JSON output
./bin/branch-aware-ci -format json

# Test with YAML output
./bin/branch-aware-ci -format yaml

# Initialize config file
./bin/branch-aware-ci -init
```

---

## ✅ STEP 5: Verify Everything (1 minute)

```bash
# Check module path
head -1 go.mod

# Build
make build

# Run tests
make test

# Check for old references (should be empty)
grep -r "nadeesha_medagama" --include="*.go" --include="*.md" . 2>/dev/null | grep -v ".git"

# Verify new references
grep -r "NadeeshaMedagama/branch_aware_ci" README.md | head -3
```

---

## 📝 GitHub Secrets to Add

After creating the repository, add these secrets:

**Path**: `Settings → Secrets and variables → Actions → New repository secret`

```
Secret 1:
Name: DOCKER_USERNAME
Value: NadeeshaMedagama

Secret 2:
Name: DOCKER_PASSWORD
Value: [Your Docker Hub Access Token]
```

---

## 🎯 Test in Another Repository

Create `.github/workflows/test.yml`:

```yaml
name: Test Branch Aware CI

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Branch Detection
        id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Show Results
        run: |
          echo "Branch: ${{ steps.branch.outputs.branch_name }}"
          echo "Environment: ${{ steps.branch.outputs.environment }}"
          echo "Deploy: ${{ steps.branch.outputs.should_deploy }}"
```

---

## 🔧 Useful Development Commands

```bash
# Clean build cache
go clean -cache

# Update dependencies
go get -u ./...
go mod tidy

# Format code
go fmt ./...

# Run linter (if installed)
golangci-lint run

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o branch-aware-ci-linux

# Run with custom config
./bin/branch-aware-ci -config .branchci.yml

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📊 Monitor Your Release

After pushing the tag, monitor these:

```bash
# Check workflow status
# Visit: https://github.com/NadeeshaMedagama/branch_aware_ci/actions

# View release
# Visit: https://github.com/NadeeshaMedagama/branch_aware_ci/releases

# Check Docker Hub
# Visit: https://hub.docker.com/u/nadeeshamedagama
```

---

## 🌟 Promote Your Project

```bash
# Add topics to repository (in GitHub UI):
# - github-actions
# - cicd
# - devops
# - automation
# - branch-detection
# - deployment
# - golang
# - docker
# - microservices

# Star your own repository
# Visit: https://github.com/NadeeshaMedagama/branch_aware_ci

# Share on social media
# Tweet: "Just released Branch Aware CI 🚀 - Automatic branch detection and intelligent CI/CD decisions for GitHub Actions! Check it out: https://github.com/NadeeshaMedagama/branch_aware_ci"
```

---

## 📦 Docker Hub Repositories to Create

1. https://hub.docker.com/repository/create
   - Name: `branch-aware-gateway`
   
2. https://hub.docker.com/repository/create
   - Name: `branch-aware-branch-detector`
   
3. https://hub.docker.com/repository/create
   - Name: `branch-aware-policy-engine`

---

## 🎓 After Release

```bash
# Install via Go
go install github.com/NadeeshaMedagama/branch_aware_ci@latest

# Verify installation
branch-aware-ci -version

# Use in any repository
cd ~/projects/my-app
branch-aware-ci
```

---

## 🆘 Emergency Rollback

If something goes wrong:

```bash
# Delete tag locally
git tag -d v1.0.0

# Delete tag remotely
git push --delete origin v1.0.0

# Delete release in GitHub UI
# Visit: https://github.com/NadeeshaMedagama/branch_aware_ci/releases

# Fix issues, then re-tag
git tag -a v1.0.0 -m "Initial release (fixed)"
git push origin v1.0.0
```

---

**Everything ready! 🎉 Just follow the steps above!**

