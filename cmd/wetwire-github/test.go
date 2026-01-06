package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
)

var testFormat string
var testPersona string

var testCmd = &cobra.Command{
	Use:   "test <path>",
	Short: "Run persona-based workflow tests",
	Long: `Test runs persona-based tests against workflow declarations.

Personas simulate different GitHub Actions scenarios to validate
workflow behavior without running actual workflows.

Example:
  wetwire-github test .
  wetwire-github test ./my-workflows --format json
  wetwire-github test . --persona push`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// TODO: Implement persona-based testing
		result := wetwire.TestResult{
			Success: true,
			Tests:   []wetwire.TestCase{},
		}

		_ = path        // Will be used when implemented
		_ = testPersona // Will be used when implemented

		if testFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if len(result.Tests) == 0 {
			fmt.Println("No tests found")
			return nil
		}

		passed := 0
		failed := 0
		for _, t := range result.Tests {
			if t.Passed {
				passed++
				fmt.Printf("✓ %s\n", t.Name)
			} else {
				failed++
				fmt.Printf("✗ %s: %s\n", t.Name, t.Error)
			}
		}

		fmt.Printf("\n%d passed, %d failed\n", passed, failed)

		if !result.Success {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	testCmd.Flags().StringVar(&testFormat, "format", "text", "output format (text, json)")
	testCmd.Flags().StringVar(&testPersona, "persona", "", "run specific persona (push, pull_request, etc.)")
}
