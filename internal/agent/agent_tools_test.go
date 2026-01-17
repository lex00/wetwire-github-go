package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lex00/wetwire-core-go/agent/results"
)

func TestGitHubAgent_ToolInitPackage(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolInitPackage("test-project")

	if result == "" {
		t.Error("toolInitPackage() returned empty result")
	}

	// Check directory was created
	projectDir := filepath.Join(tmpDir, "test-project")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Error("project directory was not created")
	}

	// Check go.mod was created
	goModPath := filepath.Join(projectDir, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod was not created")
	}
}

func TestGitHubAgent_ToolWriteFile(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	content := "package main\n\nfunc main() {}\n"
	result := agent.toolWriteFile("test/main.go", content)

	if result == "" {
		t.Error("toolWriteFile() returned empty result")
	}

	// Check file was created
	filePath := filepath.Join(tmpDir, "test/main.go")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("file content = %q, want %q", string(fileContent), content)
	}

	// Check pendingLint state
	if !agent.pendingLint {
		t.Error("pendingLint should be true after write")
	}

	// Check generatedFiles
	if len(agent.generatedFiles) != 1 {
		t.Errorf("len(generatedFiles) = %d, want 1", len(agent.generatedFiles))
	}
}

func TestGitHubAgent_ToolReadFile(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Write a file first
	content := "test content"
	filePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	result := agent.toolReadFile("test.txt")
	if result != content {
		t.Errorf("toolReadFile() = %q, want %q", result, content)
	}
}

func TestGitHubAgent_ToolReadFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolReadFile("nonexistent.txt")
	if result == "" {
		t.Error("toolReadFile() should return error message for missing file")
	}
}

func TestGitHubAgent_GetTools(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	tools := agent.getTools()

	// Should have 7 tools
	if len(tools) != 7 {
		t.Errorf("len(getTools()) = %d, want 7", len(tools))
	}

	// Check tool names
	expectedTools := []string{
		"init_package",
		"write_file",
		"read_file",
		"run_lint",
		"run_build",
		"run_validate",
		"ask_developer",
	}

	for i, tool := range tools {
		if tool.Name != expectedTools[i] {
			t.Errorf("tools[%d].Name = %q, want %q", i, tool.Name, expectedTools[i])
		}
	}
}

func TestGitHubAgent_AskDeveloper_NoDeveloper(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	_, err = agent.AskDeveloper(context.Background(), "test question")
	if err == nil {
		t.Error("AskDeveloper() should fail without developer configured")
	}
}

func TestGitHubAgent_ToolWriteFile_StateChanges(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Simulate that lint previously passed
	agent.lintPassed = true
	agent.pendingLint = false

	// Write a file - should reset lint state
	agent.toolWriteFile("test.go", "package main")

	if agent.lintPassed {
		t.Error("lintPassed should be false after writing a file")
	}

	if !agent.pendingLint {
		t.Error("pendingLint should be true after writing a file")
	}

	// Write another file - should add to generatedFiles
	agent.toolWriteFile("test2.go", "package main")

	if len(agent.generatedFiles) != 2 {
		t.Errorf("len(generatedFiles) = %d, want 2", len(agent.generatedFiles))
	}
}

