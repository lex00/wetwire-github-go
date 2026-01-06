package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initOutput string

var initCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Create a new workflow project",
	Long: `Init creates a new wetwire-github workflow project with example declarations.

Example:
  wetwire-github init my-workflows
  wetwire-github init my-ci -o ./projects/`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		outputDir := filepath.Join(initOutput, name)

		// Check if directory already exists
		if _, err := os.Stat(outputDir); !os.IsNotExist(err) {
			return fmt.Errorf("directory already exists: %s", outputDir)
		}

		// Create directory structure
		dirs := []string{
			outputDir,
			filepath.Join(outputDir, "cmd"),
		}
		for _, dir := range dirs {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating directory %s: %w", dir, err)
			}
		}

		// Create go.mod
		goMod := fmt.Sprintf(`module %s

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0
`, name)
		if err := os.WriteFile(filepath.Join(outputDir, "go.mod"), []byte(goMod), 0644); err != nil {
			return fmt.Errorf("writing go.mod: %w", err)
		}

		// Create workflows.go
		workflowsGo := fmt.Sprintf(`package %s

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// CI workflow runs on push and pull requests to main
var CI = workflow.Workflow{
	Name: "CI",
	On:   CITriggers,
}

var CITriggers = workflow.Triggers{
	Push:        &workflow.PushTrigger{Branches: workflow.List("main")},
	PullRequest: &workflow.PullRequestTrigger{Branches: workflow.List("main")},
}
`, name)
		if err := os.WriteFile(filepath.Join(outputDir, "workflows.go"), []byte(workflowsGo), 0644); err != nil {
			return fmt.Errorf("writing workflows.go: %w", err)
		}

		// Create jobs.go
		jobsGo := fmt.Sprintf(`package %s

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// Build job compiles and tests the code
var Build = workflow.Job{
	Name:   "build",
	RunsOn: "ubuntu-latest",
	Steps:  BuildSteps,
}

var BuildSteps = workflow.List(
	workflow.Step{Uses: "actions/checkout@v4"},
	workflow.Step{
		Uses: "actions/setup-go@v5",
		With: workflow.With{"go-version": "1.23"},
	},
	workflow.Step{Run: "go build ./..."},
	workflow.Step{Run: "go test ./..."},
)
`, name)
		if err := os.WriteFile(filepath.Join(outputDir, "jobs.go"), []byte(jobsGo), 0644); err != nil {
			return fmt.Errorf("writing jobs.go: %w", err)
		}

		// Create cmd/main.go
		mainGo := `package main

import "fmt"

func main() {
	fmt.Println("wetwire-github workflow project")
	fmt.Println("")
	fmt.Println("Build workflows:")
	fmt.Println("  wetwire-github build .")
	fmt.Println("")
	fmt.Println("List workflows:")
	fmt.Println("  wetwire-github list .")
}
`
		if err := os.WriteFile(filepath.Join(outputDir, "cmd", "main.go"), []byte(mainGo), 0644); err != nil {
			return fmt.Errorf("writing cmd/main.go: %w", err)
		}

		// Create README.md
		readme := fmt.Sprintf(`# %s

GitHub workflow declarations using wetwire-github-go.

## Build

Generate YAML workflows:

`+"```"+`bash
wetwire-github build .
`+"```"+`

## Files

- `+"`workflows.go`"+` - Workflow declarations
- `+"`jobs.go`"+` - Job declarations
`, name)
		if err := os.WriteFile(filepath.Join(outputDir, "README.md"), []byte(readme), 0644); err != nil {
			return fmt.Errorf("writing README.md: %w", err)
		}

		fmt.Printf("Created project: %s\n", outputDir)
		fmt.Printf("  %s/go.mod\n", outputDir)
		fmt.Printf("  %s/workflows.go\n", outputDir)
		fmt.Printf("  %s/jobs.go\n", outputDir)
		fmt.Printf("  %s/cmd/main.go\n", outputDir)
		fmt.Printf("  %s/README.md\n", outputDir)
		fmt.Println("")
		fmt.Println("Next steps:")
		fmt.Printf("  cd %s\n", outputDir)
		fmt.Println("  wetwire-github build .")

		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&initOutput, "output", "o", ".", "output directory")
}
