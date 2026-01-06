package codegen

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestActionURL(t *testing.T) {
	tests := []struct {
		owner, repo string
		want        string
	}{
		{"actions", "checkout", "https://raw.githubusercontent.com/actions/checkout/main/action.yml"},
		{"actions", "setup-go", "https://raw.githubusercontent.com/actions/setup-go/main/action.yml"},
		{"owner", "repo", "https://raw.githubusercontent.com/owner/repo/main/action.yml"},
	}

	for _, tt := range tests {
		got := ActionURL(tt.owner, tt.repo)
		if got != tt.want {
			t.Errorf("ActionURL(%q, %q) = %q, want %q", tt.owner, tt.repo, got, tt.want)
		}
	}
}

func TestFetcher_Fetch(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/success":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"test": "data"}`))
		case "/not-found":
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	f := NewFetcher()
	f.RetryDelay = 10 * time.Millisecond // Speed up tests

	// Test successful fetch
	data, err := f.Fetch(server.URL + "/success")
	if err != nil {
		t.Errorf("Fetch() error = %v", err)
	}
	if string(data) != `{"test": "data"}` {
		t.Errorf("Fetch() = %q, want %q", string(data), `{"test": "data"}`)
	}

	// Test 404
	_, err = f.Fetch(server.URL + "/not-found")
	if err == nil {
		t.Error("Fetch() expected error for 404, got nil")
	}
}

func TestFetcher_FetchSchema(t *testing.T) {
	f := NewFetcher()

	// Test unknown schema type
	_, err := f.FetchSchema("unknown")
	if err == nil {
		t.Error("FetchSchema() expected error for unknown type, got nil")
	}
}

func TestFetcher_FetchAction(t *testing.T) {
	f := NewFetcher()

	// Test unknown action
	_, err := f.FetchAction("unknown-action")
	if err == nil {
		t.Error("FetchAction() expected error for unknown action, got nil")
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

func TestFetcher_FetchAll(t *testing.T) {
	// Create a test server that returns mock data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"mock": "data"}`))
	}))
	defer server.Close()

	// Override URLs for testing
	originalSchemaURLs := SchemaURLs
	SchemaURLs = map[SchemaType]string{
		SchemaWorkflow: server.URL + "/workflow.json",
	}
	defer func() { SchemaURLs = originalSchemaURLs }()

	originalPopularActions := PopularActions
	PopularActions = map[string]struct{ Owner, Repo string }{
		"test-action": {"test", "action"},
	}
	defer func() { PopularActions = originalPopularActions }()

	// Temporarily override ActionURL behavior by mocking the fetch
	f := NewFetcher()
	f.RetryDelay = 10 * time.Millisecond

	// Create a temp directory
	tmpDir := t.TempDir()

	// This will fail because ActionURL still points to real GitHub
	// For a full test we'd need to mock the HTTP client entirely
	// For now, we test the local server URLs

	// Test with just schemas
	PopularActions = map[string]struct{ Owner, Repo string }{}

	manifest, err := f.FetchAll(tmpDir)
	if err != nil {
		t.Fatalf("FetchAll() error = %v", err)
	}

	if manifest.Version != "1.0" {
		t.Errorf("manifest.Version = %q, want %q", manifest.Version, "1.0")
	}

	if len(manifest.Schemas) != 1 {
		t.Errorf("len(manifest.Schemas) = %d, want 1", len(manifest.Schemas))
	}

	// Check manifest file was created
	manifestPath := filepath.Join(tmpDir, "manifest.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Error("manifest.json was not created")
	}

	// Check schema file was created
	schemaPath := filepath.Join(tmpDir, "workflow.json")
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		t.Error("workflow.json was not created")
	}
}

func TestNewFetcher(t *testing.T) {
	f := NewFetcher()

	if f.Client == nil {
		t.Error("NewFetcher() Client is nil")
	}
	if f.MaxRetries != 3 {
		t.Errorf("NewFetcher() MaxRetries = %d, want 3", f.MaxRetries)
	}
	if f.RetryDelay != 1*time.Second {
		t.Errorf("NewFetcher() RetryDelay = %v, want 1s", f.RetryDelay)
	}
}
