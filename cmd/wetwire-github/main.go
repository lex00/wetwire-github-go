// Package main provides the wetwire-github CLI.
package main

import (
	"fmt"
	"os"

	"github.com/lex00/wetwire-github-go/domain"
)

// Version information (set via ldflags at build time)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// Set version in domain
	domain.Version = version

	// Create the domain instance and get root command with standard tools
	d := &domain.GitHubDomain{}
	root := domain.CreateRootCommand(d)

	// Add GitHub-specific commands
	root.AddCommand(designCmd)
	root.AddCommand(testCmd)
	root.AddCommand(importCmd)
	root.AddCommand(diffCmd)
	root.AddCommand(watchCmd)
	root.AddCommand(mcpCmd)

	return root.Execute()
}
