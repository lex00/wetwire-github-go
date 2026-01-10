# Roadmap

Feature matrix and implementation status for wetwire-github-go.

**Last Updated:** 2026-01-09

---

## CLI Commands

| Command | Status | Notes |
|---------|--------|-------|
| `build` | ✅ Complete | Generates .github/workflows/*.yml |
| `lint` | ✅ Complete | Includes `--fix` for WAG001 |
| `import` | ✅ Complete | Supports workflow, dependabot, issue-template, discussion-template |
| `validate` | ✅ Complete | Uses actionlint |
| `list` | ✅ Complete | Lists workflows, jobs, triggers |
| `init` | ✅ Complete | Scaffolds new projects |
| `graph` | ✅ Complete | Mermaid and DOT output |
| `design` | ✅ Complete | AI-assisted workflow generation via wetwire-core-go |
| `test` | ✅ Complete | Structural tests + 5 personas + 5-dimension scoring |
| `mcp` | ❌ Missing | MCP server for IDE integration (#68) |

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

---

## Configuration Types

| Type | Build | Import | Output Location |
|------|-------|--------|-----------------|
| GitHub Actions Workflows | ✅ | ✅ | `.github/workflows/*.yml` |
| Dependabot | ✅ | ✅ | `.github/dependabot.yml` |
| Issue Templates | ✅ | ✅ | `.github/ISSUE_TEMPLATE/*.yml` |
| Discussion Templates | ✅ | ✅ | `.github/DISCUSSION_TEMPLATE/*.yml` |

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
| actions/cache | `actions/cache` | ✅ |
| actions/upload-artifact | `actions/upload_artifact` | ✅ |
| actions/download-artifact | `actions/download_artifact` | ✅ |

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

### Phase 4: Agent Integration 🚧
- [x] wetwire-core-go dependency
- [x] Tool definitions (7 tools)
- [x] Design command implementation (#66)
- [x] Test command personas (#67)
- [ ] MCP server support (#68)

### Phase 5: Polish
- [x] Lint --fix implementation (#65)
- [ ] Additional action wrappers
- [ ] Performance optimization

---

## References

- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md)
- [Feature Matrix](https://github.com/lex00/wetwire/blob/main/docs/FEATURE_MATRIX.md)
- [Domain Guide](https://github.com/lex00/wetwire/blob/main/docs/DOMAIN_GUIDE.md)
