package validation

import (
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

func TestNewActionlintValidator(t *testing.T) {
	v := NewActionlintValidator()
	if v == nil {
		t.Error("NewActionlintValidator() returned nil")
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
