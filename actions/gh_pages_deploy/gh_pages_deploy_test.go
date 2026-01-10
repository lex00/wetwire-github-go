package gh_pages_deploy

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGitHubPagesDeploy_Action(t *testing.T) {
	g := GitHubPagesDeploy{}
	if got := g.Action(); got != "JamesIves/github-pages-deploy-action@v4" {
		t.Errorf("Action() = %q, want %q", got, "JamesIves/github-pages-deploy-action@v4")
	}
}

func TestGitHubPagesDeploy_Inputs(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:        "build",
		Branch:        "main",
		CommitMessage: "Deploy to GitHub Pages",
	}

	inputs := g.Inputs()

	if inputs["folder"] != "build" {
		t.Errorf("inputs[folder] = %v, want %q", inputs["folder"], "build")
	}

	if inputs["branch"] != "main" {
		t.Errorf("inputs[branch] = %v, want %q", inputs["branch"], "main")
	}

	if inputs["commit-message"] != "Deploy to GitHub Pages" {
		t.Errorf("inputs[commit-message] = %v, want %q", inputs["commit-message"], "Deploy to GitHub Pages")
	}
}

func TestGitHubPagesDeploy_Inputs_Empty(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder: ".", // Only required field
	}
	inputs := g.Inputs()

	// Should only have folder
	if len(inputs) != 1 {
		t.Errorf("GitHubPagesDeploy.Inputs() has %d entries, want 1", len(inputs))
	}

	if inputs["folder"] != "." {
		t.Errorf("inputs[folder] = %v, want %q", inputs["folder"], ".")
	}
}

func TestGitHubPagesDeploy_Inputs_SSHKey(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder: "build",
		SSHKey: "ssh-rsa AAAAB3...",
	}

	inputs := g.Inputs()

	if inputs["ssh-key"] != "ssh-rsa AAAAB3..." {
		t.Errorf("inputs[ssh-key] = %v, want %q", inputs["ssh-key"], "ssh-rsa AAAAB3...")
	}
}

func TestGitHubPagesDeploy_Inputs_Token(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder: "dist",
		Token:  "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := g.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestGitHubPagesDeploy_Inputs_TargetFolder(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:       "build",
		TargetFolder: "docs",
	}

	inputs := g.Inputs()

	if inputs["target-folder"] != "docs" {
		t.Errorf("inputs[target-folder] = %v, want %q", inputs["target-folder"], "docs")
	}
}

func TestGitHubPagesDeploy_Inputs_BoolFields(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:       "build",
		Clean:        true,
		DryRun:       true,
		Force:        true,
		SingleCommit: true,
		Silent:       true,
	}

	inputs := g.Inputs()

	if inputs["clean"] != true {
		t.Errorf("inputs[clean] = %v, want true", inputs["clean"])
	}

	if inputs["dry-run"] != true {
		t.Errorf("inputs[dry-run] = %v, want true", inputs["dry-run"])
	}

	if inputs["force"] != true {
		t.Errorf("inputs[force] = %v, want true", inputs["force"])
	}

	if inputs["single-commit"] != true {
		t.Errorf("inputs[single-commit] = %v, want true", inputs["single-commit"])
	}

	if inputs["silent"] != true {
		t.Errorf("inputs[silent] = %v, want true", inputs["silent"])
	}
}

func TestGitHubPagesDeploy_Inputs_FalseBoolFields(t *testing.T) {
	// Test that false boolean values are not included in inputs
	g := GitHubPagesDeploy{
		Folder:       "build",
		Clean:        false,
		DryRun:       false,
		Force:        false,
		SingleCommit: false,
		Silent:       false,
	}

	inputs := g.Inputs()

	// Only folder should be in the inputs map
	if len(inputs) != 1 {
		t.Errorf("inputs for false bools has %d entries, want 1. Got: %v", len(inputs), inputs)
	}

	if _, exists := inputs["clean"]; exists {
		t.Errorf("inputs[clean] should not exist for false value")
	}
}

