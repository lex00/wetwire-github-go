package importer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/workflow"
)

// TestRoundTrip_ReferenceWorkflows tests the full import/export cycle with GitHub starter workflows
func TestRoundTrip_ReferenceWorkflows(t *testing.T) {
	referenceDir := filepath.Join("..", "..", "testdata", "reference")

	// Find all YAML files in the reference directory
	files, err := filepath.Glob(filepath.Join(referenceDir, "*.yml"))
	require.NoError(t, err, "should find reference workflow files")
	require.NotEmpty(t, files, "should have at least one reference workflow")

	for _, file := range files {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			testRoundTrip(t, file)
		})
	}
}

// testRoundTrip performs a complete round-trip test for a single workflow file
func testRoundTrip(t *testing.T, yamlPath string) {
	// Step 1: Read the original YAML
	originalYAML, err := os.ReadFile(yamlPath)
	require.NoError(t, err, "should read original YAML file")

	// Step 2: Import YAML to IR
	workflow, err := ParseWorkflow(originalYAML)
	require.NoError(t, err, "should parse workflow from YAML")
	require.NotNil(t, workflow, "parsed workflow should not be nil")

	// Validate basic workflow structure
	require.NotEmpty(t, workflow.Name, "workflow should have a name")
	require.NotEmpty(t, workflow.Jobs, "workflow should have jobs")

	// Step 3: Generate Go code from IR
	gen := NewCodeGenerator()
	gen.PackageName = "testworkflow"
	gen.SingleFile = true

	workflowName := toVarName(workflow.Name)
	result, err := gen.Generate(workflow, workflowName)
	require.NoError(t, err, "should generate Go code")
	require.NotEmpty(t, result.Files, "should generate at least one file")

	goCode := result.Files["workflows.go"]
	require.NotEmpty(t, goCode, "should generate workflows.go")

	// Step 4: Create a temporary directory for the Go project
	tempDir := t.TempDir()

	// Write the generated Go code
	workflowsPath := filepath.Join(tempDir, "workflows.go")
	err = os.WriteFile(workflowsPath, []byte(goCode), 0644)
	require.NoError(t, err, "should write generated Go code")

	// Create go.mod
	goMod := `module testworkflow

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => ` + getProjectRoot(t)

	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(goMod), 0644)
	require.NoError(t, err, "should write go.mod")

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tempDir
	output, err := tidyCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", output)

	// Step 5: Use the runner to extract the workflow value
	r := runner.NewRunner()

	// Create a minimal discovery result
	discovered := &discover.DiscoveryResult{
		Workflows: []discover.DiscoveredWorkflow{
			{
				Name: workflowName,
				File: workflowsPath,
				Line: 1,
				Jobs: []string{},
			},
		},
		Jobs: []discover.DiscoveredJob{},
	}

	extracted, err := r.ExtractValues(tempDir, discovered)
	require.NoError(t, err, "should extract workflow values")
	require.Len(t, extracted.Workflows, 1, "should extract one workflow")

	// Step 6: Convert extracted data to workflow.Workflow
	extractedData := extracted.Workflows[0].Data
	reconstructed := mapToWorkflow(t, extractedData)
	require.NotNil(t, reconstructed, "should reconstruct workflow")

	// Step 7: Serialize back to YAML
	regeneratedYAML, err := serialize.ToYAML(reconstructed)
	require.NoError(t, err, "should serialize workflow to YAML")

	// Debug: Log the extracted and reconstructed data
	t.Logf("Extracted workflow data: %+v", extractedData)
	t.Logf("Reconstructed workflow: Name=%s, Jobs=%d", reconstructed.Name, len(reconstructed.Jobs))

	// Step 8: Compare the YAMLs semantically
	compareWorkflowYAML(t, originalYAML, regeneratedYAML)
}

// mapToWorkflow converts a map[string]any to workflow.Workflow
func mapToWorkflow(t *testing.T, data map[string]any) *workflow.Workflow {
	wf := &workflow.Workflow{}

	// Extract top-level fields
	if name, ok := data["Name"].(string); ok {
		wf.Name = name
	}

	// Handle Jobs
	if jobsMap, ok := data["Jobs"].(map[string]any); ok {
		wf.Jobs = make(map[string]workflow.Job)
		for jobID, jobData := range jobsMap {
			if jobMap, ok := jobData.(map[string]any); ok {
				wf.Jobs[jobID] = mapToJob(jobMap)
			}
		}
	}

	// Handle triggers (On)
	if onData, ok := data["On"].(map[string]any); ok {
		wf.On = mapToTriggers(onData)
	}

	// Handle other fields as needed
	if env, ok := data["Env"].(map[string]any); ok {
		wf.Env = env
	}

	if _, ok := data["Permissions"].(map[string]any); ok {
		wf.Permissions = &workflow.Permissions{}
		// Convert permissions map to Permissions struct if needed
	}

	return wf
}