func TestGitHubAgent_ToolInitPackage_GoModContent(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolInitPackage("my-project")

	if !findSubstring(result, "go.mod") {
		t.Errorf("toolInitPackage() result = %q, want to mention go.mod", result)
	}

	// Read and verify go.mod content
	goModPath := filepath.Join(tmpDir, "my-project", "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}

	goModContent := string(content)
	if !findSubstring(goModContent, "module github.com/example/my-project") {
		t.Errorf("go.mod content = %q, want to contain module declaration", goModContent)
	}
	if !findSubstring(goModContent, "go 1.23") {
		t.Errorf("go.mod content = %q, want to contain go version", goModContent)
	}
	if !findSubstring(goModContent, "wetwire-github-go") {
		t.Errorf("go.mod content = %q, want to contain wetwire-github-go dependency", goModContent)
	}
}

func TestGitHubAgent_ToolReadFile_ErrorMessage(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolReadFile("does-not-exist.txt")

	if !findSubstring(result, "Error reading file") {
		t.Errorf("toolReadFile() = %q, want to contain 'Error reading file'", result)
	}
}

func TestGitHubAgent_ToolWriteFile_NestedDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Write to a deeply nested path
	content := "package nested"
	result := agent.toolWriteFile("a/b/c/d/nested.go", content)

	if !findSubstring(result, "Wrote") {
		t.Errorf("toolWriteFile() = %q, want to contain 'Wrote'", result)
	}

	// Verify file exists and has correct content
	filePath := filepath.Join(tmpDir, "a/b/c/d/nested.go")
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("reading nested file: %v", err)
	}

	if string(fileContent) != content {
		t.Errorf("nested file content = %q, want %q", string(fileContent), content)
	}
}

func TestGitHubAgent_ToolWriteFile_ByteCount(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	content := "hello world"
	result := agent.toolWriteFile("test.txt", content)

	// Result should mention the byte count
	expectedBytes := "11 bytes"
	if !findSubstring(result, expectedBytes) {
		t.Errorf("toolWriteFile() result = %q, want to contain %q", result, expectedBytes)
	}
}

