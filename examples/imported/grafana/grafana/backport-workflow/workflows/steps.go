package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BackportSteps = []any{
	workflow.Step{
		ID:   "secrets",
		Name: "Get vault secrets",
		Uses: "grafana/shared-workflows/actions/get-vault-secrets@main",
		With: map[string]any{
			"export_env":   false,
			"repo_secrets": "APP_PEM=delivery-bot-app:PRIVATE_KEY\n",
		},
	},
	workflow.Step{
		ID:   "generate_token",
		Name: "Generate token",
		Uses: "tibdex/github-app-token@b62528385c34dbc9f38e5f4225ac829252d1ea92",
		With: map[string]any{
			"app_id":      "${{ vars.DELIVERY_BOT_APP_ID }}",
			"private_key": "${{ fromJSON(steps.secrets.outputs.secrets).APP_PEM }}",
		},
	},
	workflow.Step{
		ID:   "download-pr-info",
		Name: "Download PR info artifact",
		Uses: "actions/download-artifact@v6",
		With: map[string]any{
			"github-token": "${{ github.token }}",
			"name":         "pr_info",
			"run-id":       "${{ github.event.workflow_run.id }}",
		},
	},
	workflow.Step{
		ID:   "pr-info",
		Name: "Get PR info",
		Run:  "jq -r 'to_entries[] | select(.value | type != \"object\") | \"\\(.key)=\\(.value)\"' \"$PR_INFO_FILE\" >> \"$GITHUB_OUTPUT\"",
		Env: map[string]any{
			"PR_INFO_FILE": "${{ steps.download-pr-info.outputs.download-path }}/pr_info.json",
		},
	},
	workflow.Step{
		Name: "Print PR info",
		Run: `echo "PR action: $PR_ACTION"
echo "PR label: $PR_LABEL"
echo "PR number: $PR_NUMBER"
`,
		Env: map[string]any{
			"PR_ACTION": "${{ steps.pr-info.outputs.action }}",
			"PR_LABEL":  "${{ steps.pr-info.outputs.label }}",
			"PR_NUMBER": "${{ steps.pr-info.outputs.pr_number }}",
		},
	},
	workflow.Step{
		Name: "Checkout Grafana",
		Uses: "actions/checkout@v5",
		With: map[string]any{
			"fetch-depth":         2,
			"fetch-tags":          false,
			"persist-credentials": true,
			"ref":                 "${{ github.event.repository.default_branch }}",
			"token":               "${{ steps.generate_token.outputs.token }}",
		},
	},
	workflow.Step{
		Name: "Configure git user",
		Run: `git config --local user.name "github-actions[bot]"
git config --local user.email "github-actions[bot]@users.noreply.github.com"
git config --local --add --bool push.autoSetupRemote true
`,
	},
	workflow.Step{
		Name: "Run backport",
		Uses: "grafana/grafana-github-actions-go/backport@dev",
		With: map[string]any{
			"pr_label":   "${{ steps.pr-info.outputs.action == 'labeled' && steps.pr-info.outputs.label || '' }}",
			"pr_number":  "${{ steps.pr-info.outputs.pr_number }}",
			"repo_name":  "${{ github.event.repository.name }}",
			"repo_owner": "${{ github.repository_owner }}",
			"token":      "${{ steps.generate_token.outputs.token }}",
		},
	},
}
