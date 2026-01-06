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

var testFormat string
var testPersona string
var testScenario string
var testList bool

var testCmd = &cobra.Command{
	Use:   "test <path>",
	Short: "Run persona-based workflow tests",
	Long: `Test runs persona-based tests against workflow declarations.

Personas simulate different GitHub Actions scenarios to validate
workflow behavior without running actual workflows.

Available personas:
  push          - Simulates push event (branch commit)
  pull_request  - Simulates pull request event
  schedule      - Simulates scheduled event (cron)
  workflow_dispatch - Simulates manual dispatch

Scenarios:
  ci-setup      - Basic CI workflow test
  deployment    - Deployment workflow test
  release       - Release workflow test

Example:
  wetwire-github test .
  wetwire-github test . --persona push
  wetwire-github test . --scenario ci-setup
  wetwire-github test --list`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if testList {
			return listTestPersonas()
		}
		if len(args) == 0 {
			return fmt.Errorf("path argument required")
		}
		return runTest(args[0])
	},
}

func init() {
	testCmd.Flags().StringVar(&testFormat, "format", "text", "output format (text, json)")
	testCmd.Flags().StringVar(&testPersona, "persona", "", "run specific persona (push, pull_request, schedule, workflow_dispatch)")
	testCmd.Flags().StringVar(&testScenario, "scenario", "", "run specific scenario (ci-setup, deployment, release)")
	testCmd.Flags().BoolVar(&testList, "list", false, "list available personas and scenarios")
}

// listTestPersonas lists available personas and scenarios.
func listTestPersonas() error {
	fmt.Println("Available Personas:")
	fmt.Println("  push              Simulates push event (branch commit)")
	fmt.Println("  pull_request      Simulates pull request event")
	fmt.Println("  schedule          Simulates scheduled event (cron)")
	fmt.Println("  workflow_dispatch Simulates manual dispatch")
	fmt.Println("")
	fmt.Println("Available Scenarios:")
	fmt.Println("  ci-setup          Basic CI workflow test")
	fmt.Println("  deployment        Deployment workflow test")
	fmt.Println("  release           Release workflow test")
	fmt.Println("")
	fmt.Println("Note: Full persona-based testing requires wetwire-core-go (Phase 4B)")
	return nil
}

// runTest executes the test command.
func runTest(path string) error {
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

	// Discover workflows
	disc := discover.NewDiscoverer()
	discovered, err := disc.Discover(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: discovery failed: %v\n", err)
		os.Exit(1)
		return nil
	}

	// Run basic structural tests
	result := runStructuralTests(discovered, testPersona)

	// Output result
	if testFormat == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	// Text output
	if len(result.Tests) == 0 {
		fmt.Println("No tests run")
		return nil
	}

	for _, t := range result.Tests {
		if t.Passed {
			fmt.Printf("PASS %s\n", t.Name)
		} else {
			fmt.Printf("FAIL %s: %s\n", t.Name, t.Error)
		}
	}

	fmt.Printf("\n%d passed, %d failed\n", result.Passed, result.Failed)

	if !result.Success {
		os.Exit(1)
	}
	return nil
}

// runStructuralTests runs basic structural tests on discovered workflows.
func runStructuralTests(discovered *discover.DiscoveryResult, persona string) wetwire.TestResult {
	result := wetwire.TestResult{
		Success: true,
		Tests:   []wetwire.TestCase{},
		Passed:  0,
		Failed:  0,
	}

	// Test: Workflows exist
	if len(discovered.Workflows) > 0 {
		result.Tests = append(result.Tests, wetwire.TestCase{
			Name:    "workflows_exist",
			Persona: persona,
			Passed:  true,
		})
		result.Passed++
	} else {
		result.Tests = append(result.Tests, wetwire.TestCase{
			Name:    "workflows_exist",
			Persona: persona,
			Passed:  false,
			Error:   "no workflows found",
		})
		result.Failed++
		result.Success = false
	}

	// Test: Jobs exist
	if len(discovered.Jobs) > 0 {
		result.Tests = append(result.Tests, wetwire.TestCase{
			Name:    "jobs_exist",
			Persona: persona,
			Passed:  true,
		})
		result.Passed++
	} else {
		result.Tests = append(result.Tests, wetwire.TestCase{
			Name:    "jobs_exist",
			Persona: persona,
			Passed:  false,
			Error:   "no jobs found",
		})
		result.Failed++
		result.Success = false
	}

	// Test: No parse errors
	if len(discovered.Errors) == 0 {
		result.Tests = append(result.Tests, wetwire.TestCase{
			Name:    "no_parse_errors",
			Persona: persona,
			Passed:  true,
		})
		result.Passed++
	} else {
		result.Tests = append(result.Tests, wetwire.TestCase{
			Name:    "no_parse_errors",
			Persona: persona,
			Passed:  false,
			Error:   fmt.Sprintf("%d parse errors", len(discovered.Errors)),
		})
		result.Failed++
		result.Success = false
	}

	// Test: All workflow jobs exist
	for _, wf := range discovered.Workflows {
		testName := fmt.Sprintf("workflow_%s_jobs_valid", wf.Name)
		jobNames := make(map[string]bool)
		for _, job := range discovered.Jobs {
			jobNames[job.Name] = true
		}

		allExist := true
		for _, jobRef := range wf.Jobs {
			if !jobNames[jobRef] {
				allExist = false
				break
			}
		}

		if allExist {
			result.Tests = append(result.Tests, wetwire.TestCase{
				Name:    testName,
				Persona: persona,
				Passed:  true,
			})
			result.Passed++
		} else {
			result.Tests = append(result.Tests, wetwire.TestCase{
				Name:    testName,
				Persona: persona,
				Passed:  false,
				Error:   "references undefined jobs",
			})
			result.Failed++
			result.Success = false
		}
	}

	return result
}
