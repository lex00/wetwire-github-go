package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/importer"
)

var importOutput string
var importSingleFile bool
var importNoScaffold bool
var importType string
var importFormat string

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Convert existing config files to Go code",
	Long: `Import converts existing GitHub configuration files to typed Go declarations.

Supported types:
  workflow          GitHub Actions YAML files (default)
  dependabot        Dependabot configuration
  issue-template    Issue template YAML files
  discussion-template Discussion template YAML files
  codeowners        CODEOWNERS file

Example:
  wetwire-github import .github/workflows/ci.yml -o my-workflows/
  wetwire-github import ci.yml --single-file
  wetwire-github import ci.yml --no-scaffold
  wetwire-github import .github/CODEOWNERS --type codeowners`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runImport(args[0])
	},
}

func init() {
	importCmd.Flags().StringVarP(&importOutput, "output", "o", ".", "output directory")
	importCmd.Flags().BoolVar(&importSingleFile, "single-file", false, "generate all code in a single file")
	importCmd.Flags().BoolVar(&importNoScaffold, "no-scaffold", false, "skip generating go.mod, README, etc.")
	importCmd.Flags().StringVar(&importType, "type", "workflow", "config type (workflow, dependabot, issue-template, discussion-template, codeowners)")
	importCmd.Flags().StringVar(&importFormat, "format", "text", "output format (text, json)")
}

// runImport executes the import command.
func runImport(path string) error {
	// Dispatch based on type
	switch importType {
	case "codeowners":
		return runImportCodeowners(path)
	default:
		return runImportWorkflow(path)
	}
}

// runImportCodeowners handles CODEOWNERS file import.
func runImportCodeowners(path string) error {
	result := wetwire.ImportResult{
		Success:   false,
		OutputDir: importOutput,
		Files:     []string{},
		Errors:    []string{},
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("file not found: %s", path))
		outputImportResult(result)
		return nil
	}

	// Parse the CODEOWNERS file
	codeowners, err := importer.ParseCodeownersFile(path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("parsing CODEOWNERS: %v", err))
		outputImportResult(result)
		return nil
	}

	// Create code generator
	gen := &importer.CodeownersCodeGenerator{PackageName: "workflows"}

	// Generate code
	code, err := gen.Generate(codeowners)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("generating code: %v", err))
		outputImportResult(result)
		return nil
	}

	// Resolve output directory
	absOutput, err := filepath.Abs(importOutput)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("resolving output path: %v", err))
		outputImportResult(result)
		return nil
	}
	result.OutputDir = absOutput

	// Create output directory and workflows subdirectory
	workflowsDir := filepath.Join(absOutput, "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
		outputImportResult(result)
		return nil
	}

	// Write generated code to workflows/ subdirectory
	for filename, content := range code.Files {
		filePath := filepath.Join(workflowsDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing %s: %v", filename, err))
			outputImportResult(result)
			return nil
		}
		result.Files = append(result.Files, filePath)
	}

	// Generate scaffold files if not disabled
	if !importNoScaffold {
		projectName := filepath.Base(absOutput)
		if projectName == "." {
			projectName = "codeowners"
		}
		modulePath := "github.com/example/" + strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))

		scaffold := importer.NewScaffold(modulePath, projectName)
		scaffoldFiles := scaffold.Generate()

		if err := importer.WriteScaffold(absOutput, scaffoldFiles); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing scaffold: %v", err))
			outputImportResult(result)
			return nil
		}

		for filename := range scaffoldFiles.Files {
			result.Files = append(result.Files, filepath.Join(absOutput, filename))
		}
	}

	result.Success = true
	// CODEOWNERS doesn't have workflows/jobs/steps, but we report rules count
	// We'll use a custom output for CODEOWNERS

	outputImportCodeownersResult(result, code.Rules)
	return nil
}

