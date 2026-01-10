package setup_ruby

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSetupRuby_Action(t *testing.T) {
	a := SetupRuby{}
	if got := a.Action(); got != "ruby/setup-ruby@v1" {
		t.Errorf("Action() = %q, want %q", got, "ruby/setup-ruby@v1")
	}
}

func TestSetupRuby_Inputs(t *testing.T) {
	a := SetupRuby{
		RubyVersion: "3.3",
	}

	inputs := a.Inputs()

	if a.Action() != "ruby/setup-ruby@v1" {
		t.Errorf("Action() = %q, want %q", a.Action(), "ruby/setup-ruby@v1")
	}

	if inputs["ruby-version"] != "3.3" {
		t.Errorf("inputs[ruby-version] = %v, want %q", inputs["ruby-version"], "3.3")
	}
}

func TestSetupRuby_Inputs_Empty(t *testing.T) {
	a := SetupRuby{}
	inputs := a.Inputs()

	if a.Action() != "ruby/setup-ruby@v1" {
		t.Errorf("Action() = %q, want %q", a.Action(), "ruby/setup-ruby@v1")
	}

	if _, ok := inputs["ruby-version"]; ok {
		t.Error("Empty ruby-version should not be in inputs")
	}
}

func TestSetupRuby_Inputs_WithBundler(t *testing.T) {
	a := SetupRuby{
		RubyVersion:    "3.3",
		BundlerCache:   true,
		Bundler:        "true",
		BundlerVersion: "latest",
	}

	inputs := a.Inputs()

	if inputs["ruby-version"] != "3.3" {
		t.Errorf("ruby-version = %v, want %q", inputs["ruby-version"], "3.3")
	}
	if inputs["bundler-cache"] != true {
		t.Errorf("bundler-cache = %v, want %v", inputs["bundler-cache"], true)
	}
	if inputs["bundler"] != "true" {
		t.Errorf("bundler = %v, want %q", inputs["bundler"], "true")
	}
	if inputs["bundler-version"] != "latest" {
		t.Errorf("bundler-version = %v, want %q", inputs["bundler-version"], "latest")
	}
}

func TestSetupRuby_Inputs_Versions(t *testing.T) {
	versions := []string{"3.2", "3.3", "ruby-head", "jruby-9.4", "truffleruby"}

	for _, ver := range versions {
		t.Run(ver, func(t *testing.T) {
			a := SetupRuby{
				RubyVersion: ver,
			}

			inputs := a.Inputs()
			if inputs["ruby-version"] != ver {
				t.Errorf("ruby-version = %v, want %q", inputs["ruby-version"], ver)
			}
		})
	}
}

func TestSetupRuby_Inputs_WorkingDirectory(t *testing.T) {
	a := SetupRuby{
		RubyVersion:      "3.3",
		WorkingDirectory: "./app",
	}

	inputs := a.Inputs()

	if inputs["working-directory"] != "./app" {
		t.Errorf("working-directory = %v, want %q", inputs["working-directory"], "./app")
	}
}

func TestSetupRuby_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = SetupRuby{}
}
