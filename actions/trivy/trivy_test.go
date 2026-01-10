package trivy

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestTrivy_Action(t *testing.T) {
	tr := Trivy{}
	if got := tr.Action(); got != "aquasecurity/trivy-action@0.28.0" {
		t.Errorf("Action() = %q, want %q", got, "aquasecurity/trivy-action@0.28.0")
	}
}

func TestTrivy_Inputs(t *testing.T) {
	tr := Trivy{
		ImageRef: "alpine:3.18",
		ScanType: "image",
		Format:   "sarif",
		Severity: "CRITICAL,HIGH",
	}

	inputs := tr.Inputs()

	if inputs["image-ref"] != "alpine:3.18" {
		t.Errorf("inputs[image-ref] = %v, want %q", inputs["image-ref"], "alpine:3.18")
	}

	if inputs["scan-type"] != "image" {
		t.Errorf("inputs[scan-type] = %v, want %q", inputs["scan-type"], "image")
	}

	if inputs["format"] != "sarif" {
		t.Errorf("inputs[format] = %v, want %q", inputs["format"], "sarif")
	}

	if inputs["severity"] != "CRITICAL,HIGH" {
		t.Errorf("inputs[severity] = %v, want %q", inputs["severity"], "CRITICAL,HIGH")
	}
}

func TestTrivy_Inputs_Empty(t *testing.T) {
	tr := Trivy{}
	inputs := tr.Inputs()

	// Empty Trivy should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Trivy.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestTrivy_Inputs_BoolFields(t *testing.T) {
	tr := Trivy{
		IgnoreUnfixed: true,
		ExitCode:      1,
	}

	inputs := tr.Inputs()

	if inputs["ignore-unfixed"] != true {
		t.Errorf("inputs[ignore-unfixed] = %v, want true", inputs["ignore-unfixed"])
	}

	if inputs["exit-code"] != 1 {
		t.Errorf("inputs[exit-code] = %v, want 1", inputs["exit-code"])
	}
}

func TestTrivy_Inputs_AllFields(t *testing.T) {
	tr := Trivy{
		ImageRef:      "myapp:latest",
		ScanType:      "image",
		Format:        "json",
		Severity:      "CRITICAL,HIGH,MEDIUM",
		ExitCode:      1,
		IgnoreUnfixed: true,
		VulnType:      "os,library",
		Scanners:      "vuln,config,secret",
		Template:      "@/path/to/template.tpl",
		Output:        "trivy-results.json",
	}

	inputs := tr.Inputs()

	if inputs["image-ref"] != "myapp:latest" {
		t.Errorf("inputs[image-ref] = %v, want %q", inputs["image-ref"], "myapp:latest")
	}

	if inputs["scan-type"] != "image" {
		t.Errorf("inputs[scan-type] = %v, want %q", inputs["scan-type"], "image")
	}

	if inputs["format"] != "json" {
		t.Errorf("inputs[format] = %v, want %q", inputs["format"], "json")
	}

	if inputs["severity"] != "CRITICAL,HIGH,MEDIUM" {
		t.Errorf("inputs[severity] = %v, want %q", inputs["severity"], "CRITICAL,HIGH,MEDIUM")
	}

	if inputs["exit-code"] != 1 {
		t.Errorf("inputs[exit-code] = %v, want 1", inputs["exit-code"])
	}

	if inputs["ignore-unfixed"] != true {
		t.Errorf("inputs[ignore-unfixed] = %v, want true", inputs["ignore-unfixed"])
	}

	if inputs["vuln-type"] != "os,library" {
		t.Errorf("inputs[vuln-type] = %v, want %q", inputs["vuln-type"], "os,library")
	}

	if inputs["scanners"] != "vuln,config,secret" {
		t.Errorf("inputs[scanners] = %v, want %q", inputs["scanners"], "vuln,config,secret")
	}

	if inputs["template"] != "@/path/to/template.tpl" {
		t.Errorf("inputs[template] = %v, want %q", inputs["template"], "@/path/to/template.tpl")
	}

	if inputs["output"] != "trivy-results.json" {
		t.Errorf("inputs[output] = %v, want %q", inputs["output"], "trivy-results.json")
	}
}

func TestTrivy_Inputs_FilesystemScan(t *testing.T) {
	tr := Trivy{
		ScanType: "fs",
		Format:   "table",
		Severity: "HIGH,CRITICAL",
	}

	inputs := tr.Inputs()

	if inputs["scan-type"] != "fs" {
		t.Errorf("inputs[scan-type] = %v, want %q", inputs["scan-type"], "fs")
	}

	if inputs["format"] != "table" {
		t.Errorf("inputs[format] = %v, want %q", inputs["format"], "table")
	}
}

func TestTrivy_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = Trivy{}
}
