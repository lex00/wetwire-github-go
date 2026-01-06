package workflow

import "fmt"

// Expression wraps a GitHub Actions expression string.
// When serialized to YAML, becomes ${{ expression }}.
type Expression string

// String returns the expression wrapped in ${{ }}.
func (e Expression) String() string {
	return fmt.Sprintf("${{ %s }}", string(e))
}

// Raw returns the raw expression without the ${{ }} wrapper.
func (e Expression) Raw() string {
	return string(e)
}

// And combines this expression with another using &&.
func (e Expression) And(other Expression) Expression {
	return Expression(fmt.Sprintf("(%s) && (%s)", e.Raw(), other.Raw()))
}

// Or combines this expression with another using ||.
func (e Expression) Or(other Expression) Expression {
	return Expression(fmt.Sprintf("(%s) || (%s)", e.Raw(), other.Raw()))
}

// Not negates this expression.
func (e Expression) Not() Expression {
	return Expression(fmt.Sprintf("!(%s)", e.Raw()))
}

// Context accessors for GitHub Actions expressions.
var (
	// GitHub provides access to github.* context.
	GitHub = githubContext{}

	// Runner provides access to runner.* context.
	Runner = runnerContext{}

	// Secrets provides access to secrets.* context.
	Secrets = secretsContext{}

	// MatrixContext provides access to matrix.* context.
	MatrixContext = matrixContext{}

	// Steps provides access to steps.* context.
	Steps = stepsContext{}

	// Needs provides access to needs.* context.
	Needs = needsContext{}

	// Inputs provides access to inputs.* context.
	Inputs = inputsContext{}

	// Vars provides access to vars.* context.
	Vars = varsContext{}

	// EnvContext provides access to env.* context.
	EnvContext = envContext{}
)

// githubContext provides typed access to github.* expressions.
type githubContext struct{}

func (githubContext) Ref() Expression          { return Expression("github.ref") }
func (githubContext) RefName() Expression      { return Expression("github.ref_name") }
func (githubContext) RefType() Expression      { return Expression("github.ref_type") }
func (githubContext) SHA() Expression          { return Expression("github.sha") }
func (githubContext) Actor() Expression        { return Expression("github.actor") }
func (githubContext) Repository() Expression   { return Expression("github.repository") }
func (githubContext) RepositoryOwner() Expression { return Expression("github.repository_owner") }
func (githubContext) EventName() Expression    { return Expression("github.event_name") }
func (githubContext) Workspace() Expression    { return Expression("github.workspace") }
func (githubContext) RunID() Expression        { return Expression("github.run_id") }
func (githubContext) RunNumber() Expression    { return Expression("github.run_number") }
func (githubContext) RunAttempt() Expression   { return Expression("github.run_attempt") }
func (githubContext) Job() Expression          { return Expression("github.job") }
func (githubContext) Token() Expression        { return Expression("github.token") }
func (githubContext) ServerURL() Expression    { return Expression("github.server_url") }
func (githubContext) APIURL() Expression       { return Expression("github.api_url") }
func (githubContext) GraphQLURL() Expression   { return Expression("github.graphql_url") }
func (githubContext) HeadRef() Expression      { return Expression("github.head_ref") }
func (githubContext) BaseRef() Expression      { return Expression("github.base_ref") }

// Event returns an expression for github.event.<path>.
func (githubContext) Event(path string) Expression {
	return Expression(fmt.Sprintf("github.event.%s", path))
}

// runnerContext provides typed access to runner.* expressions.
type runnerContext struct{}

func (runnerContext) OS() Expression      { return Expression("runner.os") }
func (runnerContext) Arch() Expression    { return Expression("runner.arch") }
func (runnerContext) Name() Expression    { return Expression("runner.name") }
func (runnerContext) Temp() Expression    { return Expression("runner.temp") }
func (runnerContext) ToolCache() Expression { return Expression("runner.tool_cache") }

// secretsContext provides typed access to secrets.* expressions.
type secretsContext struct{}

// Get returns an expression for secrets.<name>.
func (secretsContext) Get(name string) Expression {
	return Expression(fmt.Sprintf("secrets.%s", name))
}

// GITHUB_TOKEN returns the expression for the default GitHub token.
func (secretsContext) GITHUB_TOKEN() Expression {
	return Expression("secrets.GITHUB_TOKEN")
}

// matrixContext provides typed access to matrix.* expressions.
type matrixContext struct{}

// Get returns an expression for matrix.<name>.
func (matrixContext) Get(name string) Expression {
	return Expression(fmt.Sprintf("matrix.%s", name))
}

// stepsContext provides typed access to steps.* expressions.
type stepsContext struct{}

