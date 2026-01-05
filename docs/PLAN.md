# wetwire-github-go Implementation Plan

## Overview

Build `wetwire-github-go` following the same patterns as `wetwire-aws-go` — a synthesis library that generates GitHub YAML configurations from typed Go declarations.

## Scope

wetwire-github-go generates typed Go declarations for three GitHub YAML configuration types:

| Config Type | Output Location | Schema Source |
|-------------|-----------------|---------------|
| **GitHub Actions** | `.github/workflows/*.yml` | `json.schemastore.org/github-workflow.json` |
| **Dependabot** | `.github/dependabot.yml` | `json.schemastore.org/dependabot-2.0.json` |
| **Issue/Discussion Templates** | `.github/ISSUE_TEMPLATE/*.yml` | `json.schemastore.org/github-issue-forms.json` |

## Key Decisions

- **Action versions**: Hardcode major version in generated wrappers (e.g., `@v4`)
- **Validation**: Use actionlint as Go library (`github.com/rhysd/actionlint`)
- **Approach**: Full feature parity with wetwire-aws-go
- **Config types**: Support Actions, Dependabot, and Issue/Discussion Templates via `--type` flag

---

## Feature Matrix: wetwire-aws-go → wetwire-github-go

| Feature | wetwire-aws-go | wetwire-github-go |
|---------|----------------|-------------------|
| **Schema Source** | CloudFormation spec JSON | Workflow JSON schema + action.yml files |
| **Schema URL** | AWS CF spec URL | `json.schemastore.org/github-workflow.json` |
| **Secondary Source** | — | Popular action.yml files (checkout, setup-python, etc.) |
| **Output Format** | CloudFormation JSON/YAML | GitHub Actions workflow YAML |
| **Generated Types** | 262 AWS service packages | Workflow, Job, Step, Matrix, Triggers + Action wrappers |
| **Intrinsics** | Ref, GetAtt, Sub, Join, etc. | Expression contexts (github, runner, env, secrets, matrix) |
| **Validation** | cfn-lint integration | actionlint integration |

---

## Schema Sources

### 1. Workflow Schema
- **URL**: `https://json.schemastore.org/github-workflow.json`
- **Raw**: `https://raw.githubusercontent.com/SchemaStore/schemastore/master/src/schemas/json/github-workflow.json`
- **Provides**: Triggers, jobs, steps, matrix, concurrency, permissions, environments

### 2. Action Metadata (action.yml)
- **Pattern**: `https://raw.githubusercontent.com/{owner}/{repo}/main/action.yml`
- **Popular actions to generate wrappers for**:
  - `actions/checkout`
  - `actions/setup-python`, `setup-node`, `setup-go`, `setup-java`
  - `actions/cache`
  - `actions/upload-artifact`, `download-artifact`
  - `docker/build-push-action`
  - `codecov/codecov-action`
  - `pypa/gh-action-pypi-publish`
  - (extensible list in config)

### 3. Dependabot Schema
- **URL**: `https://json.schemastore.org/dependabot-2.0.json`
- **Provides**: Package ecosystems, update schedules, registries, groups, ignore patterns

### 4. Issue Forms Schema
- **URL**: `https://json.schemastore.org/github-issue-forms.json`
- **Provides**: Form body elements (input, textarea, dropdown, checkboxes, markdown)

---

## Directory Structure

