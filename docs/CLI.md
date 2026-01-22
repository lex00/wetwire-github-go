<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

Complete reference for the `wetwire-github` command-line interface.

## Installation

See [README.md](../README.md#installation) for installation instructions.

## Commands

### `wetwire-github init`

Create a new workflow project.

```bash
wetwire-github init <name> [flags]
```

**Flags:**
- `-o, --output <dir>` — Output directory (default: current directory)
- `--format <format>` — Output format: `text` or `json` (default: `text`)

**Example:**
```bash
wetwire-github init my-workflows
wetwire-github init my-workflows -o ./projects/
wetwire-github init my-workflows --format json
```

### `wetwire-github build`

Generate YAML from Go workflow declarations.

```bash
wetwire-github build <path> [flags]
```

**Flags:**
- `-o, --output <dir>` — Output directory (default: `.github/workflows/`)
- `--format <format>` — Output format: `yaml` or `json` (default: `yaml`)
- `--type <type>` — Config type: `workflow`, `dependabot`, `issue-template`, `discussion-template` (default: `workflow`)
- `--dry-run` — Show what would be written without writing

**Example:**
```bash
wetwire-github build .
wetwire-github build ./my-workflows -o ./output/
wetwire-github build . --type dependabot
```

### `wetwire-github import`

Convert existing configuration files to Go code.

```bash
wetwire-github import <file> [flags]
```

**Flags:**
- `-o, --output <dir>` — Output directory (default: current directory)
- `--single-file` — Generate all code in a single file
- `--no-scaffold` — Skip generating go.mod, README, etc.
- `--type <type>` — Config type: `workflow`, `dependabot`, `issue-template`, `discussion-template`, `codeowners`

**Supported Types:**
- `workflow` — GitHub Actions YAML files (default)
- `dependabot` — Dependabot configuration
- `issue-template` — Issue template YAML files
- `discussion-template` — Discussion template YAML files
- `codeowners` — CODEOWNERS file

**Example:**
```bash
wetwire-github import .github/workflows/ci.yml -o my-workflows/
wetwire-github import ci.yml --single-file
wetwire-github import .github/CODEOWNERS --type codeowners -o my-project/
```

See [Import Workflow](IMPORT_WORKFLOW.md) for detailed import documentation.

### `wetwire-github validate`

Validate YAML using actionlint.

```bash
wetwire-github validate <workflow.yml> [flags]
```

**Flags:**
- `--format <format>` — Output format: `text` or `json` (default: `text`)

**Example:**
```bash
wetwire-github validate .github/workflows/ci.yml
wetwire-github validate ci.yml --format json
```

### `wetwire-github lint`

Check Go code for wetwire best practices.

```bash
wetwire-github lint <path> [flags]
```

**Flags:**
- `--format <format>` — Output format: `text` or `json` (default: `text`)
- `--fix` — Automatically fix issues where possible

**Example:**
```bash
wetwire-github lint .
wetwire-github lint ./my-workflows --format json
wetwire-github lint . --fix
```

### `wetwire-github list`

List discovered workflows and jobs.

```bash
wetwire-github list <path> [flags]
```

**Flags:**
- `--format <format>` — Output format: `text` or `json` (default: `text`)

**Example:**
```bash
wetwire-github list .
wetwire-github list ./my-workflows --format json
```

### `wetwire-github graph`

Generate DAG visualization of workflow jobs.

```bash
wetwire-github graph <path> [flags]
```

**Flags:**
- `--format <format>` — Output format: `dot` or `mermaid` (default: `mermaid`)
- `-o, --output <file>` — Output file (default: stdout)

**Example:**
```bash
wetwire-github graph . --format mermaid
wetwire-github graph . --format dot -o workflow.dot
```

### `wetwire-github design`

AI-assisted workflow design (requires wetwire-core-go).

```bash
wetwire-github design [prompt] [flags]
```

**Flags:**
- `--provider <provider>` — LLM provider: `anthropic` or `kiro` (default: anthropic)
- `--stream` — Stream output tokens
- `--max-lint-cycles <n>` — Maximum lint/fix cycles (default: 5)
- `--model <model>` — Claude model to use (default: claude-sonnet-4-20250514)
- `-w, --workdir <dir>` — Working directory for generated files (default: .)

**MCP Server Mode:**

To run as an MCP server for IDE integration (e.g., Kiro, Claude Desktop):

```bash
wetwire-github design --mcp-server
```

This starts the MCP protocol over stdio with the following tools:
- `wetwire_init` — Initialize a new wetwire-github project
- `wetwire_lint` — Lint Go packages for WAG001-WAG008 rules
- `wetwire_build` — Generate YAML workflows from Go declarations
- `wetwire_validate` — Validate YAML using actionlint

**Example:**
```bash
wetwire-github design "Create a CI workflow for a Go project"
wetwire-github design --stream "Add matrix testing for Go 1.22 and 1.23"
wetwire-github design --provider kiro "Create a deployment workflow"
```

### `wetwire-github test`

Run persona-based tests against workflow declarations.

```bash
wetwire-github test <path> [flags]
```

**Flags:**
- `--provider <provider>` — LLM provider: `anthropic` or `kiro` (default: anthropic)
- `--format <format>` — Output format: `text` or `json` (default: `text`)
- `--persona <name>` — Run specific persona: `beginner`, `intermediate`, `expert`
- `--scenario <name>` — Run specific scenario: `ci-workflow`, `deployment`, `release`, `matrix`
- `--score` — Show scoring breakdown
- `--list` — List available personas and scenarios

**Developer Personas:**
- `beginner` — New to GitHub Actions, needs guidance
- `intermediate` — Some experience, knows basics but misses details
- `expert` — Deep CI/CD knowledge, precise requirements

Custom personas can be registered for domain-specific testing.

**Scoring Dimensions (0-3 each):**
- Completeness — Were all required workflows generated?
- Lint Quality — Did the code pass linting?
- Code Quality — Does the code follow idiomatic patterns?
- Output Validity — Is the generated YAML valid?
- Question Efficiency — Appropriate number of clarifying questions?

**Thresholds:** 0-5 Failure, 6-9 Partial, 10-12 Success, 13-15 Excellent

**Example:**
```bash
wetwire-github test .
wetwire-github test . --persona beginner --score
wetwire-github test . --scenario ci-workflow --format json
wetwire-github test . --provider kiro --persona beginner
wetwire-github test --list
```

### `wetwire-github version`

Print version information.

```bash
wetwire-github version
```

## Exit Codes

| Command | Exit 0 | Exit 1 | Exit 2 |
|---------|--------|--------|--------|
| `build` | Success | Error (parse, generation) | — |
| `lint` | No issues | Issues found | Error (parse failure) |
| `import` | Success | Error (parse, generation) | — |
| `validate` | Valid | Invalid (actionlint errors) | Error (file not found) |
| `list` | Success | Error | — |
| `init` | Success | Error (dir exists, write fail) | — |

## Environment Variables

| Variable | Description |
|----------|-------------|
| `WETWIRE_OUTPUT_DIR` | Default output directory |
| `WETWIRE_FORMAT` | Default output format |

## Configuration

wetwire-github reads configuration from:

1. Command-line flags (highest priority)
2. Environment variables
3. `.wetwire.yaml` in project root (if present)

Example `.wetwire.yaml`:

```yaml
output: .github/workflows
format: yaml
lint:
  rules:
    - WAG001
    - WAG002
    - WAG003
```