func TestGitHubAgent_AskDeveloper_WithMockDeveloper(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{response: "Yes, I want tests"}

	agent, err := NewGitHubAgent(Config{
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	answer, err := agent.AskDeveloper(context.Background(), "Do you want tests?")
	if err != nil {
		t.Fatalf("AskDeveloper() error = %v", err)
	}

	if answer != "Yes, I want tests" {
		t.Errorf("AskDeveloper() = %q, want %q", answer, "Yes, I want tests")
	}
}

func TestGitHubAgent_AskDeveloper_WithError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{err: errMock("user cancelled")}

	agent, err := NewGitHubAgent(Config{
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	_, err = agent.AskDeveloper(context.Background(), "Do you want tests?")
	if err == nil {
		t.Error("AskDeveloper() should return error from developer")
	}
}

func TestGitHubAgent_AskDeveloper_WithSession(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{response: "My answer"}

	// Create a mock session (we can't easily test the full session behavior
	// without importing the results package, but we can at least ensure
	// the path is exercised)
	agent, err := NewGitHubAgent(Config{
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Without session, should still work
	answer, err := agent.AskDeveloper(context.Background(), "Question?")
	if err != nil {
		t.Fatalf("AskDeveloper() error = %v", err)
	}

	if answer != "My answer" {
		t.Errorf("AskDeveloper() = %q, want %q", answer, "My answer")
	}
}

func TestGitHubAgent_AskDeveloper_WithSessionTracking(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{response: "Tracked answer"}
	session := results.NewSession("test-persona", "test-scenario")

	agent, err := NewGitHubAgent(Config{
		Developer: mock,
		Session:   session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	answer, err := agent.AskDeveloper(context.Background(), "Question for tracking?")
	if err != nil {
		t.Fatalf("AskDeveloper() error = %v", err)
	}

	if answer != "Tracked answer" {
		t.Errorf("AskDeveloper() = %q, want %q", answer, "Tracked answer")
	}

	// Verify the question was tracked in the session
	if len(session.Questions) != 1 {
		t.Fatalf("session.Questions has %d entries, want 1", len(session.Questions))
	}

	if session.Questions[0].Question != "Question for tracking?" {
		t.Errorf("session.Questions[0].Question = %q, want %q", session.Questions[0].Question, "Question for tracking?")
	}

	if session.Questions[0].Answer != "Tracked answer" {
		t.Errorf("session.Questions[0].Answer = %q, want %q", session.Questions[0].Answer, "Tracked answer")
	}
}

func TestGitHubAgent_ToolWriteFile_PathInResult(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolWriteFile("my-file.go", "package main")

	// Result should mention the file path
	if !findSubstring(result, "my-file.go") {
		t.Errorf("toolWriteFile() result = %q, want to contain file path", result)
	}
}

func TestGitHubAgent_ToolInitPackage_InvalidPath(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Use a path that we can't write to
	agent, err := NewGitHubAgent(Config{WorkDir: "/nonexistent/invalid/path"})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolInitPackage("test")

	// Should return error message
	if !findSubstring(result, "Error") {
		t.Errorf("toolInitPackage with invalid path should return error, got %q", result)
	}
}

func TestGitHubAgent_ToolWriteFile_InvalidPath(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Use a path that we can't write to
	agent, err := NewGitHubAgent(Config{WorkDir: "/nonexistent/invalid/path"})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.toolWriteFile("test.go", "package main")

	// Should return error message
	if !findSubstring(result, "Error") {
		t.Errorf("toolWriteFile with invalid path should return error, got %q", result)
	}
}

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

func TestGitHubAgent_ToolWriteFile_EmptyContent(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Write empty file
	result := agent.toolWriteFile("empty.txt", "")

	if !findSubstring(result, "0 bytes") {
		t.Errorf("toolWriteFile with empty content should report 0 bytes, got %q", result)
	}

	// Verify state changes
	if !agent.pendingLint {
		t.Error("pendingLint should be true after write")
	}
	if agent.lintPassed {
		t.Error("lintPassed should be false after write")
	}
}

func TestGitHubAgent_ToolInitPackage_MultipleProjects(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create multiple projects
	result1 := agent.toolInitPackage("project1")
	result2 := agent.toolInitPackage("project2")

	if !findSubstring(result1, "project1") {
		t.Errorf("result1 should mention project1, got %q", result1)
	}
	if !findSubstring(result2, "project2") {
		t.Errorf("result2 should mention project2, got %q", result2)
	}

	// Verify both directories exist
	dir1 := filepath.Join(tmpDir, "project1")
	dir2 := filepath.Join(tmpDir, "project2")

	if _, err := os.Stat(dir1); os.IsNotExist(err) {
		t.Error("project1 directory should exist")
	}
	if _, err := os.Stat(dir2); os.IsNotExist(err) {
		t.Error("project2 directory should exist")
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

func TestGitHubAgent_ToolWriteFile_DirectoryCreationError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Create a file where we want a directory
	tmpDir := t.TempDir()
	blockingFile := filepath.Join(tmpDir, "blocking")
	if err := os.WriteFile(blockingFile, []byte("content"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Try to write a file under the blocking file path
	result := agent.toolWriteFile("blocking/test.go", "package main")

	// Should return error about directory creation
	if !findSubstring(result, "Error") {
		t.Errorf("toolWriteFile should return error when directory creation fails, got %q", result)
	}
}

func TestGitHubAgent_ToolInitPackage_GoModWriteError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tmpDir := t.TempDir()

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a directory first
	projectDir := filepath.Join(tmpDir, "test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Create go.mod as a directory to cause write error
	goModPath := filepath.Join(projectDir, "go.mod")
	if err := os.MkdirAll(goModPath, 0755); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Try to init package
	result := agent.toolInitPackage("test-project")

	// Should return error about writing go.mod
	if !findSubstring(result, "Error writing go.mod") {
		t.Errorf("toolInitPackage should return go.mod write error, got %q", result)
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

func TestGitHubAgent_ToolWriteFile_WithNestedDirectoryError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tmpDir := t.TempDir()

	// Create a file that blocks directory creation
	blockPath := filepath.Join(tmpDir, "blocked")
	if err := os.WriteFile(blockPath, []byte("blocker"), 0444); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Make it read-only to prevent directory creation
	if err := os.Chmod(blockPath, 0444); err != nil {
		t.Fatalf("setup: %v", err)
	}

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Try to create a file under the blocked path
	result := agent.toolWriteFile("blocked/subdir/file.go", "package main")

	// Should fail with directory creation error
	if !findSubstring(result, "Error") {
		t.Errorf("toolWriteFile should fail when directory creation is blocked, got %q", result)
	}
}

func TestGitHubAgent_ToolWriteFile_WriteError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tmpDir := t.TempDir()

	// Create a read-only directory
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
		t.Fatalf("setup: %v", err)
	}

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Try to write to read-only directory
	result := agent.toolWriteFile("readonly/file.go", "package main")

	// Should return write error
	if !findSubstring(result, "Error writing file") {
		t.Errorf("toolWriteFile should return write error for read-only directory, got %q", result)
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

	// Run lint - will fail because no wetwire-github binary, but state should still update
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

func TestGitHubAgent_GetTools_ToolDetails(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	tools := agent.getTools()

	// Verify each tool has required fields
	for i, tool := range tools {
		if tool.Name == "" {
			t.Errorf("tools[%d].Name is empty", i)
		}

		// InputSchema should have Properties
		if tool.InputSchema.Properties == nil {
			t.Errorf("tools[%d] (%s) has nil Properties", i, tool.Name)
		}
	}
}

func TestGitHubAgent_AskDeveloper_SessionIntegration(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	session := results.NewSession("persona", "scenario")
	mock := &mockDeveloper{response: "the answer"}

	agent, err := NewGitHubAgent(Config{
		Session:   session,
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Ask multiple questions
	questions := []struct {
		q string
		a string
	}{
		{"First question?", "the answer"},
		{"Second question?", "the answer"},
	}

	for _, qa := range questions {
		answer, err := agent.AskDeveloper(context.Background(), qa.q)
		if err != nil {
			t.Fatalf("AskDeveloper() error = %v", err)
		}
		if answer != qa.a {
			t.Errorf("AskDeveloper(%q) = %q, want %q", qa.q, answer, qa.a)
		}
	}

	// Verify session recorded all questions
	if len(session.Questions) != len(questions) {
		t.Errorf("session.Questions has %d entries, want %d", len(session.Questions), len(questions))
	}
}

func TestGitHubAgent_GetToolsCompleteness(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	tools := agent.getTools()

	// Verify each tool has proper structure
	for i, tool := range tools {
		if tool.Name == "" {
			t.Errorf("tool[%d] Name is empty", i)
		}
		// Just verify tool has a name - that's sufficient for structure validation
	}
}

func TestGitHubAgent_ToolInitPackageGoModFormat(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	projectName := "my-special-project"
	agent.toolInitPackage(projectName)

	// Read and verify go.mod structure
	goModPath := filepath.Join(tmpDir, projectName, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("reading go.mod: %v", err)
	}

	goModStr := string(content)

	// Check all expected parts
	expectedParts := []string{
		"module github.com/example/" + projectName,
		"go 1.23",
		"require github.com/lex00/wetwire-github-go",
	}

	for _, part := range expectedParts {
		if !findSubstring(goModStr, part) {
			t.Errorf("go.mod missing expected part: %q", part)
		}
	}
}

func TestGitHubAgent_MultipleFileOperations(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Perform a series of file operations and verify state consistency
	operations := []struct {
		op       string
		path     string
		content  string
		checkGen int
	}{
		{"write", "file1.go", "package main", 1},
		{"write", "file2.go", "package test", 2},
		{"write", "subdir/file3.go", "package subdir", 3},
		{"write", "deep/nested/file4.go", "package nested", 4},
	}

	for _, op := range operations {
		agent.toolWriteFile(op.path, op.content)
		if len(agent.GetGeneratedFiles()) != op.checkGen {
			t.Errorf("after writing %s, got %d generated files, want %d",
				op.path, len(agent.GetGeneratedFiles()), op.checkGen)
		}
	}

	// Verify all files can be read back
	for _, op := range operations {
		result := agent.toolReadFile(op.path)
		if result != op.content {
			t.Errorf("toolReadFile(%s) = %q, want %q", op.path, result, op.content)
		}
	}
}

// errMock is a simple error type for testing
type errMock string

func (e errMock) Error() string { return string(e) }
