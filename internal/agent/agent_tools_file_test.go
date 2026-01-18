package agent

import (
	"os"
	"path/filepath"
	"testing"
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
