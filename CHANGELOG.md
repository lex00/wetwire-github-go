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
  - `internal/runner`: 20.9% → 89.9% (69+ percentage point improvement) (#121, #135, #146)
  - `internal/serialize`: 31.7% → 92.6% (61+ percentage point improvement) (#121)
  - `internal/validation`: 55.3% → 92.1% (37+ percentage point improvement) (#130, #145)
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
  - `internal/agent`: 24.7% → 55.2% (30+ percentage point improvement) (#139, #155)
    - Added 76+ test functions covering checkCompletionGate, executeTool routing, AskDeveloper, state management
    - Functions at 100%: NewGitHubAgent, checkLintEnforcement, checkCompletionGate, getTools, executeTool, toolReadFile, AskDeveloper
  - `internal/linter`: 83.9% → 91.3% (7+ percentage point improvement) (#143, #149)
  - `codegen`: 85.2% → 95.7% (10+ percentage point improvement) (#143, #158)
  - `internal/runner`: 89.9% → 90.8% (additional edge case tests) (#158)
- **Security Lint Rules (WAG017-WAG018)** - 2 new security-focused lint rules (#160, #161)
  - WAG017: Suggest adding explicit permissions scope for workflow security
  - WAG018: Detect dangerous pull_request_target patterns with checkout actions
- **Agent Test Coverage Improvement** - Increased from 55.2% to 58.8% (#162)
  - Added 97 tests covering tool routing, error handling, state management, completion gate logic
- **Additional Lint Rules (WAG013-WAG016)** - 4 new lint rules for code quality (#149)
  - WAG013: Detect pointer assignments (&Type{}) - wetwire uses value semantics
  - WAG014: Flag jobs without timeout-minutes setting
  - WAG015: Suggest caching for setup-go/node/python actions
  - WAG016: Validate concurrency settings (cancel-in-progress without group)
- **Reference Example Testing** - Round-trip testing with GitHub starter workflows (#151)
  - `testdata/reference/` directory with 4 official GitHub starter workflows
  - Round-trip tests: import YAML → generate Go → compile → execute → export YAML → compare
  - Tests for go.yml, nodejs.yml, docker-publish.yml, codeql.yml
- **Additional Action Wrappers** - 3 typed wrappers for deployment and notifications (#152)
  - `JamesIves/github-pages-deploy-action@v4` - Deploy to GitHub Pages
  - `slackapi/slack-github-action@v1` - Slack workflow notifications
  - `hashicorp/setup-terraform@v3` - Terraform CLI setup
- **Additional Action Wrappers** - 4 typed wrappers for GitHub automation (#154)
  - `peaceiris/actions-gh-pages@v4` - GitHub Pages deployment (alternative)
  - `actions/first-interaction@v1` - First-time contributor greeting
  - `actions/add-to-project@v1` - Add issues/PRs to GitHub Projects
  - `actions/create-github-app-token@v1` - GitHub App token creation
- **Additional Action Wrappers** - 3 typed wrappers for CI/CD (#164, #165, #166)
  - `actions/attest-build-provenance@v1` - SLSA build provenance attestation for supply chain security
  - `docker/metadata-action@v5` - Extract metadata for Docker builds (tags, labels, annotations)
  - `mikepenz/action-junit-report@v4` - Publish JUnit test results as GitHub check runs
- **Additional Action Wrappers** - 3 typed wrappers for releases and Rust (#168, #169)
  - `actions/create-release@v1` - Create GitHub releases (legacy, with migration guide)
  - `actions/upload-release-asset@v1` - Upload assets to GitHub releases (legacy)
  - `actions-rs/toolchain@v1` - Rust toolchain setup with components and targets
- **Security Workflow Example** - Comprehensive security workflow demonstrating WAG017/WAG018 (#170)
  - CodeQL static analysis with security-extended queries
  - Trivy vulnerability scanning with SARIF upload
  - SLSA build provenance attestation
  - Explicit permissions at workflow and job levels
- **Additional Action Wrappers** - 3 typed wrappers for GitHub Pages deployment (#174)
  - `actions/configure-pages@v5` - Configure GitHub Pages settings
  - `actions/deploy-pages@v4` - Deploy artifacts to GitHub Pages
  - `actions/upload-pages-artifact@v3` - Upload artifacts for Pages deployment
- **Scheduled Workflow Example** - Example demonstrating schedule and dispatch triggers (#172)
  - Schedule trigger with cron patterns
  - Workflow dispatch with typed inputs (choice, boolean, string)
  - Conditional job execution based on event type
- **Reusable Workflow Example** - Example demonstrating workflow_call trigger (#176)
  - Reusable workflow with typed inputs, outputs, and secrets
  - Caller workflow showing how to invoke reusable workflows
  - Job outputs mapping to workflow outputs
- **Additional Action Wrappers** - 3 typed wrappers for code quality tools (#177)
  - `reviewdog/action-setup@v1` - Setup reviewdog for inline code review comments
  - `super-linter/super-linter@v7` - GitHub Super-Linter for comprehensive code linting
  - `github/codeql-action/upload-sarif@v3` - Upload SARIF results to GitHub Security
- **Additional Action Wrappers** - 3 typed wrappers for supply chain security and automation (#179)
  - `sigstore/cosign-installer@v3` - Cosign for signing and verifying container images
  - `actions-rs/cargo@v1` - Cargo command runner for Rust CI/CD
  - `EndBug/add-and-commit@v9` - Git automation for automatic commits
- **Issue Automation Example** - Example demonstrating issue and PR automation triggers (#180)
  - Issues trigger with auto-labeling patterns
  - IssueComment trigger for issue management
  - PullRequestReview trigger for review automation
- **Expression Contexts Guide** - Comprehensive documentation for GitHub expression contexts (#181)
  - Secrets, GitHub, Matrix, Env, Needs, Steps contexts
  - Condition builders and string functions
  - Security considerations with WAG017/WAG018 integration
- **Additional Action Wrappers** - 4 typed wrappers for security, code quality, and releases (#183)
  - `crazy-max/ghaction-import-gpg@v6` - Import GPG keys for commit/tag signing
  - `pre-commit/action@v3.0.1` - Run pre-commit hooks in CI
  - `SonarSource/sonarcloud-github-action@v3` - SonarCloud code quality analysis
  - `ncipollo/release-action@v1` - Enhanced GitHub release creation
- **Publishing Workflow Example** - Comprehensive example for release publishing (#184)
  - Docker image build and push to GHCR with multi-platform support
  - GitHub release creation with auto-generated notes
  - Multi-platform binary builds using matrix strategy
  - Conditional steps for stable vs prerelease handling
- **Security Patterns Guide** - Documentation for workflow security best practices (#185)
  - Permission scoping with WAG017 integration
  - Secrets handling with WAG003 patterns
  - Pull request target safety with WAG018 examples
  - Supply chain security (SLSA, cosign, dependency review)
  - Security scanning setup (CodeQL, Trivy, Scorecard)
- **Performance Benchmarks** - Benchmarks for key operations (#157)
  - `internal/discover` - AST discovery benchmarks
  - `internal/serialize` - YAML serialization benchmarks
  - `internal/importer` - Import and code generation benchmarks
- **Monorepo Workflow Example** - Comprehensive example for monorepo CI (#187)
  - Path filters on triggers to only run on relevant changes
  - Change detection using dorny/paths-filter action
  - Conditional job execution based on detected changes
  - Parallel service builds for API (Go), Web (Node.js), and Shared (Go)
  - Job outputs for passing change detection results to downstream jobs
- **Deployment Workflow Example** - Multi-environment deployment example (#188)
  - Staging and production deployment jobs with environment protection
  - workflow.Environment with name and URL for deployment tracking
  - Conditional deployment based on trigger (push vs workflow_dispatch)
  - Health check verification steps
- **Additional Action Wrappers** - 3 typed wrappers for Kubernetes and compliance (#189)
  - `helm/kind-action@v1` - Create Kubernetes clusters with KinD for testing
  - `fossas/fossa-action@v1` - FOSSA license compliance scanning
  - `azure/k8s-set-context@v4` - Set Kubernetes context for AKS/Arc deployments
- **Additional Example Workflows** - 3 new workflow examples (#148)
  - `examples/docker-workflow` - Docker build and push to GHCR with multi-stage CI
  - `examples/release-workflow` - Automated GitHub releases on version tags
  - `examples/matrix-workflow` - Multi-OS and Go version matrix testing
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
- Fixed trailing comma syntax in generated extraction programs for PR templates and Codeowners (#146)

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
