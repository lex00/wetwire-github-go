package agent

import (
	"os"
	"testing"

	"github.com/lex00/wetwire-core-go/providers"
)

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

			enforcement := agent.checkCompletionGate(resp.toMessageResponse())
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
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
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
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
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
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
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
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{},
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
	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "Your CI workflow is complete! The files have been generated."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	if enforcement != "" {
		t.Errorf("checkCompletionGate() should return empty when all conditions met, got %q", enforcement)
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

	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
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

	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
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

	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "Complete!"},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should mention that lint must be called
	if !findSubstring(enforcement, "cannot complete without running the linter") {
		t.Errorf("enforcement should mention lint requirement, got %q", enforcement)
	}
}

func TestGitHubAgent_CheckCompletionGateWithFiles(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with files but no completion attempt
	agent.generatedFiles = []string{"file1.go", "file2.go"}
	agent.lintCalled = false

	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "Here are the files I generated."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should enforce lint requirement even without completion keywords
	// because files were generated
	if enforcement == "" {
		t.Error("checkCompletionGate should enforce lint check when files exist")
	}
}

func TestGitHubAgent_CheckLintEnforcementOrder(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test that order doesn't matter - if both lint and write are called,
	// no enforcement
	enforcement := agent.checkLintEnforcement([]string{"run_lint", "write_file", "read_file"})
	if enforcement != "" {
		t.Error("checkLintEnforcement should not enforce when both write and lint are present")
	}

	// Test write in the middle
	enforcement = agent.checkLintEnforcement([]string{"read_file", "write_file", "run_lint", "read_file"})
	if enforcement != "" {
		t.Error("checkLintEnforcement should not enforce when both write and lint are present (any order)")
	}
}

func TestGitHubAgent_CheckLintEnforcement_MultipleToolsWithWrite(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Test with many tools but write_file without run_lint
	tools := []string{"init_package", "write_file", "read_file", "write_file", "ask_developer"}
	enforcement := agent.checkLintEnforcement(tools)

	if enforcement == "" {
		t.Error("checkLintEnforcement should enforce when write_file is called without run_lint")
	}
}

func TestGitHubAgent_CheckCompletionGate_CaseSensitivity(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	agent.generatedFiles = []string{"file.go"}
	agent.lintCalled = false

	tests := []struct {
		name string
		text string
	}{
		{"all caps DONE", "I'M DONE WITH THE WORK"},
		{"mixed case FiNiShEd", "The task is FiNiShEd"},
		{"mixed Complete", "This is now Complete."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &providers.MessageResponse{
				Content: []providers.ContentBlock{
					{Type: "text", Text: tt.text},
				},
			}

			enforcement := agent.checkCompletionGate(resp)

			// Should detect completion keywords case-insensitively
			if enforcement == "" {
				t.Errorf("checkCompletionGate should detect completion keyword in %q", tt.text)
			}
		})
	}
}

func TestGitHubAgent_CheckCompletionGate_WithGeneratedFilesNoCompletionKeywords(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Files exist but no completion keywords in text
	agent.generatedFiles = []string{"workflow.go"}
	agent.lintCalled = false

	resp := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "I wrote the workflow file for you."},
		},
	}

	enforcement := agent.checkCompletionGate(resp)

	// Should enforce because files exist even without completion keywords
	if enforcement == "" {
		t.Error("checkCompletionGate should enforce when files exist but lint not called")
	}
	if !findSubstring(enforcement, "MUST call run_lint") {
		t.Errorf("enforcement message should require lint, got %q", enforcement)
	}
}

