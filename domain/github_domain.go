// Package domain provides the GitHubDomain implementation for wetwire-core-go.
package domain

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/linter"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/template"
	"github.com/lex00/wetwire-github-go/internal/validation"
	"github.com/spf13/cobra"
)

// Version is set at build time
var Version = "dev"

// Re-export core types for convenience
type (
	Context      = coredomain.Context
	BuildOpts    = coredomain.BuildOpts
	LintOpts     = coredomain.LintOpts
	InitOpts     = coredomain.InitOpts
	ValidateOpts = coredomain.ValidateOpts
	ListOpts     = coredomain.ListOpts
	GraphOpts    = coredomain.GraphOpts
	Result       = coredomain.Result
	Error        = coredomain.Error
)

var (
	NewResult              = coredomain.NewResult
	NewResultWithData      = coredomain.NewResultWithData
	NewErrorResult         = coredomain.NewErrorResult
	NewErrorResultMultiple = coredomain.NewErrorResultMultiple
)

// GitHubDomain implements the Domain interface for GitHub Actions workflows.
type GitHubDomain struct{}

// Compile-time checks
var (
	_ coredomain.Domain        = (*GitHubDomain)(nil)
	_ coredomain.ListerDomain  = (*GitHubDomain)(nil)
	_ coredomain.GrapherDomain = (*GitHubDomain)(nil)
)

// Name returns "github"
func (d *GitHubDomain) Name() string {
	return "github"
}

// Version returns the current version
func (d *GitHubDomain) Version() string {
	return Version
}

// Builder returns the GitHub builder implementation
func (d *GitHubDomain) Builder() coredomain.Builder {
	return &githubBuilder{}
}

// Linter returns the GitHub linter implementation
func (d *GitHubDomain) Linter() coredomain.Linter {
	return &githubLinter{}
}

// Initializer returns the GitHub initializer implementation
func (d *GitHubDomain) Initializer() coredomain.Initializer {
	return &githubInitializer{}
}

// Validator returns the GitHub validator implementation
func (d *GitHubDomain) Validator() coredomain.Validator {
	return &githubValidator{}
}

// Lister returns the GitHub lister implementation
func (d *GitHubDomain) Lister() coredomain.Lister {
	return &githubLister{}
}

// Grapher returns the GitHub grapher implementation
func (d *GitHubDomain) Grapher() coredomain.Grapher {
	return &githubGrapher{}
}

// CreateRootCommand creates the root command using the domain interface.
func CreateRootCommand(d coredomain.Domain) *cobra.Command {
	return coredomain.Run(d)
}

// githubBuilder implements domain.Builder
type githubBuilder struct{}

func (b *githubBuilder) Build(ctx *Context, path string, opts BuildOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Default output directory for GitHub workflows
	outputDir := opts.Output
	if outputDir == "" {
		outputDir = ".github/workflows"
	}

	// Discover workflows and jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	if len(discovered.Workflows) == 0 {
		return NewErrorResult("no workflows found", Error{
			Path:    absPath,
			Message: "no workflows found",
		}), nil
	}

	// Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractValues(absPath, discovered)
	if err != nil {
		return nil, fmt.Errorf("extraction failed: %w", err)
	}

	if extracted.Error != "" {
		return NewErrorResult("extraction failed", Error{
			Path:    absPath,
			Message: extracted.Error,
		}), nil
	}

	// Build templates
	builder := template.NewBuilder()
	built, err := builder.Build(discovered, extracted)
	if err != nil {
		return nil, fmt.Errorf("template build failed: %w", err)
	}

	if len(built.Errors) > 0 {
		errs := make([]Error, len(built.Errors))
		for i, e := range built.Errors {
			errs[i] = Error{
				Path:    absPath,
				Message: e,
			}
		}
		return NewErrorResultMultiple("template build failed", errs), nil
	}

	// Resolve output directory
	absOutputDir := outputDir
	if !filepath.IsAbs(outputDir) {
		absOutputDir = filepath.Join(absPath, outputDir)
	}

	// Create output directory if needed
	if !opts.DryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			return nil, fmt.Errorf("creating output directory: %w", err)
		}
	}

	// Write workflow files
	var files []string
	for _, wf := range built.Workflows {
		filename := toFilename(wf.Name) + ".yml"
		filePath := filepath.Join(absOutputDir, filename)

		if !opts.DryRun {
			if err := os.WriteFile(filePath, wf.YAML, 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", filename, err)
			}
		}
		files = append(files, filePath)
	}

	return NewResult(fmt.Sprintf("Built %d workflow(s) to %s", len(files), absOutputDir)), nil
}

// githubLinter implements domain.Linter
type githubLinter struct{}

