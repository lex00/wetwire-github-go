package helm_chart_releaser

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestHelmChartReleaser_Action(t *testing.T) {
	h := HelmChartReleaser{}
	if got := h.Action(); got != "helm/chart-releaser-action@v1" {
		t.Errorf("Action() = %q, want %q", got, "helm/chart-releaser-action@v1")
	}
}

func TestHelmChartReleaser_Inputs_Empty(t *testing.T) {
	h := HelmChartReleaser{}
	inputs := h.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty HelmChartReleaser.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestHelmChartReleaser_Inputs_Version(t *testing.T) {
	h := HelmChartReleaser{
		Version: "v1.6.0",
	}

	inputs := h.Inputs()

	if inputs["version"] != "v1.6.0" {
		t.Errorf("inputs[version] = %v, want %q", inputs["version"], "v1.6.0")
	}
}

func TestHelmChartReleaser_Inputs_Config(t *testing.T) {
	h := HelmChartReleaser{
		Config: ".github/cr.yaml",
	}

	inputs := h.Inputs()

	if inputs["config"] != ".github/cr.yaml" {
		t.Errorf("inputs[config] = %v, want %q", inputs["config"], ".github/cr.yaml")
	}
}

func TestHelmChartReleaser_Inputs_ChartsDir(t *testing.T) {
	h := HelmChartReleaser{
		ChartsDir: "charts",
	}

	inputs := h.Inputs()

	if inputs["charts_dir"] != "charts" {
		t.Errorf("inputs[charts_dir] = %v, want %q", inputs["charts_dir"], "charts")
	}
}

func TestHelmChartReleaser_Inputs_ChartsRepoURL(t *testing.T) {
	h := HelmChartReleaser{
		ChartsRepoURL: "https://example.github.io/helm-charts",
	}

	inputs := h.Inputs()

	if inputs["charts_repo_url"] != "https://example.github.io/helm-charts" {
		t.Errorf("inputs[charts_repo_url] = %v, want %q", inputs["charts_repo_url"], "https://example.github.io/helm-charts")
	}
}

func TestHelmChartReleaser_Inputs_SkipPackaging(t *testing.T) {
	h := HelmChartReleaser{
		SkipPackaging: true,
	}

	inputs := h.Inputs()

	if inputs["skip_packaging"] != true {
		t.Errorf("inputs[skip_packaging] = %v, want true", inputs["skip_packaging"])
	}
}

func TestHelmChartReleaser_Inputs_SkipExisting(t *testing.T) {
	h := HelmChartReleaser{
		SkipExisting: true,
	}

	inputs := h.Inputs()

	if inputs["skip_existing"] != true {
		t.Errorf("inputs[skip_existing] = %v, want true", inputs["skip_existing"])
	}
}

func TestHelmChartReleaser_Inputs_InstallDir(t *testing.T) {
	h := HelmChartReleaser{
		InstallDir: "/usr/local/bin",
	}

	inputs := h.Inputs()

	if inputs["install_dir"] != "/usr/local/bin" {
		t.Errorf("inputs[install_dir] = %v, want %q", inputs["install_dir"], "/usr/local/bin")
	}
}

func TestHelmChartReleaser_Inputs_InstallOnly(t *testing.T) {
	h := HelmChartReleaser{
		InstallOnly: true,
	}

	inputs := h.Inputs()

	if inputs["install_only"] != true {
		t.Errorf("inputs[install_only] = %v, want true", inputs["install_only"])
	}
}

func TestHelmChartReleaser_Inputs_SkipUpload(t *testing.T) {
	h := HelmChartReleaser{
		SkipUpload: true,
	}

	inputs := h.Inputs()

	if inputs["skip_upload"] != true {
		t.Errorf("inputs[skip_upload] = %v, want true", inputs["skip_upload"])
	}
}

func TestHelmChartReleaser_Inputs_MarkAsLatest(t *testing.T) {
	h := HelmChartReleaser{
		MarkAsLatest: true,
	}

	inputs := h.Inputs()

	if inputs["mark_as_latest"] != true {
		t.Errorf("inputs[mark_as_latest] = %v, want true", inputs["mark_as_latest"])
	}
}

func TestHelmChartReleaser_Inputs_PackagesWithIndex(t *testing.T) {
	h := HelmChartReleaser{
		PackagesWithIndex: true,
	}

	inputs := h.Inputs()

	if inputs["packages_with_index"] != true {
		t.Errorf("inputs[packages_with_index] = %v, want true", inputs["packages_with_index"])
	}
}

func TestHelmChartReleaser_Inputs_PagesBranch(t *testing.T) {
	h := HelmChartReleaser{
		PagesBranch: "gh-pages",
	}

	inputs := h.Inputs()

	if inputs["pages_branch"] != "gh-pages" {
		t.Errorf("inputs[pages_branch] = %v, want %q", inputs["pages_branch"], "gh-pages")
	}
}

func TestHelmChartReleaser_Inputs_AllFields(t *testing.T) {
	h := HelmChartReleaser{
		Version:           "v1.6.0",
		Config:            ".github/cr.yaml",
		ChartsDir:         "charts",
		ChartsRepoURL:     "https://example.github.io/helm-charts",
		InstallDir:        "/tmp/cr",
		InstallOnly:       true,
		SkipPackaging:     true,
		SkipExisting:      true,
		SkipUpload:        true,
		MarkAsLatest:      true,
		PackagesWithIndex: true,
		PagesBranch:       "gh-pages",
	}

	inputs := h.Inputs()

	expected := map[string]any{
		"version":             "v1.6.0",
		"config":              ".github/cr.yaml",
		"charts_dir":          "charts",
		"charts_repo_url":     "https://example.github.io/helm-charts",
		"install_dir":         "/tmp/cr",
		"install_only":        true,
		"skip_packaging":      true,
		"skip_existing":       true,
		"skip_upload":         true,
		"mark_as_latest":      true,
		"packages_with_index": true,
		"pages_branch":        "gh-pages",
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

func TestHelmChartReleaser_Inputs_FalseBoolFields(t *testing.T) {
	h := HelmChartReleaser{
		InstallOnly:       false,
		SkipPackaging:     false,
		SkipExisting:      false,
		SkipUpload:        false,
		MarkAsLatest:      false,
		PackagesWithIndex: false,
	}

	inputs := h.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestHelmChartReleaser_ImplementsStepAction(t *testing.T) {
	h := HelmChartReleaser{}
	// Verify HelmChartReleaser implements StepAction interface
	var _ workflow.StepAction = h
}
