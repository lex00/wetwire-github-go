# Roadmap

Feature matrix and implementation status for wetwire-github-go.

**Last Updated:** 2026-01-10

---

## CLI Commands

| Command | Status | Notes |
|---------|--------|-------|
| `build` | ✅ Complete | Generates .github/workflows/*.yml |
| `lint` | ✅ Complete | Includes `--fix` for WAG001 |
| `import` | ✅ Complete | Supports workflow, dependabot, issue-template, discussion-template, codeowners |
| `validate` | ✅ Complete | Uses actionlint |
| `list` | ✅ Complete | Lists workflows, jobs, triggers |
| `init` | ✅ Complete | Scaffolds new projects |
| `graph` | ✅ Complete | Mermaid and DOT output |
| `design` | ✅ Complete | AI-assisted workflow generation via wetwire-core-go |
| `test` | ✅ Complete | Structural tests + 5 personas + 5-dimension scoring |
| `mcp` | ✅ Complete | MCP server via `design --mcp-server` for IDE integration |

---

## Lint Rules (WAG)

| Rule | Status | Description |
|------|--------|-------------|
| WAG001 | ✅ | Use typed action wrappers instead of raw `uses:` strings |
| WAG002 | ✅ | Use condition builders instead of raw expression strings |
| WAG003 | ✅ | Use secrets context instead of hardcoded strings |
| WAG004 | ✅ | Use matrix builder instead of inline maps |
| WAG005 | ✅ | Extract inline structs to named variables |
| WAG006 | ✅ | Detect duplicate workflow names |
| WAG007 | ✅ | Flag oversized files (>N jobs) |
| WAG008 | ✅ | Avoid hardcoded expression strings |
| WAG009 | ✅ | Validate matrix dimension values (empty matrix detection) |
| WAG010 | ✅ | Flag missing recommended action inputs |
| WAG011 | ✅ | Detect potential unreachable jobs (undefined dependencies) |
| WAG012 | ✅ | Warn about deprecated action versions |
| WAG013 | ✅ | Detect pointer assignments (&Type{}) - use value semantics |
| WAG014 | ✅ | Flag jobs without timeout-minutes setting |
| WAG015 | ✅ | Suggest caching for setup-go/node/python actions |
| WAG016 | ✅ | Validate concurrency settings |
| WAG017 | ✅ | Suggest adding explicit permissions scope for security |
| WAG018 | ✅ | Detect dangerous pull_request_target patterns with checkout actions |

---

## Configuration Types

| Type | Build | Import | Output Location |
|------|-------|--------|-----------------|
| GitHub Actions Workflows | ✅ | ✅ | `.github/workflows/*.yml` |
| Dependabot | ✅ | ✅ | `.github/dependabot.yml` |
| Issue Templates | ✅ | ✅ | `.github/ISSUE_TEMPLATE/*.yml` |
| Discussion Templates | ✅ | ✅ | `.github/DISCUSSION_TEMPLATE/*.yml` |
| PR Templates | ✅ | ✅ | `.github/PULL_REQUEST_TEMPLATE.md` |
| CODEOWNERS | ✅ | ✅ | `.github/CODEOWNERS` |

---

## Agent Integration (wetwire-core-go)

| Feature | Status | Notes |
|---------|--------|-------|
| Tool definitions | ✅ | `init_package`, `write_file`, `run_lint`, `run_build`, `run_validate`, `read_file`, `ask_developer` |
| System prompt | ✅ | GitHub Actions domain knowledge |
| GitHubAgent integration | ✅ | AI-assisted workflow generation |
| Streaming support | ✅ | `--stream` flag for token streaming |
| Lint enforcement | ✅ | Automatic lint after file writes |
| ConsoleDeveloper | ✅ | Interactive question/answer |
| Persona testing | ✅ | 5 standard personas (beginner, intermediate, expert, terse, verbose) |
| 5-dimension scoring | ✅ | Completeness, Lint, Code, Output, Questions |

---

## Action Wrappers

Type-safe wrappers for popular GitHub Actions:

| Action | Package | Status |
|--------|---------|--------|
| actions/checkout | `actions/checkout` | ✅ |
| actions/setup-go | `actions/setup_go` | ✅ |
| actions/setup-node | `actions/setup_node` | ✅ |
| actions/setup-python | `actions/setup_python` | ✅ |
| actions/setup-java | `actions/setup_java` | ✅ |
| actions/setup-dotnet | `actions/setup_dotnet` | ✅ |
| actions/setup-ruby | `actions/setup_ruby` | ✅ |
| dtolnay/rust-toolchain | `actions/setup_rust` | ✅ |
| actions-rs/toolchain | `actions/actions_rs_toolchain` | ✅ |
| actions/cache | `actions/cache` | ✅ |
| actions/upload-artifact | `actions/upload_artifact` | ✅ |
| actions/download-artifact | `actions/download_artifact` | ✅ |
| softprops/action-gh-release | `actions/gh_release` | ✅ |
| actions/github-script | `actions/github_script` | ✅ |
| docker/build-push-action | `actions/docker_build_push` | ✅ |
| docker/login-action | `actions/docker_login` | ✅ |
| docker/setup-buildx-action | `actions/docker_setup_buildx` | ✅ |
| docker/metadata-action | `actions/docker_metadata` | ✅ |
| codecov/codecov-action | `actions/codecov` | ✅ |
| aws-actions/configure-aws-credentials | `actions/aws_configure_credentials` | ✅ |
| aws-actions/amazon-ecr-login | `actions/aws_ecr_login` | ✅ |
| azure/login | `actions/azure_login` | ✅ |
| azure/webapps-deploy | `actions/azure_webapps_deploy` | ✅ |
| azure/docker-login | `actions/azure_docker_login` | ✅ |
| google-github-actions/auth | `actions/gcp_auth` | ✅ |
| google-github-actions/deploy-cloudrun | `actions/gcp_deploy_cloudrun` | ✅ |
| google-github-actions/setup-gcloud | `actions/gcp_setup_gcloud` | ✅ |
| golangci/golangci-lint-action | `actions/golangci_lint` | ✅ |
| github/codeql-action/init | `actions/codeql_init` | ✅ |
| github/codeql-action/analyze | `actions/codeql_analyze` | ✅ |
| aquasecurity/trivy-action | `actions/trivy` | ✅ |
| ossf/scorecard-action | `actions/scorecard` | ✅ |
| actions/labeler | `actions/labeler` | ✅ |
| actions/stale | `actions/stale` | ✅ |
| peter-evans/create-pull-request | `actions/create_pull_request` | ✅ |
| actions/dependency-review-action | `actions/dependency_review` | ✅ |
| JamesIves/github-pages-deploy-action | `actions/gh_pages_deploy` | ✅ |
| slackapi/slack-github-action | `actions/slack` | ✅ |
| hashicorp/setup-terraform | `actions/setup_terraform` | ✅ |
| peaceiris/actions-gh-pages | `actions/gh_pages_peaceiris` | ✅ |
| actions/first-interaction | `actions/first_interaction` | ✅ |
| actions/add-to-project | `actions/add_to_project` | ✅ |
| actions/create-github-app-token | `actions/create_github_app_token` | ✅ |
| actions/attest-build-provenance | `actions/attest_build_provenance` | ✅ |
| mikepenz/action-junit-report | `actions/junit_report` | ✅ |
| actions/create-release | `actions/create_release` | ✅ |
| actions/upload-release-asset | `actions/upload_release_asset` | ✅ |
| actions/configure-pages | `actions/configure_pages` | ✅ |
| actions/deploy-pages | `actions/deploy_pages` | ✅ |
| actions/upload-pages-artifact | `actions/upload_pages_artifact` | ✅ |
| reviewdog/action-setup | `actions/reviewdog` | ✅ |
| super-linter/super-linter | `actions/super_linter` | ✅ |
| github/codeql-action/upload-sarif | `actions/upload_sarif` | ✅ |
| sigstore/cosign-installer | `actions/cosign_installer` | ✅ |
| actions-rs/cargo | `actions/cargo` | ✅ |
| EndBug/add-and-commit | `actions/add_and_commit` | ✅ |
| crazy-max/ghaction-import-gpg | `actions/import_gpg` | ✅ |
| pre-commit/action | `actions/pre_commit` | ✅ |
| SonarSource/sonarcloud-github-action | `actions/sonarcloud` | ✅ |
| ncipollo/release-action | `actions/ncipollo_release` | ✅ |
| helm/kind-action | `actions/kind` | ✅ |
| fossas/fossa-action | `actions/fossa` | ✅ |
| azure/k8s-set-context | `actions/k8s_set_context` | ✅ |
| dawidd6/action-download-artifact | `actions/dawidd6_download_artifact` | ✅ |
| anothrNick/github-tag-action | `actions/github_tag_action` | ✅ |
| peaceiris/actions-hugo | `actions/hugo` | ✅ |

---

## Documentation

| Document | Status | Path |
|----------|--------|------|
| README | ✅ | `README.md` |
| CLAUDE.md | ✅ | `CLAUDE.md` |
| CHANGELOG | ✅ | `CHANGELOG.md` |
| Quick Start | ✅ | `docs/QUICK_START.md` |
| CLI Reference | ✅ | `docs/CLI.md` |
| Import Workflow | ✅ | `docs/IMPORT_WORKFLOW.md` |
| FAQ | ✅ | `docs/FAQ.md` |
| Roadmap | ✅ | `docs/ROADMAP.md` |
| Developers | ✅ | `docs/DEVELOPERS.md` |
| Internals | ✅ | `docs/INTERNALS.md` |
| Examples | ✅ | `docs/EXAMPLES.md` |
| Versioning | ✅ | `docs/VERSIONING.md` |
| Adoption | ✅ | `docs/ADOPTION.md` |

---

## Implementation Phases

### Phase 1: Core Types ✅
- Workflow, Job, Step types
- Trigger configurations
- Expression contexts (secrets, github, matrix, env, needs)

### Phase 2: CLI Commands ✅
- build, lint, import, validate, list, init, graph, test

### Phase 3: Extended Types ✅
- Dependabot configuration
- Issue templates
- Discussion templates

### Phase 4: Agent Integration ✅
- [x] wetwire-core-go dependency
- [x] Tool definitions (7 tools)
- [x] Design command implementation (#66)
- [x] Test command personas (#67)
- [x] MCP server support (#68)

### Phase 5: Polish ✅
- [x] Lint --fix implementation (#65)
- [x] Additional action wrappers (setup_java, setup_dotnet, setup_ruby, setup_rust, gh_release)
- [x] Additional lint rules (WAG009-WAG012)
- [x] CODEOWNERS support (#84)
- [x] PR Templates support (#83)
- [x] More action wrappers (github-script #91, docker actions #92, codecov #93)
- [x] Documentation expansion (DEVELOPERS.md #94, INTERNALS.md #95, EXAMPLES.md #96)

### Phase 6: Cloud Platform Wrappers ✅
- [x] PR Template import support (#103)
- [x] CODEOWNERS import support (#104)
- [x] AWS action wrappers (configure-aws-credentials, ecr-login) (#105)
- [x] Azure action wrappers (login, webapp-deploy, docker-login) (#106)
- [x] GCP action wrappers (auth, deploy-cloudrun, setup-gcloud) (#107)
- [x] Remove .ToStep() refactor - pure struct-literal pattern (#110)

### Phase 7: Future Enhancements
- [x] Type-safe intrinsics for GitHub expression contexts (already implemented in workflow/expressions.go)
- [x] Dedicated graph package for visualization (already implemented - DOT, Mermaid, JSON via CLI)
- [x] Additional documentation (VERSIONING.md, ADOPTION.md) (#117)
- [x] Additional lint rules (WAG013-WAG016) (#149)
- [x] Additional example workflows (docker, release, matrix) (#148)
- [x] Reference example testing with round-trip tests (#151)
- [x] Additional action wrappers (gh_pages_deploy, slack, setup_terraform) (#152)
- [x] Additional action wrappers (gh_pages_peaceiris, first_interaction, add_to_project, create_github_app_token) (#154)
- [x] Agent test coverage improvements (24.7% → 55.2%) (#155)
- [x] Performance benchmarks (discover, serialize, importer, runner, linter) (#157)
- [x] Codegen test coverage improvements (89.5% → 95.7%) (#158)
- [x] Reusable workflow example with workflow_call trigger (#176)
- [x] Advanced trigger pattern examples (workflow_run, repository_dispatch) (#195)
- [x] CLI integration tests for build, lint, list commands (#196)
- [x] Hugo action wrapper (#197)

---

## References

- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md)
- [Feature Matrix](https://github.com/lex00/wetwire/blob/main/docs/FEATURE_MATRIX.md)
- [Domain Guide](https://github.com/lex00/wetwire/blob/main/docs/DOMAIN_GUIDE.md)
