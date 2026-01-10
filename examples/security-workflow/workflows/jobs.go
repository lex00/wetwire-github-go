package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// CodeQLPermissions defines minimal permissions for CodeQL analysis.
// Demonstrates WAG017 compliance with explicit permission scopes.
var CodeQLPermissions = workflow.Permissions{
	Actions:        workflow.PermissionRead,
	Contents:       workflow.PermissionRead,
	SecurityEvents: workflow.PermissionWrite,
}

// CodeQL performs static code analysis using GitHub's CodeQL.
var CodeQL = workflow.Job{
	Name:        "CodeQL Analysis",
	RunsOn:      "ubuntu-latest",
	Permissions: &CodeQLPermissions,
	Steps:       CodeQLSteps,
}

// TrivyPermissions defines minimal permissions for container scanning.
// Demonstrates WAG017 compliance with explicit permission scopes.
var TrivyPermissions = workflow.Permissions{
	Contents:       workflow.PermissionRead,
	SecurityEvents: workflow.PermissionWrite,
}

// TrivyScan performs vulnerability scanning on the repository.
var TrivyScan = workflow.Job{
	Name:        "Trivy Security Scan",
	RunsOn:      "ubuntu-latest",
	Permissions: &TrivyPermissions,
	Steps:       TrivySteps,
}

// AttestPermissions defines minimal permissions for SLSA attestation.
// Demonstrates WAG017 compliance with explicit permission scopes.
var AttestPermissions = workflow.Permissions{
	Contents: workflow.PermissionRead,
	IDToken:  workflow.PermissionWrite,
}

// BuildAttest builds artifacts and generates SLSA provenance attestation.
var BuildAttest = workflow.Job{
	Name:        "Build with Attestation",
	RunsOn:      "ubuntu-latest",
	Permissions: &AttestPermissions,
	Steps:       BuildAttestSteps,
}