func TestGitHubPagesDeploy_Inputs_CleanExclude(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:       "build",
		CleanExclude: "*.json\n*.txt",
	}

	inputs := g.Inputs()

	if inputs["clean-exclude"] != "*.json\n*.txt" {
		t.Errorf("inputs[clean-exclude] = %v, want %q", inputs["clean-exclude"], "*.json\n*.txt")
	}
}

func TestGitHubPagesDeploy_Inputs_GitConfig(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:          "build",
		GitConfigName:   "GitHub Actions Bot",
		GitConfigEmail:  "actions@github.com",
	}

	inputs := g.Inputs()

	if inputs["git-config-name"] != "GitHub Actions Bot" {
		t.Errorf("inputs[git-config-name] = %v, want %q", inputs["git-config-name"], "GitHub Actions Bot")
	}

	if inputs["git-config-email"] != "actions@github.com" {
		t.Errorf("inputs[git-config-email] = %v, want %q", inputs["git-config-email"], "actions@github.com")
	}
}

func TestGitHubPagesDeploy_Inputs_RepositoryName(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:         "build",
		RepositoryName: "owner/repo",
	}

	inputs := g.Inputs()

	if inputs["repository-name"] != "owner/repo" {
		t.Errorf("inputs[repository-name] = %v, want %q", inputs["repository-name"], "owner/repo")
	}
}

func TestGitHubPagesDeploy_Inputs_Tag(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder: "build",
		Tag:    "v1.0.0",
	}

	inputs := g.Inputs()

	if inputs["tag"] != "v1.0.0" {
		t.Errorf("inputs[tag] = %v, want %q", inputs["tag"], "v1.0.0")
	}
}

func TestGitHubPagesDeploy_Inputs_AttemptLimit(t *testing.T) {
	g := GitHubPagesDeploy{
		Folder:       "build",
		AttemptLimit: 5,
	}

	inputs := g.Inputs()

	if inputs["attempt-limit"] != 5 {
		t.Errorf("inputs[attempt-limit] = %v, want 5", inputs["attempt-limit"])
	}
}

func TestGitHubPagesDeploy_Inputs_ZeroAttemptLimit(t *testing.T) {
	// Test that AttemptLimit = 0 is not included
	g := GitHubPagesDeploy{
		Folder:       "build",
		AttemptLimit: 0,
	}

	inputs := g.Inputs()

	if _, exists := inputs["attempt-limit"]; exists {
		t.Errorf("inputs[attempt-limit] should not exist for AttemptLimit=0")
	}
}

func TestGitHubPagesDeploy_Inputs_AllFields(t *testing.T) {
	g := GitHubPagesDeploy{
		SSHKey:          "ssh-key-value",
		Token:           "token-value",
		Branch:          "gh-pages",
		Folder:          "build",
		TargetFolder:    "docs",
		CommitMessage:   "Deploy",
		Clean:           true,
		CleanExclude:    "*.keep",
		DryRun:          true,
		Force:           true,
		GitConfigName:   "Bot",
		GitConfigEmail:  "bot@example.com",
		RepositoryName:  "owner/repo",
		Tag:             "v1.0",
		SingleCommit:    true,
		Silent:          true,
		AttemptLimit:    10,
	}

	inputs := g.Inputs()

	expected := map[string]any{
		"ssh-key":          "ssh-key-value",
		"token":            "token-value",
		"branch":           "gh-pages",
		"folder":           "build",
		"target-folder":    "docs",
		"commit-message":   "Deploy",
		"clean":            true,
		"clean-exclude":    "*.keep",
		"dry-run":          true,
		"force":            true,
		"git-config-name":  "Bot",
		"git-config-email": "bot@example.com",
		"repository-name":  "owner/repo",
		"tag":              "v1.0",
		"single-commit":    true,
		"silent":           true,
		"attempt-limit":    10,
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

func TestGitHubPagesDeploy_ImplementsStepAction(t *testing.T) {
	g := GitHubPagesDeploy{}
	// Verify GitHubPagesDeploy implements StepAction interface
	var _ workflow.StepAction = g
}
