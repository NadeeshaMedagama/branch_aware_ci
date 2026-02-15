# Branch-Aware CI - Microservices Architecture

## 🏗️ Architecture Overview

This project has been transformed into a microservices architecture following SOLID principles.

### Services

```
┌─────────────────────────────────────────────────────────┐
│                     API Gateway                         │
│                  (Port: 8080)                           │
│  • Request routing                                      │
│  • Service orchestration                               │
│  • API aggregation                                     │
└────────────┬───────────────────┬──────────────┬─────────┘
             │                   │              │
             ▼                   ▼              ▼
    ┌────────────────┐  ┌────────────────┐  ┌──────────────┐
    │ Branch Detector│  │ Policy Engine  │  │Config Service│
    │  (Port: 8081)  │  │  (Port: 8082)  │  │(Port: 8083)  │
    │                │  │                │  │              │
    │ • Git ops      │  │ • Evaluation   │  │• Config mgmt │
    │ • Detection    │  │ • Rules        │  │• Validation  │
    └────────────────┘  └────────────────┘  └──────────────┘
```

### SOLID Principles Implementation

**Single Responsibility Principle**
- Each service has one clear purpose
- Branch Detector: Only Git operations
- Policy Engine: Only policy evaluation
- API Gateway: Only routing and orchestration

**Open/Closed Principle**
- Services are open for extension via interfaces
- New branch patterns can be added without modifying core logic
- New policy rules can be added without changing the engine

**Liskov Substitution Principle**
- All services implement defined interfaces
- Mock implementations can replace real ones for testing

**Interface Segregation Principle**
- Small, focused interfaces (IBranchDetector, IPolicyEngine, etc.)
- Services depend only on what they need

**Dependency Inversion Principle**
- Services depend on abstractions (interfaces), not concrete implementations
- Easy to swap implementations without changing consumers

## 🚀 Quick Start

### Using Docker Compose (Recommended)

```bash
# Start all services
docker-compose up -d

# Check service health
curl http://localhost:8080/health
curl http://localhost:8081/health
curl http://localhost:8082/health

# Test the API
curl -X POST http://localhost:8080/api/v1/analyze \
  -H "Content-Type: application/json" \
  -d '{"repo_path": "."}'

# Stop services
docker-compose down
```

### Using Make

```bash
# Build all services
make build-services

# Run services individually
make run-gateway
make run-branch-detector
make run-policy-engine
```

## 📦 Service Details

### API Gateway (Port 8080)

**Endpoints:**
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `POST /api/v1/analyze` - Full branch analysis (orchestrates all services)

**Environment Variables:**
- `PORT` - HTTP port (default: 8080)
- `BRANCH_DETECTOR_URL` - Branch detector service URL
- `POLICY_ENGINE_URL` - Policy engine service URL
- `CONFIG_SERVICE_URL` - Config service URL

### Branch Detector (Ports 8081, 50051)

**Endpoints:**
- `GET /health` - Health check
- `POST /api/v1/detect` - Detect branch information

**Environment Variables:**
- `HTTP_PORT` - HTTP port (default: 8081)
- `GRPC_PORT` - gRPC port (default: 50051)

### Policy Engine (Ports 8082, 50052)

**Endpoints:**
- `GET /health` - Health check
- `POST /api/v1/evaluate` - Evaluate policy

**Environment Variables:**
- `HTTP_PORT` - HTTP port (default: 8082)
- `GRPC_PORT` - gRPC port (default: 50052)

## 🔄 GitHub Workflows

### CodeQL Security Analysis
- **File:** `.github/workflows/codeql.yml`
- **Runs:** On push to main/develop, PRs, and weekly
- **Purpose:** Automated security vulnerability scanning

### Copilot Code Review
- **File:** `.github/workflows/copilot-review.yml`
- **Runs:** On pull requests
- **Purpose:** Automated code quality review
- **Checks:**
  - Go formatting (`go fmt`)
  - Go vet issues
  - Static analysis (staticcheck)
  - Security issues (gosec)

### Docker Build & Push
- **File:** `.github/workflows/docker-build-push.yml`
- **Runs:** On push to main/develop and tags
- **Purpose:** Build and push Docker images to Docker Hub
- **Features:**
  - Multi-platform builds (amd64, arm64)
  - Vulnerability scanning with Trivy
  - Integration testing with Docker Compose

### Release
- **File:** `.github/workflows/release.yml`
- **Runs:** On version tags (v*)
- **Purpose:** Create GitHub releases
- **Features:**
  - Auto-generated changelog
  - Multi-platform binaries
  - Docker image publishing
  - GitHub Release creation

### Integration Tests
- **File:** `.github/workflows/integration-tests.yml`
- **Runs:** On push and PRs
- **Purpose:** End-to-end testing
- **Tests:**
  - Unit tests with coverage
  - Service integration tests
  - Health check validation

### Dependabot
- **File:** `.github/dependabot.yml`
- **Purpose:** Automated dependency updates
- **Updates:**
  - Go modules (weekly)
  - GitHub Actions (weekly)
  - Docker base images (weekly)

## 🐳 Docker Images

All services are published to Docker Hub:

```bash
# Pull images
docker pull NadeeshaMedagama/branch-aware-gateway:latest
docker pull NadeeshaMedagama/branch-aware-branch-detector:latest
docker pull NadeeshaMedagama/branch-aware-policy-engine:latest

# Or use specific versions
docker pull NadeeshaMedagama/branch-aware-gateway:v1.0.0
```

## 🔐 Security Features

1. **CodeQL Analysis** - Automated security scanning
2. **Dependabot** - Automated dependency updates
3. **Gosec** - Go security checker
4. **Trivy** - Container vulnerability scanning
5. **SARIF Upload** - Security findings in GitHub Security tab

## 🧪 Testing

```bash
# Run unit tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run integration tests
make test-integration

# Test with Docker Compose
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📊 Monitoring

Prometheus metrics are exposed on each service:
- Gateway: `http://localhost:8080/metrics`
- Branch Detector: `http://localhost:8081/metrics`
- Policy Engine: `http://localhost:8082/metrics`

## 🚀 Deployment

### Local Development
```bash
docker-compose up -d
```

### Kubernetes
```bash
kubectl apply -f k8s/
```

### Production
See deployment guides in `docs/deployment/`

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📝 License

MIT License - see [LICENSE](../../LICENSE) file for details.

---

**Built with microservices architecture and SOLID principles** 🏗️✨

