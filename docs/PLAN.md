# wetwire-github-go Implementation Plan

## Overview

Build `wetwire-github-go` following the same patterns as `wetwire-aws-go` — a synthesis library that generates GitHub YAML configurations from typed Go declarations.

## The "No Parens" Pattern

Resources are declared as Go variables using struct literals — no function calls or registration needed:

```go
// Flat variables for optional nested structs (pointer types declared once with &)
var CIPush = workflow.PushTrigger{Branches: List("main")}
var CIPullRequest = workflow.PullRequestTrigger{Branches: List("main")}

// Workflow declaration - clean references, no & at usage site
var CI = workflow.Workflow{
    Name: "CI",
    On: workflow.Triggers{
        Push:        CIPush,
        PullRequest: CIPullRequest,
    },
}

// Job declaration - automatically associated via AST discovery
var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps: []workflow.Step{
        checkout.Checkout{}.ToStep(),
        setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
        {Run: "go build ./..."},
        {Run: "go test ./..."},
    },
}

// Cross-references via direct field access
var Deploy = workflow.Job{
    Needs: []any{Build, Test},  // Automatic dependency resolution
    Steps: []workflow.Step{
        {
            If:  workflow.Branch("main"),
            Run: "deploy.sh",
            Env: map[string]any{
                "TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
            },
        },
    },
}
```

**Key principles:**
- Variables declared with struct literals (no function calls)
- Pointer types declared once with `&`, referenced cleanly without `&`
- Cross-resource references via direct field access
- AST-based discovery — no registration needed
- Type-safe action wrappers with `.ToStep()` conversion
- Expression contexts as typed accessors (`workflow.Secrets.Get(...)`)

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

## Generated Package Structure (User Projects)

When users run `wetwire-github import` or `wetwire-github init`, the importer scaffolds a project package:

```
my-ci/                           # User's workflow package
├── go.mod                       # Module: my-ci
├── go.sum
├── README.md                    # Generated docs
├── CLAUDE.md                    # AI assistant context
├── .gitignore
│
├── cmd/
│   └── main.go                  # Usage instructions
│
├── workflows.go                 # Workflow declarations
├── jobs.go                      # Job declarations
├── steps.go                     # Reusable step declarations
├── triggers.go                  # Trigger configurations
└── matrix.go                    # Matrix configurations
```

**Example generated files:**

**go.mod:**
```go
module my-ci

go 1.23

require github.com/lex00/wetwire-github-go v0.1.0
```

**workflows.go:**
```go
package my_ci

import (
    "github.com/lex00/wetwire-github-go/workflow"
)

var CIPush = workflow.PushTrigger{Branches: List("main")}
var CIPullRequest = workflow.PullRequestTrigger{Branches: List("main")}

var CI = workflow.Workflow{
    Name: "CI",
    On: workflow.Triggers{
        Push:        CIPush,
        PullRequest: CIPullRequest,
    },
}
```

**jobs.go:**
```go
package my_ci

import (
    "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
)

var BuildSteps = []workflow.Step{
    checkout.Checkout{}.ToStep(),
    setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
    {Run: "go build ./..."},
    {Run: "go test ./..."},
}

var Build = workflow.Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}
```

**cmd/main.go:**
```go
package main

import "fmt"

func main() {
    // Build workflows using the wetwire-github CLI:
    //   wetwire-github build .
    //
    // This generates .github/workflows/*.yml from Go definitions.
    fmt.Println("Usage: wetwire-github build .")
}
```

**Key patterns (same as wetwire-aws-go):**
- Single package per project
- Flat variables for all nested structs
- Importer generates correct `&` based on field types
- Cross-file references work via Go's package scope
- Variables reference each other directly (e.g., `Steps: BuildSteps`)

---

## Library Directory Structure (This Repository)

The wetwire-github-go library itself:

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

## Contracts (contracts.go)

Core interfaces and types (mirroring wetwire-aws-go pattern):

