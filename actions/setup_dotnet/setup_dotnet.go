// Package setup_dotnet provides a typed wrapper for actions/setup-dotnet.
package setup_dotnet

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// SetupDotnet wraps the actions/setup-dotnet@v4 action.
// Set up a specific version of the .NET SDK and add it to PATH.
type SetupDotnet struct {
	// Optional SDK version(s) to use. If not provided, will install global.json version if available.
	// Examples: "6.0.x", "7.0.x", "8.0.x", "8.0.100"
	DotnetVersion string `yaml:"dotnet-version,omitempty"`

	// Optional quality of the build to download from the dotnet channel.
	// The possible values are: daily, signed, validated, preview, ga.
	DotnetQuality string `yaml:"dotnet-quality,omitempty"`

	// Optional global.json location, if your global.json isn't in the root of the repo.
	GlobalJsonFile string `yaml:"global-json-file,omitempty"`

	// Whether prerelease versions should be matched with non-exact versions (eg. 6.0.0-preview)
	IncludePrerelease bool `yaml:"include-prerelease,omitempty"`

	// Used to specify the source for the .NET SDK
	Source string `yaml:"source,omitempty"`

	// Optional NUGET_AUTH_TOKEN for authentication to private NuGet repositories
	Token string `yaml:"token,omitempty"`

	// Optional NuGet.config location for custom NuGet sources
	ConfigFile string `yaml:"config-file,omitempty"`

	// Used to specify whether caching is needed. Set to true to enable caching
	Cache bool `yaml:"cache,omitempty"`

	// Used to specify the path to a dependency file. Supports packages.lock.json
	CacheDependencyPath string `yaml:"cache-dependency-path,omitempty"`
}

// Action returns the action reference.
func (a SetupDotnet) Action() string {
	return "actions/setup-dotnet@v4"
}

// ToStep converts this action to a workflow step.
func (a SetupDotnet) ToStep() workflow.Step {
	with := make(workflow.With)

	if a.DotnetVersion != "" {
		with["dotnet-version"] = a.DotnetVersion
	}
	if a.DotnetQuality != "" {
		with["dotnet-quality"] = a.DotnetQuality
	}
	if a.GlobalJsonFile != "" {
		with["global-json-file"] = a.GlobalJsonFile
	}
	if a.IncludePrerelease {
		with["include-prerelease"] = a.IncludePrerelease
	}
	if a.Source != "" {
		with["source"] = a.Source
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.ConfigFile != "" {
		with["config-file"] = a.ConfigFile
	}
	if a.Cache {
		with["cache"] = a.Cache
	}
	if a.CacheDependencyPath != "" {
		with["cache-dependency-path"] = a.CacheDependencyPath
	}

	return workflow.Step{
		Uses: a.Action(),
		With: with,
	}
}
