Repo: github.com/example/webapp. Create these files:

**expected/workflows/workflow.go:**
- CI workflow, triggers: push+PR to main
- Jobs: build, test, deploy-staging, deploy-production

**expected/workflows/build.go:**
- Matrix: go=[1.23,1.24], os=[ubuntu-latest,macos-latest]
- Steps: checkout@v4, setup-go@v5, cache modules, build

**expected/workflows/test.go:**
- Ubuntu-latest, coverage report
- Steps: checkout, setup-go, test with -race -coverprofile

**expected/workflows/deploy.go:**
- DeployStaging: needs=[build,test], if=main, env=staging
- DeployProduction: needs=[build,test], if=main, env=production (with URL), manual approval
