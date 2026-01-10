package setup_python

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupPython_Action(t *testing.T) {
	a := SetupPython{}
	if got := a.Action(); got != "actions/setup-python@v5" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-python@v5")
	}
}

func TestSetupPython_Inputs(t *testing.T) {
	a := SetupPython{
		PythonVersion: "3.12",
	}

	inputs := a.Inputs()

	if a.Action() != "actions/setup-python@v5" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/setup-python@v5")
	}

	if inputs["python-version"] != "3.12" {
		t.Errorf("inputs[python-version] = %v, want %q", inputs["python-version"], "3.12")
	}
}

func TestSetupPython_Inputs_Empty(t *testing.T) {
	a := SetupPython{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty SetupPython.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSetupPython_Inputs_AllFields(t *testing.T) {
	a := SetupPython{
		PythonVersion:       "3.11",
		PythonVersionFile:   ".python-version",
		Cache:               "pip",
		Architecture:        "x64",
		CheckLatest:         true,
		Token:               "token",
		CacheDependencyPath: "requirements.txt",
		UpdateEnvironment:   true,
		AllowPrereleases:    true,
	}

	inputs := a.Inputs()

	if inputs["python-version"] != "3.11" {
		t.Errorf("inputs[python-version] = %v, want %q", inputs["python-version"], "3.11")
	}

	if inputs["python-version-file"] != ".python-version" {
		t.Errorf("inputs[python-version-file] = %v, want %q", inputs["python-version-file"], ".python-version")
	}

	if inputs["cache"] != "pip" {
		t.Errorf("inputs[cache] = %v, want %q", inputs["cache"], "pip")
	}

	if inputs["architecture"] != "x64" {
		t.Errorf("inputs[architecture] = %v, want %q", inputs["architecture"], "x64")
	}

	if inputs["check-latest"] != true {
		t.Errorf("inputs[check-latest] = %v, want true", inputs["check-latest"])
	}

	if inputs["token"] != "token" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "token")
	}

	if inputs["cache-dependency-path"] != "requirements.txt" {
		t.Errorf("inputs[cache-dependency-path] = %v, want %q", inputs["cache-dependency-path"], "requirements.txt")
	}

	if inputs["update-environment"] != true {
		t.Errorf("inputs[update-environment] = %v, want true", inputs["update-environment"])
	}

	if inputs["allow-prereleases"] != true {
		t.Errorf("inputs[allow-prereleases] = %v, want true", inputs["allow-prereleases"])
	}
}

func TestSetupPython_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = SetupPython{}
}
