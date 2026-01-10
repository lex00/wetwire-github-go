// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// SecurityPermissions defines workflow-level permissions with minimal scopes.
// Demonstrates WAG017 compliance with explicit permission constraints.
var SecurityPermissions = workflow.Permissions{
	Contents:       workflow.PermissionRead,
	SecurityEvents: workflow.PermissionWrite,
	IDToken:        workflow.PermissionWrite,
}

// Security is the main security workflow with CodeQL, Trivy, and SLSA attestation.
// It runs on push and pull request to main branch.
// Demonstrates WAG018 compliance: uses pull_request (not pull_request_target)
// with checkout, which is safe as it runs in the context of the PR head.
var Security = workflow.Workflow{
	Name:        "Security",
	On:          SecurityTriggers,
	Permissions: &SecurityPermissions,
	Jobs: map[string]workflow.Job{
		"codeql":       CodeQL,
		"trivy":        TrivyScan,
		"build-attest": BuildAttest,
	},
}