```go
// WorkflowResource represents a GitHub workflow resource.
// All resource types (Workflow, Job) implement this interface.
type WorkflowResource interface {
    ResourceType() string  // e.g., "workflow", "job"
}

// OutputRef represents a reference to a step output.
// When serialized to YAML, becomes: ${{ steps.step_id.outputs.name }}
type OutputRef struct {
    StepID string
    Output string
}

func (o OutputRef) String() string {
    return fmt.Sprintf("${{ steps.%s.outputs.%s }}", o.StepID, o.Output)
}

// DiscoveredWorkflow represents a workflow found by AST parsing.
type DiscoveredWorkflow struct {
    Name         string   // Variable name
    File         string   // Source file path
    Line         int      // Line number
    Jobs         []string // Job variable names in this workflow
}

// DiscoveredJob represents a job found by AST parsing.
type DiscoveredJob struct {
    Name         string   // Variable name
    File         string   // Source file path
    Line         int      // Line number
    Dependencies []string // Referenced job names (Needs field)
}

// Result types for CLI JSON output
type BuildResult struct {
    Success   bool     `json:"success"`
    Workflows []string `json:"workflows,omitempty"`
    Files     []string `json:"files,omitempty"`
    Errors    []string `json:"errors,omitempty"`
}

type LintResult struct {
    Success bool        `json:"success"`
    Issues  []LintIssue `json:"issues,omitempty"`
}

type LintIssue struct {
    File     string `json:"file"`
    Line     int    `json:"line"`
    Column   int    `json:"column"`
    Severity string `json:"severity"`
    Message  string `json:"message"`
    Rule     string `json:"rule"`
    Fixable  bool   `json:"fixable"`
}

type ValidateResult struct {
    Success  bool     `json:"success"`
    Errors   []string `json:"errors,omitempty"`
    Warnings []string `json:"warnings,omitempty"`
}

type ListResult struct {
    Workflows []ListWorkflow `json:"workflows"`
}

type ListWorkflow struct {
    Name string `json:"name"`
    File string `json:"file"`
    Line int    `json:"line"`
    Jobs int    `json:"jobs"`
}
```

## Flat Variable Pattern (Importer-Generated)

The **importer automatically generates correct syntax** based on field types in the schema:

```go
// Type definitions determine pointer vs value
type Triggers struct {
    Push        *PushTrigger        `yaml:"push,omitempty"`  // pointer → importer adds &
    PullRequest *PullRequestTrigger `yaml:"pull_request,omitempty"`
}

type Job struct {
    Steps []Step `yaml:"steps"`  // value type → importer omits &
}
```

