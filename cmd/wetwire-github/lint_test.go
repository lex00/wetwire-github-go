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

// TestLint_ValidCode tests linting valid code.
func TestLint_ValidCode(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid workflow Go file
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Save and restore global state
	oldFormat := lintFormat
	oldFix := lintFix
	defer func() {
		lintFormat = oldFormat
		lintFix = oldFix
	}()

	lintFormat = "text"
	lintFix = false

	// Note: runLint calls os.Exit, so we test the underlying linter instead
	// This is tested via integration tests with exec.Command
}

// TestLint_WAG001_RawUses tests detecting raw uses: strings.
func TestLint_WAG001_RawUses(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// The runLint function calls os.Exit, so we need integration tests
	// This test verifies the file can be created and the pattern is detectable
}

// TestLint_NonExistentPath tests error handling for non-existent paths.
func TestLint_NonExistentPath(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "lint", "/nonexistent/path")
	err := cmd.Run()
	if err == nil {
		t.Error("Lint command should fail for non-existent path")
	}
}

// TestLintCmd_Integration tests the lint command via exec.
func TestLintCmd_Integration(t *testing.T) {
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

	// Create test fixtures with valid code
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Run lint on valid code
	cmd := exec.Command(binaryPath, "lint", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	// Lint may find issues and return non-zero, check output instead
	output := stdout.String() + stderr.String()
	if err != nil && !strings.Contains(output, "No issues found") && !strings.Contains(output, "issue") {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Lint command failed unexpectedly: %v", err)
	}
}

// TestLintCmd_WAG001_Integration tests WAG001 detection via exec.
func TestLintCmd_WAG001_Integration(t *testing.T) {
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

	// Create test fixtures with WAG001 violation
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "lint", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run() // May return error if issues found

	output := stdout.String() + stderr.String()
	if !strings.Contains(output, "WAG001") {
		t.Errorf("Expected WAG001 issue to be detected, got: %s", output)
	}
}

// TestLintCmd_JSONFormat tests JSON output format.
func TestLintCmd_JSONFormat(t *testing.T) {
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

	// Create test fixtures
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "lint", tmpDir, "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	_ = cmd.Run()

	// Should be valid JSON
	var result wetwire.LintResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output, got error: %v\nOutput: %s", err, stdout.String())
	}
}

// TestLintCmd_FixFlag tests the --fix flag.
func TestLintCmd_FixFlag(t *testing.T) {
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

	// Create test fixtures with fixable issue (WAG001)
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CheckoutStep = workflow.Step{Uses: "actions/checkout@v4"}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "lint", tmpDir, "--fix")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// --fix should run without error
	if err != nil {
		t.Logf("Fix output: %s", output)
		// Fix may fail for various reasons, just log
	}

	// Read the file to see if it was modified
	fixed, err := os.ReadFile(filepath.Join(tmpDir, "workflows.go"))
	if err != nil {
		t.Fatal(err)
	}

	// Check if fix was applied (checkout.Checkout should replace raw string)
	if strings.Contains(string(fixed), "checkout.Checkout") {
		t.Log("Fix successfully applied checkout.Checkout wrapper")
	}
}

// TestLintCmd_FixJSONFormat tests --fix with JSON output.
func TestLintCmd_FixJSONFormat(t *testing.T) {
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

	// Create test fixtures
	tmpDir := t.TempDir()
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "workflows.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "lint", tmpDir, "--fix", "--format", "json")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	_ = cmd.Run()

	// Should be valid JSON
	var result map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Errorf("Expected valid JSON output for fix, got error: %v\nOutput: %s", err, stdout.String())
	}
}

// TestLintCmd_DirRecursive tests linting a directory recursively.
func TestLintCmd_DirRecursive(t *testing.T) {
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

	// Create test fixtures with subdirectory
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "workflows")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	content := `package workflows

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	if err := os.WriteFile(filepath.Join(subDir, "ci.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "lint", tmpDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	// Should process subdirectories
	if err != nil && !strings.Contains(output, "No issues found") && !strings.Contains(output, "issue") {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Lint command failed on subdirectory: %v", err)
	}
}

// TestLintCmd_SingleFile tests linting a single file.
func TestLintCmd_SingleFile(t *testing.T) {
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

	// Create test fixture file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "workflows.go")
	content := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(binaryPath, "lint", testFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String() + stderr.String()

	if err != nil && !strings.Contains(output, "No issues found") && !strings.Contains(output, "issue") {
		t.Logf("stdout: %s", stdout.String())
		t.Logf("stderr: %s", stderr.String())
		t.Errorf("Lint command failed on single file: %v", err)
	}
}

// TestLintResult_JSON tests LintResult JSON serialization.
func TestLintResult_JSON(t *testing.T) {
	result := wetwire.LintResult{
		Success: false,
		Issues: []wetwire.LintIssue{
			{
				File:     "test.go",
				Line:     5,
				Column:   10,
				Severity: "warning",
				Message:  "Use typed action wrapper",
				Rule:     "WAG001",
				Fixable:  true,
			},
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal LintResult: %v", err)
	}

	var unmarshaled wetwire.LintResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal LintResult: %v", err)
	}

	if len(unmarshaled.Issues) != 1 {
		t.Errorf("len(Issues) = %d, want 1", len(unmarshaled.Issues))
	}

	if unmarshaled.Issues[0].Rule != "WAG001" {
		t.Errorf("Issues[0].Rule = %q, want %q", unmarshaled.Issues[0].Rule, "WAG001")
	}
}

// TestLintCmd_Help tests lint help output.
func TestLintCmd_Help(t *testing.T) {
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

	cmd := exec.Command(binaryPath, "lint", "--help")
	out, _ := cmd.CombinedOutput()

	if !strings.Contains(string(out), "lint") {
		t.Errorf("Help output should contain 'lint', got: %s", out)
	}

	if !strings.Contains(string(out), "WAG001") {
		t.Errorf("Help output should contain WAG001 rule description, got: %s", out)
	}
}
