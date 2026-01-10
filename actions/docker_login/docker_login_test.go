package docker_login

import (
	"testing"
)

func TestDockerLogin_Action(t *testing.T) {
	d := DockerLogin{}
	if got := d.Action(); got != "docker/login-action@v3" {
		t.Errorf("Action() = %q, want %q", got, "docker/login-action@v3")
	}
}

func TestDockerLogin_ToStep(t *testing.T) {
	d := DockerLogin{
		Registry: "ghcr.io",
		Username: "${{ github.actor }}",
		Password: "${{ secrets.GITHUB_TOKEN }}",
	}

	step := d.ToStep()

	if step.Uses != "docker/login-action@v3" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "docker/login-action@v3")
	}

	if step.With["registry"] != "ghcr.io" {
		t.Errorf("step.With[registry] = %v, want %q", step.With["registry"], "ghcr.io")
	}

	if step.With["username"] != "${{ github.actor }}" {
		t.Errorf("step.With[username] = %v, want %q", step.With["username"], "${{ github.actor }}")
	}

	if step.With["password"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("step.With[password] = %v, want %q", step.With["password"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestDockerLogin_ToStep_EmptyWithMaps(t *testing.T) {
	d := DockerLogin{}
	step := d.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty DockerLogin.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestDockerLogin_ToStep_ECR(t *testing.T) {
	d := DockerLogin{
		Registry: "123456789012.dkr.ecr.us-east-1.amazonaws.com",
		ECR:      "auto",
	}

	step := d.ToStep()

	if step.With["ecr"] != "auto" {
		t.Errorf("step.With[ecr] = %v, want %q", step.With["ecr"], "auto")
	}
}

func TestDockerLogin_ToStep_Logout(t *testing.T) {
	d := DockerLogin{
		Registry: "docker.io",
		Logout:   true,
	}

	step := d.ToStep()

	if step.With["logout"] != true {
		t.Errorf("step.With[logout] = %v, want true", step.With["logout"])
	}
}
