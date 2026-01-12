package kiro

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed configs/wetwire-github-runner.json
var configFS embed.FS

// AgentConfig represents a Kiro agent configuration.
type AgentConfig struct {
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Prompt      string               `json:"prompt"`
	Model       string               `json:"model"`
	MCPServers  map[string]MCPServer `json:"mcpServers"`
	Tools       []string             `json:"tools"`
}

// MCPServer represents an MCP server in the embedded config.
type MCPServer struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

// MCPEntry represents an MCP server configuration.
type MCPEntry struct {
	Name    string   `json:"name"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Cwd     string   `json:"cwd,omitempty"`
}

// EnsureInstalled installs the Kiro agent configuration if not already present.
func EnsureInstalled() error {
	return EnsureInstalledWithForce(false)
}

// EnsureInstalledWithForce installs the Kiro agent configuration.
// If force is true, overwrites any existing configuration.
func EnsureInstalledWithForce(force bool) error {
	// Get user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("get home directory: %w", err)
	}

	// Create agent config directory
	agentDir := filepath.Join(homeDir, ".kiro", "agents")
	if err := os.MkdirAll(agentDir, 0755); err != nil {
		return fmt.Errorf("create agent directory: %w", err)
	}

	// Install wetwire-github-runner agent config
	agentPath := filepath.Join(agentDir, "wetwire-github-runner.json")
	if force || !fileExists(agentPath) {
		if err := installAgentConfig(agentPath); err != nil {
			return fmt.Errorf("install agent config: %w", err)
		}
	}

	// Install project-level MCP config
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	mcpDir := filepath.Join(cwd, ".kiro")
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return fmt.Errorf("create .kiro directory: %w", err)
	}

	mcpPath := filepath.Join(mcpDir, "mcp.json")
	if force || !fileExists(mcpPath) {
		if err := installMCPConfig(mcpPath); err != nil {
			return fmt.Errorf("install mcp config: %w", err)
		}
	}

	return nil
}

// installAgentConfig installs the wetwire-github-runner agent configuration.
func installAgentConfig(path string) error {
	// Read embedded config
	data, err := configFS.ReadFile("configs/wetwire-github-runner.json")
	if err != nil {
		return fmt.Errorf("reading embedded config: %w", err)
	}

	// Parse and update with full MCP binary path
	var config AgentConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("parsing embedded config: %w", err)
	}

	// Update MCP server with correct command (now using mcp subcommand)
	if config.MCPServers == nil {
		config.MCPServers = make(map[string]MCPServer)
	}
	config.MCPServers["wetwire"] = MCPServer{
		Command: MCPCommand,
		Args:    []string{"mcp"},
	}

	// Write config
	output, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// installMCPConfig installs the project-level MCP configuration.
func installMCPConfig(path string) error {
	cwd, _ := os.Getwd()

	// Determine command (prefer binary, fallback to go run)
	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "go"
	}

	var mcpConfig map[string]MCPEntry
	if binaryPath == "go" {
		mcpConfig = map[string]MCPEntry{
			"wetwire": {
				Command: "go",
				Args:    []string{"run", "./cmd/wetwire-github", "mcp"},
				Cwd:     cwd,
			},
		}
	} else {
		mcpConfig = map[string]MCPEntry{
			"wetwire": {
				Command: MCPCommand,
				Args:    []string{"mcp"},
				Cwd:     cwd,
			},
		}
	}

	output, err := json.MarshalIndent(mcpConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal mcp config: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("write mcp config: %w", err)
	}

	return nil
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
