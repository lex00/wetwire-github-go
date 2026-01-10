package scorecard

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestScorecard_Action(t *testing.T) {
	s := Scorecard{}
	if got := s.Action(); got != "ossf/scorecard-action@v2.4.0" {
		t.Errorf("Action() = %q, want %q", got, "ossf/scorecard-action@v2.4.0")
	}
}

func TestScorecard_Inputs(t *testing.T) {
	s := Scorecard{
		ResultsFile:    "results.sarif",
		ResultsFormat:  "sarif",
		PublishResults: true,
		RepoToken:      "${{ secrets.SCORECARD_TOKEN }}",
	}

	inputs := s.Inputs()

	if inputs["results_file"] != "results.sarif" {
		t.Errorf("inputs[results_file] = %v, want %q", inputs["results_file"], "results.sarif")
	}

	if inputs["results_format"] != "sarif" {
		t.Errorf("inputs[results_format] = %v, want %q", inputs["results_format"], "sarif")
	}

	if inputs["publish_results"] != true {
		t.Errorf("inputs[publish_results] = %v, want true", inputs["publish_results"])
	}

	if inputs["repo_token"] != "${{ secrets.SCORECARD_TOKEN }}" {
		t.Errorf("inputs[repo_token] = %v, want %q", inputs["repo_token"], "${{ secrets.SCORECARD_TOKEN }}")
	}
}

func TestScorecard_Inputs_Empty(t *testing.T) {
	s := Scorecard{}
	inputs := s.Inputs()

	// Empty scorecard should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Scorecard.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestScorecard_Inputs_MinimalConfig(t *testing.T) {
	s := Scorecard{
		ResultsFile:   "results.json",
		ResultsFormat: "json",
	}

	inputs := s.Inputs()

	if inputs["results_file"] != "results.json" {
		t.Errorf("inputs[results_file] = %v, want %q", inputs["results_file"], "results.json")
	}

	if inputs["results_format"] != "json" {
		t.Errorf("inputs[results_format] = %v, want %q", inputs["results_format"], "json")
	}

	if len(inputs) != 2 {
		t.Errorf("minimal Scorecard should have 2 input entries, got %d", len(inputs))
	}
}

func TestScorecard_Inputs_BoolFields(t *testing.T) {
	s := Scorecard{
		PublishResults: true,
	}

	inputs := s.Inputs()

	if inputs["publish_results"] != true {
		t.Errorf("inputs[publish_results] = %v, want true", inputs["publish_results"])
	}
}

func TestScorecard_ImplementsStepAction(t *testing.T) {
	s := Scorecard{}
	// Verify Scorecard implements StepAction interface
	var _ workflow.StepAction = s
}
