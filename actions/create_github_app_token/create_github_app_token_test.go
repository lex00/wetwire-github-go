package create_github_app_token

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCreateGithubAppToken_Action(t *testing.T) {
	c := CreateGithubAppToken{}
	if got := c.Action(); got != "actions/create-github-app-token@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions/create-github-app-token@v1")
	}
}

func TestCreateGithubAppToken_Inputs(t *testing.T) {
	c := CreateGithubAppToken{
		AppID:      "12345",
		PrivateKey: "${{ secrets.APP_PRIVATE_KEY }}",
		Owner:      "my-org",
		Repositories: "repo1,repo2",
	}

	inputs := c.Inputs()

	if inputs["app-id"] != "12345" {
		t.Errorf("inputs[app-id] = %v, want %q", inputs["app-id"], "12345")
	}

	if inputs["private-key"] != "${{ secrets.APP_PRIVATE_KEY }}" {
		t.Errorf("inputs[private-key] = %v, want %q", inputs["private-key"], "${{ secrets.APP_PRIVATE_KEY }}")
	}

	if inputs["owner"] != "my-org" {
		t.Errorf("inputs[owner] = %v, want %q", inputs["owner"], "my-org")
	}

	if inputs["repositories"] != "repo1,repo2" {
		t.Errorf("inputs[repositories] = %v, want %q", inputs["repositories"], "repo1,repo2")
	}
}

func TestCreateGithubAppToken_Inputs_Empty(t *testing.T) {
	c := CreateGithubAppToken{}
	inputs := c.Inputs()

	// Empty struct should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty CreateGithubAppToken.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestCreateGithubAppToken_Inputs_RequiredOnly(t *testing.T) {
	c := CreateGithubAppToken{
		AppID:      "67890",
		PrivateKey: "${{ secrets.MY_APP_KEY }}",
	}

	inputs := c.Inputs()

	if len(inputs) != 2 {
		t.Errorf("inputs has %d entries, want 2", len(inputs))
	}

	if inputs["app-id"] != "67890" {
		t.Errorf("inputs[app-id] = %v, want %q", inputs["app-id"], "67890")
	}

	if inputs["private-key"] != "${{ secrets.MY_APP_KEY }}" {
		t.Errorf("inputs[private-key] = %v, want %q", inputs["private-key"], "${{ secrets.MY_APP_KEY }}")
	}
}

func TestCreateGithubAppToken_Inputs_SkipTokenRevoke(t *testing.T) {
	c := CreateGithubAppToken{
		SkipTokenRevoke: true,
	}

	inputs := c.Inputs()

	if inputs["skip-token-revoke"] != true {
		t.Errorf("inputs[skip-token-revoke] = %v, want true", inputs["skip-token-revoke"])
	}
}

func TestCreateGithubAppToken_Inputs_GithubAPIURL(t *testing.T) {
	c := CreateGithubAppToken{
		GithubAPIURL: "https://api.github.enterprise.com",
	}

	inputs := c.Inputs()

	if inputs["github-api-url"] != "https://api.github.enterprise.com" {
		t.Errorf("inputs[github-api-url] = %v, want %q", inputs["github-api-url"], "https://api.github.enterprise.com")
	}
}

func TestCreateGithubAppToken_Inputs_Permissions(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionContents:     "write",
		PermissionPullRequests: "write",
		PermissionIssues:       "read",
		PermissionActions:      "read",
	}

	inputs := c.Inputs()

	if inputs["permission-contents"] != "write" {
		t.Errorf("inputs[permission-contents] = %v, want %q", inputs["permission-contents"], "write")
	}

	if inputs["permission-pull-requests"] != "write" {
		t.Errorf("inputs[permission-pull-requests] = %v, want %q", inputs["permission-pull-requests"], "write")
	}

	if inputs["permission-issues"] != "read" {
		t.Errorf("inputs[permission-issues] = %v, want %q", inputs["permission-issues"], "read")
	}

	if inputs["permission-actions"] != "read" {
		t.Errorf("inputs[permission-actions] = %v, want %q", inputs["permission-actions"], "read")
	}
}

func TestCreateGithubAppToken_ImplementsStepAction(t *testing.T) {
	c := CreateGithubAppToken{}
	// Verify CreateGithubAppToken implements StepAction interface
	var _ workflow.StepAction = c
}

func TestCreateGithubAppToken_Inputs_MultipleRepositories(t *testing.T) {
	c := CreateGithubAppToken{
		Repositories: "repo1,repo2,repo3",
	}

	inputs := c.Inputs()

	if inputs["repositories"] != "repo1,repo2,repo3" {
		t.Errorf("inputs[repositories] = %v, want %q", inputs["repositories"], "repo1,repo2,repo3")
	}
}

