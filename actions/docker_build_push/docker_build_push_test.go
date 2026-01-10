package docker_build_push

import (
	"testing"
)

func TestDockerBuildPush_Action(t *testing.T) {
	d := DockerBuildPush{}
	if got := d.Action(); got != "docker/build-push-action@v6" {
		t.Errorf("Action() = %q, want %q", got, "docker/build-push-action@v6")
	}
}

func TestDockerBuildPush_ToStep(t *testing.T) {
	d := DockerBuildPush{
		Context:    ".",
		File:       "./Dockerfile",
		Push:       true,
		Tags:       "user/app:latest",
	}

	step := d.ToStep()

	if step.Uses != "docker/build-push-action@v6" {
		t.Errorf("step.Uses = %q, want %q", step.Uses, "docker/build-push-action@v6")
	}

	if step.With["context"] != "." {
		t.Errorf("step.With[context] = %v, want %q", step.With["context"], ".")
	}

	if step.With["file"] != "./Dockerfile" {
		t.Errorf("step.With[file] = %v, want %q", step.With["file"], "./Dockerfile")
	}

	if step.With["push"] != true {
		t.Errorf("step.With[push] = %v, want true", step.With["push"])
	}

	if step.With["tags"] != "user/app:latest" {
		t.Errorf("step.With[tags] = %v, want %q", step.With["tags"], "user/app:latest")
	}
}

func TestDockerBuildPush_ToStep_EmptyWithMaps(t *testing.T) {
	d := DockerBuildPush{}
	step := d.ToStep()

	if len(step.With) != 0 {
		t.Errorf("empty DockerBuildPush.ToStep() has %d with entries, want 0", len(step.With))
	}
}

func TestDockerBuildPush_ToStep_AllFields(t *testing.T) {
	d := DockerBuildPush{
		Context:     ".",
		File:        "Dockerfile.prod",
		Push:        true,
		Load:        false,
		Tags:        "ghcr.io/user/app:latest\nghcr.io/user/app:v1.0.0",
		BuildArgs:   "VERSION=1.0.0\nDEBUG=false",
		Platforms:   "linux/amd64,linux/arm64",
		CacheFrom:   "type=gha",
		CacheTo:     "type=gha,mode=max",
		Target:      "production",
		NoCache:     false,
		Pull:        true,
		Secrets:     "mysecret",
		Labels:      "org.opencontainers.image.source=https://github.com/user/repo",
		Outputs:     "type=registry",
		Provenance:  "false",
		SBOM:        "false",
	}

	step := d.ToStep()

	if step.With["context"] != "." {
		t.Errorf("step.With[context] = %v, want %q", step.With["context"], ".")
	}

	if step.With["platforms"] != "linux/amd64,linux/arm64" {
		t.Errorf("step.With[platforms] = %v, want %q", step.With["platforms"], "linux/amd64,linux/arm64")
	}

	if step.With["cache-from"] != "type=gha" {
		t.Errorf("step.With[cache-from] = %v, want %q", step.With["cache-from"], "type=gha")
	}

	if step.With["cache-to"] != "type=gha,mode=max" {
		t.Errorf("step.With[cache-to] = %v, want %q", step.With["cache-to"], "type=gha,mode=max")
	}

	if step.With["target"] != "production" {
		t.Errorf("step.With[target] = %v, want %q", step.With["target"], "production")
	}

	if step.With["pull"] != true {
		t.Errorf("step.With[pull] = %v, want true", step.With["pull"])
	}

	if step.With["provenance"] != "false" {
		t.Errorf("step.With[provenance] = %v, want %q", step.With["provenance"], "false")
	}

	if step.With["sbom"] != "false" {
		t.Errorf("step.With[sbom] = %v, want %q", step.With["sbom"], "false")
	}
}

func TestDockerBuildPush_ToStep_MultiPlatform(t *testing.T) {
	d := DockerBuildPush{
		Context:   ".",
		Push:      true,
		Tags:      "user/app:latest",
		Platforms: "linux/amd64,linux/arm64,linux/arm/v7",
	}

	step := d.ToStep()

	if step.With["platforms"] != "linux/amd64,linux/arm64,linux/arm/v7" {
		t.Errorf("step.With[platforms] = %v, want multi-platform string", step.With["platforms"])
	}
}
