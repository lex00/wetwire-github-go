// Package k8s_set_context provides a typed wrapper for azure/k8s-set-context.
package k8s_set_context

// K8sSetContext wraps the azure/k8s-set-context@v4 action.
// Sets the Kubernetes context for deploying to AKS or any Kubernetes cluster.
type K8sSetContext struct {
	// Authentication method: kubeconfig or service-account
	Method string `yaml:"method,omitempty"`

	// Contents of kubeconfig file (for kubeconfig method)
	Kubeconfig string `yaml:"kubeconfig,omitempty"`

	// Context name to use from kubeconfig
	Context string `yaml:"context,omitempty"`

	// Cluster type: generic, arc, or aks
	ClusterType string `yaml:"cluster-type,omitempty"`

	// Azure resource group containing the cluster (for AKS/Arc)
	ResourceGroup string `yaml:"resource-group,omitempty"`

	// Name of the AKS/Arc cluster
	ClusterName string `yaml:"cluster-name,omitempty"`
}

// Action returns the action reference.
func (a K8sSetContext) Action() string {
	return "azure/k8s-set-context@v4"
}

// Inputs returns the action inputs as a map.
func (a K8sSetContext) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Method != "" {
		with["method"] = a.Method
	}
	if a.Kubeconfig != "" {
		with["kubeconfig"] = a.Kubeconfig
	}
	if a.Context != "" {
		with["context"] = a.Context
	}
	if a.ClusterType != "" {
		with["cluster-type"] = a.ClusterType
	}
	if a.ResourceGroup != "" {
		with["resource-group"] = a.ResourceGroup
	}
	if a.ClusterName != "" {
		with["cluster-name"] = a.ClusterName
	}

	return with
}