```
wetwire-github-go/
├── .github/
│   └── workflows/
│       ├── ci.yml              # Build/test on push/PR
│       └── codebot.yml         # Claude Code integration
│
├── scripts/
│   ├── ci.sh                   # Local CI runner
│   └── import_samples.sh       # Round-trip testing
│
├── cmd/wetwire-github/         # CLI application
│   ├── main.go
│   ├── build.go                # Generate YAML (--type workflow|dependabot|issue-template)
│   ├── validate.go             # Validate via actionlint
│   ├── list.go                 # List discovered resources
│   ├── lint.go                 # Code quality rules
│   ├── import.go               # YAML → Go conversion
│   ├── init.go                 # Project scaffolding
│   └── version.go
│
├── internal/
│   ├── discover/               # AST-based resource discovery
│   ├── importer/               # YAML to Go code conversion
│   ├── linter/                 # Go code lint rules (WAG001-WAG0XX)
│   ├── template/               # YAML builder
│   ├── serialize/              # YAML serialization
│   └── validation/             # actionlint integration
│
├── workflow/                   # Core workflow types (hand-written)
│   ├── workflow.go             # Workflow struct
│   ├── job.go                  # Job struct
│   ├── step.go                 # Step struct
│   ├── matrix.go               # Matrix builder
│   ├── triggers.go             # Push, PullRequest, Schedule, etc.
│   ├── conditions.go           # Condition builders
│   └── expressions.go          # github, runner, env, secrets contexts
│
├── dependabot/                 # Dependabot types (hand-written)
│   ├── dependabot.go           # Dependabot struct
│   ├── update.go               # Update struct
│   ├── schedule.go             # Schedule struct
│   ├── registries.go           # Registry types
│   └── groups.go               # Grouping configuration
│
├── templates/                  # Issue/Discussion template types (hand-written)
│   ├── issue.go                # IssueTemplate struct
│   ├── discussion.go           # DiscussionTemplate struct
│   ├── form.go                 # FormBody struct
│   └── elements.go             # Input, Textarea, Dropdown, Checkboxes, Markdown
│
├── actions/                    # GENERATED action wrappers
│   ├── checkout/
│   │   └── checkout.go         # Typed wrapper for actions/checkout
│   ├── setup_python/
│   │   └── setup_python.go
│   ├── cache/
│   └── ... (top 20+ actions)
│
├── codegen/                    # Code generation tooling
│   ├── fetch.go                # Download schemas + action.yml
│   ├── parse.go                # Parse schemas
│   └── generate.go             # Generate Go types
│
├── examples/                   # Example configs for import testing
│   └── (fetched from real repos)
│
├── docs/
│   ├── PLAN.md                 # This file
│   ├── ROADMAP.md              # Detailed feature matrix and phased roadmap
│   ├── QUICK_START.md
│   ├── CLI.md
│   └── IMPORT_WORKFLOW.md
│
├── specs/                      # .gitignore'd (fetched schemas)
│   ├── .gitkeep
│   ├── manifest.json
│   ├── workflow-schema.json
│   ├── dependabot-schema.json
│   ├── issue-forms-schema.json
│   └── actions/
│       ├── checkout.yml
│       └── ...
│
├── contracts.go                # Interfaces and types
├── go.mod
└── README.md
```

---

## CLI Commands

### 1. `wetwire-github build`
- Discover workflow declarations from Go packages
- Serialize to GitHub Actions YAML
- Output to `.github/workflows/` or custom path

### 2. `wetwire-github validate`
- Run actionlint on generated YAML (via Go library)
- Report errors in structured format (text/JSON)

### 3. `wetwire-github list`
- List discovered workflows and jobs
- Show file locations and dependencies

### 4. `wetwire-github lint`
- Go code quality rules (WAG001-WAG0XX)
- Examples:
  - WAG001: Use typed action wrappers instead of raw `uses:` strings
  - WAG002: Use condition builders instead of raw expression strings
  - WAG003: Use secrets context instead of hardcoded strings
  - WAG004: Use matrix builder instead of inline maps

### 5. `wetwire-github import`
- Convert existing workflow YAML to Go code
- Generate typed declarations
- Scaffold project structure

### 6. `wetwire-github init`
- Create new project with example workflow
- Generate go.mod, main.go, workflow definitions

---

## Core Types (workflow/ package)

