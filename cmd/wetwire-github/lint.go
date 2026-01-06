package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
)

var lintFormat string
var lintFix bool

var lintCmd = &cobra.Command{
	Use:   "lint <path>",
	Short: "Check Go code for wetwire best practices",
	Long: `Lint checks Go workflow declarations for wetwire best practices.

Rules:
  WAG001: Use typed action wrappers instead of raw uses: strings
  WAG002: Use condition builders instead of raw expression strings
  WAG003: Use secrets context instead of hardcoded strings
  WAG004: Use matrix builder instead of inline maps
  WAG005: Extract inline structs to named variables
  WAG006: Detect duplicate workflow names
  WAG007: Flag oversized files (>N jobs)
  WAG008: Avoid hardcoded expression strings

Example:
  wetwire-github lint .
  wetwire-github lint ./my-workflows --format json
  wetwire-github lint . --fix`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// TODO: Implement linter framework
		result := wetwire.LintResult{
			Success: true,
			Issues:  []wetwire.LintIssue{},
		}

		_ = path   // Will be used when implemented
		_ = lintFix // Will be used when implemented

		if lintFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if len(result.Issues) == 0 {
			fmt.Println("No issues found")
			return nil
		}

		for _, issue := range result.Issues {
			fmt.Printf("%s:%d:%d: %s [%s]\n",
				issue.File, issue.Line, issue.Column, issue.Message, issue.Rule)
		}

		if !result.Success {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	lintCmd.Flags().StringVar(&lintFormat, "format", "text", "output format (text, json)")
	lintCmd.Flags().BoolVar(&lintFix, "fix", false, "automatically fix issues where possible")
}
