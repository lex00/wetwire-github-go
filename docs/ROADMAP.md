# wetwire-github-go Implementation Roadmap

## The "No Parens" Pattern

All declarations use struct literals — no function calls or registration:

```go
// Workflows, jobs, steps as variables
var CI = workflow.Workflow{Name: "CI", On: workflow.Triggers{...}}
var Build = workflow.Job{Name: "build", RunsOn: "ubuntu-latest", Steps: [...]}

// Cross-references via field access
var Deploy = workflow.Job{Needs: []any{Build, Test}}

// Type-safe action wrappers
checkout.Checkout{FetchDepth: 0}.ToStep()

// Expression contexts
workflow.Secrets.Get("TOKEN")
workflow.Matrix.Get("os")
workflow.GitHub.Ref
```

AST-based discovery — no registration needed.

---

## Scope

wetwire-github-go generates typed Go declarations for three GitHub YAML configuration types:

| Config Type | Output Location | Schema Source |
|-------------|-----------------|---------------|
| **GitHub Actions** | `.github/workflows/*.yml` | `json.schemastore.org/github-workflow.json` |
| **Dependabot** | `.github/dependabot.yml` | `json.schemastore.org/dependabot-2.0.json` |
| **Issue/Discussion Templates** | `.github/ISSUE_TEMPLATE/*.yml` | `json.schemastore.org/github-issue-forms.json` |

---

## Development Infrastructure (Phase 0)

Project setup and CI/CD infrastructure (mirroring wetwire-aws-go):

| Feature | File/Location | Status |
|---------|---------------|--------|
| **Repository Setup** | | |
| ├─ go.mod with dependencies | `go.mod` | [ ] |
| ├─ .gitignore (Go project) | `.gitignore` | [ ] |
| ├─ specs/.gitkeep | `specs/.gitkeep` | [ ] |
| └─ specs in .gitignore | `.gitignore` | [ ] |
| **GitHub Actions CI** | | |
| ├─ ci.yml (build/test on push/PR) | `.github/workflows/ci.yml` | [ ] |
| └─ codebot.yml (Claude Code integration) | `.github/workflows/codebot.yml` | [ ] |
| **Development Scripts** | | |
| ├─ ci.sh (local CI runner) | `scripts/ci.sh` | [ ] |
| └─ import_samples.sh (round-trip testing) | `scripts/import_samples.sh` | [ ] |
| **Documentation** | | |
| ├─ README.md (full) | `README.md` | [ ] |
| ├─ QUICK_START.md | `docs/QUICK_START.md` | [ ] |
| ├─ CLI.md | `docs/CLI.md` | [ ] |
| └─ IMPORT_WORKFLOW.md | `docs/IMPORT_WORKFLOW.md` | [ ] |
| **Build Artifacts** | | |
| └─ codegen-bin (pre-built binary) | `codegen-bin` | [ ] |

---

## Feature Matrix

