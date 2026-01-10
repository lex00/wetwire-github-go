package pulumi

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestPulumi_Action(t *testing.T) {
	p := Pulumi{}
	if got := p.Action(); got != "pulumi/actions@v6" {
		t.Errorf("Action() = %q, want %q", got, "pulumi/actions@v6")
	}
}

func TestPulumi_Inputs_Empty(t *testing.T) {
	p := Pulumi{}
	inputs := p.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Pulumi.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestPulumi_Inputs_Command(t *testing.T) {
	p := Pulumi{
		Command: "up",
	}

	inputs := p.Inputs()

	if inputs["command"] != "up" {
		t.Errorf("inputs[command] = %v, want %q", inputs["command"], "up")
	}
}

func TestPulumi_Inputs_CommandPreview(t *testing.T) {
	p := Pulumi{
		Command: "preview",
	}

	inputs := p.Inputs()

	if inputs["command"] != "preview" {
		t.Errorf("inputs[command] = %v, want %q", inputs["command"], "preview")
	}
}

func TestPulumi_Inputs_CommandDestroy(t *testing.T) {
	p := Pulumi{
		Command: "destroy",
	}

	inputs := p.Inputs()

	if inputs["command"] != "destroy" {
		t.Errorf("inputs[command] = %v, want %q", inputs["command"], "destroy")
	}
}

func TestPulumi_Inputs_StackName(t *testing.T) {
	p := Pulumi{
		StackName: "org/project/production",
	}

	inputs := p.Inputs()

	if inputs["stack-name"] != "org/project/production" {
		t.Errorf("inputs[stack-name] = %v, want %q", inputs["stack-name"], "org/project/production")
	}
}

func TestPulumi_Inputs_WorkDir(t *testing.T) {
	p := Pulumi{
		WorkDir: "./infra",
	}

	inputs := p.Inputs()

	if inputs["work-dir"] != "./infra" {
		t.Errorf("inputs[work-dir] = %v, want %q", inputs["work-dir"], "./infra")
	}
}

func TestPulumi_Inputs_CloudURL(t *testing.T) {
	p := Pulumi{
		CloudURL: "https://api.pulumi.com",
	}

	inputs := p.Inputs()

	if inputs["cloud-url"] != "https://api.pulumi.com" {
		t.Errorf("inputs[cloud-url] = %v, want %q", inputs["cloud-url"], "https://api.pulumi.com")
	}
}

func TestPulumi_Inputs_ConfigMap(t *testing.T) {
	p := Pulumi{
		ConfigMap: `{"aws:region": "us-east-1"}`,
	}

	inputs := p.Inputs()

	if inputs["config-map"] != `{"aws:region": "us-east-1"}` {
		t.Errorf("inputs[config-map] = %v, want %q", inputs["config-map"], `{"aws:region": "us-east-1"}`)
	}
}

func TestPulumi_Inputs_SecretsProvider(t *testing.T) {
	p := Pulumi{
		SecretsProvider: "awskms://alias/my-key",
	}

	inputs := p.Inputs()

	if inputs["secrets-provider"] != "awskms://alias/my-key" {
		t.Errorf("inputs[secrets-provider] = %v, want %q", inputs["secrets-provider"], "awskms://alias/my-key")
	}
}

func TestPulumi_Inputs_Color(t *testing.T) {
	p := Pulumi{
		Color: "always",
	}

	inputs := p.Inputs()

	if inputs["color"] != "always" {
		t.Errorf("inputs[color] = %v, want %q", inputs["color"], "always")
	}
}

func TestPulumi_Inputs_ColorNever(t *testing.T) {
	p := Pulumi{
		Color: "never",
	}

	inputs := p.Inputs()

	if inputs["color"] != "never" {
		t.Errorf("inputs[color] = %v, want %q", inputs["color"], "never")
	}
}

func TestPulumi_Inputs_Diff(t *testing.T) {
	p := Pulumi{
		Diff: true,
	}

	inputs := p.Inputs()

	if inputs["diff"] != true {
		t.Errorf("inputs[diff] = %v, want true", inputs["diff"])
	}
}

func TestPulumi_Inputs_CommentOnPR(t *testing.T) {
	p := Pulumi{
		CommentOnPR: true,
	}

	inputs := p.Inputs()

	if inputs["comment-on-pr"] != true {
		t.Errorf("inputs[comment-on-pr] = %v, want true", inputs["comment-on-pr"])
	}
}

func TestPulumi_Inputs_EditPRComment(t *testing.T) {
	p := Pulumi{
		EditPRComment: true,
	}

	inputs := p.Inputs()

	if inputs["edit-pr-comment"] != true {
		t.Errorf("inputs[edit-pr-comment] = %v, want true", inputs["edit-pr-comment"])
	}
}

func TestPulumi_Inputs_AllFields(t *testing.T) {
	p := Pulumi{
		Command:         "up",
		StackName:       "org/project/production",
		WorkDir:         "./infra",
		CloudURL:        "https://api.pulumi.com",
		ConfigMap:       `{"aws:region": "us-east-1"}`,
		SecretsProvider: "awskms://alias/my-key",
		Color:           "auto",
		Diff:            true,
		CommentOnPR:     true,
		EditPRComment:   true,
	}

	inputs := p.Inputs()

	expected := map[string]any{
		"command":          "up",
		"stack-name":       "org/project/production",
		"work-dir":         "./infra",
		"cloud-url":        "https://api.pulumi.com",
		"config-map":       `{"aws:region": "us-east-1"}`,
		"secrets-provider": "awskms://alias/my-key",
		"color":            "auto",
		"diff":             true,
		"comment-on-pr":    true,
		"edit-pr-comment":  true,
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestPulumi_Inputs_FalseBoolFields(t *testing.T) {
	p := Pulumi{
		Diff:          false,
		CommentOnPR:   false,
		EditPRComment: false,
	}

	inputs := p.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestPulumi_ImplementsStepAction(t *testing.T) {
	p := Pulumi{}
	// Verify Pulumi implements StepAction interface
	var _ workflow.StepAction = p
}
