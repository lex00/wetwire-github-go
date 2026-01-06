package validation

import (
	"io"

	"github.com/rhysd/actionlint"
)

// ActionlintValidator validates workflows using actionlint.
type ActionlintValidator struct {
	// OnlineMode enables checking for actions that reference latest versions
	OnlineMode bool
}

// NewActionlintValidator creates a new ActionlintValidator.
func NewActionlintValidator() *ActionlintValidator {
	return &ActionlintValidator{}
}

// Validate validates a workflow file using actionlint.
func (v *ActionlintValidator) Validate(path string, content []byte) (*ValidationResult, error) {
	// Use actionlint's direct parsing and checking
	opts := &actionlint.LinterOptions{}
	linter, err := actionlint.NewLinter(io.Discard, opts)
	if err != nil {
		return nil, err
	}

	// Lint the workflow content
	errs, err := linter.Lint(path, content, nil)
	if err != nil {
		return nil, err
	}

	result := &ValidationResult{
		Success: len(errs) == 0,
		Issues:  make([]ValidationIssue, 0, len(errs)),
	}

	for _, e := range errs {
		result.Issues = append(result.Issues, ValidationIssue{
			File:    e.Filepath,
			Line:    e.Line,
			Column:  e.Column,
			Message: e.Message,
			RuleID:  e.Kind,
		})
	}

	return result, nil
}

// ValidateFile validates a workflow file from disk.
func (v *ActionlintValidator) ValidateFile(path string) (*ValidationResult, error) {
	opts := &actionlint.LinterOptions{}
	linter, err := actionlint.NewLinter(io.Discard, opts)
	if err != nil {
		return nil, err
	}

	// LintFile reads the file from disk
	errs, err := linter.LintFile(path, nil)
	if err != nil {
		return nil, err
	}

	result := &ValidationResult{
		Success: len(errs) == 0,
		Issues:  make([]ValidationIssue, 0, len(errs)),
	}

	for _, e := range errs {
		result.Issues = append(result.Issues, ValidationIssue{
			File:    e.Filepath,
			Line:    e.Line,
			Column:  e.Column,
			Message: e.Message,
			RuleID:  e.Kind,
		})
	}

	return result, nil
}

// ValidateWorkflow is a convenience function for one-off validation.
func ValidateWorkflow(path string, content []byte) (*ValidationResult, error) {
	v := NewActionlintValidator()
	return v.Validate(path, content)
}

// ValidateWorkflowFile is a convenience function for validating a file on disk.
func ValidateWorkflowFile(path string) (*ValidationResult, error) {
	v := NewActionlintValidator()
	return v.ValidateFile(path)
}
