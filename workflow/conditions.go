package workflow

// Condition represents a conditional expression that can be used in If fields.
// This interface is satisfied by Expression and string types.
type Condition interface {
	// condition is a marker method to identify condition types.
	condition()
}

// Ensure Expression implements Condition.
func (Expression) condition() {}

// StringCondition wraps a raw condition string.
type StringCondition string

func (StringCondition) condition() {}

// String returns the condition string.
func (s StringCondition) String() string {
	return string(s)
}

// ConditionBuilder provides fluent API for building complex conditions.
type ConditionBuilder struct {
	expr Expression
}

// NewCondition creates a new condition builder from an expression.
func NewCondition(expr Expression) *ConditionBuilder {
	return &ConditionBuilder{expr: expr}
}

// And adds an AND condition.
func (c *ConditionBuilder) And(other Expression) *ConditionBuilder {
	c.expr = c.expr.And(other)
	return c
}

// Or adds an OR condition.
func (c *ConditionBuilder) Or(other Expression) *ConditionBuilder {
	c.expr = c.expr.Or(other)
	return c
}

// Not negates the condition.
func (c *ConditionBuilder) Not() *ConditionBuilder {
	c.expr = c.expr.Not()
	return c
}

// Build returns the final Expression.
func (c *ConditionBuilder) Build() Expression {
	return c.expr
}

// Common condition patterns.

// OnMainBranch returns a condition that checks if running on the main branch.
func OnMainBranch() Expression {
	return Branch("main")
}

// OnDefaultBranch returns a condition that checks if running on the default branch.
func OnDefaultBranch() Expression {
	return Expression("github.ref == format('refs/heads/{0}', github.event.repository.default_branch)")
}

// IsPullRequest returns a condition that checks if the event is a pull request.
func IsPullRequest() Expression {
	return PullRequest()
}

// IsPush returns a condition that checks if the event is a push.
func IsPush() Expression {
	return Push()
}

// IsRelease returns a condition that checks if the event is a release.
func IsRelease() Expression {
	return Expression("github.event_name == 'release'")
}

// IsTag returns a condition that checks if the ref is a tag.
func IsTag() Expression {
	return Expression("startsWith(github.ref, 'refs/tags/')")
}

// PreviousJobSucceeded returns a condition checking if a previous job succeeded.
func PreviousJobSucceeded(jobID string) Expression {
	return Needs.Result(jobID).And(Expression("'success'"))
}

// PreviousJobFailed returns a condition checking if a previous job failed.
func PreviousJobFailed(jobID string) Expression {
	return Expression("needs." + jobID + ".result == 'failure'")
}
