package codecov

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCodecov_Action(t *testing.T) {
	c := Codecov{}
	if got := c.Action(); got != "codecov/codecov-action@v5" {
		t.Errorf("Action() = %q, want %q", got, "codecov/codecov-action@v5")
	}
}

func TestCodecov_Inputs(t *testing.T) {
	c := Codecov{
		Token: "${{ secrets.CODECOV_TOKEN }}",
		Files: "./coverage.xml",
		Flags: "unittests",
	}

	inputs := c.Inputs()

	if c.Action() != "codecov/codecov-action@v5" {
		t.Errorf("Action() = %q, want %q", c.Action(), "codecov/codecov-action@v5")
	}

	if inputs["token"] != "${{ secrets.CODECOV_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "${{ secrets.CODECOV_TOKEN }}")
	}

	if inputs["files"] != "./coverage.xml" {
		t.Errorf("inputs[files] = %v, want %q", inputs["files"], "./coverage.xml")
	}

	if inputs["flags"] != "unittests" {
		t.Errorf("inputs[flags] = %v, want %q", inputs["flags"], "unittests")
	}
}

func TestCodecov_Inputs_EmptyWithMaps(t *testing.T) {
	c := Codecov{}
	inputs := c.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Codecov.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCodecov_Inputs_AllFields(t *testing.T) {
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

	inputs := c.Inputs()

	if inputs["token"] != "${{ secrets.CODECOV_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want token", inputs["token"])
	}

	if inputs["files"] != "./coverage1.xml,./coverage2.xml" {
		t.Errorf("inputs[files] = %v, want files", inputs["files"])
	}

	if inputs["directory"] != "./coverage" {
		t.Errorf("inputs[directory] = %v, want directory", inputs["directory"])
	}

	if inputs["flags"] != "backend,frontend" {
		t.Errorf("inputs[flags] = %v, want flags", inputs["flags"])
	}

	if inputs["name"] != "codecov-umbrella" {
		t.Errorf("inputs[name] = %v, want name", inputs["name"])
	}

	if inputs["fail_ci_if_error"] != true {
		t.Errorf("inputs[fail_ci_if_error] = %v, want true", inputs["fail_ci_if_error"])
	}

	if inputs["verbose"] != true {
		t.Errorf("inputs[verbose] = %v, want true", inputs["verbose"])
	}

	if inputs["working-directory"] != "./src" {
		t.Errorf("inputs[working-directory] = %v, want working-directory", inputs["working-directory"])
	}

	if inputs["env_vars"] != "OS,PYTHON" {
		t.Errorf("inputs[env_vars] = %v, want env_vars", inputs["env_vars"])
	}

	if inputs["os"] != "linux" {
		t.Errorf("inputs[os] = %v, want os", inputs["os"])
	}

	if inputs["slug"] != "owner/repo" {
		t.Errorf("inputs[slug] = %v, want slug", inputs["slug"])
	}

	if inputs["use_oidc"] != true {
		t.Errorf("inputs[use_oidc] = %v, want true", inputs["use_oidc"])
	}
}

func TestCodecov_Inputs_MinimalConfig(t *testing.T) {
	c := Codecov{
		Token: "${{ secrets.CODECOV_TOKEN }}",
	}

	inputs := c.Inputs()

	if inputs["token"] != "${{ secrets.CODECOV_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want token", inputs["token"])
	}

	if len(inputs) != 1 {
		t.Errorf("minimal Codecov should have 1 input entry, got %d", len(inputs))
	}
}

func TestCodecov_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = Codecov{}
}
