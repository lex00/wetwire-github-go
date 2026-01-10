// Package kustomize provides a typed wrapper for stefanprodan/kustomize-action.
package kustomize

// Kustomize wraps the stefanprodan/kustomize-action@master action.
// Run kustomize build and apply the resulting manifests for GitOps deployments.
type Kustomize struct {
	// Path to the kustomization directory
	Kustomization string `yaml:"kustomization,omitempty"`
}

// Action returns the action reference.
func (a Kustomize) Action() string {
	return "stefanprodan/kustomize-action@master"
}

// Inputs returns the action inputs as a map.
func (a Kustomize) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Kustomization != "" {
		with["kustomization"] = a.Kustomization
	}

	return with
}
