package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
)

var validateFormat string

var validateCmd = &cobra.Command{
	Use:   "validate <workflow.yml>",
	Short: "Validate YAML using actionlint",
	Long: `Validate checks GitHub Actions workflow YAML files using actionlint.

Example:
  wetwire-github validate .github/workflows/ci.yml
  wetwire-github validate ci.yml --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if validateFormat == "json" {
				result := wetwire.ValidateResult{
					Success: false,
					Errors: []wetwire.ValidationError{
						{File: path, Message: "file not found"},
					},
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				_ = enc.Encode(result)
			} else {
				fmt.Fprintf(os.Stderr, "error: file not found: %s\n", path)
			}
			os.Exit(2)
		}

		// TODO: Implement actionlint integration
		result := wetwire.ValidateResult{
			Success: false,
			Errors: []wetwire.ValidationError{
				{File: path, Message: "validate command not yet implemented"},
			},
		}

		if validateFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if !result.Success {
			for _, err := range result.Errors {
				fmt.Fprintf(os.Stderr, "%s:%d:%d: %s\n", err.File, err.Line, err.Column, err.Message)
			}
			os.Exit(1)
		}

		fmt.Println("valid")
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateFormat, "format", "text", "output format (text, json)")
}
