// MCP server implementation for IDE integration.
//
// When 'wetwire-github mcp' is called, this runs the MCP protocol over stdio,
// providing wetwire_init, wetwire_lint, wetwire_build, and wetwire_validate tools.
//
// This implementation uses domain.BuildMCPServer() for automatic MCP server generation.
package main

import (
	"context"

	"github.com/lex00/wetwire-core-go/domain"
	"github.com/spf13/cobra"

	githubdomain "github.com/lex00/wetwire-github-go/domain"
)

// newMCPCmd creates the "mcp" subcommand for running the MCP server.
func newMCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Run MCP server for IDE integration",
		Long: `Run the Model Context Protocol (MCP) server on stdio transport.

This command starts an MCP server that exposes wetwire-github tools:
  - wetwire_init: Initialize a new wetwire-github project
  - wetwire_lint: Lint Go packages for wetwire-github issues
  - wetwire_build: Generate GitHub Actions YAML from Go packages
  - wetwire_validate: Validate GitHub Actions YAML via actionlint
  - wetwire_list: List all discovered resources
  - wetwire_graph: Visualize resource relationships

This is typically used by IDEs and AI assistants to integrate with wetwire-github.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	}

	return cmd
}

// runMCPServer starts the MCP server on stdio transport.
func runMCPServer() error {
	server := domain.BuildMCPServer(&githubdomain.GitHubDomain{})
	return server.Start(context.Background())
}

var mcpCmd = newMCPCmd()