func TestCreateGithubAppToken_Inputs_NewlineRepositories(t *testing.T) {
	c := CreateGithubAppToken{
		Repositories: "repo1\nrepo2\nrepo3",
	}

	inputs := c.Inputs()

	if inputs["repositories"] != "repo1\nrepo2\nrepo3" {
		t.Errorf("inputs[repositories] = %v, want %q", inputs["repositories"], "repo1\nrepo2\nrepo3")
	}
}

func TestCreateGithubAppToken_Inputs_AllCoreFields(t *testing.T) {
	c := CreateGithubAppToken{
		AppID:           "123",
		PrivateKey:      "${{ secrets.KEY }}",
		Owner:           "test-org",
		Repositories:    "repo1,repo2",
		SkipTokenRevoke: true,
		GithubAPIURL:    "https://api.github.com",
	}

	inputs := c.Inputs()

	expected := map[string]any{
		"app-id":            "123",
		"private-key":       "${{ secrets.KEY }}",
		"owner":             "test-org",
		"repositories":      "repo1,repo2",
		"skip-token-revoke": true,
		"github-api-url":    "https://api.github.com",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestCreateGithubAppToken_Inputs_CommonPermissions(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionActions:             "read",
		PermissionChecks:              "write",
		PermissionContents:            "write",
		PermissionDeployments:         "write",
		PermissionEnvironments:        "read",
		PermissionIssues:              "write",
		PermissionPackages:            "read",
		PermissionPages:               "write",
		PermissionPullRequests:        "write",
		PermissionRepositoryProjects:  "write",
		PermissionSecurityEvents:      "read",
		PermissionStatuses:            "write",
		PermissionWorkflows:           "write",
	}

	inputs := c.Inputs()

	expectedPermissions := map[string]string{
		"permission-actions":             "read",
		"permission-checks":              "write",
		"permission-contents":            "write",
		"permission-deployments":         "write",
		"permission-environments":        "read",
		"permission-issues":              "write",
		"permission-packages":            "read",
		"permission-pages":               "write",
		"permission-pull-requests":       "write",
		"permission-repository-projects": "write",
		"permission-security-events":     "read",
		"permission-statuses":            "write",
		"permission-workflows":           "write",
	}

	for key, want := range expectedPermissions {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestCreateGithubAppToken_Inputs_OrgPermissions(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionMembers:                 "read",
		PermissionOrganizationAdministration: "write",
		PermissionOrganizationHooks:       "write",
		PermissionOrganizationProjects:    "write",
		PermissionOrganizationSecrets:     "write",
	}

	inputs := c.Inputs()

	if inputs["permission-members"] != "read" {
		t.Errorf("inputs[permission-members] = %v, want %q", inputs["permission-members"], "read")
	}

	if inputs["permission-organization-administration"] != "write" {
		t.Errorf("inputs[permission-organization-administration] = %v, want %q", inputs["permission-organization-administration"], "write")
	}

	if inputs["permission-organization-hooks"] != "write" {
		t.Errorf("inputs[permission-organization-hooks] = %v, want %q", inputs["permission-organization-hooks"], "write")
	}

	if inputs["permission-organization-projects"] != "write" {
		t.Errorf("inputs[permission-organization-projects] = %v, want %q", inputs["permission-organization-projects"], "write")
	}

	if inputs["permission-organization-secrets"] != "write" {
		t.Errorf("inputs[permission-organization-secrets] = %v, want %q", inputs["permission-organization-secrets"], "write")
	}
}

func TestCreateGithubAppToken_Inputs_FalseBoolFields(t *testing.T) {
	// Test that false boolean values are not included in inputs
	c := CreateGithubAppToken{
		SkipTokenRevoke: false,
	}

	inputs := c.Inputs()

	// False bools should not be in the inputs map
	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestCreateGithubAppToken_Inputs_AdditionalPermissions(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionAdministration:      "write",
		PermissionCodespaces:          "write",
		PermissionDependabotSecrets:   "write",
		PermissionMetadata:            "read",
		PermissionSecrets:             "write",
		PermissionSecretScanningAlerts: "read",
		PermissionVulnerabilityAlerts: "read",
	}

	inputs := c.Inputs()

	if inputs["permission-administration"] != "write" {
		t.Errorf("inputs[permission-administration] = %v, want %q", inputs["permission-administration"], "write")
	}

	if inputs["permission-codespaces"] != "write" {
		t.Errorf("inputs[permission-codespaces] = %v, want %q", inputs["permission-codespaces"], "write")
	}

	if inputs["permission-dependabot-secrets"] != "write" {
		t.Errorf("inputs[permission-dependabot-secrets] = %v, want %q", inputs["permission-dependabot-secrets"], "write")
	}

	if inputs["permission-metadata"] != "read" {
		t.Errorf("inputs[permission-metadata] = %v, want %q", inputs["permission-metadata"], "read")
	}

	if inputs["permission-secrets"] != "write" {
		t.Errorf("inputs[permission-secrets] = %v, want %q", inputs["permission-secrets"], "write")
	}

	if inputs["permission-secret-scanning-alerts"] != "read" {
		t.Errorf("inputs[permission-secret-scanning-alerts] = %v, want %q", inputs["permission-secret-scanning-alerts"], "read")
	}

	if inputs["permission-vulnerability-alerts"] != "read" {
		t.Errorf("inputs[permission-vulnerability-alerts] = %v, want %q", inputs["permission-vulnerability-alerts"], "read")
	}
}

func TestCreateGithubAppToken_Inputs_UserPermissions(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionEmailAddresses:    "read",
		PermissionFollowers:         "write",
		PermissionGitSSHKeys:        "write",
		PermissionGPGKeys:           "write",
		PermissionInteractionLimits: "write",
		PermissionProfile:           "write",
		PermissionStarring:          "write",
	}

	inputs := c.Inputs()

	if inputs["permission-email-addresses"] != "read" {
		t.Errorf("inputs[permission-email-addresses] = %v, want %q", inputs["permission-email-addresses"], "read")
	}

	if inputs["permission-followers"] != "write" {
		t.Errorf("inputs[permission-followers] = %v, want %q", inputs["permission-followers"], "write")
	}

	if inputs["permission-git-ssh-keys"] != "write" {
		t.Errorf("inputs[permission-git-ssh-keys] = %v, want %q", inputs["permission-git-ssh-keys"], "write")
	}

	if inputs["permission-gpg-keys"] != "write" {
		t.Errorf("inputs[permission-gpg-keys] = %v, want %q", inputs["permission-gpg-keys"], "write")
	}

	if inputs["permission-interaction-limits"] != "write" {
		t.Errorf("inputs[permission-interaction-limits] = %v, want %q", inputs["permission-interaction-limits"], "write")
	}

	if inputs["permission-profile"] != "write" {
		t.Errorf("inputs[permission-profile] = %v, want %q", inputs["permission-profile"], "write")
	}

	if inputs["permission-starring"] != "write" {
		t.Errorf("inputs[permission-starring] = %v, want %q", inputs["permission-starring"], "write")
	}
}

func TestCreateGithubAppToken_Inputs_MoreOrgPermissions(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionOrganizationEvents:            "read",
		PermissionOrganizationPackages:          "write",
		PermissionOrganizationPlan:              "read",
		PermissionOrganizationSelfHostedRunners: "write",
		PermissionOrganizationUserBlocking:      "write",
		PermissionTeamDiscussions:               "write",
	}

	inputs := c.Inputs()

	if inputs["permission-organization-events"] != "read" {
		t.Errorf("inputs[permission-organization-events] = %v, want %q", inputs["permission-organization-events"], "read")
	}

	if inputs["permission-organization-packages"] != "write" {
		t.Errorf("inputs[permission-organization-packages] = %v, want %q", inputs["permission-organization-packages"], "write")
	}

	if inputs["permission-organization-plan"] != "read" {
		t.Errorf("inputs[permission-organization-plan] = %v, want %q", inputs["permission-organization-plan"], "read")
	}

	if inputs["permission-organization-self-hosted-runners"] != "write" {
		t.Errorf("inputs[permission-organization-self-hosted-runners] = %v, want %q", inputs["permission-organization-self-hosted-runners"], "write")
	}

	if inputs["permission-organization-user-blocking"] != "write" {
		t.Errorf("inputs[permission-organization-user-blocking] = %v, want %q", inputs["permission-organization-user-blocking"], "write")
	}

	if inputs["permission-team-discussions"] != "write" {
		t.Errorf("inputs[permission-team-discussions] = %v, want %q", inputs["permission-team-discussions"], "write")
	}
}

func TestCreateGithubAppToken_Inputs_RepositoryHooks(t *testing.T) {
	c := CreateGithubAppToken{
		PermissionRepositoryHooks: "write",
	}

	inputs := c.Inputs()

	if inputs["permission-repository-hooks"] != "write" {
		t.Errorf("inputs[permission-repository-hooks] = %v, want %q", inputs["permission-repository-hooks"], "write")
	}
}
