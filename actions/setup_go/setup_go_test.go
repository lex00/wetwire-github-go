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

func TestSetupGo_Inputs_Empty(t *testing.T) {
	s := SetupGo{}
	inputs := s.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty SetupGo.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSetupGo_Inputs_CheckLatest(t *testing.T) {
	s := SetupGo{
		GoVersion:   "1.23",
		CheckLatest: true,
	}

	inputs := s.Inputs()

	if inputs["check-latest"] != true {
		t.Errorf("inputs[check-latest] = %v, want true", inputs["check-latest"])
	}
}

func TestSetupGo_Inputs_Token(t *testing.T) {
	s := SetupGo{
		GoVersion: "1.23",
		Token:     "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := s.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want secret reference", inputs["token"])
	}
}

func TestSetupGo_Inputs_CacheDependencyPath(t *testing.T) {
	s := SetupGo{
		GoVersion:           "1.23",
		Cache:               true,
		CacheDependencyPath: "go.sum",
	}

	inputs := s.Inputs()

	if inputs["cache-dependency-path"] != "go.sum" {
		t.Errorf("inputs[cache-dependency-path] = %v, want go.sum", inputs["cache-dependency-path"])
	}
}

func TestSetupGo_Inputs_Architecture(t *testing.T) {
	s := SetupGo{
		GoVersion:    "1.23",
		Architecture: "x64",
	}

	inputs := s.Inputs()

	if inputs["architecture"] != "x64" {
		t.Errorf("inputs[architecture] = %v, want x64", inputs["architecture"])
	}
}

func TestSetupGo_ImplementsStepAction(t *testing.T) {
	s := SetupGo{}
	// Verify SetupGo implements StepAction interface
	var _ workflow.StepAction = s
}
