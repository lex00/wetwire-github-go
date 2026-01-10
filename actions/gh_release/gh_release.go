// Package gh_release provides a typed wrapper for softprops/action-gh-release.
package gh_release

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// GHRelease wraps the softprops/action-gh-release@v2 action.
// Create and upload assets to a GitHub Release.
type GHRelease struct {
	// Body of the release. Can include markdown.
	Body string `yaml:"body,omitempty"`

	// Path to a file with the release body content.
	BodyPath string `yaml:"body_path,omitempty"`

	// Name of the release. If not specified, uses tag name.
	Name string `yaml:"name,omitempty"`

	// Tag name for the release. If not specified, uses GITHUB_REF.
	TagName string `yaml:"tag_name,omitempty"`

	// Commitish value to tag. Defaults to the repository's default branch.
	TargetCommitish string `yaml:"target_commitish,omitempty"`

	// Whether this is a draft release. Drafts are not visible to users.
	Draft bool `yaml:"draft,omitempty"`

	// Whether this is a prerelease.
	Prerelease bool `yaml:"prerelease,omitempty"`

	// Whether to automatically generate the name and body for this release.
	GenerateReleaseNotes bool `yaml:"generate_release_notes,omitempty"`

	// Newline-separated list of glob patterns for files to upload.
	Files string `yaml:"files,omitempty"`

	// Whether to fail if no files are matched for upload.
	FailOnUnmatchedFiles bool `yaml:"fail_on_unmatched_files,omitempty"`

	// GitHub token for authentication.
	Token string `yaml:"token,omitempty"`

	// Repository to release to (format: owner/repo).
	Repository string `yaml:"repository,omitempty"`

	// Whether to append body content to existing release.
	AppendBody bool `yaml:"append_body,omitempty"`

	// Whether to only create release if none exists for the tag.
	MakeLatest string `yaml:"make_latest,omitempty"`

	// Discussion category name for the release.
	DiscussionCategoryName string `yaml:"discussion_category_name,omitempty"`
}

// Action returns the action reference.
func (a GHRelease) Action() string {
	return "softprops/action-gh-release@v2"
}

// ToStep converts this action to a workflow step.
func (a GHRelease) ToStep() workflow.Step {
	with := make(workflow.With)

	if a.Body != "" {
		with["body"] = a.Body
	}
	if a.BodyPath != "" {
		with["body_path"] = a.BodyPath
	}
	if a.Name != "" {
		with["name"] = a.Name
	}
	if a.TagName != "" {
		with["tag_name"] = a.TagName
	}
	if a.TargetCommitish != "" {
		with["target_commitish"] = a.TargetCommitish
	}
	if a.Draft {
		with["draft"] = a.Draft
	}
	if a.Prerelease {
		with["prerelease"] = a.Prerelease
	}
	if a.GenerateReleaseNotes {
		with["generate_release_notes"] = a.GenerateReleaseNotes
	}
	if a.Files != "" {
		with["files"] = a.Files
	}
	if a.FailOnUnmatchedFiles {
		with["fail_on_unmatched_files"] = a.FailOnUnmatchedFiles
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.Repository != "" {
		with["repository"] = a.Repository
	}
	if a.AppendBody {
		with["append_body"] = a.AppendBody
	}
	if a.MakeLatest != "" {
		with["make_latest"] = a.MakeLatest
	}
	if a.DiscussionCategoryName != "" {
		with["discussion_category_name"] = a.DiscussionCategoryName
	}

	return workflow.Step{
		Uses: a.Action(),
		With: with,
	}
}
