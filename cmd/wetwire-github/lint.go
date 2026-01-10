package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/linter"
)

var lintFormat string
var lintFix bool

var lintCmd = &cobra.Command{
	Use:   "lint <path>",
	Short: "Check Go code for wetwire best practices",
	Long: `Lint checks Go workflow declarations for wetwire best practices.

Rules:
  WAG001: Use typed action wrappers instead of raw uses: strings
  WAG002: Use condition builders instead of raw expression strings
  WAG003: Use secrets context instead of hardcoded strings
  WAG004: Use matrix builder instead of inline maps
  WAG005: Extract inline structs to named variables
  WAG006: Detect duplicate workflow names
  WAG007: Flag oversized files (>N jobs)
  WAG008: Avoid hardcoded expression strings

Example:
  wetwire-github lint .
  wetwire-github lint ./my-workflows --format json
  wetwire-github lint . --fix

Exit codes:
  0 - No issues found
  1 - Issues found
  2 - Parse error or other failure`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLint(args[0])
	},
}

func init() {
	lintCmd.Flags().StringVar(&lintFormat, "format", "text", "output format (text, json)")
	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "automatically fix issues where possible")
}

// runLint executes the linter on the given path.
func runLint(path string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		outputLintError(path, fmt.Sprintf("resolving path: %v", err))
		return nil
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			outputLintError(path, "path not found")
		} else {
			outputLintError(path, fmt.Sprintf("accessing path: %v", err))
		}
		return nil
	}

	// Create linter with default rules
	l := linter.DefaultLinter()

	// Handle --fix flag
	if lintFix {
		return runLintWithFix(l, absPath, info.IsDir())
	}

	var lintResult *linter.LintResult

	// Lint file or directory
	if info.IsDir() {
		lintResult, err = l.LintDir(absPath)
	} else {
		lintResult, err = l.LintFile(absPath)
	}

	if err != nil {
		outputLintError(path, fmt.Sprintf("linting failed: %v", err))
		return nil
	}

	// Convert to CLI result type
	result := wetwire.LintResult{
		Success: lintResult.Success,
		Issues:  make([]wetwire.LintIssue, len(lintResult.Issues)),
	}

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

	// Output result
	outputLintResult(result)
	return nil
}

// runLintWithFix applies automatic fixes to the code.
func runLintWithFix(l *linter.Linter, absPath string, isDir bool) error {
	var fixResult *linter.FixResult
	var fixDirResult *linter.FixDirResult
	var err error

	if isDir {
		fixDirResult, err = l.FixDir(absPath)
		if err != nil {
			outputLintError(absPath, fmt.Sprintf("fix failed: %v", err))
			return nil
		}

		if lintFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(map[string]any{
				"success":     true,
				"fixed_count": fixDirResult.TotalFixed,
				"files":       fixDirResult.Files,
			})
		} else {
			if fixDirResult.TotalFixed > 0 {
				fmt.Printf("Fixed %d issue(s) in %d file(s):\n", fixDirResult.TotalFixed, len(fixDirResult.Files))
				for _, f := range fixDirResult.Files {
					fmt.Printf("  %s\n", f)
				}
			} else {
				fmt.Println("No fixable issues found")
			}
		}
	} else {
		fixResult, err = l.FixFile(absPath)
		if err != nil {
			outputLintError(absPath, fmt.Sprintf("fix failed: %v", err))
			return nil
		}

		if lintFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(map[string]any{
				"success":     true,
				"fixed_count": fixResult.FixedCount,
				"remaining":   len(fixResult.Issues),
			})
		} else {
			if fixResult.FixedCount > 0 {
				fmt.Printf("Fixed %d issue(s) in %s\n", fixResult.FixedCount, absPath)
			} else {
				fmt.Println("No fixable issues found")
			}

			// Report remaining unfixable issues
			if len(fixResult.Issues) > 0 {
				fmt.Printf("\n%d issue(s) could not be fixed:\n", len(fixResult.Issues))
				for _, issue := range fixResult.Issues {
					fmt.Printf("  %s:%d:%d: %s [%s]\n",
						issue.File, issue.Line, issue.Column, issue.Message, issue.Rule)
				}
			}
		}
	}

	return nil
}

// outputLintError outputs an error in the appropriate format.
func outputLintError(file, message string) {
	if lintFormat == "json" {
		result := wetwire.LintResult{
			Success: false,
			Issues: []wetwire.LintIssue{
				{File: file, Line: 1, Column: 1, Severity: "error", Message: message, Rule: "error"},
			},
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		os.Exit(2)
		return
	}

	fmt.Fprintf(os.Stderr, "error: %s: %s\n", file, message)
	os.Exit(2)
}

// outputLintResult outputs the lint result in the appropriate format.
func outputLintResult(result wetwire.LintResult) {
	if lintFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		if !result.Success {
			os.Exit(1)
		}
		return
	}

	// Text output
	if len(result.Issues) == 0 {
		fmt.Println("No issues found")
		return
	}

	// Group issues by severity for better output
	errorCount := 0
	warningCount := 0

	for _, issue := range result.Issues {
		// Format: file:line:column: [rule] message
		severity := issue.Severity
		if severity == "" {
			severity = "warning"
		}

		switch severity {
		case "error":
			errorCount++
		case "warning":
			warningCount++
		}

		fmt.Printf("%s:%d:%d: %s: %s [%s]\n",
			issue.File, issue.Line, issue.Column, severity, issue.Message, issue.Rule)
	}

	// Print summary
	fmt.Printf("\n%d issue(s): %d error(s), %d warning(s)\n",
		len(result.Issues), errorCount, warningCount)

	if !result.Success {
		os.Exit(1)
	}
}
