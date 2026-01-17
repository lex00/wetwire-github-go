package agent

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

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

func TestGitHubAgent_ExecuteTool_LintBuildValidateRouting(t *testing.T) {
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
		checkState func(*testing.T, *GitHubAgent)
	}{
		{
			name:     "run_lint routing",
			toolName: "run_lint",
			input:    `{"path": "."}`,
			checkState: func(t *testing.T, a *GitHubAgent) {
				if !a.lintCalled {
					t.Error("lintCalled should be true after run_lint")
				}
			},
		},
		{
			name:     "run_build routing",
			toolName: "run_build",
			input:    `{"path": "."}`,
			checkState: func(t *testing.T, a *GitHubAgent) {
				// Build doesn't modify state, just verify it executed
			},
		},
		{
			name:     "run_validate routing",
			toolName: "run_validate",
			input:    `{"path": "."}`,
			checkState: func(t *testing.T, a *GitHubAgent) {
				// Validate doesn't modify state, just verify it executed
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset agent state for each test
			agent.lintCalled = false

			_ = agent.executeTool(context.Background(), tt.toolName, []byte(tt.input))

			// Verify using state checks instead of output
			tt.checkState(t, agent)
		})
	}
}

func TestGitHubAgent_ExecuteTool_AskDeveloperWithError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{err: errMockExecute("developer unavailable")}

	agent, err := NewGitHubAgent(Config{
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.executeTool(context.Background(), "ask_developer", []byte(`{"question": "Are you there?"}`))

	// Should return error message
	if !findSubstring(result, "Error:") {
		t.Errorf("executeTool(ask_developer) should return error, got %q", result)
	}
}

func TestGitHubAgent_ExecuteTool_InvalidJSON(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with completely invalid JSON
	result := agent.executeTool(context.Background(), "write_file", []byte(`{invalid json`))

	if !findSubstring(result, "Error parsing input") {
		t.Errorf("should return parsing error, got %q", result)
	}
}

func TestGitHubAgent_ExecuteTool_AllPathsCovered(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Setup test file for read_file
	testFile := filepath.Join(tmpDir, "readable.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	mock := &mockDeveloper{response: "answer"}

	agent, err := NewGitHubAgent(Config{
		WorkDir:   tmpDir,
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test all tool routes
	testCases := []struct {
		tool    string
		input   string
		checkFn func(*testing.T, string)
	}{
		{
			"init_package",
			`{"name": "test"}`,
			func(t *testing.T, result string) {
				if !findSubstring(result, "Created") {
					t.Error("init_package should report creation")
				}
			},
		},
		{
			"write_file",
			`{"path": "file.go", "content": "package main"}`,
			func(t *testing.T, result string) {
				if !findSubstring(result, "Wrote") {
					t.Error("write_file should report write")
				}
			},
		},
		{
			"read_file",
			`{"path": "readable.txt"}`,
			func(t *testing.T, result string) {
				if result != "content" {
					t.Errorf("read_file returned %q, want %q", result, "content")
				}
			},
		},
		{
			"run_lint",
			`{"path": "."}`,
			func(t *testing.T, result string) {
				// Just verify it executed
				if !agent.lintCalled {
					t.Error("run_lint should set lintCalled")
				}
			},
		},
		{
			"run_build",
			`{"path": "."}`,
			func(t *testing.T, result string) {
				// Verify non-empty result
				if result == "" {
					t.Error("run_build should return non-empty result")
				}
			},
		},
		{
			"run_validate",
			`{"path": "file.yml"}`,
			func(t *testing.T, result string) {
				// Verify non-empty result
				if result == "" {
					t.Error("run_validate should return non-empty result")
				}
			},
		},
		{
			"ask_developer",
			`{"question": "test?"}`,
			func(t *testing.T, result string) {
				if result != "answer" {
					t.Errorf("ask_developer returned %q, want %q", result, "answer")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.tool, func(t *testing.T) {
			result := agent.executeTool(context.Background(), tc.tool, []byte(tc.input))
			tc.checkFn(t, result)
		})
	}
}

func TestGitHubAgent_ExecuteTool_AllToolsWithValidInput(t *testing.T) {
	tmpDir := t.TempDir()

	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Create a test file for read operations
	testFile := filepath.Join(tmpDir, "existing.txt")
	if err := os.WriteFile(testFile, []byte("existing content"), 0644); err != nil {
		t.Fatalf("creating test file: %v", err)
	}

	mock := &mockDeveloper{response: "developer answer"}

	agent, err := NewGitHubAgent(Config{
		WorkDir:   tmpDir,
		Developer: mock,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	tests := []struct {
		name           string
		toolName       string
		input          string
		wantEmpty      bool
		skipEmptyCheck bool
		checkState     func(*testing.T, *GitHubAgent)
	}{
		{
			name:     "init_package creates project",
			toolName: "init_package",
			input:    `{"name": "new-project"}`,
			checkState: func(t *testing.T, a *GitHubAgent) {
				projectDir := filepath.Join(tmpDir, "new-project")
				if _, err := os.Stat(projectDir); os.IsNotExist(err) {
					t.Error("project directory should exist")
				}
			},
		},
		{
			name:     "write_file creates file",
			toolName: "write_file",
			input:    `{"path": "new-file.go", "content": "package main"}`,
			checkState: func(t *testing.T, a *GitHubAgent) {
				if !a.pendingLint {
					t.Error("pendingLint should be true after write")
				}
			},
		},
		{
			name:      "read_file returns content",
			toolName:  "read_file",
			input:     `{"path": "existing.txt"}`,
			wantEmpty: false,
		},
		{
			name:           "run_lint updates state",
			toolName:       "run_lint",
			input:          `{"path": "."}`,
			skipEmptyCheck: true, // lint may return empty if binary not found
			checkState: func(t *testing.T, a *GitHubAgent) {
				if !a.lintCalled {
					t.Error("lintCalled should be true")
				}
			},
		},
		{
			name:           "run_build returns output",
			toolName:       "run_build",
			input:          `{"path": "."}`,
			skipEmptyCheck: true, // build may return empty if binary not found
		},
		{
			name:           "run_validate returns output",
			toolName:       "run_validate",
			input:          `{"path": "."}`,
			skipEmptyCheck: true, // validate may return empty if binary not found
		},
		{
			name:      "ask_developer returns answer",
			toolName:  "ask_developer",
			input:     `{"question": "test question"}`,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := agent.executeTool(context.Background(), tt.toolName, []byte(tt.input))

			if !tt.skipEmptyCheck {
				if tt.wantEmpty && result != "" {
					t.Errorf("expected empty result, got %q", result)
				}
				if !tt.wantEmpty && result == "" {
					t.Error("expected non-empty result")
				}
			}

			if tt.checkState != nil {
				tt.checkState(t, agent)
			}
		})
	}
}

func TestGitHubAgent_ExecuteTool_AskDeveloperError(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	mock := &mockDeveloper{err: errMockExecute("developer unavailable")}

	agent, err := NewGitHubAgent(Config{Developer: mock})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	result := agent.executeTool(context.Background(), "ask_developer", []byte(`{"question": "test"}`))

	if !findSubstring(result, "Error") {
		t.Errorf("result should contain error, got %q", result)
	}
}

// errMockExecute is a simple error type for testing in execute tests
type errMockExecute string

func (e errMockExecute) Error() string { return string(e) }
