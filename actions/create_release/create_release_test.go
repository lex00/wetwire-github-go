package create_release

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestCreateRelease_Action(t *testing.T) {
	a := CreateRelease{}
	if got := a.Action(); got != "actions/create-release@v1" {
		t.Errorf("Action() = %q, want %q", got, "actions/create-release@v1")
	}
}

func TestCreateRelease_Inputs_Empty(t *testing.T) {
	a := CreateRelease{}
	inputs := a.Inputs()

	if a.Action() != "actions/create-release@v1" {
		t.Errorf("Action() = %q, want %q", a.Action(), "actions/create-release@v1")
	}

	// Empty inputs should not include optional fields
	if _, ok := inputs["body"]; ok {
		t.Error("Empty body should not be in inputs")
	}
	if _, ok := inputs["body_path"]; ok {
		t.Error("Empty body_path should not be in inputs")
	}
	if _, ok := inputs["commitish"]; ok {
		t.Error("Empty commitish should not be in inputs")
	}
	if _, ok := inputs["owner"]; ok {
		t.Error("Empty owner should not be in inputs")
	}
	if _, ok := inputs["repo"]; ok {
		t.Error("Empty repo should not be in inputs")
	}
}

func TestCreateRelease_Inputs_RequiredOnly(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
	}

	inputs := a.Inputs()

	if inputs["tag_name"] != "v1.0.0" {
		t.Errorf("tag_name = %v, want %q", inputs["tag_name"], "v1.0.0")
	}
	if inputs["release_name"] != "Release 1.0.0" {
		t.Errorf("release_name = %v, want %q", inputs["release_name"], "Release 1.0.0")
	}

	// Optional fields should not be present
	if _, ok := inputs["body"]; ok {
		t.Error("Empty body should not be in inputs")
	}
}

func TestCreateRelease_Inputs_WithBody(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Body:        "## Release Notes\n\n- Feature 1\n- Fix 1",
	}

	inputs := a.Inputs()

	if inputs["body"] != "## Release Notes\n\n- Feature 1\n- Fix 1" {
		t.Errorf("body = %v, want expected value", inputs["body"])
	}
	if inputs["tag_name"] != "v1.0.0" {
		t.Errorf("tag_name = %v, want %q", inputs["tag_name"], "v1.0.0")
	}
	if inputs["release_name"] != "Release 1.0.0" {
		t.Errorf("release_name = %v, want %q", inputs["release_name"], "Release 1.0.0")
	}
}

func TestCreateRelease_Inputs_WithBodyPath(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		BodyPath:    "./CHANGELOG.md",
	}

	inputs := a.Inputs()

	if inputs["body_path"] != "./CHANGELOG.md" {
		t.Errorf("body_path = %v, want %q", inputs["body_path"], "./CHANGELOG.md")
	}
}

func TestCreateRelease_Inputs_Draft(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Draft:       true,
	}

	inputs := a.Inputs()

	if inputs["draft"] != true {
		t.Errorf("draft = %v, want %v", inputs["draft"], true)
	}
}

func TestCreateRelease_Inputs_Prerelease(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Prerelease:  true,
	}

	inputs := a.Inputs()

	if inputs["prerelease"] != true {
		t.Errorf("prerelease = %v, want %v", inputs["prerelease"], true)
	}
}

func TestCreateRelease_Inputs_Commitish(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Commitish:   "main",
	}

	inputs := a.Inputs()

	if inputs["commitish"] != "main" {
		t.Errorf("commitish = %v, want %q", inputs["commitish"], "main")
	}
}

func TestCreateRelease_Inputs_OwnerRepo(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Owner:       "octocat",
		Repo:        "hello-world",
	}

	inputs := a.Inputs()

	if inputs["owner"] != "octocat" {
		t.Errorf("owner = %v, want %q", inputs["owner"], "octocat")
	}
	if inputs["repo"] != "hello-world" {
		t.Errorf("repo = %v, want %q", inputs["repo"], "hello-world")
	}
}

func TestCreateRelease_Inputs_AllFields(t *testing.T) {
	a := CreateRelease{
		TagName:     "v2.0.0",
		ReleaseName: "Major Release 2.0.0",
		Body:        "Release body content",
		BodyPath:    "./RELEASE_NOTES.md",
		Draft:       true,
		Prerelease:  true,
		Commitish:   "develop",
		Owner:       "myorg",
		Repo:        "myrepo",
	}

	inputs := a.Inputs()

	if inputs["tag_name"] != "v2.0.0" {
		t.Errorf("tag_name = %v, want %q", inputs["tag_name"], "v2.0.0")
	}
	if inputs["release_name"] != "Major Release 2.0.0" {
		t.Errorf("release_name = %v, want %q", inputs["release_name"], "Major Release 2.0.0")
	}
	if inputs["body"] != "Release body content" {
		t.Errorf("body = %v, want %q", inputs["body"], "Release body content")
	}
	if inputs["body_path"] != "./RELEASE_NOTES.md" {
		t.Errorf("body_path = %v, want %q", inputs["body_path"], "./RELEASE_NOTES.md")
	}
	if inputs["draft"] != true {
		t.Errorf("draft = %v, want %v", inputs["draft"], true)
	}
	if inputs["prerelease"] != true {
		t.Errorf("prerelease = %v, want %v", inputs["prerelease"], true)
	}
	if inputs["commitish"] != "develop" {
		t.Errorf("commitish = %v, want %q", inputs["commitish"], "develop")
	}
	if inputs["owner"] != "myorg" {
		t.Errorf("owner = %v, want %q", inputs["owner"], "myorg")
	}
	if inputs["repo"] != "myrepo" {
		t.Errorf("repo = %v, want %q", inputs["repo"], "myrepo")
	}
}

func TestCreateRelease_Inputs_BooleanDefaults(t *testing.T) {
	a := CreateRelease{
		TagName:     "v1.0.0",
		ReleaseName: "Release 1.0.0",
		Draft:       false,
		Prerelease:  false,
	}

	inputs := a.Inputs()

	// False boolean values should not be included (omitempty behavior)
	if _, ok := inputs["draft"]; ok {
		t.Error("False draft should not be in inputs")
	}
	if _, ok := inputs["prerelease"]; ok {
		t.Error("False prerelease should not be in inputs")
	}
}

func TestCreateRelease_ImplementsStepAction(t *testing.T) {
	var _ workflow.StepAction = CreateRelease{}
}
