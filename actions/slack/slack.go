// Package slack provides a typed wrapper for slackapi/slack-github-action.
package slack

// Slack wraps the slackapi/slack-github-action@v1 action.
// Send notifications to Slack from your GitHub Actions workflow.
type Slack struct {
	// Slack channel ID where message will be posted. Needed if using bot token
	ChannelID string `yaml:"channel-id,omitempty"`

	// Message to post into Slack. Needed if using bot token
	SlackMessage string `yaml:"slack-message,omitempty"`

	// JSON payload to send to Slack if webhook route
	Payload string `yaml:"payload,omitempty"`

	// Custom delimiter used to flatten nested values in the JSON payload
	PayloadDelimiter string `yaml:"payload-delimiter,omitempty"`

	// Path to JSON payload file for webhook route
	PayloadFilePath string `yaml:"payload-file-path,omitempty"`

	// Replace templated variables in payload file with values from GitHub context
	PayloadFilePathParsed bool `yaml:"payload-file-path-parsed,omitempty"`

	// The timestamp of a previous message to update instead of posting new
	UpdateTs string `yaml:"update-ts,omitempty"`
}

// Action returns the action reference.
func (a Slack) Action() string {
	return "slackapi/slack-github-action@v1"
}

// Inputs returns the action inputs as a map.
func (a Slack) Inputs() map[string]any {
	with := make(map[string]any)

	if a.ChannelID != "" {
		with["channel-id"] = a.ChannelID
	}
	if a.SlackMessage != "" {
		with["slack-message"] = a.SlackMessage
	}
	if a.Payload != "" {
		with["payload"] = a.Payload
	}
	if a.PayloadDelimiter != "" {
		with["payload-delimiter"] = a.PayloadDelimiter
	}
	if a.PayloadFilePath != "" {
		with["payload-file-path"] = a.PayloadFilePath
	}
	if a.PayloadFilePathParsed {
		with["payload-file-path-parsed"] = a.PayloadFilePathParsed
	}
	if a.UpdateTs != "" {
		with["update-ts"] = a.UpdateTs
	}

	return with
}
