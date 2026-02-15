# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-13

### Added
- Initial release of Branch-Aware CI
- Automatic Git branch detection
- Configurable branch-to-environment mappings
- Policy-based decision engine
- Multiple output formats (JSON, YAML, env, GitHub Actions)
- GitHub Action integration
- Docker container support
- CLI tool for local usage
- Default configuration templates
- Comprehensive documentation
- Example workflows for common use cases

### Features
- Smart branch type detection (feature, hotfix, release, main, etc.)
- Metadata extraction from branch names (ticket numbers, feature names)
- Protected branch recognition
- Priority-based pattern matching
- Environment-specific deployment rules
- Approval requirements for production deployments
- Auto-deploy configuration
- Warning system for policy violations

### Supported Environments
- Production (main, master branches)
- Staging (staging, develop branches)
- Development (feature/*, bugfix/*, hotfix/* branches)

## [Unreleased]

### Planned
- Support for GitLab CI/CD
- Support for Bitbucket Pipelines
- Web UI for configuration visualization
- Monorepo path filtering
- Parallel environment deployments
- Rollback strategies
- Integration with Slack/Discord notifications
- Advanced branch naming patterns

