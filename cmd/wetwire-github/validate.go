package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/validation"
)

var validateFormat string

var validateCmd = &cobra.Command{
	Use:   "validate <workflow.yml> [workflow2.yml...]",
	Short: "Validate YAML using actionlint",
	Long: `Validate checks GitHub Actions workflow YAML files using actionlint.

Example:
  wetwire-github validate .github/workflows/ci.yml
  wetwire-github validate ci.yml --format json
  wetwire-github validate .github/workflows/*.yml

Exit codes:
  0 - All files valid
  1 - Validation errors found
  2 - File not found or other error`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runValidate(args)
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateFormat, "format", "text", "output format (text, json)")
}

// runValidate validates one or more workflow files.
func runValidate(paths []string) error {
	result := wetwire.ValidateResult{
		Success:  true,
		Errors:   []wetwire.ValidationError{},
		Warnings: []string{},
	}

	// Expand globs and validate each file
	var files []string
	for _, pattern := range paths {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, wetwire.ValidationError{
				File:    pattern,
				Message: fmt.Sprintf("invalid pattern: %v", err),
			})
			continue
		}

		if len(matches) == 0 {
			// Check if it's a literal path that doesn't exist
			if _, err := os.Stat(pattern); os.IsNotExist(err) {
				result.Success = false
				result.Errors = append(result.Errors, wetwire.ValidationError{
					File:    pattern,
					Message: "file not found",
				})
				outputResult(result, true)
				return nil
			}
		}

		files = append(files, matches...)
	}

	// Validate each file
	validator := validation.NewActionlintValidator()
	hasErrors := false

	for _, path := range files {
		// Verify file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			result.Success = false
			result.Errors = append(result.Errors, wetwire.ValidationError{
				File:    path,
				Message: "file not found",
			})
			hasErrors = true
			continue
		}

		// Validate using actionlint
		validationResult, err := validator.ValidateFile(path)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, wetwire.ValidationError{
				File:    path,
				Message: fmt.Sprintf("validation error: %v", err),
			})
			hasErrors = true
			continue
		}

		// Convert validation issues to errors
		for _, issue := range validationResult.Issues {
			result.Success = false
			result.Errors = append(result.Errors, wetwire.ValidationError{
				File:    issue.File,
				Line:    issue.Line,
				Column:  issue.Column,
				Message: issue.Message,
				RuleID:  issue.RuleID,
			})
			hasErrors = true
		}
	}

	// Output result
	outputResult(result, hasErrors)
	return nil
}

// outputResult outputs the validation result in the appropriate format.
func outputResult(result wetwire.ValidateResult, hasErrors bool) {
	if validateFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(result)
		if hasErrors {
			os.Exit(1)
		}
		return
	}

	// Text output
	if result.Success && len(result.Errors) == 0 {
		fmt.Println("valid")
		return
	}

	// Check for file-not-found errors (exit code 2)
	fileNotFound := false
	for _, err := range result.Errors {
		if err.Message == "file not found" {
			fmt.Fprintf(os.Stderr, "error: file not found: %s\n", err.File)
			fileNotFound = true
		} else {
			// Standard error format: file:line:column: message
			if err.Line > 0 {
				fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", err.File, err.Line, err.Column, err.Message)
			} else {
				fmt.Fprintf(os.Stderr, "%s: %s\n", err.File, err.Message)
			}
		}
	}

	if fileNotFound {
		os.Exit(2)
	}
	os.Exit(1)
}
