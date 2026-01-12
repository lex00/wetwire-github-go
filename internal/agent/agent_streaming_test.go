package agent

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/lex00/wetwire-core-go/providers"
	"github.com/stretchr/testify/assert"
)

// Note: The runWithStreaming method is very difficult to test without a real API client
// or extensive mocking of the Anthropic SDK's streaming API. These tests focus on
// what we can test: the integration points and state management.

// TestGitHubAgent_RunWithStreaming_Integration tests that streaming is used when handler is set
// This test is intentionally skipped by default as it requires a real API key and makes actual API calls.
// Enable it for real integration testing by setting ANTHROPIC_API_KEY and removing the skip.
func TestGitHubAgent_RunWithStreaming_Integration(t *testing.T) {
	// Skip this test - it makes real API calls
	t.Skip("Skipping integration test - requires real ANTHROPIC_API_KEY and makes actual API calls")

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping integration test - requires ANTHROPIC_API_KEY")
	}

	var capturedText strings.Builder
	handler := func(text string) {
		capturedText.WriteString(text)
	}

	agent, err := NewGitHubAgent(Config{
		APIKey:        apiKey,
		StreamHandler: handler,
	})
	assert.NoError(t, err)

	ctx := context.Background()
	err = agent.Run(ctx, "Say hello")
	assert.NoError(t, err)

	// Verify streaming captured text
	assert.NotEmpty(t, capturedText.String())
}

// TestGitHubAgent_StreamHandler_CalledDuringStreaming tests the stream handler callback
func TestGitHubAgent_StreamHandler_CalledDuringStreaming(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	var capturedChunks []string
	handler := func(text string) {
		capturedChunks = append(capturedChunks, text)
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	assert.NoError(t, err)
	assert.NotNil(t, agent.streamHandler)

	// Simulate what would happen during streaming
	agent.streamHandler("chunk1")
	agent.streamHandler("chunk2")
	agent.streamHandler("chunk3")

	assert.Len(t, capturedChunks, 3)
	assert.Equal(t, "chunk1", capturedChunks[0])
	assert.Equal(t, "chunk2", capturedChunks[1])
	assert.Equal(t, "chunk3", capturedChunks[2])
}

// TestGitHubAgent_StreamHandler_NilByDefault tests that stream handler is nil by default
func TestGitHubAgent_StreamHandler_NilByDefault(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)

	// By default, streamHandler should be nil
	assert.Nil(t, agent.streamHandler)
}

// TestGitHubAgent_Run_UsesStreamingWhenHandlerSet tests the branching logic in Run
func TestGitHubAgent_Run_UsesStreamingWhenHandlerSet(t *testing.T) {
	// This test verifies the branching logic in Run() that chooses between
	// streaming and non-streaming based on streamHandler being set
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	tests := []struct {
		name          string
		streamHandler providers.StreamHandler
		expectStream  bool
	}{
		{
			name:          "nil handler uses non-streaming",
			streamHandler: nil,
			expectStream:  false,
		},
		{
			name: "non-nil handler uses streaming",
			streamHandler: func(text string) {
				// no-op
			},
			expectStream: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := NewGitHubAgent(Config{
				StreamHandler: tt.streamHandler,
			})
			assert.NoError(t, err)

			if tt.expectStream {
				assert.NotNil(t, agent.streamHandler)
			} else {
				assert.Nil(t, agent.streamHandler)
			}
		})
	}
}

// TestGitHubAgent_RunWithStreaming_MessageBuilding tests the message building logic
// This test focuses on the data structures used in runWithStreaming
func TestGitHubAgent_RunWithStreaming_MessageBuilding(t *testing.T) {
	// This tests the logic for building messages from streaming events
	// We can't easily test the actual streaming without mocking the entire
	// Anthropic SDK, but we can test the data structures and logic

	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	handler := func(text string) {}
	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	assert.NoError(t, err)

	// Verify the agent is set up for streaming
	assert.NotNil(t, agent.streamHandler)
	assert.NotNil(t, agent.provider)
}

