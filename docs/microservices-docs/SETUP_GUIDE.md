# Quick Setup Guide - Branch Aware CI

## 🎯 Repository Information
- **GitHub Username**: NadeeshaMedagama
- **Repository Name**: branch_aware_ci
- **Repository URL**: https://github.com/NadeeshaMedagama/branch_aware_ci

## ✅ Changes Completed

All references have been updated from:
- ❌ `nadeesha_medagama/branch-aware-ci`
- ✅ `NadeeshaMedagama/branch_aware_ci`

**Status**: ✅ Build successful, tests passing, ready to deploy!

---

## 🚀 Quick Start - Set Up Your GitHub Repository

### Step 1: Create GitHub Repository
1. Go to https://github.com/new
2. **Repository name**: `branch_aware_ci`
3. **Description**: Branch-aware CI/CD automation for GitHub Actions
4. **Visibility**: Public (recommended for GitHub Actions)
5. Click "Create repository"

### Step 2: Update Git Remote and Push
```bash
cd /Users/nadeesha_medagama/GolandProjects/branch-aware-ci

# Update remote URL
git remote set-url origin https://github.com/NadeeshaMedagama/branch_aware_ci.git

# Add all changes
git add .

# Commit the updates
git commit -m "chore: update repository references to NadeeshaMedagama/branch_aware_ci"

# Push to GitHub
git push -u origin main
```

---

## 🐳 Set Up Docker Hub

### Create Docker Hub Repositories
Go to https://hub.docker.com and create these repositories:
1. `NadeeshaMedagama/branch-aware-gateway`
2. `NadeeshaMedagama/branch-aware-branch-detector`
3. `NadeeshaMedagama/branch-aware-policy-engine`