func (l *githubLinter) Lint(ctx *Context, path string, opts LintOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("accessing path: %w", err)
	}

	// Create linter with default rules
	lntr := linter.DefaultLinter()

	var lintResult *linter.LintResult
	if info.IsDir() {
		lintResult, err = lntr.LintDir(absPath)
	} else {
		lintResult, err = lntr.LintFile(absPath)
	}

	if err != nil {
		return nil, fmt.Errorf("linting failed: %w", err)
	}

	if len(lintResult.Issues) == 0 {
		return NewResult("No lint issues found"), nil
	}

	// Convert to domain errors
	errs := make([]Error, 0, len(lintResult.Issues))
	for _, issue := range lintResult.Issues {
		errs = append(errs, Error{
			Path:     issue.File,
			Line:     issue.Line,
			Severity: issue.Severity.String(),
			Message:  issue.Message,
			Code:     issue.Rule,
		})
	}

	return NewErrorResultMultiple("lint issues found", errs), nil
}

// githubInitializer implements domain.Initializer
type githubInitializer struct{}

func (i *githubInitializer) Init(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Use opts.Path if provided, otherwise fall back to path argument
	targetPath := opts.Path
	if targetPath == "" || targetPath == "." {
		targetPath = path
	}

	// Handle scenario initialization
	if opts.Scenario {
		return i.initScenario(ctx, targetPath, opts)
	}

	// Basic project initialization
	return i.initProject(ctx, targetPath, opts)
}

// initScenario creates a full scenario structure with prompts and expected outputs
func (i *githubInitializer) initScenario(ctx *Context, path string, opts InitOpts) (*Result, error) {
	name := opts.Name
	if name == "" {
		name = filepath.Base(path)
	}

	description := opts.Description
	if description == "" {
		description = "GitHub Actions workflow scenario"
	}

	// Use core's scenario scaffolding
	scenario := coredomain.ScaffoldScenario(name, description, "github")
	created, err := coredomain.WriteScenario(path, scenario)
	if err != nil {
		return nil, fmt.Errorf("write scenario: %w", err)
	}

	// Create github-specific expected directories
	expectedDirs := []string{
		filepath.Join(path, "expected", "workflows"),
	}
	for _, dir := range expectedDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	// Create example workflow in expected/workflows/
	exampleWorkflow := `package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// CI workflow runs on push and pull requests to main
var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

var CITriggers = workflow.Triggers{
	Push: &workflow.PushTrigger{
		Branches: []string{"main"},
	},
	PullRequest: &workflow.PullRequestTrigger{
		Branches: []string{"main"},
	},
}

// Build job compiles and tests the code
var Build = workflow.Job{
	RunsOn: "ubuntu-latest",
	Steps: []workflow.Step{
		{Uses: "actions/checkout@v4"},
		{
			Uses: "actions/setup-go@v5",
			With: map[string]any{"go-version": "1.24"},
		},
		{Run: "go build ./..."},
		{Run: "go test ./..."},
	},
}
`
	workflowPath := filepath.Join(path, "expected", "workflows", "workflows.go")
	if err := os.WriteFile(workflowPath, []byte(exampleWorkflow), 0644); err != nil {
		return nil, fmt.Errorf("write example workflow: %w", err)
	}
	created = append(created, "expected/workflows/workflows.go")

	return NewResultWithData(
		fmt.Sprintf("Created scenario %s with %d files", name, len(created)),
		created,
	), nil
}

