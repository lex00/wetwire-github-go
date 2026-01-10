package setup_python

import (
	"testing"
)

func TestSetupPython_Action(t *testing.T) {
	a := SetupPython{}
	if got := a.Action(); got != "actions/setup-python@v5" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-python@v5")
	}
}

func TestSetupPython_ToStep(t *testing.T) {
	a := SetupPython{
		PythonVersion: "3.12",
	}

	step := a.ToStep()

	if step.Uses != "actions/setup-python@v5" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/setup-python@v5")
	}

	if step.With["python-version"] != "3.12" {
		t.Errorf("step.With[python-version] = %v, want %q", step.With["python-version"], "3.12")
	}
}

func TestSetupPython_ToStep_Empty(t *testing.T) {
	a := SetupPython{}
	step := a.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty SetupPython.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestSetupPython_ToStep_AllFields(t *testing.T) {
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

	step := a.ToStep()

	if step.With["python-version"] != "3.11" {
		t.Errorf("step.With[python-version] = %v, want %q", step.With["python-version"], "3.11")
	}

	if step.With["python-version-file"] != ".python-version" {
		t.Errorf("step.With[python-version-file] = %v, want %q", step.With["python-version-file"], ".python-version")
	}

	if step.With["cache"] != "pip" {
		t.Errorf("step.With[cache] = %v, want %q", step.With["cache"], "pip")
	}

	if step.With["architecture"] != "x64" {
		t.Errorf("step.With[architecture] = %v, want %q", step.With["architecture"], "x64")
	}

	if step.With["check-latest"] != true {
		t.Errorf("step.With[check-latest] = %v, want true", step.With["check-latest"])
	}

	if step.With["token"] != "token" {
		t.Errorf("step.With[token] = %v, want %q", step.With["token"], "token")
	}

	if step.With["cache-dependency-path"] != "requirements.txt" {
		t.Errorf("step.With[cache-dependency-path] = %v, want %q", step.With["cache-dependency-path"], "requirements.txt")
	}

	if step.With["update-environment"] != true {
		t.Errorf("step.With[update-environment] = %v, want true", step.With["update-environment"])
	}

	if step.With["allow-prereleases"] != true {
		t.Errorf("step.With[allow-prereleases] = %v, want true", step.With["allow-prereleases"])
	}
}
