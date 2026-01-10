# Roadmap

Feature matrix and implementation status for wetwire-github-go.

**Last Updated:** 2026-01-09

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
| actions/cache | `actions/cache` | ✅ |
| actions/upload-artifact | `actions/upload_artifact` | ✅ |
| actions/download-artifact | `actions/download_artifact` | ✅ |
| softprops/action-gh-release | `actions/gh_release` | ✅ |
| actions/github-script | `actions/github_script` | ✅ |
| docker/build-push-action | `actions/docker_build_push` | ✅ |
| docker/login-action | `actions/docker_login` | ✅ |
| docker/setup-buildx-action | `actions/docker_setup_buildx` | ✅ |
| codecov/codecov-action | `actions/codecov` | ✅ |
| aws-actions/configure-aws-credentials | `actions/aws_configure_credentials` | ✅ |
| aws-actions/amazon-ecr-login | `actions/aws_ecr_login` | ✅ |
| azure/login | `actions/azure_login` | ✅ |
| azure/webapps-deploy | `actions/azure_webapps_deploy` | ✅ |
| azure/docker-login | `actions/azure_docker_login` | ✅ |
| google-github-actions/auth | `actions/gcp_auth` | ✅ |
| google-github-actions/deploy-cloudrun | `actions/gcp_deploy_cloudrun` | ✅ |
| google-github-actions/setup-gcloud | `actions/gcp_setup_gcloud` | ✅ |

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
- [ ] Performance optimization

---

## References

- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md)
- [Feature Matrix](https://github.com/lex00/wetwire/blob/main/docs/FEATURE_MATRIX.md)
- [Domain Guide](https://github.com/lex00/wetwire/blob/main/docs/DOMAIN_GUIDE.md)
