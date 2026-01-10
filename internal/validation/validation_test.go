package validation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestActionlintValidator_Validate_Valid(t *testing.T) {
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "Hello"
`)

	v := NewActionlintValidator()
	result, err := v.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if !result.Success {
		t.Errorf("Validate() Success = false, want true")
		for _, issue := range result.Issues {
			t.Logf("Issue: %s:%d:%d: %s", issue.File, issue.Line, issue.Column, issue.Message)
		}
	}
}

func TestActionlintValidator_Validate_Invalid(t *testing.T) {
	// Invalid YAML - unknown key
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    invalid-key: true
    steps:
      - uses: actions/checkout@v4
`)

	v := NewActionlintValidator()
	result, err := v.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if result.Success {
		t.Error("Validate() Success = true, want false for invalid workflow")
	}

	if len(result.Issues) == 0 {
		t.Error("Validate() expected issues for invalid workflow")
	}
}

func TestActionlintValidator_Validate_MissingOn(t *testing.T) {
	content := []byte(`name: CI
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "Hello"
`)

	v := NewActionlintValidator()
	result, err := v.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if result.Success {
		t.Error("Validate() Success = true, want false for missing 'on'")
	}
}

func TestActionlintValidator_Validate_MalformedYAML(t *testing.T) {
	// Completely malformed YAML
	content := []byte(`{
invalid yaml: [[[
this is not valid
`)

	v := NewActionlintValidator()
	result, err := v.Validate("malformed.yml", content)

	// actionlint should return an error for malformed YAML
	if err == nil && result.Success {
		t.Error("Validate() expected error or failed validation for malformed YAML")
	}
}

func TestActionlintValidator_Validate_EmptyWorkflow(t *testing.T) {
	content := []byte(``)

	v := NewActionlintValidator()
	result, err := v.Validate("empty.yml", content)

	// Empty workflow should fail validation or return error
	if err == nil && result.Success {
		t.Error("Validate() expected error or failed validation for empty workflow")
	}
}

