package k8s_set_context

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestK8sSetContext_Action(t *testing.T) {
	k := K8sSetContext{}
	if got := k.Action(); got != "azure/k8s-set-context@v4" {
		t.Errorf("Action() = %q, want %q", got, "azure/k8s-set-context@v4")
	}
}

func TestK8sSetContext_Inputs_Empty(t *testing.T) {
	k := K8sSetContext{}
	inputs := k.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty K8sSetContext.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestK8sSetContext_Inputs_Method(t *testing.T) {
	k := K8sSetContext{
		Method: "kubeconfig",
	}

	inputs := k.Inputs()

	if inputs["method"] != "kubeconfig" {
		t.Errorf("inputs[method] = %v, want %q", inputs["method"], "kubeconfig")
	}
}

func TestK8sSetContext_Inputs_MethodServiceAccount(t *testing.T) {
	k := K8sSetContext{
		Method: "service-account",
	}

	inputs := k.Inputs()

	if inputs["method"] != "service-account" {
		t.Errorf("inputs[method] = %v, want %q", inputs["method"], "service-account")
	}
}

func TestK8sSetContext_Inputs_Kubeconfig(t *testing.T) {
	kubeconfig := `apiVersion: v1
kind: Config
clusters:
- name: my-cluster
  cluster:
    server: https://my-cluster.example.com`

	k := K8sSetContext{
		Kubeconfig: kubeconfig,
	}

	inputs := k.Inputs()

	if inputs["kubeconfig"] != kubeconfig {
		t.Errorf("inputs[kubeconfig] = %v, want %q", inputs["kubeconfig"], kubeconfig)
	}
}

func TestK8sSetContext_Inputs_Context(t *testing.T) {
	k := K8sSetContext{
		Context: "my-k8s-context",
	}

	inputs := k.Inputs()

	if inputs["context"] != "my-k8s-context" {
		t.Errorf("inputs[context] = %v, want %q", inputs["context"], "my-k8s-context")
	}
}

func TestK8sSetContext_Inputs_ClusterType(t *testing.T) {
	k := K8sSetContext{
		ClusterType: "aks",
	}

	inputs := k.Inputs()

	if inputs["cluster-type"] != "aks" {
		t.Errorf("inputs[cluster-type] = %v, want %q", inputs["cluster-type"], "aks")
	}
}

func TestK8sSetContext_Inputs_ClusterTypeArc(t *testing.T) {
	k := K8sSetContext{
		ClusterType: "arc",
	}

	inputs := k.Inputs()

	if inputs["cluster-type"] != "arc" {
		t.Errorf("inputs[cluster-type] = %v, want %q", inputs["cluster-type"], "arc")
	}
}

func TestK8sSetContext_Inputs_ClusterTypeGeneric(t *testing.T) {
	k := K8sSetContext{
		ClusterType: "generic",
	}

	inputs := k.Inputs()

	if inputs["cluster-type"] != "generic" {
		t.Errorf("inputs[cluster-type] = %v, want %q", inputs["cluster-type"], "generic")
	}
}

func TestK8sSetContext_Inputs_ResourceGroup(t *testing.T) {
	k := K8sSetContext{
		ResourceGroup: "my-resource-group",
	}

	inputs := k.Inputs()

	if inputs["resource-group"] != "my-resource-group" {
		t.Errorf("inputs[resource-group] = %v, want %q", inputs["resource-group"], "my-resource-group")
	}
}

func TestK8sSetContext_Inputs_ClusterName(t *testing.T) {
	k := K8sSetContext{
		ClusterName: "my-aks-cluster",
	}

	inputs := k.Inputs()

	if inputs["cluster-name"] != "my-aks-cluster" {
		t.Errorf("inputs[cluster-name] = %v, want %q", inputs["cluster-name"], "my-aks-cluster")
	}
}

func TestK8sSetContext_ImplementsStepAction(t *testing.T) {
	k := K8sSetContext{}
	// Verify K8sSetContext implements StepAction interface
	var _ workflow.StepAction = k
}

func TestK8sSetContext_Inputs_AllFields(t *testing.T) {
	k := K8sSetContext{
		Method:        "kubeconfig",
		Kubeconfig:    "kubeconfig-content",
		Context:       "production",
		ClusterType:   "aks",
		ResourceGroup: "prod-rg",
		ClusterName:   "prod-cluster",
	}

	inputs := k.Inputs()

	expected := map[string]any{
		"method":         "kubeconfig",
		"kubeconfig":     "kubeconfig-content",
		"context":        "production",
		"cluster-type":   "aks",
		"resource-group": "prod-rg",
		"cluster-name":   "prod-cluster",
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

func TestK8sSetContext_Inputs_KubeconfigMethod(t *testing.T) {
	// Test typical kubeconfig authentication method
	k := K8sSetContext{
		Method:     "kubeconfig",
		Kubeconfig: "${{ secrets.KUBECONFIG }}",
		Context:    "default",
	}

	inputs := k.Inputs()

	if len(inputs) != 3 {
		t.Errorf("inputs has %d entries, want 3", len(inputs))
	}

	if inputs["method"] != "kubeconfig" {
		t.Errorf("inputs[method] = %v, want %q", inputs["method"], "kubeconfig")
	}
	if inputs["kubeconfig"] != "${{ secrets.KUBECONFIG }}" {
		t.Errorf("inputs[kubeconfig] = %v, want %q", inputs["kubeconfig"], "${{ secrets.KUBECONFIG }}")
	}
	if inputs["context"] != "default" {
		t.Errorf("inputs[context] = %v, want %q", inputs["context"], "default")
	}
}

func TestK8sSetContext_Inputs_AKSMethod(t *testing.T) {
	// Test AKS cluster configuration
	k := K8sSetContext{
		ClusterType:   "aks",
		ResourceGroup: "my-rg",
		ClusterName:   "my-aks",
	}

	inputs := k.Inputs()

	if len(inputs) != 3 {
		t.Errorf("inputs has %d entries, want 3", len(inputs))
	}

	if inputs["cluster-type"] != "aks" {
		t.Errorf("inputs[cluster-type] = %v, want %q", inputs["cluster-type"], "aks")
	}
	if inputs["resource-group"] != "my-rg" {
		t.Errorf("inputs[resource-group] = %v, want %q", inputs["resource-group"], "my-rg")
	}
	if inputs["cluster-name"] != "my-aks" {
		t.Errorf("inputs[cluster-name] = %v, want %q", inputs["cluster-name"], "my-aks")
	}
}

func TestK8sSetContext_Inputs_ServiceAccountMethod(t *testing.T) {
	// Test service account method
	k := K8sSetContext{
		Method: "service-account",
	}

	inputs := k.Inputs()

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}

	if inputs["method"] != "service-account" {
		t.Errorf("inputs[method] = %v, want %q", inputs["method"], "service-account")
	}
}
