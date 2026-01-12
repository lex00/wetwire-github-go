// MCP server implementation for IDE integration.
//
// When design --mcp-server is called, this runs the MCP protocol over stdio,
// providing wetwire_init, wetwire_lint, wetwire_build, and wetwire_validate tools.
//
// This implementation uses wetwire-core-go/mcp infrastructure for the server
// and protocol handling, with GitHub-specific tool handlers.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lex00/wetwire-core-go/mcp"

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
	server := mcp.NewServer(mcp.Config{
		Name:    "wetwire-github",
		Version: getVersion(),
	})

	// Register standard tools with GitHub-specific handlers
	registerTools(server)

	// Run on stdio transport
	return server.Start(context.Background())
}

// Tool schemas for standard wetwire tools
var (
	initSchema = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "Project name",
			},
			"path": map[string]any{
				"type":        "string",
				"description": "Output directory (default: current directory)",
			},
		},
	}

	buildSchema = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to discover resources from",
			},
			"output": map[string]any{
				"type":        "string",
				"description": "Output directory for generated files",
			},
			"dry_run": map[string]any{
				"type":        "boolean",
				"description": "Return content without writing files",
			},
		},
	}

	lintSchema = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"package": map[string]any{
				"type":        "string",
				"description": "Package path to lint",
			},
		},
	}

	validateSchema = map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path": map[string]any{
				"type":        "string",
				"description": "Path to file or directory to validate",
			},
		},
	}
)

// registerTools registers all GitHub-specific MCP tools.
func registerTools(server *mcp.Server) {
	server.RegisterToolWithSchema(
		"wetwire_init",
		"Initialize a new wetwire-github project with example code",
		handleInit,
		initSchema,
	)

	server.RegisterToolWithSchema(
		"wetwire_build",
		"Generate GitHub Actions YAML workflows from wetwire declarations",
		handleBuild,
		buildSchema,
	)

	server.RegisterToolWithSchema(
		"wetwire_lint",
		"Check code quality and style (WAG001-WAG008 lint rules)",
		handleLint,
		lintSchema,
	)

	server.RegisterToolWithSchema(
		"wetwire_validate",
		"Validate generated workflows using actionlint",
		handleValidate,
		validateSchema,
	)
}

// validMCPProjectName matches valid Go module/project names
var validMCPProjectName = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// handleInit implements the wetwire_init tool.
func handleInit(_ context.Context, args map[string]any) (string, error) {
	type InitResult struct {
		Success bool     `json:"success"`
		Path    string   `json:"path"`
		Files   []string `json:"files"`
		Error   string   `json:"error,omitempty"`
	}

	result := InitResult{}

	name, _ := args["name"].(string)
	if name == "" {
		result.Error = "name is required"
		return toJSON(result)
	}

	// Validate project name
	if !validMCPProjectName.MatchString(name) {
		result.Error = fmt.Sprintf("invalid project name %q: must start with a letter and contain only letters, numbers, hyphens, or underscores", name)
		return toJSON(result)
	}

	// Default path to current directory
	workspaceDir, _ := args["path"].(string)
	if workspaceDir == "" {
		workspaceDir = "."
	}

	// Create project as subdirectory of workspace
	projectPath := filepath.Join(workspaceDir, name)
	result.Path = projectPath

	// Check if project already exists
	if _, err := os.Stat(projectPath); err == nil {
		result.Error = fmt.Sprintf("project already exists: %s", projectPath)
		return toJSON(result)
	}

	// Create project directory
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		result.Error = fmt.Sprintf("creating project directory: %v", err)
		return toJSON(result)
	}

	// Use project name as module name
	moduleName := name

	// Write go.mod
	goMod := fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-github-go v0.1.0
`, moduleName)

	goModPath := filepath.Join(projectPath, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		result.Error = fmt.Sprintf("writing go.mod: %v", err)
		return toJSON(result)
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
var BuildSteps = []any{
	checkout.Checkout{},
	setup_go.SetupGo{GoVersion: "1.23"},
	workflow.Step{Name: "Build", Run: "go build ./..."},
	workflow.Step{Name: "Test", Run: "go test ./..."},
}
`
	workflowsGoPath := filepath.Join(projectPath, "workflows.go")
	if err := os.WriteFile(workflowsGoPath, []byte(workflowsGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing workflows.go: %v", err)
		return toJSON(result)
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
		return toJSON(result)
	}
	result.Files = append(result.Files, ".gitignore")

	result.Success = true
	return toJSON(result)
}

