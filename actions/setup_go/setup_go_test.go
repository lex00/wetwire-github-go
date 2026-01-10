package setup_go

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupGo_Action(t *testing.T) {
	s := SetupGo{}
	if got := s.Action(); got != "actions/setup-go@v5" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-go@v5")
	}
}

func TestSetupGo_Inputs(t *testing.T) {
	s := SetupGo{
		GoVersion: "1.23",
		Cache:     true,
	}

	inputs := s.Inputs()

	if s.Action() != "actions/setup-go@v5" {
		t.Errorf("Action() = %q, want %q", s.Action(), "actions/setup-go@v5")
	}

	if inputs["go-version"] != "1.23" {
		t.Errorf("inputs[go-version] = %v, want %q", inputs["go-version"], "1.23")
	}

	if inputs["cache"] != true {
		t.Errorf("inputs[cache] = %v, want true", inputs["cache"])
	}
}

func TestSetupGo_Inputs_GoVersionFile(t *testing.T) {
	s := SetupGo{
		GoVersionFile: "go.mod",
	}

	inputs := s.Inputs()

	if inputs["go-version-file"] != "go.mod" {
		t.Errorf("inputs[go-version-file] = %v, want %q", inputs["go-version-file"], "go.mod")
	}
}

func TestSetupGo_ImplementsStepAction(t *testing.T) {
	s := SetupGo{}
	// Verify SetupGo implements StepAction interface
	var _ workflow.StepAction = s
}
