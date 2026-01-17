package codegen

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetcher_FetchSchema(t *testing.T) {
	f := NewFetcher()

	// Test unknown schema type
	_, err := f.FetchSchema("unknown")
	if err == nil {
		t.Error("FetchSchema() expected error for unknown type, got nil")
	}
}

func TestPopularActions(t *testing.T) {
	expected := []string{
		"checkout",
		"setup-go",
		"setup-node",
		"setup-python",
		"cache",
		"upload-artifact",
		"download-artifact",
	}

	for _, name := range expected {
		if _, ok := PopularActions[name]; !ok {
			t.Errorf("PopularActions missing %q", name)
		}
	}
}

func TestSchemaURLs(t *testing.T) {
	if _, ok := SchemaURLs[SchemaWorkflow]; !ok {
		t.Error("SchemaURLs missing workflow")
	}
	if _, ok := SchemaURLs[SchemaDependabot]; !ok {
		t.Error("SchemaURLs missing dependabot")
	}
	if _, ok := SchemaURLs[SchemaIssueForms]; !ok {
		t.Error("SchemaURLs missing issue-forms")
	}
}

func TestFetcher_FetchSchema_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"title": "Workflow Schema"}`))
	}))
	defer server.Close()

	// Override SchemaURLs temporarily
	originalURLs := SchemaURLs
	SchemaURLs = map[SchemaType]string{
		SchemaWorkflow: server.URL,
	}
	defer func() { SchemaURLs = originalURLs }()

	f := NewFetcher()
	data, err := f.FetchSchema(SchemaWorkflow)
	if err != nil {
		t.Errorf("FetchSchema() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("FetchSchema() returned empty data")
	}
}

func TestManifestTypes(t *testing.T) {
	// Test that manifest types are properly initialized
	manifest := &Manifest{
		Version:   "1.0",
		Schemas:   []ManifestSchema{},
		Actions:   []ManifestAction{},
		FetchedAt: "2024-01-01T00:00:00Z",
	}

	if manifest.Version != "1.0" {
		t.Errorf("manifest.Version = %q, want %q", manifest.Version, "1.0")
	}

	schema := ManifestSchema{
		Type: SchemaWorkflow,
		URL:  "https://example.com/schema.json",
		File: "schema.json",
	}
	if schema.Type != SchemaWorkflow {
		t.Errorf("schema.Type = %v, want %v", schema.Type, SchemaWorkflow)
	}

	action := ManifestAction{
		Name:  "test",
		Owner: "owner",
		Repo:  "repo",
		URL:   "https://example.com/action.yml",
		File:  "action.yml",
	}
	if action.Name != "test" {
		t.Errorf("action.Name = %q, want %q", action.Name, "test")
	}
}

func TestSchemaTypes(t *testing.T) {
	// Test all schema type constants
	types := []SchemaType{
		SchemaWorkflow,
		SchemaDependabot,
		SchemaIssueForms,
		SchemaAction,
	}

	for _, st := range types {
		if st == "" {
			t.Error("Schema type should not be empty")
		}
	}
}
