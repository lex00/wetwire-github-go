// Package aws_ecr_login provides a typed wrapper for aws-actions/amazon-ecr-login.
package aws_ecr_login

// AWSECRLogin wraps the aws-actions/amazon-ecr-login@v2 action.
// Authenticate to Amazon ECR Private or Public registries.
type AWSECRLogin struct {
	// Comma-separated list of AWS account IDs for ECR Private registries.
	// If not provided, assumes default ECR Private registry.
	Registries string `yaml:"registries,omitempty"`

	// ECR registry type: "private" or "public".
	RegistryType string `yaml:"registry-type,omitempty"`

	// Prevents docker password from appearing in action logs during debug mode.
	MaskPassword bool `yaml:"mask-password,omitempty"`

	// Bypass explicit logout during post-job cleanup.
	SkipLogout bool `yaml:"skip-logout,omitempty"`

	// Proxy for the AWS SDK agent.
	HTTPProxy string `yaml:"http-proxy,omitempty"`
}

// Action returns the action reference.
func (a AWSECRLogin) Action() string {
	return "aws-actions/amazon-ecr-login@v2"
}

// Inputs returns the action inputs as a map.
func (a AWSECRLogin) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Registries != "" {
		with["registries"] = a.Registries
	}
	if a.RegistryType != "" {
		with["registry-type"] = a.RegistryType
	}
	if a.MaskPassword {
		with["mask-password"] = a.MaskPassword
	}
	if a.SkipLogout {
		with["skip-logout"] = a.SkipLogout
	}
	if a.HTTPProxy != "" {
		with["http-proxy"] = a.HTTPProxy
	}

	return with
}
