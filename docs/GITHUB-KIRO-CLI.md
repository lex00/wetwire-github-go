# Kiro CLI Integration

Use Kiro CLI with wetwire-github for AI-assisted GitHub Actions workflow design in corporate GitHub environments.

## Prerequisites

- Go 1.23+ installed
- Kiro CLI installed ([installation guide](https://kiro.dev/docs/cli/installation/))
- AWS Builder ID or GitHub/Google account (for Kiro authentication)

---

## Step 1: Install wetwire-github

### Option A: Using Go (recommended)

```bash
go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest
```

### Option B: Pre-built binaries

Download from [GitHub Releases](https://github.com/lex00/wetwire-github-go/releases):

```bash
# macOS (Apple Silicon)
curl -LO https://github.com/lex00/wetwire-github-go/releases/latest/download/wetwire-github-darwin-arm64
chmod +x wetwire-github-darwin-arm64
sudo mv wetwire-github-darwin-arm64 /usr/local/bin/wetwire-github

# macOS (Intel)
curl -LO https://github.com/lex00/wetwire-github-go/releases/latest/download/wetwire-github-darwin-amd64
chmod +x wetwire-github-darwin-amd64
sudo mv wetwire-github-darwin-amd64 /usr/local/bin/wetwire-github

# Linux (x86-64)
curl -LO https://github.com/lex00/wetwire-github-go/releases/latest/download/wetwire-github-linux-amd64
chmod +x wetwire-github-linux-amd64
sudo mv wetwire-github-linux-amd64 /usr/local/bin/wetwire-github
```

### Verify installation

```bash
wetwire-github --version
```

---

## Step 2: Install Kiro CLI

```bash
# Install Kiro CLI
curl -fsSL https://cli.kiro.dev/install | bash

# Verify installation
kiro-cli --version

# Sign in (opens browser)
kiro-cli login
```

---

## Step 3: Configure Kiro for wetwire-github

Run the design command with `--provider kiro` to auto-configure:

```bash
# Create a project directory
mkdir my-workflows && cd my-workflows

# Initialize Go module
go mod init my-workflows

# Run design with Kiro provider (auto-installs configs on first run)
wetwire-github design --provider kiro "Create a CI workflow for a Go project"
```

This automatically installs:

| File | Purpose |
|------|---------|
| `~/.kiro/agents/wetwire-github-runner.json` | Kiro agent configuration |
| `.kiro/mcp.json` | Project MCP server configuration |

### Manual configuration (optional)

The MCP server is provided as a subcommand `wetwire-github mcp`. If you prefer to configure manually:

**~/.kiro/agents/wetwire-github-runner.json:**
```json
{
  "name": "wetwire-github-runner",
  "description": "GitHub Actions workflow generator using wetwire-github",
  "prompt": "You are a GitHub Actions workflow design assistant...",
  "model": "claude-sonnet-4",
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-github",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  },
  "tools": ["*"]
}
```

**.kiro/mcp.json:**
```json
{
  "mcpServers": {
    "wetwire": {
      "command": "wetwire-github",
      "args": ["mcp"],
      "cwd": "/path/to/your/project"
    }
  }
}
```

> **Note:** The `cwd` field ensures MCP tools resolve paths correctly in your project directory. When using `wetwire-github design --provider kiro`, this is configured automatically.

---

## Step 4: Run Kiro with wetwire design

### Using the wetwire-github CLI

```bash
# Start Kiro design session
wetwire-github design --provider kiro "Create a CI workflow with Go matrix testing"
```

This launches Kiro CLI with the wetwire-github-runner agent and your prompt.

### Using Kiro CLI directly

```bash
# Start chat with wetwire-github-runner agent
kiro-cli chat --agent wetwire-github-runner

# Or with an initial prompt
kiro-cli chat --agent wetwire-github-runner "Create a release workflow with semantic versioning"
```

---

## Available MCP Tools

The wetwire-github MCP server exposes four tools to Kiro:

| Tool | Description | Example |
|------|-------------|---------|
| `wetwire_init` | Initialize a new project | `wetwire_init(path="./myapp")` |
| `wetwire_lint` | Lint code for issues | `wetwire_lint(path="./workflows/...")` |
| `wetwire_build` | Generate GitHub Actions workflows | `wetwire_build(path="./workflows/...", format="yaml")` |
| `wetwire_validate` | Validate workflows using actionlint | `wetwire_validate(path=".github/workflows/ci.yml")` |

---

## Example Session

```
$ wetwire-github design --provider kiro "Create a CI workflow with Go matrix testing"

Installed Kiro agent config: ~/.kiro/agents/wetwire-github-runner.json
Installed project MCP config: .kiro/mcp.json
Starting Kiro CLI design session...

> I'll help you create a CI workflow with Go matrix testing.

Let me initialize the project and create the workflow code.

[Calling wetwire_init...]
[Calling wetwire_lint...]
[Calling wetwire_build...]

I've created the following files:
- workflows.go

The CI workflow includes:
- Matrix testing for Go 1.22 and 1.23
- Ubuntu and macOS runners
- Build and test steps
- Code coverage reporting

Would you like me to add any additional configurations?
```

---

## Workflow

The Kiro agent follows this workflow:

1. **Explore** - Understand your requirements
2. **Plan** - Design the GitHub Actions architecture
3. **Implement** - Generate Go code using wetwire-github patterns
4. **Lint** - Run `wetwire_lint` to check for issues
5. **Build** - Run `wetwire_build` to generate GitHub Actions YAML

---

## Deploying Generated Workflows

After Kiro generates your workflow code:

```bash
# Build the GitHub Actions workflows
wetwire-github build ./workflows

# Commit and push the generated workflows
git add .github/workflows/*.yml
git commit -m "Add CI workflow"
git push
```

The workflows will automatically run according to their configured triggers (push, pull_request, etc.).

---

## Troubleshooting

### MCP server not found

```
Mcp error: -32002: No such file or directory
```

**Solution:** Ensure `wetwire-github` is in your PATH:

```bash
which wetwire-github

# If not found, add to PATH or reinstall
go install github.com/lex00/wetwire-github-go/cmd/wetwire-github@latest
```

### Kiro CLI not found

```
kiro-cli not found in PATH
```

**Solution:** Install Kiro CLI:

```bash
curl -fsSL https://cli.kiro.dev/install | bash
```

### Authentication issues

```
Error: Not authenticated
```

**Solution:** Sign in to Kiro:

```bash
kiro-cli login
```

---

## Known Limitations

### Automated Testing

When using `wetwire-github test --provider kiro`, tests run in non-interactive mode (`--no-interactive`). This means:

- The agent runs autonomously without waiting for user input
- Persona simulation is limited - all personas behave similarly
- The agent won't ask clarifying questions

For true persona simulation with multi-turn conversations, use the Anthropic provider:

```bash
wetwire-github test --provider anthropic --persona expert "Create a CI workflow"
```

### Interactive Design Mode

Interactive design mode (`wetwire-github design --provider kiro`) works fully as expected:

- Real-time conversation with the agent
- Agent can ask clarifying questions
- Lint loop executes as specified in the agent prompt

---

## See Also

- [CLI Reference](CLI.md) - Full wetwire-github CLI documentation
- [Quick Start](QUICK_START.md) - Getting started with wetwire-github
- [Kiro CLI Installation](https://kiro.dev/docs/cli/installation/) - Official installation guide
- [Kiro CLI Docs](https://kiro.dev/docs/cli/) - Official Kiro documentation
