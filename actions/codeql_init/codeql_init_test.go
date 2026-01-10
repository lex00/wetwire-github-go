package codeql_init

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCodeQLInit_Action(t *testing.T) {
	a := CodeQLInit{}
	if got := a.Action(); got != "github/codeql-action/init@v3" {
		t.Errorf("Action() = %q, want %q", got, "github/codeql-action/init@v3")
	}
}

func TestCodeQLInit_Inputs(t *testing.T) {
	a := CodeQLInit{
		Languages: "go,javascript",
		Queries:   "security-extended",
	}

	inputs := a.Inputs()

	if inputs["languages"] != "go,javascript" {
		t.Errorf("inputs[languages] = %v, want go,javascript", inputs["languages"])
	}
	if inputs["queries"] != "security-extended" {
		t.Errorf("inputs[queries] = %v, want security-extended", inputs["queries"])
	}
}

func TestCodeQLInit_Inputs_Empty(t *testing.T) {
	a := CodeQLInit{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty CodeQLInit.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCodeQLInit_Inputs_Config(t *testing.T) {
	a := CodeQLInit{
		ConfigFile: ".github/codeql/config.yml",
		BuildMode:  "autobuild",
	}

	inputs := a.Inputs()

	if inputs["config-file"] != ".github/codeql/config.yml" {
		t.Errorf("inputs[config-file] = %v, want config path", inputs["config-file"])
	}
	if inputs["build-mode"] != "autobuild" {
		t.Errorf("inputs[build-mode] = %v, want autobuild", inputs["build-mode"])
	}
}

func TestCodeQLInit_ImplementsStepAction(t *testing.T) {
	a := CodeQLInit{}
	var _ workflow.StepAction = a
}
