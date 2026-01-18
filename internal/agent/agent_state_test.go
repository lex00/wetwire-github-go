package agent

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lex00/wetwire-core-go/providers"
)

func TestGitHubAgent_GetGeneratedFiles(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	files := agent.GetGeneratedFiles()
	// Initially nil is acceptable - will be populated as files are generated
	if len(files) != 0 {
		t.Errorf("len(GetGeneratedFiles()) = %d, want 0", len(files))
	}
}

func TestGitHubAgent_GetLintCycles(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	cycles := agent.GetLintCycles()
	if cycles != 0 {
		t.Errorf("GetLintCycles() = %d, want 0", cycles)
	}
}

func TestGitHubAgent_LintPassed(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	passed := agent.LintPassed()
	if passed {
		t.Error("LintPassed() = true, want false (initial)")
	}
}

func TestGitHubAgent_GetGeneratedFiles_AfterWrites(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Write multiple files
	agent.toolWriteFile("file1.go", "package main")
	agent.toolWriteFile("file2.go", "package main")
	agent.toolWriteFile("subdir/file3.go", "package subdir")

	files := agent.GetGeneratedFiles()

	if len(files) != 3 {
		t.Errorf("GetGeneratedFiles() returned %d files, want 3", len(files))
	}

	// Verify all files are tracked
	expected := map[string]bool{
		"file1.go":        false,
		"file2.go":        false,
		"subdir/file3.go": false,
	}

	for _, f := range files {
		if _, ok := expected[f]; ok {
			expected[f] = true
		}
	}

	for file, found := range expected {
		if !found {
			t.Errorf("GetGeneratedFiles() missing %q", file)
		}
	}
}

func TestGitHubAgent_StateTransitions(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Initial state
	if agent.lintCalled {
		t.Error("initial lintCalled should be false")
	}
	if agent.lintPassed {
		t.Error("initial lintPassed should be false")
	}
	if agent.pendingLint {
		t.Error("initial pendingLint should be false")
	}

	// Write a file - should set pendingLint
	agent.toolWriteFile("test.go", "package main")
	if !agent.pendingLint {
		t.Error("pendingLint should be true after write")
	}
	if agent.lintPassed {
		t.Error("lintPassed should remain false after write")
	}

	// Run lint - should clear pendingLint and set lintCalled
	agent.toolRunLint(".")
	if agent.pendingLint {
		t.Error("pendingLint should be false after lint")
	}
	if !agent.lintCalled {
		t.Error("lintCalled should be true after lint")
	}

	// Write another file - should reset lintPassed and set pendingLint
	agent.lintPassed = true // Simulate successful lint
	agent.toolWriteFile("test2.go", "package main")
	if agent.lintPassed {
		t.Error("lintPassed should be reset to false after new write")
	}
	if !agent.pendingLint {
		t.Error("pendingLint should be true after new write")
	}
}

func TestGitHubAgent_LintCyclesIncrement(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir, MaxLintCycles: 3})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Check maxLintCycles is set
	if agent.maxLintCycles != 3 {
		t.Errorf("maxLintCycles = %d, want 3", agent.maxLintCycles)
	}

	// Run lint multiple times and verify cycle count
	for i := 1; i <= 5; i++ {
		agent.toolRunLint(".")
		if agent.GetLintCycles() != i {
			t.Errorf("after %d lints, GetLintCycles() = %d, want %d", i, agent.GetLintCycles(), i)
		}
	}
}

func TestGitHubAgent_GetGeneratedFiles_Order(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Write files in specific order
	files := []string{"first.go", "second.go", "third.go"}
	for _, f := range files {
		agent.toolWriteFile(f, "package main")
	}

	generated := agent.GetGeneratedFiles()

	// Verify all files are tracked in order
	if len(generated) != len(files) {
		t.Fatalf("GetGeneratedFiles() returned %d files, want %d", len(generated), len(files))
	}

	for i, expected := range files {
		if generated[i] != expected {
			t.Errorf("GetGeneratedFiles()[%d] = %q, want %q", i, generated[i], expected)
		}
	}
}

