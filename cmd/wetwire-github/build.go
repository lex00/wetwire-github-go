package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/template"
)

var buildOutput string
var buildFormat string
var buildType string
var buildDryRun bool

var buildCmd = &cobra.Command{
	Use:   "build <path>",
	Short: "Generate YAML from Go workflow declarations",
	Long: `Build reads Go workflow declarations and generates GitHub YAML files.

Example:
  wetwire-github build .
  wetwire-github build ./my-workflows -o .github/workflows/
  wetwire-github build . --type dependabot
  wetwire-github build . --dry-run`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Resolve absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("resolving path: %w", err)
		}

		result := runBuild(absPath, buildOutput, buildDryRun)

		if buildFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if !result.Success {
			for _, err := range result.Errors {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
			return fmt.Errorf("build failed")
		}

		for _, file := range result.Files {
			fmt.Printf("wrote %s\n", file)
		}
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", ".github/workflows", "output directory")
	buildCmd.Flags().StringVar(&buildFormat, "format", "yaml", "output format (yaml, json)")
	buildCmd.Flags().StringVar(&buildType, "type", "workflow", "config type (workflow, dependabot, issue-template)")
	buildCmd.Flags().BoolVar(&buildDryRun, "dry-run", false, "show what would be written without writing")
}

// runBuild executes the build pipeline.
func runBuild(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Step 1: Discover workflows and jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	// Add any discovery errors
	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Workflows) == 0 {
		result.Errors = append(result.Errors, "no workflows found in "+sourcePath)
		return result
	}

	// Step 2: Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractValues(sourcePath, discovered)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("extraction failed: %v", err))
		return result
	}

	if extracted.Error != "" {
		result.Errors = append(result.Errors, extracted.Error)
		return result
	}

	// Step 3: Build templates
	builder := template.NewBuilder()
	built, err := builder.Build(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Step 4: Resolve output directory
	absOutputDir := outputDir
	if !filepath.IsAbs(outputDir) {
		absOutputDir = filepath.Join(sourcePath, outputDir)
	}

	// Step 5: Create output directory if needed
	if !dryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return result
		}
	}

	// Step 6: Write workflow files
	for _, wf := range built.Workflows {
		// Generate filename from workflow name
		filename := toFilename(wf.Name) + ".yml"
		filePath := filepath.Join(absOutputDir, filename)

		if dryRun {
			result.Workflows = append(result.Workflows, wf.Name)
			result.Files = append(result.Files, filePath)
			continue
		}

		if err := os.WriteFile(filePath, wf.YAML, 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing %s: %v", filename, err))
			continue
		}

		result.Workflows = append(result.Workflows, wf.Name)
		result.Files = append(result.Files, filePath)
	}

	// Success if we wrote at least one file and have no errors
	result.Success = len(result.Files) > 0 && len(result.Errors) == 0

	return result
}

// toFilename converts a workflow name to a valid filename.
// "CI" -> "ci", "MyWorkflow" -> "my-workflow"
func toFilename(name string) string {
	// Convert camelCase/PascalCase to kebab-case
	var result strings.Builder
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('-')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