### Get Docker Hub Access Token
1. Go to https://hub.docker.com/settings/security
2. Click "New Access Token"
3. Name: `branch_aware_ci_github_actions`
4. Permissions: Read, Write, Delete
5. Click "Generate"
6. **Copy the token** (you won't see it again!)

---

## 🔐 Configure GitHub Secrets

In your GitHub repository (https://github.com/NadeeshaMedagama/branch_aware_ci):

1. Go to **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**
3. Add these secrets:

| Name | Value |
|------|-------|
| `DOCKER_USERNAME` | `NadeeshaMedagama` |
| `DOCKER_PASSWORD` | Your Docker Hub access token from above |

---

## 🏷️ Create Your First Release

```bash
# Tag the current commit
git tag -a v1.0.0 -m "Initial release"

# Push the tag
git push origin v1.0.0
```

This will automatically trigger:
- ✅ Build binaries for multiple platforms
- ✅ Create GitHub Release
- ✅ Build and push Docker images
- ✅ Update documentation

---

## 🧪 Test the GitHub Action

### In Another Repository
Create a `.github/workflows/test-branch-aware.yml`:

```yaml
name: Test Branch Aware CI

on:
  push:
    branches:
      - '**'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Detect branch and environment
        id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Show results
        run: |
          echo "Branch: ${{ steps.branch.outputs.branch_name }}"
          echo "Type: ${{ steps.branch.outputs.branch_type }}"
          echo "Environment: ${{ steps.branch.outputs.environment }}"
          echo "Should Deploy: ${{ steps.branch.outputs.should_deploy }}"
          echo "Actions: ${{ steps.branch.outputs.actions }}"
```

---

## 📦 Install as CLI Tool

### Using Go
```bash
go install github.com/NadeeshaMedagama/branch_aware_ci@latest
```

### Using Homebrew (Future)
```bash
# After creating a Homebrew tap
brew install NadeeshaMedagama/tap/branch-aware-ci
```

### Download Binary
```bash
# Linux (amd64)
wget https://github.com/NadeeshaMedagama/branch_aware_ci/releases/latest/download/branch-aware-ci-linux-amd64
chmod +x branch-aware-ci-linux-amd64
sudo mv branch-aware-ci-linux-amd64 /usr/local/bin/branch-aware-ci

# macOS (Apple Silicon)
wget https://github.com/NadeeshaMedagama/branch_aware_ci/releases/latest/download/branch-aware-ci-darwin-arm64
chmod +x branch-aware-ci-darwin-arm64
sudo mv branch-aware-ci-darwin-arm64 /usr/local/bin/branch-aware-ci
```

---

## 🐳 Run with Docker

### Pull and Run
```bash
# Pull the gateway
docker pull NadeeshaMedagama/branch-aware-gateway:latest

# Pull branch detector
docker pull NadeeshaMedagama/branch-aware-branch-detector:latest

# Pull policy engine
docker pull NadeeshaMedagama/branch-aware-policy-engine:latest

# Run with docker-compose
docker-compose up -d
```

---

## 🔍 Verify Everything Works

### 1. Check Build
```bash
cd /Users/nadeesha_medagama/GolandProjects/branch-aware-ci
make build
```

### 2. Run Tests
```bash
make test
```

### 3. Run Locally
```bash
./bin/branch-aware-ci
```

### 4. Test with Docker
```bash
docker-compose up -d
curl http://localhost:8080/health
```

---

## 📚 Available Workflows

Your repository includes these GitHub Actions workflows:

| Workflow | File | Trigger | Purpose |
|----------|------|---------|---------|
| CodeQL | `codeql.yml` | Push, PR | Security analysis |
| Copilot Review | `copilot-review.yml` | PR | AI code review |
| Docker Build | `docker-build-push.yml` | Push to main | Build & push images |
| Release | `release.yml` | Tag push | Create releases |
| Integration Tests | `integration-tests.yml` | Push, PR | Test all services |
| Branch Aware CI | `branch-aware-ci.yml` | Push | Self-test |

---

## 🎨 Customize Configuration

Edit `.branchci.yml` in your projects:

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches:
      - main
    variables:
      ENV: production
      API_URL: https://api.production.com

  staging:
    name: staging
    requires_approval: false
    allowed_branches:
      - staging
    variables:
      ENV: staging
      API_URL: https://api.staging.com

branch_mappings:
  - pattern: main
    environment: production
    actions: [test, deploy, notify]
    priority: 100
  
  - pattern: feature/*
    environment: development
    actions: [test]
    priority: 50
```

---

## 📖 Documentation

- **Main README**: `/README.md`
- **Architecture**: `/docs/monolithic-docs/ARCHITECTURE.md`
- **Configuration**: `/docs/monolithic-docs/CONFIGURATION.md`
- **Use Cases**: `/docs/monolithic-docs/USE-CASES.md`
- **Microservices**: `/docs/microservices-docs/MICROSERVICES_README.md`

---

## 🆘 Troubleshooting

### Build Fails
```bash
# Clean and rebuild
make clean
make build
```

### Import Errors
```bash
# Update dependencies
go mod tidy
go mod download
```

### Docker Issues
```bash
# Rebuild images
docker-compose build --no-cache
docker-compose up -d
```

### GitHub Action Not Found
- Make sure you've pushed to GitHub
- Tag must be created for versioned use
- For latest: `uses: NadeeshaMedagama/branch_aware_ci@main`

---

## 🎯 Next Steps

1. ✅ Create GitHub repository
2. ✅ Push code to GitHub
3. ✅ Set up Docker Hub repositories
4. ✅ Configure GitHub secrets
5. ✅ Create first release (v1.0.0)
6. ✅ Test the GitHub Action in another repo
7. ✅ Star your own repository 🌟
8. ✅ Share with the community!

---

## 📞 Support

- **Issues**: https://github.com/NadeeshaMedagama/branch_aware_ci/issues
- **Discussions**: https://github.com/NadeeshaMedagama/branch_aware_ci/discussions
- **Author**: Nadeesha Medagama

---

**Happy CI/CD Automation! 🚀**

