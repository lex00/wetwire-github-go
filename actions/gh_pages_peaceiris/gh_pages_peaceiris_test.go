package gh_pages_peaceiris

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGHPagesPeaceiris_Action(t *testing.T) {
	g := GHPagesPeaceiris{}
	if got := g.Action(); got != "peaceiris/actions-gh-pages@v4" {
		t.Errorf("Action() = %q, want %q", got, "peaceiris/actions-gh-pages@v4")
	}
}

func TestGHPagesPeaceiris_Inputs(t *testing.T) {
	g := GHPagesPeaceiris{
		GithubToken:   "${{ secrets.GITHUB_TOKEN }}",
		PublishBranch: "main",
		PublishDir:    "./dist",
		CNAME:         "example.com",
	}

	inputs := g.Inputs()

	if inputs["github_token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[github_token] = %v, want %q", inputs["github_token"], "${{ secrets.GITHUB_TOKEN }}")
	}

	if inputs["publish_branch"] != "main" {
		t.Errorf("inputs[publish_branch] = %v, want %q", inputs["publish_branch"], "main")
	}

	if inputs["publish_dir"] != "./dist" {
		t.Errorf("inputs[publish_dir] = %v, want %q", inputs["publish_dir"], "./dist")
	}

	if inputs["cname"] != "example.com" {
		t.Errorf("inputs[cname] = %v, want %q", inputs["cname"], "example.com")
	}
}

func TestGHPagesPeaceiris_Inputs_Empty(t *testing.T) {
	g := GHPagesPeaceiris{}
	inputs := g.Inputs()

	// Empty struct should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty GHPagesPeaceiris.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestGHPagesPeaceiris_Inputs_BoolFields(t *testing.T) {
	g := GHPagesPeaceiris{
		AllowEmptyCommit: true,
		KeepFiles:        true,
		ForceOrphan:      true,
		EnableJekyll:     true,
		DisableNojekyll:  true,
	}

	inputs := g.Inputs()

	if inputs["allow_empty_commit"] != true {
		t.Errorf("inputs[allow_empty_commit] = %v, want true", inputs["allow_empty_commit"])
	}

	if inputs["keep_files"] != true {
		t.Errorf("inputs[keep_files] = %v, want true", inputs["keep_files"])
	}

	if inputs["force_orphan"] != true {
		t.Errorf("inputs[force_orphan] = %v, want true", inputs["force_orphan"])
	}

	if inputs["enable_jekyll"] != true {
		t.Errorf("inputs[enable_jekyll] = %v, want true", inputs["enable_jekyll"])
	}

	if inputs["disable_nojekyll"] != true {
		t.Errorf("inputs[disable_nojekyll] = %v, want true", inputs["disable_nojekyll"])
	}
}

func TestGHPagesPeaceiris_Inputs_DeployKey(t *testing.T) {
	g := GHPagesPeaceiris{
		DeployKey: "${{ secrets.DEPLOY_KEY }}",
	}

	inputs := g.Inputs()

	if inputs["deploy_key"] != "${{ secrets.DEPLOY_KEY }}" {
		t.Errorf("inputs[deploy_key] = %v, want %q", inputs["deploy_key"], "${{ secrets.DEPLOY_KEY }}")
	}
}

func TestGHPagesPeaceiris_Inputs_PersonalToken(t *testing.T) {
	g := GHPagesPeaceiris{
		PersonalToken: "${{ secrets.PAT }}",
	}

	inputs := g.Inputs()

	if inputs["personal_token"] != "${{ secrets.PAT }}" {
		t.Errorf("inputs[personal_token] = %v, want %q", inputs["personal_token"], "${{ secrets.PAT }}")
	}
}

func TestGHPagesPeaceiris_Inputs_DestinationDir(t *testing.T) {
	g := GHPagesPeaceiris{
		DestinationDir: "docs",
	}

	inputs := g.Inputs()

	if inputs["destination_dir"] != "docs" {
		t.Errorf("inputs[destination_dir] = %v, want %q", inputs["destination_dir"], "docs")
	}
}

func TestGHPagesPeaceiris_Inputs_ExternalRepository(t *testing.T) {
	g := GHPagesPeaceiris{
		ExternalRepository: "owner/repo",
	}

	inputs := g.Inputs()

	if inputs["external_repository"] != "owner/repo" {
		t.Errorf("inputs[external_repository] = %v, want %q", inputs["external_repository"], "owner/repo")
	}
}

