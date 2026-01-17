package codegen

import (
	"net/http"
	"net/http/httptest"
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

func TestFetcher_Fetch_RetryOnError(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	f := NewFetcher()
	f.RetryDelay = 10 * time.Millisecond
	f.MaxRetries = 3

	data, err := f.Fetch(server.URL)
	if err != nil {
		t.Errorf("Fetch() error = %v, expected success after retries", err)
	}
	if string(data) != "success" {
		t.Errorf("Fetch() = %q, want %q", string(data), "success")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestFetcher_Fetch_AllRetriesFail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	f := NewFetcher()
	f.RetryDelay = 10 * time.Millisecond
	f.MaxRetries = 2

	_, err := f.Fetch(server.URL)
	if err == nil {
		t.Error("Fetch() expected error after all retries fail")
	}
}

func TestFetcher_Fetch_ConnectionError(t *testing.T) {
	f := NewFetcher()
	f.RetryDelay = 10 * time.Millisecond
	f.MaxRetries = 1

	// Use an invalid URL that will fail connection
	_, err := f.Fetch("http://localhost:99999/invalid")
	if err == nil {
		t.Error("Fetch() expected error for invalid URL")
	}
}

func TestFetcher_Fetch_ReadBodyError(t *testing.T) {
	// Test successful fetch with valid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response body"))
	}))
	defer server.Close()

	f := NewFetcher()
	f.RetryDelay = 10 * time.Millisecond

	data, err := f.Fetch(server.URL)
	if err != nil {
		t.Errorf("Fetch() error = %v", err)
	}
	if string(data) != "response body" {
		t.Errorf("Fetch() = %q, want %q", string(data), "response body")
	}
}

func TestFetcher_Fetch_ZeroRetries(t *testing.T) {
	// Test with zero retries to ensure first failure returns error immediately
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer server.Close()

	f := &Fetcher{
		Client:     &http.Client{Timeout: 30 * time.Second},
		MaxRetries: 0,
		RetryDelay: 10 * time.Millisecond,
	}

	_, err := f.Fetch(server.URL)
	if err == nil {
		t.Error("Fetch() expected error with zero retries and failing server")
	}
}

func TestFetcher_Fetch_FirstAttemptSuccess(t *testing.T) {
	// Test that first successful attempt returns immediately without retry
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))
	defer server.Close()

	f := &Fetcher{
		Client:     &http.Client{Timeout: 30 * time.Second},
		MaxRetries: 3,
		RetryDelay: 10 * time.Millisecond,
	}

	data, err := f.Fetch(server.URL)
	if err != nil {
		t.Errorf("Fetch() error = %v", err)
	}
	if string(data) != "success" {
		t.Errorf("Fetch() = %q, want %q", string(data), "success")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt for successful first request, got %d", attempts)
	}
}

// mockHTTPTransport allows mocking HTTP responses for testing
type mockHTTPTransport struct {
	response *http.Response
	err      error
}

func (m *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.response, m.err
}
