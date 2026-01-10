package ncipollo_release

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestNcipolloRelease_Action(t *testing.T) {
	a := NcipolloRelease{}
	if got := a.Action(); got != "ncipollo/release-action@v1" {
		t.Errorf("Action() = %q, want %q", got, "ncipollo/release-action@v1")
	}
}

func TestNcipolloRelease_Inputs_Empty(t *testing.T) {
	a := NcipolloRelease{}
	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("empty NcipolloRelease.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestNcipolloRelease_ImplementsStepAction(t *testing.T) {
	a := NcipolloRelease{}
	var _ workflow.StepAction = a
}

func TestNcipolloRelease_Inputs_Artifacts(t *testing.T) {
	a := NcipolloRelease{
		Artifacts: "dist/*.zip",
	}

	inputs := a.Inputs()

	if inputs["artifacts"] != "dist/*.zip" {
		t.Errorf("inputs[artifacts] = %v, want %q", inputs["artifacts"], "dist/*.zip")
	}
}

func TestNcipolloRelease_Inputs_ArtifactContentType(t *testing.T) {
	a := NcipolloRelease{
		ArtifactContentType: "application/zip",
	}

	inputs := a.Inputs()

	if inputs["artifactContentType"] != "application/zip" {
		t.Errorf("inputs[artifactContentType] = %v, want %q", inputs["artifactContentType"], "application/zip")
	}
}

func TestNcipolloRelease_Inputs_ArtifactErrorsFailBuild(t *testing.T) {
	a := NcipolloRelease{
		ArtifactErrorsFailBuild: true,
	}

	inputs := a.Inputs()

	if inputs["artifactErrorsFailBuild"] != true {
		t.Errorf("inputs[artifactErrorsFailBuild] = %v, want true", inputs["artifactErrorsFailBuild"])
	}
}

func TestNcipolloRelease_Inputs_Body(t *testing.T) {
	a := NcipolloRelease{
		Body: "Release notes for v1.0.0",
	}

	inputs := a.Inputs()

	if inputs["body"] != "Release notes for v1.0.0" {
		t.Errorf("inputs[body] = %v, want %q", inputs["body"], "Release notes for v1.0.0")
	}
}

func TestNcipolloRelease_Inputs_BodyFile(t *testing.T) {
	a := NcipolloRelease{
		BodyFile: "RELEASE_NOTES.md",
	}

	inputs := a.Inputs()

	if inputs["bodyFile"] != "RELEASE_NOTES.md" {
		t.Errorf("inputs[bodyFile] = %v, want %q", inputs["bodyFile"], "RELEASE_NOTES.md")
	}
}

func TestNcipolloRelease_Inputs_Commit(t *testing.T) {
	a := NcipolloRelease{
		Commit: "main",
	}

	inputs := a.Inputs()

	if inputs["commit"] != "main" {
		t.Errorf("inputs[commit] = %v, want %q", inputs["commit"], "main")
	}
}

func TestNcipolloRelease_Inputs_DiscussionCategory(t *testing.T) {
	a := NcipolloRelease{
		DiscussionCategory: "Announcements",
	}

	inputs := a.Inputs()

	if inputs["discussionCategory"] != "Announcements" {
		t.Errorf("inputs[discussionCategory] = %v, want %q", inputs["discussionCategory"], "Announcements")
	}
}

func TestNcipolloRelease_Inputs_Draft(t *testing.T) {
	a := NcipolloRelease{
		Draft: true,
	}

	inputs := a.Inputs()

	if inputs["draft"] != true {
		t.Errorf("inputs[draft] = %v, want true", inputs["draft"])
	}
}

func TestNcipolloRelease_Inputs_GenerateReleaseNotes(t *testing.T) {
	a := NcipolloRelease{
		GenerateReleaseNotes: true,
	}

	inputs := a.Inputs()

	if inputs["generateReleaseNotes"] != true {
		t.Errorf("inputs[generateReleaseNotes] = %v, want true", inputs["generateReleaseNotes"])
	}
}

func TestNcipolloRelease_Inputs_MakeLatest(t *testing.T) {
	a := NcipolloRelease{
		MakeLatest: "true",
	}

	inputs := a.Inputs()

	if inputs["makeLatest"] != "true" {
		t.Errorf("inputs[makeLatest] = %v, want %q", inputs["makeLatest"], "true")
	}
}

func TestNcipolloRelease_Inputs_MakeLatest_Legacy(t *testing.T) {
	a := NcipolloRelease{
		MakeLatest: "legacy",
	}

	inputs := a.Inputs()

	if inputs["makeLatest"] != "legacy" {
		t.Errorf("inputs[makeLatest] = %v, want %q", inputs["makeLatest"], "legacy")
	}
}

func TestNcipolloRelease_Inputs_Name(t *testing.T) {
	a := NcipolloRelease{
		Name: "Release v1.0.0",
	}

	inputs := a.Inputs()

	if inputs["name"] != "Release v1.0.0" {
		t.Errorf("inputs[name] = %v, want %q", inputs["name"], "Release v1.0.0")
	}
}

func TestNcipolloRelease_Inputs_OmitBody(t *testing.T) {
	a := NcipolloRelease{
		OmitBody: true,
	}

	inputs := a.Inputs()

	if inputs["omitBody"] != true {
		t.Errorf("inputs[omitBody] = %v, want true", inputs["omitBody"])
	}
}

func TestNcipolloRelease_Inputs_OmitBodyDuringUpdate(t *testing.T) {
	a := NcipolloRelease{
		OmitBodyDuringUpdate: true,
	}

	inputs := a.Inputs()

	if inputs["omitBodyDuringUpdate"] != true {
		t.Errorf("inputs[omitBodyDuringUpdate] = %v, want true", inputs["omitBodyDuringUpdate"])
	}
}

func TestNcipolloRelease_Inputs_OmitDraftDuringUpdate(t *testing.T) {
	a := NcipolloRelease{
		OmitDraftDuringUpdate: true,
	}

	inputs := a.Inputs()

	if inputs["omitDraftDuringUpdate"] != true {
		t.Errorf("inputs[omitDraftDuringUpdate] = %v, want true", inputs["omitDraftDuringUpdate"])
	}
}

func TestNcipolloRelease_Inputs_OmitName(t *testing.T) {
	a := NcipolloRelease{
		OmitName: true,
	}

	inputs := a.Inputs()

	if inputs["omitName"] != true {
		t.Errorf("inputs[omitName] = %v, want true", inputs["omitName"])
	}
}

func TestNcipolloRelease_Inputs_OmitNameDuringUpdate(t *testing.T) {
	a := NcipolloRelease{
		OmitNameDuringUpdate: true,
	}

	inputs := a.Inputs()

	if inputs["omitNameDuringUpdate"] != true {
		t.Errorf("inputs[omitNameDuringUpdate] = %v, want true", inputs["omitNameDuringUpdate"])
	}
}

func TestNcipolloRelease_Inputs_OmitPrereleaseDuringUpdate(t *testing.T) {
	a := NcipolloRelease{
		OmitPrereleaseDuringUpdate: true,
	}

	inputs := a.Inputs()

	if inputs["omitPrereleaseDuringUpdate"] != true {
		t.Errorf("inputs[omitPrereleaseDuringUpdate] = %v, want true", inputs["omitPrereleaseDuringUpdate"])
	}
}

func TestNcipolloRelease_Inputs_Owner(t *testing.T) {
	a := NcipolloRelease{
		Owner: "myorg",
	}

	inputs := a.Inputs()

	if inputs["owner"] != "myorg" {
		t.Errorf("inputs[owner] = %v, want %q", inputs["owner"], "myorg")
	}
}

func TestNcipolloRelease_Inputs_Prerelease(t *testing.T) {
	a := NcipolloRelease{
		Prerelease: true,
	}

	inputs := a.Inputs()

	if inputs["prerelease"] != true {
		t.Errorf("inputs[prerelease] = %v, want true", inputs["prerelease"])
	}
}

func TestNcipolloRelease_Inputs_RemoveArtifacts(t *testing.T) {
	a := NcipolloRelease{
		RemoveArtifacts: true,
	}

	inputs := a.Inputs()

	if inputs["removeArtifacts"] != true {
		t.Errorf("inputs[removeArtifacts] = %v, want true", inputs["removeArtifacts"])
	}
}

func TestNcipolloRelease_Inputs_ReplacesArtifacts(t *testing.T) {
	a := NcipolloRelease{
		ReplacesArtifacts: true,
	}

	inputs := a.Inputs()

	if inputs["replacesArtifacts"] != true {
		t.Errorf("inputs[replacesArtifacts] = %v, want true", inputs["replacesArtifacts"])
	}
}

func TestNcipolloRelease_Inputs_Repo(t *testing.T) {
	a := NcipolloRelease{
		Repo: "myrepo",
	}

	inputs := a.Inputs()

	if inputs["repo"] != "myrepo" {
		t.Errorf("inputs[repo] = %v, want %q", inputs["repo"], "myrepo")
	}
}

func TestNcipolloRelease_Inputs_SkipIfReleaseExists(t *testing.T) {
	a := NcipolloRelease{
		SkipIfReleaseExists: true,
	}

	inputs := a.Inputs()

	if inputs["skipIfReleaseExists"] != true {
		t.Errorf("inputs[skipIfReleaseExists] = %v, want true", inputs["skipIfReleaseExists"])
	}
}

func TestNcipolloRelease_Inputs_Tag(t *testing.T) {
	a := NcipolloRelease{
		Tag: "v1.0.0",
	}

	inputs := a.Inputs()

	if inputs["tag"] != "v1.0.0" {
		t.Errorf("inputs[tag] = %v, want %q", inputs["tag"], "v1.0.0")
	}
}

func TestNcipolloRelease_Inputs_Token(t *testing.T) {
	a := NcipolloRelease{
		Token: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if inputs["token"] != "${{ secrets.GITHUB_TOKEN }}" {
		t.Errorf("inputs[token] = %v, want %q", inputs["token"], "${{ secrets.GITHUB_TOKEN }}")
	}
}

func TestNcipolloRelease_Inputs_UpdateOnlyUnreleased(t *testing.T) {
	a := NcipolloRelease{
		UpdateOnlyUnreleased: true,
	}

	inputs := a.Inputs()

	if inputs["updateOnlyUnreleased"] != true {
		t.Errorf("inputs[updateOnlyUnreleased] = %v, want true", inputs["updateOnlyUnreleased"])
	}
}

func TestNcipolloRelease_Inputs_AllowUpdates(t *testing.T) {
	a := NcipolloRelease{
		AllowUpdates: true,
	}

	inputs := a.Inputs()

	if inputs["allowUpdates"] != true {
		t.Errorf("inputs[allowUpdates] = %v, want true", inputs["allowUpdates"])
	}
}

func TestNcipolloRelease_Inputs_AllFields(t *testing.T) {
	a := NcipolloRelease{
		Artifacts:                  "dist/*.zip",
		ArtifactContentType:        "application/zip",
		ArtifactErrorsFailBuild:    true,
		Body:                       "Release notes",
		BodyFile:                   "CHANGELOG.md",
		Commit:                     "main",
		DiscussionCategory:         "Announcements",
		Draft:                      true,
		GenerateReleaseNotes:       true,
		MakeLatest:                 "true",
		Name:                       "v1.0.0",
		OmitBody:                   true,
		OmitBodyDuringUpdate:       true,
		OmitDraftDuringUpdate:      true,
		OmitName:                   true,
		OmitNameDuringUpdate:       true,
		OmitPrereleaseDuringUpdate: true,
		Owner:                      "myorg",
		Prerelease:                 true,
		RemoveArtifacts:            true,
		ReplacesArtifacts:          true,
		Repo:                       "myrepo",
		SkipIfReleaseExists:        true,
		Tag:                        "v1.0.0",
		Token:                      "${{ secrets.GITHUB_TOKEN }}",
		UpdateOnlyUnreleased:       true,
		AllowUpdates:               true,
	}

	inputs := a.Inputs()

	expected := map[string]any{
		"artifacts":                  "dist/*.zip",
		"artifactContentType":        "application/zip",
		"artifactErrorsFailBuild":    true,
		"body":                       "Release notes",
		"bodyFile":                   "CHANGELOG.md",
		"commit":                     "main",
		"discussionCategory":         "Announcements",
		"draft":                      true,
		"generateReleaseNotes":       true,
		"makeLatest":                 "true",
		"name":                       "v1.0.0",
		"omitBody":                   true,
		"omitBodyDuringUpdate":       true,
		"omitDraftDuringUpdate":      true,
		"omitName":                   true,
		"omitNameDuringUpdate":       true,
		"omitPrereleaseDuringUpdate": true,
		"owner":                      "myorg",
		"prerelease":                 true,
		"removeArtifacts":            true,
		"replacesArtifacts":          true,
		"repo":                       "myrepo",
		"skipIfReleaseExists":        true,
		"tag":                        "v1.0.0",
		"token":                      "${{ secrets.GITHUB_TOKEN }}",
		"updateOnlyUnreleased":       true,
		"allowUpdates":               true,
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

func TestNcipolloRelease_Inputs_FalseBoolFields(t *testing.T) {
	a := NcipolloRelease{
		ArtifactErrorsFailBuild:    false,
		Draft:                      false,
		GenerateReleaseNotes:       false,
		OmitBody:                   false,
		OmitBodyDuringUpdate:       false,
		OmitDraftDuringUpdate:      false,
		OmitName:                   false,
		OmitNameDuringUpdate:       false,
		OmitPrereleaseDuringUpdate: false,
		Prerelease:                 false,
		RemoveArtifacts:            false,
		ReplacesArtifacts:          false,
		SkipIfReleaseExists:        false,
		UpdateOnlyUnreleased:       false,
		AllowUpdates:               false,
	}

	inputs := a.Inputs()

	if len(inputs) != 0 {
		t.Errorf("inputs for false bools has %d entries, want 0. Got: %v", len(inputs), inputs)
	}
}

func TestNcipolloRelease_Inputs_CommonUsage(t *testing.T) {
	// Test common usage pattern: simple release with artifacts
	a := NcipolloRelease{
		Artifacts:            "dist/*.tar.gz",
		Tag:                  "${{ github.ref_name }}",
		Token:                "${{ secrets.GITHUB_TOKEN }}",
		GenerateReleaseNotes: true,
	}

	inputs := a.Inputs()

	if len(inputs) != 4 {
		t.Errorf("common usage has %d entries, want 4", len(inputs))
	}

	if inputs["artifacts"] != "dist/*.tar.gz" {
		t.Errorf("inputs[artifacts] = %v, want dist/*.tar.gz", inputs["artifacts"])
	}
}

func TestNcipolloRelease_Inputs_DraftRelease(t *testing.T) {
	// Test draft release pattern
	a := NcipolloRelease{
		Tag:   "v1.0.0-beta.1",
		Draft: true,
		Name:  "Beta Release v1.0.0-beta.1",
		Token: "${{ secrets.GITHUB_TOKEN }}",
	}

	inputs := a.Inputs()

	if len(inputs) != 4 {
		t.Errorf("draft release has %d entries, want 4", len(inputs))
	}

	if inputs["draft"] != true {
		t.Errorf("inputs[draft] = %v, want true", inputs["draft"])
	}
}

func TestNcipolloRelease_Inputs_CrossRepoRelease(t *testing.T) {
	// Test cross-repo release pattern
	a := NcipolloRelease{
		Owner: "other-org",
		Repo:  "other-repo",
		Tag:   "v1.0.0",
		Token: "${{ secrets.PAT }}",
	}

	inputs := a.Inputs()

	if len(inputs) != 4 {
		t.Errorf("cross-repo release has %d entries, want 4", len(inputs))
	}

	if inputs["owner"] != "other-org" {
		t.Errorf("inputs[owner] = %v, want other-org", inputs["owner"])
	}

	if inputs["repo"] != "other-repo" {
		t.Errorf("inputs[repo] = %v, want other-repo", inputs["repo"])
	}
}
