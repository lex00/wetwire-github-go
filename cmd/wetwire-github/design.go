// Command design provides AI-assisted workflow design.
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/lex00/wetwire-github-go/internal/kiro"
	"github.com/lex00/wetwire-core-go/agent/agents"
	"github.com/lex00/wetwire-core-go/agent/orchestrator"
	"github.com/lex00/wetwire-core-go/agent/results"
	"github.com/spf13/cobra"
)

// GitHubDomain returns the domain configuration for GitHub Actions workflows.
func GitHubDomain() agents.DomainConfig {
	return agents.DomainConfig{
		Name:       "github",
		CLICommand: "wetwire-github",
		SystemPrompt: `You are a GitHub Actions workflow generator using the wetwire-github framework.
Your job is to generate Go code that defines GitHub Actions workflows.

Use the workflow pattern:
    var CIWorkflow = workflow.Workflow{
        Name: "CI",
        On: workflow.On{Push: &workflow.Push{Branches: []string{"main"}}},
        Jobs: map[string]workflow.Job{"build": BuildJob},
    }

Available tools: init_package, write_file, read_file, run_lint, run_build, ask_developer

Always run_lint after writing files, and fix any issues before running build.`,
		OutputFormat: "GitHub Actions YAML",
	}
}

// newDesignCmd creates the "design" subcommand for AI-assisted workflow design.
// It supports both Anthropic API and Kiro CLI providers for interactive code generation.
func newDesignCmd() *cobra.Command {
	var outputDir string
	var maxLintCycles int
	var stream bool
	var provider string
	var mcpServerMode bool

	cmd := &cobra.Command{
		Use:   "design [prompt]",
		Short: "AI-assisted workflow design",
		Long: `Start an interactive AI-assisted session to design and generate workflow code.

The AI agent will:
1. Ask clarifying questions about your requirements
2. Generate Go code using wetwire-github patterns
3. Run the linter and fix any issues
4. Build the GitHub Actions YAML

Providers:
    anthropic (default) - Uses Anthropic API directly (requires prompt)
    kiro                - Uses Kiro CLI with wetwire-github-runner agent

With the Kiro provider, you can omit the prompt and the agent will ask what
you'd like to create. The Anthropic provider requires an initial prompt.

Example:
    wetwire-github design "Create a CI workflow for a Go project"
    wetwire-github design --provider kiro "Create a release workflow"
    wetwire-github design --provider kiro`,
		Args: cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Hidden mode: run as MCP server (used by Kiro internally)
			if mcpServerMode {
				return runMCPServer()
			}

			prompt := strings.Join(args, " ")
			if prompt == "" && provider != "kiro" {
				return fmt.Errorf("prompt is required for the %s provider", provider)
			}
			return runDesign(prompt, outputDir, maxLintCycles, stream, provider)
		},
	}

	cmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory for generated files")
	cmd.Flags().IntVarP(&maxLintCycles, "max-lint-cycles", "l", 3, "Maximum lint/fix cycles")
	cmd.Flags().BoolVarP(&stream, "stream", "s", true, "Stream AI responses")
	cmd.Flags().StringVar(&provider, "provider", "anthropic", "AI provider: 'anthropic' or 'kiro'")
	cmd.Flags().BoolVar(&mcpServerMode, "mcp-server", false, "Run as MCP server (internal use)")
	_ = cmd.Flags().MarkHidden("mcp-server")

	return cmd
}

// runDesign starts an AI-assisted design session with the specified provider.
// It dispatches to either Kiro CLI or Anthropic API based on the provider parameter.
func runDesign(prompt, outputDir string, maxLintCycles int, stream bool, provider string) error {
	switch provider {
	case "kiro":
		return runDesignKiro(prompt, outputDir)
	case "anthropic":
		return runDesignAnthropic(prompt, outputDir, maxLintCycles, stream)
	default:
		return fmt.Errorf("unknown provider: %s (use 'anthropic' or 'kiro')", provider)
	}
}

// runDesignKiro launches an interactive Kiro CLI session with the wetwire-github-runner agent.
func runDesignKiro(prompt, outputDir string) error {
	// Change to output directory if specified
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("creating output directory: %w", err)
		}
		if err := os.Chdir(outputDir); err != nil {
			return fmt.Errorf("changing to output directory: %w", err)
		}
	}

	fmt.Println("Starting Kiro CLI design session...")
	fmt.Println()

	// Launch Kiro CLI chat (handles config installation internally)
	return kiro.LaunchChat("wetwire-github-runner", prompt)
}

// runDesignAnthropic runs an interactive design session using the Anthropic API directly.
// It creates a runner agent that generates code, runs the linter, and fixes issues.
func runDesignAnthropic(prompt, outputDir string, maxLintCycles int, stream bool) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nInterrupted, cleaning up...")
		cancel()
	}()

	// Create session for tracking
	session := results.NewSession("human", "design")

	// Create human developer (reads from stdin)
	reader := bufio.NewReader(os.Stdin)
	developer := orchestrator.NewHumanDeveloper(func() (string, error) {
		return reader.ReadString('\n')
	})

	// Create stream handler if streaming enabled
	var streamHandler agents.StreamHandler
	if stream {
		streamHandler = func(text string) {
			fmt.Print(text)
		}
	}

	// Create runner agent
	runner, err := agents.NewRunnerAgent(agents.RunnerConfig{
		WorkDir:       outputDir,
		MaxLintCycles: maxLintCycles,
		Session:       session,
		Developer:     developer,
		StreamHandler: streamHandler,
		Domain:        GitHubDomain(),
	})
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	fmt.Println("Starting AI-assisted design session...")
	fmt.Println("The AI will ask questions and generate workflow code.")
	fmt.Println("Press Ctrl+C to stop.")
	fmt.Println()

	// Run the agent
	if err := runner.Run(ctx, prompt); err != nil {
		return fmt.Errorf("design session failed: %w", err)
	}

	// Print summary
	fmt.Println("\n--- Session Summary ---")
	fmt.Printf("Generated files: %d\n", len(runner.GetGeneratedFiles()))
	for _, f := range runner.GetGeneratedFiles() {
		fmt.Printf("  - %s\n", f)
	}
	fmt.Printf("Lint cycles: %d\n", runner.GetLintCycles())
	fmt.Printf("Lint passed: %v\n", runner.LintPassed())

	return nil
}

var designCmd = newDesignCmd()
