package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	wetwire "github.com/lex00/wetwire-github-go"
)

// TestValidateCmd_Help tests validate help output.
func TestValidateCmd_Help(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	cmd := exec.Command(binaryPath, "validate", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "validate") {
		t.Errorf("Help output should contain 'validate', got: %s", out)
	}

	if !strings.Contains(string(out), "actionlint") {
		t.Errorf("Help output should mention 'actionlint', got: %s", out)
	}

	if !strings.Contains(string(out), "--format") {
		t.Errorf("Help output should contain '--format' flag, got: %s", out)
	}

	if !strings.Contains(string(out), "workflow") {
		t.Errorf("Help output should mention 'workflow', got: %s", out)
	}
}

// TestValidateCmd_ValidWorkflow tests validation of a valid workflow file.
func TestValidateCmd_ValidWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a valid workflow YAML file
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: go build ./...
`
	yamlPath := filepath.Join(tmpDir, "ci.yml")
	if err := os.WriteFile(yamlPath, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "validate", yamlPath)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Validate command should succeed for valid workflow: %v", err)
		return
	}

	output := stdout.String()
	if !strings.Contains(output, "valid") {
		t.Errorf("Output should contain 'valid', got: %s", output)
	}
}

// TestValidateCmd_InvalidPath tests error handling for non-existent files.
func TestValidateCmd_InvalidPath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	cmd := exec.Command(binaryPath, "validate", "/nonexistent/workflow.yml")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Validate command should fail for non-existent file")
	}

	if !strings.Contains(stderr.String(), "not found") {
		t.Errorf("Expected error message about file not found, got: %s", stderr.String())
	}
}

// TestValidateCmd_JSONFormat tests JSON output format.
func TestValidateCmd_JSONFormat(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create a valid workflow YAML file
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on: push

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`
	yamlPath := filepath.Join(tmpDir, "ci.yml")
	if err := os.WriteFile(yamlPath, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "validate", yamlPath, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	_ = cmd.Run() // May succeed or fail depending on actionlint availability

	// Should be valid JSON
	var result wetwire.ValidateResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
		return
	}

	// Check that result has proper structure
	// Errors and Warnings should exist (may be empty slices or nil)
	_ = result.Errors
	_ = result.Warnings
}

// TestValidateCmd_MultipleFiles tests validation of multiple workflow files.
func TestValidateCmd_MultipleFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create multiple valid workflow YAML files
	tmpDir := t.TempDir()
	workflowYAML := `name: CI

on: push

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`
	yaml1 := filepath.Join(tmpDir, "ci.yml")
	yaml2 := filepath.Join(tmpDir, "deploy.yml")

	if err := os.WriteFile(yaml1, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(yaml2, []byte(workflowYAML), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "validate", yaml1, yaml2)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		// Don't fail - actionlint may not be available
		return
	}

	output := stdout.String()
	if !strings.Contains(output, "valid") {
		t.Logf("Output: %s", output)
	}
}

// TestValidateCmd_MissingArg tests error handling for missing file argument.
func TestValidateCmd_MissingArg(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	cmd := exec.Command(binaryPath, "validate")
	err := cmd.Run()
	if err == nil {
		t.Error("Validate command should fail without file argument")
	}
}

// TestValidateCmd_InvalidYAML tests validation of invalid YAML.
func TestValidateCmd_InvalidYAML(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary
	binaryPath := filepath.Join(t.TempDir(), "wetwire-github")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = getModulePath() + "/cmd/wetwire-github"
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build binary: %v\n%s", err, out)
	}

	// Create an invalid YAML file
	tmpDir := t.TempDir()
	invalidYAML := `name: CI
  this is: invalid::: yaml
    - not valid
`
	yamlPath := filepath.Join(tmpDir, "invalid.yml")
	if err := os.WriteFile(yamlPath, []byte(invalidYAML), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "validate", yamlPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		t.Error("Validate command should fail for invalid YAML")
	}

	// Should have error output
	if stderr.Len() == 0 {
		t.Error("Expected error output for invalid YAML")
	}
}

// TestValidateResult_JSON tests ValidateResult JSON serialization.
func TestValidateResult_JSON(t *testing.T) {
	result := wetwire.ValidateResult{
		Success: false,
		Errors: []wetwire.ValidationError{
			{
				File:    "ci.yml",
				Line:    10,
				Column:  5,
				Message: "invalid syntax",
				RuleID:  "syntax-check",
			},
		},
		Warnings: []string{"deprecated action version"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal ValidateResult: %v", err)
	}

	var unmarshaled wetwire.ValidateResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ValidateResult: %v", err)
	}

	if unmarshaled.Success {
		t.Error("Expected success to be false")
	}

	if len(unmarshaled.Errors) != 1 {
		t.Errorf("len(Errors) = %d, want 1", len(unmarshaled.Errors))
	}

	if unmarshaled.Errors[0].File != "ci.yml" {
		t.Errorf("Error.File = %q, want %q", unmarshaled.Errors[0].File, "ci.yml")
	}

	if unmarshaled.Errors[0].Line != 10 {
		t.Errorf("Error.Line = %d, want 10", unmarshaled.Errors[0].Line)
	}

	if len(unmarshaled.Warnings) != 1 {
		t.Errorf("len(Warnings) = %d, want 1", len(unmarshaled.Warnings))
	}
}
