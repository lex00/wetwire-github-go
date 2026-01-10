package setup_node

import (
	"testing"
)

func TestSetupNode_Action(t *testing.T) {
	a := SetupNode{}
	if got := a.Action(); got != "actions/setup-node@v4" {
		t.Errorf("Action() = %q, want %q", got, "actions/setup-node@v4")
	}
}

func TestSetupNode_ToStep(t *testing.T) {
	a := SetupNode{
		NodeVersion: "20",
	}

	step := a.ToStep()

	if step.Uses != "actions/setup-node@v4" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "actions/setup-node@v4")
	}

	if step.With["node-version"] != "20" {
		t.Errorf("step.With[node-version] = %v, want %q", step.With["node-version"], "20")
	}
}

func TestSetupNode_ToStep_Empty(t *testing.T) {
	a := SetupNode{}
	step := a.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty SetupNode.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestSetupNode_ToStep_AllFields(t *testing.T) {
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

	step := a.ToStep()

	if step.With["node-version"] != "18.x" {
		t.Errorf("step.With[node-version] = %v, want %q", step.With["node-version"], "18.x")
	}

	if step.With["node-version-file"] != ".nvmrc" {
		t.Errorf("step.With[node-version-file] = %v, want %q", step.With["node-version-file"], ".nvmrc")
	}

	if step.With["architecture"] != "x64" {
		t.Errorf("step.With[architecture] = %v, want %q", step.With["architecture"], "x64")
	}

	if step.With["check-latest"] != true {
		t.Errorf("step.With[check-latest] = %v, want true", step.With["check-latest"])
	}

	if step.With["registry-url"] != "https://npm.pkg.github.com" {
		t.Errorf("step.With[registry-url] = %v, want %q", step.With["registry-url"], "https://npm.pkg.github.com")
	}

	if step.With["scope"] != "@myorg" {
		t.Errorf("step.With[scope] = %v, want %q", step.With["scope"], "@myorg")
	}

	if step.With["token"] != "token" {
		t.Errorf("step.With[token] = %v, want %q", step.With["token"], "token")
	}

	if step.With["cache"] != "npm" {
		t.Errorf("step.With[cache] = %v, want %q", step.With["cache"], "npm")
	}

	if step.With["cache-dependency-path"] != "package-lock.json" {
		t.Errorf("step.With[cache-dependency-path] = %v, want %q", step.With["cache-dependency-path"], "package-lock.json")
	}

	if step.With["always-auth"] != true {
		t.Errorf("step.With[always-auth] = %v, want true", step.With["always-auth"])
	}
}
