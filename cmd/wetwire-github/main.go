// Package main provides the wetwire-github CLI.
package main

import (
	"fmt"
	"os"

	corecmd "github.com/lex00/wetwire-core-go/cmd"
)

// Version information (set via ldflags at build time)
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// rootCmd uses the core command framework for consistent CLI structure.
var rootCmd = corecmd.NewRootCommand(
	"wetwire-github",
	"Generate GitHub YAML from typed Go declarations",
)

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
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(watchCmd)
}
