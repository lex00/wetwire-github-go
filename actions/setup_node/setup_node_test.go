package setup_node

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupNode_Action(t *testing.T) {
	a := SetupNode{}
	if got := a.Action(); got != "actions/setup-node@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-node@v4")
	}
}

func TestSetupNode_Inputs(t *testing.T) {
	a := SetupNode{
		NodeVersion: "20",
	}

	inputs := a.Inputs()

	if a.Action() != "actions/setup-node@v4" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/setup-node@v4")
	}

	if inputs["node-version"] != "20" {
		t.Errorf("inputs[node-version] = %v, want %q", inputs["node-version"], "20")
	}
}

func TestSetupNode_Inputs_Empty(t *testing.T) {
	a := SetupNode{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty SetupNode.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSetupNode_Inputs_AllFields(t *testing.T) {
	a := SetupNode{
		NodeVersion:         "18.x",
		NodeVersionFile:     ".nvmrc",
		Architecture:        "x64",
		CheckLatest:         true,
		RegistryURL:         "https://npm.pkg.github.com",
		Scope:               "@myorg",
		Token:               "token",
		Cache:               "npm",
		CacheDependencyPath: "package-lock.json",
		AlwaysAuth:          true,
	}

	inputs := a.Inputs()

	if inputs["node-version"] != "18.x" {
		t.Errorf("inputs[node-version] = %v, want %q", inputs["node-version"], "18.x")
	}

	if inputs["node-version-file"] != ".nvmrc" {
		t.Errorf("inputs[node-version-file] = %v, want %q", inputs["node-version-file"], ".nvmrc")
	}

	if inputs["architecture"] != "x64" {
		t.Errorf("inputs[architecture] = %v, want %q", inputs["architecture"], "x64")
	}

	if inputs["check-latest"] != true {
		t.Errorf("inputs[check-latest] = %v, want true", inputs["check-latest"])
	}

	if inputs["registry-url"] != "https://npm.pkg.github.com" {
		t.Errorf("inputs[registry-url] = %v, want %q", inputs["registry-url"], "https://npm.pkg.github.com")
	}

	if inputs["scope"] != "@myorg" {
		t.Errorf("inputs[scope] = %v, want %q", inputs["scope"], "@myorg")
	}

	if inputs["token"] != "token" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "token")
	}

	if inputs["cache"] != "npm" {
		t.Errorf("inputs[cache] = %v, want %q", inputs["cache"], "npm")
	}

	if inputs["cache-dependency-path"] != "package-lock.json" {
		t.Errorf("inputs[cache-dependency-path] = %v, want %q", inputs["cache-dependency-path"], "package-lock.json")
	}

	if inputs["always-auth"] != true {
		t.Errorf("inputs[always-auth] = %v, want true", inputs["always-auth"])
	}
}

func TestSetupNode_ImplementsStepAction(t *testing.T) {
	a := SetupNode{}
	// Verify SetupNode implements StepAction interface
	var _ workflow.StepAction = a
}
