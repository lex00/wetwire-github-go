// Package ncipollo_release provides a typed wrapper for ncipollo/release-action.
package ncipollo_release

// NcipolloRelease wraps the ncipollo/release-action@v1 action.
// Create GitHub releases with ease.
type NcipolloRelease struct {
	// Artifacts is a glob pattern for files to upload as release assets.
	Artifacts string `yaml:"artifacts,omitempty"`

	// ArtifactContentType sets the content type for uploaded artifacts.
	ArtifactContentType string `yaml:"artifactContentType,omitempty"`

	// ArtifactErrorsFailBuild fails the build if artifact upload fails.
	ArtifactErrorsFailBuild bool `yaml:"artifactErrorsFailBuild,omitempty"`

	// Body is the release body/description text.
	Body string `yaml:"body,omitempty"`

	// BodyFile is the path to a file containing the release body.
	BodyFile string `yaml:"bodyFile,omitempty"`

	// Commit is the commitish value for the release tag.
	Commit string `yaml:"commit,omitempty"`

	// DiscussionCategory creates a discussion linked to the release.
	DiscussionCategory string `yaml:"discussionCategory,omitempty"`

	// Draft creates the release as a draft.
	Draft bool `yaml:"draft,omitempty"`

	// GenerateReleaseNotes auto-generates release notes.
	GenerateReleaseNotes bool `yaml:"generateReleaseNotes,omitempty"`

	// MakeLatest marks this release as the latest.
	// Valid values: "true", "false", "legacy"
	MakeLatest string `yaml:"makeLatest,omitempty"`

	// Name is the release name/title.
	Name string `yaml:"name,omitempty"`

	// OmitBody omits the body from the release.
	OmitBody bool `yaml:"omitBody,omitempty"`

	// OmitBodyDuringUpdate omits the body when updating an existing release.
	OmitBodyDuringUpdate bool `yaml:"omitBodyDuringUpdate,omitempty"`

	// OmitDraftDuringUpdate omits the draft setting when updating.
	OmitDraftDuringUpdate bool `yaml:"omitDraftDuringUpdate,omitempty"`

	// OmitName omits the name from the release.
	OmitName bool `yaml:"omitName,omitempty"`

	// OmitNameDuringUpdate omits the name when updating an existing release.
	OmitNameDuringUpdate bool `yaml:"omitNameDuringUpdate,omitempty"`

	// OmitPrereleaseDuringUpdate omits the prerelease setting when updating.
	OmitPrereleaseDuringUpdate bool `yaml:"omitPrereleaseDuringUpdate,omitempty"`

	// Owner is the repository owner (defaults to current repository).
	Owner string `yaml:"owner,omitempty"`

	// Prerelease marks this release as a prerelease.
	Prerelease bool `yaml:"prerelease,omitempty"`

	// RemoveArtifacts removes existing artifacts before uploading new ones.
	RemoveArtifacts bool `yaml:"removeArtifacts,omitempty"`

	// ReplacesArtifacts replaces artifacts with the same name.
	ReplacesArtifacts bool `yaml:"replacesArtifacts,omitempty"`

	// Repo is the repository name (defaults to current repository).
	Repo string `yaml:"repo,omitempty"`

	// SkipIfReleaseExists skips creation if a release already exists.
	SkipIfReleaseExists bool `yaml:"skipIfReleaseExists,omitempty"`

	// Tag is the release tag name.
	Tag string `yaml:"tag,omitempty"`

	// Token is the GitHub token for authentication.
	Token string `yaml:"token,omitempty"`

	// UpdateOnlyUnreleased only updates releases that are not published.
	UpdateOnlyUnreleased bool `yaml:"updateOnlyUnreleased,omitempty"`

	// AllowUpdates allows updating an existing release.
	AllowUpdates bool `yaml:"allowUpdates,omitempty"`
}

// Action returns the action reference.
func (a NcipolloRelease) Action() string {
	return "ncipollo/release-action@v1"
}

// Inputs returns the action inputs as a map.
func (a NcipolloRelease) Inputs() map[string]any {
	with := make(map[string]any)

	if a.Artifacts != "" {
		with["artifacts"] = a.Artifacts
	}
	if a.ArtifactContentType != "" {
		with["artifactContentType"] = a.ArtifactContentType
	}
	if a.ArtifactErrorsFailBuild {
		with["artifactErrorsFailBuild"] = a.ArtifactErrorsFailBuild
	}
	if a.Body != "" {
		with["body"] = a.Body
	}
	if a.BodyFile != "" {
		with["bodyFile"] = a.BodyFile
	}
	if a.Commit != "" {
		with["commit"] = a.Commit
	}
	if a.DiscussionCategory != "" {
		with["discussionCategory"] = a.DiscussionCategory
	}
	if a.Draft {
		with["draft"] = a.Draft
	}
	if a.GenerateReleaseNotes {
		with["generateReleaseNotes"] = a.GenerateReleaseNotes
	}
	if a.MakeLatest != "" {
		with["makeLatest"] = a.MakeLatest
	}
	if a.Name != "" {
		with["name"] = a.Name
	}
	if a.OmitBody {
		with["omitBody"] = a.OmitBody
	}
	if a.OmitBodyDuringUpdate {
		with["omitBodyDuringUpdate"] = a.OmitBodyDuringUpdate
	}
	if a.OmitDraftDuringUpdate {
		with["omitDraftDuringUpdate"] = a.OmitDraftDuringUpdate
	}
	if a.OmitName {
		with["omitName"] = a.OmitName
	}
	if a.OmitNameDuringUpdate {
		with["omitNameDuringUpdate"] = a.OmitNameDuringUpdate
	}
	if a.OmitPrereleaseDuringUpdate {
		with["omitPrereleaseDuringUpdate"] = a.OmitPrereleaseDuringUpdate
	}
	if a.Owner != "" {
		with["owner"] = a.Owner
	}
	if a.Prerelease {
		with["prerelease"] = a.Prerelease
	}
	if a.RemoveArtifacts {
		with["removeArtifacts"] = a.RemoveArtifacts
	}
	if a.ReplacesArtifacts {
		with["replacesArtifacts"] = a.ReplacesArtifacts
	}
	if a.Repo != "" {
		with["repo"] = a.Repo
	}
	if a.SkipIfReleaseExists {
		with["skipIfReleaseExists"] = a.SkipIfReleaseExists
	}
	if a.Tag != "" {
		with["tag"] = a.Tag
	}
	if a.Token != "" {
		with["token"] = a.Token
	}
	if a.UpdateOnlyUnreleased {
		with["updateOnlyUnreleased"] = a.UpdateOnlyUnreleased
	}
	if a.AllowUpdates {
		with["allowUpdates"] = a.AllowUpdates
	}

	return with
}