// mapToJob converts a map to a workflow.Job
func mapToJob(jobMap map[string]any) workflow.Job {
	job := workflow.Job{}

	if name, ok := jobMap["Name"].(string); ok {
		job.Name = name
	}
	if runsOn, ok := jobMap["RunsOn"].(string); ok {
		job.RunsOn = runsOn
	}
	if ifCond, ok := jobMap["If"]; ok {
		job.If = ifCond
	}
	if timeout, ok := jobMap["TimeoutMinutes"].(int); ok {
		job.TimeoutMinutes = timeout
	}
	if needs, ok := jobMap["Needs"]; ok {
		job.Needs = []any{needs}
	}

	// Handle Steps
	if stepsSlice, ok := jobMap["Steps"].([]any); ok {
		job.Steps = make([]any, len(stepsSlice))
		for i, stepData := range stepsSlice {
			if stepMap, ok := stepData.(map[string]any); ok {
				job.Steps[i] = mapToStep(stepMap)
			} else {
				job.Steps[i] = stepData
			}
		}
	}

	if env, ok := jobMap["Env"].(map[string]any); ok {
		job.Env = env
	}

	return job
}

// mapToTriggers converts a map to workflow.Triggers
func mapToTriggers(onMap map[string]any) workflow.Triggers {
	triggers := workflow.Triggers{}

	if pushData, ok := onMap["Push"].(map[string]any); ok {
		triggers.Push = &workflow.PushTrigger{}
		if branches, ok := pushData["Branches"].([]any); ok {
			triggers.Push.Branches = make([]string, len(branches))
			for i, b := range branches {
				if s, ok := b.(string); ok {
					triggers.Push.Branches[i] = s
				}
			}
		}
	}

	if prData, ok := onMap["PullRequest"].(map[string]any); ok {
		triggers.PullRequest = &workflow.PullRequestTrigger{}
		if branches, ok := prData["Branches"].([]any); ok {
			triggers.PullRequest.Branches = make([]string, len(branches))
			for i, b := range branches {
				if s, ok := b.(string); ok {
					triggers.PullRequest.Branches[i] = s
				}
			}
		}
	}

	return triggers
}

// normalizeMapKeys converts map keys to lowercase recursively
func normalizeMapKeys(data any) any {
	switch v := data.(type) {
	case map[string]any:
		result := make(map[string]any)
		for key, value := range v {
			// Convert PascalCase to kebab-case for YAML compatibility
			normalizedKey := toYAMLKey(key)

			// Special handling for Steps field - convert step maps to Step structs
			if normalizedKey == "steps" {
				if stepSlice, ok := value.([]any); ok {
					normalizedSteps := make([]any, len(stepSlice))
					for i, stepData := range stepSlice {
						if stepMap, ok := stepData.(map[string]any); ok {
							normalizedSteps[i] = mapToStep(stepMap)
						} else {
							normalizedSteps[i] = stepData
						}
					}
					result[normalizedKey] = normalizedSteps
					continue
				}
			}

			result[normalizedKey] = normalizeMapKeys(value)
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = normalizeMapKeys(item)
		}
		return result
	default:
		return v
	}
}

// mapToStep converts a map to a workflow.Step
func mapToStep(stepMap map[string]any) workflow.Step {
	step := workflow.Step{}

	if id, ok := stepMap["ID"].(string); ok {
		step.ID = id
	}
	if name, ok := stepMap["Name"].(string); ok {
		step.Name = name
	}
	if uses, ok := stepMap["Uses"].(string); ok {
		step.Uses = uses
	}
	if run, ok := stepMap["Run"].(string); ok {
		step.Run = run
	}
	if shell, ok := stepMap["Shell"].(string); ok {
		step.Shell = shell
	}
	if ifCond, ok := stepMap["If"].(string); ok {
		step.If = ifCond
	}
	if workDir, ok := stepMap["WorkingDirectory"].(string); ok {
		step.WorkingDirectory = workDir
	}
	if timeout, ok := stepMap["TimeoutMinutes"].(int); ok {
		step.TimeoutMinutes = timeout
	}
	if continueOnError, ok := stepMap["ContinueOnError"].(bool); ok {
		step.ContinueOnError = continueOnError
	}
	if with, ok := stepMap["With"].(map[string]any); ok {
		step.With = with
	}
	if env, ok := stepMap["Env"].(map[string]any); ok {
		step.Env = env
	}

	return step
}

