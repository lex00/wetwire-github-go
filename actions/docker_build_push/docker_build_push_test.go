package docker_build_push

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDockerBuildPush_Action(t *testing.T) {
	d := DockerBuildPush{}
	if got := d.Action(); got != "docker/build-push-action@v6" {
		t.Errorf("Action() = %q, want %q", got, "docker/build-push-action@v6")
	}
}

func TestDockerBuildPush_Inputs(t *testing.T) {
	d := DockerBuildPush{
		Context: ".",
		File:    "./Dockerfile",
		Push:    true,
		Tags:    "user/app:latest",
	}

	inputs := d.Inputs()

	if d.Action() != "docker/build-push-action@v6" {
		t.Errorf("Action() = %q, want %q", d.Action(), "docker/build-push-action@v6")
	}

	if inputs["context"] != "." {
		t.Errorf("inputs[context] = %v, want %q", inputs["context"], ".")
	}

	if inputs["file"] != "./Dockerfile" {
		t.Errorf("inputs[file] = %v, want %q", inputs["file"], "./Dockerfile")
	}

	if inputs["push"] != true {
		t.Errorf("inputs[push] = %v, want true", inputs["push"])
	}

	if inputs["tags"] != "user/app:latest" {
		t.Errorf("inputs[tags] = %v, want %q", inputs["tags"], "user/app:latest")
	}
}

func TestDockerBuildPush_Inputs_EmptyWithMaps(t *testing.T) {
	d := DockerBuildPush{}
	inputs := d.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DockerBuildPush.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDockerBuildPush_Inputs_AllFields(t *testing.T) {
	d := DockerBuildPush{
		Context:    ".",
		File:       "Dockerfile.prod",
		Push:       true,
		Load:       false,
		Tags:       "ghcr.io/user/app:latest\nghcr.io/user/app:v1.0.0",
		BuildArgs:  "VERSION=1.0.0\nDEBUG=false",
		Platforms:  "linux/amd64,linux/arm64",
		CacheFrom:  "type=gha",
		CacheTo:    "type=gha,mode=max",
		Target:     "production",
		NoCache:    false,
		Pull:       true,
		Secrets:    "mysecret",
		Labels:     "org.opencontainers.image.source=https://github.com/user/repo",
		Outputs:    "type=registry",
		Provenance: "false",
		SBOM:       "false",
	}

	inputs := d.Inputs()

	if inputs["context"] != "." {
		t.Errorf("inputs[context] = %v, want %q", inputs["context"], ".")
	}

	if inputs["platforms"] != "linux/amd64,linux/arm64" {
		t.Errorf("inputs[platforms] = %v, want %q", inputs["platforms"], "linux/amd64,linux/arm64")
	}

	if inputs["cache-from"] != "type=gha" {
		t.Errorf("inputs[cache-from] = %v, want %q", inputs["cache-from"], "type=gha")
	}

	if inputs["cache-to"] != "type=gha,mode=max" {
		t.Errorf("inputs[cache-to] = %v, want %q", inputs["cache-to"], "type=gha,mode=max")
	}

	if inputs["target"] != "production" {
		t.Errorf("inputs[target] = %v, want %q", inputs["target"], "production")
	}

	if inputs["pull"] != true {
		t.Errorf("inputs[pull] = %v, want true", inputs["pull"])
	}

	if inputs["provenance"] != "false" {
		t.Errorf("inputs[provenance] = %v, want %q", inputs["provenance"], "false")
	}

	if inputs["sbom"] != "false" {
		t.Errorf("inputs[sbom] = %v, want %q", inputs["sbom"], "false")
	}
}

func TestDockerBuildPush_Inputs_MultiPlatform(t *testing.T) {
	d := DockerBuildPush{
		Context:   ".",
		Push:      true,
		Tags:      "user/app:latest",
		Platforms: "linux/amd64,linux/arm64,linux/arm/v7",
	}

	inputs := d.Inputs()

	if inputs["platforms"] != "linux/amd64,linux/arm64,linux/arm/v7" {
		t.Errorf("inputs[platforms] = %v, want multi-platform string", inputs["platforms"])
	}
}

func TestDockerBuildPush_Inputs_LoadAndNoCache(t *testing.T) {
	d := DockerBuildPush{
		Context: ".",
		Load:    true,
		NoCache: true,
	}

	inputs := d.Inputs()

	if inputs["load"] != true {
		t.Errorf("inputs[load] = %v, want true", inputs["load"])
	}
	if inputs["no-cache"] != true {
		t.Errorf("inputs[no-cache] = %v, want true", inputs["no-cache"])
	}
}

func TestDockerBuildPush_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DockerBuildPush{}
}
