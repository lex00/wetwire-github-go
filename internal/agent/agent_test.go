package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitHubAgent_NoAPIKey(t *testing.T) {
	// Clear API key for test
	origKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer func() {
		if origKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", origKey)
		}
	}()

	_, err := NewGitHubAgent(Config{})
	if err == nil {
		t.Error("NewGitHubAgent() should fail without API key")
	}
}

func TestNewGitHubAgent_WithAPIKey(t *testing.T) {
	// Set a dummy API key for test
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent == nil {
		t.Error("NewGitHubAgent() returned nil")
	}
}

func TestNewGitHubAgent_Defaults(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent.maxLintCycles != 5 {
		t.Errorf("maxLintCycles = %d, want 5", agent.maxLintCycles)
	}

	if agent.workDir != "." {
		t.Errorf("workDir = %q, want %q", agent.workDir, ".")
	}
}

func TestNewGitHubAgent_CustomConfig(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{
		Model:         "custom-model",
		WorkDir:       "/tmp/test",
		MaxLintCycles: 10,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent.model != "custom-model" {
		t.Errorf("model = %q, want %q", agent.model, "custom-model")
	}

	if agent.workDir != "/tmp/test" {
		t.Errorf("workDir = %q, want %q", agent.workDir, "/tmp/test")
	}

	if agent.maxLintCycles != 10 {
		t.Errorf("maxLintCycles = %d, want 10", agent.maxLintCycles)
	}
}

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

func TestGitHubAgent_CheckLintEnforcement(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	tests := []struct {
		name        string
		toolsCalled []string
		wantEnforce bool
	}{
		{
			name:        "no tools called",
			toolsCalled: []string{},
			wantEnforce: false,
		},
		{
			name:        "only lint called",
			toolsCalled: []string{"run_lint"},
			wantEnforce: false,
		},
		{
			name:        "write and lint called",
			toolsCalled: []string{"write_file", "run_lint"},
			wantEnforce: false,
		},
		{
			name:        "only write called",
			toolsCalled: []string{"write_file"},
			wantEnforce: true,
		},
		{
			name:        "write without lint",
			toolsCalled: []string{"init_package", "write_file", "read_file"},
			wantEnforce: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enforcement := agent.checkLintEnforcement(tt.toolsCalled)
			gotEnforce := enforcement != ""
			if gotEnforce != tt.wantEnforce {
				t.Errorf("checkLintEnforcement() enforce = %v, want %v", gotEnforce, tt.wantEnforce)
			}
		})
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
		if tool.OfTool.Name != expectedTools[i] {
			t.Errorf("tools[%d].Name = %q, want %q", i, tool.OfTool.Name, expectedTools[i])
		}
	}
}
