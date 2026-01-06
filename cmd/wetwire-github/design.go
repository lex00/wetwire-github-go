package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var designStream bool
var designMaxLintCycles int
var designModel string

var designCmd = &cobra.Command{
	Use:   "design",
	Short: "AI-assisted workflow design (requires wetwire-core-go)",
	Long: `Design provides AI-assisted workflow creation using wetwire-core-go.

The design command starts an interactive session where an AI assistant
helps you create and modify GitHub Actions workflows. It uses the
wetwire-github tools to:

  - Create new workflow projects (init)
  - Write workflow code (write_file)
  - Run linting checks (lint)
  - Build YAML output (build)
  - Validate generated YAML (validate)

The assistant can automatically fix lint errors through multiple cycles.

Example:
  wetwire-github design
  wetwire-github design --stream
  wetwire-github design --max-lint-cycles 5

Note: This feature requires wetwire-core-go to be configured.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDesign()
	},
}

func init() {
	designCmd.Flags().BoolVar(&designStream, "stream", false, "stream output tokens")
	designCmd.Flags().IntVar(&designMaxLintCycles, "max-lint-cycles", 5, "maximum lint/fix cycles")
	designCmd.Flags().StringVar(&designModel, "model", "claude-sonnet-4-20250514", "model to use")
}

// runDesign executes the design command.
func runDesign() error {
	// Check for ANTHROPIC_API_KEY
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: ANTHROPIC_API_KEY environment variable not set")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "The design command requires an Anthropic API key.")
		fmt.Fprintln(os.Stderr, "Set your API key with:")
		fmt.Fprintln(os.Stderr, "  export ANTHROPIC_API_KEY=sk-ant-...")
		os.Exit(1)
		return nil
	}

	// Print status
	fmt.Println("wetwire-github design")
	fmt.Println("")
	fmt.Printf("Model: %s\n", designModel)
	fmt.Printf("Stream: %t\n", designStream)
	fmt.Printf("Max lint cycles: %d\n", designMaxLintCycles)
	fmt.Println("")
	fmt.Println("This feature requires wetwire-core-go integration.")
	fmt.Println("Implementation planned for Phase 4B.")
	fmt.Println("")
	fmt.Println("Available tools for agent:")
	fmt.Println("  - init_package: Create new workflow project")
	fmt.Println("  - write_file: Write Go code files")
	fmt.Println("  - read_file: Read existing files")
	fmt.Println("  - run_lint: Check code with linter")
	fmt.Println("  - run_build: Generate YAML from Go")
	fmt.Println("  - run_validate: Validate generated YAML")
	fmt.Println("  - ask_developer: Ask for user input")

	return nil
}
