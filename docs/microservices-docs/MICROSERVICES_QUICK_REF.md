# 🚀 Microservices Quick Reference

## Services Overview

| Service | HTTP Port | gRPC Port | Purpose |
|---------|-----------|-----------|---------|
| Gateway | 8080 | - | API routing & orchestration |
| Branch Detector | 8081 | 50051 | Git branch detection |
| Policy Engine | 8082 | 50052 | Policy evaluation |

## Quick Commands

### Docker Compose
```bash
docker-compose up -d      # Start all services
docker-compose down       # Stop all services
docker-compose logs -f    # View logs
docker-compose ps         # List services
```

### Makefile
```bash
make build-services       # Build all services
make run-gateway         # Run API Gateway
make run-branch-detector # Run Branch Detector
make run-policy-engine   # Run Policy Engine
make test-integration    # Run integration tests
make docker-build-all    # Build all Docker images
make docker-compose-up   # Start with compose
make docker-compose-down # Stop compose
```

### Health Checks
```bash
curl http://localhost:8080/health  # Gateway
curl http://localhost:8081/health  # Branch Detector
curl http://localhost:8082/health  # Policy Engine
```

## API Endpoints

### Gateway (Port 8080)

**Analyze Branch (Full Workflow)**
```bash
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "repo_path": ".",
    "config_path": ".branchci.yml"
  }'
```

**Health & Readiness**
```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

### Branch Detector (Port 8081)

**Detect Branch**
```bash
curl -X POST http://localhost:8081/api/v1/detect \
  -H "Content-Type: application/json" \
  -d '{"repo_path": "."}'
```

### Policy Engine (Port 8082)

**Evaluate Policy**
```bash
curl -X POST http://localhost:8082/api/v1/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "branch_info": {
      "name": "main",
      "type": "main"
    },
    "config": {...}
  }'
```

## GitHub Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| CodeQL | Push, PR, Weekly | Security scanning |
| Copilot Review | PR | Code quality checks |
| Docker Build/Push | Push, Tags | Build & publish images |
| Release | Tags (v*) | Create releases |
| Integration Tests | Push, PR | Test suite |
| Dependabot | Weekly | Dependency updates |

## Docker Images

### Pull Images
```bash
docker pull NadeeshaMedagama/branch-aware-gateway:latest
docker pull NadeeshaMedagama/branch-aware-branch-detector:latest
docker pull NadeeshaMedagama/branch-aware-policy-engine:latest
```

### Build Images
```bash
docker build -f services/gateway/Dockerfile -t branch-aware-gateway .
docker build -f services/branch-detector/Dockerfile -t branch-aware-branch-detector .
docker build -f services/policy-engine/Dockerfile -t branch-aware-policy-engine .
```

## Environment Variables

### Gateway
```bash
PORT=8080
BRANCH_DETECTOR_URL=http://branch-detector:8081
POLICY_ENGINE_URL=http://policy-engine:8082
CONFIG_SERVICE_URL=http://config-service:8083
```

### Branch Detector
```bash
HTTP_PORT=8081
GRPC_PORT=50051
```

### Policy Engine
```bash
HTTP_PORT=8082
GRPC_PORT=50052
```

## SOLID Principles Quick Reference

**Single Responsibility** - One service, one purpose
```go
type BranchDetector struct {}  // Only does branch detection
type PolicyEngine struct {}    // Only does policy evaluation
```

**Open/Closed** - Extend without modification
```go
type IBranchDetector interface {
    DetectBranch(ctx context.Context, path string) (*BranchInfo, error)
}
```

**Liskov Substitution** - Interfaces are substitutable
```go
var detector interfaces.IBranchDetector = detector.NewBranchDetector()
// Can swap with mock implementation
```

**Interface Segregation** - Small, focused interfaces
```go
type IBranchDetector interface { /* only branch detection methods */ }
type IPolicyEngine interface   { /* only policy methods */ }
```

**Dependency Inversion** - Depend on abstractions
```go
type Handler struct {
    detector interfaces.IBranchDetector  // Interface, not concrete
}
```

## Testing

### Unit Tests
```bash
go test ./...
go test -v ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
```

### Integration Tests
```bash
make test-integration
```

### Coverage
```bash
make coverage
open coverage.html
```

## Troubleshooting

### Service Won't Start
```bash
# Check if port is in use
lsof -i :8080
lsof -i :8081
lsof -i :8082

# Check Docker logs
docker-compose logs gateway
docker-compose logs branch-detector
docker-compose logs policy-engine
```

### Service Not Responding
```bash
# Restart service
docker-compose restart gateway

# Check health
curl http://localhost:8080/health
```

### Build Failures
```bash
# Clean and rebuild
make clean
go mod tidy
make build-services
```

## Release Process

### Create Release
```bash
# Tag version
git tag v1.0.0

# Push tag (triggers release workflow)
git push --tags

# Workflow automatically:
# 1. Builds binaries
# 2. Publishes Docker images
# 3. Creates GitHub release
# 4. Generates changelog
```

### Manual Release
```bash
# Build binaries
GOOS=linux GOARCH=amd64 go build -o dist/branch-aware-ci-linux-amd64 .

# Build Docker images
make docker-build-all

# Push to Docker Hub
docker push NadeeshaMedagama/branch-aware-gateway:v1.0.0
```

## Monitoring

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f gateway
docker-compose logs -f branch-detector
docker-compose logs -f policy-engine
```

### Prometheus (Optional)
```bash
# Access Prometheus
open http://localhost:9090

# Query metrics
http://localhost:8080/metrics
http://localhost:8081/metrics
http://localhost:8082/metrics
```

## Security

### Run Security Scans
```bash
# Go security
gosec ./...

# Static analysis
staticcheck ./...

# Vulnerability scan
go list -json -m all | nancy sleuth
```

### Container Scanning
```bash
# Trivy scan
trivy image branch-aware-gateway:latest
trivy image branch-aware-branch-detector:latest
trivy image branch-aware-policy-engine:latest
```

## GitHub Secrets

Required secrets for workflows:
```
DOCKER_USERNAME  # Docker Hub username
DOCKER_PASSWORD  # Docker Hub token
CODECOV_TOKEN   # (Optional) Codecov token
```

## Common Issues

| Issue | Solution |
|-------|----------|
| Port already in use | Change port in docker-compose.yml |
| Service not responding | Check health endpoint |
| Docker build fails | Run `go mod tidy` first |
| Can't connect to service | Ensure service is running |
| Tests failing | Check service dependencies |

## Development Workflow

```bash
# 1. Make changes to code
vim services/gateway/main.go

# 2. Build services
make build-services

# 3. Run tests
make test

# 4. Test with Docker Compose
docker-compose up -d

# 5. Test endpoints
curl http://localhost:8080/api/v1/analyze -X POST \
  -H "Content-Type: application/json" \
  -d '{"repo_path": "."}'

# 6. Stop services
docker-compose down

# 7. Commit and push
git add .
git commit -m "feat: add new feature"
git push
```

## Resources

- **Main Docs:** MICROSERVICES_README.md
- **Architecture:** docs/ARCHITECTURE.md
- **Original Docs:** README.md
- **Quick Start:** QUICKSTART.md

---

**Need help?** Check the full documentation or open an issue!

**Quick Test:** `docker-compose up -d && curl http://localhost:8080/health`

🚀 **Your microservices are ready to rock!**

