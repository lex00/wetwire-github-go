Create GitHub Actions CI/CD workflow for Go web application:

**Jobs:**
- Build: matrix test on Go 1.23/1.24, ubuntu/macos
- Test: run tests with coverage on ubuntu-latest
- DeployStaging: deploy to staging environment (runs on main branch)
- DeployProduction: deploy to production with manual approval (environment gate)

**Triggers:**
- Push to main
- Pull requests to main

**Requirements:**
- Deploy jobs need both build and test to pass
- Use actions/checkout@v4, actions/setup-go@v5
- Cache Go modules for performance
- Production requires environment approval