func TestGitHubAgent_CheckCompletionGate_AllEdgeCases(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name          string
		setupAgent    func(*GitHubAgent)
		responseText  string
		wantEnforce   bool
		wantSubstring string
	}{
		{
			name: "no files no keywords",
			setupAgent: func(a *GitHubAgent) {
				a.generatedFiles = nil
			},
			responseText: "Working on it",
			wantEnforce:  false,
		},
		{
			name: "files exist with completion - lint needed",
			setupAgent: func(a *GitHubAgent) {
				a.generatedFiles = []string{"file.go"}
				a.lintCalled = false
			},
			responseText:  "Done with everything!",
			wantEnforce:   true,
			wantSubstring: "MUST call run_lint",
		},
		{
			name: "files exist - pending lint",
			setupAgent: func(a *GitHubAgent) {
				a.generatedFiles = []string{"file.go"}
				a.lintCalled = true
				a.pendingLint = true
			},
			responseText:  "Complete!",
			wantEnforce:   true,
			wantSubstring: "written code since the last lint",
		},
		{
			name: "files exist - lint failed",
			setupAgent: func(a *GitHubAgent) {
				a.generatedFiles = []string{"file.go"}
				a.lintCalled = true
				a.pendingLint = false
				a.lintPassed = false
			},
			responseText:  "Finished!",
			wantEnforce:   true,
			wantSubstring: "linter found issues",
		},
		{
			name: "everything good",
			setupAgent: func(a *GitHubAgent) {
				a.generatedFiles = []string{"file.go"}
				a.lintCalled = true
				a.pendingLint = false
				a.lintPassed = true
			},
			responseText: "All done!",
			wantEnforce:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}

			tt.setupAgent(agent)

			resp := &providers.MessageResponse{
				Content: []providers.ContentBlock{
					{Type: "text", Text: tt.responseText},
				},
			}

			enforcement := agent.checkCompletionGate(resp)
			gotEnforce := enforcement != ""

			if gotEnforce != tt.wantEnforce {
				t.Errorf("checkCompletionGate() enforce = %v, want %v (enforcement: %q)",
					gotEnforce, tt.wantEnforce, enforcement)
			}

			if tt.wantSubstring != "" && gotEnforce {
				if !findSubstring(enforcement, tt.wantSubstring) {
					t.Errorf("checkCompletionGate() = %q, want to contain %q",
						enforcement, tt.wantSubstring)
				}
			}
		})
	}
}

func TestGitHubAgent_CheckCompletionGate_AllPaths(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name           string
		generatedFiles []string
		lintCalled     bool
		pendingLint    bool
		lintPassed     bool
		responseText   string
		wantEnforce    bool
		wantMsgPart    string
	}{
		{
			name:           "no files, no completion - passes",
			generatedFiles: nil,
			responseText:   "Working on it...",
			wantEnforce:    false,
		},
		{
			name:           "files exist, lint not called",
			generatedFiles: []string{"file.go"},
			lintCalled:     false,
			responseText:   "done",
			wantEnforce:    true,
			wantMsgPart:    "cannot complete without running the linter",
		},
		{
			name:           "files exist, pending lint",
			generatedFiles: []string{"file.go"},
			lintCalled:     true,
			pendingLint:    true,
			responseText:   "complete",
			wantEnforce:    true,
			wantMsgPart:    "written code since the last lint",
		},
		{
			name:           "files exist, lint failed",
			generatedFiles: []string{"file.go"},
			lintCalled:     true,
			pendingLint:    false,
			lintPassed:     false,
			responseText:   "finished",
			wantEnforce:    true,
			wantMsgPart:    "linter found issues",
		},
		{
			name:           "files exist, lint passed - allows completion",
			generatedFiles: []string{"file.go"},
			lintCalled:     true,
			pendingLint:    false,
			lintPassed:     true,
			responseText:   "all set",
			wantEnforce:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}

			agent.generatedFiles = tt.generatedFiles
			agent.lintCalled = tt.lintCalled
			agent.pendingLint = tt.pendingLint
			agent.lintPassed = tt.lintPassed

			resp := &providers.MessageResponse{
				Content: []providers.ContentBlock{
					{Type: "text", Text: tt.responseText},
				},
			}

			enforcement := agent.checkCompletionGate(resp)
			gotEnforce := enforcement != ""

			if gotEnforce != tt.wantEnforce {
				t.Errorf("checkCompletionGate() enforce = %v, want %v", gotEnforce, tt.wantEnforce)
			}

			if tt.wantMsgPart != "" && !findSubstring(enforcement, tt.wantMsgPart) {
				t.Errorf("enforcement = %q, want to contain %q", enforcement, tt.wantMsgPart)
			}
		})
	}
}
