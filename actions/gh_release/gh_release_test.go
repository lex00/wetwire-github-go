package gh_release

import (
	"testing"
)

func TestGHRelease_Action(t *testing.T) {
	a := GHRelease{}
	if got := a.Action(); got != "softprops/action-gh-release@v2" {
		t.Errorf("Action() = %q, want %q", got, "softprops/action-gh-release@v2")
	}
}

func TestGHRelease_ToStep(t *testing.T) {
	a := GHRelease{
		Files: "*.tar.gz",
	}

	step := a.ToStep()

	if step.Uses != "softprops/action-gh-release@v2" {
		t.Errorf("Uses = %q, want %q", step.Uses, "softprops/action-gh-release@v2")
	}

	if step.With["files"] != "*.tar.gz" {
		t.Errorf("With[files] = %v, want %q", step.With["files"], "*.tar.gz")
	}
}

func TestGHRelease_ToStep_Empty(t *testing.T) {
	a := GHRelease{}
	step := a.ToStep()

	if step.Uses != "softprops/action-gh-release@v2" {
		t.Errorf("Uses = %q, want %q", step.Uses, "softprops/action-gh-release@v2")
	}

	if _, ok := step.With["files"]; ok {
		t.Error("Empty files should not be in With")
	}
}

func TestGHRelease_ToStep_WithBody(t *testing.T) {
	a := GHRelease{
		Body:    "## Release Notes\n\n- Feature 1\n- Fix 1",
		TagName: "v1.0.0",
		Name:    "Release 1.0.0",
	}

	step := a.ToStep()

	if step.With["body"] != "## Release Notes\n\n- Feature 1\n- Fix 1" {
		t.Errorf("body = %v, want expected value", step.With["body"])
	}
	if step.With["tag_name"] != "v1.0.0" {
		t.Errorf("tag_name = %v, want %q", step.With["tag_name"], "v1.0.0")
	}
	if step.With["name"] != "Release 1.0.0" {
		t.Errorf("name = %v, want %q", step.With["name"], "Release 1.0.0")
	}
}

func TestGHRelease_ToStep_Draft(t *testing.T) {
	a := GHRelease{
		Draft: true,
	}

	step := a.ToStep()

	if step.With["draft"] != true {
		t.Errorf("draft = %v, want %v", step.With["draft"], true)
	}
}

func TestGHRelease_ToStep_Prerelease(t *testing.T) {
	a := GHRelease{
		Prerelease: true,
	}

	step := a.ToStep()

	if step.With["prerelease"] != true {
		t.Errorf("prerelease = %v, want %v", step.With["prerelease"], true)
	}
}

func TestGHRelease_ToStep_GenerateNotes(t *testing.T) {
	a := GHRelease{
		GenerateReleaseNotes: true,
	}

	step := a.ToStep()

	if step.With["generate_release_notes"] != true {
		t.Errorf("generate_release_notes = %v, want %v", step.With["generate_release_notes"], true)
	}
}

func TestGHRelease_ToStep_AllFields(t *testing.T) {
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

	step := a.ToStep()

	if step.With["body"] != "Release body" {
		t.Errorf("body = %v, want %q", step.With["body"], "Release body")
	}
	if step.With["name"] != "v2.0.0" {
		t.Errorf("name = %v, want %q", step.With["name"], "v2.0.0")
	}
	if step.With["tag_name"] != "v2.0.0" {
		t.Errorf("tag_name = %v, want %q", step.With["tag_name"], "v2.0.0")
	}
	if step.With["target_commitish"] != "main" {
		t.Errorf("target_commitish = %v, want %q", step.With["target_commitish"], "main")
	}
	if step.With["draft"] != true {
		t.Errorf("draft = %v, want %v", step.With["draft"], true)
	}
	if step.With["prerelease"] != true {
		t.Errorf("prerelease = %v, want %v", step.With["prerelease"], true)
	}
	if step.With["generate_release_notes"] != true {
		t.Errorf("generate_release_notes = %v, want %v", step.With["generate_release_notes"], true)
	}
	if step.With["fail_on_unmatched_files"] != true {
		t.Errorf("fail_on_unmatched_files = %v, want %v", step.With["fail_on_unmatched_files"], true)
	}
	if step.With["append_body"] != true {
		t.Errorf("append_body = %v, want %v", step.With["append_body"], true)
	}
	if step.With["make_latest"] != "true" {
		t.Errorf("make_latest = %v, want %q", step.With["make_latest"], "true")
	}
	if step.With["discussion_category_name"] != "Releases" {
		t.Errorf("discussion_category_name = %v, want %q", step.With["discussion_category_name"], "Releases")
	}
}

func TestGHRelease_ToStep_MultipleFiles(t *testing.T) {
	a := GHRelease{
		Files: "dist/*.tar.gz\ndist/*.zip\nbin/*",
	}

	step := a.ToStep()

	if step.With["files"] != "dist/*.tar.gz\ndist/*.zip\nbin/*" {
		t.Errorf("files = %v, want expected value", step.With["files"])
	}
}
