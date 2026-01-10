package importer

import (
	"strings"
	"testing"
)

func TestToVarName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// Basic cases
		{"CI", "CI"},
		{"build", "Build"},
		{"my-workflow", "MyWorkflow"},
		{"my_workflow", "MyWorkflow"},
		{"my workflow", "MyWorkflow"},

		// Special characters that need sanitization
		{"C/C++ CI", "CCppCI"},
		{"C++ Build", "CppBuild"},
		{"Node.js CI", "NodeJsCI"},
		{"iOS Build", "IOSBuild"},
		{"D CI", "DCI"},

		// Parentheses
		{"Build (Linux)", "BuildLinux"},

		// Multiple special chars
		{"C/C++/Objective-C", "CCppObjectiveC"},

		// Reserved words
		{"type", "TypeJob"},
		{"go", "GoJob"},
		{"import", "ImportJob"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toVarName(tt.input)
			if result != tt.expected {
				t.Errorf("toVarName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToVarNameProducesValidIdentifier(t *testing.T) {
	// These are real workflow names from actions/starter-workflows
	names := []string{
		"C/C++ CI",
		"D",
		"iOS",
		"Objective-C Xcode",
		".NET",
		"Node.js",
		"R",
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			result := toVarName(name)

			// Check it's not empty
			if result == "" {
				t.Errorf("toVarName(%q) returned empty string", name)
				return
			}

			// Check first char is letter (Go identifier requirement)
			first := rune(result[0])
			if !((first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z')) {
				t.Errorf("toVarName(%q) = %q, first char is not a letter", name, result)
			}

			// Check all chars are valid Go identifier chars
			for i, r := range result {
				valid := (r >= 'A' && r <= 'Z') ||
				         (r >= 'a' && r <= 'z') ||
				         (r >= '0' && r <= '9' && i > 0) ||
				         r == '_'
				if !valid {
					t.Errorf("toVarName(%q) = %q, char %q at pos %d is invalid", name, result, string(r), i)
				}
			}
		})
	}
}

func TestNewCodeGenerator(t *testing.T) {
	gen := NewCodeGenerator()
	if gen == nil {
		t.Fatal("NewCodeGenerator() returned nil")
	}
	if gen.PackageName != "workflows" {
		t.Errorf("PackageName = %q, want %q", gen.PackageName, "workflows")
	}
}

func TestCodeGenerator_Generate_SingleFile(t *testing.T) {
	gen := &CodeGenerator{
		PackageName: "workflows",
		SingleFile:  true,
	}

	workflow := &IRWorkflow{
		Name: "CI",
		On: IRTriggers{
			Push:        &IRPushTrigger{Branches: []string{"main"}},
			PullRequest: &IRPullRequestTrigger{Branches: []string{"main"}},
		},
		Jobs: map[string]*IRJob{
			"build": {
				Name:   "Build",
				RunsOn: "ubuntu-latest",
				Steps: []IRStep{
					{Uses: "actions/checkout@v4"},
					{Run: "go build ./..."},
				},
			},
		},
	}

	result, err := gen.Generate(workflow, "ci")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if result.Workflows != 1 {
		t.Errorf("Workflows = %d, want 1", result.Workflows)
	}
	if result.Jobs != 1 {
		t.Errorf("Jobs = %d, want 1", result.Jobs)
	}
	if result.Steps != 2 {
		t.Errorf("Steps = %d, want 2", result.Steps)
	}

	code := result.Files["workflows.go"]
	if code == "" {
		t.Fatal("workflows.go is empty")
	}

	// Verify the code contains expected elements
	if !strings.Contains(code, "package workflows") {
		t.Error("Missing package declaration")
	}
	if !strings.Contains(code, "var Ci = workflow.Workflow{") {
		t.Error("Missing workflow variable")
	}
	if !strings.Contains(code, "var CiPush = workflow.PushTrigger{") {
		t.Error("Missing push trigger variable")
	}
	if !strings.Contains(code, "var Build = workflow.Job{") {
		t.Error("Missing job variable")
	}
}

func TestCodeGenerator_Generate_SeparateFiles(t *testing.T) {
	gen := &CodeGenerator{
		PackageName: "workflows",
		SingleFile:  false,
	}

	workflow := &IRWorkflow{
		Name: "CI",
		On: IRTriggers{
			Push: &IRPushTrigger{Branches: []string{"main", "develop"}},
		},
		Jobs: map[string]*IRJob{
			"test": {
				Name:   "Test",
				RunsOn: "ubuntu-latest",
				Steps: []IRStep{
					{Run: "go test ./..."},
				},
			},
		},
	}

	result, err := gen.Generate(workflow, "ci")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Check that separate files were created
	expectedFiles := []string{"workflows.go", "triggers.go", "jobs.go", "steps.go"}
	for _, filename := range expectedFiles {
		if _, ok := result.Files[filename]; !ok {
			t.Errorf("Missing expected file: %s", filename)
		}
	}
}

func TestCodeGenerator_GenerateWorkflow(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	workflow := &IRWorkflow{
		Name: "Test Workflow",
		Jobs: map[string]*IRJob{
			"build": {Name: "Build"},
			"test":  {Name: "Test"},
		},
	}

	code := gen.generateWorkflow(workflow, "TestWorkflow")

	if !strings.Contains(code, "var TestWorkflow = workflow.Workflow{") {
		t.Error("Missing workflow variable declaration")
	}
	if !strings.Contains(code, `Name: "Test Workflow"`) {
		t.Error("Missing workflow name")
	}
	if !strings.Contains(code, "On:   TestWorkflowTriggers") {
		t.Error("Missing triggers reference")
	}
	if !strings.Contains(code, `"build": Build`) {
		t.Error("Missing build job reference")
	}
	if !strings.Contains(code, `"test": Test`) {
		t.Error("Missing test job reference")
	}
}

func TestCodeGenerator_GenerateTriggers(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	workflow := &IRWorkflow{
		On: IRTriggers{
			Push: &IRPushTrigger{
				Branches:       []string{"main", "develop"},
				BranchesIgnore: []string{"staging"},
				Tags:           []string{"v*"},
				Paths:          []string{"src/**"},
			},
			PullRequest: &IRPullRequestTrigger{
				Branches: []string{"main"},
				Types:    []string{"opened", "synchronize"},
				Paths:    []string{"src/**", "tests/**"},
			},
			WorkflowDispatch: &IRWorkflowDispatch{},
			Schedule: []IRSchedule{
				{Cron: "0 0 * * *"},
			},
		},
	}

	code := gen.generateTriggers(workflow, "CI")

	// Check push trigger
	if !strings.Contains(code, "var CIPush = workflow.PushTrigger{") {
		t.Error("Missing push trigger variable")
	}
	if !strings.Contains(code, `Branches: []string{"main", "develop"}`) {
		t.Error("Missing push branches")
	}
	if !strings.Contains(code, `BranchesIgnore: []string{"staging"}`) {
		t.Error("Missing push branches-ignore")
	}
	if !strings.Contains(code, `Tags: []string{"v*"}`) {
		t.Error("Missing push tags")
	}

	// Check PR trigger
	if !strings.Contains(code, "var CIPullRequest = workflow.PullRequestTrigger{") {
		t.Error("Missing PR trigger variable")
	}
	if !strings.Contains(code, `Types: []string{"opened", "synchronize"}`) {
		t.Error("Missing PR types")
	}

	// Check main triggers struct
	if !strings.Contains(code, "var CITriggers = workflow.Triggers{") {
		t.Error("Missing triggers variable")
	}
	if !strings.Contains(code, "WorkflowDispatch: &workflow.WorkflowDispatchTrigger{}") {
		t.Error("Missing workflow_dispatch")
	}
	if !strings.Contains(code, `{Cron: "0 0 * * *"}`) {
		t.Error("Missing schedule cron")
	}
}

func TestCodeGenerator_GenerateJob(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	job := &IRJob{
		Name:           "Test Job",
		RunsOn:         "ubuntu-latest",
		Needs:          []any{"build"},
		If:             "success()",
		TimeoutMinutes: 30,
		Steps: []IRStep{
			{Run: "echo test"},
		},
	}

	code := gen.generateJob("test", job)

	if !strings.Contains(code, "var Test = workflow.Job{") {
		t.Error("Missing job variable")
	}
	if !strings.Contains(code, `Name: "Test Job"`) {
		t.Error("Missing job name")
	}
	if !strings.Contains(code, `RunsOn: "ubuntu-latest"`) {
		t.Error("Missing runs-on")
	}
	if !strings.Contains(code, `Needs: []any{"build"}`) {
		t.Error("Missing needs")
	}
	if !strings.Contains(code, `If: "success()"`) {
		t.Error("Missing if condition")
	}
	if !strings.Contains(code, "TimeoutMinutes: 30") {
		t.Error("Missing timeout")
	}
	if !strings.Contains(code, "Steps: TestSteps") {
		t.Error("Missing steps reference")
	}
}

func TestCodeGenerator_GenerateSteps(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	steps := []IRStep{
		{
			ID:   "checkout",
			Name: "Checkout",
			Uses: "actions/checkout@v4",
			With: map[string]any{
				"fetch-depth": 0,
				"submodules":  "recursive",
			},
		},
		{
			Name: "Run tests",
			Run:  "go test ./...",
			Env: map[string]any{
				"GO_VERSION": "1.23",
				"CGO_ENABLED": 0,
			},
			If:               "success()",
			WorkingDirectory: "./src",
			TimeoutMinutes:   10,
			Shell:            "bash",
		},
		{
			Name: "Multiline script",
			Run:  "echo 'line 1'\necho 'line 2'",
		},
	}

	code := gen.generateSteps("TestSteps", steps)

	if !strings.Contains(code, "var TestSteps = []workflow.Step{") {
		t.Error("Missing steps variable")
	}
	if !strings.Contains(code, `ID: "checkout"`) {
		t.Error("Missing step ID")
	}
	if !strings.Contains(code, `Uses: "actions/checkout@v4"`) {
		t.Error("Missing uses")
	}
	if !strings.Contains(code, `"fetch-depth": 0`) {
		t.Error("Missing with parameter")
	}
	if !strings.Contains(code, `Env: map[string]any{`) {
		t.Error("Missing env map")
	}
	if !strings.Contains(code, `If: "success()"`) {
		t.Error("Missing if condition")
	}
	if !strings.Contains(code, `WorkingDirectory: "./src"`) {
		t.Error("Missing working directory")
	}
	if !strings.Contains(code, "TimeoutMinutes: 10") {
		t.Error("Missing timeout")
	}
	if !strings.Contains(code, `Shell: "bash"`) {
		t.Error("Missing shell")
	}
	// Multiline strings should use backticks
	if !strings.Contains(code, "`echo 'line 1'\necho 'line 2'`") {
		t.Error("Missing multiline run with backticks")
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{"hello", `"hello"`},
		{123, "123"},
		{int64(456), "456"},
		{3.14, "3.14"},
		{true, "true"},
		{false, "false"},
	}

	for _, tt := range tests {
		result := formatValue(tt.input)
		if result != tt.expected {
			t.Errorf("formatValue(%v) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"CI", "ci"},
		{"Build Workflow", "build-workflow"},
		{"C/C++ CI", "c-c-ci"},
		{"Node.js CI", "node-js-ci"},
		{"--test--", "test"},
	}

	for _, tt := range tests {
		result := toFilename(tt.input)
		if result != tt.expected {
			t.Errorf("toFilename(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestCodeGenerator_GenerateWorkflowFile(t *testing.T) {
	gen := &CodeGenerator{PackageName: "myworkflows"}

	workflow := &IRWorkflow{
		Name: "Test",
		Jobs: map[string]*IRJob{"build": {Name: "Build"}},
	}

	code := gen.generateWorkflowFile(workflow, "Test")

	if !strings.Contains(code, "package myworkflows") {
		t.Error("Missing custom package name")
	}
}

func TestCodeGenerator_GenerateTriggersFile(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	workflow := &IRWorkflow{
		On: IRTriggers{
			WorkflowCall: &IRWorkflowCall{},
		},
	}

	code := gen.generateTriggersFile(workflow, "Reusable")

	if !strings.Contains(code, "package workflows") {
		t.Error("Missing package declaration")
	}
	if !strings.Contains(code, "WorkflowCall: &workflow.WorkflowCallTrigger{}") {
		t.Error("Missing workflow_call trigger")
	}
}

func TestCodeGenerator_GenerateJobsFile(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	workflow := &IRWorkflow{
		Jobs: map[string]*IRJob{
			"build": {Name: "Build", RunsOn: "ubuntu-latest"},
			"test":  {Name: "Test", RunsOn: "ubuntu-latest"},
		},
	}

	code := gen.generateJobsFile(workflow)

	if !strings.Contains(code, "package workflows") {
		t.Error("Missing package declaration")
	}
	// Jobs should be sorted alphabetically
	buildIdx := strings.Index(code, "var Build")
	testIdx := strings.Index(code, "var Test")
	if buildIdx == -1 || testIdx == -1 || buildIdx > testIdx {
		t.Error("Jobs not properly sorted or missing")
	}
}

func TestCodeGenerator_GenerateStepsFile(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	workflow := &IRWorkflow{
		Jobs: map[string]*IRJob{
			"build": {
				Steps: []IRStep{
					{Run: "echo build"},
				},
			},
			"test": {
				Steps: []IRStep{
					{Run: "echo test"},
				},
			},
		},
	}

	code := gen.generateStepsFile(workflow)

	if !strings.Contains(code, "package workflows") {
		t.Error("Missing package declaration")
	}
	if !strings.Contains(code, "var BuildSteps = []workflow.Step{") {
		t.Error("Missing BuildSteps")
	}
	if !strings.Contains(code, "var TestSteps = []workflow.Step{") {
		t.Error("Missing TestSteps")
	}
}

func TestCodeGenerator_EmptyJobs(t *testing.T) {
	gen := &CodeGenerator{PackageName: "workflows"}

	workflow := &IRWorkflow{
		Name: "Empty",
		Jobs: map[string]*IRJob{},
	}

	code := gen.generateWorkflow(workflow, "Empty")

	// Should handle empty jobs gracefully
	if !strings.Contains(code, "var Empty = workflow.Workflow{") {
		t.Error("Missing workflow declaration")
	}
}
