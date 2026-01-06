package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
)

var initOutput string
var initFormat string

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Create a new workflow project",
	Long: `Init creates a new wetwire-github workflow project with example declarations.

Example:
  wetwire-github init my-workflows
  wetwire-github init my-ci -o ./projects/
  wetwire-github init my-ci --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return runInit(name)
	},
}

func init() {
	initCmd.Flags().StringVarP(&initOutput, "output", "o", ".", "output directory")
	initCmd.Flags().StringVar(&initFormat, "format", "text", "output format (text, json)")
}

// runInit creates a new workflow project.
func runInit(name string) error {
	result := wetwire.InitResult{
		Success: false,
		Files:   []string{},
	}

	outputDir := filepath.Join(initOutput, name)
	result.OutputDir = outputDir

	// Check if directory already exists
	if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
		result.Error = fmt.Sprintf("directory already exists: %s", outputDir)
		return outputInitResult(result)
	}

	// Create directory structure
	dirs := []string{
		outputDir,
		filepath.Join(outputDir, "workflows"),
		filepath.Join(outputDir, "cmd"),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			result.Error = fmt.Sprintf("creating directory %s: %v", dir, err)
			return outputInitResult(result)
		}
	}

	// Create go.mod
	modulePath := "github.com/example/" + name
	goMod := fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0
`, modulePath)
	goModPath := filepath.Join(outputDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goMod), 0644); err != nil {
		result.Error = fmt.Sprintf("writing go.mod: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, goModPath)

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
	workflowsPath := filepath.Join(outputDir, "workflows", "workflows.go")
	if err := os.WriteFile(workflowsPath, []byte(workflowsGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing workflows.go: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, workflowsPath)

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
	triggersPath := filepath.Join(outputDir, "workflows", "triggers.go")
	if err := os.WriteFile(triggersPath, []byte(triggersGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing triggers.go: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, triggersPath)

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
	jobsPath := filepath.Join(outputDir, "workflows", "jobs.go")
	if err := os.WriteFile(jobsPath, []byte(jobsGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing jobs.go: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, jobsPath)

	// Create workflows/steps.go
	stepsGo := `package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BuildSteps = []workflow.Step{
	{Uses: "actions/checkout@v4"},
	{
		Uses: "actions/setup-go@v5",
		With: map[string]any{"go-version": "1.23"},
	},
	{Run: "go build ./..."},
	{Run: "go test ./..."},
}
`
	stepsPath := filepath.Join(outputDir, "workflows", "steps.go")
	if err := os.WriteFile(stepsPath, []byte(stepsGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing steps.go: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, stepsPath)

	// Create cmd/main.go
	mainGo := fmt.Sprintf(`package main

import (
	"fmt"

	// Import workflows to ensure they compile
	_ "%s/workflows"
)

func main() {
	fmt.Println("wetwire-github workflow project")
	fmt.Println("")
	fmt.Println("Build workflows:")
	fmt.Println("  wetwire-github build .")
	fmt.Println("")
	fmt.Println("List workflows:")
	fmt.Println("  wetwire-github list .")
}
`, modulePath)
	mainPath := filepath.Join(outputDir, "cmd", "main.go")
	if err := os.WriteFile(mainPath, []byte(mainGo), 0644); err != nil {
		result.Error = fmt.Sprintf("writing cmd/main.go: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, mainPath)

	// Create README.md
	readme := fmt.Sprintf(`# %s

GitHub workflow declarations using wetwire-github-go.

## Build

Generate YAML workflows:

`+"```"+`bash
wetwire-github build .
`+"```"+`

## Files

- `+"`workflows/workflows.go`"+` - Workflow declarations
- `+"`workflows/triggers.go`"+` - Trigger configurations
- `+"`workflows/jobs.go`"+` - Job declarations
- `+"`workflows/steps.go`"+` - Step lists
`, name)
	readmePath := filepath.Join(outputDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readme), 0644); err != nil {
		result.Error = fmt.Sprintf("writing README.md: %v", err)
		return outputInitResult(result)
	}
	result.Files = append(result.Files, readmePath)

	result.Success = true
	return outputInitResult(result)
}

// outputInitResult outputs the init result in the appropriate format.
func outputInitResult(result wetwire.InitResult) error {
	if initFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Text output
	if !result.Success {
		fmt.Fprintf(os.Stderr, "error: %s\n", result.Error)
		os.Exit(1)
		return nil
	}

	fmt.Printf("Created project: %s\n", result.OutputDir)
	for _, file := range result.Files {
		relPath, _ := filepath.Rel(".", file)
		fmt.Printf("  %s\n", relPath)
	}
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", result.OutputDir)
	fmt.Println("  wetwire-github build .")

	return nil
}
