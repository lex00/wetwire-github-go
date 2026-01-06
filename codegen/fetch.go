// Package codegen provides schema fetching and code generation for GitHub Actions.
package codegen

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// SchemaType represents a type of schema that can be fetched.
type SchemaType string

const (
	SchemaWorkflow   SchemaType = "workflow"
	SchemaDependabot SchemaType = "dependabot"
	SchemaIssueForms SchemaType = "issue-forms"
	SchemaAction     SchemaType = "action"
)

// SchemaURLs contains the URLs for various JSON schemas.
var SchemaURLs = map[SchemaType]string{
	SchemaWorkflow:   "https://json.schemastore.org/github-workflow.json",
	SchemaDependabot: "https://json.schemastore.org/dependabot-2.0.json",
	SchemaIssueForms: "https://json.schemastore.org/github-issue-forms.json",
}

// ActionURL returns the URL for an action's action.yml file.
// Pattern: https://raw.githubusercontent.com/{owner}/{repo}/main/action.yml
func ActionURL(owner, repo string) string {
	return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/action.yml", owner, repo)
}

// PopularActions contains commonly used GitHub Actions.
var PopularActions = map[string]struct{ Owner, Repo string }{
	"checkout":          {"actions", "checkout"},
	"setup-go":          {"actions", "setup-go"},
	"setup-node":        {"actions", "setup-node"},
	"setup-python":      {"actions", "setup-python"},
	"cache":             {"actions", "cache"},
	"upload-artifact":   {"actions", "upload-artifact"},
	"download-artifact": {"actions", "download-artifact"},
}

// Fetcher handles HTTP requests for schemas with retry logic.
type Fetcher struct {
	Client     *http.Client
	MaxRetries int
	RetryDelay time.Duration
}

// NewFetcher creates a new Fetcher with default settings.
func NewFetcher() *Fetcher {
	return &Fetcher{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
	}
}

// Fetch retrieves content from the given URL with retry logic.
func (f *Fetcher) Fetch(url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= f.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(f.RetryDelay * time.Duration(attempt))
		}

		resp, err := f.Client.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			continue
		}

		if err != nil {
			lastErr = fmt.Errorf("reading response: %w", err)
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("after %d retries: %w", f.MaxRetries, lastErr)
}

// FetchSchema fetches a JSON schema by type.
func (f *Fetcher) FetchSchema(schemaType SchemaType) ([]byte, error) {
	url, ok := SchemaURLs[schemaType]
	if !ok {
		return nil, fmt.Errorf("unknown schema type: %s", schemaType)
	}
	return f.Fetch(url)
}

// FetchAction fetches an action.yml file for the given action.
func (f *Fetcher) FetchAction(name string) ([]byte, error) {
	action, ok := PopularActions[name]
	if !ok {
		return nil, fmt.Errorf("unknown action: %s", name)
	}
	return f.Fetch(ActionURL(action.Owner, action.Repo))
}

// FetchActionByOwnerRepo fetches an action.yml file for the given owner/repo.
func (f *Fetcher) FetchActionByOwnerRepo(owner, repo string) ([]byte, error) {
	return f.Fetch(ActionURL(owner, repo))
}

// Manifest represents the specs/manifest.json file.
type Manifest struct {
	Version  string           `json:"version"`
	Schemas  []ManifestSchema `json:"schemas"`
	Actions  []ManifestAction `json:"actions"`
	FetchedAt string          `json:"fetched_at"`
}

// ManifestSchema represents a schema entry in the manifest.
type ManifestSchema struct {
	Type SchemaType `json:"type"`
	URL  string     `json:"url"`
	File string     `json:"file"`
}

// ManifestAction represents an action entry in the manifest.
type ManifestAction struct {
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	URL   string `json:"url"`
	File  string `json:"file"`
}

// FetchAll fetches all schemas and actions to the specified output directory.
// Returns the manifest of fetched files.
func (f *Fetcher) FetchAll(outputDir string) (*Manifest, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	manifest := &Manifest{
		Version:   "1.0",
		Schemas:   []ManifestSchema{},
		Actions:   []ManifestAction{},
		FetchedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// Fetch JSON schemas
	for schemaType, url := range SchemaURLs {
		filename := fmt.Sprintf("%s.json", schemaType)
		filepath := filepath.Join(outputDir, filename)

		data, err := f.Fetch(url)
		if err != nil {
			return nil, fmt.Errorf("fetching %s schema: %w", schemaType, err)
		}

		if err := os.WriteFile(filepath, data, 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", filepath, err)
		}

		manifest.Schemas = append(manifest.Schemas, ManifestSchema{
			Type: schemaType,
			URL:  url,
			File: filename,
		})
	}

	// Fetch action.yml files
	for name, action := range PopularActions {
		url := ActionURL(action.Owner, action.Repo)
		filename := fmt.Sprintf("%s.yml", name)
		filepath := filepath.Join(outputDir, filename)

		data, err := f.Fetch(url)
		if err != nil {
			return nil, fmt.Errorf("fetching %s action: %w", name, err)
		}

		if err := os.WriteFile(filepath, data, 0644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", filepath, err)
		}

		manifest.Actions = append(manifest.Actions, ManifestAction{
			Name:  name,
			Owner: action.Owner,
			Repo:  action.Repo,
			URL:   url,
			File:  filename,
		})
	}

	// Write manifest
	manifestPath := filepath.Join(outputDir, "manifest.json")
	manifestData, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshaling manifest: %w", err)
	}
	if err := os.WriteFile(manifestPath, manifestData, 0644); err != nil {
		return nil, fmt.Errorf("writing manifest: %w", err)
	}

	return manifest, nil
}
