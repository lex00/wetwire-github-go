---
title: "Roadmap"
---

Feature matrix and implementation status for wetwire-github-go.

**Last Updated:** 2026-01-10

---

## CLI Commands

| Command | Status | Notes |
|---------|--------|-------|
| `build` | Complete | Generates .github/workflows/*.yml |
| `lint` | Complete | Includes `--fix` for WAG001 |
| `import` | Complete | Supports workflow, dependabot, issue-template, discussion-template, codeowners |
| `validate` | Complete | Uses actionlint |
| `list` | Complete | Lists workflows, jobs, triggers |
| `init` | Complete | Scaffolds new projects |
| `graph` | Complete | Mermaid and DOT output |
| `design` | Complete | AI-assisted workflow generation |
| `test` | Complete | Structural tests + 5 personas + 5-dimension scoring |
| `mcp` | Complete | MCP server via `design --mcp-server` for IDE integration |

---

## Lint Rules (WAG)

| Rule | Status | Description |
|------|--------|-------------|
| WAG001 | Complete | Use typed action wrappers instead of raw `uses:` strings |
| WAG002 | Complete | Use condition builders instead of raw expression strings |
| WAG003 | Complete | Use secrets context instead of hardcoded strings |
| WAG004 | Complete | Use matrix builder instead of inline maps |
| WAG005 | Complete | Extract inline structs to named variables |
| WAG006 | Complete | Detect duplicate workflow names |
| WAG007 | Complete | Flag oversized files (>N jobs) |
| WAG008 | Complete | Avoid hardcoded expression strings |
| WAG009 | Complete | Validate matrix dimension values (empty matrix detection) |
| WAG010 | Complete | Flag missing recommended action inputs |
| WAG011 | Complete | Detect potential unreachable jobs (undefined dependencies) |
| WAG012 | Complete | Warn about deprecated action versions |
| WAG013 | Complete | Detect pointer assignments (&Type{}) - use value semantics |
| WAG014 | Complete | Flag jobs without timeout-minutes setting |
| WAG015 | Complete | Suggest caching for setup-go/node/python actions |
| WAG016 | Complete | Validate concurrency settings |
| WAG017 | Complete | Suggest adding explicit permissions scope for security |
| WAG018 | Complete | Detect dangerous pull_request_target patterns with checkout actions |

---

## Configuration Types

| Type | Build | Import | Output Location |
|------|-------|--------|-----------------|
| GitHub Actions Workflows | Complete | Complete | `.github/workflows/*.yml` |
| Dependabot | Complete | Complete | `.github/dependabot.yml` |
| Issue Templates | Complete | Complete | `.github/ISSUE_TEMPLATE/*.yml` |
| Discussion Templates | Complete | Complete | `.github/DISCUSSION_TEMPLATE/*.yml` |
| PR Templates | Complete | Complete | `.github/PULL_REQUEST_TEMPLATE.md` |
| CODEOWNERS | Complete | Complete | `.github/CODEOWNERS` |

---

## Agent Integration

| Feature | Status | Notes |
|---------|--------|-------|
| Tool definitions | Complete | `init_package`, `write_file`, `run_lint`, `run_build`, `run_validate`, `read_file`, `ask_developer` |
| System prompt | Complete | GitHub Actions domain knowledge |
| GitHubAgent integration | Complete | AI-assisted workflow generation |
| Streaming support | Complete | `--stream` flag for token streaming |
| Lint enforcement | Complete | Automatic lint after file writes |
| ConsoleDeveloper | Complete | Interactive question/answer |
| Persona testing | Complete | 3 built-in personas (beginner, intermediate, expert) + custom |
| 5-dimension scoring | Complete | Completeness, Lint, Code, Output, Questions |

---

## Action Wrappers

Type-safe wrappers for popular GitHub Actions:

| Action | Package | Status |
|--------|---------|--------|
| actions/checkout | `actions/checkout` | Complete |
| actions/setup-go | `actions/setup_go` | Complete |
| actions/setup-node | `actions/setup_node` | Complete |
| actions/setup-python | `actions/setup_python` | Complete |
| actions/setup-java | `actions/setup_java` | Complete |
| actions/setup-dotnet | `actions/setup_dotnet` | Complete |
| actions/setup-ruby | `actions/setup_ruby` | Complete |
| dtolnay/rust-toolchain | `actions/setup_rust` | Complete |
| actions-rs/toolchain | `actions/actions_rs_toolchain` | Complete |
| actions/cache | `actions/cache` | Complete |
| actions/upload-artifact | `actions/upload_artifact` | Complete |
| actions/download-artifact | `actions/download_artifact` | Complete |
| softprops/action-gh-release | `actions/gh_release` | Complete |
| actions/github-script | `actions/github_script` | Complete |
| docker/build-push-action | `actions/docker_build_push` | Complete |
| docker/login-action | `actions/docker_login` | Complete |
| docker/setup-buildx-action | `actions/docker_setup_buildx` | Complete |
| docker/metadata-action | `actions/docker_metadata` | Complete |
| codecov/codecov-action | `actions/codecov` | Complete |
| aws-actions/configure-aws-credentials | `actions/aws_configure_credentials` | Complete |
| aws-actions/amazon-ecr-login | `actions/aws_ecr_login` | Complete |
| azure/login | `actions/azure_login` | Complete |
| azure/webapps-deploy | `actions/azure_webapps_deploy` | Complete |
| azure/docker-login | `actions/azure_docker_login` | Complete |
| google-github-actions/auth | `actions/gcp_auth` | Complete |
| google-github-actions/deploy-cloudrun | `actions/gcp_deploy_cloudrun` | Complete |
| google-github-actions/setup-gcloud | `actions/gcp_setup_gcloud` | Complete |
| golangci/golangci-lint-action | `actions/golangci_lint` | Complete |
| github/codeql-action/init | `actions/codeql_init` | Complete |
| github/codeql-action/analyze | `actions/codeql_analyze` | Complete |
| aquasecurity/trivy-action | `actions/trivy` | Complete |
| ossf/scorecard-action | `actions/scorecard` | Complete |
| actions/labeler | `actions/labeler` | Complete |
| actions/stale | `actions/stale` | Complete |
| peter-evans/create-pull-request | `actions/create_pull_request` | Complete |
| actions/dependency-review-action | `actions/dependency_review` | Complete |
| JamesIves/github-pages-deploy-action | `actions/gh_pages_deploy` | Complete |
| slackapi/slack-github-action | `actions/slack` | Complete |
| hashicorp/setup-terraform | `actions/setup_terraform` | Complete |
| peaceiris/actions-gh-pages | `actions/gh_pages_peaceiris` | Complete |
| actions/first-interaction | `actions/first_interaction` | Complete |
| actions/add-to-project | `actions/add_to_project` | Complete |
| actions/create-github-app-token | `actions/create_github_app_token` | Complete |
| actions/attest-build-provenance | `actions/attest_build_provenance` | Complete |
| mikepenz/action-junit-report | `actions/junit_report` | Complete |
| actions/create-release | `actions/create_release` | Complete |
| actions/upload-release-asset | `actions/upload_release_asset` | Complete |
| actions/configure-pages | `actions/configure_pages` | Complete |
| actions/deploy-pages | `actions/deploy_pages` | Complete |
| actions/upload-pages-artifact | `actions/upload_pages_artifact` | Complete |
| reviewdog/action-setup | `actions/reviewdog` | Complete |
| super-linter/super-linter | `actions/super_linter` | Complete |
| github/codeql-action/upload-sarif | `actions/upload_sarif` | Complete |
| sigstore/cosign-installer | `actions/cosign_installer` | Complete |
| actions-rs/cargo | `actions/cargo` | Complete |
| EndBug/add-and-commit | `actions/add_and_commit` | Complete |
| crazy-max/ghaction-import-gpg | `actions/import_gpg` | Complete |
| pre-commit/action | `actions/pre_commit` | Complete |
| SonarSource/sonarcloud-github-action | `actions/sonarcloud` | Complete |
| ncipollo/release-action | `actions/ncipollo_release` | Complete |
| helm/kind-action | `actions/kind` | Complete |
| fossas/fossa-action | `actions/fossa` | Complete |
| azure/k8s-set-context | `actions/k8s_set_context` | Complete |
| dawidd6/action-download-artifact | `actions/dawidd6_download_artifact` | Complete |
| anothrNick/github-tag-action | `actions/github_tag_action` | Complete |
| peaceiris/actions-hugo | `actions/hugo` | Complete |

---

## Documentation

| Document | Status | Path |
|----------|--------|------|
| README | Complete | `README.md` |
| CLAUDE.md | Complete | `CLAUDE.md` |
| CHANGELOG | Complete | `CHANGELOG.md` |
| Quick Start | Complete | `docs/QUICK_START.md` |
| CLI Reference | Complete | `docs/CLI.md` |
| Import Workflow | Complete | `docs/IMPORT_WORKFLOW.md` |
| FAQ | Complete | `docs/FAQ.md` |
| Roadmap | Complete | `docs/ROADMAP.md` |
| Developers | Complete | `docs/DEVELOPERS.md` |
| Internals | Complete | `docs/INTERNALS.md` |
| Examples | Complete | `docs/EXAMPLES.md` |
| Versioning | Complete | `docs/VERSIONING.md` |
| Adoption | Complete | `docs/ADOPTION.md` |

---

