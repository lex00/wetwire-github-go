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
	buildCmd.Flags().StringVar(&buildType, "type", "workflow", "config type (workflow, dependabot, issue-template, discussion-template, pr-template, codeowners)")
	buildCmd.Flags().BoolVar(&buildDryRun, "dry-run", false, "show what would be written without writing")
}

// runBuild executes the build pipeline.
func runBuild(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	// Handle different config types
	switch buildType {
	case "dependabot":
		return runBuildDependabot(sourcePath, outputDir, dryRun)
	case "issue-template":
		return runBuildIssueTemplate(sourcePath, outputDir, dryRun)
	case "discussion-template":
		return runBuildDiscussionTemplate(sourcePath, outputDir, dryRun)
	case "pr-template":
		return runBuildPRTemplate(sourcePath, outputDir, dryRun)
	case "codeowners":
		return runBuildCodeowners(sourcePath, outputDir, dryRun)
	default:
		return runBuildWorkflow(sourcePath, outputDir, dryRun)
	}
}

// runBuildWorkflow executes the workflow build pipeline.
func runBuildWorkflow(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
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

// runBuildDependabot executes the dependabot build pipeline.
func runBuildDependabot(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Step 1: Discover Dependabot configs
	disc := discover.NewDiscoverer()
	discovered, err := disc.DiscoverDependabot(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	// Add any discovery errors
	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Configs) == 0 {
		result.Errors = append(result.Errors, "no dependabot configs found in "+sourcePath)
		return result
	}

	// Step 2: Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractDependabot(sourcePath, discovered)
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
	built, err := builder.BuildDependabot(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Step 4: Resolve output directory (dependabot goes in .github/)
	absOutputDir := outputDir
	if outputDir == ".github/workflows" {
		// Default for dependabot is .github/
		absOutputDir = ".github"
	}
	if !filepath.IsAbs(absOutputDir) {
		absOutputDir = filepath.Join(sourcePath, absOutputDir)
	}

	// Step 5: Create output directory if needed
	if !dryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return result
		}
	}

	// Step 6: Write dependabot file (always named dependabot.yml)
	for _, cfg := range built.Configs {
		filePath := filepath.Join(absOutputDir, "dependabot.yml")

		if dryRun {
			result.Workflows = append(result.Workflows, cfg.Name)
			result.Files = append(result.Files, filePath)
			continue
		}

		if err := os.WriteFile(filePath, cfg.YAML, 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing dependabot.yml: %v", err))
			continue
		}

		result.Workflows = append(result.Workflows, cfg.Name)
		result.Files = append(result.Files, filePath)
	}

	// Success if we wrote at least one file and have no errors
	result.Success = len(result.Files) > 0 && len(result.Errors) == 0

	return result
}

// runBuildIssueTemplate executes the issue template build pipeline.
func runBuildIssueTemplate(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Step 1: Discover IssueTemplates
	disc := discover.NewDiscoverer()
	discovered, err := disc.DiscoverIssueTemplates(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	// Add any discovery errors
	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Templates) == 0 {
		result.Errors = append(result.Errors, "no issue templates found in "+sourcePath)
		return result
	}

	// Step 2: Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractIssueTemplates(sourcePath, discovered)
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
	built, err := builder.BuildIssueTemplates(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Step 4: Resolve output directory (issue templates go in .github/ISSUE_TEMPLATE/)
	absOutputDir := outputDir
	if outputDir == ".github/workflows" {
		// Default for issue templates is .github/ISSUE_TEMPLATE/
		absOutputDir = ".github/ISSUE_TEMPLATE"
	}
	if !filepath.IsAbs(absOutputDir) {
		absOutputDir = filepath.Join(sourcePath, absOutputDir)
	}

	// Step 5: Create output directory if needed
	if !dryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return result
		}
	}

	// Step 6: Write issue template files
	for _, tmpl := range built.Templates {
		// Generate filename from template name
		filename := toFilename(tmpl.Name) + ".yml"
		filePath := filepath.Join(absOutputDir, filename)

		if dryRun {
			result.Workflows = append(result.Workflows, tmpl.Name)
			result.Files = append(result.Files, filePath)
			continue
		}

		if err := os.WriteFile(filePath, tmpl.YAML, 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing %s: %v", filename, err))
			continue
		}

		result.Workflows = append(result.Workflows, tmpl.Name)
		result.Files = append(result.Files, filePath)
	}

	// Success if we wrote at least one file and have no errors
	result.Success = len(result.Files) > 0 && len(result.Errors) == 0

	return result
}

