// Package agent provides AI-assisted workflow generation using wetwire-core-go.
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/lex00/wetwire-core-go/agent/orchestrator"
	"github.com/lex00/wetwire-core-go/agent/results"
)

// GitHubAgent generates GitHub Actions workflows using the Anthropic API.
type GitHubAgent struct {
	client         anthropic.Client
	model          string
	session        *results.Session
	developer      orchestrator.Developer
	workDir        string
	generatedFiles []string
	maxLintCycles  int
	streamHandler  StreamHandler

	// Lint enforcement state
	lintCalled  bool
	lintPassed  bool
	pendingLint bool
	lintCycles  int
}

// StreamHandler is called for each text chunk during streaming.
type StreamHandler func(text string)

// Config configures the GitHubAgent.
type Config struct {
	// APIKey for Anthropic (defaults to ANTHROPIC_API_KEY env var)
	APIKey string

	// Model to use (defaults to claude-sonnet-4-20250514)
	Model string

	// WorkDir is the directory to write generated files
	WorkDir string

	// MaxLintCycles is the maximum number of lint/fix attempts
	MaxLintCycles int

	// Session for tracking results
	Session *results.Session

	// Developer to ask clarifying questions
	Developer orchestrator.Developer

	// StreamHandler is called for each text chunk during streaming
	StreamHandler StreamHandler
}

// NewGitHubAgent creates a new GitHubAgent.
func NewGitHubAgent(config Config) (*GitHubAgent, error) {
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	if config.WorkDir == "" {
		config.WorkDir = "."
	}
	if config.MaxLintCycles == 0 {
		config.MaxLintCycles = 5
	}

	model := config.Model
	if model == "" {
		model = string(anthropic.ModelClaudeSonnet4_20250514)
	}

	return &GitHubAgent{
		client:        client,
		model:         model,
		session:       config.Session,
		developer:     config.Developer,
		workDir:       config.WorkDir,
		maxLintCycles: config.MaxLintCycles,
		streamHandler: config.StreamHandler,
	}, nil
}

const systemPrompt = `You are a GitHub Actions workflow generator using the wetwire-github framework.
Your job is to generate Go code that defines GitHub Actions workflows.

The user will describe what CI/CD workflows they need. You will:
1. Ask clarifying questions if the requirements are unclear
2. Generate Go code using the wetwire-github patterns
3. Run the linter and fix any issues
4. Build the YAML output
5. Validate the generated YAML with actionlint

Use the wetwire-github patterns for all workflows:

    var CI = workflow.Workflow{
        Name: "CI",
        On:   CITriggers,
        Jobs: map[string]workflow.Job{
            "build": Build,
        },
    }

    var CITriggers = workflow.Triggers{
        Push:        &workflow.PushTrigger{Branches: []string{"main"}},
        PullRequest: &workflow.PullRequestTrigger{Branches: []string{"main"}},
    }

    var Build = workflow.Job{
        RunsOn: "ubuntu-latest",
        Steps:  BuildSteps,
    }

Use typed action wrappers instead of raw uses strings:
    checkout.Checkout{}
    setup_go.SetupGo{GoVersion: "1.23"}

Available tools:
- init_package: Create a new workflow project
- write_file: Write a Go file
- read_file: Read a file's contents
- run_lint: Run the wetwire-github linter
- run_build: Build the YAML workflows
- run_validate: Validate generated YAML with actionlint
- ask_developer: Ask the developer a clarifying question

Always run_lint after writing files, and fix any issues before running build.`

