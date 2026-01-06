// Package validation provides YAML validation using actionlint.
package validation

// ValidationResult contains the results of validating a workflow.
type ValidationResult struct {
	Success bool              `json:"success"`
	Issues  []ValidationIssue `json:"issues,omitempty"`
}

// ValidationIssue represents a single validation issue.
type ValidationIssue struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	RuleID  string `json:"rule_id,omitempty"`
}

// Severity returns the severity level of the issue.
// Currently all actionlint issues are errors.
func (v ValidationIssue) Severity() string {
	return "error"
}

// Validator is the interface for workflow validators.
type Validator interface {
	Validate(path string, content []byte) (*ValidationResult, error)
}

// ValidatorPipeline chains multiple validators together.
type ValidatorPipeline struct {
	validators []Validator
}

// NewValidatorPipeline creates a new validator pipeline.
func NewValidatorPipeline(validators ...Validator) *ValidatorPipeline {
	return &ValidatorPipeline{
		validators: validators,
	}
}

// Validate runs all validators and combines their results.
func (p *ValidatorPipeline) Validate(path string, content []byte) (*ValidationResult, error) {
	combined := &ValidationResult{
		Success: true,
		Issues:  []ValidationIssue{},
	}

	for _, v := range p.validators {
		result, err := v.Validate(path, content)
		if err != nil {
			return nil, err
		}

		if !result.Success {
			combined.Success = false
		}

		combined.Issues = append(combined.Issues, result.Issues...)
	}

	return combined, nil
}
