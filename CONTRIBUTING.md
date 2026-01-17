# Contributing to wetwire-github-go

Thank you for your interest in contributing to wetwire-github-go! This document provides guidelines and instructions for contributing.

## Welcome

wetwire-github-go generates GitHub YAML configurations (Actions workflows, Dependabot configs, Issue Templates) from typed Go declarations. We welcome contributions of all kinds: bug fixes, new features, documentation improvements, and action wrapper additions.

## Development Setup

### Prerequisites

- **Go 1.23+** (required)
- **git** (version control)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/lex00/wetwire-github-go.git
cd wetwire-github-go

# Download dependencies
go mod download

# Build the CLI
go build -o wetwire-github ./cmd/wetwire-github

# Verify the installation
./wetwire-github version
```

## How to Submit Issues

When submitting issues, please include:

1. **Clear title** - Summarize the issue in a few words
2. **Description** - Detailed explanation of the bug or feature request
3. **Environment** - Go version, OS, and wetwire-github-go version
4. **Reproduction steps** - For bugs, provide minimal steps to reproduce
5. **Expected behavior** - What you expected to happen
6. **Actual behavior** - What actually happened
7. **Code samples** - Include relevant Go code or YAML output if applicable

Use appropriate labels:
- `bug` - Something is not working correctly
- `enhancement` - New feature or improvement
- `documentation` - Documentation updates
- `action-wrapper` - New action wrapper request

## Pull Request Process

1. **Create a feature branch** from `main`
   ```bash
   git checkout -b feat/my-feature
   ```

2. **Write tests first** (TDD approach recommended)

3. **Implement your changes**

4. **Ensure all tests pass**
   ```bash
   go test ./...
   ```

5. **Run static analysis**
   ```bash
   go vet ./...
   ```

6. **Update documentation**
   - Update `CHANGELOG.md` under the "Unreleased" section
   - Update `docs/ROADMAP.md` if adding features

7. **Commit with conventional format**
   ```
   feat: add action wrapper for owner/action
   fix: correct YAML serialization for matrix
   docs: update CLI reference
   test: add tests for WAG012
   ```

8. **Create a Pull Request** with:
   - Clear description of changes
   - Reference to related issues
   - Screenshots or examples if applicable

## Code Style Guidelines

### General

- Run `go fmt` before committing
- Run `go vet` to catch common issues
- Follow existing patterns in the codebase

### wetwire Pattern

When writing Go declarations, follow the "No Parens" pattern:

```go
// Use flat variables, not nested structs
var MyWorkflow = workflow.Workflow{
    Name: "CI",
    On:   MyTriggers,  // Reference another variable
}

var MyTriggers = workflow.Triggers{
    Push: MyPush,
}

var MyPush = workflow.PushTrigger{Branches: List("main")}
```

Key principles:
- **Flat variables** - Extract nested structs into named variables
- **No pointers** - Never use `&` or `*` in declarations
- **Direct references** - Variables reference each other by name
- **Struct literals only** - No function calls in declarations

### Commit Messages

Use conventional commit format:

| Prefix | Description |
|--------|-------------|
| `feat:` | New feature |
| `fix:` | Bug fix |
| `docs:` | Documentation only |
| `test:` | Adding or updating tests |
| `refactor:` | Code restructuring |
| `chore:` | Maintenance tasks |

## Testing Requirements

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/lint/...

# Run specific test
go test -v ./internal/lint/... -run TestWAG001

# Run with race detection
go test -race ./...
```

### Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Coverage targets:**
- New code should include tests
- Aim for meaningful coverage of logic paths
- Action wrappers require tests for `Action()` method and input handling

### Test Requirements for PRs

- All existing tests must pass
- New features require corresponding tests
- Bug fixes should include regression tests

## Adding Action Wrappers

See the [Developer Guide](docs/DEVELOPERS.md#adding-action-wrappers) for detailed instructions on creating type-safe action wrappers.

Quick overview:

1. Create package in `actions/my_action/`
2. Implement `Action()` and `ToStep()` methods
3. Add tests in `my_action_test.go`
4. Update documentation

## Adding Lint Rules

See the [Developer Guide](docs/DEVELOPERS.md#adding-lint-rules) for detailed instructions on adding new lint rules (WAG001-WAGXXX).

## Technical Documentation

For detailed technical documentation, architecture details, and advanced contribution topics, see:

- [Developer Guide](docs/DEVELOPERS.md) - Comprehensive development guide
- [Internals](docs/INTERNALS.md) - Architecture and implementation details
- [CLI Reference](docs/CLI.md) - CLI commands and options
- [FAQ](docs/FAQ.md) - Common questions

## License

By contributing to wetwire-github-go, you agree that your contributions will be licensed under the MIT License.

## Questions?

If you have questions about contributing, feel free to:
- Open a discussion on GitHub
- Check the [FAQ](docs/FAQ.md) for common questions
- Review existing issues and PRs for context

Thank you for contributing!
