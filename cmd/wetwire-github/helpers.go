package main

import (
	"fmt"
	"os"
	"path/filepath"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/linter"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/template"
)

// getVersion returns the version string.
func getVersion() string {
	return version
}

// runBuild executes the workflow build pipeline.
// This is a simplified version used by watch and MCP commands.
func runBuild(sourcePath, outputDir string, dryRun bool) wetwire.BuildResult {
	result := wetwire.BuildResult{
		Success:   false,
		Workflows: []string{},
		Files:     []string{},
		Errors:    []string{},
	}

	// Discover workflows and jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(sourcePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("discovery failed: %v", err))
		return result
	}

	result.Errors = append(result.Errors, discovered.Errors...)

	if len(discovered.Workflows) == 0 {
		result.Errors = append(result.Errors, "no workflows found in "+sourcePath)
		return result
	}

	// Extract values using runner
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

	// Build templates
	builder := template.NewBuilder()
	built, err := builder.Build(discovered, extracted)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("template build failed: %v", err))
		return result
	}

	result.Errors = append(result.Errors, built.Errors...)

	// Resolve output directory
	absOutputDir := outputDir
	if !filepath.IsAbs(outputDir) {
		absOutputDir = filepath.Join(sourcePath, outputDir)
	}

	// Create output directory if needed
	if !dryRun {
		if err := os.MkdirAll(absOutputDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("creating output directory: %v", err))
			return result
		}
	}

	// Write workflow files
	for _, wf := range built.Workflows {
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

	result.Success = len(result.Files) > 0 && len(result.Errors) == 0
	return result
}

// runLintPath runs the linter on the given path and returns a result.
// This is used by the watch command to check lint status without exiting.
func runLintPath(path string) LintPathResult {
	result := LintPathResult{Success: false}

	absPath, err := filepath.Abs(path)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("resolving path: %v", err))
		return result
	}

	info, err := os.Stat(absPath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("accessing path: %v", err))
		return result
	}

	l := linter.DefaultLinter()

	var lintResult *linter.LintResult
	if info.IsDir() {
		lintResult, err = l.LintDir(absPath)
	} else {
		lintResult, err = l.LintFile(absPath)
	}

	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("linting failed: %v", err))
		return result
	}

	result.Success = lintResult.Success
	result.Issues = make([]wetwire.LintIssue, len(lintResult.Issues))

	for i, issue := range lintResult.Issues {
		result.Issues[i] = wetwire.LintIssue{
			File:     issue.File,
			Line:     issue.Line,
			Column:   issue.Column,
			Severity: issue.Severity,
			Message:  issue.Message,
			Rule:     issue.Rule,
			Fixable:  issue.Fixable,
		}
	}

	return result
}

// LintPathResult is an extended lint result with an Errors field for watch mode.
type LintPathResult struct {
	Success bool
	Issues  []wetwire.LintIssue
	Errors  []string
}

// toFilename converts a workflow name to a valid filename.
// "CI" -> "ci", "MyWorkflow" -> "my-workflow"
func toFilename(name string) string {
	var result string
	for i, r := range name {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "-"
		}
		result += string(r)
	}
	return toFilenameCase(result)
}

func toFilenameCase(s string) string {
	return toLowerString(s)
}

func toLowerString(s string) string {
	result := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			result += string(r + 32)
		} else {
			result += string(r)
		}
	}
	return result
}
