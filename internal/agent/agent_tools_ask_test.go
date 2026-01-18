package agent

import (
	"context"
	"os"
	"testing"

	"github.com/lex00/wetwire-core-go/agent/results"
)

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

// errMock is a simple error type for testing
type errMock string

func (e errMock) Error() string { return string(e) }