func TestActionlintValidator_Validate_MultipleIssues(t *testing.T) {
	// Workflow with multiple validation issues
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    invalid-key: true
    steps:
      - uses: nonexistent/action@v1
      - run: ${{ invalid.expression }}
`)

	v := NewActionlintValidator()
	result, err := v.Validate("multi-issue.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if result.Success {
		t.Error("Validate() Success = true, want false for workflow with multiple issues")
	}

	// Verify that issues are properly populated with all required fields
	for _, issue := range result.Issues {
		if issue.File == "" {
			t.Error("Issue missing File field")
		}
		if issue.Message == "" {
			t.Error("Issue missing Message field")
		}
		// Line and Column can be 0 for some issues, but should be set for most
	}
}

func TestActionlintValidator_ValidateFile(t *testing.T) {
	// Create a temporary file with valid workflow
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "test-workflow.yml")

	content := []byte(`name: Test
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "test"
`)

	if err := os.WriteFile(workflowPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	v := NewActionlintValidator()
	result, err := v.ValidateFile(workflowPath)
	if err != nil {
		t.Fatalf("ValidateFile() error = %v", err)
	}

	if !result.Success {
		t.Errorf("ValidateFile() Success = false, want true")
		for _, issue := range result.Issues {
			t.Logf("Issue: %s:%d:%d: %s", issue.File, issue.Line, issue.Column, issue.Message)
		}
	}

	if len(result.Issues) > 0 {
		t.Errorf("ValidateFile() expected no issues, got %d", len(result.Issues))
	}
}

func TestActionlintValidator_ValidateFile_Invalid(t *testing.T) {
	// Create a temporary file with invalid workflow
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "invalid-workflow.yml")

	content := []byte(`name: Test
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    invalid-field: true
    steps:
      - run: echo "test"
`)

	if err := os.WriteFile(workflowPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	v := NewActionlintValidator()
	result, err := v.ValidateFile(workflowPath)
	if err != nil {
		t.Fatalf("ValidateFile() error = %v", err)
	}

	if result.Success {
		t.Error("ValidateFile() Success = true, want false for invalid workflow")
	}

	if len(result.Issues) == 0 {
		t.Error("ValidateFile() expected issues for invalid workflow")
	}
}

func TestActionlintValidator_ValidateFile_NonExistent(t *testing.T) {
	v := NewActionlintValidator()
	result, err := v.ValidateFile("/nonexistent/path/to/workflow.yml")

	// Should return an error for non-existent file
	if err == nil {
		t.Error("ValidateFile() expected error for non-existent file")
	}

	// Result could be nil when error is returned
	if result != nil && result.Success {
		t.Error("ValidateFile() Success = true, want false or error for non-existent file")
	}
}

func TestNewActionlintValidator(t *testing.T) {
	v := NewActionlintValidator()
	if v == nil {
		t.Error("NewActionlintValidator() returned nil")
	}

	// Verify default values
	if v.OnlineMode {
		t.Error("NewActionlintValidator() OnlineMode should default to false")
	}
}

func TestActionlintValidator_OnlineMode(t *testing.T) {
	v := NewActionlintValidator()
	v.OnlineMode = true

	if !v.OnlineMode {
		t.Error("Failed to set OnlineMode to true")
	}

	// OnlineMode doesn't affect basic validation
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	result, err := v.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if !result.Success {
		t.Error("Validate() with OnlineMode failed")
	}
}

func TestValidationIssue_Severity(t *testing.T) {
	issue := ValidationIssue{
		File:    "test.yml",
		Line:    1,
		Column:  1,
		Message: "test error",
	}

	if issue.Severity() != "error" {
		t.Errorf("Severity() = %q, want %q", issue.Severity(), "error")
	}
}

func TestValidationIssue_Fields(t *testing.T) {
	issue := ValidationIssue{
		File:    "workflow.yml",
		Line:    42,
		Column:  10,
		Message: "invalid syntax",
		RuleID:  "syntax-check",
	}

	if issue.File != "workflow.yml" {
		t.Errorf("File = %q, want %q", issue.File, "workflow.yml")
	}
	if issue.Line != 42 {
		t.Errorf("Line = %d, want %d", issue.Line, 42)
	}
	if issue.Column != 10 {
		t.Errorf("Column = %d, want %d", issue.Column, 10)
	}
	if issue.Message != "invalid syntax" {
		t.Errorf("Message = %q, want %q", issue.Message, "invalid syntax")
	}
	if issue.RuleID != "syntax-check" {
		t.Errorf("RuleID = %q, want %q", issue.RuleID, "syntax-check")
	}
}

func TestValidatorPipeline_Validate(t *testing.T) {
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	pipeline := NewValidatorPipeline(NewActionlintValidator())
	result, err := pipeline.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if !result.Success {
		t.Error("Pipeline.Validate() Success = false, want true")
	}
}

func TestValidatorPipeline_Empty(t *testing.T) {
	pipeline := NewValidatorPipeline()
	result, err := pipeline.Validate("ci.yml", []byte{})
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	// Empty pipeline should always succeed
	if !result.Success {
		t.Error("Empty pipeline.Validate() Success = false, want true")
	}

	if len(result.Issues) != 0 {
		t.Errorf("Empty pipeline expected 0 issues, got %d", len(result.Issues))
	}
}

func TestValidatorPipeline_Multiple(t *testing.T) {
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`)

	// Create pipeline with multiple validators
	pipeline := NewValidatorPipeline(
		NewActionlintValidator(),
		NewActionlintValidator(), // Same validator twice to test multiple validators
	)

	result, err := pipeline.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if !result.Success {
		t.Error("Pipeline.Validate() with multiple validators failed")
	}
}

func TestValidatorPipeline_MultipleWithErrors(t *testing.T) {
	// Invalid workflow
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    invalid-key: true
    steps:
      - uses: actions/checkout@v4
`)

	// Create pipeline with multiple validators
	pipeline := NewValidatorPipeline(
		NewActionlintValidator(),
		NewActionlintValidator(),
	)

	result, err := pipeline.Validate("ci.yml", content)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if result.Success {
		t.Error("Pipeline.Validate() Success = true, want false for invalid workflow")
	}

	// Should have issues from both validators
	if len(result.Issues) == 0 {
		t.Error("Pipeline.Validate() expected issues from validators")
	}
}

func TestValidatorPipeline_ErrorPropagation(t *testing.T) {
	// Create a pipeline with a validator
	pipeline := NewValidatorPipeline(NewActionlintValidator())

	// Try to validate malformed YAML that should cause an error
	content := []byte(`{invalid}`)

	result, err := pipeline.Validate("bad.yml", content)

	// The pipeline should propagate errors from validators
	if err == nil && (result == nil || result.Success) {
		t.Error("Pipeline.Validate() expected error or failed validation for malformed YAML")
	}
}

func TestValidateWorkflow(t *testing.T) {
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo "test"
`)

	result, err := ValidateWorkflow("test.yml", content)
	if err != nil {
		t.Fatalf("ValidateWorkflow() error = %v", err)
	}

	if !result.Success {
		t.Error("ValidateWorkflow() Success = false, want true")
	}
}

