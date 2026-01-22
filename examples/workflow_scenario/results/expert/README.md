<picture>
  <source media="(prefers-color-scheme: dark)" srcset="../../../../docs/wetwire-dark.svg">
  <img src="../../../../docs/wetwire-light.svg" width="100" height="67">
</picture>

This project defines a comprehensive CI/CD workflow using wetwire-github Go syntax.

## Workflow Structure

### Triggers
- Push to `main` branch
- Pull requests to `main` branch

### Jobs

#### 1. Build Job
- **Matrix strategy**: Tests across multiple Go versions and operating systems
  - Go versions: 1.23, 1.24
  - OS: ubuntu-latest, macos-latest
- **Steps**:
  - Checkout code (v4)
  - Setup Go (v5) with matrix Go version
  - Cache Go modules
  - Build all packages

#### 2. Test Job
- **Runs on**: ubuntu-latest
- **Steps**:
  - Checkout code
  - Setup Go 1.23
  - Run tests with race detection and coverage output

#### 3. Deploy Staging Job
- **Dependencies**: Requires build and test jobs to complete
- **Condition**: Only runs on main branch
- **Environment**: staging
- **Steps**: Deploy to staging environment

#### 4. Deploy Production Job
- **Dependencies**: Requires build and test jobs to complete
- **Condition**: Only runs on main branch
- **Environment**: production (https://example.com)
- **Steps**: Deploy to production environment

## Building

Generate the GitHub Actions workflow YAML:

```bash
wetwire-github build .
```

This outputs to `.github/workflows/ci.yml`.

## Key Features

- **Type-safe**: All workflow components are strongly typed Go structs
- **Declarative**: Pure struct literals, no function calls
- **Cross-references**: Jobs reference each other directly
- **Matrix builds**: Test across multiple Go versions and platforms
- **Environment gates**: Separate staging and production deployments
- **Conditional execution**: Deployments only run on main branch
