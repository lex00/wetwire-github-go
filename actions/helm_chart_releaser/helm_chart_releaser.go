// Package helm_chart_releaser provides a typed wrapper for helm/chart-releaser-action.
package helm_chart_releaser

// HelmChartReleaser wraps the helm/chart-releaser-action@v1 action.
// Turn your GitHub repo into a self-hosted Helm chart repository.
type HelmChartReleaser struct {
	// The version of chart-releaser to use
	Version string `yaml:"version,omitempty"`

	// Path to cr config file
	Config string `yaml:"config,omitempty"`

	// The directory containing the charts to be released
	ChartsDir string `yaml:"charts_dir,omitempty"`

	// The URL to the charts repository
	ChartsRepoURL string `yaml:"charts_repo_url,omitempty"`

	// Where to install the cr tool
	InstallDir string `yaml:"install_dir,omitempty"`

	// Just install cr tool
	InstallOnly bool `yaml:"install_only,omitempty"`

	// Skip the packaging step (useful if charts are already packaged)
	SkipPackaging bool `yaml:"skip_packaging,omitempty"`

	// Skip the upload step for releases that already exist
	SkipExisting bool `yaml:"skip_existing,omitempty"`

	// Skip package upload
	SkipUpload bool `yaml:"skip_upload,omitempty"`

	// Mark the release as latest
	MarkAsLatest bool `yaml:"mark_as_latest,omitempty"`

	// Upload chart packages directly into publishing branch
	PackagesWithIndex bool `yaml:"packages_with_index,omitempty"`

	// Name of the branch to be used to push the index and artifacts
	PagesBranch string `yaml:"pages_branch,omitempty"`
}

// Action returns the action reference.
func (a HelmChartReleaser) Action() string {
	return "helm/chart-releaser-action@v1"
}

// Inputs returns the action inputs as a map.
func (a HelmChartReleaser) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Version != "" {
		with["version"] = a.Version
	}
	if a.Config != "" {
		with["config"] = a.Config
	}
	if a.ChartsDir != "" {
		with["charts_dir"] = a.ChartsDir
	}
	if a.ChartsRepoURL != "" {
		with["charts_repo_url"] = a.ChartsRepoURL
	}
	if a.InstallDir != "" {
		with["install_dir"] = a.InstallDir
	}
	if a.InstallOnly {
		with["install_only"] = a.InstallOnly
	}
	if a.SkipPackaging {
		with["skip_packaging"] = a.SkipPackaging
	}
	if a.SkipExisting {
		with["skip_existing"] = a.SkipExisting
	}
	if a.SkipUpload {
		with["skip_upload"] = a.SkipUpload
	}
	if a.MarkAsLatest {
		with["mark_as_latest"] = a.MarkAsLatest
	}
	if a.PackagesWithIndex {
		with["packages_with_index"] = a.PackagesWithIndex
	}
	if a.PagesBranch != "" {
		with["pages_branch"] = a.PagesBranch
	}

	return with
}
