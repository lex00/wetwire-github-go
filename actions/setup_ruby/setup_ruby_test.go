package setup_ruby

import (
	"testing"
)

func TestSetupRuby_Action(t *testing.T) {
	a := SetupRuby{}
	if got := a.Action(); got != "ruby/setup-ruby@v1" {
		t.Errorf("Action() = %q, want %q", got, "ruby/setup-ruby@v1")
	}
}

func TestSetupRuby_ToStep(t *testing.T) {
	a := SetupRuby{
		RubyVersion: "3.3",
	}

	step := a.ToStep()

	if step.Uses != "ruby/setup-ruby@v1" {
		t.Errorf("Uses = %q, want %q", step.Uses, "ruby/setup-ruby@v1")
	}

	if step.With["ruby-version"] != "3.3" {
		t.Errorf("With[ruby-version] = %v, want %q", step.With["ruby-version"], "3.3")
	}
}

func TestSetupRuby_ToStep_Empty(t *testing.T) {
	a := SetupRuby{}
	step := a.ToStep()

	if step.Uses != "ruby/setup-ruby@v1" {
		t.Errorf("Uses = %q, want %q", step.Uses, "ruby/setup-ruby@v1")
	}

	if _, ok := step.With["ruby-version"]; ok {
		t.Error("Empty ruby-version should not be in With")
	}
}

func TestSetupRuby_ToStep_WithBundler(t *testing.T) {
	a := SetupRuby{
		RubyVersion:   "3.3",
		BundlerCache:  true,
		Bundler:       "true",
		BundlerVersion: "latest",
	}

	step := a.ToStep()

	if step.With["ruby-version"] != "3.3" {
		t.Errorf("ruby-version = %v, want %q", step.With["ruby-version"], "3.3")
	}
	if step.With["bundler-cache"] != true {
		t.Errorf("bundler-cache = %v, want %v", step.With["bundler-cache"], true)
	}
	if step.With["bundler"] != "true" {
		t.Errorf("bundler = %v, want %q", step.With["bundler"], "true")
	}
	if step.With["bundler-version"] != "latest" {
		t.Errorf("bundler-version = %v, want %q", step.With["bundler-version"], "latest")
	}
}

func TestSetupRuby_ToStep_Versions(t *testing.T) {
	versions := []string{"3.2", "3.3", "ruby-head", "jruby-9.4", "truffleruby"}

	for _, ver := range versions {
		t.Run(ver, func(t *testing.T) {
			a := SetupRuby{
				RubyVersion: ver,
			}

			step := a.ToStep()
			if step.With["ruby-version"] != ver {
				t.Errorf("ruby-version = %v, want %q", step.With["ruby-version"], ver)
			}
		})
	}
}

func TestSetupRuby_ToStep_WorkingDirectory(t *testing.T) {
	a := SetupRuby{
		RubyVersion:      "3.3",
		WorkingDirectory: "./app",
	}

	step := a.ToStep()

	if step.With["working-directory"] != "./app" {
		t.Errorf("working-directory = %v, want %q", step.With["working-directory"], "./app")
	}
}
