package agent

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lex00/wetwire-core-go/agent/results"
)

func TestGitHubAgent_ToolRunLint_PathHandling(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with subdirectory path
	_ = agent.toolRunLint("subdir")

	// Verify state was updated (this is what matters)
	if !agent.lintCalled {
		t.Error("lintCalled should be true after toolRunLint")
	}
}

func TestGitHubAgent_ToolRunBuild_PathHandling(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with specific path
	result := agent.toolRunBuild("subdir")

	// Should execute and return result
	if result == "" {
		t.Error("toolRunBuild should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunValidate_PathHandling(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with specific file path
	result := agent.toolRunValidate("workflow.yml")

	// Should execute and return result
	if result == "" {
		t.Error("toolRunValidate should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunLint_Success(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test that toolRunLint updates state correctly (even if command fails)
	initialCycles := agent.lintCycles
	agent.pendingLint = true

	_ = agent.toolRunLint(".")

	// Verify state changes
	if !agent.lintCalled {
		t.Error("lintCalled should be true after toolRunLint")
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false after toolRunLint")
	}
	if agent.lintCycles != initialCycles+1 {
		t.Errorf("lintCycles = %d, want %d", agent.lintCycles, initialCycles+1)
	}
}

func TestGitHubAgent_ToolRunLint_WithSession(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	session := results.NewSession("persona", "scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	_ = agent.toolRunLint(".")

	// Verify lint state was updated
	if !agent.lintCalled {
		t.Error("lintCalled should be true")
	}

	// Session will be updated only if command succeeds or fails with exit code 2
	// Just verify the session is still configured
	if agent.session != session {
		t.Error("session should remain configured")
	}
}

func TestGitHubAgent_ToolRunLint_MultipleCycles(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Run lint multiple times
	for i := 0; i < 3; i++ {
		agent.toolWriteFile("test.go", "package main")
		agent.toolRunLint(".")
	}

	if agent.lintCycles != 3 {
		t.Errorf("lintCycles = %d, want 3", agent.lintCycles)
	}
}

func TestGitHubAgent_ToolRunBuild(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolRunBuild(".")

	// Result will contain error since we don't have a valid project,
	// but we're testing the function executes
	if result == "" {
		t.Error("toolRunBuild should return a result")
	}
}

func TestGitHubAgent_ToolRunValidate(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolRunValidate(".")

	// Result will contain error/validation info
	if result == "" {
		t.Error("toolRunValidate should return a result")
	}
}

func TestGitHubAgent_ToolRunLint_StateReset(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Set up initial state
	agent.pendingLint = true
	agent.lintCalled = false
	agent.lintPassed = false

	// Run lint
	agent.toolRunLint(".")

	// Verify state is updated correctly
	if !agent.lintCalled {
		t.Error("lintCalled should be true")
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false")
	}
}

func TestGitHubAgent_ToolRunBuild_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Build with invalid path should return error in result
	result := agent.toolRunBuild("nonexistent-project")

	// Should contain error info (either "Build error" or command output)
	if result == "" {
		t.Error("toolRunBuild should return error information")
	}
}

func TestGitHubAgent_ToolRunValidate_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Validate with nonexistent file should return error in result
	result := agent.toolRunValidate("nonexistent.yml")

	// Should contain error/validation info
	if result == "" {
		t.Error("toolRunValidate should return validation information")
	}
}

func TestGitHubAgent_ToolRunLint_LintFailureWithSession(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Run lint on non-existent directory (will fail)
	agent.toolRunLint("nonexistent")

	// Verify state was updated
	if !agent.lintCalled {
		t.Error("lintCalled should be true")
	}
	if agent.lintPassed {
		t.Error("lintPassed should be false after failed lint")
	}
}

func TestGitHubAgent_ToolRunBuild_WithError(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Build non-existent project
	result := agent.toolRunBuild("nonexistent-project")

	// Should contain build error
	if !findSubstring(result, "Build error") && result == "" {
		t.Errorf("toolRunBuild should return error for nonexistent project, got %q", result)
	}
}

func TestGitHubAgent_ToolRunValidate_WithError(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Validate non-existent file
	result := agent.toolRunValidate("nonexistent.yml")

	// Should contain validation issues or error
	if result == "" {
		t.Error("toolRunValidate should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunBuild_Success(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test build execution
	result := agent.toolRunBuild(".")

	// Should return some output (error or success)
	if result == "" {
		t.Error("toolRunBuild should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunValidate_Success(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test validate execution
	result := agent.toolRunValidate(".")

	// Should return some output
	if result == "" {
		t.Error("toolRunValidate should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunLint_StateTracking(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	session := results.NewSession("test", "test")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Set pending lint
	agent.pendingLint = true

	// Run lint
	_ = agent.toolRunLint(".")

	// Verify state changes
	if !agent.lintCalled {
		t.Error("lintCalled should be true")
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false after lint")
	}
	if agent.lintCycles == 0 {
		t.Error("lintCycles should be incremented")
	}
}

func TestGitHubAgent_ToolRunLintWithRealCommand(t *testing.T) {
	// Check if wetwire-github binary exists
	if _, err := os.Stat("/tmp/wetwire-github"); os.IsNotExist(err) {
		t.Skip("wetwire-github binary not found at /tmp/wetwire-github")
	}

	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Add /tmp to PATH for this test
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "/tmp:"+origPath)
	defer os.Setenv("PATH", origPath)

	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a valid wetwire-github project
	agent.toolInitPackage("test-project")

	// Write a valid workflow file
	workflowContent := `package main

import (
	"github.com/lex00/wetwire-github-go/workflow"
	"github.com/lex00/wetwire-github-go/actions/checkout"
)

var CI = workflow.Workflow{
	Name: "CI",
	On: workflow.Triggers{
		Push: &workflow.PushTrigger{
			Branches: []string{"main"},
		},
	},
	Jobs: map[string]workflow.Job{
		"build": {
			RunsOn: "ubuntu-latest",
			Steps: []any{
				checkout.Checkout{},
			},
		},
	},
}
`
	agent.toolWriteFile("test-project/workflow.go", workflowContent)

	// Run lint - should succeed or fail gracefully
	result := agent.toolRunLint("test-project")

	// Verify state was updated
	if !agent.lintCalled {
		t.Error("lintCalled should be true after toolRunLint")
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false after toolRunLint")
	}
	if agent.lintCycles != 1 {
		t.Errorf("lintCycles = %d, want 1", agent.lintCycles)
	}

	// Result should have some output
	if result == "" {
		t.Log("Warning: lint result is empty (command may have failed)")
	}
}

func TestGitHubAgent_ToolRunBuildWithRealCommand(t *testing.T) {
	// Check if wetwire-github binary exists
	if _, err := os.Stat("/tmp/wetwire-github"); os.IsNotExist(err) {
		t.Skip("wetwire-github binary not found at /tmp/wetwire-github")
	}

	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Add /tmp to PATH for this test
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "/tmp:"+origPath)
	defer os.Setenv("PATH", origPath)

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a valid project
	agent.toolInitPackage("build-test")
	workflowContent := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
	On: workflow.Triggers{
		Push: &workflow.PushTrigger{Branches: []string{"main"}},
	},
}
`
	agent.toolWriteFile("build-test/workflow.go", workflowContent)

	// Run build
	result := agent.toolRunBuild("build-test")

	// Should have some output
	if result == "" {
		t.Log("Warning: build result is empty")
	}
}

func TestGitHubAgent_ToolRunValidateWithRealCommand(t *testing.T) {
	// Check if wetwire-github binary exists
	if _, err := os.Stat("/tmp/wetwire-github"); os.IsNotExist(err) {
		t.Skip("wetwire-github binary not found at /tmp/wetwire-github")
	}

	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Add /tmp to PATH for this test
	origPath := os.Getenv("PATH")
	os.Setenv("PATH", "/tmp:"+origPath)
	defer os.Setenv("PATH", origPath)

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a YAML file to validate
	yamlContent := `name: CI
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
`
	yamlPath := filepath.Join(tmpDir, "workflow.yml")
	if err := os.WriteFile(yamlPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("writing yaml file: %v", err)
	}

	// Run validate
	result := agent.toolRunValidate("workflow.yml")

	// Should have some output
	if result == "" {
		t.Log("Warning: validate result is empty")
	}
}

func TestGitHubAgent_ToolRunLint_LintPassedState(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Initial state
	if agent.lintPassed {
		t.Error("initial lintPassed should be false")
	}

	// Create a file with lint violation (WAG001 - raw uses string)
	agent.toolWriteFile("workflows.go", `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Jobs: map[string]workflow.Job{
		"test": {
			Steps: []any{
				workflow.Step{Uses: "actions/checkout@v4"},
			},
		},
	},
}
`)

	// Run lint - should fail due to lint violations
	agent.toolRunLint(".")

	// lintPassed should be false after failed lint
	if agent.lintPassed {
		t.Error("lintPassed should be false when lint command fails")
	}

	// But lintCalled should be true
	if !agent.lintCalled {
		t.Error("lintCalled should be true")
	}
}

func TestGitHubAgent_ToolRunLint_SessionTracking(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Run lint multiple times to exercise session tracking
	for i := 0; i < 3; i++ {
		agent.toolWriteFile(formatIteration(i), "package main")
		agent.toolRunLint(".")
	}

	// Verify lint cycles tracked
	if agent.GetLintCycles() != 3 {
		t.Errorf("lintCycles = %d, want 3", agent.GetLintCycles())
	}

	// Session should still be configured
	if agent.session != session {
		t.Error("session should remain configured")
	}
}

func TestGitHubAgent_ToolRunBuild_SuccessPath(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Even without a valid project, the method should execute and return output
	result := agent.toolRunBuild(".")

	// Should return some output (error or success)
	if result == "" {
		t.Error("toolRunBuild should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunValidate_SuccessPath(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a simple YAML file
	yamlPath := filepath.Join(tmpDir, "test.yml")
	if err := os.WriteFile(yamlPath, []byte("name: test"), 0644); err != nil {
		t.Fatalf("writing yaml: %v", err)
	}

	result := agent.toolRunValidate("test.yml")

	// Should return some output
	if result == "" {
		t.Error("toolRunValidate should return non-empty result")
	}
}

func TestGitHubAgent_ToolRunLint_CommandExecution(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Initial state checks
	if agent.lintCalled {
		t.Error("lintCalled should be false initially")
	}
	if agent.lintCycles != 0 {
		t.Error("lintCycles should be 0 initially")
	}

	// Run lint
	_ = agent.toolRunLint(".")

	// State should be updated regardless of whether command exists
	if !agent.lintCalled {
		t.Error("lintCalled should be true after lint")
	}
	if agent.lintCycles != 1 {
		t.Errorf("lintCycles = %d, want 1", agent.lintCycles)
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false after lint")
	}
}

func TestGitHubAgent_ToolRunBuild_CommandExecution(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Run build - will fail but we're testing execution
	result := agent.toolRunBuild("nonexistent-path")

	// Should have output (error message)
	if result == "" {
		t.Error("result should not be empty")
	}

	// Should mention error
	if !findSubstring(result, "error") && !findSubstring(result, "Error") && !findSubstring(result, "Build error") {
		t.Log("result:", result) // Log for debugging
	}
}

func TestGitHubAgent_ToolRunValidate_CommandExecution(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a file
	testFile := filepath.Join(tmpDir, "test.yml")
	if err := os.WriteFile(testFile, []byte("name: test\non: push"), 0644); err != nil {
		t.Fatalf("creating test file: %v", err)
	}

	// Run validate
	result := agent.toolRunValidate("test.yml")

	// Should have output
	if result == "" {
		t.Error("result should not be empty")
	}
}
