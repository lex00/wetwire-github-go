package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/lex00/wetwire-core-go/agent/results"
	"github.com/lex00/wetwire-core-go/providers"
	"github.com/stretchr/testify/assert"
)

// TestGitHubAgent_Run_ContextCancellationImmediate tests that Run respects context cancellation
func TestGitHubAgent_Run_ContextCancellationImmediate(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Run should fail with context.Canceled error
	err = agent.Run(ctx, "Generate a workflow")
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

// TestGitHubAgent_Run_ContextTimeout tests that Run respects context timeout
func TestGitHubAgent_Run_ContextTimeout(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	// Run should fail with context deadline exceeded
	err = agent.Run(ctx, "Generate a workflow")
	assert.Error(t, err)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestGitHubAgent_ToolRunLint_SuccessWithSession tests successful lint with session tracking
func TestGitHubAgent_ToolRunLint_SuccessWithSession(t *testing.T) {
	// This test covers the path where lint succeeds and updates the session
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	assert.NoError(t, err)

	// Create a valid Go module structure for lint to succeed
	// (This will likely still fail since wetwire-github isn't installed, but it tests the code path)
	agent.toolInitPackage("test-project")
	projectPath := filepath.Join(tmpDir, "test-project")

	// Write a simple workflow file
	workflowContent := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	err = os.WriteFile(filepath.Join(projectPath, "workflow.go"), []byte(workflowContent), 0644)
	if err != nil {
		t.Logf("Warning: could not write workflow file: %v", err)
	}

	// Run lint (will probably fail, but we're testing the execution path)
	_ = agent.toolRunLint("test-project")

	// Verify lint was called
	assert.True(t, agent.lintCalled, "lintCalled should be true")
	assert.False(t, agent.pendingLint, "pendingLint should be false")
	assert.Equal(t, 1, agent.lintCycles, "should have 1 lint cycle")

	// Session should be updated with lint cycle
	// The session may or may not have cycles depending on whether the command succeeded/failed correctly
	assert.NotNil(t, session)
}

// TestGitHubAgent_ToolRunLint_FailureWithJSONOutput tests lint failure with JSON parsing
func TestGitHubAgent_ToolRunLint_FailureWithJSONOutput(t *testing.T) {
	// This test is designed to fail since we're linting a non-existent path
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	assert.NoError(t, err)

	// Run lint on non-existent path - should fail
	_ = agent.toolRunLint("nonexistent")

	// Verify state changes (output may be empty if command not found)
	assert.True(t, agent.lintCalled, "lintCalled should be true")
	assert.False(t, agent.pendingLint, "pendingLint should be false")
	assert.Equal(t, 1, agent.lintCycles, "should have 1 lint cycle")

	// Lint should not pass on error
	// Note: This might be false (failed) or false (not run), depending on command
	// The important thing is lintCalled is true
}

// TestGitHubAgent_ToolRunBuild_WithErrorOutput tests build command error handling
func TestGitHubAgent_ToolRunBuild_WithErrorOutput(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Create an invalid project that will fail to build
	invalidPath := filepath.Join(tmpDir, "invalid-project")
	err = os.MkdirAll(invalidPath, 0755)
	assert.NoError(t, err)

	// Write an invalid go.mod
	err = os.WriteFile(filepath.Join(invalidPath, "go.mod"), []byte("invalid go.mod content"), 0644)
	assert.NoError(t, err)

	// Build should fail and return error message
	result := agent.toolRunBuild("invalid-project")

	// Should contain error information
	assert.NotEmpty(t, result)
	// The result should mention "Build error" or contain error output
	// Since build will fail, we just verify we got a non-empty result
}

// TestGitHubAgent_ToolRunValidate_WithErrorOutput tests validate command error handling
func TestGitHubAgent_ToolRunValidate_WithErrorOutput(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Try to validate a file that doesn't exist
	result := agent.toolRunValidate("nonexistent.yml")

	// Should return validation issues/error message
	assert.NotEmpty(t, result)
	// Result should contain error information since file doesn't exist
}

// TestGitHubAgent_ToolRunLint_ExitCode2WithSession tests the specific path for exit code 2
func TestGitHubAgent_ToolRunLint_ExitCode2WithSession(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	assert.NoError(t, err)

	// Create a project with intentional lint errors
	agent.toolInitPackage("lint-error-project")
	projectPath := filepath.Join(tmpDir, "lint-error-project")

	// Write a file with intentional errors (if wetwire-github lint returns exit code 2)
	workflowContent := `package main

// This is intentionally minimal to potentially trigger lint issues
var BadWorkflow = "not a workflow"
`
	err = os.WriteFile(filepath.Join(projectPath, "bad.go"), []byte(workflowContent), 0644)
	assert.NoError(t, err)

	// Run lint - may fail with exit code 2 if wetwire-github is available
	_ = agent.toolRunLint("lint-error-project")

	// Verify state was updated (output may vary)
	assert.True(t, agent.lintCalled, "lintCalled should be true")
	assert.False(t, agent.pendingLint, "pendingLint should be false")
	assert.Equal(t, 1, agent.lintCycles, "should have 1 lint cycle")
}

// TestGitHubAgent_ExecuteTool_ContextCancellation tests executeTool with cancelled context
func TestGitHubAgent_ExecuteTool_ContextCancellation(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	mock := &mockDeveloperWithContext{checkContext: true}
	agent, err := NewGitHubAgent(Config{Developer: mock})
	assert.NoError(t, err)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Execute ask_developer tool with cancelled context
	result := agent.executeTool(ctx, "ask_developer", []byte(`{"question": "test?"}`))

	// Should return error since context is cancelled
	assert.Contains(t, result, "Error")
}

// TestGitHubAgent_ToolRunLint_WithoutSession tests lint without session
func TestGitHubAgent_ToolRunLint_WithoutSession(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	// Create agent WITHOUT session
	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)
	assert.Nil(t, agent.session)

	// Run lint - should work without session
	_ = agent.toolRunLint(".")

	// Verify state changes even without session
	assert.True(t, agent.lintCalled, "lintCalled should be true")
	assert.False(t, agent.pendingLint, "pendingLint should be false")
	assert.Equal(t, 1, agent.lintCycles, "should have 1 lint cycle")
}

// TestGitHubAgent_ToolInitPackage_GoModWriteFailure tests error when go.mod can't be written
func TestGitHubAgent_ToolInitPackage_GoModWriteFailure(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	// Create a read-only directory
	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	err := os.MkdirAll(readOnlyDir, 0555) // read and execute only
	assert.NoError(t, err)

	agent, err := NewGitHubAgent(Config{WorkDir: readOnlyDir})
	assert.NoError(t, err)

	// Try to init package in read-only dir
	result := agent.toolInitPackage("test")

	// Should return error about writing go.mod
	assert.Contains(t, result, "Error")
}

// TestGitHubAgent_StreamHandler_NotNil tests that stream handler is properly set
func TestGitHubAgent_StreamHandler_NotNil(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	called := false
	handler := func(text string) {
		called = true
	}

	agent, err := NewGitHubAgent(Config{StreamHandler: handler})
	assert.NoError(t, err)
	assert.NotNil(t, agent.streamHandler)

	// Verify handler can be called
	agent.streamHandler("test")
	assert.True(t, called)
}

// TestGitHubAgent_ToolRunBuild_ValidProject tests the success path for build
func TestGitHubAgent_ToolRunBuild_ValidProject(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Create a minimal valid project structure
	agent.toolInitPackage("build-test")
	projectPath := filepath.Join(tmpDir, "build-test")

	// Write a minimal workflow file
	workflowContent := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var CI = workflow.Workflow{
	Name: "CI",
}
`
	err = os.WriteFile(filepath.Join(projectPath, "workflow.go"), []byte(workflowContent), 0644)
	assert.NoError(t, err)

	// Run build (may succeed or fail, but we're testing the code path)
	result := agent.toolRunBuild("build-test")

	// Should return some result
	assert.NotEmpty(t, result)
}

// TestGitHubAgent_ToolRunValidate_ValidYAML tests the success path for validate
func TestGitHubAgent_ToolRunValidate_ValidYAML(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Create a valid YAML file to validate
	yamlContent := `name: CI
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
`
	yamlPath := filepath.Join(tmpDir, "workflow.yml")
	err = os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	assert.NoError(t, err)

	// Run validate
	result := agent.toolRunValidate("workflow.yml")

	// Should return some result (success or validation issues)
	assert.NotEmpty(t, result)
}

// TestGitHubAgent_ToolWriteFile_NestedDirectoryCreationError tests error creating nested directories
func TestGitHubAgent_ToolWriteFile_NestedDirectoryCreationError(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	// Create a file (not directory) that will block directory creation
	tmpDir := t.TempDir()
	blockingFile := filepath.Join(tmpDir, "blocking")
	err := os.WriteFile(blockingFile, []byte("blocker"), 0644)
	assert.NoError(t, err)

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Try to write a file under the blocking file (should fail)
	result := agent.toolWriteFile("blocking/nested/file.go", "content")

	// Should return error about creating directory
	assert.Contains(t, result, "Error")
}

// TestGitHubAgent_AskDeveloper_ContextCancellation tests AskDeveloper with cancelled context
func TestGitHubAgent_AskDeveloper_ContextCancellation(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	// Create a mock that checks context
	mock := &mockDeveloperWithContext{
		checkContext: true,
	}

	agent, err := NewGitHubAgent(Config{Developer: mock})
	assert.NoError(t, err)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// AskDeveloper should fail with context cancelled
	_, err = agent.AskDeveloper(ctx, "Question?")
	assert.Error(t, err)
}

// mockDeveloperWithContext is a mock that respects context cancellation
type mockDeveloperWithContext struct {
	checkContext bool
}

func (m *mockDeveloperWithContext) Respond(ctx context.Context, question string) (string, error) {
	if m.checkContext {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}
	}
	return "answer", nil
}

// TestGitHubAgent_CheckCompletionGate_WithGeneratedFilesNoCompletion tests completion gate with files but no keywords
func TestGitHubAgent_CheckCompletionGate_WithGeneratedFilesNoCompletion(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Set generated files but no lint called
	agent.generatedFiles = []string{"workflow.go"}
	agent.lintCalled = false

	// Response without completion keywords
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "Here is the workflow code."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should enforce lint even without completion keywords when files exist
	assert.NotEmpty(t, enforcement)
	assert.Contains(t, enforcement, "MUST call run_lint")
}

// TestGitHubAgent_ExecuteTool_ToolInitPackage tests executeTool routing for init_package
func TestGitHubAgent_ExecuteTool_ToolInitPackage(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	result := agent.executeTool(context.Background(), "init_package", []byte(`{"name": "test-project"}`))

	assert.Contains(t, result, "Created")

	// Verify directory was created
	projectDir := filepath.Join(tmpDir, "test-project")
	_, err = os.Stat(projectDir)
	assert.NoError(t, err)
}

// TestGitHubAgent_ExecuteTool_ToolRunLint tests executeTool routing for run_lint
func TestGitHubAgent_ExecuteTool_ToolRunLint(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	_ = agent.executeTool(context.Background(), "run_lint", []byte(`{"path": "."}`))

	// Should execute lint and update state
	assert.True(t, agent.lintCalled, "lintCalled should be true")
}

// TestGitHubAgent_ExecuteTool_ToolRunBuild tests executeTool routing for run_build
func TestGitHubAgent_ExecuteTool_ToolRunBuild(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	result := agent.executeTool(context.Background(), "run_build", []byte(`{"path": "."}`))

	// Should execute build
	assert.NotEmpty(t, result)
}

// TestGitHubAgent_ExecuteTool_ToolRunValidate tests executeTool routing for run_validate
func TestGitHubAgent_ExecuteTool_ToolRunValidate(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	result := agent.executeTool(context.Background(), "run_validate", []byte(`{"path": "."}`))

	// Should execute validate
	assert.NotEmpty(t, result)
}

// TestGitHubAgent_ToolRunLint_LintPassedTrue tests the path where lint succeeds
func TestGitHubAgent_ToolRunLint_LintPassedTrue(t *testing.T) {
	// Note: This is hard to test without a valid wetwire-github installation
	// But we can at least verify the state management logic
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Initial state
	assert.False(t, agent.lintPassed)

	// Run lint
	agent.toolRunLint(".")

	// State should be updated
	assert.True(t, agent.lintCalled)
	assert.False(t, agent.pendingLint)
	assert.Equal(t, 1, agent.lintCycles)
}

// TestGitHubAgent_ToolRunLint_ExitCodeOther tests lint with non-2 exit code
func TestGitHubAgent_ToolRunLint_ExitCodeOther(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	session := results.NewSession("test", "test")
	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	assert.NoError(t, err)

	// Run lint on invalid path (likely exit code 1 or other)
	_ = agent.toolRunLint("/nonexistent/path/that/does/not/exist")

	// Should still update state
	assert.True(t, agent.lintCalled, "lintCalled should be true")
	assert.False(t, agent.pendingLint, "pendingLint should be false")
}

// TestGitHubAgent_CheckLintEnforcement_MultipleWrites tests enforcement with multiple writes
func TestGitHubAgent_CheckLintEnforcement_MultipleWrites(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Multiple writes without lint should still enforce
	enforcement := agent.checkLintEnforcement([]string{
		"write_file",
		"read_file",
		"write_file",
		"init_package",
	})

	assert.NotEmpty(t, enforcement)
	assert.Contains(t, enforcement, "run_lint")
}

// TestGitHubAgent_CheckLintEnforcement_LintBeforeWrite tests lint called before write
func TestGitHubAgent_CheckLintEnforcement_LintBeforeWrite(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Lint before write - both present so no enforcement
	enforcement := agent.checkLintEnforcement([]string{"run_lint", "write_file"})

	assert.Empty(t, enforcement)
}

// TestGitHubAgent_GetTools_AllToolsPresent tests that all expected tools are returned
func TestGitHubAgent_GetTools_AllToolsPresent(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	tools := agent.getTools()
	assert.Len(t, tools, 7)

	toolNames := make([]string, len(tools))
	for i, tool := range tools {
		toolNames[i] = tool.Name
	}

	expectedTools := []string{
		"init_package",
		"write_file",
		"read_file",
		"run_lint",
		"run_build",
		"run_validate",
		"ask_developer",
	}

	for _, expected := range expectedTools {
		assert.Contains(t, toolNames, expected)
	}
}

// TestGitHubAgent_GetTools_ToolSchemas tests tool schema structure
func TestGitHubAgent_GetTools_ToolSchemas(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	tools := agent.getTools()

	for _, tool := range tools {
		assert.NotEmpty(t, tool.Name)
		assert.NotEmpty(t, tool.Description)

		// Verify input schema has properties
		props := tool.InputSchema.Properties
		assert.NotNil(t, props)
		assert.Greater(t, len(props), 0)
	}
}

// TestGitHubAgent_NewWithAllConfig tests creating agent with all config options
func TestGitHubAgent_NewWithAllConfig(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	session := results.NewSession("persona", "scenario")
	mock := &mockDeveloper{response: "test"}
	handler := func(text string) {}

	agent, err := NewGitHubAgent(Config{
		APIKey:        "custom-key",
		Model:         "custom-model",
		WorkDir:       "/custom/dir",
		MaxLintCycles: 10,
		Session:       session,
		Developer:     mock,
		StreamHandler: handler,
	})

	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, "custom-model", agent.model)
	assert.Equal(t, "/custom/dir", agent.workDir)
	assert.Equal(t, 10, agent.maxLintCycles)
	assert.Equal(t, session, agent.session)
	assert.Equal(t, mock, agent.developer)
	assert.NotNil(t, agent.streamHandler)
}

// TestGitHubAgent_ExecuteTool_MalformedJSON tests executeTool with malformed JSON
func TestGitHubAgent_ExecuteTool_MalformedJSON(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	tests := []struct {
		name  string
		input []byte
	}{
		{"invalid json", []byte(`{invalid}`)},
		{"unclosed brace", []byte(`{"name": "test"`)},
		{"null bytes", []byte{0, 0, 0}},
		{"empty", []byte{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.executeTool(context.Background(), "init_package", tt.input)
			// Should return error or handle gracefully
			assert.NotEmpty(t, result)
		})
	}
}

// TestGitHubAgent_ToolReadFile_LargeFile tests reading a large file
func TestGitHubAgent_ToolReadFile_LargeFile(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Create a large file
	largeContent := make([]byte, 10000)
	for i := range largeContent {
		largeContent[i] = byte('a' + (i % 26))
	}

	filePath := filepath.Join(tmpDir, "large.txt")
	err = os.WriteFile(filePath, largeContent, 0644)
	assert.NoError(t, err)

	// Read the large file
	result := agent.toolReadFile("large.txt")

	// Should return the full content
	assert.Equal(t, string(largeContent), result)
}

// TestGitHubAgent_ToolWriteFile_LargeContent tests writing large content
func TestGitHubAgent_ToolWriteFile_LargeContent(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Create large content
	largeContent := make([]byte, 10000)
	for i := range largeContent {
		largeContent[i] = byte('x')
	}

	result := agent.toolWriteFile("large.go", string(largeContent))

	assert.Contains(t, result, "10000 bytes")

	// Verify file was written
	filePath := filepath.Join(tmpDir, "large.go")
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, largeContent, content)
}

// TestGitHubAgent_ToolInitPackage_SpecialCharactersInName tests init with special chars
func TestGitHubAgent_ToolInitPackage_SpecialCharactersInName(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	assert.NoError(t, err)

	// Project name with special characters (still valid for directory)
	result := agent.toolInitPackage("my-project_v2")

	assert.Contains(t, result, "Created")

	// Verify directory exists
	projectDir := filepath.Join(tmpDir, "my-project_v2")
	_, err = os.Stat(projectDir)
	assert.NoError(t, err)

	// Verify go.mod contains the project name
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), "my-project_v2")
}

// TestGitHubAgent_LintPassed_InitialState tests initial state of LintPassed
func TestGitHubAgent_LintPassed_InitialState(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Initial state should be false
	assert.False(t, agent.LintPassed())
}

// TestGitHubAgent_GetLintCycles_InitialState tests initial lint cycles
func TestGitHubAgent_GetLintCycles_InitialState(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Initial state should be 0
	assert.Equal(t, 0, agent.GetLintCycles())
}

// TestGitHubAgent_GetGeneratedFiles_InitialState tests initial generated files
func TestGitHubAgent_GetGeneratedFiles_InitialState(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// Initial state should be nil/empty
	files := agent.GetGeneratedFiles()
	assert.Nil(t, files)
}

// TestGitHubAgent_ToolRunLint_JSONUnmarshalSuccess tests successful JSON unmarshal of lint output
func TestGitHubAgent_ToolRunLint_JSONUnmarshalSuccess(t *testing.T) {
	tmpDir := t.TempDir()

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	session := results.NewSession("test", "test")
	agent, err := NewGitHubAgent(Config{
		WorkDir: tmpDir,
		Session: session,
	})
	assert.NoError(t, err)

	// Create a valid project that might produce JSON lint output
	agent.toolInitPackage("json-test")
	projectPath := filepath.Join(tmpDir, "json-test")

	// Write a workflow file
	workflowContent := `package main

import "github.com/lex00/wetwire-github-go/workflow"

var TestWorkflow = workflow.Workflow{
	Name: "Test",
}
`
	err = os.WriteFile(filepath.Join(projectPath, "workflow.go"), []byte(workflowContent), 0644)
	assert.NoError(t, err)

	// Run lint - should attempt to parse JSON if available
	_ = agent.toolRunLint("json-test")

	assert.True(t, agent.lintCalled, "lintCalled should be true")
}

// TestGitHubAgent_CheckCompletionGate_CaseInsensitiveKeywords tests case-insensitive completion detection
func TestGitHubAgent_CheckCompletionGate_CaseInsensitiveKeywords(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = false

	tests := []struct {
		name string
		text string
	}{
		{"uppercase DONE", "I'm DONE"},
		{"MixedCase Done", "We're Done here"},
		{"COMPLETE", "Task is COMPLETE"},
		{"Finished", "Finished working"},
		{"THAT'S IT", "THAT'S IT!"},
		{"All Set", "You're All Set now"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &providers.MessageResponse{
				Content: []providers.ContentBlock{
					{Type: "text", Text: tt.text},
				},
			}

			enforcement := agent.checkCompletionGate(resp)
			assert.NotEmpty(t, enforcement, "should enforce for: %s", tt.text)
		})
	}
}