| Feature | Phase | Status | Dependencies |
|---------|-------|--------|--------------|
| **Core Types** | | | |
| ├─ Workflow struct | 1A | [ ] | — |
| ├─ Job struct | 1A | [ ] | — |
| ├─ Step struct | 1A | [ ] | — |
| ├─ Matrix struct | 1A | [ ] | — |
| ├─ Triggers (all 30+ types) | 1A | [ ] | — |
| ├─ Conditions interface | 1A | [ ] | — |
| ├─ Expression contexts | 1A | [ ] | — |
| ├─ Helper types (Env, With, List, Any, Strings) | 1A | [ ] | — |
| ├─ contracts.go (interfaces) | 1A | [ ] | — |
| ├─ Result types (Build/Lint/Validate/List) | 1A | [ ] | — |
| └─ DiscoveredWorkflow/Job structs | 1A | [ ] | — |
| **Serialization** | | | |
| ├─ Workflow → YAML | 1B | [ ] | Core Types |
| ├─ Job → YAML | 1B | [ ] | Core Types |
| ├─ Step → YAML | 1B | [ ] | Core Types |
| ├─ Matrix → YAML | 1B | [ ] | Core Types |
| ├─ Triggers → YAML | 1B | [ ] | Core Types |
| ├─ Conditions → expression strings | 1B | [ ] | Core Types |
| ├─ Struct → map serialization | 1B | [ ] | Core Types |
| ├─ PascalCase/snake_case conversion | 1B | [ ] | — |
| └─ Zero value omission | 1B | [ ] | — |
| **CLI Framework** | | | |
| ├─ main.go + cobra setup | 1C | [ ] | — |
| ├─ `init` command | 1C | [ ] | — |
| ├─ `build` command (stub) | 1C | [ ] | — |
| ├─ `validate` command (stub) | 1C | [ ] | — |
| ├─ `list` command (stub) | 1C | [ ] | — |
| ├─ `lint` command (stub) | 1C | [ ] | — |
| ├─ `import` command (stub) | 1C | [ ] | — |
| ├─ `design` command (stub) | 1C | [ ] | — |
| ├─ `test` command (stub) | 1C | [ ] | — |
| ├─ `version` command | 1C | [ ] | — |
| └─ Exit code semantics (0/1/2) | 1C | [ ] | — |
| **Schema Fetching** | | | |
| ├─ HTTP fetcher with retry | 1D | [ ] | — |
| ├─ Workflow schema fetch | 1D | [ ] | — |
| ├─ Action.yml fetch (checkout) | 1D | [ ] | — |
| ├─ Action.yml fetch (setup-python) | 1D | [ ] | — |
| ├─ Action.yml fetch (setup-node) | 1D | [ ] | — |
| ├─ Action.yml fetch (setup-go) | 1D | [ ] | — |
| ├─ Action.yml fetch (cache) | 1D | [ ] | — |
| ├─ Action.yml fetch (upload-artifact) | 1D | [ ] | — |
| ├─ Action.yml fetch (download-artifact) | 1D | [ ] | — |
| └─ specs/manifest.json | 1D | [ ] | — |
| **Schema Parsing** | | | |
| ├─ Workflow schema parser | 2A | [ ] | Schema Fetching |
| ├─ Action.yml parser | 2A | [ ] | Schema Fetching |
| └─ Intermediate representation | 2A | [ ] | Schema Fetching |
| **Action Codegen** | | | |
| ├─ Generator templates | 2B | [ ] | Schema Parsing |
| ├─ Code formatting (go/format) | 2B | [ ] | — |
| ├─ actions/checkout wrapper | 2B | [ ] | Schema Parsing |
| ├─ actions/setup_python wrapper | 2B | [ ] | Schema Parsing |
| ├─ actions/setup_node wrapper | 2B | [ ] | Schema Parsing |
| ├─ actions/setup_go wrapper | 2B | [ ] | Schema Parsing |
| ├─ actions/cache wrapper | 2B | [ ] | Schema Parsing |
| ├─ actions/upload_artifact wrapper | 2B | [ ] | Schema Parsing |
| └─ actions/download_artifact wrapper | 2B | [ ] | Schema Parsing |
| **AST Discovery** | | | |
| ├─ Package scanning (go/parser) | 2C | [ ] | Core Types |
| ├─ Workflow variable detection | 2C | [ ] | Core Types |
| ├─ Job variable detection | 2C | [ ] | Core Types |
| ├─ Dependency graph building | 2C | [ ] | Core Types |
| ├─ Reference validation | 2C | [ ] | Core Types |
| ├─ Recursive directory traversal | 2C | [ ] | — |
| └─ Vendor/hidden directory exclusion | 2C | [ ] | — |
| **Runner (Value Extraction)** | | | |
| ├─ Temp Go program generation | 2G | [ ] | AST Discovery |
| ├─ go.mod parsing | 2G | [ ] | — |
| ├─ Replace directive handling | 2G | [ ] | — |
| ├─ Go binary discovery | 2G | [ ] | — |
| └─ JSON value extraction | 2G | [ ] | — |
| **Template Builder** | | | |
| ├─ Topological sort (Kahn's algorithm) | 2H | [ ] | AST Discovery |
| ├─ Cycle detection | 2H | [ ] | AST Discovery |
| └─ Dependency ordering | 2H | [ ] | AST Discovery |
| **Build Command (full)** | | | |
| ├─ Discovery integration | 3A | [ ] | AST Discovery |
| ├─ Runner integration | 3A | [ ] | Runner |
| ├─ Template builder integration | 3A | [ ] | Template Builder |
| ├─ Multi-workflow support | 3A | [ ] | Serialization |
| ├─ Output to `.github/workflows/` | 3A | [ ] | Serialization |
| ├─ --format json/yaml | 3A | [ ] | — |
| └─ --output flag | 3A | [ ] | — |
| **Validation (actionlint)** | | | |
| ├─ actionlint Go library integration | 2D | [ ] | — |
| ├─ ValidationResult types | 2D | [ ] | — |
| ├─ Multiple validator pipeline | 2D | [ ] | — |
| └─ `validate` command (full) | 3B | [ ] | actionlint integration |
| **Linting Rules** | | | |
| ├─ Linter framework | 2E | [ ] | — |
| ├─ Rule interface (ID, Description, Check) | 2E | [ ] | — |
| ├─ WAG001: typed action wrappers | 2E | [ ] | — |
| ├─ WAG002: condition builders | 2E | [ ] | — |
| ├─ WAG003: secrets context | 2E | [ ] | — |
| ├─ WAG004: matrix builder | 2E | [ ] | — |
| ├─ WAG005: inline structs → named vars | 2E | [ ] | — |
| ├─ WAG006: duplicate workflow names | 2E | [ ] | — |
| ├─ WAG007: file too large (>N jobs) | 2E | [ ] | — |
| ├─ WAG008: hardcoded expression strings | 2E | [ ] | — |
| ├─ Recursive package scanning | 2E | [ ] | — |
| ├─ --fix flag support | 2E | [ ] | — |
| └─ `lint` command (full) | 3C | [ ] | Linter framework |
| **Import (YAML → Go)** | | | |
| ├─ YAML parser | 2F | [ ] | — |
| ├─ IR (intermediate representation) | 2F | [ ] | — |
| ├─ IRWorkflow, IRJob, IRStep structs | 2F | [ ] | — |
| ├─ Reference graph tracking | 2F | [ ] | — |
| ├─ Go code generator | 3D | [ ] | YAML parser, Core Types |
| ├─ Field name transformation | 3D | [ ] | — |
| ├─ Reserved name handling | 3D | [ ] | — |
| ├─ Scaffold: go.mod | 3D | [ ] | — |
| ├─ Scaffold: cmd/main.go | 3D | [ ] | — |
| ├─ Scaffold: README.md | 3D | [ ] | — |
| ├─ Scaffold: CLAUDE.md | 3D | [ ] | — |
| ├─ Scaffold: .gitignore | 3D | [ ] | — |
| ├─ --single-file flag | 3D | [ ] | — |
| ├─ --no-scaffold flag | 3D | [ ] | — |
| └─ `import` command (full) | 3D | [ ] | Go code generator |
| **List Command** | | | |
| ├─ Table output format | 3E | [ ] | AST Discovery |
| ├─ --format json | 3E | [ ] | — |
| └─ `list` command (full) | 3E | [ ] | AST Discovery |
| **Design Command (AI-assisted)** | | | |
| ├─ wetwire-core-go orchestrator | 3F | [ ] | All CLI commands |
| ├─ Interactive session | 3F | [ ] | — |
| ├─ --stream flag | 3F | [ ] | — |
| ├─ --max-lint-cycles flag | 3F | [ ] | — |
| └─ `design` command (full) | 3F | [ ] | wetwire-core-go |
| **Test Command (Persona-based)** | | | |
| ├─ Persona selection | 3G | [ ] | wetwire-core-go |
| ├─ Scenario configuration | 3G | [ ] | — |
| ├─ Result writing | 3G | [ ] | — |
| ├─ Session tracking | 3G | [ ] | — |
| └─ `test` command (full) | 3G | [ ] | wetwire-core-go |
| **Examples & Testing** | | | |
| ├─ Fetch starter-workflows | 4A | [ ] | Schema Fetching |
| ├─ Import/rebuild cycle tests | 4A | [ ] | Import, Build |
| ├─ Round-trip validation | 4A | [ ] | Import, Build, Validate |
| └─ Success rate tracking | 4A | [ ] | — |
| **wetwire-core-go Integration** | | | |
| ├─ RunnerAgent tool definitions | 4B | [ ] | All CLI commands |
| ├─ Tool handlers implementation | 4B | [ ] | — |
| ├─ Stream handler support | 4B | [ ] | — |
| ├─ Session result writing | 4B | [ ] | — |
| ├─ Scoring integration | 4B | [ ] | — |
| └─ Agent testing with personas | 4B | [ ] | — |
| **Dependabot Types** | | | |
| ├─ Dependabot struct | 5A | [ ] | — |
| ├─ Update struct | 5A | [ ] | — |
| ├─ Schedule struct | 5A | [ ] | — |
| ├─ PackageEcosystem enum | 5A | [ ] | — |
| ├─ Registries struct | 5A | [ ] | — |
| ├─ Groups struct | 5A | [ ] | — |
| └─ contracts (Dependabot) | 5A | [ ] | — |
| **Dependabot Schema** | | | |
| ├─ Fetch dependabot-2.0.json | 5A | [ ] | — |
| ├─ Parse dependabot schema | 5A | [ ] | — |
| └─ Dependabot → YAML serialization | 5A | [ ] | — |
| **Dependabot CLI** | | | |
| ├─ `build --type dependabot` | 5A | [ ] | Dependabot Types |
| ├─ `import --type dependabot` | 5A | [ ] | Dependabot Types |
| ├─ `validate --type dependabot` | 5A | [ ] | — |
| └─ AST discovery for Dependabot | 5A | [ ] | Dependabot Types |
| **Issue Template Types** | | | |
| ├─ IssueTemplate struct | 5B | [ ] | — |
| ├─ FormBody struct | 5B | [ ] | — |
| ├─ FormElement interface | 5B | [ ] | — |
| ├─ Input element | 5B | [ ] | — |
| ├─ Textarea element | 5B | [ ] | — |
| ├─ Dropdown element | 5B | [ ] | — |
| ├─ Checkboxes element | 5B | [ ] | — |
| ├─ Markdown element | 5B | [ ] | — |
| └─ contracts (IssueTemplate) | 5B | [ ] | — |
| **Issue Template Schema** | | | |
| ├─ Fetch github-issue-forms.json | 5B | [ ] | — |
| ├─ Parse issue forms schema | 5B | [ ] | — |
| └─ IssueTemplate → YAML serialization | 5B | [ ] | — |
| **Issue Template CLI** | | | |
| ├─ `build --type issue-template` | 5B | [ ] | IssueTemplate Types |
| ├─ `import --type issue-template` | 5B | [ ] | IssueTemplate Types |
| ├─ `validate --type issue-template` | 5B | [ ] | — |
| └─ AST discovery for IssueTemplate | 5B | [ ] | IssueTemplate Types |
| **Discussion Template Types** | | | |
| ├─ DiscussionTemplate struct | 5C | [ ] | FormBody (from 5B) |
| └─ Discussion category handling | 5C | [ ] | — |
| **Discussion Template CLI** | | | |
| ├─ `build --type discussion-template` | 5C | [ ] | DiscussionTemplate Types |
| └─ `import --type discussion-template` | 5C | [ ] | DiscussionTemplate Types |

---

## Phased Implementation

### Phase 1: Foundation (Parallel Streams)

All Phase 1 work streams can be developed **in parallel** with no dependencies on each other.

```
┌─────────────────────────────────────────────────────────────────────────┐
│ PHASE 1: Foundation                                                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐ │
│  │     1A       │  │     1B       │  │     1C       │  │     1D       │ │
│  │  Core Types  │  │ Serialization│  │     CLI      │  │Schema Fetch  │ │
│  │              │  │              │  │  Framework   │  │              │ │
│  │  workflow/   │  │  internal/   │  │    cmd/      │  │  codegen/    │ │
│  │  *.go        │  │  serialize/  │  │  wetwire-    │  │  fetch.go    │ │
│  │              │  │              │  │  github/     │  │              │ │
│  └──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘ │
│         │                 │                 │                 │         │
│         │                 │                 │                 │         │
│         ▼                 ▼                 ▼                 ▼         │
│    No deps           Needs 1A          No deps           No deps       │
│                      (can stub)                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 1A: Core Types (`workflow/`)
- [ ] `workflow.go` — Workflow struct with ResourceType()
- [ ] `job.go` — Job struct with OutputRef fields
- [ ] `step.go` — Step struct with Output() method
- [ ] `matrix.go` — Matrix, Strategy structs
- [ ] `triggers.go` — All trigger types (Push, PullRequest, Schedule, etc.)
- [ ] `conditions.go` — Condition interface + builders
- [ ] `expressions.go` — Context accessors (GitHub, Runner, Env, Secrets, Matrix)
- [ ] `helpers.go` — Type aliases (Env, With) + helpers (List, Any, Strings)
- [ ] `contracts.go` — Root-level interfaces (WorkflowResource, OutputRef, etc.)

**Output:** Compilable types, no runtime behavior yet

#### 1B: Serialization (`internal/serialize/`)
- [ ] `workflow.go` — Workflow → YAML map
- [ ] `job.go` — Job → YAML map
- [ ] `step.go` — Step → YAML map
- [ ] `matrix.go` — Matrix → YAML map
- [ ] `triggers.go` — Triggers → YAML map
- [ ] `conditions.go` — Condition → expression string
- [ ] `yaml.go` — Final YAML output with proper formatting

**Note:** Can start with stub types, integrate with 1A when ready

#### 1C: CLI Framework (`cmd/wetwire-github/`)
- [ ] `main.go` — Cobra root command
- [ ] `init.go` — Project scaffolding (functional)
- [ ] `build.go` — Stub returning "not implemented"
- [ ] `validate.go` — Stub returning "not implemented"
- [ ] `list.go` — Stub returning "not implemented"
- [ ] `lint.go` — Stub returning "not implemented"
- [ ] `import.go` — Stub returning "not implemented"
- [ ] `version.go` — Version output

**Output:** Working CLI with `init` functional, other commands stubbed

#### 1D: Schema Fetching (`codegen/`)
- [ ] `fetch.go` — HTTP fetcher with retry
- [ ] `manifest.go` — Manifest tracking
- [ ] `config.go` — Action URLs configuration
- [ ] Fetch workflow schema from schemastore
- [ ] Fetch action.yml for each configured action
- [ ] Write to `specs/` directory

**Output:** `specs/` populated with workflow-schema.json and action.yml files

---

### Phase 2: Core Capabilities (Parallel Streams)

Phase 2 streams have dependencies on Phase 1 but are **parallel to each other**.

```
┌─────────────────────────────────────────────────────────────────────────┐
│ PHASE 2: Core Capabilities                                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐     │
│  │   2A   │ │   2B   │ │   2C   │ │   2D   │ │   2E   │ │   2F   │     │
│  │ Schema │ │ Action │ │  AST   │ │action- │ │ Linter │ │  YAML  │     │
│  │ Parse  │ │Codegen │ │Discover│ │  lint  │ │ Rules  │ │ Parser │     │
│  └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘     │
│      │          │          │          │          │          │          │
│      ▼          ▼          ▼          ▼          ▼          ▼          │
│   Needs 1D  Needs 2A   Needs 1A   No deps    No deps    No deps        │
│                                                                          │
│  ┌────────┐ ┌────────┐                                                  │
│  │   2G   │ │   2H   │                                                  │
│  │ Runner │ │Template│                                                  │
│  │  Value │ │Builder │                                                  │
│  └───┬────┘ └───┬────┘                                                  │
│      │          │                                                        │
│      ▼          ▼                                                        │
│   Needs 2C  Needs 2C                                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2A: Schema Parsing (`codegen/`)
- [ ] `parse.go` — Parse workflow-schema.json
- [ ] `parse_action.go` — Parse action.yml files
- [ ] `schema.go` — Intermediate representation types

**Depends on:** 1D (Schema Fetching)

#### 2B: Action Codegen (`codegen/`)
- [ ] `generate.go` — Go code generation from parsed actions
- [ ] `templates/action.go.tmpl` — Template for action wrappers
- [ ] Code formatting with go/format
- [ ] Generate `actions/checkout/checkout.go`
- [ ] Generate `actions/setup_python/setup_python.go`
- [ ] Generate `actions/setup_node/setup_node.go`
- [ ] Generate `actions/setup_go/setup_go.go`
- [ ] Generate `actions/cache/cache.go`
- [ ] Generate `actions/upload_artifact/upload_artifact.go`
- [ ] Generate `actions/download_artifact/download_artifact.go`

**Depends on:** 2A (Schema Parsing)

#### 2C: AST Discovery (`internal/discover/`)
- [ ] `discover.go` — Package scanning with go/parser
- [ ] `workflow.go` — Workflow variable detection
- [ ] `job.go` — Job variable detection
- [ ] `graph.go` — Dependency graph building
- [ ] `validate.go` — Reference validation
- [ ] Recursive directory traversal
- [ ] Vendor/hidden directory exclusion

**Depends on:** 1A (Core Types)

#### 2D: Actionlint Integration (`internal/validation/`)
- [ ] `actionlint.go` — Library wrapper
- [ ] `types.go` — ValidationResult, ValidationIssue
- [ ] `pipeline.go` — Multiple validator pipeline

**Depends on:** None (external library)

#### 2E: Linter Rules (`internal/linter/`)
- [ ] `linter.go` — Framework and runner
- [ ] `rule.go` — Rule interface (ID, Description, Check)
- [ ] `wag001.go` — Use typed action wrappers
- [ ] `wag002.go` — Use condition builders
- [ ] `wag003.go` — Use secrets context
- [ ] `wag004.go` — Use matrix builder
- [ ] `wag005.go` — Inline structs → named vars
- [ ] `wag006.go` — Duplicate workflow names
- [ ] `wag007.go` — File too large
- [ ] `wag008.go` — Hardcoded expression strings
- [ ] Recursive package scanning
- [ ] `--fix` auto-remediation support

**Depends on:** None

#### 2F: YAML Parser (`internal/importer/`)
- [ ] `parser.go` — Parse workflow YAML files
- [ ] `ir.go` — Intermediate representation (IRWorkflow, IRJob, IRStep)
- [ ] Reference graph tracking

**Depends on:** None

#### 2G: Runner / Value Extraction (`internal/runner/`)
- [ ] `runner.go` — Temp Go program generation
- [ ] `module.go` — go.mod parsing
- [ ] `replace.go` — Replace directive handling
- [ ] `gobinary.go` — Go binary discovery (PATH + common locations)
- [ ] JSON value extraction and parsing

**Depends on:** 2C (AST Discovery)

#### 2H: Template Builder (`internal/template/`)
- [ ] `builder.go` — Template construction from discovered resources
- [ ] `sort.go` — Topological sort (Kahn's algorithm)
- [ ] `cycle.go` — Cycle detection with error messages
- [ ] Dependency ordering

**Depends on:** 2C (AST Discovery)

---

### Phase 3: Command Integration (Parallel Streams)

Phase 3 integrates Phase 2 capabilities into CLI commands. Each command can be completed **in parallel**.

```
┌─────────────────────────────────────────────────────────────────────────┐
│ PHASE 3: Command Integration                                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │   3A    │ │   3B    │ │   3C    │ │   3D    │ │   3E    │           │
│  │  build  │ │validate │ │  lint   │ │ import  │ │  list   │           │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘           │
│       │           │           │           │           │                 │
│       ▼           ▼           ▼           ▼           ▼                 │
│   2C,2G,2H      2D          2E        1A,2F        2C                  │
│                                                                          │
│  ┌─────────┐ ┌─────────┐                                                │
│  │   3F    │ │   3G    │                                                │
│  │ design  │ │  test   │                                                │
│  └────┬────┘ └────┬────┘                                                │
│       │           │                                                      │
│       ▼           ▼                                                      │
│  All Phase 3  All Phase 3                                               │
│   + core-go    + core-go                                                │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 3A: Build Command (full)
- [ ] Discovery integration (use 2C)
- [ ] Runner integration (use 2G)
- [ ] Template builder integration (use 2H)
- [ ] Multi-workflow support
- [ ] Output to `.github/workflows/`
- [ ] `--format json/yaml` for agent integration
- [ ] `--output` flag

**Depends on:** 1B (Serialization), 2C (AST Discovery), 2G (Runner), 2H (Template Builder)

#### 3B: Validate Command (full)
- [ ] Integrate actionlint (use 2D)
- [ ] Validate generated YAML
- [ ] `--format json/text` output

**Depends on:** 2D (Actionlint Integration)

#### 3C: Lint Command (full)
- [ ] Integrate linter framework (use 2E)
- [ ] Recursive package scanning (`./...` pattern)
- [ ] `--fix` flag support
- [ ] `--format json/text` output
- [ ] Exit code 2 for lint issues

**Depends on:** 2E (Linter Rules)

#### 3D: Import Command (full)
- [ ] YAML parsing (use 2F)
- [ ] Go code generation from IR
- [ ] Field name transformation (kebab-case → snake_case)
- [ ] Reserved name handling
- [ ] Scaffold: go.mod with module path
- [ ] Scaffold: cmd/main.go
- [ ] Scaffold: README.md
- [ ] Scaffold: CLAUDE.md (AI context)
- [ ] Scaffold: .gitignore
- [ ] `--single-file` flag
- [ ] `--no-scaffold` flag
- [ ] `--package` flag
- [ ] `--module` flag

**Depends on:** 1A (Core Types), 2F (YAML Parser)

#### 3E: List Command (full)
- [ ] Discovery integration (use 2C)
- [ ] Table output (name, type, file, line)
- [ ] `--format json` output

**Depends on:** 2C (AST Discovery)

#### 3F: Design Command (AI-assisted)
- [ ] wetwire-core-go orchestrator integration
- [ ] Interactive session with developer
- [ ] Lint feedback loop
- [ ] `--stream` flag for response streaming
- [ ] `--max-lint-cycles` flag
- [ ] `--output` working directory

**Depends on:** 3A-3E (All basic commands), wetwire-core-go

#### 3G: Test Command (Persona-based)
- [ ] Persona selection (beginner, intermediate, expert, terse, verbose)
- [ ] Scenario configuration
- [ ] Result writing (RESULTS.md, session.json, score.json)
- [ ] Session tracking (questions, lint cycles)
- [ ] `--persona` flag
- [ ] `--scenario` flag
- [ ] `--output` directory

**Depends on:** 3A-3E (All basic commands), wetwire-core-go

---

### Phase 4: Polish & Integration (Parallel Streams)

```
┌─────────────────────────────────────────────────────────────────────────┐
│ PHASE 4: Polish & Integration                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌────────────────────────────┐  ┌────────────────────────────────────┐ │
│  │            4A              │  │              4B                    │ │
│  │    Examples & Testing      │  │     wetwire-core-go Integration   │ │
│  │                            │  │                                    │ │
│  │  - Fetch starter-workflows │  │  - RunnerAgent tool definitions   │ │
│  │  - Import/rebuild tests    │  │  - Agent testing with personas    │ │
│  │  - Round-trip validation   │  │  - Scoring integration            │ │
│  └────────────────────────────┘  └────────────────────────────────────┘ │
│              │                                │                         │
│              ▼                                ▼                         │
│       Needs Phase 3                    Needs Phase 3                    │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 4A: Examples & Testing
- [ ] Fetch `actions/starter-workflows` examples
- [ ] Import each example → Go code
- [ ] Build Go code → YAML
- [ ] Diff original vs rebuilt
- [ ] Validate with actionlint
- [ ] Track success rate

**Depends on:** All Phase 3 commands

#### 4B: wetwire-core-go Integration
- [ ] Define RunnerAgent tools
- [ ] Implement tool handlers
- [ ] Test with DeveloperAgent + personas
- [ ] Scoring integration

**Depends on:** All Phase 3 commands

---

### Phase 5: Extended Config Types (Parallel Streams)

Phase 5 adds support for additional GitHub YAML configuration types. All streams can be developed **in parallel** with each other and with Phase 4.

```
┌─────────────────────────────────────────────────────────────────────────┐
│ PHASE 5: Extended Config Types                                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│  ┌─────────────────────┐  ┌─────────────────────┐  ┌─────────────────┐  │
│  │         5A          │  │         5B          │  │       5C        │  │
│  │     Dependabot      │  │   Issue Templates   │  │   Discussion    │  │
│  │                     │  │                     │  │   Templates     │  │
│  │  dependabot/        │  │  templates/         │  │                 │  │
│  │  *.go               │  │  issue.go           │  │  templates/     │  │
│  │                     │  │  form.go            │  │  discussion.go  │  │
│  │                     │  │  elements.go        │  │                 │  │
│  └─────────────────────┘  └─────────────────────┘  └─────────────────┘  │
│           │                        │                       │            │
│           ▼                        ▼                       ▼            │
│      No deps from             No deps from          Needs 5B           │
│      Phase 1-4                Phase 1-4             (FormBody)          │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 5A: Dependabot (`dependabot/`)
- [ ] `dependabot.go` — Dependabot struct (version, registries, updates)
- [ ] `update.go` — Update struct (package-ecosystem, directory, schedule)
- [ ] `schedule.go` — Schedule struct (interval, day, time, timezone)
- [ ] `registries.go` — Registry types (npm, docker, maven, etc.)
- [ ] `groups.go` — Groups struct (patterns, dependency-type)
- [ ] `contracts.go` — Dependabot interface additions
- [ ] Fetch `dependabot-2.0.json` schema
- [ ] Parse dependabot schema
- [ ] Dependabot → YAML serialization
- [ ] AST discovery for Dependabot variables
- [ ] `build --type dependabot` output to `.github/dependabot.yml`
- [ ] `import --type dependabot` YAML → Go conversion
- [ ] `validate --type dependabot` schema validation

**Depends on:** 1B (Serialization patterns), 1C (CLI framework)

#### 5B: Issue Templates (`templates/`)
- [ ] `issue.go` — IssueTemplate struct (name, description, title, labels, assignees)
- [ ] `form.go` — FormBody struct (body array)
- [ ] `elements.go` — FormElement interface + implementations:
  - [ ] Input element (label, description, placeholder, value, required)
  - [ ] Textarea element (label, description, placeholder, value, render)
  - [ ] Dropdown element (label, description, options, default, multiple)
  - [ ] Checkboxes element (label, description, options)
  - [ ] Markdown element (value)
- [ ] `contracts.go` — IssueTemplate interface additions
- [ ] Fetch `github-issue-forms.json` schema
- [ ] Parse issue forms schema
- [ ] IssueTemplate → YAML serialization
- [ ] AST discovery for IssueTemplate variables
- [ ] `build --type issue-template` output to `.github/ISSUE_TEMPLATE/`
- [ ] `import --type issue-template` YAML → Go conversion
- [ ] `validate --type issue-template` schema validation

**Depends on:** 1B (Serialization patterns), 1C (CLI framework)

#### 5C: Discussion Templates (`templates/`)
- [ ] `discussion.go` — DiscussionTemplate struct (title, labels, body)
- [ ] Discussion category handling
- [ ] AST discovery for DiscussionTemplate variables
- [ ] `build --type discussion-template` output to `.github/DISCUSSION_TEMPLATE/`
- [ ] `import --type discussion-template` YAML → Go conversion

**Depends on:** 5B (FormBody reuse from Issue Templates)

---

## Parallel Development Matrix

This matrix shows which work can happen simultaneously:

| Week | Stream 1 | Stream 2 | Stream 3 | Stream 4 | Stream 5 |
|------|----------|----------|----------|----------|----------|
| 1 | 1A: Core Types | 1C: CLI Framework | 1D: Schema Fetch | — | — |
| 2 | 1A: (cont.) | 1B: Serialization | 2D: actionlint | 2E: Linter Rules | 2F: YAML Parser |
| 3 | 2C: AST Discovery | 2A: Schema Parse | 2E: (cont.) | 2F: (cont.) | 5A: Dependabot |
| 4 | 2G: Runner | 2H: Template | 2B: Action Codegen | 2E: (cont.) | 5A: (cont.) |
| 5 | 3A: build | 3B: validate | 3C: lint | 3D: import | 5B: Issue Templates |
| 6 | 3A: (cont.) | 3D: (cont.) | 3C: (cont.) | 3E: list | 5B: (cont.) |
| 7 | 4A: Examples | 4B: core-go | 3F: design | 3G: test | 5C: Discussion |
| 8 | 4A: (cont.) | 4B: (cont.) | 3F: (cont.) | 3G: (cont.) | 5C: (cont.) |

---

## Progress Tracking

### Phase 0 Progress
- [ ] Repository Setup (0/4)
- [ ] GitHub Actions CI (0/2)
- [ ] Development Scripts (0/2)
- [ ] Documentation (0/4)
- [ ] Build Artifacts (0/1)

### Phase 1 Progress
- [ ] 1A: Core Types (0/11)
- [ ] 1B: Serialization (0/9)
- [ ] 1C: CLI Framework (0/11)
- [ ] 1D: Schema Fetching (0/10)

### Phase 2 Progress
- [ ] 2A: Schema Parsing (0/3)
- [ ] 2B: Action Codegen (0/10)
- [ ] 2C: AST Discovery (0/7)
- [ ] 2D: Actionlint Integration (0/3)
- [ ] 2E: Linter Rules (0/12)
- [ ] 2F: YAML Parser (0/3)
- [ ] 2G: Runner/Value Extraction (0/5)
- [ ] 2H: Template Builder (0/4)

### Phase 3 Progress
- [ ] 3A: Build Command (0/7)
- [ ] 3B: Validate Command (0/3)
- [ ] 3C: Lint Command (0/5)
- [ ] 3D: Import Command (0/13)
- [ ] 3E: List Command (0/3)
- [ ] 3F: Design Command (0/6)
- [ ] 3G: Test Command (0/7)

### Phase 4 Progress
- [ ] 4A: Examples & Testing (0/4)
- [ ] 4B: wetwire-core-go Integration (0/6)

### Phase 5 Progress
- [ ] 5A: Dependabot (0/13)
- [ ] 5B: Issue Templates (0/14)
- [ ] 5C: Discussion Templates (0/5)

---

## Critical Path

The minimum sequence to reach a working `build` command:

```
1A (Core Types) → 1B (Serialization) ─┐
                                      ├→ 2C (AST Discovery) → 2G (Runner) ─┐
                                      │                      2H (Template) ┼→ 3A (build)
                                      └────────────────────────────────────┘
```

The minimum sequence to reach a working `import` command:

```
1A (Core Types) ─┐
                 ├→ 3D (import)
2F (YAML Parser) ┘
```

The minimum sequence to reach a working `validate` command:

```
2D (actionlint) → 3B (validate)
```

The minimum sequence to reach a working `lint` command:

```
2E (Linter Rules) → 3C (lint)
```

**Key insight:** `import`, `validate`, and `lint` can all proceed in parallel with `build`, meeting at Phase 4 for round-trip testing.

The minimum sequence to reach Dependabot support:

```
1B (Serialization) ─┐
                    ├→ 5A (Dependabot)
1C (CLI Framework) ─┘
```

The minimum sequence to reach Issue Template support:

```
1B (Serialization) ─┐
                    ├→ 5B (Issue Templates) → 5C (Discussion Templates)
1C (CLI Framework) ─┘
```

**Key insight:** Phase 5 streams (Dependabot, Issue Templates) can begin as early as Phase 1 completes, running in parallel with Phases 2-4.

## Feature Count Summary

| Phase | Streams | Features |
|-------|---------|----------|
| Phase 0 | 5 | 13 |
| Phase 1 | 4 | 41 |
| Phase 2 | 8 | 47 |
| Phase 3 | 7 | 44 |
| Phase 4 | 2 | 10 |
| Phase 5 | 3 | 32 |
| **Total** | **29** | **187** |
