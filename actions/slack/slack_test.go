package slack

import (
	"testing"

	"github.com/lex00/wetwire-github-go/workflow"
)

func TestSlack_Action(t *testing.T) {
	s := Slack{}
	if got := s.Action(); got != "slackapi/slack-github-action@v1" {
		t.Errorf("Action() = %q, want %q", got, "slackapi/slack-github-action@v1")
	}
}

func TestSlack_Inputs(t *testing.T) {
	s := Slack{
		ChannelID:    "C1234567890",
		SlackMessage: "Deployment successful",
	}

	inputs := s.Inputs()

	if inputs["channel-id"] != "C1234567890" {
		t.Errorf("inputs[channel-id] = %v, want %q", inputs["channel-id"], "C1234567890")
	}

	if inputs["slack-message"] != "Deployment successful" {
		t.Errorf("inputs[slack-message] = %v, want %q", inputs["slack-message"], "Deployment successful")
	}
}

func TestSlack_Inputs_Empty(t *testing.T) {
	s := Slack{}
	inputs := s.Inputs()

	// Empty slack should have no inputs
	if len(inputs) != 0 {
		t.Errorf("empty Slack.Inputs() has %d entries, want 0", len(inputs))
	}
}

func TestSlack_Inputs_Payload(t *testing.T) {
	s := Slack{
		Payload: `{"text": "Hello World"}`,
	}

	inputs := s.Inputs()

	if inputs["payload"] != `{"text": "Hello World"}` {
		t.Errorf("inputs[payload] = %v, want %q", inputs["payload"], `{"text": "Hello World"}`)
	}
}

func TestSlack_Inputs_PayloadDelimiter(t *testing.T) {
	s := Slack{
		Payload:          `{"nested": {"key": "value"}}`,
		PayloadDelimiter: "_",
	}

	inputs := s.Inputs()

	if inputs["payload-delimiter"] != "_" {
		t.Errorf("inputs[payload-delimiter] = %v, want %q", inputs["payload-delimiter"], "_")
	}
}

func TestSlack_Inputs_PayloadFilePath(t *testing.T) {
	s := Slack{
		PayloadFilePath: "slack-payload.json",
	}

	inputs := s.Inputs()

	if inputs["payload-file-path"] != "slack-payload.json" {
		t.Errorf("inputs[payload-file-path] = %v, want %q", inputs["payload-file-path"], "slack-payload.json")
	}
}

func TestSlack_Inputs_PayloadFilePathParsed(t *testing.T) {
	s := Slack{
		PayloadFilePath:       "slack-payload.json",
		PayloadFilePathParsed: true,
	}

	inputs := s.Inputs()

	if inputs["payload-file-path-parsed"] != true {
		t.Errorf("inputs[payload-file-path-parsed] = %v, want true", inputs["payload-file-path-parsed"])
	}
}

func TestSlack_Inputs_PayloadFilePathParsedFalse(t *testing.T) {
	// Test that false boolean values are not included
	s := Slack{
		PayloadFilePath:       "slack-payload.json",
		PayloadFilePathParsed: false,
	}

	inputs := s.Inputs()

	// Should only have payload-file-path
	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1. Got: %v", len(inputs), inputs)
	}

	if _, exists := inputs["payload-file-path-parsed"]; exists {
		t.Errorf("inputs[payload-file-path-parsed] should not exist for false value")
	}
}

func TestSlack_Inputs_UpdateTs(t *testing.T) {
	s := Slack{
		ChannelID:    "C1234567890",
		SlackMessage: "Updated message",
		UpdateTs:     "1234567890.123456",
	}

	inputs := s.Inputs()

	if inputs["update-ts"] != "1234567890.123456" {
		t.Errorf("inputs[update-ts] = %v, want %q", inputs["update-ts"], "1234567890.123456")
	}
}

func TestSlack_Inputs_AllFields(t *testing.T) {
	s := Slack{
		ChannelID:             "C1234567890",
		SlackMessage:          "Test message",
		Payload:               `{"key": "value"}`,
		PayloadDelimiter:      "-",
		PayloadFilePath:       "payload.json",
		PayloadFilePathParsed: true,
		UpdateTs:              "1234567890.123456",
	}

	inputs := s.Inputs()

	expected := map[string]any{
		"channel-id":               "C1234567890",
		"slack-message":            "Test message",
		"payload":                  `{"key": "value"}`,
		"payload-delimiter":        "-",
		"payload-file-path":        "payload.json",
		"payload-file-path-parsed": true,
		"update-ts":                "1234567890.123456",
	}

	if len(inputs) != len(expected) {
		t.Errorf("inputs has %d entries, want %d", len(inputs), len(expected))
	}

	for key, want := range expected {
		if got := inputs[key]; got != want {
			t.Errorf("inputs[%q] = %v, want %v", key, got, want)
		}
	}
}

func TestSlack_ImplementsStepAction(t *testing.T) {
	s := Slack{}
	// Verify Slack implements StepAction interface
	var _ workflow.StepAction = s
}

func TestSlack_Inputs_ChannelIDOnly(t *testing.T) {
	s := Slack{
		ChannelID: "C1234567890",
	}

	inputs := s.Inputs()

	if inputs["channel-id"] != "C1234567890" {
		t.Errorf("inputs[channel-id] = %v, want %q", inputs["channel-id"], "C1234567890")
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestSlack_Inputs_MessageOnly(t *testing.T) {
	s := Slack{
		SlackMessage: "Test notification",
	}

	inputs := s.Inputs()

	if inputs["slack-message"] != "Test notification" {
		t.Errorf("inputs[slack-message] = %v, want %q", inputs["slack-message"], "Test notification")
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}

func TestSlack_Inputs_PayloadOnly(t *testing.T) {
	s := Slack{
		Payload: `{"text": "webhook payload"}`,
	}

	inputs := s.Inputs()

	if inputs["payload"] != `{"text": "webhook payload"}` {
		t.Errorf("inputs[payload] = %v, want %q", inputs["payload"], `{"text": "webhook payload"}`)
	}

	if len(inputs) != 1 {
		t.Errorf("inputs has %d entries, want 1", len(inputs))
	}
}