// Run executes the agent workflow.
func (a *GitHubAgent) Run(ctx context.Context, prompt string) error {
	tools := a.getTools()

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}

	// Agentic loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		params := anthropic.MessageNewParams{
			Model:     anthropic.Model(a.model),
			MaxTokens: 4096,
			System:    []anthropic.TextBlockParam{{Text: systemPrompt}},
			Messages:  messages,
			Tools:     tools,
		}

		var resp *anthropic.Message
		var err error

		if a.streamHandler != nil {
			resp, err = a.runWithStreaming(ctx, params)
		} else {
			resp, err = a.client.Messages.New(ctx, params)
		}
		if err != nil {
			return fmt.Errorf("API call failed: %w", err)
		}

		messages = append(messages, resp.ToParam())

		if resp.StopReason == anthropic.StopReasonEndTurn {
			if enforcement := a.checkCompletionGate(resp); enforcement != "" {
				messages = append(messages, anthropic.NewUserMessage(
					anthropic.NewTextBlock(enforcement),
				))
				continue
			}
			break
		}

		if resp.StopReason == anthropic.StopReasonToolUse {
			var toolResults []anthropic.ContentBlockParamUnion
			var toolsCalled []string

			for _, block := range resp.Content {
				if block.Type == "tool_use" {
					result := a.executeTool(ctx, block.Name, block.Input)
					toolResults = append(toolResults, anthropic.NewToolResultBlock(
						block.ID,
						result,
						false,
					))
					toolsCalled = append(toolsCalled, block.Name)
				}
			}

			messages = append(messages, anthropic.NewUserMessage(toolResults...))

			if enforcement := a.checkLintEnforcement(toolsCalled); enforcement != "" {
				messages = append(messages, anthropic.NewUserMessage(
					anthropic.NewTextBlock(enforcement),
				))
			}
		}
	}

	return nil
}

func (a *GitHubAgent) runWithStreaming(ctx context.Context, params anthropic.MessageNewParams) (*anthropic.Message, error) {
	stream := a.client.Messages.NewStreaming(ctx, params)

	var message *anthropic.Message
	var contentBlocks []anthropic.ContentBlockUnion
	currentTextContent := make(map[int64]*strings.Builder)
	currentToolInput := make(map[int64]*strings.Builder)

	for stream.Next() {
		event := stream.Current()

		switch event.Type {
		case "message_start":
			startEvent := event.AsMessageStart()
			message = &startEvent.Message
			contentBlocks = nil
			currentTextContent = make(map[int64]*strings.Builder)

		case "content_block_start":
			startEvent := event.AsContentBlockStart()
			if startEvent.ContentBlock.Type == "text" {
				currentTextContent[startEvent.Index] = &strings.Builder{}
			} else if startEvent.ContentBlock.Type == "tool_use" {
				currentToolInput[startEvent.Index] = &strings.Builder{}
			}
			block := anthropic.ContentBlockUnion{
				Type: startEvent.ContentBlock.Type,
				ID:   startEvent.ContentBlock.ID,
				Name: startEvent.ContentBlock.Name,
				Text: startEvent.ContentBlock.Text,
			}
			contentBlocks = append(contentBlocks, block)

		case "content_block_delta":
			deltaEvent := event.AsContentBlockDelta()
			if deltaEvent.Delta.Type == "text_delta" && deltaEvent.Delta.Text != "" {
				a.streamHandler(deltaEvent.Delta.Text)
				if builder, ok := currentTextContent[deltaEvent.Index]; ok {
					builder.WriteString(deltaEvent.Delta.Text)
				}
			}
			if deltaEvent.Delta.Type == "input_json_delta" && deltaEvent.Delta.PartialJSON != "" {
				if builder, ok := currentToolInput[deltaEvent.Index]; ok {
					builder.WriteString(deltaEvent.Delta.PartialJSON)
				}
			}

		case "content_block_stop":
			stopEvent := event.AsContentBlockStop()
			idx := int(stopEvent.Index)
			if idx < len(contentBlocks) {
				if builder, ok := currentTextContent[stopEvent.Index]; ok {
					contentBlocks[idx].Text = builder.String()
				}
				if builder, ok := currentToolInput[stopEvent.Index]; ok {
					contentBlocks[idx].Input = json.RawMessage(builder.String())
				}
			}

		case "message_delta":
			deltaEvent := event.AsMessageDelta()
			if message != nil {
				message.StopReason = deltaEvent.Delta.StopReason
				message.StopSequence = deltaEvent.Delta.StopSequence
			}
		}
	}

	if err := stream.Err(); err != nil {
		return nil, err
	}

	if message != nil {
		message.Content = contentBlocks
	}

	return message, nil
}

