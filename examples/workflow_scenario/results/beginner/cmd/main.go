package main

import "fmt"

func main() {
	fmt.Println("CI/CD Workflow Configuration")
	fmt.Println("============================")
	fmt.Println()
	fmt.Println("To generate GitHub Actions workflow files:")
	fmt.Println("  wetwire-github build .")
	fmt.Println()
	fmt.Println("To validate the configuration:")
	fmt.Println("  wetwire-github lint .")
	fmt.Println()
	fmt.Println("To validate generated YAML:")
	fmt.Println("  wetwire-github validate .")
	fmt.Println()
	fmt.Println("Workflow Features:")
	fmt.Println("- Build and test on Go 1.23 and 1.24")
	fmt.Println("- Test on Ubuntu and macOS")
	fmt.Println("- Linter checks with golangci-lint")
	fmt.Println("- Test coverage reporting")
	fmt.Println("- Auto-deploy to staging on main branch")
	fmt.Println("- Manual approval for production deployment")
}
