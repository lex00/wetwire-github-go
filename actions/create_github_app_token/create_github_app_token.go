// Package create_github_app_token provides a typed wrapper for actions/create-github-app-token.
package create_github_app_token

// CreateGithubAppToken wraps the actions/create-github-app-token@v1 action.
// Create GitHub App installation access tokens.
type CreateGithubAppToken struct {
	// GitHub App ID (required)
	AppID string `yaml:"app-id,omitempty"`

	// GitHub App private key (required)
	PrivateKey string `yaml:"private-key,omitempty"`

	// Owner of the GitHub App installation (defaults to current repository owner)
	Owner string `yaml:"owner,omitempty"`

	// Comma or newline-separated list of repositories to install the GitHub App on
	Repositories string `yaml:"repositories,omitempty"`

	// If true, the token will not be revoked when the current job is complete
	SkipTokenRevoke bool `yaml:"skip-token-revoke,omitempty"`

	// URL of the GitHub REST API
	GithubAPIURL string `yaml:"github-api-url,omitempty"`

	// Repository permissions
	PermissionActions            string `yaml:"permission-actions,omitempty"`
	PermissionAdministration     string `yaml:"permission-administration,omitempty"`
	PermissionChecks             string `yaml:"permission-checks,omitempty"`
	PermissionCodespaces         string `yaml:"permission-codespaces,omitempty"`
	PermissionContents           string `yaml:"permission-contents,omitempty"`
	PermissionDependabotSecrets  string `yaml:"permission-dependabot-secrets,omitempty"`
	PermissionDeployments        string `yaml:"permission-deployments,omitempty"`
	PermissionEnvironments       string `yaml:"permission-environments,omitempty"`
	PermissionIssues             string `yaml:"permission-issues,omitempty"`
	PermissionMetadata           string `yaml:"permission-metadata,omitempty"`
	PermissionPackages           string `yaml:"permission-packages,omitempty"`
	PermissionPages              string `yaml:"permission-pages,omitempty"`
	PermissionPullRequests       string `yaml:"permission-pull-requests,omitempty"`
	PermissionRepositoryHooks    string `yaml:"permission-repository-hooks,omitempty"`
	PermissionRepositoryProjects string `yaml:"permission-repository-projects,omitempty"`
	PermissionSecretScanningAlerts string `yaml:"permission-secret-scanning-alerts,omitempty"`
	PermissionSecrets            string `yaml:"permission-secrets,omitempty"`
	PermissionSecurityEvents     string `yaml:"permission-security-events,omitempty"`
	PermissionStatuses           string `yaml:"permission-statuses,omitempty"`
	PermissionVulnerabilityAlerts string `yaml:"permission-vulnerability-alerts,omitempty"`
	PermissionWorkflows          string `yaml:"permission-workflows,omitempty"`

	// Organization permissions
	PermissionMembers                 string `yaml:"permission-members,omitempty"`
	PermissionOrganizationAdministration string `yaml:"permission-organization-administration,omitempty"`
	PermissionOrganizationEvents      string `yaml:"permission-organization-events,omitempty"`
	PermissionOrganizationHooks       string `yaml:"permission-organization-hooks,omitempty"`
	PermissionOrganizationPackages    string `yaml:"permission-organization-packages,omitempty"`
	PermissionOrganizationPlan        string `yaml:"permission-organization-plan,omitempty"`
	PermissionOrganizationProjects    string `yaml:"permission-organization-projects,omitempty"`
	PermissionOrganizationSecrets     string `yaml:"permission-organization-secrets,omitempty"`
	PermissionOrganizationSelfHostedRunners string `yaml:"permission-organization-self-hosted-runners,omitempty"`
	PermissionOrganizationUserBlocking string `yaml:"permission-organization-user-blocking,omitempty"`
	PermissionTeamDiscussions         string `yaml:"permission-team-discussions,omitempty"`

	// User permissions
	PermissionEmailAddresses string `yaml:"permission-email-addresses,omitempty"`
	PermissionFollowers      string `yaml:"permission-followers,omitempty"`
	PermissionGitSSHKeys     string `yaml:"permission-git-ssh-keys,omitempty"`
	PermissionGPGKeys        string `yaml:"permission-gpg-keys,omitempty"`
	PermissionInteractionLimits string `yaml:"permission-interaction-limits,omitempty"`
	PermissionProfile        string `yaml:"permission-profile,omitempty"`
	PermissionStarring       string `yaml:"permission-starring,omitempty"`
}

// Action returns the action reference.
func (a CreateGithubAppToken) Action() string {
	return "actions/create-github-app-token@v1"
}

