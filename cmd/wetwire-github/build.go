package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
)

var buildOutput string
var buildFormat string
var buildType string

var buildCmd = &cobra.Command{
	Use:   "build <path>",
	Short: "Generate YAML from Go workflow declarations",
	Long: `Build reads Go workflow declarations and generates GitHub YAML files.

Example:
  wetwire-github build .
  wetwire-github build ./my-workflows -o .github/workflows/
  wetwire-github build . --type dependabot`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// TODO: Implement discovery, runner, and serialization pipeline
		// For now, return a stub result
		result := wetwire.BuildResult{
			Success:   false,
			Errors:    []string{"build command not yet implemented"},
			Workflows: []string{},
			Files:     []string{},
		}

		_ = path // Will be used when implemented

		if buildFormat == "json" {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(result)
		}

		if !result.Success {
			for _, err := range result.Errors {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
			}
			return fmt.Errorf("build failed")
		}

		for _, file := range result.Files {
			fmt.Printf("wrote %s\n", file)
		}
		return nil
	},
}

func init() {
	buildCmd.Flags().StringVarP(&buildOutput, "output", "o", ".github/workflows", "output directory")
	buildCmd.Flags().StringVar(&buildFormat, "format", "yaml", "output format (yaml, json)")
	buildCmd.Flags().StringVar(&buildType, "type", "workflow", "config type (workflow, dependabot, issue-template)")
}
