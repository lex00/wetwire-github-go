package kiro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureAgentConfig(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	tempHome := filepath.Join(tempDir, "home")
	tempProject := filepath.Join(tempDir, "project")
	os.MkdirAll(tempHome, 0755)
	os.MkdirAll(tempProject, 0755)

	// Set environment
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	originalDir, _ := os.Getwd()
	os.Chdir(tempProject)
	defer os.Chdir(originalDir)

	// First install
	err := EnsureInstalled()
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}

	// Verify file exists
	agentPath := filepath.Join(tempHome, ".kiro", "agents", "wetwire-github-runner.json")
	if _, err := os.Stat(agentPath); err != nil {
		t.Fatalf("agent config file not found: %v", err)
	}

	// Read and verify config can be parsed
	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("reading agent config: %v", err)
	}

	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("parsing agent config: %v", err)
	}

	if config.Name != "wetwire-github-runner" {
		t.Errorf("agent name = %q, want %q", config.Name, "wetwire-github-runner")
	}

	// Verify MCP server is configured
	if _, ok := config.MCPServers["wetwire"]; !ok {
		t.Error("wetwire MCP server not found in agent config")
	}
}

func TestEnsureProjectMCPConfig(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	tempHome := filepath.Join(tempDir, "home")
	tempProject := filepath.Join(tempDir, "project")
	os.MkdirAll(tempHome, 0755)
	os.MkdirAll(tempProject, 0755)

	// Set environment
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	originalDir, _ := os.Getwd()
	os.Chdir(tempProject)
	defer os.Chdir(originalDir)

	// First install
	err := EnsureInstalled()
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}

	// Verify file exists
	mcpPath := filepath.Join(tempProject, ".kiro", "mcp.json")
	if _, err := os.Stat(mcpPath); err != nil {
		t.Fatalf("MCP config file not found: %v", err)
	}

	// Verify config structure
	data, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("reading MCP config: %v", err)
	}

	var config map[string]MCPEntry
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("parsing MCP config: %v", err)
	}

	if config == nil {
		t.Fatal("MCP config is nil")
	}

	if _, ok := config["wetwire"]; !ok {
		t.Error("wetwire MCP server not found in config")
	}

	// Verify command is set
	if config["wetwire"].Command == "" {
		t.Error("MCP server command is empty")
	}

	// Verify args contain "mcp"
	foundMcp := false
	for _, arg := range config["wetwire"].Args {
		if arg == "mcp" {
			foundMcp = true
			break
		}
	}
	if !foundMcp {
		t.Error("MCP server args should contain 'mcp'")
	}
}

func TestMCPCommand(t *testing.T) {
	// Test that MCPCommand constant is set correctly
	if MCPCommand == "" {
		t.Error("MCPCommand constant is empty")
	}
	if MCPCommand != "wetwire-github" {
		t.Errorf("MCPCommand = %q, want %q", MCPCommand, "wetwire-github")
	}
}

func TestEnsureInstalled(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	tempHome := filepath.Join(tempDir, "home")
	tempProject := filepath.Join(tempDir, "project")
	os.MkdirAll(tempHome, 0755)
	os.MkdirAll(tempProject, 0755)

	// Set environment
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	originalDir, _ := os.Getwd()
	os.Chdir(tempProject)
	defer os.Chdir(originalDir)

	// Run EnsureInstalled
	err := EnsureInstalled()
	if err != nil {
		t.Fatalf("EnsureInstalled failed: %v", err)
	}

	// Verify agent config exists
	agentPath := filepath.Join(tempHome, ".kiro", "agents", "wetwire-github-runner.json")
	if _, err := os.Stat(agentPath); err != nil {
		t.Errorf("agent config not found after EnsureInstalled: %v", err)
	}

	// Verify MCP config exists
	mcpPath := filepath.Join(tempProject, ".kiro", "mcp.json")
	if _, err := os.Stat(mcpPath); err != nil {
		t.Errorf("MCP config not found after EnsureInstalled: %v", err)
	}
}

func TestEnsureInstalledWithForce(t *testing.T) {
	// Create temp directories
	tempDir := t.TempDir()
	tempHome := filepath.Join(tempDir, "home")
	tempProject := filepath.Join(tempDir, "project")
	os.MkdirAll(tempHome, 0755)
	os.MkdirAll(tempProject, 0755)

	// Set environment
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", originalHome)

	originalDir, _ := os.Getwd()
	os.Chdir(tempProject)
	defer os.Chdir(originalDir)

	// First install
	err := EnsureInstalled()
	if err != nil {
		t.Fatalf("first EnsureInstalled failed: %v", err)
	}

	// Force reinstall
	err = EnsureInstalledWithForce(true)
	if err != nil {
		t.Fatalf("EnsureInstalledWithForce failed: %v", err)
	}
}

func TestAgentConfigContent(t *testing.T) {
	// Verify embedded config can be parsed
	data, err := configFS.ReadFile("configs/wetwire-github-runner.json")
	if err != nil {
		t.Fatalf("reading embedded config: %v", err)
	}

	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("parsing embedded config: %v", err)
	}

	// Verify required fields
	if config.Name == "" {
		t.Error("agent name is empty")
	}
	if config.Name != "wetwire-github-runner" {
		t.Errorf("agent name = %q, want %q", config.Name, "wetwire-github-runner")
	}
	if config.Description == "" {
		t.Error("agent description is empty")
	}
	if config.Prompt == "" {
		t.Error("agent prompt is empty")
	}
	if config.Model == "" {
		t.Error("agent model is empty")
	}
}
