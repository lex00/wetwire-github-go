// MCP server implementation for IDE integration.
//
// When design --mcp-server is called, this runs the MCP protocol over stdio,
// providing wetwire_init, wetwire_lint, wetwire_build, and wetwire_validate tools.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/linter"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/template"
	"github.com/lex00/wetwire-github-go/internal/validation"
)

// runMCPServer starts the MCP server on stdio transport.
// This is called when design --mcp-server is invoked.
func runMCPServer() error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "wetwire-github",
		Version: getVersion(),
	}, nil)

	// Register tools
	registerInitTool(server)
	registerLintTool(server)
	registerBuildTool(server)
	registerValidateTool(server)

	// Run on stdio transport
	return server.Run(context.Background(), &mcp.StdioTransport{})
}

// InitArgs are the arguments for the wetwire_init tool.
type InitArgs struct {
	Name string `json:"name" jsonschema:"required,Project name (e.g. my-workflows, ci-project)"`
	Path string `json:"path,omitempty" jsonschema:"Workspace directory to create project in (defaults to current directory)"`
}

// InitResult is the result of the wetwire_init tool.
type InitResult struct {
	Success bool     `json:"success"`
	Path    string   `json:"path"`
	Files   []string `json:"files"`
	Error   string   `json:"error,omitempty"`
}

// validProjectName matches valid Go module/project names
var validMCPProjectName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

func registerInitTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wetwire_init",
		Description: "Initialize a new wetwire-github project in a subdirectory. Creates {path}/{name}/ with go.mod and workflow declarations.",
	}, handleInit)
}

func handleInit(_ context.Context, _ *mcp.CallToolRequest, args InitArgs) (*mcp.CallToolResult, any, error) {
	result := InitResult{}

	if args.Name == "" {
		result.Error = "name is required"
		return toolResult(result)
	}

	// Validate project name
	if !validMCPProjectName.MatchString(args.Name) {
		result.Error = fmt.Sprintf("invalid project name %q: must start with a letter and contain only letters, numbers, hyphens, or underscores", args.Name)
		return toolResult(result)
	}

	// Default path to current directory
	workspaceDir := args.Path
	if workspaceDir == "" {
		workspaceDir = "."
	}

	// Create project as subdirectory of workspace
	projectPath := filepath.Join(workspaceDir, args.Name)
	result.Path = projectPath

	// Check if project already exists
	if _, err := os.Stat(projectPath); err == nil {
		result.Error = fmt.Sprintf("project already exists: %s", projectPath)
		return toolResult(result)
	}

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		result.Error = fmt.Sprintf("creating project directory: %v", err)
		return toolResult(result)
	}

	// Use project name as module name
	moduleName := args.Name

	// Write go.mod
	goMod := fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-github-go v0.1.0
`, moduleName)

	goModPath := filepath.Join(projectPath, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		result.Error = fmt.Sprintf("writing go.mod: %v", err)
		return toolResult(result)
	}
	result.Files = append(result.Files, "go.mod")

	// Write workflows.go
	workflowsGo := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/setup_go"
)

// CI is the main continuous integration workflow
var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
	Jobs: map[string]workflow.Job{
		"build": Build,
	},
}

// CITriggers defines when the CI workflow runs
var CITriggers = workflow.Triggers{
	Push:        workflow.PushTrigger{Branches: workflow.List("main")},
	PullRequest: workflow.PullRequestTrigger{Branches: workflow.List("main")},
}

// Build is the main build job
var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}

// BuildSteps are the steps for the build job
var BuildSteps = workflow.List(
	checkout.Checkout{}.ToStep(),
	setup_go.SetupGo{GoVersion: "1.23"}.ToStep(),
	workflow.Step{Name: "Build", Run: "go build ./..."},
	workflow.Step{Name: "Test", Run: "go test ./..."},
)
`
	workflowsGoPath := filepath.Join(projectPath, "workflows.go")
	if err := os.WriteFile(workflowsGoPath, []byte(workflowsGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing workflows.go: %v", err)
		return toolResult(result)
	}
	result.Files = append(result.Files, "workflows.go")

	// Write .gitignore
	gitignore := `# Build output
.github/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
`
	gitignorePath := filepath.Join(projectPath, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignore), 0644); err != nil {
		result.Error = fmt.Sprintf("writing .gitignore: %v", err)
		return toolResult(result)
	}
	result.Files = append(result.Files, ".gitignore")

	result.Success = true
	return toolResult(result)
}

// LintArgs are the arguments for the wetwire_lint tool.
type LintArgs struct {
	Path string `json:"path" jsonschema:"Path to the Go package to lint (e.g. . or ./workflows)"`
}

func registerLintTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wetwire_lint",
		Description: "Lint Go packages for wetwire-github style issues (WAG001-WAG008 rules)",
	}, handleLint)
}

