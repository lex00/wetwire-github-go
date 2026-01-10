package codecov

import (
	"testing"
)

func TestCodecov_Action(t *testing.T) {
	c := Codecov{}
	if got := c.Action(); got != "codecov/codecov-action@v5" {
		t.Errorf("Action() = %q, want %q", got, "codecov/codecov-action@v5")
	}
}

func TestCodecov_ToStep(t *testing.T) {
	c := Codecov{
		Token: "${{ secrets.CODECOV_TOKEN }}",
		Files: "./coverage.xml",
		Flags: "unittests",
	}

	step := c.ToStep()

	if step.Uses != "codecov/codecov-action@v5" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "codecov/codecov-action@v5")
	}

	if step.With["token"] != "${{ secrets.CODECOV_TOKEN }}" {
		t.Errorf("step.With[token] = %v, want %q", step.With["token"], "${{ secrets.CODECOV_TOKEN }}")
	}

	if step.With["files"] != "./coverage.xml" {
		t.Errorf("step.With[files] = %v, want %q", step.With["files"], "./coverage.xml")
	}

	if step.With["flags"] != "unittests" {
		t.Errorf("step.With[flags] = %v, want %q", step.With["flags"], "unittests")
	}
}

func TestCodecov_ToStep_EmptyWithMaps(t *testing.T) {
	c := Codecov{}
	step := c.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty Codecov.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestCodecov_ToStep_AllFields(t *testing.T) {
	c := Codecov{
		Token:            "${{ secrets.CODECOV_TOKEN }}",
		Files:            "./coverage1.xml,./coverage2.xml",
		Directory:        "./coverage",
		Flags:            "backend,frontend",
		Name:             "codecov-umbrella",
		FailCIIfError:    true,
		Verbose:          true,
		WorkingDirectory: "./src",
		EnvVars:          "OS,PYTHON",
		OS:               "linux",
		Slug:             "owner/repo",
		Version:          "v0.6.0",
		DryRun:           false,
		UseOIDC:          true,
		CodecovYMLPath:   "./codecov.yml",
		Plugin:           "noop",
	}

	step := c.ToStep()

	if step.With["token"] != "${{ secrets.CODECOV_TOKEN }}" {
		t.Errorf("step.With[token] = %v, want token", step.With["token"])
	}

	if step.With["files"] != "./coverage1.xml,./coverage2.xml" {
		t.Errorf("step.With[files] = %v, want files", step.With["files"])
	}

	if step.With["directory"] != "./coverage" {
		t.Errorf("step.With[directory] = %v, want directory", step.With["directory"])
	}

	if step.With["flags"] != "backend,frontend" {
		t.Errorf("step.With[flags] = %v, want flags", step.With["flags"])
	}

	if step.With["name"] != "codecov-umbrella" {
		t.Errorf("step.With[name] = %v, want name", step.With["name"])
	}

	if step.With["fail_ci_if_error"] != true {
		t.Errorf("step.With[fail_ci_if_error] = %v, want true", step.With["fail_ci_if_error"])
	}

	if step.With["verbose"] != true {
		t.Errorf("step.With[verbose] = %v, want true", step.With["verbose"])
	}

	if step.With["working-directory"] != "./src" {
		t.Errorf("step.With[working-directory] = %v, want working-directory", step.With["working-directory"])
	}

	if step.With["env_vars"] != "OS,PYTHON" {
		t.Errorf("step.With[env_vars] = %v, want env_vars", step.With["env_vars"])
	}

	if step.With["os"] != "linux" {
		t.Errorf("step.With[os] = %v, want os", step.With["os"])
	}

	if step.With["slug"] != "owner/repo" {
		t.Errorf("step.With[slug] = %v, want slug", step.With["slug"])
	}

	if step.With["use_oidc"] != true {
		t.Errorf("step.With[use_oidc] = %v, want true", step.With["use_oidc"])
	}
}

func TestCodecov_ToStep_MinimalConfig(t *testing.T) {
	c := Codecov{
		Token: "${{ secrets.CODECOV_TOKEN }}",
	}

	step := c.ToStep()

	if step.With["token"] != "${{ secrets.CODECOV_TOKEN }}" {
		t.Errorf("step.With[token] = %v, want token", step.With["token"])
	}

	if len(step.With) != 1 {
		t.Errorf("minimal Codecov should have 1 with entry, got %d", len(step.With))
	}
}
