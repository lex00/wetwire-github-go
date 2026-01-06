package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/discover"
)

var listFormat string

var listCmd = &cobra.Command{
	Use:   "list <path>",
	Short: "List discovered workflows and jobs",
	Long: `List discovers and displays workflows and jobs from Go declarations.

Example:
  wetwire-github list .
  wetwire-github list ./my-workflows --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runList(args[0])
	},
}

func init() {
	listCmd.Flags().StringVar(&listFormat, "format", "text", "output format (text, json)")
}

// runList executes the list command.
func runList(path string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: resolving path: %v\n", err)
		os.Exit(1)
		return nil
	}

	// Check if path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "error: path not found: %s\n", path)
		os.Exit(1)
		return nil
	}

	// Discover workflows and jobs
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: discovery failed: %v\n", err)
		os.Exit(1)
		return nil
	}

	// Build job count map
	jobCount := make(map[string]int)
	for _, wf := range discovered.Workflows {
		jobCount[wf.Name] = len(wf.Jobs)
	}

	// Convert to result type
	result := wetwire.ListResult{
		Workflows: make([]wetwire.ListWorkflow, len(discovered.Workflows)),
	}

	for i, wf := range discovered.Workflows {
		// Make file path relative to input path for cleaner output
		relPath := wf.File
		if rel, err := filepath.Rel(absPath, filepath.Join(absPath, wf.File)); err == nil {
			relPath = rel
		}

		result.Workflows[i] = wetwire.ListWorkflow{
			Name: wf.Name,
			File: relPath,
			Line: wf.Line,
			Jobs: len(wf.Jobs),
		}
	}

	// Output result
	if listFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Text output
	if len(result.Workflows) == 0 {
		fmt.Println("No workflows found")
		return nil
	}

	fmt.Printf("%-20s %-30s %-6s %s\n", "WORKFLOW", "FILE", "LINE", "JOBS")
	for _, w := range result.Workflows {
		fmt.Printf("%-20s %-30s %-6d %d\n", w.Name, w.File, w.Line, w.Jobs)
	}

	return nil
}