// runBuildDiscussionTemplate executes the discussion template build pipeline.
func runBuildDiscussionTemplate(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Step 1: Discover DiscussionTemplates
	disc := discover.NewDiscoverer()
	discovered, err := disc.DiscoverDiscussionTemplates(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	// Add any discovery errors
	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Templates) == 0 {
		result.Errors = append(result.Errors, "no discussion templates found in "+sourcePath)
		return result
	}

	// Step 2: Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractDiscussionTemplates(sourcePath, discovered)
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
	built, err := builder.BuildDiscussionTemplates(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Step 4: Resolve output directory (discussion templates go in .github/DISCUSSION_TEMPLATE/)
	absOutputDir := outputDir
	if outputDir == ".github/workflows" {
		// Default for discussion templates is .github/DISCUSSION_TEMPLATE/
		absOutputDir = ".github/DISCUSSION_TEMPLATE"
	}
	if !filepath.IsAbs(absOutputDir) {
		absOutputDir = filepath.Join(sourcePath, absOutputDir)
	}

	// Step 5: Create output directory if needed
	if !dryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return result
		}
	}

	// Step 6: Write discussion template files
	for _, tmpl := range built.Templates {
		// Generate filename from template name
		filename := toFilename(tmpl.Name) + ".yml"
		filePath := filepath.Join(absOutputDir, filename)

		if dryRun {
			result.Workflows = append(result.Workflows, tmpl.Name)
			result.Files = append(result.Files, filePath)
			continue
		}

		if err := os.WriteFile(filePath, tmpl.YAML, 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing %s: %v", filename, err))
			continue
		}

		result.Workflows = append(result.Workflows, tmpl.Name)
		result.Files = append(result.Files, filePath)
	}

	// Success if we wrote at least one file and have no errors
	result.Success = len(result.Files) > 0 && len(result.Errors) == 0

	return result
}

// runBuildPRTemplate executes the PR template build pipeline.
func runBuildPRTemplate(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Step 1: Discover PRTemplates
	disc := discover.NewDiscoverer()
	discovered, err := disc.DiscoverPRTemplates(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	// Add any discovery errors
	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Templates) == 0 {
		result.Errors = append(result.Errors, "no PR templates found in "+sourcePath)
		return result
	}

	// Step 2: Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractPRTemplates(sourcePath, discovered)
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
	built, err := builder.BuildPRTemplates(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Step 4: Resolve output directory (PR templates go in .github/)
	absOutputDir := outputDir
	if outputDir == ".github/workflows" {
		// Default for PR templates is .github/
		absOutputDir = ".github"
	}
	if !filepath.IsAbs(absOutputDir) {
		absOutputDir = filepath.Join(sourcePath, absOutputDir)
	}

	// Step 5: Write PR template files
	for _, tmpl := range built.Templates {
		// Use the Filename from the template (handles PULL_REQUEST_TEMPLATE/ subdirectory)
		filePath := filepath.Join(absOutputDir, tmpl.Filename)

		// Create parent directory if needed (for named templates in PULL_REQUEST_TEMPLATE/)
		if !dryRun {
			parentDir := filepath.Dir(filePath)
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("creating directory: %v", err))
				continue
			}
		}

		if dryRun {
			result.Workflows = append(result.Workflows, tmpl.Name)
			result.Files = append(result.Files, filePath)
			continue
		}

		if err := os.WriteFile(filePath, tmpl.Content, 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing %s: %v", tmpl.Filename, err))
			continue
		}

		result.Workflows = append(result.Workflows, tmpl.Name)
		result.Files = append(result.Files, filePath)
	}

	// Success if we wrote at least one file and have no errors
	result.Success = len(result.Files) > 0 && len(result.Errors) == 0

	return result
}

// runBuildCodeowners executes the CODEOWNERS build pipeline.
func runBuildCodeowners(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Step 1: Discover Codeowners configs
	disc := discover.NewDiscoverer()
	discovered, err := disc.DiscoverCodeowners(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	// Add any discovery errors
	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Configs) == 0 {
		result.Errors = append(result.Errors, "no codeowners configs found in "+sourcePath)
		return result
	}

	// Step 2: Extract values using runner
	run := runner.NewRunner()
	extracted, err := run.ExtractCodeowners(sourcePath, discovered)
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
	built, err := builder.BuildCodeowners(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	// Add template builder errors
	result.Errors = append(result.Errors, built.Errors...)

	// Step 4: Resolve output directory (CODEOWNERS goes in .github/)
	absOutputDir := outputDir
	if outputDir == ".github/workflows" {
		// Default for codeowners is .github/
		absOutputDir = ".github"
	}
	if !filepath.IsAbs(absOutputDir) {
		absOutputDir = filepath.Join(sourcePath, absOutputDir)
	}

	// Step 5: Create output directory if needed
	if !dryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return result
		}
	}

	// Step 6: Write CODEOWNERS file (always named CODEOWNERS)
	for _, cfg := range built.Configs {
		filePath := filepath.Join(absOutputDir, "CODEOWNERS")

		if dryRun {
			result.Workflows = append(result.Workflows, cfg.Name)
			result.Files = append(result.Files, filePath)
			continue
		}

		if err := os.WriteFile(filePath, cfg.Content, 0644); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("writing CODEOWNERS: %v", err))
			continue
		}

		result.Workflows = append(result.Workflows, cfg.Name)
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
