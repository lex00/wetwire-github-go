package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/dawidd6_download_artifact"
	"github.com/lex00/wetwire-github-go/workflow"
)

// DeploySteps are the steps for deploying after CI passes.
var DeploySteps = []any{
	// Display triggering workflow information
	workflow.Step{
		Name: "Show Workflow Run Info",
		Run: `echo "Triggered by workflow: ${{ github.event.workflow_run.name }}"
echo "Workflow run ID: ${{ github.event.workflow_run.id }}"
echo "Conclusion: ${{ github.event.workflow_run.conclusion }}"
echo "Head SHA: ${{ github.event.workflow_run.head_sha }}"
echo "Head branch: ${{ github.event.workflow_run.head_branch }}"`,
	},
	// Checkout the commit that triggered the CI workflow
	checkout.Checkout{
		Ref: "${{ github.event.workflow_run.head_sha }}",
	},
	// Download artifacts from the triggering workflow
	dawidd6_download_artifact.DownloadArtifact{
		GitHubToken:       "${{ secrets.GITHUB_TOKEN }}",
		RunID:             "${{ github.event.workflow_run.id }}",
		Name:              "build-artifacts",
		Path:              "./artifacts",
		IfNoArtifactFound: "warn",
	},
	// List downloaded artifacts
	workflow.Step{
		Name: "List Artifacts",
		Run:  "ls -la ./artifacts 2>/dev/null || echo 'No artifacts found'",
	},
	// Perform deployment
	workflow.Step{
		Name: "Deploy to Production",
		Run: `echo "Deploying commit ${{ github.event.workflow_run.head_sha }}"
echo "From branch: ${{ github.event.workflow_run.head_branch }}"
echo "Deployment successful!"`,
	},
}

// NotifySteps are the steps for sending notifications about workflow status.
var NotifySteps = []any{
	workflow.Step{
		Name: "Prepare Notification",
		ID:   "prepare",
		Run: `if [ "${{ github.event.workflow_run.conclusion }}" = "success" ]; then
  echo "status=passed" >> $GITHUB_OUTPUT
  echo "emoji=:white_check_mark:" >> $GITHUB_OUTPUT
else
  echo "status=failed" >> $GITHUB_OUTPUT
  echo "emoji=:x:" >> $GITHUB_OUTPUT
fi`,
	},
	workflow.Step{
		Name: "Send Notification",
		Run: `echo "Sending notification..."
echo "Workflow: ${{ github.event.workflow_run.name }}"
echo "Status: ${{ steps.prepare.outputs.status }}"
echo "Actor: ${{ github.event.workflow_run.actor.login }}"
echo "URL: ${{ github.event.workflow_run.html_url }}"`,
	},
}