func (a *GitHubAgent) checkLintEnforcement(toolsCalled []string) string {
	wroteFile := false
	ranLint := false

	for _, tool := range toolsCalled {
		if tool == "write_file" {
			wroteFile = true
		}
		if tool == "run_lint" {
			ranLint = true
		}
	}

	if wroteFile && !ranLint {
		return `ENFORCEMENT: You wrote a file but did not call run_lint in the same turn.
You MUST call run_lint immediately after writing code to check for issues.
Call run_lint now before proceeding.`
	}

	return ""
}

func (a *GitHubAgent) checkCompletionGate(resp *anthropic.Message) string {
	var responseText string
	for _, block := range resp.Content {
		if block.Type == "text" {
			responseText += block.Text
		}
	}

	lowerText := strings.ToLower(responseText)
	isCompletionAttempt := strings.Contains(lowerText, "done") ||
		strings.Contains(lowerText, "complete") ||
		strings.Contains(lowerText, "finished") ||
		strings.Contains(lowerText, "that's it") ||
		strings.Contains(lowerText, "all set")

	if !isCompletionAttempt && len(a.generatedFiles) == 0 {
		return ""
	}

	if !a.lintCalled {
		return `ENFORCEMENT: You cannot complete without running the linter.
You MUST call run_lint to validate your code before finishing.
Call run_lint now.`
	}

	if a.pendingLint {
		return `ENFORCEMENT: You have written code since the last lint run.
You MUST call run_lint to validate your latest changes before finishing.
Call run_lint now.`
	}

	if !a.lintPassed {
		return `ENFORCEMENT: The linter found issues that have not been resolved.
You MUST fix the lint errors and run_lint again until it passes.
Review the lint output and fix the issues.`
	}

	return ""
}