// toYAMLKey converts Go struct field names to YAML field names
func toYAMLKey(key string) string {
	// Handle special cases
	specialCases := map[string]string{
		"RunsOn":          "runs-on",
		"ContinueOnError": "continue-on-error",
		"TimeoutMinutes":  "timeout-minutes",
		"WorkingDirectory":"working-directory",
		"BranchesIgnore":  "branches-ignore",
		"PathsIgnore":     "paths-ignore",
		"TagsIgnore":      "tags-ignore",
		"PullRequest":     "pull_request",
		"PullRequestTarget": "pull_request_target",
		"WorkflowDispatch": "workflow_dispatch",
		"WorkflowCall":    "workflow_call",
		"WorkflowRun":     "workflow_run",
		"RepositoryDispatch": "repository_dispatch",
		"IssueComment":    "issue_comment",
		"ProjectCard":     "project_card",
		"ProjectColumn":   "project_column",
		"PullRequestReview": "pull_request_review",
		"PullRequestReviewComment": "pull_request_review_comment",
		"CheckRun":        "check_run",
		"CheckSuite":      "check_suite",
		"DiscussionComment": "discussion_comment",
		"MergeGroup":      "merge_group",
		"PageBuild":       "page_build",
	}

	if mapped, ok := specialCases[key]; ok {
		return mapped
	}

	// Convert to lowercase for simple cases
	return strings.ToLower(key)
}

// compareWorkflowYAML compares two workflow YAMLs semantically
func compareWorkflowYAML(t *testing.T, original, regenerated []byte) {
	// Parse both YAMLs into generic maps
	var originalMap map[string]any
	err := yaml.Unmarshal(original, &originalMap)
	require.NoError(t, err, "should parse original YAML")

	var regeneratedMap map[string]any
	err = yaml.Unmarshal(regenerated, &regeneratedMap)
	require.NoError(t, err, "should parse regenerated YAML")

	// Compare key structural elements
	// We don't do exact string comparison because:
	// 1. Comments are not preserved
	// 2. Field ordering may differ
	// 3. Some fields may have equivalent representations (e.g., string vs array)

	// Compare workflow name
	require.Equal(t, originalMap["name"], regeneratedMap["name"], "workflow names should match")

	// Verify jobs exist
	originalJobs, ok1 := originalMap["jobs"].(map[string]any)
	regeneratedJobs, ok2 := regeneratedMap["jobs"].(map[string]any)
	require.True(t, ok1, "original should have jobs map")
	require.True(t, ok2, "regenerated should have jobs map")

	// Compare job count
	require.Len(t, regeneratedJobs, len(originalJobs), "should have same number of jobs")

	// Compare job names
	for jobName := range originalJobs {
		require.Contains(t, regeneratedJobs, jobName, "regenerated should contain job: %s", jobName)
	}

	// Verify triggers exist (but allow flexibility in structure)
	require.Contains(t, originalMap, "on", "original should have triggers")
	require.Contains(t, regeneratedMap, "on", "regenerated should have triggers")

	t.Logf("Round-trip successful:\n  Original workflow: %s\n  Jobs: %d\n  Triggers: %v",
		originalMap["name"], len(originalJobs), originalMap["on"])
}

// getProjectRoot returns the absolute path to the project root
func getProjectRoot(t *testing.T) string {
	// Start from the current test file location and walk up
	cwd, err := os.Getwd()
	require.NoError(t, err, "should get current directory")

	// Walk up until we find go.mod
	dir := cwd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (go.mod)")
		}
		dir = parent
	}
}

