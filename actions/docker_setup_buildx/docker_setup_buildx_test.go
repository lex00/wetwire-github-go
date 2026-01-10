package docker_setup_buildx

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDockerSetupBuildx_Action(t *testing.T) {
	d := DockerSetupBuildx{}
	if got := d.Action(); got != "docker/setup-buildx-action@v3" {
		t.Errorf("Action() = %q, want %q", got, "docker/setup-buildx-action@v3")
	}
}

func TestDockerSetupBuildx_Inputs(t *testing.T) {
	d := DockerSetupBuildx{
		Version: "latest",
		Driver:  "docker-container",
	}

	inputs := d.Inputs()

	if d.Action() != "docker/setup-buildx-action@v3" {
		t.Errorf("Action() = %q, want %q", d.Action(), "docker/setup-buildx-action@v3")
	}

	if inputs["version"] != "latest" {
		t.Errorf("inputs[version] = %v, want %q", inputs["version"], "latest")
	}

	if inputs["driver"] != "docker-container" {
		t.Errorf("inputs[driver] = %v, want %q", inputs["driver"], "docker-container")
	}
}

func TestDockerSetupBuildx_Inputs_EmptyWithMaps(t *testing.T) {
	d := DockerSetupBuildx{}
	inputs := d.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DockerSetupBuildx.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDockerSetupBuildx_Inputs_AllFields(t *testing.T) {
	d := DockerSetupBuildx{
		Version:        "v0.12.0",
		Driver:         "docker-container",
		DriverOpts:     "network=host",
		BuildkitdFlags: "--allow-insecure-entitlement security.insecure",
		Install:        true,
		Use:            true,
		Endpoint:       "tcp://localhost:2375",
		Platforms:      "linux/amd64,linux/arm64",
		Config:         "/path/to/buildkitd.toml",
		ConfigInline:   "[worker.oci]\nmax-parallelism = 4",
		Append:         "name=builder2,endpoint=tcp://remote:1234",
		Cleanup:        true,
	}

	inputs := d.Inputs()

	if inputs["version"] != "v0.12.0" {
		t.Errorf("inputs[version] = %v, want %q", inputs["version"], "v0.12.0")
	}

	if inputs["driver"] != "docker-container" {
		t.Errorf("inputs[driver] = %v, want %q", inputs["driver"], "docker-container")
	}

	if inputs["driver-opts"] != "network=host" {
		t.Errorf("inputs[driver-opts] = %v, want %q", inputs["driver-opts"], "network=host")
	}

	if inputs["buildkitd-flags"] != "--allow-insecure-entitlement security.insecure" {
		t.Errorf("inputs[buildkitd-flags] = %v, want buildkitd flags string", inputs["buildkitd-flags"])
	}

	if inputs["install"] != true {
		t.Errorf("inputs[install] = %v, want true", inputs["install"])
	}

	if inputs["use"] != true {
		t.Errorf("inputs[use] = %v, want true", inputs["use"])
	}

	if inputs["platforms"] != "linux/amd64,linux/arm64" {
		t.Errorf("inputs[platforms] = %v, want platforms string", inputs["platforms"])
	}

	if inputs["cleanup"] != true {
		t.Errorf("inputs[cleanup] = %v, want true", inputs["cleanup"])
	}
}

func TestDockerSetupBuildx_Inputs_ConfigInline(t *testing.T) {
	d := DockerSetupBuildx{
		ConfigInline: "[worker.oci]\nmax-parallelism = 4",
	}

	inputs := d.Inputs()

	if inputs["config-inline"] != "[worker.oci]\nmax-parallelism = 4" {
		t.Errorf("inputs[config-inline] = %v, want config string", inputs["config-inline"])
	}
}

func TestDockerSetupBuildx_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DockerSetupBuildx{}
}
