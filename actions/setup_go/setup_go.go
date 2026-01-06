// Package setup_go provides a typed wrapper for actions/setup-go.
package setup_go

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// SetupGo wraps the actions/setup-go@v5 action.
// Setup a Go environment and add it to PATH.
type SetupGo struct {
	// The Go version to download (if necessary) and use
	GoVersion string `yaml:"go-version,omitempty"`

	// Path to the go.mod or go.work file
	GoVersionFile string `yaml:"go-version-file,omitempty"`

	// Set this option to true if you want the action to always check for latest version
	CheckLatest bool `yaml:"check-latest,omitempty"`

	// Used to pull Go distributions. Since there's a default, this is typically not needed
	Token string `yaml:"token,omitempty"`

	// Used to specify whether caching is needed. Set to true if you'd like to enable caching
	Cache bool `yaml:"cache,omitempty"`

	// Used to specify the path to a dependency file: go.sum
	CacheDependencyPath string `yaml:"cache-dependency-path,omitempty"`

	// Target architecture for Go to use. Examples: x86, x64
	Architecture string `yaml:"architecture,omitempty"`
}

// Action returns the action reference.
func (a SetupGo) Action() string {
	return "actions/setup-go@v5"
}

// ToStep converts this action to a workflow step.
func (a SetupGo) ToStep() workflow.Step {
	with := make(workflow.With)

	if a.GoVersion != "" {
		with["go-version"] = a.GoVersion
	}
	if a.GoVersionFile != "" {
		with["go-version-file"] = a.GoVersionFile
	}
	if a.CheckLatest {
		with["check-latest"] = a.CheckLatest
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.Cache {
		with["cache"] = a.Cache
	}
	if a.CacheDependencyPath != "" {
		with["cache-dependency-path"] = a.CacheDependencyPath
	}
	if a.Architecture != "" {
		with["architecture"] = a.Architecture
	}

	return workflow.Step{
		Uses: a.Action(),
		With: with,
	}
}