// Get returns an expression for steps.<id>.outputs.<name>.
func (stepsContext) Get(stepID, outputName string) Expression {
	return Expression(fmt.Sprintf("steps.%s.outputs.%s", stepID, outputName))
}

// Outcome returns an expression for steps.<id>.outcome.
func (stepsContext) Outcome(stepID string) Expression {
	return Expression(fmt.Sprintf("steps.%s.outcome", stepID))
}

// Conclusion returns an expression for steps.<id>.conclusion.
func (stepsContext) Conclusion(stepID string) Expression {
	return Expression(fmt.Sprintf("steps.%s.conclusion", stepID))
}

// needsContext provides typed access to needs.* expressions.
type needsContext struct{}

// Get returns an expression for needs.<job>.outputs.<name>.
func (needsContext) Get(jobID, outputName string) Expression {
	return Expression(fmt.Sprintf("needs.%s.outputs.%s", jobID, outputName))
}

// Result returns an expression for needs.<job>.result.
func (needsContext) Result(jobID string) Expression {
	return Expression(fmt.Sprintf("needs.%s.result", jobID))
}

// inputsContext provides typed access to inputs.* expressions.
type inputsContext struct{}

// Get returns an expression for inputs.<name>.
func (inputsContext) Get(name string) Expression {
	return Expression(fmt.Sprintf("inputs.%s", name))
}

// varsContext provides typed access to vars.* expressions.
type varsContext struct{}

// Get returns an expression for vars.<name>.
func (varsContext) Get(name string) Expression {
	return Expression(fmt.Sprintf("vars.%s", name))
}

// envContext provides typed access to env.* expressions.
type envContext struct{}

// Get returns an expression for env.<name>.
func (envContext) Get(name string) Expression {
	return Expression(fmt.Sprintf("env.%s", name))
}

// Condition builders for common workflow conditions.

// Always returns an expression that always evaluates to true.
func Always() Expression { return Expression("always()") }

// Failure returns an expression that is true when any previous step failed.
func Failure() Expression { return Expression("failure()") }

// Success returns an expression that is true when all previous steps succeeded.
func Success() Expression { return Expression("success()") }

// Cancelled returns an expression that is true when the workflow was cancelled.
func Cancelled() Expression { return Expression("cancelled()") }

// Branch returns an expression that checks if the ref is a specific branch.
func Branch(name string) Expression {
	return Expression(fmt.Sprintf("github.ref == 'refs/heads/%s'", name))
}

// Tag returns an expression that checks if the ref is a specific tag.
func Tag(name string) Expression {
	return Expression(fmt.Sprintf("github.ref == 'refs/tags/%s'", name))
}

// TagPrefix returns an expression that checks if the ref starts with a tag prefix.
func TagPrefix(prefix string) Expression {
	return Expression(fmt.Sprintf("startsWith(github.ref, 'refs/tags/%s')", prefix))
}

// Push returns an expression that checks if the event is a push.
func Push() Expression {
	return Expression("github.event_name == 'push'")
}

// PullRequest returns an expression that checks if the event is a pull_request.
func PullRequest() Expression {
	return Expression("github.event_name == 'pull_request'")
}

// Contains returns an expression that checks if a value contains a substring.
func Contains(haystack, needle Expression) Expression {
	return Expression(fmt.Sprintf("contains(%s, %s)", haystack, needle))
}

// StartsWith returns an expression that checks if a value starts with a prefix.
func StartsWith(value, prefix Expression) Expression {
	return Expression(fmt.Sprintf("startsWith(%s, %s)", value, prefix))
}

// EndsWith returns an expression that checks if a value ends with a suffix.
func EndsWith(value, suffix Expression) Expression {
	return Expression(fmt.Sprintf("endsWith(%s, %s)", value, suffix))
}

// Format returns an expression that formats a string with arguments.
func Format(formatStr string, args ...Expression) Expression {
	argStrs := make([]any, len(args))
	for i, arg := range args {
		argStrs[i] = string(arg)
	}
	return Expression(fmt.Sprintf("format('%s', %s)", formatStr, fmt.Sprint(argStrs...)))
}

// Join returns an expression that joins an array with a separator.
func Join(array Expression, separator string) Expression {
	return Expression(fmt.Sprintf("join(%s, '%s')", array, separator))
}

// ToJSON returns an expression that converts a value to JSON.
func ToJSON(value Expression) Expression {
	return Expression(fmt.Sprintf("toJSON(%s)", value))
}

// FromJSON returns an expression that parses JSON.
func FromJSON(json Expression) Expression {
	return Expression(fmt.Sprintf("fromJSON(%s)", json))
}