func TestGHPagesPeaceiris_Inputs_GitConfig(t *testing.T) {
	g := GHPagesPeaceiris{
		UserName:  "GitHub Actions Bot",
		UserEmail: "actions@github.com",
	}

	inputs := g.Inputs()

	if inputs["user_name"] != "GitHub Actions Bot" {
		t.Errorf("inputs[user_name] = %v, want %q", inputs["user_name"], "GitHub Actions Bot")
	}

	if inputs["user_email"] != "actions@github.com" {
		t.Errorf("inputs[user_email] = %v, want %q", inputs["user_email"], "actions@github.com")
	}
}

func TestGHPagesPeaceiris_Inputs_CommitMessage(t *testing.T) {
	g := GHPagesPeaceiris{
		CommitMessage:     "Deploy to GitHub Pages",
		FullCommitMessage: "Full deployment message",
	}

	inputs := g.Inputs()

	if inputs["commit_message"] != "Deploy to GitHub Pages" {
		t.Errorf("inputs[commit_message] = %v, want %q", inputs["commit_message"], "Deploy to GitHub Pages")
	}

	if inputs["full_commit_message"] != "Full deployment message" {
		t.Errorf("inputs[full_commit_message] = %v, want %q", inputs["full_commit_message"], "Full deployment message")
	}
}

func TestGHPagesPeaceiris_Inputs_TagInfo(t *testing.T) {
	g := GHPagesPeaceiris{
		TagName:    "v1.0.0",
		TagMessage: "Release v1.0.0",
	}

	inputs := g.Inputs()

	if inputs["tag_name"] != "v1.0.0" {
		t.Errorf("inputs[tag_name] = %v, want %q", inputs["tag_name"], "v1.0.0")
	}

	if inputs["tag_message"] != "Release v1.0.0" {
		t.Errorf("inputs[tag_message] = %v, want %q", inputs["tag_message"], "Release v1.0.0")
	}
}

func TestGHPagesPeaceiris_Inputs_ExcludeAssets(t *testing.T) {
	g := GHPagesPeaceiris{
		ExcludeAssets: ".github,.gitignore",
	}

	inputs := g.Inputs()

	if inputs["exclude_assets"] != ".github,.gitignore" {
		t.Errorf("inputs[exclude_assets] = %v, want %q", inputs["exclude_assets"], ".github,.gitignore")
	}
}

func TestGHPagesPeaceiris_ImplementsStepAction(t *testing.T) {
	g := GHPagesPeaceiris{}
	// Verify GHPagesPeaceiris implements StepAction interface
	var _ workflow.StepAction = g
}

func TestGHPagesPeaceiris_Inputs_AllFields(t *testing.T) {
	g := GHPagesPeaceiris{
		DeployKey:          "${{ secrets.DEPLOY_KEY }}",
		GithubToken:        "${{ secrets.GITHUB_TOKEN }}",
		PersonalToken:      "${{ secrets.PAT }}",
		PublishBranch:      "gh-pages",
		PublishDir:         "./build",
		DestinationDir:     "docs",
		ExternalRepository: "org/repo",
		AllowEmptyCommit:   true,
		KeepFiles:          true,
		ForceOrphan:        true,
		UserName:           "Bot",
		UserEmail:          "bot@example.com",
		CommitMessage:      "Deploy",
		FullCommitMessage:  "Full deploy",
		TagName:            "v1.0.0",
		TagMessage:         "Release",
		EnableJekyll:       true,
		DisableNojekyll:    true,
		CNAME:              "example.com",
		ExcludeAssets:      ".git",
	}

	inputs := g.Inputs()

	expected := map[string]any{
		"deploy_key":          "${{ secrets.DEPLOY_KEY }}",
		"github_token":        "${{ secrets.GITHUB_TOKEN }}",
		"personal_token":      "${{ secrets.PAT }}",
		"publish_branch":      "gh-pages",
		"publish_dir":         "./build",
		"destination_dir":     "docs",
		"external_repository": "org/repo",
		"allow_empty_commit":  true,
		"keep_files":          true,
		"force_orphan":        true,
		"user_name":           "Bot",
		"user_email":          "bot@example.com",
		"commit_message":      "Deploy",
		"full_commit_message": "Full deploy",
		"tag_name":            "v1.0.0",
		"tag_message":         "Release",
		"enable_jekyll":       true,
		"disable_nojekyll":    true,
		"cname":               "example.com",
		"exclude_assets":      ".git",
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

func TestGHPagesPeaceiris_Inputs_FalseBoolFields(t *testing.T) {
	// Test that false boolean values are not included in inputs
	g := GHPagesPeaceiris{
		AllowEmptyCommit: false,
		KeepFiles:        false,
		ForceOrphan:      false,
		EnableJekyll:     false,
		DisableNojekyll:  false,
	}

	inputs := g.Inputs()

	// None of these should be in the inputs map
	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}