// handleLint implements the wetwire_lint tool.
func handleLint(_ context.Context, args map[string]any) (string, error) {
	result := wetwire.LintResult{}

	// Support both "path" and "package" parameter names for compatibility
	path, _ := args["package"].(string)
	if path == "" {
		path, _ = args["path"].(string)
	}

	if path == "" {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "error",
			Message:  "package is required",
			Rule:     "internal",
		})
		return toJSON(result)
	}

	// Discover workflows (validates references)
	disc := discover.NewDiscoverer()
	discoverResult, err := disc.Discover(path)
	if err != nil {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "error",
			Message:  fmt.Sprintf("discovery failed: %v", err),
			Rule:     "internal",
		})
		return toJSON(result)
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
	lintResult, err := l.LintDir(path)
	if err != nil {
		result.Issues = append(result.Issues, wetwire.LintIssue{
			Severity: "warning",
			Message:  fmt.Sprintf("failed to lint %s: %v", path, err),
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
	return toJSON(result)
}

// handleBuild implements the wetwire_build tool.
func handleBuild(_ context.Context, args map[string]any) (string, error) {
	type BuildResult struct {
		Success   bool              `json:"success"`
		Workflows map[string]string `json:"workflows,omitempty"`
		Files     []string          `json:"files,omitempty"`
		Errors    []string          `json:"errors,omitempty"`
	}

	result := BuildResult{
		Workflows: make(map[string]string),
	}

	// Support both "path" and "package" parameter names for compatibility
	path, _ := args["package"].(string)
	if path == "" {
		path, _ = args["path"].(string)
	}

	if path == "" {
		result.Errors = append(result.Errors, "package is required")
		return toJSON(result)
	}

	outputDir, _ := args["output"].(string)
	if outputDir == "" {
		outputDir = ".github/workflows"
	}

	dryRun, _ := args["dry_run"].(bool)

	// Discover workflows
	disc := discover.NewDiscoverer()
	discoverResult, err := disc.Discover(path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return toJSON(result)
	}

	// Check for discovery errors
	if len(discoverResult.Errors) > 0 {
		for _, e := range discoverResult.Errors {
			result.Errors = append(result.Errors, e)
		}
		return toJSON(result)
	}

	// Check we have workflows to build
	if len(discoverResult.Workflows) == 0 {
		result.Errors = append(result.Errors, "no workflows found")
		return toJSON(result)
	}

	// Extract values using runner
	r := runner.NewRunner()
	extracted, err := r.ExtractValues(path, discoverResult)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("extracting values: %v", err))
		return toJSON(result)
	}

	if extracted.Error != "" {
		result.Errors = append(result.Errors, extracted.Error)
		return toJSON(result)
	}

	// Build templates
	builder := template.NewBuilder()
	built, err := builder.Build(discoverResult, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return toJSON(result)
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Store YAML content for each workflow
	for _, wf := range built.Workflows {
		result.Workflows[wf.Name] = string(wf.YAML)
	}

	// Write files if not dry run
	if !dryRun {
		// Create output directory
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return toJSON(result)
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
	return toJSON(result)
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

// handleValidate implements the wetwire_validate tool.
func handleValidate(_ context.Context, args map[string]any) (string, error) {
	type ValidateResult struct {
		Success bool     `json:"success"`
		Errors  []string `json:"errors,omitempty"`
	}

	result := ValidateResult{}

	path, _ := args["path"].(string)
	if path == "" {
		result.Errors = append(result.Errors, "path is required")
		return toJSON(result)
	}

	// Validate the workflow file
	validationResult, err := validation.ValidateWorkflowFile(path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("validation failed: %v", err))
		return toJSON(result)
	}

	for _, issue := range validationResult.Issues {
		result.Errors = append(result.Errors, fmt.Sprintf("%s:%d:%d: %s", issue.File, issue.Line, issue.Column, issue.Message))
	}

	result.Success = len(result.Errors) == 0
	return toJSON(result)
}

// toJSON marshals the given value to a JSON string.
func toJSON(v any) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshaling result: %w", err)
	}
	return string(data), nil
}
