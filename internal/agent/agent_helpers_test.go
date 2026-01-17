package agent

import (
	"context"
	"fmt"

	"github.com/lex00/wetwire-core-go/providers"
)

// mockMessage helps create providers.MessageResponse for testing
type mockMessage struct {
	text string
}

func (m *mockMessage) toMessageResponse() *providers.MessageResponse {
	return &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: m.text},
		},
	}
}

// contains checks if s contains substr (case-insensitive would be needed for some checks)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// mockDeveloper implements orchestrator.Developer for testing
type mockDeveloper struct {
	response string
	err      error
}

func (m *mockDeveloper) Respond(ctx context.Context, question string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.response, nil
}

// formatIteration is a helper to generate iteration-specific file names
func formatIteration(i int) string {
	return fmt.Sprintf("file%d.go", i)
}
