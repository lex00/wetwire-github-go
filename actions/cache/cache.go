// Package cache provides a typed wrapper for actions/cache.
package cache

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

// Cache wraps the actions/cache@v4 action.
// Cache dependencies and build outputs.
type Cache struct {
	// A list of files, directories, and wildcard patterns to cache and restore
	Path string `yaml:"path,omitempty"`

	// An explicit key for restoring and saving the cache
	Key string `yaml:"key,omitempty"`

	// An ordered list of keys to use for restoring stale cache if no hit for key
	RestoreKeys string `yaml:"restore-keys,omitempty"`

	// The chunk size used to split up large files during upload (in bytes)
	UploadChunkSize int `yaml:"upload-chunk-size,omitempty"`

	// An optional boolean to enable cross-os archive support
	EnableCrossOsArchive bool `yaml:"enableCrossOsArchive,omitempty"`

	// Fail the workflow if cache entry is not found
	FailOnCacheMiss bool `yaml:"fail-on-cache-miss,omitempty"`

	// Check if a cache entry exists without downloading
	LookupOnly bool `yaml:"lookup-only,omitempty"`

	// Run the post step to save the cache even if another step fails
	SaveAlways bool `yaml:"save-always,omitempty"`
}

// Action returns the action reference.
func (a Cache) Action() string {
	return "actions/cache@v4"
}

// ToStep converts this action to a workflow step.
func (a Cache) ToStep() workflow.Step {
	with := make(workflow.With)

	if a.Path != "" {
		with["path"] = a.Path
	}
	if a.Key != "" {
		with["key"] = a.Key
	}
	if a.RestoreKeys != "" {
		with["restore-keys"] = a.RestoreKeys
	}
	if a.UploadChunkSize != 0 {
		with["upload-chunk-size"] = a.UploadChunkSize
	}
	if a.EnableCrossOsArchive {
		with["enableCrossOsArchive"] = a.EnableCrossOsArchive
	}
	if a.FailOnCacheMiss {
		with["fail-on-cache-miss"] = a.FailOnCacheMiss
	}
	if a.LookupOnly {
		with["lookup-only"] = a.LookupOnly
	}
	if a.SaveAlways {
		with["save-always"] = a.SaveAlways
	}

	return workflow.Step{
		Uses: a.Action(),
		With: with,
	}
}
