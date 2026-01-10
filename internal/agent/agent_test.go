package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/lex00/wetwire-core-go/agent/results"
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

func TestGitHubAgent_CheckCompletionGate(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name           string
		responseText   string
		generatedFiles []string
		lintCalled     bool
		pendingLint    bool
		lintPassed     bool
		wantEnforce    bool
		wantContains   string
	}{
		{
			name:           "no completion attempt and no files",
			responseText:   "Here is some information for you.",
			generatedFiles: nil,
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    false,
		},
		{
			name:           "done keyword without lint called",
			responseText:   "I'm done with the task.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "complete keyword without lint called",
			responseText:   "The task is complete.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "finished keyword without lint called",
			responseText:   "I have finished generating the workflow.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "thats it keyword without lint called",
			responseText:   "That's it for the implementation.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "all set keyword without lint called",
			responseText:   "You're all set now!",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "done with lint called but pending lint",
			responseText:   "Done!",
			generatedFiles: []string{"file.go"},
			lintCalled:     true,
			pendingLint:    true,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "written code since the last lint",
		},
		{
			name:           "done with lint called but not passed",
			responseText:   "Done!",
			generatedFiles: []string{"file.go"},
			lintCalled:     true,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "linter found issues",
		},
		{
			name:           "done with lint called and passed",
			responseText:   "Done!",
			generatedFiles: []string{"file.go"},
			lintCalled:     true,
			pendingLint:    false,
			lintPassed:     true,
			wantEnforce:    false,
		},
		{
			name:           "generated files but no completion keywords - enforces lint check",
			responseText:   "Here is the generated file.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "uppercase DONE keyword",
			responseText:   "I'm DONE with everything.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
		{
			name:           "mixed case Complete keyword",
			responseText:   "The workflow is Complete now.",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			pendingLint:    false,
			lintPassed:     false,
			wantEnforce:    true,
			wantContains:   "MUST call run_lint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}

			// Set agent state
			agent.generatedFiles = tt.generatedFiles
			agent.lintCalled = tt.lintCalled
			agent.pendingLint = tt.pendingLint
			agent.lintPassed = tt.lintPassed

			// Create mock response
			resp := &mockMessage{text: tt.responseText}

			enforcement := agent.checkCompletionGate(resp.toAnthropicMessage())
			gotEnforce := enforcement != ""

			if gotEnforce != tt.wantEnforce {
				t.Errorf("checkCompletionGate() enforce = %v, want %v", gotEnforce, tt.wantEnforce)
			}

			if tt.wantContains != "" && gotEnforce {
				if !contains(enforcement, tt.wantContains) {
					t.Errorf("checkCompletionGate() message = %q, want to contain %q", enforcement, tt.wantContains)
				}
			}
		})
	}
}

// mockMessage helps create anthropic.Message for testing
type mockMessage struct {
	text string
}

func (m *mockMessage) toAnthropicMessage() *anthropic.Message {
	return &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "text", Text: m.text},
		},
	}
}

// contains checks if s contains substr (case-insensitive would be needed for some checks)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGitHubAgent_ExecuteTool_UnknownTool(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.executeTool(context.Background(), "unknown_tool", []byte(`{"foo": "bar"}`))

	if result != "Unknown tool: unknown_tool" {
		t.Errorf("executeTool() = %q, want %q", result, "Unknown tool: unknown_tool")
	}
}

func TestGitHubAgent_ExecuteTool_ParseError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.executeTool(context.Background(), "write_file", []byte(`invalid json`))

	if !findSubstring(result, "Error parsing input") {
		t.Errorf("executeTool() = %q, want to contain %q", result, "Error parsing input")
	}
}