// runImportWorkflow handles workflow file import (original behavior).
func runImportWorkflow(path string) error {
	result := wetwire.ImportResult{
		Success:   false,
		OutputDir: importOutput,
		Files:     []string{},
		Errors:    []string{},
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		result.Errors = append(result.Errors, fmt.Sprintf("file not found: %s", path))
		outputImportResult(result)
		return nil
	}

	// Parse the YAML file
	workflow, err := importer.ParseWorkflowFile(path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("parsing YAML: %v", err))
		outputImportResult(result)
		return nil
	}

	// Derive workflow name from filename
	workflowName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if workflow.Name != "" {
		workflowName = workflow.Name
	}

	// Create code generator
	gen := importer.NewCodeGenerator()
	gen.SingleFile = importSingleFile
	// Keep default "workflows" package name for library-style imports

	// Generate code
	code, err := gen.Generate(workflow, workflowName)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("generating code: %v", err))
		outputImportResult(result)
		return nil
	}

	// Resolve output directory
	absOutput, err := filepath.Abs(importOutput)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("resolving output path: %v", err))
		outputImportResult(result)
		return nil
	}
	result.OutputDir = absOutput

	// Create output directory and workflows subdirectory
	workflowsDir := filepath.Join(absOutput, "workflows")
	if err := os.MkdirAll(workflowsDir, 0755); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
		outputImportResult(result)
		return nil
	}

	// Write generated code to workflows/ subdirectory
	if err := importer.WriteGeneratedCode(workflowsDir, code); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("writing code: %v", err))
		outputImportResult(result)
		return nil
	}

	for filename := range code.Files {
		result.Files = append(result.Files, filepath.Join(workflowsDir, filename))
	}

	// Generate scaffold files if not disabled
	if !importNoScaffold {
		projectName := filepath.Base(absOutput)
		if projectName == "." {
			projectName = workflowName
		}
		modulePath := "github.com/example/" + strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))

		scaffold := importer.NewScaffold(modulePath, projectName)
		scaffoldFiles := scaffold.Generate()

		if err := importer.WriteScaffold(absOutput, scaffoldFiles); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing scaffold: %v", err))
			outputImportResult(result)
			return nil
		}

		for filename := range scaffoldFiles.Files {
			result.Files = append(result.Files, filepath.Join(absOutput, filename))
		}
	}

	result.Success = true
	result.Workflows = code.Workflows
	result.Jobs = code.Jobs
	result.Steps = code.Steps

	outputImportResult(result)
	return nil
}

// outputImportResult outputs the import result in the appropriate format.
func outputImportResult(result wetwire.ImportResult) {
	if importFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		if !result.Success {
			os.Exit(1)
		}
		return
	}

	// Text output
	if !result.Success {
		for _, err := range result.Errors {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		os.Exit(1)
		return
	}

	fmt.Printf("Imported %d workflow(s), %d job(s), %d step(s)\n",
		result.Workflows, result.Jobs, result.Steps)
	fmt.Printf("Output: %s\n", result.OutputDir)
	for _, file := range result.Files {
		relPath, _ := filepath.Rel(result.OutputDir, file)
		fmt.Printf("  %s\n", relPath)
	}
}

// outputImportCodeownersResult outputs CODEOWNERS import result.
func outputImportCodeownersResult(result wetwire.ImportResult, rules int) {
	if importFormat == "json" {
		// Extend result with rules count for JSON output
		type codeownersResult struct {
			wetwire.ImportResult
			Rules int `json:"rules"`
		}
		extended := codeownersResult{
			ImportResult: result,
			Rules:        rules,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(extended)
		if !result.Success {
			os.Exit(1)
		}
		return
	}

	// Text output
	if !result.Success {
		for _, err := range result.Errors {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		os.Exit(1)
		return
	}

	fmt.Printf("Imported CODEOWNERS with %d rule(s)\n", rules)
	fmt.Printf("Output: %s\n", result.OutputDir)
	for _, file := range result.Files {
		relPath, _ := filepath.Rel(result.OutputDir, file)
		fmt.Printf("  %s\n", relPath)
	}
}
