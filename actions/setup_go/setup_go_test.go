package setup_go

import (
	"testing"
)

func TestSetupGo_Action(t *testing.T) {
	s := SetupGo{}
	if got := s.Action(); got != "actions/setup-go@v5" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-go@v5")
	}
}

func TestSetupGo_ToStep(t *testing.T) {
	s := SetupGo{
		GoVersion: "1.23",
		Cache:     true,
	}

	step := s.ToStep()

	if step.Uses != "actions/setup-go@v5" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/setup-go@v5")
	}

	if step.With["go-version"] != "1.23" {
		t.Errorf("step.With[go-version] = %v, want %q", step.With["go-version"], "1.23")
	}

	if step.With["cache"] != true {
		t.Errorf("step.With[cache] = %v, want true", step.With["cache"])
	}
}

func TestSetupGo_ToStep_GoVersionFile(t *testing.T) {
	s := SetupGo{
		GoVersionFile: "go.mod",
	}

	step := s.ToStep()

	if step.With["go-version-file"] != "go.mod" {
		t.Errorf("step.With[go-version-file] = %v, want %q", step.With["go-version-file"], "go.mod")
	}
}
