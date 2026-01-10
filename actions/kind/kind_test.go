package kind

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestKind_Action(t *testing.T) {
	k := Kind{}
	if got := k.Action(); got != "helm/kind-action@v1" {
		t.Errorf("Action() = %q, want %q", got, "helm/kind-action@v1")
	}
}

func TestKind_Inputs_Empty(t *testing.T) {
	k := Kind{}
	inputs := k.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty Kind.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestKind_Inputs_Version(t *testing.T) {
	k := Kind{
		Version: "v0.20.0",
	}

	inputs := k.Inputs()

	if inputs["version"] != "v0.20.0" {
		t.Errorf("inputs[version] = %v, want %q", inputs["version"], "v0.20.0")
	}
}

func TestKind_Inputs_Config(t *testing.T) {
	k := Kind{
		Config: "./kind-config.yaml",
	}

	inputs := k.Inputs()

	if inputs["config"] != "./kind-config.yaml" {
		t.Errorf("inputs[config] = %v, want %q", inputs["config"], "./kind-config.yaml")
	}
}

func TestKind_Inputs_ClusterName(t *testing.T) {
	k := Kind{
		ClusterName: "test-cluster",
	}

	inputs := k.Inputs()

	if inputs["cluster_name"] != "test-cluster" {
		t.Errorf("inputs[cluster_name] = %v, want %q", inputs["cluster_name"], "test-cluster")
	}
}

func TestKind_Inputs_WaitDuration(t *testing.T) {
	k := Kind{
		WaitDuration: "10m",
	}

	inputs := k.Inputs()

	if inputs["wait"] != "10m" {
		t.Errorf("inputs[wait] = %v, want %q", inputs["wait"], "10m")
	}
}

func TestKind_Inputs_Verbosity(t *testing.T) {
	k := Kind{
		Verbosity: 5,
	}

	inputs := k.Inputs()

	if inputs["verbosity"] != 5 {
		t.Errorf("inputs[verbosity] = %v, want 5", inputs["verbosity"])
	}
}

func TestKind_Inputs_Kubeconfig(t *testing.T) {
	k := Kind{
		Kubeconfig: "/tmp/kubeconfig",
	}

	inputs := k.Inputs()

	if inputs["kubeconfig"] != "/tmp/kubeconfig" {
		t.Errorf("inputs[kubeconfig] = %v, want %q", inputs["kubeconfig"], "/tmp/kubeconfig")
	}
}

func TestKind_Inputs_Registry(t *testing.T) {
	k := Kind{
		Registry: true,
	}

	inputs := k.Inputs()

	if inputs["registry"] != true {
		t.Errorf("inputs[registry] = %v, want true", inputs["registry"])
	}
}

func TestKind_Inputs_KubectlVersion(t *testing.T) {
	k := Kind{
		KubectlVersion: "v1.28.0",
	}

	inputs := k.Inputs()

	if inputs["kubectl_version"] != "v1.28.0" {
		t.Errorf("inputs[kubectl_version] = %v, want %q", inputs["kubectl_version"], "v1.28.0")
	}
}

func TestKind_Inputs_InstallOnly(t *testing.T) {
	k := Kind{
		InstallOnly: true,
	}

	inputs := k.Inputs()

	if inputs["install_only"] != true {
		t.Errorf("inputs[install_only] = %v, want true", inputs["install_only"])
	}
}

func TestKind_Inputs_IgnoreFailedClean(t *testing.T) {
	k := Kind{
		IgnoreFailedClean: true,
	}

	inputs := k.Inputs()

	if inputs["ignore_failed_clean"] != true {
		t.Errorf("inputs[ignore_failed_clean] = %v, want true", inputs["ignore_failed_clean"])
	}
}

func TestKind_ImplementsStepAction(t *testing.T) {
	k := Kind{}
	// Verify Kind implements StepAction interface
	var _ workflow.StepAction = k
}

func TestKind_Inputs_AllFields(t *testing.T) {
	k := Kind{
		Version:           "v0.20.0",
		Config:            "./kind-config.yaml",
		ClusterName:       "my-cluster",
		WaitDuration:      "15m",
		Verbosity:         3,
		Kubeconfig:        "/home/runner/.kube/config",
		Registry:          true,
		KubectlVersion:    "v1.28.0",
		InstallOnly:       true,
		IgnoreFailedClean: true,
	}

	inputs := k.Inputs()

	expected := map[string]any{
		"version":            "v0.20.0",
		"config":             "./kind-config.yaml",
		"cluster_name":       "my-cluster",
		"wait":               "15m",
		"verbosity":          3,
		"kubeconfig":         "/home/runner/.kube/config",
		"registry":           true,
		"kubectl_version":    "v1.28.0",
		"install_only":       true,
		"ignore_failed_clean": true,
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

func TestKind_Inputs_FalseBoolFields(t *testing.T) {
	k := Kind{
		Registry:          false,
		InstallOnly:       false,
		IgnoreFailedClean: false,
	}

	inputs := k.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestKind_Inputs_ZeroVerbosity(t *testing.T) {
	k := Kind{
		Verbosity: 0,
	}

	inputs := k.Inputs()

	if _, exists := inputs["verbosity"]; exists {
		t.Errorf("inputs[verbosity] should not exist for Verbosity=0")
	}
}

func TestKind_Inputs_NegativeVerbosity(t *testing.T) {
	k := Kind{
		Verbosity: -1,
	}

	inputs := k.Inputs()

	if inputs["verbosity"] != -1 {
		t.Errorf("inputs[verbosity] = %v, want -1", inputs["verbosity"])
	}
}

func TestKind_Inputs_Combined(t *testing.T) {
	// Test typical use case: version + cluster name + config
	k := Kind{
		Version:     "v0.20.0",
		ClusterName: "integration-test",
		Config:      "kind-config.yaml",
	}

	inputs := k.Inputs()

	if len(inputs) != 3 {
		t.Errorf("inputs has %d entries, want 3", len(inputs))
	}

	if inputs["version"] != "v0.20.0" {
		t.Errorf("inputs[version] = %v, want %q", inputs["version"], "v0.20.0")
	}
	if inputs["cluster_name"] != "integration-test" {
		t.Errorf("inputs[cluster_name] = %v, want %q", inputs["cluster_name"], "integration-test")
	}
	if inputs["config"] != "kind-config.yaml" {
		t.Errorf("inputs[config] = %v, want %q", inputs["config"], "kind-config.yaml")
	}
}
