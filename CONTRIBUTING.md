# Contributing to GoMeta Redirector

Thank you for your interest in contributing to GoMeta Redirector! This document provides guidelines and instructions for contributing.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Community](#community)

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it to understand what behavior is expected.

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the bug**
- **Provide specific examples**
- **Include Go version and OS information**
- **Include relevant logs or screenshots**

### 💡 Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

- **Clear use case description**
- **Expected behavior**
- **Potential implementation approach** (optional)

### 🔧 Pull Requests

1. Fork the repo and create your branch from `main`
2. Make your changes
3. Ensure tests pass (if any)
4. Update documentation as needed
5. Submit the pull request

## Development Setup

### Prerequisites

- **Go 1.21+** - [Install Go](https://go.dev/doc/install)
- **Docker** (optional) - [Install Docker](https://docs.docker.com/get-docker/)
- **Git** - [Install Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Local Development

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/go-meta-redirector.git
cd go-meta-redirector

# Add upstream remote
git remote add upstream https://github.com/pilab-dev/go-meta-redirector.git

# Build
go build -o go-meta-redirector .

# Run locally
./go-meta-redirector :8080

# Test changes
curl -H "Host: go.pilab.hu" "http://localhost:8080/cloud/log?go-get=1"
```

### Running Tests

```bash
# Run tests (when available)
go test ./...

# Run with verbose output
go test -v ./...
```

## Pull Request Process

1. **Update Documentation** - Update README.md or relevant docs if needed
2. **Add Tests** - Add tests for new functionality
3. **Update CHANGELOG** - Add entry to CHANGELOG.md (if maintained)
4. **Follow Go Conventions** - Use `gofmt` and follow standard Go idioms
5. **One Feature Per PR** - Keep PRs focused on a single feature or fix
6. **Link Issues** - Reference relevant issues in PR description

### PR Title Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add new domain support`
- `fix: correct fallback pattern matching`
- `docs: update configuration examples`
- `refactor: simplify handler logic`

## Coding Standards

### Go Style Guide

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Run `go vet` to catch common mistakes
- Keep functions small and focused
- Write meaningful comments for exported functions

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

Example:
```
feat(config): add support for multiple fallback patterns

- Implement wildcard matching for fallback rules
- Add validation for pattern syntax
- Update documentation with examples

Closes #123
```

## Community

- **GitHub Issues** - Report bugs or suggest features
- **Pull Requests** - Submit contributions
- **Discussions** - Join conversations (if enabled)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to GoMeta Redirector! 🚀
