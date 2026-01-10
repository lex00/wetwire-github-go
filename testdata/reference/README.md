# Reference GitHub Workflows

This directory contains reference GitHub Actions workflow files used for testing the import/export functionality of wetwire-github-go.

## Purpose

These workflows are used in round-trip tests to verify that:
1. YAML workflows can be correctly imported to Go intermediate representation
2. Go code can be generated from the IR
3. The generated Go code can be compiled and executed
4. The executed code produces YAML output that is semantically equivalent to the original

## Workflow Sources

These workflows are based on official GitHub starter workflows from:
https://github.com/actions/starter-workflows

They have been slightly modified to:
- Use concrete branch names instead of placeholders (e.g., `main` instead of `$default-branch`)
- Use concrete cron expressions instead of placeholders
- Simplify matrix configurations for testing purposes

## Workflows Included

### go.yml
A simple Go CI workflow that:
- Checks out code
- Sets up Go
- Builds and tests the project

### nodejs.yml
A Node.js CI workflow that:
- Tests across multiple Node.js versions using matrix strategy
- Runs build and test commands

### docker-publish.yml
A Docker workflow that:
- Builds and publishes Docker images to GHCR
- Signs images with cosign
- Includes multiple complex steps with conditionals

### codeql.yml
A CodeQL security scanning workflow that:
- Runs on a schedule
- Analyzes code for security vulnerabilities
- Uses matrix strategy for multiple languages

## Usage in Tests

See `/internal/importer/roundtrip_test.go` for the test implementation that uses these reference workflows.