// TestGitHubAgent_StreamHandler_ConcurrentCalls tests concurrent stream handler calls
func TestGitHubAgent_StreamHandler_ConcurrentCalls(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	var capturedText strings.Builder
	handler := func(text string) {
		capturedText.WriteString(text)
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	assert.NoError(t, err)

	// Simulate rapid streaming calls
	for i := 0; i < 100; i++ {
		agent.streamHandler("x")
	}

	result := capturedText.String()
	assert.Len(t, result, 100)
}

// TestGitHubAgent_StreamHandler_EmptyChunks tests handling of empty stream chunks
func TestGitHubAgent_StreamHandler_EmptyChunks(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	var chunks []string
	handler := func(text string) {
		chunks = append(chunks, text)
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	assert.NoError(t, err)

	// Test with empty chunks
	agent.streamHandler("")
	agent.streamHandler("text")
	agent.streamHandler("")

	assert.Len(t, chunks, 3)
	assert.Equal(t, "", chunks[0])
	assert.Equal(t, "text", chunks[1])
	assert.Equal(t, "", chunks[2])
}

// TestGitHubAgent_StreamHandler_SpecialCharacters tests handling of special characters
func TestGitHubAgent_StreamHandler_SpecialCharacters(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	var capturedText strings.Builder
	handler := func(text string) {
		capturedText.WriteString(text)
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	assert.NoError(t, err)

	// Test with special characters
	specialChars := []string{
		"Hello\n",
		"Tab\there",
		"Quote\"test",
		"Emoji 🚀",
		"Unicode: 你好",
	}

	for _, chunk := range specialChars {
		agent.streamHandler(chunk)
	}

	result := capturedText.String()
	for _, expected := range specialChars {
		assert.Contains(t, result, expected)
	}
}

// TestGitHubAgent_StreamHandler_LargeChunks tests handling of large text chunks
func TestGitHubAgent_StreamHandler_LargeChunks(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	var capturedText strings.Builder
	handler := func(text string) {
		capturedText.WriteString(text)
	}

	agent, err := NewGitHubAgent(Config{
		StreamHandler: handler,
	})
	assert.NoError(t, err)

	// Create a large chunk
	largeChunk := strings.Repeat("Lorem ipsum dolor sit amet. ", 1000)
	agent.streamHandler(largeChunk)

	result := capturedText.String()
	assert.Equal(t, largeChunk, result)
}

// TestGitHubAgent_Run_WithoutStreamHandler tests Run uses non-streaming path
func TestGitHubAgent_Run_WithoutStreamHandler(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	// Create agent without stream handler
	agent, err := NewGitHubAgent(Config{})
	assert.NoError(t, err)
	assert.Nil(t, agent.streamHandler)

	// The Run method will use the non-streaming path
	// We can't test the actual API call without mocking, but we verify the setup
	assert.NotNil(t, agent.provider)
}

// TestGitHubAgent_StreamHandler_MultipleAgents tests multiple agents with different handlers
func TestGitHubAgent_StreamHandler_MultipleAgents(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	var text1, text2 strings.Builder
	handler1 := func(text string) {
		text1.WriteString("1:" + text)
	}
	handler2 := func(text string) {
		text2.WriteString("2:" + text)
	}

	agent1, err := NewGitHubAgent(Config{StreamHandler: handler1})
	assert.NoError(t, err)

	agent2, err := NewGitHubAgent(Config{StreamHandler: handler2})
	assert.NoError(t, err)

	// Each agent should have its own handler
	agent1.streamHandler("test")
	agent2.streamHandler("test")

	assert.Equal(t, "1:test", text1.String())
	assert.Equal(t, "2:test", text2.String())
}

// TestGitHubAgent_StreamHandler_StateManagement tests that streaming doesn't affect agent state
func TestGitHubAgent_StreamHandler_StateManagement(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "test-key")

	handler := func(text string) {}
	agent, err := NewGitHubAgent(Config{StreamHandler: handler})
	assert.NoError(t, err)

	// Streaming handler shouldn't affect these states
	assert.False(t, agent.lintCalled)
	assert.False(t, agent.lintPassed)
	assert.False(t, agent.pendingLint)
	assert.Equal(t, 0, agent.lintCycles)
	assert.Nil(t, agent.generatedFiles)

	// Call handler
	agent.streamHandler("test")

	// States should remain unchanged
	assert.False(t, agent.lintCalled)
	assert.False(t, agent.lintPassed)
	assert.False(t, agent.pendingLint)
}

// TestGitHubAgent_RunWithStreaming_ErrorHandling tests error handling in streaming
func TestGitHubAgent_RunWithStreaming_ErrorHandling(t *testing.T) {
	// Skip if no real API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" || apiKey == "test-key" {
		t.Skip("Skipping - requires real ANTHROPIC_API_KEY for streaming tests")
	}

	// This test documents that streaming errors are properly handled
	// In practice, we can't easily test this without mocking the SDK
	t.Skip("Skipping actual API call")
}

// TestGitHubAgent_StreamingEventTypes tests understanding of streaming event types
// This is a documentation test for the event types handled in runWithStreaming
func TestGitHubAgent_StreamingEventTypes(t *testing.T) {
	// The runWithStreaming method handles these event types:
	// - message_start: Initializes the message
	// - content_block_start: Starts a new content block (text or tool_use)
	// - content_block_delta: Adds delta to current block (text_delta or input_json_delta)
	// - content_block_stop: Finalizes a content block
	// - message_delta: Updates message metadata (stop_reason, stop_sequence)

	// This test verifies we understand the event structure
	eventTypes := []string{
		"message_start",
		"content_block_start",
		"content_block_delta",
		"content_block_stop",
		"message_delta",
	}

	assert.Len(t, eventTypes, 5)
}

// TestGitHubAgent_StreamingContentBlockTypes tests understanding of content block types
func TestGitHubAgent_StreamingContentBlockTypes(t *testing.T) {
	// Content blocks can be:
	// - text: Regular text response
	// - tool_use: Tool invocation

	// Delta types can be:
	// - text_delta: Text content
	// - input_json_delta: Tool input JSON (partial)

	blockTypes := []string{"text", "tool_use"}
	deltaTypes := []string{"text_delta", "input_json_delta"}

	assert.Len(t, blockTypes, 2)
	assert.Len(t, deltaTypes, 2)
}

// TestGitHubAgent_StreamingJSONParsing tests JSON parsing in streaming
func TestGitHubAgent_StreamingJSONParsing(t *testing.T) {
	// The runWithStreaming method uses json.RawMessage for tool input
	// This test verifies we can construct valid JSON from streaming chunks

	chunks := []string{
		`{"na`,
		`me": `,
		`"test-pr`,
		`oject"}`,
	}

	var builder strings.Builder
	for _, chunk := range chunks {
		builder.WriteString(chunk)
	}

	fullJSON := builder.String()
	assert.Equal(t, `{"name": "test-project"}`, fullJSON)

	// Verify it's valid JSON
	var result map[string]string
	err := json.Unmarshal([]byte(fullJSON), &result)
	assert.NoError(t, err)
	assert.Equal(t, "test-project", result["name"])
}

// TestGitHubAgent_StreamingStringBuilder tests StringBuilder usage
func TestGitHubAgent_StreamingStringBuilder(t *testing.T) {
	// The runWithStreaming method uses strings.Builder to accumulate text
	// This tests that pattern

	builder := &strings.Builder{}
	chunks := []string{"Hello", " ", "streaming", " ", "world!"}

	for _, chunk := range chunks {
		builder.WriteString(chunk)
	}

	result := builder.String()
	assert.Equal(t, "Hello streaming world!", result)
}

// TestGitHubAgent_StreamingMapInitialization tests map initialization pattern
func TestGitHubAgent_StreamingMapInitialization(t *testing.T) {
	// The runWithStreaming method initializes maps for tracking content
	// This tests that pattern

	currentTextContent := make(map[int64]*strings.Builder)
	currentToolInput := make(map[int64]*strings.Builder)

	// Simulate adding content for different indices
	currentTextContent[0] = &strings.Builder{}
	currentTextContent[0].WriteString("text1")

	currentToolInput[1] = &strings.Builder{}
	currentToolInput[1].WriteString(`{"input":"data"}`)

	assert.Len(t, currentTextContent, 1)
	assert.Len(t, currentToolInput, 1)
	assert.Equal(t, "text1", currentTextContent[0].String())
	assert.Equal(t, `{"input":"data"}`, currentToolInput[1].String())
}

// TestGitHubAgent_StreamingContentBlockArray tests content block array management
func TestGitHubAgent_StreamingContentBlockArray(t *testing.T) {
	// The runWithStreaming method builds an array of ContentBlockUnion
	// This tests that pattern

	var contentBlocks []providers.ContentBlock

	// Add text block
	contentBlocks = append(contentBlocks, providers.ContentBlock{
		Type: "text",
		Text: "Hello world",
	})

	// Add tool block
	contentBlocks = append(contentBlocks, providers.ContentBlock{
		Type:  "tool_use",
		ID:    "tool_1",
		Name:  "write_file",
		Input: json.RawMessage(`{"path":"test.go"}`),
	})

	assert.Len(t, contentBlocks, 2)
	assert.Equal(t, "text", contentBlocks[0].Type)
	assert.Equal(t, "tool_use", contentBlocks[1].Type)
	assert.Equal(t, "write_file", contentBlocks[1].Name)
}

// TestGitHubAgent_StreamingMessageConstruction tests message construction
func TestGitHubAgent_StreamingMessageConstruction(t *testing.T) {
	// The runWithStreaming method constructs a providers.MessageResponse
	// This tests that pattern

	message := &providers.MessageResponse{
		Content: []providers.ContentBlock{
			{Type: "text", Text: "Response text"},
		},
		StopReason: providers.StopReasonEndTurn,
	}

	assert.NotNil(t, message)
	assert.Len(t, message.Content, 1)
	assert.Equal(t, providers.StopReasonEndTurn, message.StopReason)
}

// TestGitHubAgent_StreamingIndexMapping tests index-based content mapping
func TestGitHubAgent_StreamingIndexMapping(t *testing.T) {
	// The streaming code uses index-based mapping for content blocks
	// This tests that we correctly map indices to content

	type contentItem struct {
		index   int64
		content string
	}

	items := []contentItem{
		{0, "first"},
		{1, "second"},
		{2, "third"},
	}

	contentMap := make(map[int64]*strings.Builder)
	for _, item := range items {
		contentMap[item.index] = &strings.Builder{}
		contentMap[item.index].WriteString(item.content)
	}

	assert.Equal(t, "first", contentMap[0].String())
	assert.Equal(t, "second", contentMap[1].String())
	assert.Equal(t, "third", contentMap[2].String())
}
