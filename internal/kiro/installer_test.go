package kiro

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureAgentConfig(t *testing.T) {
	// Create a temp home directory
	tempDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	t.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// First install
	installed, err := ensureAgentConfig(false)
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}
	if !installed {
		t.Error("expected agent config to be installed on first run")
	}

	// Verify file exists
	agentPath := filepath.Join(tempDir, ".kiro", "agents", "wetwire-runner.json")
	if _, err := os.Stat(agentPath); err != nil {
		t.Fatalf("agent config file not found: %v", err)
	}

	// Second install without force should not overwrite
	installed, err = ensureAgentConfig(false)
	if err != nil {
		t.Fatalf("second install failed: %v", err)
	}
	if installed {
		t.Error("expected agent config NOT to be reinstalled without force")
	}

	// Third install with force should overwrite
	installed, err = ensureAgentConfig(true)
	if err != nil {
		t.Fatalf("force install failed: %v", err)
	}
	if !installed {
		t.Error("expected agent config to be reinstalled with force")
	}
}

func TestEnsureProjectMCPConfig(t *testing.T) {
	// Create a temp directory and change to it
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(originalDir)

	// First install
	installed, err := ensureProjectMCPConfig(false)
	if err != nil {
		t.Fatalf("first install failed: %v", err)
	}
	if !installed {
		t.Error("expected MCP config to be installed on first run")
	}

	// Verify file exists
	mcpPath := filepath.Join(tempDir, ".kiro", "mcp.json")
	if _, err := os.Stat(mcpPath); err != nil {
		t.Fatalf("MCP config file not found: %v", err)
	}

	// Verify config structure
	data, err := os.ReadFile(mcpPath)
	if err != nil {
		t.Fatalf("reading MCP config: %v", err)
	}

	var config mcpConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("parsing MCP config: %v", err)
	}

	if config.MCPServers == nil {
		t.Fatal("MCPServers is nil")
	}

	if _, ok := config.MCPServers["wetwire"]; !ok {
		t.Error("wetwire MCP server not found in config")
	}

	// Second install without force should not overwrite
	installed, err = ensureProjectMCPConfig(false)
	if err != nil {
		t.Fatalf("second install failed: %v", err)
	}
	if installed {
		t.Error("expected MCP config NOT to be reinstalled without force")
	}
}

func TestFindWetwireBinaryPath(t *testing.T) {
	// This test just verifies the function doesn't panic
	// The actual path depends on the runtime environment
	path := findWetwireBinaryPath()
	// path may be empty (for go run fallback) or a valid path
	_ = path
}

func TestGetMCPServerConfig(t *testing.T) {
	config := getMCPServerConfig()

	// Verify config has required fields
	if config.Command == "" {
		t.Error("MCP server command is empty")
	}

	// Args should contain "mcp"
	foundMcp := false
	for _, arg := range config.Args {
		if arg == "mcp" {
			foundMcp = true
			break
		}
	}
	if !foundMcp {
		t.Error("MCP server args should contain 'mcp'")
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
	agentPath := filepath.Join(tempHome, ".kiro", "agents", "wetwire-runner.json")
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
	data, err := configFS.ReadFile("configs/wetwire-runner.json")
	if err != nil {
		t.Fatalf("reading embedded config: %v", err)
	}

	var config agentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		t.Fatalf("parsing embedded config: %v", err)
	}

	// Verify required fields
	if config.Name == "" {
		t.Error("agent name is empty")
	}
	if config.Name != "wetwire-runner" {
		t.Errorf("agent name = %q, want %q", config.Name, "wetwire-runner")
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
