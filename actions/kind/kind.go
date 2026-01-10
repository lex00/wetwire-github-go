// Package kind provides a typed wrapper for helm/kind-action.
package kind

// Kind wraps the helm/kind-action@v1 action.
// Creates a Kubernetes cluster using kind (Kubernetes IN Docker).
type Kind struct {
	// Version of kind to use (e.g., v0.20.0)
	Version string `yaml:"version,omitempty"`

	// Path to kind config file
	Config string `yaml:"config,omitempty"`

	// Name of the kind cluster
	ClusterName string `yaml:"cluster_name,omitempty"`

	// How long to wait for control plane to become ready (default: 5m)
	WaitDuration string `yaml:"wait,omitempty"`

	// Log verbosity level for kind (0-9)
	Verbosity int `yaml:"verbosity,omitempty"`

	// Path to write kubeconfig file
	Kubeconfig string `yaml:"kubeconfig,omitempty"`

	// Enable local registry for the cluster
	Registry bool `yaml:"registry,omitempty"`

	// Version of kubectl to install
	KubectlVersion string `yaml:"kubectl_version,omitempty"`

	// Only install kind without creating cluster
	InstallOnly bool `yaml:"install_only,omitempty"`

	// Ignore cluster cleanup failures
	IgnoreFailedClean bool `yaml:"ignore_failed_clean,omitempty"`
}

// Action returns the action reference.
func (a Kind) Action() string {
	return "helm/kind-action@v1"
}

// Inputs returns the action inputs as a map.
func (a Kind) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.Config != "" {
		with["config"] = a.Config
	}
	if a.ClusterName != "" {
		with["cluster_name"] = a.ClusterName
	}
	if a.WaitDuration != "" {
		with["wait"] = a.WaitDuration
	}
	if a.Verbosity != 0 {
		with["verbosity"] = a.Verbosity
	}
	if a.Kubeconfig != "" {
		with["kubeconfig"] = a.Kubeconfig
	}
	if a.Registry {
		with["registry"] = a.Registry
	}
	if a.KubectlVersion != "" {
		with["kubectl_version"] = a.KubectlVersion
	}
	if a.InstallOnly {
		with["install_only"] = a.InstallOnly
	}
	if a.IgnoreFailedClean {
		with["ignore_failed_clean"] = a.IgnoreFailedClean
	}

	return with
}
