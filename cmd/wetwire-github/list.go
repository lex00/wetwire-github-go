package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
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
		path := args[0]

		// TODO: Implement AST discovery
		result := wetwire.ListResult{
			Workflows: []wetwire.ListWorkflow{},
		}

		_ = path // Will be used when implemented

		if listFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if len(result.Workflows) == 0 {
			fmt.Println("No workflows found")
			return nil
		}

		fmt.Printf("%-20s %-30s %-6s %s\n", "WORKFLOW", "FILE", "LINE", "JOBS")
		for _, w := range result.Workflows {
			fmt.Printf("%-20s %-30s %-6d %d\n", w.Name, w.File, w.Line, w.Jobs)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().StringVar(&listFormat, "format", "text", "output format (text, json)")
}
