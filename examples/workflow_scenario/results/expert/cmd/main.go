package main

// This file provides usage instructions for building the workflow.
//
// To generate the GitHub Actions workflow YAML:
//
//   wetwire-github build .
//
// This will output the workflow to:
//   .github/workflows/ci.yml
//
// The workflow includes:
// - Build job with matrix strategy (Go 1.23, 1.24 on ubuntu-latest, macos-latest)
// - Test job with race detection and coverage
// - Deploy-staging job (requires build and test, only on main branch)
// - Deploy-production job (requires build and test, only on main branch, with environment URL)
