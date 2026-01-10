// Package azure_docker_login provides a typed wrapper for azure/docker-login.
package azure_docker_login

// AzureDockerLogin wraps the azure/docker-login@v2 action.
// Login to Azure Container Registry or other Docker registries.
type AzureDockerLogin struct {
	// Container registry server URL (required).
	LoginServer string `yaml:"login-server,omitempty"`

	// Container registry username (required).
	Username string `yaml:"username,omitempty"`

	// Container registry password (required).
	Password string `yaml:"password,omitempty"`
}

// Action returns the action reference.
func (a AzureDockerLogin) Action() string {
	return "azure/docker-login@v2"
}

// Inputs returns the action inputs as a map.
func (a AzureDockerLogin) Inputs() map[string]any {
	with := make(map[string]any)

	if a.LoginServer != "" {
		with["login-server"] = a.LoginServer
	}
	if a.Username != "" {
		with["username"] = a.Username
	}
	if a.Password != "" {
		with["password"] = a.Password
	}

	return with
}
