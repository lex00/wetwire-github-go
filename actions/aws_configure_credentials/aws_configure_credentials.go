// Package aws_configure_credentials provides a typed wrapper for aws-actions/configure-aws-credentials.
package aws_configure_credentials

// AWSConfigureCredentials wraps the aws-actions/configure-aws-credentials@v4 action.
// Configure AWS credentials for use in subsequent steps.
type AWSConfigureCredentials struct {
	// AWS Region (required). (e.g., "us-east-1", "eu-west-1")
	AWSRegion string `yaml:"aws-region,omitempty"`

	// ARN of IAM role to assume using OIDC or access keys.
	RoleToAssume string `yaml:"role-to-assume,omitempty"`

	// AWS Access Key ID for assuming a role with access keys.
	AWSAccessKeyID string `yaml:"aws-access-key-id,omitempty"`

	// AWS Secret Access Key. Required if AWSAccessKeyID is provided.
	AWSSecretAccessKey string `yaml:"aws-secret-access-key,omitempty"`

	// AWS Session Token.
	AWSSessionToken string `yaml:"aws-session-token,omitempty"`

	// Path to web identity token file.
	WebIdentityTokenFile string `yaml:"web-identity-token-file,omitempty"`

	// Use environment credentials to assume a new role.
	RoleChaining bool `yaml:"role-chaining,omitempty"`

	// Audience for the OIDC provider.
	Audience string `yaml:"audience,omitempty"`

	// Proxy for the AWS SDK agent.
	HTTPProxy string `yaml:"http-proxy,omitempty"`

	// Role duration in seconds.
	RoleDurationSeconds int `yaml:"role-duration-seconds,omitempty"`

	// External ID of role to assume.
	RoleExternalID string `yaml:"role-external-id,omitempty"`

	// Role session name (default: GitHubActions).
	RoleSessionName string `yaml:"role-session-name,omitempty"`

	// Skip session tagging during role assumption.
	RoleSkipSessionTagging bool `yaml:"role-skip-session-tagging,omitempty"`

	// Inline policy document for role assumption.
	InlineSessionPolicy string `yaml:"inline-session-policy,omitempty"`

	// List of managed policy ARNs for role assumption.
	ManagedSessionPolicies string `yaml:"managed-session-policies,omitempty"`

	// Set credentials as step output.
	OutputCredentials bool `yaml:"output-credentials,omitempty"`

	// Mask AWS account ID in logs.
	MaskAWSAccountID bool `yaml:"mask-aws-account-id,omitempty"`

	// Unset existing runner credentials.
	UnsetCurrentCredentials bool `yaml:"unset-current-credentials,omitempty"`

	// Disable retry mechanism for assume role.
	DisableRetry bool `yaml:"disable-retry,omitempty"`

	// Maximum retry attempts for assume role.
	RetryMaxAttempts int `yaml:"retry-max-attempts,omitempty"`

	// Retry until secret key lacks special characters.
	SpecialCharactersWorkaround bool `yaml:"special-characters-workaround,omitempty"`
}

// Action returns the action reference.
func (a AWSConfigureCredentials) Action() string {
	return "aws-actions/configure-aws-credentials@v4"
}

// Inputs returns the action inputs as a map.
func (a AWSConfigureCredentials) Inputs() map[string]any {
	with := make(map[string]any)

	if a.AWSRegion != "" {
		with["aws-region"] = a.AWSRegion
	}
	if a.RoleToAssume != "" {
		with["role-to-assume"] = a.RoleToAssume
	}
	if a.AWSAccessKeyID != "" {
		with["aws-access-key-id"] = a.AWSAccessKeyID
	}
	if a.AWSSecretAccessKey != "" {
		with["aws-secret-access-key"] = a.AWSSecretAccessKey
	}
	if a.AWSSessionToken != "" {
		with["aws-session-token"] = a.AWSSessionToken
	}
	if a.WebIdentityTokenFile != "" {
		with["web-identity-token-file"] = a.WebIdentityTokenFile
	}
	if a.RoleChaining {
		with["role-chaining"] = a.RoleChaining
	}
	if a.Audience != "" {
		with["audience"] = a.Audience
	}
	if a.HTTPProxy != "" {
		with["http-proxy"] = a.HTTPProxy
	}
	if a.RoleDurationSeconds != 0 {
		with["role-duration-seconds"] = a.RoleDurationSeconds
	}
	if a.RoleExternalID != "" {
		with["role-external-id"] = a.RoleExternalID
	}
	if a.RoleSessionName != "" {
		with["role-session-name"] = a.RoleSessionName
	}
	if a.RoleSkipSessionTagging {
		with["role-skip-session-tagging"] = a.RoleSkipSessionTagging
	}
	if a.InlineSessionPolicy != "" {
		with["inline-session-policy"] = a.InlineSessionPolicy
	}
	if a.ManagedSessionPolicies != "" {
		with["managed-session-policies"] = a.ManagedSessionPolicies
	}
	if a.OutputCredentials {
		with["output-credentials"] = a.OutputCredentials
	}
	if a.MaskAWSAccountID {
		with["mask-aws-account-id"] = a.MaskAWSAccountID
	}
	if a.UnsetCurrentCredentials {
		with["unset-current-credentials"] = a.UnsetCurrentCredentials
	}
	if a.DisableRetry {
		with["disable-retry"] = a.DisableRetry
	}
	if a.RetryMaxAttempts != 0 {
		with["retry-max-attempts"] = a.RetryMaxAttempts
	}
	if a.SpecialCharactersWorkaround {
		with["special-characters-workaround"] = a.SpecialCharactersWorkaround
	}

	return with
}
