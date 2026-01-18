# Imported Workflows from Major Open Source Projects

This directory contains GitHub Actions workflows imported from major open source repositories, converted to Go code using `wetwire-github import`.

## Purpose

These examples demonstrate:
- Real-world workflow patterns from production repositories
- Complex CI/CD pipelines with multiple jobs and steps
- Matrix builds, conditional execution, and advanced GitHub Actions features
- How existing YAML workflows translate to typed Go declarations

## Imported Workflows

### Docker Compose (docker/compose)

Repository: [docker/compose](https://github.com/docker/compose)
License: Apache-2.0

**Workflows:**
- `ci.yml` - Comprehensive CI pipeline with multi-platform builds and testing
- `merge.yml` - Merge workflow for automated PR processing

### Grafana (grafana/grafana)

Repository: [grafana/grafana](https://github.com/grafana/grafana)
License: AGPL-3.0

**Workflows:**
- `backport-workflow.yml` - Automated backporting to stable branches
- `codeql-analysis.yml` - Security analysis with CodeQL

### HashiCorp Terraform (hashicorp/terraform)

Repository: [hashicorp/terraform](https://github.com/hashicorp/terraform)
License: BUSL-1.1

**Workflows:**
- `build.yml` - Multi-platform build pipeline with cross-compilation
- `checks.yml` - Code quality checks, linting, and validation

### Prometheus (prometheus/prometheus)

Repository: [prometheus/prometheus](https://github.com/prometheus/prometheus)
License: Apache-2.0

**Workflows:**
- `ci.yml` - Extensive CI pipeline with matrix testing across platforms
- `codeql-analysis.yml` - Security scanning with GitHub CodeQL

## Structure

Each imported workflow is organized as:

```
examples/imported/{org}/{repo}/{workflow-name}/
├── workflows/
│   ├── workflows.go  # Main workflow definitions
│   ├── jobs.go       # Job configurations
│   ├── steps.go      # Step definitions
│   └── triggers.go   # Event triggers and conditions
```

## License Information

The workflows in this directory are imported from open source projects and retain their original licenses:

- **Docker Compose**: Apache License 2.0
- **Grafana**: GNU Affero General Public License v3.0
- **HashiCorp Terraform**: Business Source License 1.1
- **Prometheus**: Apache License 2.0

These examples are provided for educational and demonstration purposes. The imported Go code is licensed under MIT (the license of wetwire-github-go), but the original workflow logic and structure are subject to their respective licenses.

## Usage

These imports are **read-only examples**. They demonstrate the import functionality but are not intended to be executed directly. To use these as templates:

1. Copy the relevant workflow directory
2. Modify the Go code to suit your needs
3. Build with `wetwire-github build` to generate YAML

## Import Command

These workflows were imported using:

```bash
wetwire-github import <workflow.yml> --output examples/imported/{org}/{repo}/{workflow-name} --no-scaffold
```

## Statistics

Total imports:
- 8 workflows
- 42 jobs
- 201 steps

Last updated: 2026-01-17
