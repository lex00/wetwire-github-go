package sonarcloud

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSonarCloud_Action(t *testing.T) {
	a := SonarCloud{}
	if got := a.Action(); got != "SonarSource/sonarcloud-github-action@v3" {
		t.Errorf("Action() = %q, want %q", got, "SonarSource/sonarcloud-github-action@v3")
	}
}

func TestSonarCloud_Inputs_Empty(t *testing.T) {
	a := SonarCloud{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty SonarCloud.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSonarCloud_ImplementsStepAction(t *testing.T) {
	a := SonarCloud{}
	var _ workflow.StepAction = a
}

func TestSonarCloud_Inputs_ProjectBaseDir(t *testing.T) {
	a := SonarCloud{
		ProjectBaseDir: "./src",
	}

	inputs := a.Inputs()

	if inputs["projectBaseDir"] != "./src" {
		t.Errorf("inputs[projectBaseDir] = %v, want %q", inputs["projectBaseDir"], "./src")
	}
}

func TestSonarCloud_Inputs_Args(t *testing.T) {
	a := SonarCloud{
		Args: "-Dsonar.verbose=true",
	}

	inputs := a.Inputs()

	if inputs["args"] != "-Dsonar.verbose=true" {
		t.Errorf("inputs[args] = %v, want %q", inputs["args"], "-Dsonar.verbose=true")
	}
}

func TestSonarCloud_Inputs_Args_Multiple(t *testing.T) {
	a := SonarCloud{
		Args: "-Dsonar.verbose=true -Dsonar.javascript.lcov.reportPaths=coverage/lcov.info",
	}

	inputs := a.Inputs()

	expected := "-Dsonar.verbose=true -Dsonar.javascript.lcov.reportPaths=coverage/lcov.info"
	if inputs["args"] != expected {
		t.Errorf("inputs[args] = %v, want %q", inputs["args"], expected)
	}
}

func TestSonarCloud_Inputs_AllFields(t *testing.T) {
	a := SonarCloud{
		ProjectBaseDir: "./app",
		Args:           "-Dsonar.sources=src -Dsonar.tests=test",
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"projectBaseDir": "./app",
		"args":           "-Dsonar.sources=src -Dsonar.tests=test",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestSonarCloud_Inputs_CommonUsage(t *testing.T) {
	// Test common usage pattern: just run with defaults
	a := SonarCloud{}

	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("common usage (defaults) has %d entries, want 0", len(inputs))
	}
}

func TestSonarCloud_Inputs_MonorepoUsage(t *testing.T) {
	// Test monorepo usage with specific project directory
	a := SonarCloud{
		ProjectBaseDir: "./packages/backend",
	}

	inputs := a.Inputs()

	if len(inputs) != 1 {
		t.Errorf("monorepo usage has %d entries, want 1", len(inputs))
	}

	if inputs["projectBaseDir"] != "./packages/backend" {
		t.Errorf("inputs[projectBaseDir] = %v, want ./packages/backend", inputs["projectBaseDir"])
	}
}

func TestSonarCloud_Inputs_CoverageReportUsage(t *testing.T) {
	// Test coverage report configuration
	a := SonarCloud{
		Args: "-Dsonar.javascript.lcov.reportPaths=coverage/lcov.info -Dsonar.coverage.exclusions=**/*.test.js",
	}

	inputs := a.Inputs()

	if len(inputs) != 1 {
		t.Errorf("coverage report usage has %d entries, want 1", len(inputs))
	}
}

func TestSonarCloud_Inputs_EmptyProjectBaseDir(t *testing.T) {
	a := SonarCloud{
		ProjectBaseDir: "",
	}

	inputs := a.Inputs()

	if _, exists := inputs["projectBaseDir"]; exists {
		t.Error("inputs[projectBaseDir] should not exist for empty ProjectBaseDir")
	}
}

func TestSonarCloud_Inputs_EmptyArgs(t *testing.T) {
	a := SonarCloud{
		Args: "",
	}

	inputs := a.Inputs()

	if _, exists := inputs["args"]; exists {
		t.Error("inputs[args] should not exist for empty Args")
	}
}
