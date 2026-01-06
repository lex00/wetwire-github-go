// Package setup_python provides a typed wrapper for actions/setup-python.
package setup_python

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// SetupPython wraps the actions/setup-python@v5 action.
// Setup a Python environment and add it to PATH.
type SetupPython struct {
	// Version range or exact version of Python to use
	PythonVersion string `yaml:"python-version,omitempty"`

	// File containing the Python version to use
	PythonVersionFile string `yaml:"python-version-file,omitempty"`

	// Used to specify a package manager for caching (pip, pipenv, poetry)
	Cache string `yaml:"cache,omitempty"`

	// Target architecture (x86, x64)
	Architecture string `yaml:"architecture,omitempty"`

	// Set this option if you want the action to check for latest available version
	CheckLatest bool `yaml:"check-latest,omitempty"`

	// Used to pull Python distributions
	Token string `yaml:"token,omitempty"`

	// The path to a dependency file (requirements.txt, Pipfile.lock, etc.)
	CacheDependencyPath string `yaml:"cache-dependency-path,omitempty"`

	// Set this option to true to update environment variables
	UpdateEnvironment bool `yaml:"update-environment,omitempty"`

	// Allow pre-release versions of Python to be installed
	AllowPrereleases bool `yaml:"allow-prereleases,omitempty"`
}

// Action returns the action reference.
func (a SetupPython) Action() string {
	return "actions/setup-python@v5"
}

// ToStep converts this action to a workflow step.
func (a SetupPython) ToStep() workflow.Step {
	with := make(workflow.With)

	if a.PythonVersion != "" {
		with["python-version"] = a.PythonVersion
	}
	if a.PythonVersionFile != "" {
		with["python-version-file"] = a.PythonVersionFile
	}
	if a.Cache != "" {
		with["cache"] = a.Cache
	}
	if a.Architecture != "" {
		with["architecture"] = a.Architecture
	}
	if a.CheckLatest {
		with["check-latest"] = a.CheckLatest
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.CacheDependencyPath != "" {
		with["cache-dependency-path"] = a.CacheDependencyPath
	}
	if a.UpdateEnvironment {
		with["update-environment"] = a.UpdateEnvironment
	}
	if a.AllowPrereleases {
		with["allow-prereleases"] = a.AllowPrereleases
	}

	return workflow.Step{
		Uses: a.Action(),
		With: with,
	}
}
