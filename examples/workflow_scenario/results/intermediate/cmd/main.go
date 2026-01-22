package main

import "fmt"

func main() {
	fmt.Println("CI/CD Workflow Configuration")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("To build the GitHub Actions workflow YAML:")
	fmt.Println("  wetwire-github build .")
	fmt.Println()
	fmt.Println("To validate the generated YAML:")
	fmt.Println("  wetwire-github validate")
	fmt.Println()
	fmt.Println("Workflow includes:")
	fmt.Println("  - Build job with matrix testing (Go 1.23/1.24, ubuntu/macos)")
	fmt.Println("  - Test job with coverage on ubuntu-latest")
	fmt.Println("  - Deploy staging (auto-deploy on main branch)")
	fmt.Println("  - Deploy production (manual approval required)")
}
