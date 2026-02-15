# Contributing to Branch-Aware CI

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## 🚀 Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR-USERNAME/branch-aware-ci.git
   cd branch-aware-ci
   ```
3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## 🛠️ Development Setup

### Prerequisites
- Go 1.23 or higher
- Git
- Docker (optional, for testing the action)

### Installation

```bash
# Install dependencies
go mod download

# Build the project
go build -o branch-aware-ci

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -cover ./...
```

## 📝 Code Guidelines

### Go Code Style
- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Testing
- Write unit tests for new functionality
- Maintain or improve code coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

Example:
```go
func TestDetectBranch(t *testing.T) {
	tests := []struct {
		name     string
		branch   string
		expected string
	}{
		{"main branch", "main", "production"},
		{"feature branch", "feature/auth", "development"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test implementation
		})
	}
}
```

## 🔄 Pull Request Process

1. **Update documentation** if you're changing functionality
2. **Add tests** for new features
3. **Run tests** and ensure they pass
4. **Update CHANGELOG.md** with your changes
5. **Submit PR** with a clear description

### PR Title Format
```
type(scope): description

Examples:
feat(git): add support for GitLab repositories
fix(policy): correct priority sorting logic
docs(readme): update configuration examples
```

### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

## 🧪 Testing Your Changes

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
# Test with a real Git repository
./branch-aware-ci -repo /path/to/test/repo
```

### Testing the GitHub Action
```bash
# Build Docker image
docker build -t branch-aware-ci:test .

# Run in container
docker run -v $(pwd):/repo branch-aware-ci:test -repo /repo
```

## 📚 Documentation

When adding features:
- Update README.md with usage examples
- Add comments to exported functions
- Update configuration examples
- Add workflow examples if relevant

## 🐛 Reporting Bugs

### Before Submitting
- Check existing issues
- Test with the latest version
- Gather relevant information

### Bug Report Template
```markdown
**Describe the bug**
A clear description of the bug.

**To Reproduce**
Steps to reproduce:
1. Run command '...'
2. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.23]
- Branch-Aware CI version: [e.g., 1.0.0]

**Additional context**
Any other relevant information.
```

## 💡 Suggesting Features

We love feature suggestions! Please:
1. Check if it's already been suggested
2. Clearly describe the use case
3. Provide examples of how it would work
4. Explain why it would be useful

## 📋 Code of Conduct

### Our Standards
- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Assume good intentions

### Unacceptable Behavior
- Harassment or discrimination
- Trolling or insulting comments
- Personal or political attacks
- Publishing others' private information

## 🎯 Areas for Contribution

Looking for ideas? Here are some areas:

### High Priority
- [ ] Add unit tests for policy engine
- [ ] Support for GitLab CI/CD
- [ ] Add more output formats
- [ ] Improve error messages

### Medium Priority
- [ ] Web UI for configuration
- [ ] Integration tests
- [ ] Performance optimizations
- [ ] Additional examples

### Good First Issues
- [ ] Add more branch patterns
- [ ] Improve documentation
- [ ] Add workflow examples
- [ ] Fix typos

## 📞 Questions?

- Open a GitHub issue
- Start a discussion
- Check existing documentation

## 🙏 Thank You!

Your contributions make this project better for everyone!

---

**Happy coding! 🚀**