func handleLint(_ context.Context, _ *mcp.CallToolRequest, args LintArgs) (*mcp.CallToolResult, any, error) {
	result := wetwire.LintResult{}

	if args.Path == "" {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "error",
			Message:  "path is required",
			Rule:     "internal",
		})
		return toolResult(result)
	}

	// Discover workflows (validates references)
	disc := discover.NewDiscoverer()
	discoverResult, err := disc.Discover(args.Path)
	if err != nil {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "error",
			Message:  fmt.Sprintf("discovery failed: %v", err),
			Rule:     "internal",
		})
		return toolResult(result)
	}

	// Convert discovery errors to lint issues
	for _, e := range discoverResult.Errors {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "error",
			Message:  e,
			Rule:     "discovery",
		})
	}

	// Run lint rules
	l := linter.DefaultLinter()
	lintResult, err := l.LintDir(args.Path)
	if err != nil {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "warning",
			Message:  fmt.Sprintf("failed to lint %s: %v", args.Path, err),
			Rule:     "internal",
		})
	} else {
		for _, issue := range lintResult.Issues {
			result.Issues = append(result.Issues, wetwire.LintIssue{
				Severity: issue.Severity,
				Message:  issue.Message,
				Rule:     issue.Rule,
				File:     issue.File,
				Line:     issue.Line,
				Column:   issue.Column,
			})
		}
	}

	result.Success = len(result.Issues) == 0
	return toolResult(result)
}

// BuildArgs are the arguments for the wetwire_build tool.
type BuildArgs struct {
	Path   string `json:"path" jsonschema:"Path to the Go package to build (e.g. . or ./workflows)"`
	Output string `json:"output,omitempty" jsonschema:"Output directory for YAML files (default: .github/workflows)"`
	DryRun bool   `json:"dry_run,omitempty" jsonschema:"If true, return YAML content without writing files"`
}

// BuildResult is the result of the wetwire_build tool.
type BuildResult struct {
	Success   bool              `json:"success"`
	Workflows map[string]string `json:"workflows,omitempty"` // name -> YAML content
	Files     []string          `json:"files,omitempty"`     // written file paths
	Errors    []string          `json:"errors,omitempty"`
}

func registerBuildTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wetwire_build",
		Description: "Generate GitHub Actions YAML workflows from Go packages containing wetwire-github declarations",
	}, handleBuild)
}

func handleBuild(_ context.Context, _ *mcp.CallToolRequest, args BuildArgs) (*mcp.CallToolResult, any, error) {
	result := BuildResult{
		Workflows: make(map[string]string),
	}

	if args.Path == "" {
		result.Errors = append(result.Errors, "path is required")
		return toolResult(result)
	}

	outputDir := args.Output
	if outputDir == "" {
		outputDir = ".github/workflows"
	}

	// Discover workflows
	disc := discover.NewDiscoverer()
	discoverResult, err := disc.Discover(args.Path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return toolResult(result)
	}

	// Check for discovery errors
	if len(discoverResult.Errors) > 0 {
		for _, e := range discoverResult.Errors {
			result.Errors = append(result.Errors, e)
		}
		return toolResult(result)
	}

	// Check we have workflows to build
	if len(discoverResult.Workflows) == 0 {
		result.Errors = append(result.Errors, "no workflows found")
		return toolResult(result)
	}

	// Extract values using runner
	r := runner.NewRunner()
	extracted, err := r.ExtractValues(args.Path, discoverResult)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("extracting values: %v", err))
		return toolResult(result)
	}

	if extracted.Error != "" {
		result.Errors = append(result.Errors, extracted.Error)
		return toolResult(result)
	}

	// Build templates
	builder := template.NewBuilder()
	built, err := builder.Build(discoverResult, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return toolResult(result)
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Store YAML content for each workflow
	for _, wf := range built.Workflows {
		result.Workflows[wf.Name] = string(wf.YAML)
	}

	// Write files if not dry run
	if !args.DryRun {
		// Create output directory
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return toolResult(result)
		}

		for name := range result.Workflows {
			filename := filepath.Join(outputDir, toMCPFilename(name)+".yml")
			if err := os.WriteFile(filename, []byte(result.Workflows[name]), 0644); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("writing %s: %v", filename, err))
				continue
			}
			result.Files = append(result.Files, filename)
		}
	}

	result.Success = len(result.Errors) == 0
	return toolResult(result)
}

// toMCPFilename converts a workflow name to a valid filename.
func toMCPFilename(name string) string {
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// ValidateArgs are the arguments for the wetwire_validate tool.
type ValidateArgs struct {
	Path string `json:"path" jsonschema:"Path to YAML file or directory to validate"`
}

// ValidateResult is the result of the wetwire_validate tool.
type ValidateResult struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors,omitempty"`
}

func registerValidateTool(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "wetwire_validate",
		Description: "Validate GitHub Actions YAML workflows using actionlint",
	}, handleValidate)
}

func handleValidate(_ context.Context, _ *mcp.CallToolRequest, args ValidateArgs) (*mcp.CallToolResult, any, error) {
	result := ValidateResult{}

	if args.Path == "" {
		result.Errors = append(result.Errors, "path is required")
		return toolResult(result)
	}

	// Validate the workflow file
	validationResult, err := validation.ValidateWorkflowFile(args.Path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("validation failed: %v", err))
		return toolResult(result)
	}

	for _, issue := range validationResult.Issues {
		result.Errors = append(result.Errors, fmt.Sprintf("%s:%d:%d: %s", issue.File, issue.Line, issue.Column, issue.Message))
	}

	result.Success = len(result.Errors) == 0
	return toolResult(result)
}

// toolResult creates an MCP CallToolResult from any JSON-serializable value.
func toolResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, nil, fmt.Errorf("marshaling result: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}