// initProject creates a basic project with example workflows
func (i *githubInitializer) initProject(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Create directory structure
	dirs := []string{
		path,
		filepath.Join(path, "workflows"),
		filepath.Join(path, "cmd"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	// Create go.mod
	name := filepath.Base(path)
	modulePath := "github.com/example/" + name
	goMod := fmt.Sprintf(`module %s

go 1.24

require github.com/lex00/wetwire-github-go v0.0.0
`, modulePath)
	goModPath := filepath.Join(path, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		return nil, fmt.Errorf("writing go.mod: %w", err)
	}

	// Create workflows/workflows.go
	workflowsGo := `package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// CI workflow runs on push and pull requests to main
var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}
`
	workflowsPath := filepath.Join(path, "workflows", "workflows.go")
	if err := os.WriteFile(workflowsPath, []byte(workflowsGo), 0644); err != nil {
		return nil, fmt.Errorf("writing workflows.go: %w", err)
	}

	// Create workflows/triggers.go
	triggersGo := `package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CIPush = workflow.PushTrigger{
	Branches: []string{"main"},
}

var CIPullRequest = workflow.PullRequestTrigger{
	Branches: []string{"main"},
}

var CITriggers = workflow.Triggers{
	Push:        &CIPush,
	PullRequest: &CIPullRequest,
}
`
	triggersPath := filepath.Join(path, "workflows", "triggers.go")
	if err := os.WriteFile(triggersPath, []byte(triggersGo), 0644); err != nil {
		return nil, fmt.Errorf("writing triggers.go: %w", err)
	}

	// Create workflows/jobs.go
	jobsGo := `package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// Build job compiles and tests the code
var Build = workflow.Job{
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}
`
	jobsPath := filepath.Join(path, "workflows", "jobs.go")
	if err := os.WriteFile(jobsPath, []byte(jobsGo), 0644); err != nil {
		return nil, fmt.Errorf("writing jobs.go: %w", err)
	}

	// Create workflows/steps.go
	stepsGo := `package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BuildSteps = []workflow.Step{
	{Uses: "actions/checkout@v4"},
	{
		Uses: "actions/setup-go@v5",
		With: map[string]any{"go-version": "1.24"},
	},
	{Run: "go build ./..."},
	{Run: "go test ./..."},
}
`
	stepsPath := filepath.Join(path, "workflows", "steps.go")
	if err := os.WriteFile(stepsPath, []byte(stepsGo), 0644); err != nil {
		return nil, fmt.Errorf("writing steps.go: %w", err)
	}

	return NewResult(fmt.Sprintf("Created project: %s", path)), nil
}

// githubValidator implements domain.Validator
type githubValidator struct{}

func (v *githubValidator) Validate(ctx *Context, path string, opts ValidateOpts) (*Result, error) {
	// Validate using actionlint
	validator := validation.NewActionlintValidator()
	validationResult, err := validator.ValidateFile(path)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if len(validationResult.Issues) == 0 {
		return NewResult("Validation passed"), nil
	}

	// Convert validation issues to errors
	errs := make([]Error, 0, len(validationResult.Issues))
	for _, issue := range validationResult.Issues {
		errs = append(errs, Error{
			Path:     issue.File,
			Line:     issue.Line,
			Severity: "error",
			Message:  issue.Message,
			Code:     issue.RuleID,
		})
	}

	return NewErrorResultMultiple("validation failed", errs), nil
}

// githubLister implements domain.Lister
type githubLister struct{}

func (l *githubLister) List(ctx *Context, path string, opts ListOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover workflows and jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Build list
	list := make([]map[string]any, 0)
	for _, wf := range discovered.Workflows {
		relPath := wf.File
		if rel, err := filepath.Rel(absPath, wf.File); err == nil {
			relPath = rel
		}

		list = append(list, map[string]any{
			"name": wf.Name,
			"type": "workflow",
			"file": relPath,
			"line": wf.Line,
			"jobs": len(wf.Jobs),
		})
	}

	for _, job := range discovered.Jobs {
		relPath := job.File
		if rel, err := filepath.Rel(absPath, job.File); err == nil {
			relPath = rel
		}

		list = append(list, map[string]any{
			"name": job.Name,
			"type": "job",
			"file": relPath,
			"line": job.Line,
		})
	}

	return NewResultWithData(fmt.Sprintf("Discovered %d resources", len(list)), list), nil
}

// githubGrapher implements domain.Grapher
type githubGrapher struct{}

func (g *githubGrapher) Graph(ctx *Context, path string, opts GraphOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Build graph
	graph := discover.NewDependencyGraph(discovered.Jobs)

	// Generate output based on format
	var output string
	switch opts.Format {
	case "dot", "":
		output = generateDOT(graph)
	case "mermaid":
		output = generateMermaid(graph)
	default:
		return nil, fmt.Errorf("unknown format: %s", opts.Format)
	}

	return NewResultWithData("Graph generated", output), nil
}

// Helper functions

// toFilename converts a workflow name to a valid filename.
// "CI" -> "ci", "MyWorkflow" -> "my-workflow"
func toFilename(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// generateDOT generates DOT format output.
func generateDOT(graph *discover.DependencyGraph) string {
	var sb strings.Builder

	sb.WriteString("digraph workflow {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=box];\n\n")

	// Get sorted node names
	nodeNames := make([]string, 0, len(graph.Nodes))
	for name := range graph.Nodes {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	// Write nodes
	for _, name := range nodeNames {
		sb.WriteString(fmt.Sprintf("  %q;\n", name))
	}

	sb.WriteString("\n")

	// Write edges
	for _, name := range nodeNames {
		deps := graph.Edges[name]
		if len(deps) > 0 {
			sortedDeps := make([]string, len(deps))
			copy(sortedDeps, deps)
			sort.Strings(sortedDeps)
			for _, dep := range sortedDeps {
				sb.WriteString(fmt.Sprintf("  %q -> %q;\n", dep, name))
			}
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// generateMermaid generates Mermaid diagram format output.
func generateMermaid(graph *discover.DependencyGraph) string {
	var sb strings.Builder

	sb.WriteString("graph TB\n")

	// Get sorted node names
	nodeNames := make([]string, 0, len(graph.Nodes))
	for name := range graph.Nodes {
		nodeNames = append(nodeNames, name)
	}
	sort.Strings(nodeNames)

	// Write edges
	hasEdges := false
	for _, name := range nodeNames {
		deps := graph.Edges[name]
		if len(deps) > 0 {
			sortedDeps := make([]string, len(deps))
			copy(sortedDeps, deps)
			sort.Strings(sortedDeps)
			for _, dep := range sortedDeps {
				sb.WriteString(fmt.Sprintf("    %s --> %s\n", dep, name))
				hasEdges = true
			}
		}
	}

	// If no edges, just list nodes
	if !hasEdges {
		for _, name := range nodeNames {
			sb.WriteString(fmt.Sprintf("    %s\n", name))
		}
	}

	return sb.String()
}
