// Package main provides the wetwire-github CLI.
package main

import (
	"os"

	"github.com/spf13/cobra"
)

// Version information (set via ldflags at build time)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "wetwire-github",
	Short: "Generate GitHub YAML from typed Go declarations",
	Long: `wetwire-github generates GitHub Actions workflows, Dependabot configs,
and Issue Templates from typed Go declarations.

Example:
  wetwire-github build ./my-workflows
  wetwire-github import .github/workflows/ci.yml -o my-workflows/
  wetwire-github validate .github/workflows/ci.yml`,
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(lintCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(designCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(graphCmd)
	rootCmd.AddCommand(versionCmd)
}