```go
// workflow.go
type Workflow struct {
    Name        string
    On          Triggers
    Env         map[string]any
    Defaults    *Defaults
    Concurrency *Concurrency
    Permissions *Permissions
    Jobs        []Job // discovered via registry
}

// job.go
type Job struct {
    Name           string
    RunsOn         any  // string or matrix ref
    Needs          []any // Job references (type-safe)
    If             Condition
    Environment    *Environment
    Concurrency    *Concurrency
    Outputs        map[string]Expression
    Strategy       *Strategy
    Container      *Container
    Services       map[string]Service
    Steps          []Step
    TimeoutMinutes int
}

// step.go
type Step struct {
    Name    string
    ID      string
    If      Condition
    Uses    string      // or Action interface
    With    map[string]any
    Run     string
    Shell   string
    Env     map[string]any
    WorkingDirectory string
}

// matrix.go
type Matrix struct {
    Values   map[string][]any
    Include  []map[string]any
    Exclude  []map[string]any
    FailFast bool
    MaxParallel int
}

// triggers.go
type Triggers struct {
    Push              *PushTrigger
    PullRequest       *PullRequestTrigger
    PullRequestTarget *PullRequestTargetTrigger
    Schedule          []ScheduleTrigger
    WorkflowDispatch  *WorkflowDispatchTrigger
    WorkflowCall      *WorkflowCallTrigger
    // ... 30+ event types
}

// conditions.go
type Condition interface {
    String() string  // generates ${{ }} expression
}

func Push() Condition                    // github.event_name == 'push'
func PullRequest() Condition             // github.event_name == 'pull_request'
func Branch(name string) Condition       // github.ref == 'refs/heads/{name}'
func TagRef(pattern string) Condition    // startsWith(github.ref, 'refs/tags/{pattern}')
func Always() Condition                  // always()
func Failure() Condition                 // failure()
func Success() Condition                 // success()

// Operators
func (c Condition) And(other Condition) Condition  // &&
func (c Condition) Or(other Condition) Condition   // ||
func (c Condition) Not() Condition                 // !

// expressions.go - Context accessors
var GitHub = githubContext{}   // github.event_name, github.ref, etc.
var Runner = runnerContext{}   // runner.os, runner.arch
var Env = envContext{}         // env.MY_VAR
var Secrets = secretsContext{} // secrets.MY_SECRET
var Matrix = matrixContext{}   // matrix.python_version
```

---

## Generated Action Wrappers (actions/ package)

```go
// actions/checkout/checkout.go
// Generated from actions/checkout/action.yml

package checkout

type Checkout struct {
    Repository         string
    Ref                string
    Token              string
    SSHKey             string
    PersistCredentials bool
    Path               string
    Clean              bool
    FetchDepth         int
    FetchTags          bool
    Submodules         string
    LFS                bool
    // ... all inputs from action.yml
}

func (a Checkout) Action() string { return "actions/checkout@v4" }
func (a Checkout) ToStep() workflow.Step { ... }

// Output references
func (a Checkout) OutputRef() OutputRef       { return OutputRef{a, "ref"} }
func (a Checkout) OutputCommit() OutputRef   { return OutputRef{a, "commit"} }
```

---

## Integration with wetwire-core-go

Same pattern as wetwire-aws-go:

```go
// RunnerAgent tools for wetwire-github
tools := []Tool{
    "init_package",      // wetwire-github init
    "write_file",        // Write Go workflow code
    "read_file",         // Read files
    "run_lint",          // wetwire-github lint --format json
    "run_build",         // wetwire-github build --format json
    "run_validate",      // wetwire-github validate (actionlint)
    "ask_developer",     // Clarification questions
}
```

CLI must support `--format json` for agent integration.

---

## Actionlint Integration (Go Library)

```go
// internal/validation/actionlint.go
import "github.com/rhysd/actionlint"

func ValidateWorkflow(path string, content []byte) (*ValidationResult, error) {
    linter, _ := actionlint.NewLinter(io.Discard, nil)
    errs, _ := linter.Lint(path, content, nil)

    result := &ValidationResult{Success: len(errs) == 0}
    for _, e := range errs {
        result.Issues = append(result.Issues, ValidationIssue{
            File:     e.Filepath,
            Line:     e.Line,
            Column:   e.Column,
            Message:  e.Message,
            RuleID:   e.Kind,
        })
    }
    return result, nil
}
```

No external CLI dependency required.

---

## Examples Fetch & Import Workflow

