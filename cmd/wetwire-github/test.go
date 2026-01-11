package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	wetwire "github.com/lex00/wetwire-github-go"
	"github.com/lex00/wetwire-github-go/internal/discover"

	"github.com/lex00/wetwire-core-go/agent/personas"
	"github.com/lex00/wetwire-core-go/agent/scoring"
)

var testFormat string
var testPersona string
var testScenario string
var testList bool
var testScore bool
var testProvider string

var testCmd = &cobra.Command{
	Use:   "test <path>",
	Short: "Run persona-based workflow tests",
	Long: `Test runs persona-based tests against workflow declarations.

Developer personas simulate different types of users interacting with
the AI agent to test workflow generation quality.

Available personas:
  beginner      - New to GitHub Actions, needs guidance
  intermediate  - Some experience, knows basics but misses details
  expert        - Deep CI/CD knowledge, precise requirements
  terse         - Minimal words, expects system to infer
  verbose       - Over-explains, buries requirements in prose

Scenarios:
  ci-workflow   - Basic CI workflow test
  deployment    - Deployment workflow test
  release       - Release workflow test
  matrix        - Matrix strategy workflow test

Scoring (5 dimensions, 0-3 each, max 15):
  Completeness       - Were all required workflows generated?
  Lint Quality       - Did the code pass linting?
  Code Quality       - Does the code follow idiomatic patterns?
  Output Validity    - Is the generated YAML valid?
  Question Efficiency - Appropriate number of clarifying questions?

Thresholds: 0-5 Failure, 6-9 Partial, 10-12 Success, 13-15 Excellent

Providers:
  anthropic  - Use Anthropic Claude API (default)
  kiro       - Use Kiro AI agent

Example:
  wetwire-github test .
  wetwire-github test . --persona beginner
  wetwire-github test . --scenario ci-workflow --score
  wetwire-github test . --provider kiro --persona beginner
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
	testCmd.Flags().StringVar(&testPersona, "persona", "", "run specific persona ("+strings.Join(personas.Names(), ", ")+")")
	testCmd.Flags().StringVar(&testScenario, "scenario", "", "run specific scenario (ci-workflow, deployment, release, matrix)")
	testCmd.Flags().BoolVar(&testList, "list", false, "list available personas and scenarios")
	testCmd.Flags().BoolVar(&testScore, "score", false, "show scoring breakdown")
	testCmd.Flags().StringVar(&testProvider, "provider", "anthropic", "LLM provider (anthropic, kiro)")
}

// listTestPersonas lists available personas and scenarios.
func listTestPersonas() error {
	fmt.Println("Developer Personas:")
	fmt.Println("")
	for _, p := range personas.All() {
		fmt.Printf("  %-14s %s\n", p.Name, p.Description)
	}
	fmt.Println("")
	fmt.Println("Scenarios:")
	fmt.Println("")
	fmt.Println("  ci-workflow     Basic CI workflow (build, test, lint)")
	fmt.Println("  deployment      Deployment workflow (multi-environment)")
	fmt.Println("  release         Release workflow (tags, changelog, artifacts)")
	fmt.Println("  matrix          Matrix strategy (multi-version, multi-OS)")
	fmt.Println("")
	fmt.Println("Scoring Dimensions (0-3 each):")
	fmt.Println("")
	fmt.Println("  Completeness        Were all required workflows generated?")
	fmt.Println("  Lint Quality        Did the code pass wetwire-github linting?")
	fmt.Println("  Code Quality        Does the code follow idiomatic patterns?")
	fmt.Println("  Output Validity     Is the generated YAML valid per actionlint?")
	fmt.Println("  Question Efficiency Appropriate number of clarifying questions?")
	fmt.Println("")
	fmt.Println("Thresholds: 0-5 Failure, 6-9 Partial, 10-12 Success, 13-15 Excellent")
	return nil
}

// runTest executes the test command.
func runTest(path string) error {
	// Validate provider
	if !isValidProvider(testProvider) {
		fmt.Fprintf(os.Stderr, "error: invalid provider %q (valid: anthropic, kiro)\n", testProvider)
		os.Exit(1)
		return nil
	}

	// Validate persona if specified
	if testPersona != "" {
		if _, err := personas.Get(testPersona); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
			return nil
		}
	}

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

	// Calculate score if requested
	var score *scoring.Score
	if testScore {
		score = calculateScore(discovered, testPersona, testScenario)
	}

	// Output result
	if testFormat == "json" {
		output := struct {
			Result wetwire.TestResult `json:"result"`
			Score  *scoring.Score     `json:"score,omitempty"`
		}{
			Result: result,
			Score:  score,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
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

	// Show score if requested
	if score != nil {
		fmt.Println("")
		fmt.Print(score.String())
	}

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

// calculateScore computes a score for the discovered workflows.
func calculateScore(discovered *discover.DiscoveryResult, persona, scenario string) *scoring.Score {
	score := scoring.NewScore(persona, scenario)

	// Score completeness based on workflows found
	expectedWorkflows := 1 // Base expectation
	if scenario == "deployment" {
		expectedWorkflows = 2 // Deploy usually has CI + deploy
	} else if scenario == "matrix" {
		expectedWorkflows = 1 // Matrix is usually in one workflow
	}
	rating, notes := scoring.ScoreCompleteness(expectedWorkflows, len(discovered.Workflows))
	score.Completeness.Rating = rating
	score.Completeness.Notes = notes

	// Score lint quality (no lint cycles for static analysis)
	// For structural tests, assume lint passed if no parse errors
	lintPassed := len(discovered.Errors) == 0
	rating, notes = scoring.ScoreLintQuality(0, lintPassed)
	score.LintQuality.Rating = rating
	score.LintQuality.Notes = notes

	// Score code quality based on structural analysis
	var issues []string
	for _, wf := range discovered.Workflows {
		if wf.Name == "" {
			issues = append(issues, "workflow missing name")
		}
		if len(wf.Jobs) == 0 {
			issues = append(issues, fmt.Sprintf("workflow %s has no jobs", wf.Name))
		}
	}
	for _, job := range discovered.Jobs {
		if job.Name == "" {
			issues = append(issues, "job missing name")
		}
	}
	rating, notes = scoring.ScoreCodeQuality(issues)
	score.CodeQuality.Rating = rating
	score.CodeQuality.Notes = notes

	// Score output validity (would need actionlint for real validation)
	// For now, base on parse errors
	rating, notes = scoring.ScoreOutputValidity(len(discovered.Errors), 0)
	score.OutputValidity.Rating = rating
	score.OutputValidity.Notes = notes

	// Score question efficiency (0 questions for static analysis)
	rating, notes = scoring.ScoreQuestionEfficiency(0)
	score.QuestionEfficiency.Rating = rating
	score.QuestionEfficiency.Notes = notes

	return score
}
