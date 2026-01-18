# Changelog

All notable changes to wetwire-github-go are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

## [Unreleased]

### Added
- **GitHub Actions Workflow Scenario** - Added workflow_scenario example with CI/CD patterns (#279)
  - Scenario configuration with 3 persona prompts (beginner, intermediate, expert)
  - System prompt with GitHub Actions domain context and best practices
  - Expected outputs demonstrating build, test, and deploy jobs with matrix strategy
  - Validation rules requiring minimum 1 workflow and 3 jobs
  - Examples include staging and production deployment with environment gates
- **Imported Examples from Major Open Source Repositories** - Added 8 real-world workflow examples (#278)
  - Docker Compose: `ci.yml`, `merge.yml` (Apache-2.0)
  - Grafana: `backport-workflow.yml`, `codeql-analysis.yml` (AGPL-3.0)
  - HashiCorp Terraform: `build.yml`, `checks.yml` (BUSL-1.1)
  - Prometheus: `ci.yml`, `codeql-analysis.yml` (Apache-2.0)
  - Total: 8 workflows, 42 jobs, 201 steps imported
  - `examples/imported/README.md` with full attribution and license information
- **LintOpts.Fix and LintOpts.Disable Support** - Domain validator compliance for lint options (#276)
  - `opts.Disable` skips specified lint rule IDs (e.g., `["WAG001", "WAG002"]`)
  - `opts.Fix` indicates Fix mode was requested (auto-fix reserved for future implementation)
  - Added `LinterOptions` struct and `NewLinterWithOptions()` to internal/lint package
  - Domain validator now passes for both LintOpts checks

### Fixed
- **Agent Test Failures** - Fixed failing agent tests for completion requirements and lint state tracking (#272)
  - Fixed lint error tracking to handle exit code 1 (actual lint failure code) in addition to exit code 2
  - Updated tests to create actual lint violations instead of relying on command failures
  - Tests now properly verify that lintPassed is false when lint finds issues

### Changed
- **Split Large Test Files** - Improved maintainability by splitting test files over 800 lines (#271)
  - Split `internal/agent/agent_tools_test.go` (1564 lines) into 3 focused files: `agent_tools_file_test.go`, `agent_tools_exec_test.go`, `agent_tools_ask_test.go`
  - Split `internal/runner/runner_extract_errors_test.go` (1482 lines) into 2 files: `runner_extract_errors_dir_test.go`, `runner_extract_errors_exec_test.go`
  - Split `internal/template/builder_test.go` (961 lines) into 2 files: `builder_core_test.go`, `builder_reconstruction_test.go`
- **Split Large Test Files Into Focused Modules** - Major refactoring for maintainability (#266)
  - Split `internal/runner/runner_test.go` (4,312 lines) into 6 files: `runner_core_test.go` (770 lines), `runner_path_test.go` (223 lines), `runner_generate_test.go` (859 lines), `runner_extract_success_test.go` (290 lines), `runner_extract_errors_test.go` (1,482 lines), `runner_integration_test.go` (764 lines)
  - Split `internal/agent/agent_test.go` (3,870 lines) into 6 files: `agent_constructor_test.go` (375 lines), `agent_tools_test.go` (1,564 lines), `agent_state_test.go` (558 lines), `agent_enforcement_test.go` (826 lines), `agent_execute_test.go` (532 lines), `agent_helpers_test.go` (54 lines)
  - Split `internal/linter/linter_test.go` (2,292 lines) into 5 files: `linter_core_test.go` (273 lines), `rules_wag001_008_test.go` (410 lines), `rules_wag009_016_test.go` (651 lines), `rules_wag017_020_test.go` (639 lines), `linter_fix_test.go` (342 lines)
  - Split `internal/serialize/serialize_test.go` (2,051 lines) into 5 files: `serialize_workflow_test.go` (560 lines), `serialize_jobs_test.go` (746 lines), `serialize_expressions_test.go` (109 lines), `serialize_dependabot_test.go` (285 lines), `serialize_templates_test.go` (381 lines)
  - Split `codegen/fetch_test.go` (1,356 lines) into 4 files: `fetch_core_test.go` (212 lines), `fetch_schema_test.go` (122 lines), `fetch_action_test.go` (217 lines), `fetch_all_test.go` (846 lines)
- **Split Linter Rules** - Split `internal/linter/rules.go` (1,695 lines) into focused modules by rule category (#265)
  - `rules_helpers.go` (347 lines) - Shared types, helpers, and configuration maps
  - `rules_actions.go` (316 lines) - WAG001, WAG010, WAG012, WAG015
  - `rules_security.go` (252 lines) - WAG003, WAG017, WAG018, WAG020
  - `rules_structure.go` (295 lines) - WAG004-WAG008, WAG013
  - `rules_workflow.go` (239 lines) - WAG002, WAG009, WAG014, WAG016
  - `rules_dependencies.go` (214 lines) - WAG011, WAG019
- **Linter Package Rename** - Renamed `internal/linter` to `internal/lint` for consistency (#269)
- **Discover Core Migration** - Migrated discover to use `wetwire-core-go/ast` package (#268)
  - Uses coreast.ExtractTypeName and coreast.InferTypeFromValue instead of local implementations
  - Uses coreast.IsBuiltinIdent for Go builtin checks
  - Net reduction of 3 lines
- **Linter Core Migration** - Migrated linter to use `wetwire-core-go/lint` package (#267)
  - Added Severity type alias and constants from core lint package
  - Replaced string severity values with typed constants (SeverityError, SeverityWarning, SeverityInfo)
  - Maintains backward compatibility through string conversion at API boundaries
- **Core Dependency Update** - Upgraded wetwire-core-go to v1.16.0 for lint and ast packages (#267, #268)
- **MCP Migration** - Migrated to automated MCP server generation (#253)
  - Upgraded wetwire-core-go to v1.13.0 for domain.BuildMCPServer() support
  - Replaced manual MCP tool registration with automated server generation
  - MCP server now auto-generates tools from GitHubDomain implementation
  - Reduced mcp.go from 500+ lines to 48 lines
  - All MCP tools (init, build, lint, validate, list, graph) now auto-configured
- **Core Dependency Update** - Upgraded wetwire-core-go to v1.5.4 (#249)
  - Fixes provider cwd configuration in Kiro agent setup
  - Ensures MCP tools run in correct project directory

## [1.0.2] - 2026-01-11

### Added
- **Codecov Integration** - Upload coverage reports to Codecov (#228)
  - Added codecov-action@v4 to CI workflow for coverage reporting
  - Added codecov badge to README.md for visibility

## [1.0.1] - 2026-01-11

### Added
- **Release Workflow** - Automated multi-platform binary builds on version tags
  - Builds for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
  - Generates SHA256 checksums for all binaries
  - Uploads artifacts to GitHub releases automatically

## [1.0.0] - 2026-01-11

### Added
- **Provider Flag** - Add `--provider` flag to `design` and `test` commands (#220)
  - Supports `anthropic` (default) and `kiro` providers
  - Kiro provider integration pending wetwire-core-go provider abstraction
  - Added Provider field to agent.Config for future extensibility
  - Updated CLI documentation with provider examples

### Changed
- **Centralized Personas and Scoring** - Migrate to wetwire-core-go packages (#221)
  - Use `personas` package from wetwire-core-go for developer persona definitions
  - Use `scoring` package from wetwire-core-go for session scoring
  - Remove local `internal/personas/` and `internal/scoring/` packages
  - Updated wetwire-core-go dependency to v1.2.0
- **Core Command Framework** - Use wetwire-core-go command infrastructure (#222)
  - Use `cmd.NewRootCommand()` from wetwire-core-go for root command
  - Consistent CLI structure with other wetwire domain packages
  - Domain-specific commands (build, lint, etc.) preserved as-is

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
  - `internal/lint`: 83.9% → 91.3% (7+ percentage point improvement) (#143, #149)
  - `codegen`: 85.2% → 95.7% (10+ percentage point improvement) (#143, #158)
  - `internal/runner`: 89.9% → 90.8% (additional edge case tests) (#158)
- **Circular Dependency Detection (WAG019)** - New lint rule for job dependency cycles (#204)
  - WAG019: Detect circular dependencies in job needs (e.g., A -> B -> C -> A)
  - Uses DFS with recursion stack for efficient cycle detection
  - Reports all jobs involved in the cycle with file locations
- **Comprehensive Secret Detection (WAG020)** - Detect 25+ hardcoded secret patterns (#205)
  - AWS access keys, GitHub tokens (all variants), Stripe keys (live/test)
  - Private keys (RSA, EC, DSA, OpenSSH, PGP), Slack tokens, Google API keys
  - Twilio, SendGrid, Mailgun, NPM, PyPI tokens, JWT tokens
  - DigitalOcean, Heroku, Azure credentials
- **CI Coverage Reporting** - Added coverage reporting to CI workflow (#212)
  - Generates coverage profile with atomic mode
  - Reports total and per-package coverage to GitHub Actions step summary
  - Warns when coverage drops below 70% threshold
- **LINT_RULES.md Documentation** - Complete documentation of all 20 lint rules (#211)
  - Rule index with severity and auto-fix status
  - Detailed descriptions with bad/good code examples
  - Usage instructions for lint command and JSON output
- **Examples and Attribution Documentation** - Proper attribution and example documentation (#210)
  - Added examples/README.md with example categorization and usage instructions
  - Added formal attribution to testdata/reference/ for imported starter-workflows (MIT)
  - Clear distinction between hand-written examples and imported workflows
- **Diff and Watch Commands** - New CLI commands for iterative development (#209)
  - `diff`: Semantically compare workflow configurations (Go packages or YAML files)
  - `watch`: Auto-rebuild workflows on source file changes with fsnotify
  - Support for text, JSON, and markdown output formats
  - Configurable debounce duration for watch mode
  - Lint-only mode for watch command
- **Kiro CLI Integration** - AI-assisted workflow design with Kiro (#208)
  - `internal/kiro/` package for auto-installation of agent config
  - Embedded wetwire-runner.json with GitHub Actions workflow prompts
  - Project-level MCP configuration for Kiro integration
  - Automatic binary detection with go run fallback
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
- **Container Services Example** - Database integration testing workflow (#191)
  - Job.Container for running steps in a Docker container (golang:1.24)
  - Job.Services with PostgreSQL and Redis service containers
  - Health check configuration via Options field
  - Environment variables for database connection strings
- **Artifact Pipeline Example** - Multi-stage build-test-release workflow (#192)
  - Build → Test → Release pipeline with artifact passing
  - upload_artifact/download_artifact action wrappers for artifact management
  - Conditional release job with workflow.StartsWith() for tag detection
  - Cross-platform binary builds (Linux, macOS, Windows)
- **Additional Action Wrappers** - 2 typed wrappers for cross-workflow artifacts and tagging (#193)
  - `dawidd6/action-download-artifact@v6` - Download artifacts from other workflow runs
  - `anothrNick/github-tag-action@v1` - Automated semantic version tagging
- **Advanced Trigger Pattern Examples** - 2 examples for workflow_run and repository_dispatch triggers (#195)
  - `examples/workflow-run-example` - Trigger workflows on completion of other workflows
  - `examples/repository-dispatch-example` - API-triggered workflows with event type filtering
  - Demonstrates conditional job execution based on event types (deploy, deploy-staging, deploy-production)
- **CLI Integration Tests** - Comprehensive test suite for CLI commands (#196)
  - 34+ tests for build, lint, and list commands
  - Tests cover valid workflows, invalid syntax, missing directories, and flag combinations
  - Uses exec.Command to test the actual compiled binary behavior
- **Hugo Action Wrapper** - Typed wrapper for static site generation (#197)
  - `peaceiris/actions-hugo@v3` - Setup Hugo static site generator with version control
  - 100% test coverage with comprehensive input handling
- **Additional Action Wrappers** - 4 typed wrappers for Kubernetes and IaC tools (#199)
  - `helm/chart-releaser-action@v1` - Release Helm charts to GitHub Pages
  - `azure/setup-helm@v4` - Install Helm CLI on runners
  - `stefanprodan/kustomize-action@master` - Kustomize build and apply for GitOps
  - `pulumi/actions@v6` - Infrastructure deployment with Pulumi
  - All wrappers have 100% test coverage
- **CLI Integration Tests** - Extended test coverage for CLI commands (#200)
  - `validate_test.go` - 8 tests for workflow validation command
  - `test_test.go` - 9 tests for project testing command
  - `version_test.go` - 3 tests for version command
  - Extended init, graph, import tests (82+ CLI tests total)
  - All tests use exec.Command pattern for integration testing
- **Agent Test Coverage** - Improved coverage for internal/agent package (#173)
  - `agent_coverage_test.go` - 171 tests (58.8% coverage)
  - `agent_streaming_test.go` - Streaming integration tests
  - Tests cover tool routing, state management, completion gate logic
  - Note: Full 80% blocked by methods requiring real API calls
- **Environment Workflow Examples** - 2 examples for deployment environments (#201)
  - `examples/approval-gates-workflow` - Environment approval gates pattern
  - `examples/environment-promotion-workflow` - Dev/staging/prod promotion
  - Demonstrates workflow.Environment with name and URL
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

- **README Badges** - Added missing badges per WETWIRE_SPEC.md Section 12.4 (#223)
  - Go Reference badge (pkg.go.dev)
  - Go Report Card badge (goreportcard.com)
  - Badge order: CI, Go Reference, Go Report Card, License

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

[Unreleased]: https://github.com/lex00/wetwire-github-go/compare/v1.0.2...HEAD
[1.0.2]: https://github.com/lex00/wetwire-github-go/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/lex00/wetwire-github-go/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/lex00/wetwire-github-go/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/lex00/wetwire-github-go/releases/tag/v0.1.0
