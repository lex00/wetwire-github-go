// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// EnvironmentPromotion is the environment promotion workflow.
// It demonstrates a three-stage promotion pipeline:
// deploy dev -> promote staging -> promote production.
//
// Each environment can have its own approval rules configured
// through GitHub environment protection rules:
// - dev: No approval (automatic deployment)
// - staging: Optional approval (configurable)
// - production: Required approval (recommended)
var EnvironmentPromotion = workflow.Workflow{
	Name: "Environment Promotion",
	On:   PromotionTriggers,
	Jobs: map[string]workflow.Job{
		"deploy-dev":         DeployDev,
		"promote-staging":    PromoteStaging,
		"promote-production": PromoteProduction,
	},
}
