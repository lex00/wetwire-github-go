package agent

import (
	"os"
	"testing"

	"github.com/lex00/wetwire-core-go/agent/results"
	"github.com/lex00/wetwire-core-go/providers/anthropic"
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

func TestGitHubAgent_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "no api key in config or env",
			config:  Config{},
			wantErr: true,
		},
		{
			name:    "api key in config",
			config:  Config{APIKey: "test-key"},
			wantErr: false,
		},
		{
			name: "custom work dir",
			config: Config{
				APIKey:  "test-key",
				WorkDir: "/custom/path",
			},
			wantErr: false,
		},
		{
			name: "custom model",
			config: Config{
				APIKey: "test-key",
				Model:  "custom-model-id",
			},
			wantErr: false,
		},
		{
			name: "custom max lint cycles",
			config: Config{
				APIKey:        "test-key",
				MaxLintCycles: 10,
			},
			wantErr: false,
		},
	}

	// Clear env var for consistent testing
	origKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer func() {
		if origKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", origKey)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("NewGitHubAgent() should return error")
				}
				return
			}
			if err != nil {
				t.Errorf("NewGitHubAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if agent == nil {
				t.Error("NewGitHubAgent() returned nil agent")
			}
		})
	}
}

func TestGitHubAgent_DefaultModelSetting(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Verify default model is set
	if agent.model == "" {
		t.Error("model should be set to default value")
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

func TestGitHubAgent_StreamHandlerConfiguration(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	var chunks []string
	handler := func(text string) {
		chunks = append(chunks, text)
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	// Verify handler is set
	if agent.streamHandler == nil {
		t.Error("streamHandler should be configured")
	}
}

func TestGitHubAgent_ModelConfiguration(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name      string
		model     string
		wantModel string
	}{
		{
			name:      "default model",
			model:     "",
			wantModel: anthropic.DefaultModel,
		},
		{
			name:      "custom model",
			model:     "claude-custom-model",
			wantModel: "claude-custom-model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{Model: tt.model})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}

			if agent.model != tt.wantModel {
				t.Errorf("model = %q, want %q", agent.model, tt.wantModel)
			}
		})
	}
}

func TestGitHubAgent_WorkDirConfiguration(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name    string
		workDir string
		wantDir string
	}{
		{
			name:    "default work dir",
			workDir: "",
			wantDir: ".",
		},
		{
			name:    "custom work dir",
			workDir: "/custom/path",
			wantDir: "/custom/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{WorkDir: tt.workDir})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}

			if agent.workDir != tt.wantDir {
				t.Errorf("workDir = %q, want %q", agent.workDir, tt.wantDir)
			}
		})
	}
}

func TestGitHubAgent_MaxLintCyclesDefaults(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	tests := []struct {
		name          string
		maxLintCycles int
		wantCycles    int
	}{
		{
			name:          "default max lint cycles",
			maxLintCycles: 0,
			wantCycles:    5,
		},
		{
			name:          "custom max lint cycles",
			maxLintCycles: 10,
			wantCycles:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{MaxLintCycles: tt.maxLintCycles})
			if err != nil {
				t.Fatalf("NewGitHubAgent() error = %v", err)
			}

			if agent.maxLintCycles != tt.wantCycles {
				t.Errorf("maxLintCycles = %d, want %d", agent.maxLintCycles, tt.wantCycles)
			}
		})
	}
}

func TestGitHubAgent_MaxLintCyclesConfiguration(t *testing.T) {
	os.Setenv("ANTHROPIC_API_KEY", "test-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	agent, err := NewGitHubAgent(Config{
		MaxLintCycles: 7,
	})
	if err != nil {
		t.Fatalf("NewGitHubAgent() error = %v", err)
	}

	if agent.maxLintCycles != 7 {
		t.Errorf("maxLintCycles = %d, want 7", agent.maxLintCycles)
	}
}
