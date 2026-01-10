package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lex00/wetwire-github-go/internal/agent"
)

var designStream bool
var designMaxLintCycles int
var designModel string
var designWorkDir string
var designMCPServer bool

var designCmd = &cobra.Command{
	Use:   "design [prompt]",
	Short: "AI-assisted workflow design",
	Long: `Design provides AI-assisted workflow creation.

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
  wetwire-github design "Create a CI workflow for a Go project"
  wetwire-github design --stream "Add a release workflow"
  wetwire-github design --max-lint-cycles 5 "Create multi-platform build"

Requires ANTHROPIC_API_KEY environment variable to be set.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Run MCP server if requested
		if designMCPServer {
			return runMCPServer()
		}
		prompt := strings.Join(args, " ")
		return runDesign(prompt)
	},
}

func init() {
	designCmd.Flags().BoolVar(&designStream, "stream", false, "stream output tokens")
	designCmd.Flags().IntVar(&designMaxLintCycles, "max-lint-cycles", 5, "maximum lint/fix cycles")
	designCmd.Flags().StringVar(&designModel, "model", "claude-sonnet-4-20250514", "model to use")
	designCmd.Flags().StringVarP(&designWorkDir, "workdir", "w", ".", "working directory for generated files")
	designCmd.Flags().BoolVar(&designMCPServer, "mcp-server", false, "Run as MCP server (internal use)")
	_ = designCmd.Flags().MarkHidden("mcp-server")
}

// consoleDeveloper implements orchestrator.Developer for console input.
type consoleDeveloper struct {
	reader *bufio.Reader
}

func newConsoleDeveloper() *consoleDeveloper {
	return &consoleDeveloper{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (d *consoleDeveloper) Respond(ctx context.Context, question string) (string, error) {
	fmt.Printf("\n[Agent Question] %s\n> ", question)
	answer, err := d.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(answer), nil
}

// runDesign executes the design command.
func runDesign(prompt string) error {
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

	// If no prompt provided, prompt for one
	if prompt == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Describe the workflow you want to create:\n> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading input: %w", err)
		}
		prompt = strings.TrimSpace(input)
		if prompt == "" {
			return fmt.Errorf("no prompt provided")
		}
	}

	// Print status
	fmt.Println("wetwire-github design")
	fmt.Println("")
	fmt.Printf("Model: %s\n", designModel)
	fmt.Printf("Stream: %t\n", designStream)
	fmt.Printf("Max lint cycles: %d\n", designMaxLintCycles)
	fmt.Printf("Work directory: %s\n", designWorkDir)
	fmt.Println("")

	// Create agent config
	config := agent.Config{
		APIKey:        apiKey,
		Model:         designModel,
		WorkDir:       designWorkDir,
		MaxLintCycles: designMaxLintCycles,
		Developer:     newConsoleDeveloper(),
	}

	// Add stream handler if streaming is enabled
	if designStream {
		config.StreamHandler = func(text string) {
			fmt.Print(text)
		}
	}

	// Create and run the agent
	a, err := agent.NewGitHubAgent(config)
	if err != nil {
		return fmt.Errorf("creating agent: %w", err)
	}

	ctx := context.Background()
	if err := a.Run(ctx, prompt); err != nil {
		return fmt.Errorf("agent failed: %w", err)
	}

	// Print summary
	fmt.Println("")
	fmt.Println("---")
	fmt.Printf("Generated files: %d\n", len(a.GetGeneratedFiles()))
	for _, f := range a.GetGeneratedFiles() {
		fmt.Printf("  %s\n", f)
	}
	fmt.Printf("Lint cycles: %d\n", a.GetLintCycles())
	fmt.Printf("Lint passed: %t\n", a.LintPassed())

	return nil
}
