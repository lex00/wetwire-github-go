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

func TestInstall_SetsCwd(t *testing.T) {
	// Test that cwd is set in agent config MCP server so it runs in the project directory
	// Without this, wetwire_list scans the wrong directory and returns empty results

	tmpDir := t.TempDir()
	// Resolve symlinks (macOS /var -> /private/var)
	tmpDir, _ = filepath.EvalSymlinks(tmpDir)
	projectDir := filepath.Join(tmpDir, "project")
	homeDir := filepath.Join(tmpDir, "home")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Override home directory
	t.Setenv("HOME", homeDir)

	// Override working directory for the install
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(projectDir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(origWd) }()

	// Run install
	if err := EnsureInstalledWithForce(true); err != nil {
		t.Fatalf("EnsureInstalledWithForce failed: %v", err)
	}

	// Read the agent config
	agentPath := filepath.Join(homeDir, ".kiro", "agents", AgentName+".json")
	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("failed to read agent config: %v", err)
	}

	var agent map[string]any
	if err := json.Unmarshal(data, &agent); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Get mcpServers
	mcpServers, ok := agent["mcpServers"].(map[string]any)
	if !ok {
		t.Fatal("agent config must have 'mcpServers' object")
	}

	// Get wetwire server config
	server, ok := mcpServers["wetwire"].(map[string]any)
	if !ok {
		t.Fatal("mcpServers must contain 'wetwire' object")
	}

	// Must have cwd set to the project directory
	cwd, ok := server["cwd"].(string)
	if !ok {
		t.Fatal("MCP server config must have 'cwd' field")
	}
	if cwd != projectDir {
		t.Errorf("expected cwd %q, got %q", projectDir, cwd)
	}
}

func TestInstall_HasToolsArray(t *testing.T) {
	// Test that the generated config includes tools array
	// Required for kiro to enable MCP tool usage
	// See: https://github.com/aws/amazon-q-developer-cli/issues/2640

	tmpDir := t.TempDir()
	projectDir := filepath.Join(tmpDir, "project")
	homeDir := filepath.Join(tmpDir, "home")

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Override home directory
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", origHome)

	// Override working directory for the install
	origWd, _ := os.Getwd()
	os.Chdir(projectDir)
	defer os.Chdir(origWd)

	// Run install
	if err := EnsureInstalledWithForce(true); err != nil {
		t.Fatalf("EnsureInstalledWithForce failed: %v", err)
	}

	// Read the agent config
	agentPath := filepath.Join(homeDir, ".kiro", "agents", AgentName+".json")
	data, err := os.ReadFile(agentPath)
	if err != nil {
		t.Fatalf("failed to read agent config: %v", err)
	}

	var agent map[string]any
	if err := json.Unmarshal(data, &agent); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Must have tools array
	tools, ok := agent["tools"].([]any)
	if !ok {
		t.Fatal("agent config must have 'tools' array - required for kiro MCP tool usage")
	}

	// Must have at least one tool reference
	if len(tools) == 0 {
		t.Error("tools array must not be empty")
	}

	// First tool should be @server_name format
	if len(tools) > 0 {
		tool, ok := tools[0].(string)
		if !ok || len(tool) == 0 || tool[0] != '@' {
			t.Errorf("tools should use @server_name format, got: %v", tools)
		}
	}
}