func (a *GitHubAgent) getTools() []anthropic.ToolUnionParam {
	return []anthropic.ToolUnionParam{
		{
			OfTool: &anthropic.ToolParam{
				Name:        "init_package",
				Description: anthropic.String("Initialize a new wetwire-github workflow project"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"name": map[string]any{
							"type":        "string",
							"description": "Project name (directory name)",
						},
					},
					Required: []string{"name"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "write_file",
				Description: anthropic.String("Write content to a Go file"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "File path relative to work directory",
						},
						"content": map[string]any{
							"type":        "string",
							"description": "File content",
						},
					},
					Required: []string{"path", "content"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "read_file",
				Description: anthropic.String("Read a file's contents"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "File path relative to work directory",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "run_lint",
				Description: anthropic.String("Run the wetwire-github linter on the project"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Project path to lint",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "run_build",
				Description: anthropic.String("Build the YAML workflows from the Go project"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Project path to build",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "run_validate",
				Description: anthropic.String("Validate generated YAML with actionlint"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"path": map[string]any{
							"type":        "string",
							"description": "Path to YAML file or directory",
						},
					},
					Required: []string{"path"},
				},
			},
		},
		{
			OfTool: &anthropic.ToolParam{
				Name:        "ask_developer",
				Description: anthropic.String("Ask the developer a clarifying question"),
				InputSchema: anthropic.ToolInputSchemaParam{
					Properties: map[string]any{
						"question": map[string]any{
							"type":        "string",
							"description": "The question to ask",
						},
					},
					Required: []string{"question"},
				},
			},
		},
	}
}

func (a *GitHubAgent) executeTool(ctx context.Context, name string, input json.RawMessage) string {
	var params map[string]string
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Sprintf("Error parsing input: %v", err)
	}

	switch name {
	case "init_package":
		return a.toolInitPackage(params["name"])
	case "write_file":
		return a.toolWriteFile(params["path"], params["content"])
	case "read_file":
		return a.toolReadFile(params["path"])
	case "run_lint":
		return a.toolRunLint(params["path"])
	case "run_build":
		return a.toolRunBuild(params["path"])
	case "run_validate":
		return a.toolRunValidate(params["path"])
	case "ask_developer":
		answer, err := a.AskDeveloper(ctx, params["question"])
		if err != nil {
			return fmt.Sprintf("Error: %v", err)
		}
		return answer
	default:
		return fmt.Sprintf("Unknown tool: %s", name)
	}
}

func (a *GitHubAgent) toolInitPackage(name string) string {
	dir := filepath.Join(a.workDir, name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	// Create basic go.mod
	goModContent := fmt.Sprintf(`module github.com/example/%s

go 1.23

require github.com/lex00/wetwire-github-go v0.0.0
`, name)
	goModPath := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Sprintf("Error writing go.mod: %v", err)
	}

	return fmt.Sprintf("Created project directory: %s with go.mod", dir)
}

func (a *GitHubAgent) toolWriteFile(path, content string) string {
	fullPath := filepath.Join(a.workDir, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Sprintf("Error creating directory: %v", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return fmt.Sprintf("Error writing file: %v", err)
	}

	a.generatedFiles = append(a.generatedFiles, path)
	a.pendingLint = true
	a.lintPassed = false

	return fmt.Sprintf("Wrote %d bytes to %s", len(content), path)
}

func (a *GitHubAgent) toolReadFile(path string) string {
	fullPath := filepath.Join(a.workDir, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	return string(content)
}

func (a *GitHubAgent) toolRunLint(path string) string {
	fullPath := filepath.Join(a.workDir, path)
	cmd := exec.Command("wetwire-github", "lint", fullPath, "--format", "json")
	output, err := cmd.CombinedOutput()

	result := string(output)

	a.lintCalled = true
	a.pendingLint = false
	a.lintCycles++

	if err != nil {
		a.lintPassed = false
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 2 {
			var lintResult struct {
				Success bool `json:"success"`
				Issues  []struct {
					Message string `json:"message"`
				} `json:"issues"`
			}
			if json.Unmarshal(output, &lintResult) == nil && a.session != nil {
				issues := make([]string, len(lintResult.Issues))
				for i, issue := range lintResult.Issues {
					issues[i] = issue.Message
				}
				a.session.AddLintCycle(issues, a.lintCycles, false)
			}
		}
	} else {
		a.lintPassed = true
		if a.session != nil {
			a.session.AddLintCycle(nil, a.lintCycles, true)
		}
	}

	return result
}

func (a *GitHubAgent) toolRunBuild(path string) string {
	fullPath := filepath.Join(a.workDir, path)
	cmd := exec.Command("wetwire-github", "build", fullPath, "--format", "yaml")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Build error: %s\n%s", err, output)
	}
	return string(output)
}

func (a *GitHubAgent) toolRunValidate(path string) string {
	fullPath := filepath.Join(a.workDir, path)
	cmd := exec.Command("wetwire-github", "validate", fullPath, "--format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Validation issues:\n%s", output)
	}
	return string(output)
}

// AskDeveloper sends a question to the Developer.
func (a *GitHubAgent) AskDeveloper(ctx context.Context, question string) (string, error) {
	if a.developer == nil {
		return "", fmt.Errorf("no developer configured")
	}

	answer, err := a.developer.Respond(ctx, question)
	if err != nil {
		return "", err
	}

	if a.session != nil {
		a.session.AddQuestion(question, answer)
	}

	return answer, nil
}

// GetGeneratedFiles returns the list of generated file paths.
func (a *GitHubAgent) GetGeneratedFiles() []string {
	return a.generatedFiles
}

// GetLintCycles returns the number of lint attempts.
func (a *GitHubAgent) GetLintCycles() int {
	return a.lintCycles
}

// LintPassed returns whether the last lint run passed.
func (a *GitHubAgent) LintPassed() bool {
	return a.lintPassed
}
