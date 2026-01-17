package serialize_test

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-github-go/internal/serialize"
	"github.com/lex00/wetwire-github-go/workflow"
)

func TestExpressionInEnv(t *testing.T) {
	w := &workflow.Workflow{
		Name: "CI",
		On: workflow.Triggers{
			Push: &workflow.PushTrigger{},
		},
		Jobs: map[string]workflow.Job{
			"deploy": {
				RunsOn: "ubuntu-latest",
				Steps: []any{
					workflow.Step{
						Run: "deploy.sh",
						Env: workflow.Env{
							"TOKEN": workflow.Secrets.Get("DEPLOY_TOKEN"),
						},
					},
				},
			},
		},
	}

	yaml, err := serialize.ToYAML(w)
	if err != nil {
		t.Fatalf("ToYAML failed: %v", err)
	}

	yamlStr := string(yaml)

	if !strings.Contains(yamlStr, "${{ secrets.DEPLOY_TOKEN }}") {
		t.Errorf("expected secret expression, got:\n%s", yamlStr)
	}
}

// TestExpressionSerialization tests various expression types and their serialization.
func TestExpressionSerialization(t *testing.T) {
	tests := []struct {
		name     string
		expr     workflow.Expression
		expected string
	}{
		{
			name:     "secrets expression",
			expr:     workflow.Secrets.Get("MY_SECRET"),
			expected: "secrets.MY_SECRET",
		},
		{
			name:     "github context",
			expr:     workflow.GitHub.SHA(),
			expected: "github.sha",
		},
		{
			name:     "matrix context",
			expr:     workflow.MatrixContext.Get("os"),
			expected: "matrix.os",
		},
		{
			name:     "combined expression",
			expr:     workflow.Branch("main").And(workflow.Success()),
			expected: "(github.ref == 'refs/heads/main') && (success())",
		},
		{
			name:     "negated expression",
			expr:     workflow.Failure().Not(),
			expected: "!(failure())",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &workflow.Workflow{
				Name: "Test",
				On: workflow.Triggers{
					Push: &workflow.PushTrigger{},
				},
				Jobs: map[string]workflow.Job{
					"test": {
						RunsOn: "ubuntu-latest",
						Steps: []any{
							workflow.Step{
								If:  tt.expr,
								Run: "echo test",
							},
						},
					},
				},
			}

			yaml, err := serialize.ToYAML(w)
			if err != nil {
				t.Fatalf("ToYAML failed: %v", err)
			}

			yamlStr := string(yaml)
			if !strings.Contains(yamlStr, tt.expected) {
				t.Errorf("expected %q in output, got:\n%s", tt.expected, yamlStr)
			}
		})
	}
}