func TestValidateWorkflow_Invalid(t *testing.T) {
	content := []byte(`name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    unknown-field: value
    steps:
      - run: echo "test"
`)

	result, err := ValidateWorkflow("test.yml", content)
	if err != nil {
		t.Fatalf("ValidateWorkflow() error = %v", err)
	}

	if result.Success {
		t.Error("ValidateWorkflow() Success = true, want false for invalid workflow")
	}

	if len(result.Issues) == 0 {
		t.Error("ValidateWorkflow() expected issues for invalid workflow")
	}
}

func TestValidateWorkflowFile(t *testing.T) {
	// Create a temporary file with valid workflow
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "test-workflow.yml")

	content := []byte(`name: Test
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "test"
`)

	if err := os.WriteFile(workflowPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ValidateWorkflowFile(workflowPath)
	if err != nil {
		t.Fatalf("ValidateWorkflowFile() error = %v", err)
	}

	if !result.Success {
		t.Errorf("ValidateWorkflowFile() Success = false, want true")
		for _, issue := range result.Issues {
			t.Logf("Issue: %s:%d:%d: %s", issue.File, issue.Line, issue.Column, issue.Message)
		}
	}
}

func TestValidateWorkflowFile_Invalid(t *testing.T) {
	// Create a temporary file with invalid workflow
	tmpDir := t.TempDir()
	workflowPath := filepath.Join(tmpDir, "invalid.yml")

	content := []byte(`name: Test
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    bad-key: value
    steps:
      - run: echo "test"
`)

	if err := os.WriteFile(workflowPath, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ValidateWorkflowFile(workflowPath)
	if err != nil {
		t.Fatalf("ValidateWorkflowFile() error = %v", err)
	}

	if result.Success {
		t.Error("ValidateWorkflowFile() Success = true, want false for invalid workflow")
	}

	if len(result.Issues) == 0 {
		t.Error("ValidateWorkflowFile() expected issues for invalid workflow")
	}
}

func TestValidateWorkflowFile_NonExistent(t *testing.T) {
	result, err := ValidateWorkflowFile("/nonexistent/file.yml")

	// Should return an error for non-existent file
	if err == nil {
		t.Error("ValidateWorkflowFile() expected error for non-existent file")
	}

	// Result could be nil when error is returned
	if result != nil && result.Success {
		t.Error("ValidateWorkflowFile() Success = true, want error for non-existent file")
	}
}

func TestValidationResult_NoIssues(t *testing.T) {
	result := &ValidationResult{
		Success: true,
		Issues:  []ValidationIssue{},
	}

	if !result.Success {
		t.Error("ValidationResult.Success should be true")
	}

	if len(result.Issues) != 0 {
		t.Errorf("ValidationResult.Issues should be empty, got %d issues", len(result.Issues))
	}
}

func TestValidationResult_WithIssues(t *testing.T) {
	issues := []ValidationIssue{
		{
			File:    "test.yml",
			Line:    1,
			Column:  5,
			Message: "error 1",
			RuleID:  "rule1",
		},
		{
			File:    "test.yml",
			Line:    2,
			Column:  10,
			Message: "error 2",
			RuleID:  "rule2",
		},
	}

	result := &ValidationResult{
		Success: false,
		Issues:  issues,
	}

	if result.Success {
		t.Error("ValidationResult.Success should be false when there are issues")
	}

	if len(result.Issues) != 2 {
		t.Errorf("ValidationResult.Issues should have 2 issues, got %d", len(result.Issues))
	}

	// Verify issue details
	if result.Issues[0].Message != "error 1" {
		t.Errorf("First issue message = %q, want %q", result.Issues[0].Message, "error 1")
	}
	if result.Issues[1].Message != "error 2" {
		t.Errorf("Second issue message = %q, want %q", result.Issues[1].Message, "error 2")
	}
}
