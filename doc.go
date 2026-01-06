// Package wetwire provides typed Go declarations for generating GitHub YAML configurations.
//
// wetwire-github-go is a synthesis library that generates GitHub Actions workflow YAML,
// Dependabot configuration, and Issue/Discussion templates from typed Go declarations.
//
// # The "No Parens" Pattern
//
// All declarations use struct literals — no function calls or registration:
//
//	var CIPush = workflow.PushTrigger{Branches: List("main")}
//	var CI = workflow.Workflow{Name: "CI", On: workflow.Triggers{Push: CIPush}}
//	var Build = workflow.Job{Name: "build", RunsOn: "ubuntu-latest", Steps: BuildSteps}
//
// # Generated Package Structure
//
// User projects declare workflows as Go variables using struct literals.
// The wetwire-github CLI discovers these declarations via AST parsing and
// generates the corresponding YAML output.
package wetwire