// Inputs returns the action inputs as a map.
func (a CreateGithubAppToken) Inputs() map[string]any {
	with := make(map[string]any)

	// Core inputs
	if a.AppID != "" {
		with["app-id"] = a.AppID
	}
	if a.PrivateKey != "" {
		with["private-key"] = a.PrivateKey
	}
	if a.Owner != "" {
		with["owner"] = a.Owner
	}
	if a.Repositories != "" {
		with["repositories"] = a.Repositories
	}
	if a.SkipTokenRevoke {
		with["skip-token-revoke"] = a.SkipTokenRevoke
	}
	if a.GithubAPIURL != "" {
		with["github-api-url"] = a.GithubAPIURL
	}

	// Repository permissions
	if a.PermissionActions != "" {
		with["permission-actions"] = a.PermissionActions
	}
	if a.PermissionAdministration != "" {
		with["permission-administration"] = a.PermissionAdministration
	}
	if a.PermissionChecks != "" {
		with["permission-checks"] = a.PermissionChecks
	}
	if a.PermissionCodespaces != "" {
		with["permission-codespaces"] = a.PermissionCodespaces
	}
	if a.PermissionContents != "" {
		with["permission-contents"] = a.PermissionContents
	}
	if a.PermissionDependabotSecrets != "" {
		with["permission-dependabot-secrets"] = a.PermissionDependabotSecrets
	}
	if a.PermissionDeployments != "" {
		with["permission-deployments"] = a.PermissionDeployments
	}
	if a.PermissionEnvironments != "" {
		with["permission-environments"] = a.PermissionEnvironments
	}
	if a.PermissionIssues != "" {
		with["permission-issues"] = a.PermissionIssues
	}
	if a.PermissionMetadata != "" {
		with["permission-metadata"] = a.PermissionMetadata
	}
	if a.PermissionPackages != "" {
		with["permission-packages"] = a.PermissionPackages
	}
	if a.PermissionPages != "" {
		with["permission-pages"] = a.PermissionPages
	}
	if a.PermissionPullRequests != "" {
		with["permission-pull-requests"] = a.PermissionPullRequests
	}
	if a.PermissionRepositoryHooks != "" {
		with["permission-repository-hooks"] = a.PermissionRepositoryHooks
	}
	if a.PermissionRepositoryProjects != "" {
		with["permission-repository-projects"] = a.PermissionRepositoryProjects
	}
	if a.PermissionSecretScanningAlerts != "" {
		with["permission-secret-scanning-alerts"] = a.PermissionSecretScanningAlerts
	}
	if a.PermissionSecrets != "" {
		with["permission-secrets"] = a.PermissionSecrets
	}
	if a.PermissionSecurityEvents != "" {
		with["permission-security-events"] = a.PermissionSecurityEvents
	}
	if a.PermissionStatuses != "" {
		with["permission-statuses"] = a.PermissionStatuses
	}
	if a.PermissionVulnerabilityAlerts != "" {
		with["permission-vulnerability-alerts"] = a.PermissionVulnerabilityAlerts
	}
	if a.PermissionWorkflows != "" {
		with["permission-workflows"] = a.PermissionWorkflows
	}

	// Organization permissions
	if a.PermissionMembers != "" {
		with["permission-members"] = a.PermissionMembers
	}
	if a.PermissionOrganizationAdministration != "" {
		with["permission-organization-administration"] = a.PermissionOrganizationAdministration
	}
	if a.PermissionOrganizationEvents != "" {
		with["permission-organization-events"] = a.PermissionOrganizationEvents
	}
	if a.PermissionOrganizationHooks != "" {
		with["permission-organization-hooks"] = a.PermissionOrganizationHooks
	}
	if a.PermissionOrganizationPackages != "" {
		with["permission-organization-packages"] = a.PermissionOrganizationPackages
	}
	if a.PermissionOrganizationPlan != "" {
		with["permission-organization-plan"] = a.PermissionOrganizationPlan
	}
	if a.PermissionOrganizationProjects != "" {
		with["permission-organization-projects"] = a.PermissionOrganizationProjects
	}
	if a.PermissionOrganizationSecrets != "" {
		with["permission-organization-secrets"] = a.PermissionOrganizationSecrets
	}
	if a.PermissionOrganizationSelfHostedRunners != "" {
		with["permission-organization-self-hosted-runners"] = a.PermissionOrganizationSelfHostedRunners
	}
	if a.PermissionOrganizationUserBlocking != "" {
		with["permission-organization-user-blocking"] = a.PermissionOrganizationUserBlocking
	}
	if a.PermissionTeamDiscussions != "" {
		with["permission-team-discussions"] = a.PermissionTeamDiscussions
	}

	// User permissions
	if a.PermissionEmailAddresses != "" {
		with["permission-email-addresses"] = a.PermissionEmailAddresses
	}
	if a.PermissionFollowers != "" {
		with["permission-followers"] = a.PermissionFollowers
	}
	if a.PermissionGitSSHKeys != "" {
		with["permission-git-ssh-keys"] = a.PermissionGitSSHKeys
	}
	if a.PermissionGPGKeys != "" {
		with["permission-gpg-keys"] = a.PermissionGPGKeys
	}
	if a.PermissionInteractionLimits != "" {
		with["permission-interaction-limits"] = a.PermissionInteractionLimits
	}
	if a.PermissionProfile != "" {
		with["permission-profile"] = a.PermissionProfile
	}
	if a.PermissionStarring != "" {
		with["permission-starring"] = a.PermissionStarring
	}

	return with
}
