package docker_setup_buildx

import (
	"testing"
)

func TestDockerSetupBuildx_Action(t *testing.T) {
	d := DockerSetupBuildx{}
	if got := d.Action(); got != "docker/setup-buildx-action@v3" {
		t.Errorf("Action() = %q, want %q", got, "docker/setup-buildx-action@v3")
	}
}

func TestDockerSetupBuildx_ToStep(t *testing.T) {
	d := DockerSetupBuildx{
		Version: "latest",
		Driver:  "docker-container",
	}

	step := d.ToStep()

	if step.Uses != "docker/setup-buildx-action@v3" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "docker/setup-buildx-action@v3")
	}

	if step.With["version"] != "latest" {
		t.Errorf("step.With[version] = %v, want %q", step.With["version"], "latest")
	}

	if step.With["driver"] != "docker-container" {
		t.Errorf("step.With[driver] = %v, want %q", step.With["driver"], "docker-container")
	}
}

func TestDockerSetupBuildx_ToStep_EmptyWithMaps(t *testing.T) {
	d := DockerSetupBuildx{}
	step := d.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty DockerSetupBuildx.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestDockerSetupBuildx_ToStep_AllFields(t *testing.T) {
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

	step := d.ToStep()

	if step.With["version"] != "v0.12.0" {
		t.Errorf("step.With[version] = %v, want %q", step.With["version"], "v0.12.0")
	}

	if step.With["driver"] != "docker-container" {
		t.Errorf("step.With[driver] = %v, want %q", step.With["driver"], "docker-container")
	}

	if step.With["driver-opts"] != "network=host" {
		t.Errorf("step.With[driver-opts] = %v, want %q", step.With["driver-opts"], "network=host")
	}

	if step.With["buildkitd-flags"] != "--allow-insecure-entitlement security.insecure" {
		t.Errorf("step.With[buildkitd-flags] = %v, want buildkitd flags string", step.With["buildkitd-flags"])
	}

	if step.With["install"] != true {
		t.Errorf("step.With[install] = %v, want true", step.With["install"])
	}

	if step.With["use"] != true {
		t.Errorf("step.With[use] = %v, want true", step.With["use"])
	}

	if step.With["platforms"] != "linux/amd64,linux/arm64" {
		t.Errorf("step.With[platforms] = %v, want platforms string", step.With["platforms"])
	}

	if step.With["cleanup"] != true {
		t.Errorf("step.With[cleanup] = %v, want true", step.With["cleanup"])
	}
}

func TestDockerSetupBuildx_ToStep_ConfigInline(t *testing.T) {
	d := DockerSetupBuildx{
		ConfigInline: "[worker.oci]\nmax-parallelism = 4",
	}

	step := d.ToStep()

	if step.With["config-inline"] != "[worker.oci]\nmax-parallelism = 4" {
		t.Errorf("step.With[config-inline] = %v, want config string", step.With["config-inline"])
	}
}
