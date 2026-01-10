package docker_login

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDockerLogin_Action(t *testing.T) {
	d := DockerLogin{}
	if got := d.Action(); got != "docker/login-action@v3" {
		t.Errorf("Action() = %q, want %q", got, "docker/login-action@v3")
	}
}

func TestDockerLogin_Inputs(t *testing.T) {
	d := DockerLogin{
		Registry: "ghcr.io",
		Username: "${{ github.actor }}",
		Password: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := d.Inputs()

	if d.Action() != "docker/login-action@v3" {
		t.Errorf("Action() = %q, want %q", d.Action(), "docker/login-action@v3")
	}

	if inputs["registry"] != "ghcr.io" {
		t.Errorf("inputs[registry] = %v, want %q", inputs["registry"], "ghcr.io")
	}

	if inputs["username"] != "${{ github.actor }}" {
		t.Errorf("inputs[username] = %v, want %q", inputs["username"], "${{ github.actor }}")
	}

	if inputs["password"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[password] = %v, want %q", inputs["password"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestDockerLogin_Inputs_EmptyWithMaps(t *testing.T) {
	d := DockerLogin{}
	inputs := d.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DockerLogin.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDockerLogin_Inputs_ECR(t *testing.T) {
	d := DockerLogin{
		Registry: "123456789012.dkr.ecr.us-east-1.amazonaws.com",
		ECR:      "auto",
	}

	inputs := d.Inputs()

	if inputs["ecr"] != "auto" {
		t.Errorf("inputs[ecr] = %v, want %q", inputs["ecr"], "auto")
	}
}

func TestDockerLogin_Inputs_Logout(t *testing.T) {
	d := DockerLogin{
		Registry: "docker.io",
		Logout:   true,
	}

	inputs := d.Inputs()

	if inputs["logout"] != true {
		t.Errorf("inputs[logout] = %v, want true", inputs["logout"])
	}
}

func TestDockerLogin_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DockerLogin{}
}
