// Package setup_ruby provides a typed wrapper for ruby/setup-ruby.
package setup_ruby

// SetupRuby wraps the ruby/setup-ruby@v1 action.
// Set up a specific version of Ruby and add it to PATH.
type SetupRuby struct {
	// Ruby version to use. Use 'ruby-head' for the latest development version.
	// Examples: "3.2", "3.3", "3.3.0", "ruby-head", "jruby-9.4"
	RubyVersion string `yaml:"ruby-version,omitempty"`

	// The path to a .ruby-version or .tool-versions file
	RubyVersionFile string `yaml:"ruby-version-file,omitempty"`

	// Whether to run bundle install, and which command to use
	// Examples: "true", "false", "Gemfile", "none"
	Bundler string `yaml:"bundler,omitempty"`

	// The version of Bundler to install. Examples: "default", "latest", "none", "2.4.0"
	BundlerVersion string `yaml:"bundler-version,omitempty"`

	// Whether to use Bundler caching
	BundlerCache bool `yaml:"bundler-cache,omitempty"`

	// The path to the Gemfile.lock for caching
	CacheVersion string `yaml:"cache-version,omitempty"`

	// The working directory to use
	WorkingDirectory string `yaml:"working-directory,omitempty"`

	// The path to the rubygems installation directory
	RubygemsPath string `yaml:"rubygems,omitempty"`

	// Set the Bundler lockfile to allow unmet dependencies
	BundlerNoLock bool `yaml:"bundler-no-lock,omitempty"`
}

// Action returns the action reference.
func (a SetupRuby) Action() string {
	return "ruby/setup-ruby@v1"
}

// Inputs returns the action inputs as a map.
func (a SetupRuby) Inputs() map[string]any {
	with := make(map[string]any)

	if a.RubyVersion != "" {
		with["ruby-version"] = a.RubyVersion
	}
	if a.RubyVersionFile != "" {
		with["ruby-version-file"] = a.RubyVersionFile
	}
	if a.Bundler != "" {
		with["bundler"] = a.Bundler
	}
	if a.BundlerVersion != "" {
		with["bundler-version"] = a.BundlerVersion
	}
	if a.BundlerCache {
		with["bundler-cache"] = a.BundlerCache
	}
	if a.CacheVersion != "" {
		with["cache-version"] = a.CacheVersion
	}
	if a.WorkingDirectory != "" {
		with["working-directory"] = a.WorkingDirectory
	}
	if a.RubygemsPath != "" {
		with["rubygems"] = a.RubygemsPath
	}
	if a.BundlerNoLock {
		with["bundler-no-lock"] = a.BundlerNoLock
	}

	return with
}
