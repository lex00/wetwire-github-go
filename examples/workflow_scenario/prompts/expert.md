Generate workflow: .github/workflows/ci.yml

- name: CI/CD
- triggers: push+PR to main
- jobs: build, test, deploy-staging, deploy-production

build job:
- matrix: go=[1.23,1.24], os=[ubuntu-latest,macos-latest]
- steps: checkout@v4, setup-go@v5, cache modules, go build ./...

test job:
- runs-on: ubuntu-latest
- steps: checkout, setup-go, go test -race -coverprofile=coverage.out ./...

deploy-staging job:
- needs: [build, test], if: github.ref == 'refs/heads/main'
- environment: staging
- run: echo "Deploying to staging"

deploy-production job:
- needs: [build, test], if: github.ref == 'refs/heads/main'
- environment: production (with URL: https://example.com)
- run: echo "Deploying to production"
