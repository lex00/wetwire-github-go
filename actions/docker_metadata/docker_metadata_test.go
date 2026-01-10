package docker_metadata

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestDockerMetadata_Action(t *testing.T) {
	d := DockerMetadata{}
	if got := d.Action(); got != "docker/metadata-action@v5" {
		t.Errorf("Action() = %q, want %q", got, "docker/metadata-action@v5")
	}
}

func TestDockerMetadata_Inputs(t *testing.T) {
	d := DockerMetadata{
		Images: "ghcr.io/user/app",
		Tags:   "type=ref,event=branch\ntype=semver,pattern={{version}}",
	}

	inputs := d.Inputs()

	if d.Action() != "docker/metadata-action@v5" {
		t.Errorf("Action() = %q, want %q", d.Action(), "docker/metadata-action@v5")
	}

	if inputs["images"] != "ghcr.io/user/app" {
		t.Errorf("inputs[images] = %v, want %q", inputs["images"], "ghcr.io/user/app")
	}

	if inputs["tags"] != "type=ref,event=branch\ntype=semver,pattern={{version}}" {
		t.Errorf("inputs[tags] = %v, want %q", inputs["tags"], "type=ref,event=branch\ntype=semver,pattern={{version}}")
	}
}

func TestDockerMetadata_Inputs_EmptyWithMaps(t *testing.T) {
	d := DockerMetadata{}
	inputs := d.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty DockerMetadata.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestDockerMetadata_Inputs_AllFields(t *testing.T) {
	d := DockerMetadata{
		Context:        "git",
		Images:         "ghcr.io/user/app\nuser/app",
		Tags:           "type=ref,event=branch\ntype=ref,event=pr\ntype=semver,pattern={{version}}",
		Flavor:         "latest=auto\nprefix=\nsuffix=",
		Labels:         "org.opencontainers.image.source=https://github.com/user/repo\norg.opencontainers.image.description=My App",
		Annotations:    "org.opencontainers.image.source=https://github.com/user/repo",
		SepTags:        ",",
		SepLabels:      ",",
		SepAnnotations: ",",
		BakeTarget:     "my-target",
	}

	inputs := d.Inputs()

	if inputs["context"] != "git" {
		t.Errorf("inputs[context] = %v, want %q", inputs["context"], "git")
	}

	if inputs["images"] != "ghcr.io/user/app\nuser/app" {
		t.Errorf("inputs[images] = %v, want multi-line images", inputs["images"])
	}

	if inputs["tags"] != "type=ref,event=branch\ntype=ref,event=pr\ntype=semver,pattern={{version}}" {
		t.Errorf("inputs[tags] = %v, want multi-line tags", inputs["tags"])
	}

	if inputs["flavor"] != "latest=auto\nprefix=\nsuffix=" {
		t.Errorf("inputs[flavor] = %v, want %q", inputs["flavor"], "latest=auto\nprefix=\nsuffix=")
	}

	if inputs["labels"] != "org.opencontainers.image.source=https://github.com/user/repo\norg.opencontainers.image.description=My App" {
		t.Errorf("inputs[labels] = %v, want multi-line labels", inputs["labels"])
	}

	if inputs["annotations"] != "org.opencontainers.image.source=https://github.com/user/repo" {
		t.Errorf("inputs[annotations] = %v, want %q", inputs["annotations"], "org.opencontainers.image.source=https://github.com/user/repo")
	}

	if inputs["sep-tags"] != "," {
		t.Errorf("inputs[sep-tags] = %v, want %q", inputs["sep-tags"], ",")
	}

	if inputs["sep-labels"] != "," {
		t.Errorf("inputs[sep-labels] = %v, want %q", inputs["sep-labels"], ",")
	}

	if inputs["sep-annotations"] != "," {
		t.Errorf("inputs[sep-annotations] = %v, want %q", inputs["sep-annotations"], ",")
	}

	if inputs["bake-target"] != "my-target" {
		t.Errorf("inputs[bake-target] = %v, want %q", inputs["bake-target"], "my-target")
	}
}

func TestDockerMetadata_Inputs_MultiImageAndTags(t *testing.T) {
	d := DockerMetadata{
		Images: "ghcr.io/user/app\ndockerhub/user/app",
		Tags:   "type=ref,event=branch\ntype=semver,pattern={{version}}\ntype=sha",
	}

	inputs := d.Inputs()

	if inputs["images"] != "ghcr.io/user/app\ndockerhub/user/app" {
		t.Errorf("inputs[images] = %v, want multi-line images", inputs["images"])
	}

	if inputs["tags"] != "type=ref,event=branch\ntype=semver,pattern={{version}}\ntype=sha" {
		t.Errorf("inputs[tags] = %v, want multi-line tags", inputs["tags"])
	}
}

func TestDockerMetadata_Inputs_OnlyContext(t *testing.T) {
	d := DockerMetadata{
		Context: "workflow",
	}

	inputs := d.Inputs()

	if inputs["context"] != "workflow" {
		t.Errorf("inputs[context] = %v, want %q", inputs["context"], "workflow")
	}

	if len(inputs) != 1 {
		t.Errorf("expected only 1 input, got %d", len(inputs))
	}
}

func TestDockerMetadata_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = DockerMetadata{}
}