**Generated code from importer** (users don't decide `&` — tooling handles it):

```go
// Pointer field (*PushTrigger) → importer generates &
var CIPush = workflow.PushTrigger{
    Branches: List("main"),
}

// Value field ([]Step) → importer generates without &
var BuildSteps = []workflow.Step{
    checkout.Checkout{}.ToStep(),
}

// Clean references at usage site
var CI = workflow.Workflow{
    On: workflow.Triggers{
        Push: CIPush,  // Just reference the variable
    },
}

var Build = workflow.Job{
    Steps: BuildSteps,
}
```

**Key insight:** Users never manually add `&`. The importer inspects field types and generates:
- `var X = &pkg.Type{...}` for pointer fields (`*Type`)
- `var X = pkg.Type{...}` for value fields (`Type`)

**For `*bool` fields** — importer uses the `Ptr` helper:
```go
var BuildStrategy = workflow.Strategy{
    FailFast: Ptr(false),  // *bool field
}
```

---

## Core Types (workflow/ package)

Types designed for the "no parens" pattern — struct literal initialization:

```go
// workflow.go
type Workflow struct {
    Name        string                 `yaml:"name,omitempty"`
    On          Triggers               `yaml:"on"`
    Env         map[string]any         `yaml:"env,omitempty"`
    Defaults    *Defaults              `yaml:"defaults,omitempty"`
    Concurrency *Concurrency           `yaml:"concurrency,omitempty"`
    Permissions *Permissions           `yaml:"permissions,omitempty"`
    Jobs        map[string]Job         `yaml:"jobs"` // populated by discovery
}

func (w Workflow) ResourceType() string { return "workflow" }

// job.go
type Job struct {
    // Outputs for cross-job references (like AttrRef in wetwire-aws-go)
    Outputs map[string]OutputRef `yaml:"-"` // excluded from YAML, used for refs

    Name           string              `yaml:"name,omitempty"`
    RunsOn         any                 `yaml:"runs-on"`
    Needs          []any               `yaml:"needs,omitempty"` // Job references
    If             any                 `yaml:"if,omitempty"`
    Environment    *Environment        `yaml:"environment,omitempty"`
    Concurrency    *Concurrency        `yaml:"concurrency,omitempty"`
    Strategy       *Strategy           `yaml:"strategy,omitempty"`
    Container      *Container          `yaml:"container,omitempty"`
    Services       map[string]Service  `yaml:"services,omitempty"`
    Steps          []Step              `yaml:"steps"`
    TimeoutMinutes int                 `yaml:"timeout-minutes,omitempty"`
}

func (j Job) ResourceType() string { return "job" }

// step.go
type Step struct {
    // Outputs for cross-step references
    ID      string         `yaml:"id,omitempty"`
    Name    string         `yaml:"name,omitempty"`
    If      any            `yaml:"if,omitempty"`
    Uses    string         `yaml:"uses,omitempty"`
    With    map[string]any `yaml:"with,omitempty"`
    Run     string         `yaml:"run,omitempty"`
    Shell   string         `yaml:"shell,omitempty"`
    Env     map[string]any `yaml:"env,omitempty"`
    WorkingDirectory string `yaml:"working-directory,omitempty"`
}

// Output returns an OutputRef for referencing this step's outputs
func (s Step) Output(name string) OutputRef {
    return OutputRef{StepID: s.ID, Output: name}
}

// triggers.go
type Triggers struct {
    Push              *PushTrigger              `yaml:"push,omitempty"`
    PullRequest       *PullRequestTrigger       `yaml:"pull_request,omitempty"`
    PullRequestTarget *PullRequestTargetTrigger `yaml:"pull_request_target,omitempty"`
    Schedule          []ScheduleTrigger         `yaml:"schedule,omitempty"`
    WorkflowDispatch  *WorkflowDispatchTrigger  `yaml:"workflow_dispatch,omitempty"`
    WorkflowCall      *WorkflowCallTrigger      `yaml:"workflow_call,omitempty"`
    // ... 30+ event types
}

// matrix.go
type Strategy struct {
    Matrix      *Matrix `yaml:"matrix,omitempty"`
    FailFast    *bool   `yaml:"fail-fast,omitempty"`
    MaxParallel int     `yaml:"max-parallel,omitempty"`
}

type Matrix struct {
    Values  map[string][]any   `yaml:",inline"`
    Include []map[string]any   `yaml:"include,omitempty"`
    Exclude []map[string]any   `yaml:"exclude,omitempty"`
}
```

## Helper Types (workflow/helpers.go)

Convenience types for cleaner struct literals (from wetwire-aws-go pattern):

```go
// Env is a shorthand for map[string]any.
// Used for environment variable blocks.
//
// Example:
//
//	Env: Env{
//	    "NODE_ENV": "production",
//	    "TOKEN":    Secrets.Get("DEPLOY_TOKEN"),
//	}
type Env = map[string]any

// With is a shorthand for map[string]any.
// Used for action input blocks.
type With = map[string]any

// List creates a typed slice from items.
// Avoids verbose slice type annotations.
//
// Example:
//
//	// Instead of:
//	Steps: []workflow.Step{Step1, Step2, Step3},
//	// Write:
//	Steps: List(Step1, Step2, Step3),
func List[T any](items ...T) []T { return items }

// Strings creates a []string slice.
//
// Example:
//
//	Branches: Strings("main", "develop"),
func Strings(items ...string) []string { return items }

// Ptr returns a pointer to the value.
// Use for *bool and other pointer fields in struct literals.
//
// Example:
//
//	Strategy: &Strategy{FailFast: Ptr(false)},
//	PersistCredentials: Ptr(true),
func Ptr[T any](v T) *T { return &v }
```

---

## Expression Contexts (workflow/expressions.go)

Context accessors for GitHub Actions expressions (like intrinsics in wetwire-aws-go):

```go
// Expression wraps a GitHub Actions expression string
type Expression string

func (e Expression) String() string {
    return fmt.Sprintf("${{ %s }}", string(e))
}

// Context accessors - used like: workflow.Secrets.Get("TOKEN")
var GitHub = githubContext{}
var Runner = runnerContext{}
var Env = envContext{}
var Secrets = secretsContext{}
var Matrix = matrixContext{}
var Steps = stepsContext{}
var Needs = needsContext{}

type secretsContext struct{}
func (secretsContext) Get(name string) Expression {
    return Expression(fmt.Sprintf("secrets.%s", name))
}

type matrixContext struct{}
func (matrixContext) Get(name string) Expression {
    return Expression(fmt.Sprintf("matrix.%s", name))
}

type githubContext struct{}
func (githubContext) Ref() Expression     { return Expression("github.ref") }
func (githubContext) SHA() Expression     { return Expression("github.sha") }
func (githubContext) Actor() Expression   { return Expression("github.actor") }
func (githubContext) EventName() Expression { return Expression("github.event_name") }

// Condition builders
func Always() Expression  { return Expression("always()") }
func Failure() Expression { return Expression("failure()") }
func Success() Expression { return Expression("success()") }
func Cancelled() Expression { return Expression("cancelled()") }

func Branch(name string) Expression {
    return Expression(fmt.Sprintf("github.ref == 'refs/heads/%s'", name))
}

// Composable conditions
func (e Expression) And(other Expression) Expression {
    return Expression(fmt.Sprintf("(%s) && (%s)", e, other))
}

func (e Expression) Or(other Expression) Expression {
    return Expression(fmt.Sprintf("(%s) || (%s)", e, other))
}
```

---

## Generated Action Wrappers (actions/ package)

Action wrappers follow the same pattern as wetwire-aws-go resources:
- Struct fields for inputs (like resource properties)
- Output fields with `yaml:"-"` for cross-step references (like AttrRef)
- `ToStep()` method to convert to workflow.Step

```go
// actions/checkout/checkout.go
// Code generated by wetwire-github codegen. DO NOT EDIT.
// Source: actions/checkout/action.yml
// Generated: 2026-01-04T10:00:00Z

package checkout

import (
    wetwire "github.com/lex00/wetwire-github-go"
    "github.com/lex00/wetwire-github-go/workflow"
)

// Checkout represents actions/checkout@v4
type Checkout struct {
    // Outputs for cross-step references (like AttrRef in wetwire-aws-go)
    RefOutput    wetwire.OutputRef `yaml:"-"`
    CommitOutput wetwire.OutputRef `yaml:"-"`

    // Inputs (from action.yml)
    Repository         string `yaml:"repository,omitempty"`
    Ref                string `yaml:"ref,omitempty"`
    Token              string `yaml:"token,omitempty"`
    SSHKey             string `yaml:"ssh-key,omitempty"`
    PersistCredentials *bool  `yaml:"persist-credentials,omitempty"`
    Path               string `yaml:"path,omitempty"`
    Clean              *bool  `yaml:"clean,omitempty"`
    FetchDepth         int    `yaml:"fetch-depth,omitempty"`
    FetchTags          *bool  `yaml:"fetch-tags,omitempty"`
    Submodules         string `yaml:"submodules,omitempty"`
    LFS                *bool  `yaml:"lfs,omitempty"`
}

// Action returns the action reference string
func (a Checkout) Action() string { return "actions/checkout@v4" }

// ToStep converts to a workflow.Step for use in Steps slice
func (a Checkout) ToStep() workflow.Step {
    return workflow.Step{
        Uses: a.Action(),
        With: a.toWith(),
    }
}

// ToStepWithID converts to a workflow.Step with an ID for output references
func (a Checkout) ToStepWithID(id string) workflow.Step {
    return workflow.Step{
        ID:   id,
        Uses: a.Action(),
        With: a.toWith(),
    }
}
```

**Usage pattern (no parens):**

```go
// Simple - just struct initialization
var CheckoutStep = checkout.Checkout{
    FetchDepth: 0,
    Submodules: "recursive",
}.ToStep()

// With outputs - use ToStepWithID for references
var CheckoutWithRef = checkout.Checkout{
    Ref: "main",
}.ToStepWithID("checkout")

// Reference the output in another step
var NextStep = workflow.Step{
    Run: fmt.Sprintf("echo Checked out %s",
        workflow.Steps.Output("checkout", "ref")),
}
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

See [ROADMAP.md](ROADMAP.md) for detailed feature matrix with 187 features across 29 parallel streams.

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
