# Changelog

All notable changes to wetwire-github-go are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- `docs/ROADMAP.md` - Feature matrix and implementation status
- `lint --fix` - Automatic fixing for WAG001 (raw uses: strings → typed action wrappers)

### Fixed
- Documentation accuracy: project structures, CLI flags, status tables

## [0.1.0] - 2025-01-06

Initial release with full GitHub Actions workflow generation from typed Go declarations.

### Added

#### Core Features
- **Workflow Types** - Complete type system for GitHub Actions: `Workflow`, `Job`, `Step`, `Triggers`, `Matrix`, `Strategy`, `Permissions`, `Environment`, `Concurrency`
- **YAML Serialization** - Generate valid GitHub Actions YAML from Go struct declarations
- **AST Discovery** - Automatic detection of workflow/job/step variables via Go AST parsing
- **Expression Contexts** - Type-safe builders for `secrets`, `github`, `matrix`, `env`, `needs` contexts
- **Action Wrappers** - Generated typed wrappers for popular actions (`actions/checkout`, `actions/setup-go`, etc.)

#### CLI Commands
- `wetwire-github init` - Create new workflow project with scaffold
- `wetwire-github build` - Generate YAML from Go declarations
- `wetwire-github import` - Convert existing YAML to Go code
- `wetwire-github validate` - Validate generated YAML with actionlint
- `wetwire-github lint` - Check Go code for wetwire best practices (WAG001-WAG008)
- `wetwire-github list` - List discovered workflows and jobs
- `wetwire-github graph` - Generate DAG visualization (Mermaid/DOT)
- `wetwire-github design` - AI-assisted workflow design (requires wetwire-core-go)
- `wetwire-github test` - Structural testing for workflows

#### Config Types
- **GitHub Actions** - `.github/workflows/*.yml`
- **Dependabot** - `.github/dependabot.yml`
- **Issue Templates** - `.github/ISSUE_TEMPLATE/*.yml`
- **Discussion Templates** - `.github/DISCUSSION_TEMPLATE/*.yml`

#### Linter Rules
- `WAG001` - Use typed action wrappers instead of raw `uses:` strings
- `WAG002` - Use condition builders instead of raw expression strings
- `WAG003` - Use secrets context instead of hardcoded strings
- `WAG004` - Use matrix builder instead of inline maps
- `WAG005` - Extract inline structs to named variables
- `WAG006` - Detect duplicate workflow names
- `WAG007` - Flag oversized files (>N jobs)
- `WAG008` - Avoid hardcoded expression strings

#### Import Features
- Parse existing YAML workflows to Go code
- Flatten nested structures to package-level variables
- Map known actions to typed wrappers
- Convert expressions to type-safe builders
- Support for `--single-file` and `--no-scaffold` modes

#### Codegen
- Schema fetching from GitHub Actions schema repository
- Action.yml parsing for input/output discovery
- Generated typed wrappers with `.ToStep()` method

### Internal
- Runner-based value extraction using temporary Go programs
- Template builder with topological sort and cycle detection
- Intermediate representation (IR) for YAML parsing
- Reference example testing with round-trip validation

[Unreleased]: https://github.com/lex00/wetwire-github-go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/lex00/wetwire-github-go/releases/tag/v0.1.0