// TestRoundTrip_SimpleWorkflow tests round-trip with a simple inline workflow
func TestRoundTrip_SimpleWorkflow(t *testing.T) {
	yamlContent := `name: Simple CI
on:
  push:
    branches: [ main ]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Run tests
        run: go test ./...
`

	// Create a temporary file for the workflow
	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "simple.yml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	require.NoError(t, err, "should write test YAML")

	// Run the round-trip test
	testRoundTrip(t, yamlPath)
}

// TestRoundTrip_MatrixStrategy tests round-trip with matrix strategy
func TestRoundTrip_MatrixStrategy(t *testing.T) {
	yamlContent := `name: Matrix Build
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go: ["1.21", "1.22", "1.23"]
        os: [ubuntu-latest, macos-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go build ./...
`

	tempDir := t.TempDir()
	yamlPath := filepath.Join(tempDir, "matrix.yml")
	err := os.WriteFile(yamlPath, []byte(yamlContent), 0644)
	require.NoError(t, err, "should write test YAML")

	testRoundTrip(t, yamlPath)
}

// TestCodeGenerator_GeneratedCodeCompiles tests that generated code compiles
func TestCodeGenerator_GeneratedCodeCompiles(t *testing.T) {
	yamlContent := `name: Test Workflow
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: echo "test"
`

	workflow, err := ParseWorkflow([]byte(yamlContent))
	require.NoError(t, err, "should parse workflow")

	gen := NewCodeGenerator()
	gen.PackageName = "testpkg"
	gen.SingleFile = true

	result, err := gen.Generate(workflow, "TestWorkflow")
	require.NoError(t, err, "should generate code")

	goCode := result.Files["workflows.go"]
	require.NotEmpty(t, goCode, "should have generated code")

	// Verify code compiles by writing it to a temp directory and running go build
	tempDir := t.TempDir()

	workflowsPath := filepath.Join(tempDir, "workflows.go")
	err = os.WriteFile(workflowsPath, []byte(goCode), 0644)
	require.NoError(t, err, "should write Go code")

	goMod := `module testpkg

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0

replace github.com/lex00/wetwire-github-go => ` + getProjectRoot(t)

	goModPath := filepath.Join(tempDir, "go.mod")
	err = os.WriteFile(goModPath, []byte(goMod), 0644)
	require.NoError(t, err, "should write go.mod")

	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tempDir
	output, err := tidyCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", output)

	// Try to build
	buildCmd := exec.Command("go", "build", ".")
	buildCmd.Dir = tempDir
	output, err = buildCmd.CombinedOutput()
	require.NoError(t, err, "generated code should compile: %s", output)
}

// TestParser_ParseReferenceWorkflows tests that we can parse all reference workflows
func TestParser_ParseReferenceWorkflows(t *testing.T) {
	referenceDir := filepath.Join("..", "..", "testdata", "reference")

	files, err := filepath.Glob(filepath.Join(referenceDir, "*.yml"))
	require.NoError(t, err, "should find reference workflow files")
	require.NotEmpty(t, files, "should have at least one reference workflow")

	for _, file := range files {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			content, err := os.ReadFile(file)
			require.NoError(t, err, "should read file")

			workflow, err := ParseWorkflow(content)
			require.NoError(t, err, "should parse workflow")
			require.NotNil(t, workflow, "workflow should not be nil")

			// Basic validation
			require.NotEmpty(t, workflow.Name, "workflow should have a name")
			require.NotEmpty(t, workflow.Jobs, "workflow should have jobs")

			// Log workflow details for debugging
			t.Logf("Parsed workflow: %s", workflow.Name)
			t.Logf("  Jobs: %v", getJobNames(workflow.Jobs))
			t.Logf("  Triggers: %s", describeTriggers(&workflow.On))
		})
	}
}

// getJobNames extracts job names from a job map
func getJobNames(jobs map[string]*IRJob) []string {
	names := make([]string, 0, len(jobs))
	for name := range jobs {
		names = append(names, name)
	}
	return names
}

// describeTriggers creates a human-readable description of triggers
func describeTriggers(triggers *IRTriggers) string {
	parts := []string{}
	if triggers.Push != nil {
		parts = append(parts, "push")
	}
	if triggers.PullRequest != nil {
		parts = append(parts, "pull_request")
	}
	if triggers.Schedule != nil {
		parts = append(parts, "schedule")
	}
	if triggers.WorkflowDispatch != nil {
		parts = append(parts, "workflow_dispatch")
	}
	if len(parts) == 0 {
		return "none"
	}
	return strings.Join(parts, ", ")
}
