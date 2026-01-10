# Changelog

All notable changes to wetwire-github-go are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- **github-script Action Wrapper** - Typed wrapper for `actions/github-script@v7` (#91)
  - Run JavaScript in workflows with access to GitHub API
  - Supports inputs: `script`, `github-token`, `debug`, `user-agent`, `previews`, `result-encoding`, `retries`, `retry-exempt-status-codes`
- **Additional Lint Rules (WAG009-WAG012)** - 4 new lint rules for workflow quality (#86)
  - WAG009: Validate matrix dimension values (empty matrix detection)
  - WAG010: Flag missing recommended action inputs (e.g., GoVersion for setup_go)
  - WAG011: Detect potential unreachable jobs (undefined job dependencies)
  - WAG012: Warn about deprecated action versions (e.g., checkout@v2 → v4)
- **Additional Action Wrappers** - 5 new typed wrappers for popular GitHub Actions (#85)
  - `actions/setup_java` - Set up Java JDK with distribution selection
  - `actions/setup_dotnet` - Set up .NET SDK
  - `actions/setup_ruby` - Set up Ruby with Bundler support
  - `actions/setup_rust` - Set up Rust toolchain (dtolnay/rust-toolchain)
  - `actions/gh_release` - Create GitHub releases (softprops/action-gh-release)
- **CODEOWNERS** - Full support for GitHub CODEOWNERS files (#84)
  - `codeowners.Owners` type with Rules (Pattern, Owners, Comment)
  - AST-based discovery of Owners declarations
  - Runner extraction for codeowners rules
  - `build --type codeowners` to generate CODEOWNERS file
  - Serialization to CODEOWNERS text format with comments
- **PR Templates** - Full support for GitHub Pull Request templates (#83)
  - `templates.PRTemplate` type with Name, Content, and Filename() method
  - AST-based discovery of PRTemplate declarations
  - Runner extraction for PRTemplate content
  - `build --type pr-template` to generate markdown files
  - Support for default (PULL_REQUEST_TEMPLATE.md) and named templates (PULL_REQUEST_TEMPLATE/{name}.md)
- `docs/ROADMAP.md` - Feature matrix and implementation status
- `lint --fix` - Automatic fixing for WAG001 (raw uses: strings → typed action wrappers)
- Documentation examples for all 7 action wrapper packages in QUICK_START.md
- `examples/ci-workflow` - Complete CI workflow example with matrix strategy and typed action wrappers
- `design` command now fully integrated with wetwire-core-go for AI-assisted workflow generation
  - Interactive AI session using Anthropic Claude API
  - Streaming output support (`--stream`)
  - Lint enforcement with configurable max cycles (`--max-lint-cycles`)
  - 7 agent tools: init_package, write_file, read_file, run_lint, run_build, run_validate, ask_developer
- `test` command enhanced with persona-based testing and scoring (#67)
  - 5 developer personas: beginner, intermediate, expert, terse, verbose
  - 5-dimension scoring system (0-3 each, max 15): Completeness, Lint Quality, Code Quality, Output Validity, Question Efficiency
  - Thresholds: 0-5 Failure, 6-9 Partial, 10-12 Success, 13-15 Excellent
  - `--persona` flag to run specific persona
  - `--score` flag to show scoring breakdown
  - `--list` flag to list available personas and scenarios
- MCP server support for IDE integration (#68)
  - `design --mcp-server` flag starts MCP protocol server over stdio
  - 4 MCP tools: wetwire_init, wetwire_lint, wetwire_build, wetwire_validate
  - Enables integration with Kiro, Claude Desktop, and other MCP-compatible IDEs

### Fixed
- Documentation accuracy: project structures, CLI flags, status tables
- Runner now resolves relative paths in replace directives for proper module extraction

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
- **PR Templates** - `.github/PULL_REQUEST_TEMPLATE.md` or `.github/PULL_REQUEST_TEMPLATE/*.md`
- **CODEOWNERS** - `.github/CODEOWNERS`

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
