// Package docker_login provides a typed wrapper for docker/login-action.
package docker_login

// DockerLogin wraps the docker/login-action@v3 action.
// Log in to a Docker registry (Docker Hub, GitHub Container Registry, AWS ECR, etc.).
type DockerLogin struct {
	// Server address of Docker registry. Defaults to Docker Hub.
	Registry string `yaml:"registry,omitempty"`

	// Username for authentication.
	Username string `yaml:"username,omitempty"`

	// Password or personal access token for authentication.
	Password string `yaml:"password,omitempty"`

	// AWS ECR configuration. Can be "auto" to auto-detect.
	ECR string `yaml:"ecr,omitempty"`

	// Whether to logout from the registry at the end of the job.
	Logout bool `yaml:"logout,omitempty"`
}

// Action returns the action reference.
func (a DockerLogin) Action() string {
	return "docker/login-action@v3"
}

// Inputs returns the action inputs as a map.
func (a DockerLogin) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Registry != "" {
		with["registry"] = a.Registry
	}
	if a.Username != "" {
		with["username"] = a.Username
	}
	if a.Password != "" {
		with["password"] = a.Password
	}
	if a.ECR != "" {
		with["ecr"] = a.ECR
	}
	if a.Logout {
		with["logout"] = a.Logout
	}

	return with
}
