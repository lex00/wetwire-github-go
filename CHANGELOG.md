# Changelog

All notable changes to wetwire-github-go are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Changed
- **Remove .ToStep() Method** - Action wrappers now implement `StepAction` interface (#110)
  - Action wrappers can be used directly in `Job.Steps` slice
  - Changed `Job.Steps` type from `[]Step` to `[]any` to support both `Step` and action wrappers
  - Action wrappers implement `Action() string` and `Inputs() map[string]any` methods
  - Conversion to Step happens during serialization, not at declaration time
  - Pure struct-literal pattern for AI-friendly, declarative code

### Added
- **Improved Test Coverage** - Major test coverage improvements across core packages (#125, #126, #127, #130, #131, #135, #137)
  - `internal/template`: 35.7% → 90.9% (55+ percentage point improvement)
  - `internal/importer`: 37.9% → 96.0% (58+ percentage point improvement)
  - `workflow`: 43.3% → 100.0% (57+ percentage point improvement)
  - `internal/runner`: 20.9% → 80.2% (59+ percentage point improvement) (#121, #135)
  - `internal/serialize`: 31.7% → 92.6% (61+ percentage point improvement) (#121)
  - `internal/validation`: 55.3% → 89.5% (34+ percentage point improvement) (#130)
  - `internal/discover`: 61.0% → 93.8% (33+ percentage point improvement) (#131)
  - 188 test functions in workflow package covering all expression contexts, conditions, triggers
  - Action wrapper test coverage improvements (#137, #140):
    - `actions/checkout`: 70.7% → 100% (29+ percentage point improvement)
    - `actions/gcp_auth`: 69.7% → 100% (30+ percentage point improvement)
    - `actions/gcp_deploy_cloudrun`: 74.4% → 100% (26+ percentage point improvement)
    - `actions/aws_configure_credentials`: 75.6% → 100% (24+ percentage point improvement)
    - `actions/codeql_init`: 76.2% → 100% (24+ percentage point improvement)
    - `actions/setup_go`: 76.5% → 100% (24+ percentage point improvement)
    - `actions/setup_java`: 78.0% → 100% (22+ percentage point improvement)
    - `actions/setup_ruby`: 81.0% → 100% (19+ percentage point improvement) (#140)
    - `actions/azure_webapps_deploy`: 81.5% → 100% (19+ percentage point improvement) (#140)
    - `actions/golangci_lint`: 84.0% → 100% (16+ percentage point improvement) (#140)
    - `actions/azure_login`: 85.7% → 100% (14+ percentage point improvement) (#140)
    - `actions/setup_dotnet`: 85.7% → 100% (14+ percentage point improvement) (#140)
  - `internal/agent`: 24.7% → 41.8% (17+ percentage point improvement) (#139)
    - Added 46 test functions covering checkCompletionGate, executeTool routing, AskDeveloper, state management
  - `internal/linter`: 83.9% → 91.0% (7+ percentage point improvement) (#143)
  - `codegen`: 85.2% → 89.5% (4+ percentage point improvement) (#143)
  - Additional action wrapper test coverage (#142):
    - `actions/aws_ecr_login`: 92.3% → 100% (8+ percentage point improvement)
    - `actions/docker_build_push`: 94.6% → 100% (5+ percentage point improvement)
    - `actions/gh_release`: 97.0% → 100% (3+ percentage point improvement)
    - `actions/codecov`: 97.1% → 100% (3+ percentage point improvement)
- **Documentation** - Additional contributor guides (#133, #134)
  - `docs/CODEGEN.md` - Action wrapper code generation guide
  - `CONTRIBUTING.md` - Contribution guidelines with development setup, PR process, testing requirements
- **Additional Action Wrappers** - 4 new typed wrappers for popular GitHub Actions (#128)
  - `actions/labeler@v5` - Automatically label PRs based on file paths
  - `actions/stale@v9` - Mark and close stale issues/PRs
  - `peter-evans/create-pull-request@v6` - Create pull requests programmatically
  - `actions/dependency-review-action@v4` - Dependency vulnerability scanning
- **Security Scanning Action Wrappers** - 4 typed wrappers for security scanning (#120)
  - `github/codeql-action/init@v3` - Initialize CodeQL for code scanning
  - `github/codeql-action/analyze@v3` - Run CodeQL analysis and upload results
  - `aquasecurity/trivy-action@0.28.0` - Vulnerability scanning for containers and filesystems
  - `ossf/scorecard-action@v2.4.0` - OpenSSF Scorecard security assessment
- **Golangci-lint Action Wrapper** - Typed wrapper for golangci-lint-action@v6 (#119)
  - Run golangci-lint for Go code linting
  - Supports: `version`, `working-directory`, `args`, `only-new-issues`, cache options
- **GCP Action Wrappers** - 3 typed wrappers for Google Cloud GitHub Actions (#107)
  - `google-github-actions/auth@v2` - Authenticate to Google Cloud with OIDC or service account keys
  - `google-github-actions/deploy-cloudrun@v2` - Deploy containers or source to Cloud Run services/jobs
  - `google-github-actions/setup-gcloud@v2` - Set up and configure the Google Cloud SDK (gcloud)
- **Azure Action Wrappers** - 3 typed wrappers for Azure GitHub Actions (#106)
  - `azure/login@v2` - Login to Azure with service principal or OIDC
  - `azure/webapps-deploy@v3` - Deploy to Azure Web Apps or Web App for Containers
  - `azure/docker-login@v2` - Login to Azure Container Registry
- **AWS Action Wrappers** - 2 typed wrappers for AWS GitHub Actions (#105)
  - `aws-actions/configure-aws-credentials@v4` - Configure AWS credentials with OIDC or access keys
  - `aws-actions/amazon-ecr-login@v2` - Authenticate to Amazon ECR Private or Public registries
- **Documentation** - Additional guides for versioning and adoption (#117)
  - `VERSIONING.md` - Action wrapper versions, breaking changes policy, deprecation timeline
  - `ADOPTION.md` - Migration strategies, best practices, team adoption checklist
- **CODEOWNERS Import** - Import existing CODEOWNERS files to Go code (#104)
  - Parse CODEOWNERS format (pattern + owners)
  - Handle inline and full-line comments
  - Generate `codeowners.Owners` declarations with Rules slice
  - Support `--type codeowners` option in import command
- **PR Template Import** - Import existing PULL_REQUEST_TEMPLATE.md files to Go code (#103)
  - Parse markdown PR templates
  - Generate `templates.PRTemplate` declarations
  - Support multiline content with raw string literals
- **EXAMPLES.md** - Workflow examples catalog (#96)
  - Basic CI workflow
  - Multi-language matrix builds
  - Docker build and push to GHCR
  - Release workflow with changelog
  - Monorepo with path filters
  - Scheduled maintenance tasks
  - PR labeling with github-script
  - Multi-environment deployments
- **INTERNALS.md** - Architecture documentation (#95)
  - AST discovery system with Mermaid diagram
  - Template generation and multi-artifact output
  - Serialization and expression handling
  - Linter architecture (12 rules)
  - Importer architecture
  - Agent integration and scoring system
- **DEVELOPERS.md** - Comprehensive developer guide (#94)
  - Development setup and prerequisites
  - Project structure overview
  - Guide for adding action wrappers
  - Guide for adding lint rules
  - Contributing guidelines
- **Codecov Action Wrapper** - Typed wrapper for `codecov/codecov-action@v5` (#93)
  - Upload code coverage reports to Codecov
  - Supports: `token`, `files`, `directory`, `flags`, `name`, `fail_ci_if_error`, `verbose`, `working-directory`, `env_vars`, `use_oidc`
- **Docker Action Wrappers** - 3 typed wrappers for container CI/CD (#92)
  - `docker/login-action@v3` - Authenticate with Docker registries (Docker Hub, GHCR, ECR)
  - `docker/build-push-action@v6` - Build and push images with Buildx support
  - `docker/setup-buildx-action@v3` - Set up Docker Buildx for multi-platform builds
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