func TestGitHubAgent_StateResetAfterWrite(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Set up state as if lint passed
	agent.lintCalled = true
	agent.lintPassed = true
	agent.pendingLint = false

	// Write a file
	agent.toolWriteFile("new.go", "package main")

	// Verify state was reset appropriately
	if !agent.pendingLint {
		t.Error("pendingLint should be true after write")
	}
	if agent.lintPassed {
		t.Error("lintPassed should be false after new write")
	}
	// lintCalled should remain true (it was already called once)
	if !agent.lintCalled {
		t.Error("lintCalled should remain true")
	}
}

func TestGitHubAgent_MultipleStateResets(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Cycle through write -> lint -> write -> lint multiple times
	for i := 0; i < 3; i++ {
		// Write
		agent.toolWriteFile(fmt.Sprintf("file%d.go", i), "package main")

		if !agent.pendingLint {
			t.Errorf("iteration %d: pendingLint should be true after write", i)
		}
		if agent.lintPassed {
			t.Errorf("iteration %d: lintPassed should be false after write", i)
		}

		// Lint
		agent.toolRunLint(".")

		if agent.pendingLint {
			t.Errorf("iteration %d: pendingLint should be false after lint", i)
		}
		if !agent.lintCalled {
			t.Errorf("iteration %d: lintCalled should be true", i)
		}
	}

	// Verify final state
	if len(agent.GetGeneratedFiles()) != 3 {
		t.Errorf("should have 3 files, got %d", len(agent.GetGeneratedFiles()))
	}
	if agent.GetLintCycles() != 3 {
		t.Errorf("should have 3 lint cycles, got %d", agent.GetLintCycles())
	}
}

func TestGitHubAgent_LintPassedAfterSuccessfulLint(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Initially should not have passed
	if agent.LintPassed() {
		t.Error("LintPassed() should be false initially")
	}

	// After running lint (will fail but updates state)
	agent.toolRunLint(".")

	// Check cycles incremented
	if agent.GetLintCycles() != 1 {
		t.Errorf("GetLintCycles() = %d, want 1", agent.GetLintCycles())
	}
}

func TestGitHubAgent_Run_ContextCancellation(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Run should return context error
	err = agent.Run(ctx, "generate a workflow")
	if err != context.Canceled {
		t.Errorf("Run() with cancelled context = %v, want %v", err, context.Canceled)
	}
}

func TestGitHubAgent_ToolExecutionFlow(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Simulate a complete workflow through tools
	// 1. Init package
	result := agent.toolInitPackage("test-workflow")
	if !findSubstring(result, "Created project") {
		t.Errorf("toolInitPackage failed: %s", result)
	}

	// 2. Write a workflow file
	workflowContent := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	result = agent.toolWriteFile("test-workflow/workflow.go", workflowContent)
	if !findSubstring(result, "Wrote") {
		t.Errorf("toolWriteFile failed: %s", result)
	}

	// Verify state after write
	if !agent.pendingLint {
		t.Error("pendingLint should be true after write")
	}
	if agent.lintPassed {
		t.Error("lintPassed should be false after write")
	}

	// 3. Read it back
	result = agent.toolReadFile("test-workflow/workflow.go")
	if result != workflowContent {
		t.Errorf("toolReadFile returned incorrect content")
	}

	// 4. Run lint (will fail but updates state)
	result = agent.toolRunLint("test-workflow")
	// State should be updated
	if !agent.lintCalled {
		t.Error("lintCalled should be true after lint")
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false after lint")
	}

	// 5. Write another file
	result = agent.toolWriteFile("test-workflow/jobs.go", "package main")
	if len(agent.GetGeneratedFiles()) != 2 {
		t.Errorf("should have 2 generated files, got %d", len(agent.GetGeneratedFiles()))
	}

	// 6. Check lint enforcement logic
	enforcement := agent.checkLintEnforcement([]string{"write_file"})
	if enforcement == "" {
		t.Error("checkLintEnforcement should enforce after write without lint")
	}
}

