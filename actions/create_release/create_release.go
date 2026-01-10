// Package create_release provides a typed wrapper for actions/create-release.
package create_release

// CreateRelease wraps the actions/create-release@v1 action.
// Create a release for a repository tag.
//
// Note: This action is deprecated by GitHub and unmaintained.
// Consider using softprops/action-gh-release or ncipollo/release-action instead.
type CreateRelease struct {
	// TagName is the name of the tag for this release (required).
	TagName string `yaml:"tag_name,omitempty"`

	// ReleaseName is the name of the release (required).
	ReleaseName string `yaml:"release_name,omitempty"`

	// Body is text describing the contents of the release.
	// Optional, and not needed if using BodyPath.
	Body string `yaml:"body,omitempty"`

	// BodyPath is a file with contents describing the release.
	// Optional, and not needed if using Body.
	BodyPath string `yaml:"body_path,omitempty"`

	// Draft creates a draft (unpublished) release when true.
	// Defaults to false.
	Draft bool `yaml:"draft,omitempty"`

	// Prerelease identifies the release as a prerelease when true.
	// Defaults to false.
	Prerelease bool `yaml:"prerelease,omitempty"`

	// Commitish is any branch or commit SHA the Git tag is created from.
	// Unused if the Git tag already exists. Defaults to the SHA of current commit.
	Commitish string `yaml:"commitish,omitempty"`

	// Owner is the name of the owner of the repo.
	// Used when cutting releases for external repositories.
	Owner string `yaml:"owner,omitempty"`

	// Repo is the name of the repository.
	// Used when cutting releases for external repositories.
	Repo string `yaml:"repo,omitempty"`
}

// Action returns the action reference.
func (a CreateRelease) Action() string {
	return "actions/create-release@v1"
}

// Inputs returns the action inputs as a map.
func (a CreateRelease) Inputs() map[string]any {
	with := make(map[string]any)

	if a.TagName != "" {
		with["tag_name"] = a.TagName
	}
	if a.ReleaseName != "" {
		with["release_name"] = a.ReleaseName
	}
	if a.Body != "" {
		with["body"] = a.Body
	}
	if a.BodyPath != "" {
		with["body_path"] = a.BodyPath
	}
	if a.Draft {
		with["draft"] = a.Draft
	}
	if a.Prerelease {
		with["prerelease"] = a.Prerelease
	}
	if a.Commitish != "" {
		with["commitish"] = a.Commitish
	}
	if a.Owner != "" {
		with["owner"] = a.Owner
	}
	if a.Repo != "" {
		with["repo"] = a.Repo
	}

	return with
}
