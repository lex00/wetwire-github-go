# wetwire-github-go Examples

This directory contains example wetwire-github-go projects demonstrating various GitHub Actions workflow patterns.

## Example Categories

### CI/CD Workflows

| Example | Description |
|---------|-------------|
| [ci-workflow](./ci-workflow/) | Basic CI workflow with build and test jobs |
| [deployment-workflow](./deployment-workflow/) | Multi-environment deployment with staging and production |
| [release-workflow](./release-workflow/) | Automated release with changelog and assets |

### Container & Artifact Workflows

| Example | Description |
|---------|-------------|
| [docker-workflow](./docker-workflow/) | Docker build and push to container registries |
| [artifact-pipeline-workflow](./artifact-pipeline-workflow/) | Build artifacts and pass between jobs |
| [container-services-workflow](./container-services-workflow/) | Jobs using service containers (databases, etc.) |

### Matrix & Reusable Workflows

| Example | Description |
|---------|-------------|
| [matrix-workflow](./matrix-workflow/) | Multi-OS, multi-version testing with matrix strategy |
| [reusable-workflow](./reusable-workflow/) | Reusable workflows with workflow_call trigger |
| [monorepo-workflow](./monorepo-workflow/) | Path-filtered workflows for monorepo projects |

### Security & Automation

| Example | Description |
|---------|-------------|
| [security-workflow](./security-workflow/) | Security scanning with CodeQL and dependency review |
| [publishing-workflow](./publishing-workflow/) | Package publishing to npm, PyPI, Maven, etc. |
| [issue-automation-workflow](./issue-automation-workflow/) | Issue labeling and automation |

### Advanced Patterns

| Example | Description |
|---------|-------------|
| [approval-gates-workflow](./approval-gates-workflow/) | Manual approval gates for deployments |
| [environment-promotion-workflow](./environment-promotion-workflow/) | Promote deployments through environments |
| [scheduled-dispatch-workflow](./scheduled-dispatch-workflow/) | Scheduled and manual dispatch triggers |
| [repository-dispatch-example](./repository-dispatch-example/) | Cross-repository workflow triggers |
| [workflow-run-example](./workflow-run-example/) | Trigger on workflow completion events |

## Running Examples

Each example is a self-contained Go module. To build:

```bash
cd examples/ci-workflow
wetwire-github build .
```

This generates `.github/workflows/*.yml` files.

## Example Structure

Each example follows a consistent structure:

```
example-workflow/
├── go.mod              # Go module definition
├── README.md           # Example-specific documentation
├── CLAUDE.md           # AI assistant instructions
├── workflows/
│   ├── ci.go          # Workflow declarations
│   └── jobs.go        # Job and step definitions
└── .github/
    └── workflows/
        └── ci.yml     # Generated output (after build)
```

## Attribution

All examples in this directory are **original, hand-written** wetwire-github-go code created specifically for this project. They are not imports of existing YAML workflows.

For imported/adapted workflows used in testing, see:
- [testdata/reference/](../testdata/reference/) - GitHub starter workflows adapted for round-trip testing (MIT licensed)

## License

These examples are part of the wetwire-github-go project and are available under the same license as the main project.