func TestGitHubAgent_ExecuteTool_Routing(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	tests := []struct {
		name       string
		toolName   string
		input      string
		wantPrefix string
	}{
		{
			name:       "init_package routing",
			toolName:   "init_package",
			input:      `{"name": "test-proj"}`,
			wantPrefix: "Created project directory",
		},
		{
			name:       "write_file routing",
			toolName:   "write_file",
			input:      `{"path": "test.go", "content": "package main"}`,
			wantPrefix: "Wrote",
		},
		{
			name:       "read_file routing - not found",
			toolName:   "read_file",
			input:      `{"path": "nonexistent.go"}`,
			wantPrefix: "Error reading file",
		},
		{
			name:       "ask_developer routing - no developer",
			toolName:   "ask_developer",
			input:      `{"question": "What is your name?"}`,
			wantPrefix: "Error:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.executeTool(context.Background(), tt.toolName, []byte(tt.input))
			if !findSubstring(result, tt.wantPrefix) {
				t.Errorf("executeTool(%s) = %q, want to contain %q", tt.toolName, result, tt.wantPrefix)
			}
		})
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

func TestGitHubAgent_CheckLintEnforcement_EdgeCases(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with nil tools called
	enforcement := agent.checkLintEnforcement(nil)
	if enforcement != "" {
		t.Errorf("checkLintEnforcement(nil) = %q, want empty", enforcement)
	}

	// Test with multiple write_file calls but one run_lint
	enforcement = agent.checkLintEnforcement([]string{"write_file", "write_file", "run_lint"})
	if enforcement != "" {
		t.Errorf("checkLintEnforcement() with writes and lint = %q, want empty", enforcement)
	}

	// Test with lint called before write - order doesn't matter, both exist so no enforcement
	enforcement = agent.checkLintEnforcement([]string{"run_lint", "write_file"})
	if enforcement != "" {
		t.Errorf("checkLintEnforcement() with both lint and write should not enforce, got %q", enforcement)
	}
}

func TestGitHubAgent_NewWithProvidedAPIKey(t *testing.T) {
	// Clear any existing API key
	origKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer func() {
		if origKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", origKey)
		}
	}()

	// Should work with provided API key even if env var is not set
	agent, err := NewGitHubAgent(Config{
		APIKey: "provided-api-key",
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent == nil {
		t.Error("NewGitHubAgent() returned nil")
	}
}

func TestGitHubAgent_CheckCompletionGate_MultipleTextBlocks(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = false

	// Create response with multiple text blocks
	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "text", Text: "First part of response. "},
			{Type: "text", Text: "Second part. "},
			{Type: "text", Text: "I'm done now."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	if enforcement == "" {
		t.Error("checkCompletionGate() should enforce when 'done' is in any text block")
	}
}

func TestGitHubAgent_CheckCompletionGate_NonTextBlocks(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// When there are generated files, checkCompletionGate enforces lint even without completion keywords
	// So we test with no generated files to verify non-text blocks are ignored
	agent.generatedFiles = nil
	agent.lintCalled = false

	// Create response with non-text blocks (tool_use)
	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "tool_use", Name: "write_file"},
			{Type: "text", Text: "Writing file now."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should not enforce since no completion keywords AND no generated files
	if enforcement != "" {
		t.Errorf("checkCompletionGate() = %q, want empty (no completion keywords and no files)", enforcement)
	}
}

func TestGitHubAgent_CheckCompletionGate_ToolUseBlocksIgnored(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = true
	agent.pendingLint = false
	agent.lintPassed = true

	// Create response with tool_use blocks that have completion-like names
	// but only text blocks should be checked for completion keywords
	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "tool_use", Name: "done_tool"},
			{Type: "text", Text: "Still working on it."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Text doesn't contain completion keywords, but files exist so it still checks
	// Since lint passed, no enforcement needed
	if enforcement != "" {
		t.Errorf("checkCompletionGate() = %q, want empty (lint passed)", enforcement)
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

func TestGitHubAgent_ExecuteTool_ReadFileRouting(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Create a file first
	testContent := "test file content"
	testPath := filepath.Join(tmpDir, "readable.txt")
	if err := os.WriteFile(testPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("writing test file: %v", err)
	}

	result := agent.executeTool(context.Background(), "read_file", []byte(`{"path": "readable.txt"}`))

	if result != testContent {
		t.Errorf("executeTool(read_file) = %q, want %q", result, testContent)
	}
}

func TestGitHubAgent_StreamHandler(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	var capturedText string
	handler := func(text string) {
		capturedText += text
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent.streamHandler == nil {
		t.Error("streamHandler should be set")
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

// mockDeveloper implements orchestrator.Developer for testing
type mockDeveloper struct {
	response string
	err      error
}

func (m *mockDeveloper) Respond(ctx context.Context, question string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
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

	mock := &mockDeveloper{err: fmt.Errorf("user cancelled")}

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

func TestGitHubAgent_ExecuteTool_AskDeveloperWithMock(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{response: "Developer response"}

	agent, err := NewGitHubAgent(Config{
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.executeTool(context.Background(), "ask_developer", []byte(`{"question": "Test question?"}`))

	if result != "Developer response" {
		t.Errorf("executeTool(ask_developer) = %q, want %q", result, "Developer response")
	}
}

func TestGitHubAgent_CheckCompletionGate_EmptyContent(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// No generated files
	agent.generatedFiles = nil

	// Empty content
	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{},
	}

	enforcement := agent.checkCompletionGate(resp)

	// No completion keywords and no files = no enforcement
	if enforcement != "" {
		t.Errorf("checkCompletionGate() with empty content = %q, want empty", enforcement)
	}
}

func TestGitHubAgent_CheckCompletionGate_AllConditionsMet(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// All conditions met for successful completion
	agent.generatedFiles = []string{"workflow.go", "triggers.go"}
	agent.lintCalled = true
	agent.pendingLint = false
	agent.lintPassed = true

	// Response says it's complete
	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "text", Text: "Your CI workflow is complete! The files have been generated."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	if enforcement != "" {
		t.Errorf("checkCompletionGate() should return empty when all conditions met, got %q", enforcement)
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

func TestGitHubAgent_CheckLintEnforcement_EnforcementMessage(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	enforcement := agent.checkLintEnforcement([]string{"write_file"})

	// Verify the enforcement message contains expected text
	if !findSubstring(enforcement, "ENFORCEMENT") {
		t.Errorf("enforcement message should contain 'ENFORCEMENT', got %q", enforcement)
	}
	if !findSubstring(enforcement, "run_lint") {
		t.Errorf("enforcement message should mention 'run_lint', got %q", enforcement)
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

func TestGitHubAgent_NewWithSession(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	session := results.NewSession("persona", "scenario")

	agent, err := NewGitHubAgent(Config{
		Session: session,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent.session != session {
		t.Error("agent.session should be set to provided session")
	}
}

func TestGitHubAgent_CheckCompletionGate_PendingLintMessage(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = true
	agent.pendingLint = true
	agent.lintPassed = false

	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "text", Text: "All done!"},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should mention pending lint
	if !findSubstring(enforcement, "written code since the last lint") {
		t.Errorf("enforcement should mention pending lint, got %q", enforcement)
	}
}

func TestGitHubAgent_CheckCompletionGate_LintNotPassedMessage(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = true
	agent.pendingLint = false
	agent.lintPassed = false

	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "text", Text: "Finished!"},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should mention lint found issues
	if !findSubstring(enforcement, "linter found issues") {
		t.Errorf("enforcement should mention lint issues, got %q", enforcement)
	}
}

func TestGitHubAgent_CheckCompletionGate_LintNotCalledMessage(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = false

	resp := &anthropic.Message{
		Content: []anthropic.ContentBlockUnion{
			{Type: "text", Text: "Complete!"},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should mention that lint must be called
	if !findSubstring(enforcement, "cannot complete without running the linter") {
		t.Errorf("enforcement should mention lint requirement, got %q", enforcement)
	}
}

func TestGitHubAgent_ExecuteTool_EmptyInput(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with empty JSON object
	result := agent.executeTool(context.Background(), "init_package", []byte(`{}`))

	// Should still work but with empty name
	if findSubstring(result, "Error") {
		t.Errorf("executeTool with empty params should not error for init_package, got %q", result)
	}
}

func TestGitHubAgent_ExecuteTool_AllTools(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Create a test file to read
	testFilePath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFilePath, []byte("test content"), 0644); err != nil {
		t.Fatalf("creating test file: %v", err)
	}

	mock := &mockDeveloper{response: "mock response"}

	agent, err := NewGitHubAgent(Config{WorkDir: tmpDir, Developer: mock})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test all tool routing
	tests := []struct {
		name     string
		tool     string
		input    string
		contains string
	}{
		{"init_package", "init_package", `{"name": "pkg"}`, "Created"},
		{"write_file", "write_file", `{"path": "new.go", "content": "pkg"}`, "Wrote"},
		{"read_file", "read_file", `{"path": "test.txt"}`, "test content"},
		{"ask_developer", "ask_developer", `{"question": "test?"}`, "mock response"},
		{"unknown", "nonexistent_tool", `{}`, "Unknown tool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.executeTool(context.Background(), tt.tool, []byte(tt.input))
			if !findSubstring(result, tt.contains) {
				t.Errorf("executeTool(%s) = %q, want to contain %q", tt.tool, result, tt.contains)
			}
		})
	}
}