func TestGitHubAgent_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name  string
		setup func(*GitHubAgent)
		test  func(*testing.T, *GitHubAgent)
	}{
		{
			name:  "write same file twice",
			setup: func(a *GitHubAgent) {},
			test: func(t *testing.T, a *GitHubAgent) {
				a.toolWriteFile("test.go", "package main")
				a.toolWriteFile("test.go", "package main\n// modified")
				files := a.GetGeneratedFiles()
				if len(files) != 2 {
					t.Errorf("should track both writes, got %d", len(files))
				}
			},
		},
		{
			name:  "lint without any files",
			setup: func(a *GitHubAgent) {},
			test: func(t *testing.T, a *GitHubAgent) {
				a.toolRunLint(".")
				if !a.lintCalled {
					t.Error("lintCalled should be true even with no files")
				}
			},
		},
		{
			name: "multiple lint runs",
			setup: func(a *GitHubAgent) {
				a.toolWriteFile("test.go", "package main")
			},
			test: func(t *testing.T, a *GitHubAgent) {
				for i := 0; i < 10; i++ {
					a.toolRunLint(".")
				}
				if a.GetLintCycles() != 10 {
					t.Errorf("lintCycles = %d, want 10", a.GetLintCycles())
				}
			},
		},
		{
			name:  "read non-existent file",
			setup: func(a *GitHubAgent) {},
			test: func(t *testing.T, a *GitHubAgent) {
				result := a.toolReadFile("does-not-exist.go")
				if !findSubstring(result, "Error") {
					t.Errorf("should return error for non-existent file, got %q", result)
				}
			},
		},
		{
			name:  "init package with empty name",
			setup: func(a *GitHubAgent) {},
			test: func(t *testing.T, a *GitHubAgent) {
				result := a.toolInitPackage("")
				// Should still work, creates directory with empty name
				if findSubstring(result, "Error") {
					t.Errorf("unexpected error: %s", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}
			tt.setup(agent)
			tt.test(t, agent)
		})
	}
}

func TestGitHubAgent_ComprehensiveIntegration(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{
		WorkDir:       tmpDir,
		MaxLintCycles: 3,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Simulate a complete workflow without API

	// 1. Init package
	result := agent.toolInitPackage("test-workflow")
	if !findSubstring(result, "Created") {
		t.Errorf("init failed: %s", result)
	}

	// 2. Write file with lint violation (WAG001 - raw uses string)
	result = agent.toolWriteFile("test-workflow/main.go", `package main

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
	if !findSubstring(result, "Wrote") {
		t.Errorf("write failed: %s", result)
	}

	// Verify state after write
	if len(agent.GetGeneratedFiles()) != 1 {
		t.Error("generatedFiles should have 1 file")
	}
	if !agent.pendingLint {
		t.Error("pendingLint should be true")
	}

	// 3. Read file
	result = agent.toolReadFile("test-workflow/main.go")
	if !findSubstring(result, "package main") {
		t.Errorf("read returned wrong content: %s", result)
	}

	// 4. Run lint (will fail but updates state)
	result = agent.toolRunLint("test-workflow")
	if !agent.lintCalled {
		t.Error("lintCalled should be true")
	}
	if agent.pendingLint {
		t.Error("pendingLint should be false after lint")
	}

	// 5. Write another file
	result = agent.toolWriteFile("test-workflow/other.go", "package main")
	if len(agent.GetGeneratedFiles()) != 2 {
		t.Error("should have 2 generated files")
	}

	// 6. Run lint again
	result = agent.toolRunLint("test-workflow")
	if agent.GetLintCycles() != 2 {
		t.Errorf("lintCycles = %d, want 2", agent.GetLintCycles())
	}

	// 7. Test enforcement logic
	enforcement := agent.checkLintEnforcement([]string{"write_file"})
	if enforcement == "" {
		t.Error("should enforce lint after write")
	}

	// 8. Test completion gate
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "I'm done now!"},
		},
	}
	enforcement = agent.checkCompletionGate(resp)
	// Should enforce because lint didn't pass
	if enforcement == "" {
		t.Error("should enforce completion requirements")
	}
}
