package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
)

var importOutput string
var importSingleFile bool
var importNoScaffold bool
var importType string

var importCmd = &cobra.Command{
	Use:   "import <workflow.yml>",
	Short: "Convert existing YAML to Go code",
	Long: `Import converts existing GitHub Actions YAML files to typed Go declarations.

Example:
  wetwire-github import .github/workflows/ci.yml -o my-workflows/
  wetwire-github import ci.yml --single-file
  wetwire-github import .github/dependabot.yml --type dependabot`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: file not found: %s\n", path)
			os.Exit(1)
		}

		// TODO: Implement YAML parsing and Go code generation
		result := wetwire.ImportResult{
			Success:   false,
			OutputDir: importOutput,
			Errors:    []string{"import command not yet implemented"},
		}

		_ = importSingleFile // Will be used when implemented
		_ = importNoScaffold // Will be used when implemented
		_ = importType       // Will be used when implemented

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result); err != nil {
			return err
		}

		if !result.Success {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	importCmd.Flags().StringVarP(&importOutput, "output", "o", ".", "output directory")
	importCmd.Flags().BoolVar(&importSingleFile, "single-file", false, "generate all code in a single file")
	importCmd.Flags().BoolVar(&importNoScaffold, "no-scaffold", false, "skip generating go.mod, README, etc.")
	importCmd.Flags().StringVar(&importType, "type", "workflow", "config type (workflow, dependabot, issue-template)")
}