### Phase 1: Fetch Popular Workflows
```bash
# Fetch example workflows from popular repos
wetwire-github fetch-examples --output examples/
```

Sources:
- `github.com/actions/starter-workflows` (official templates)
- Popular open source projects with complex CI

### Phase 2: Import & Validate
```bash
# Import each example
for f in examples/*.yml; do
    wetwire-github import "$f" --output "imported/$(basename $f .yml)"
done

# Build back to YAML
for d in imported/*/; do
    wetwire-github build "$d" --output "rebuilt/"
done

# Validate rebuilt workflows
wetwire-github validate rebuilt/*.yml
```

### Phase 3: Improvement Cycle
1. Import workflow → generates Go code
2. Build Go code → generates YAML
3. Diff original vs rebuilt
4. Fix codegen issues
5. Repeat until 100% fidelity

---

## Implementation Phases

See [ROADMAP.md](ROADMAP.md) for detailed feature matrix with 186 features across 29 parallel streams.

### Phase 0: Development Infrastructure
- **Repository Setup** — go.mod, .gitignore, specs/.gitkeep
- **GitHub Actions CI** — ci.yml (build/test), codebot.yml (Claude Code integration)
- **Development Scripts** — ci.sh (local CI), import_samples.sh (round-trip testing)
- **Documentation** — README.md, QUICK_START.md, CLI.md, IMPORT_WORKFLOW.md

### Phase 1: Foundation (Parallel Streams)
- **1A: Core Types** — Workflow, Job, Step, Matrix, Triggers, Conditions, Expressions
- **1B: Serialization** — YAML output for all core types
- **1C: CLI Framework** — Cobra setup, command stubs, init command
- **1D: Schema Fetching** — HTTP fetcher, workflow schema, action.yml files

### Phase 2: Core Capabilities (Parallel Streams)
- **2A: Schema Parsing** — Parse workflow and action schemas
- **2B: Action Codegen** — Generate action wrappers from parsed schemas
- **2C: AST Discovery** — Package scanning, variable detection, dependency graph
- **2D: Actionlint Integration** — Library wrapper, validation pipeline
- **2E: Linter Rules** — WAG001-WAG008 rules, --fix support
- **2F: YAML Parser** — Parse existing workflows for import
- **2G: Runner/Value Extraction** — Temp program generation, JSON extraction
- **2H: Template Builder** — Topological sort, cycle detection

### Phase 3: Command Integration (Parallel Streams)
- **3A: build** — Full build command with discovery, runner, template builder
- **3B: validate** — Full validate with actionlint
- **3C: lint** — Full lint with all rules
- **3D: import** — Full import with scaffolding
- **3E: list** — Full list with table/JSON output
- **3F: design** — AI-assisted design via wetwire-core-go
- **3G: test** — Persona-based testing via wetwire-core-go

### Phase 4: Polish & Integration (Parallel Streams)
- **4A: Examples & Testing** — Fetch starter-workflows, round-trip validation
- **4B: wetwire-core-go Integration** — RunnerAgent tools, scoring

### Phase 5: Extended Config Types (Parallel Streams)
- **5A: Dependabot** — Types, schema, serialization, CLI integration (--type dependabot)
- **5B: Issue Templates** — Types, schema, serialization, CLI integration (--type issue-template)
- **5C: Discussion Templates** — Types, CLI integration (--type discussion-template)

---

## Reference Files (from wetwire-aws-go)

Key files to mirror structure from:
- `wetwire-aws-go/contracts.go` — Interface pattern
- `wetwire-aws-go/cmd/wetwire-aws/main.go` — CLI structure
- `wetwire-aws-go/codegen/generate.go` — Codegen pattern
- `wetwire-aws-go/internal/discover/discover.go` — AST discovery
- `wetwire-aws-go/internal/importer/codegen.go` — Import codegen

---

## Open Questions (Deferred)

1. **Composite actions**: Support generating wrappers for composite actions with multiple steps?

2. **Reusable workflows**: Support `workflow_call` trigger and typed inputs/outputs?

3. **Local actions**: Support `uses: ./path/to/action` syntax?
