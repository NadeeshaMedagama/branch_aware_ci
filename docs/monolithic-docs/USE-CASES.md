# Use Cases & Examples

Real-world scenarios and solutions using Branch-Aware CI.

## Table of Contents

- [Use Case 1: Auto-Deploy to Multiple Environments](#use-case-1-auto-deploy-to-multiple-environments)
- [Use Case 2: Feature Branch Preview Deployments](#use-case-2-feature-branch-preview-deployments)
- [Use Case 3: Hotfix Emergency Pipeline](#use-case-3-hotfix-emergency-pipeline)
- [Use Case 4: Release Candidate Testing](#use-case-4-release-candidate-testing)
- [Use Case 5: Monorepo Multi-Service Deployment](#use-case-5-monorepo-multi-service-deployment)

## Use Case 1: Auto-Deploy to Multiple Environments

**Scenario**: You want main → production, staging → staging, and feature → dev environments automatically.

### Configuration (`.branchci.yml`)

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main]
    variables:
      ENV: production
      API_URL: https://api.example.com
    notify_on_deploy: true

  staging:
    name: staging
    allowed_branches: [staging]
    variables:
      ENV: staging
      API_URL: https://staging-api.example.com

  development:
    name: development
    allowed_branches: [feature/*]
    variables:
      ENV: development
      API_URL: https://dev-api.example.com

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
    actions: [test, deploy]
    priority: 50

policies:
  require_tests: true
  auto_deploy_branches: [main, staging, feature/*]
```

### Workflow (`.github/workflows/deploy.yml`)

```yaml
name: Auto Deploy

on:
  push:
    branches: ['**']

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Deploy
        if: steps.branch.outputs.should_deploy == 'true'
        run: |
          ./deploy.sh ${{ steps.branch.outputs.environment }}
        env:
          API_URL: ${{ steps.branch.outputs.variables.API_URL }}
```

## Use Case 2: Feature Branch Preview Deployments

**Scenario**: Create isolated preview environments for each feature branch.

### Configuration

```yaml
environments:
  preview:
    name: preview
    allowed_branches: [feature/*]
    variables:
      ENV: preview

branch_mappings:
  - pattern: feature/*
    environment: preview
    actions: [test, deploy, create-preview-url]
    priority: 50
```

### Workflow

```yaml
name: Preview Deployment

on:
  push:
    branches: ['feature/**']

jobs:
  preview:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Extract feature name
        run: |
          FEATURE_NAME=$(echo "${{ steps.branch.outputs.branch_name }}" | sed 's/feature\///')
          echo "FEATURE_NAME=$FEATURE_NAME" >> $GITHUB_ENV
      
      - name: Deploy preview
        run: |
          # Deploy to subdomain: feature-name.preview.example.com
          ./deploy-preview.sh $FEATURE_NAME
      
      - name: Comment PR with preview URL
        uses: actions/github-script@v7
        with:
          script: |
            const url = `https://${process.env.FEATURE_NAME}.preview.example.com`;
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `🚀 Preview deployed to: ${url}`
            });
```

## Use Case 3: Hotfix Emergency Pipeline

**Scenario**: Fast-track hotfixes with minimal checks but notifications.

### Configuration

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main, hotfix/*]

  staging:
    name: staging
    allowed_branches: [hotfix/*]

branch_mappings:
  # Hotfixes: quick path to staging
  - pattern: hotfix/*
    environment: staging
    actions: [test, deploy, notify-urgent]
    priority: 95

  # Main: normal production path
  - pattern: main
    environment: production
    actions: [test, security-scan, deploy, notify]
    priority: 100

policies:
  require_tests: true
  auto_deploy_branches: [hotfix/*]
```

### Workflow

```yaml
name: Hotfix Pipeline

on:
  push:
    branches: ['hotfix/**']

jobs:
  hotfix:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Fast tests
        run: npm test -- --quick
      
      - name: Deploy to staging
        if: steps.branch.outputs.environment == 'staging'
        run: ./deploy.sh staging
      
      - name: Urgent notification
        if: contains(steps.branch.outputs.actions, 'notify-urgent')
        uses: slackapi/slack-github-action@v1
        with:
          payload: |
            {
              "text": "🚨 HOTFIX deployed to staging: ${{ github.event.head_commit.message }}"
            }
```

## Use Case 4: Release Candidate Testing

**Scenario**: Test release candidates in staging before production.

### Configuration

```yaml
environments:
  production:
    name: production
    requires_approval: true
    allowed_branches: [main]

  staging:
    name: staging
    allowed_branches: [release/*]

branch_mappings:
  - pattern: release/*
    environment: staging
    actions: [test, integration-test, deploy, smoke-test]
    priority: 85

  - pattern: main
    environment: production
    actions: [deploy, notify, tag-release]
    priority: 100
```

### Workflow

```yaml
name: Release Pipeline

on:
  push:
    branches: ['release/**']

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Run tests
        run: npm test
      
      - name: Integration tests
        if: contains(steps.branch.outputs.actions, 'integration-test')
        run: npm run test:integration

  deploy-staging:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1
      
      - name: Deploy to staging
        run: ./deploy.sh staging
      
      - name: Smoke tests
        if: contains(steps.branch.outputs.actions, 'smoke-test')
        run: npm run test:smoke -- --url https://staging.example.com
```

## Use Case 5: Monorepo Multi-Service Deployment

**Scenario**: Different services in a monorepo with service-specific branches.

### Configuration

```yaml
environments:
  production:
    name: production
    requires_approval: true

  staging:
    name: staging

branch_mappings:
  # Frontend service
  - pattern: frontend/*
    environment: staging
    actions: [test-frontend, deploy-frontend]
    priority: 70

  # Backend service
  - pattern: backend/*
    environment: staging
    actions: [test-backend, deploy-backend]
    priority: 70

  # Infrastructure
  - pattern: infra/*
    environment: staging
    actions: [validate-terraform, plan-terraform]
    priority: 70

  # Main: deploy all services
  - pattern: main
    environment: production
    actions: [test-all, deploy-all]
    priority: 100
```

### Workflow

```yaml
name: Monorepo Deploy

on:
  push:
    branches: ['**']

jobs:
  detect:
    runs-on: ubuntu-latest
    outputs:
      actions: ${{ steps.branch.outputs.actions }}
      environment: ${{ steps.branch.outputs.environment }}
    steps:
      - uses: actions/checkout@v4
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1

  deploy-frontend:
    needs: detect
    if: contains(needs.detect.outputs.actions, 'deploy-frontend')
    runs-on: ubuntu-latest
    steps:
      - name: Deploy frontend
        run: ./deploy-frontend.sh ${{ needs.detect.outputs.environment }}

  deploy-backend:
    needs: detect
    if: contains(needs.detect.outputs.actions, 'deploy-backend')
    runs-on: ubuntu-latest
    steps:
      - name: Deploy backend
        run: ./deploy-backend.sh ${{ needs.detect.outputs.environment }}

  deploy-all:
    needs: detect
    if: contains(needs.detect.outputs.actions, 'deploy-all')
    runs-on: ubuntu-latest
    steps:
      - name: Deploy all services
        run: ./deploy-all.sh ${{ needs.detect.outputs.environment }}
```

## Bonus: Security Scanning for Production Only

```yaml
name: Conditional Security Scan

on:
  push:
    branches: ['**']

jobs:
  detect:
    runs-on: ubuntu-latest
    outputs:
      environment: ${{ steps.branch.outputs.environment }}
    steps:
      - uses: actions/checkout@v4
      - id: branch
        uses: NadeeshaMedagama/branch_aware_ci@v1

  security:
    needs: detect
    if: needs.detect.outputs.environment == 'production'
    runs-on: ubuntu-latest
    steps:
      - name: Security scan
        run: npm audit --production
      
      - name: Container scan
        run: trivy image myapp:latest
```

## Best Practices

1. **Use environment protection rules** in GitHub for production
2. **Set up branch protection** for main/master
3. **Use secrets per environment** in GitHub Settings
4. **Add manual approval gates** for production deployments
5. **Send notifications** for production changes
6. **Tag releases** when deploying to production

## More Examples

See the [`.github/workflows/`](../.github/workflows/) directory for complete working examples.

