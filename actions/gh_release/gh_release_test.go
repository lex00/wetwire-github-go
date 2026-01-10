package gh_release

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestGHRelease_Action(t *testing.T) {
	a := GHRelease{}
	if got := a.Action(); got != "softprops/action-gh-release@v2" {
		t.Errorf("Action() = %q, want %q", got, "softprops/action-gh-release@v2")
	}
}

func TestGHRelease_Inputs(t *testing.T) {
	a := GHRelease{
		Files: "*.tar.gz",
	}

	inputs := a.Inputs()

	if a.Action() != "softprops/action-gh-release@v2" {
		t.Errorf("Action() = %q, want %q", a.Action(), "softprops/action-gh-release@v2")
	}

	if inputs["files"] != "*.tar.gz" {
		t.Errorf("inputs[files] = %v, want %q", inputs["files"], "*.tar.gz")
	}
}

func TestGHRelease_Inputs_Empty(t *testing.T) {
	a := GHRelease{}
	inputs := a.Inputs()

	if a.Action() != "softprops/action-gh-release@v2" {
		t.Errorf("Action() = %q, want %q", a.Action(), "softprops/action-gh-release@v2")
	}

	if _, ok := inputs["files"]; ok {
		t.Error("Empty files should not be in inputs")
	}
}

func TestGHRelease_Inputs_WithBody(t *testing.T) {
	a := GHRelease{
		Body:    "## Release Notes\n\n- Feature 1\n- Fix 1",
		TagName: "v1.0.0",
		Name:    "Release 1.0.0",
	}

	inputs := a.Inputs()

	if inputs["body"] != "## Release Notes\n\n- Feature 1\n- Fix 1" {
		t.Errorf("body = %v, want expected value", inputs["body"])
	}
	if inputs["tag_name"] != "v1.0.0" {
		t.Errorf("tag_name = %v, want %q", inputs["tag_name"], "v1.0.0")
	}
	if inputs["name"] != "Release 1.0.0" {
		t.Errorf("name = %v, want %q", inputs["name"], "Release 1.0.0")
	}
}

func TestGHRelease_Inputs_Draft(t *testing.T) {
	a := GHRelease{
		Draft: true,
	}

	inputs := a.Inputs()

	if inputs["draft"] != true {
		t.Errorf("draft = %v, want %v", inputs["draft"], true)
	}
}

func TestGHRelease_Inputs_Prerelease(t *testing.T) {
	a := GHRelease{
		Prerelease: true,
	}

	inputs := a.Inputs()

	if inputs["prerelease"] != true {
		t.Errorf("prerelease = %v, want %v", inputs["prerelease"], true)
	}
}

func TestGHRelease_Inputs_GenerateNotes(t *testing.T) {
	a := GHRelease{
		GenerateReleaseNotes: true,
	}

	inputs := a.Inputs()

	if inputs["generate_release_notes"] != true {
		t.Errorf("generate_release_notes = %v, want %v", inputs["generate_release_notes"], true)
	}
}

func TestGHRelease_Inputs_AllFields(t *testing.T) {
	a := GHRelease{
		Body:                   "Release body",
		Name:                   "v2.0.0",
		TagName:                "v2.0.0",
		TargetCommitish:        "main",
		Draft:                  true,
		Prerelease:             true,
		GenerateReleaseNotes:   true,
		Files:                  "dist/*\nbin/*",
		FailOnUnmatchedFiles:   true,
		Token:                  "${{ secrets.GITHUB_TOKEN }}",
		Repository:             "owner/repo",
		AppendBody:             true,
		MakeLatest:             "true",
		DiscussionCategoryName: "Releases",
	}

	inputs := a.Inputs()

	if inputs["body"] != "Release body" {
		t.Errorf("body = %v, want %q", inputs["body"], "Release body")
	}
	if inputs["name"] != "v2.0.0" {
		t.Errorf("name = %v, want %q", inputs["name"], "v2.0.0")
	}
	if inputs["tag_name"] != "v2.0.0" {
		t.Errorf("tag_name = %v, want %q", inputs["tag_name"], "v2.0.0")
	}
	if inputs["target_commitish"] != "main" {
		t.Errorf("target_commitish = %v, want %q", inputs["target_commitish"], "main")
	}
	if inputs["draft"] != true {
		t.Errorf("draft = %v, want %v", inputs["draft"], true)
	}
	if inputs["prerelease"] != true {
		t.Errorf("prerelease = %v, want %v", inputs["prerelease"], true)
	}
	if inputs["generate_release_notes"] != true {
		t.Errorf("generate_release_notes = %v, want %v", inputs["generate_release_notes"], true)
	}
	if inputs["fail_on_unmatched_files"] != true {
		t.Errorf("fail_on_unmatched_files = %v, want %v", inputs["fail_on_unmatched_files"], true)
	}
	if inputs["append_body"] != true {
		t.Errorf("append_body = %v, want %v", inputs["append_body"], true)
	}
	if inputs["make_latest"] != "true" {
		t.Errorf("make_latest = %v, want %q", inputs["make_latest"], "true")
	}
	if inputs["discussion_category_name"] != "Releases" {
		t.Errorf("discussion_category_name = %v, want %q", inputs["discussion_category_name"], "Releases")
	}
}

func TestGHRelease_Inputs_MultipleFiles(t *testing.T) {
	a := GHRelease{
		Files: "dist/*.tar.gz\ndist/*.zip\nbin/*",
	}

	inputs := a.Inputs()

	if inputs["files"] != "dist/*.tar.gz\ndist/*.zip\nbin/*" {
		t.Errorf("files = %v, want expected value", inputs["files"])
	}
}

func TestGHRelease_Inputs_BodyPath(t *testing.T) {
	a := GHRelease{
		BodyPath: "./CHANGELOG.md",
	}

	inputs := a.Inputs()

	if inputs["body_path"] != "./CHANGELOG.md" {
		t.Errorf("inputs[body_path] = %v, want %q", inputs["body_path"], "./CHANGELOG.md")
	}
}

func TestGHRelease_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = GHRelease{}
}
