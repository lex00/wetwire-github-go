// Package configure_pages provides a typed wrapper for actions/configure-pages.
package configure_pages

// ConfigurePages wraps the actions/configure-pages@v5 action.
// Configures GitHub Pages for deployment.
type ConfigurePages struct {
	// Static site generator to configure ("next", "nuxt", "gatsby", "jekyll", etc.)
	StaticSiteGenerator string `yaml:"static_site_generator,omitempty"`

	// Path to the generator configuration file
	GeneratorConfigFile string `yaml:"generator_config_file,omitempty"`

	// GitHub token for authentication
	Token string `yaml:"token,omitempty"`
}

// Action returns the action reference.
func (a ConfigurePages) Action() string {
	return "actions/configure-pages@v5"
}

// Inputs returns the action inputs as a map.
func (a ConfigurePages) Inputs() map[string]any {
	m := make(map[string]any)

	if a.StaticSiteGenerator != "" {
		m["static_site_generator"] = a.StaticSiteGenerator
	}
	if a.GeneratorConfigFile != "" {
		m["generator_config_file"] = a.GeneratorConfigFile
	}
	if a.Token != "" {
		m["token"] = a.Token
	}

	return m
}
