---
title: "Examples"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

Real-world workflow patterns and examples for wetwire-github-go.

## Table of Contents

- [Basic CI Workflow](#basic-ci-workflow)
- [Multi-Language Matrix](#multi-language-matrix)
- [Docker Build and Push](#docker-build-and-push)
- [Release Workflow](#release-workflow)
- [Monorepo with Path Filters](#monorepo-with-path-filters)
- [Scheduled Maintenance](#scheduled-maintenance)
- [PR Labeling](#pr-labeling)
- [Deploy to Multiple Environments](#deploy-to-multiple-environments)

---

## Basic CI Workflow

A simple CI workflow that runs tests on push and pull request.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
)

var CI = Workflow{
    Name: "CI",
    On:   CITriggers,
    Jobs: Jobs{"build": Build},
}

var CITriggers = Triggers{
    Push:        PushTrigger{Branches: List("main")},
    PullRequest: PullRequestTrigger{Branches: List("main")},
}

var Build = Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Steps:  BuildSteps,
}

var BuildSteps = []any{
    checkout.Checkout{},
    setup_go.SetupGo{GoVersion: "1.23"},
    Step{Name: "Build", Run: "go build ./..."},
    Step{Name: "Test", Run: "go test -v ./..."},
}
```

### Generated YAML

```yaml
name: CI
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    name: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - name: Build
        run: go build ./...
      - name: Test
        run: go test -v ./...
```

---

## Multi-Language Matrix

Test across multiple Go versions and operating systems.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
)

var MatrixCI = Workflow{
    Name: "Matrix CI",
    On:   MatrixTriggers,
    Jobs: Jobs{"test": MatrixTest},
}

var MatrixTriggers = Triggers{
    Push:        PushTrigger{Branches: List("main")},
    PullRequest: PullRequestTrigger{},
}

var TestMatrix = Matrix{
    Values: map[string][]any{
        "go":   {"1.22", "1.23"},
        "os":   {"ubuntu-latest", "macos-latest", "windows-latest"},
    },
}

var MatrixTest = Job{
    Name:     "test",
    RunsOn:   Matrix.Get("os"),
    Strategy: Strategy{Matrix: TestMatrix},
    Steps:    MatrixSteps,
}

var MatrixSteps = []any{
    checkout.Checkout{},
    setup_go.SetupGo{GoVersion: Matrix.Get("go")},
    Step{Name: "Test", Run: "go test -v ./..."},
}
```

### Generated YAML

```yaml
name: Matrix CI
on:
  push:
    branches: [main]
  pull_request: {}

jobs:
  test:
    name: test
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        go: ["1.22", "1.23"]
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - name: Test
        run: go test -v ./...
```

---

## Docker Build and Push

Build and push a Docker image to GitHub Container Registry.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/docker_login"
    "github.com/lex00/wetwire-github-go/actions/docker_setup_buildx"
    "github.com/lex00/wetwire-github-go/actions/docker_build_push"
)

var DockerBuild = Workflow{
    Name: "Docker Build",
    On:   DockerTriggers,
    Jobs: Jobs{"build": DockerJob},
}

var DockerTriggers = Triggers{
    Push: PushTrigger{
        Branches: List("main"),
        Tags:     List("v*"),
    },
}

var DockerJob = Job{
    Name:   "build",
    RunsOn: "ubuntu-latest",
    Permissions: &Permissions{
        Contents: "read",
        Packages: "write",
    },
    Steps: DockerSteps,
}

var DockerSteps = []any{
    checkout.Checkout{},
    docker_setup_buildx.DockerSetupBuildx{},
    docker_login.DockerLogin{
        Registry: "ghcr.io",
        Username: GitHub.Actor,
        Password: Secrets.Get("GITHUB_TOKEN"),
    },
    docker_build_push.DockerBuildPush{
        Context:   ".",
        Push:      true,
        Tags:      "ghcr.io/${{ github.repository }}:latest",
        Platforms: "linux/amd64,linux/arm64",
        CacheFrom: "type=gha",
        CacheTo:   "type=gha,mode=max",
    },
}
```

### Generated YAML

```yaml
name: Docker Build
on:
  push:
    branches: [main]
    tags: ["v*"]

jobs:
  build:
    name: build
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: ghcr.io/${{ github.repository }}:latest
          platforms: linux/amd64,linux/arm64
          cache-from: type=gha
          cache-to: type=gha,mode=max
```

---

## Release Workflow

Create GitHub releases with changelog on tag push.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/gh_release"
)

var Release = Workflow{
    Name: "Release",
    On:   ReleaseTriggers,
    Jobs: Jobs{"release": ReleaseJob},
}

var ReleaseTriggers = Triggers{
    Push: PushTrigger{Tags: List("v*")},
}

var ReleaseJob = Job{
    Name:   "release",
    RunsOn: "ubuntu-latest",
    Permissions: &Permissions{
        Contents: "write",
    },
    Steps: ReleaseSteps,
}

var ReleaseSteps = []any{
    checkout.Checkout{FetchDepth: 0},
    Step{
        Name: "Build",
        Run:  "go build -o myapp ./cmd/myapp",
    },
    gh_release.GhRelease{
        Files:         "myapp",
        GenerateReleaseNotes: true,
    },
}
```

### Generated YAML

```yaml
name: Release
on:
  push:
    tags: ["v*"]

jobs:
  release:
    name: release
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - name: Build
        run: go build -o myapp ./cmd/myapp
      - uses: softprops/action-gh-release@v2
        with:
          files: myapp
          generate_release_notes: true
```

---

## Monorepo with Path Filters

Run different jobs based on which files changed.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
    "github.com/lex00/wetwire-github-go/actions/setup_go"
    "github.com/lex00/wetwire-github-go/actions/setup_node"
)

var Monorepo = Workflow{
    Name: "Monorepo CI",
    On:   MonorepoTriggers,
    Jobs: Jobs{
        "backend":  BackendJob,
        "frontend": FrontendJob,
    },
}

var MonorepoTriggers = Triggers{
    Push: PushTrigger{
        Branches: List("main"),
        Paths:    List("backend/**", "frontend/**"),
    },
    PullRequest: PullRequestTrigger{
        Paths: List("backend/**", "frontend/**"),
    },
}

var BackendJob = Job{
    Name:   "backend",
    RunsOn: "ubuntu-latest",
    Steps: []any{
        checkout.Checkout{},
        setup_go.SetupGo{GoVersion: "1.23"},
        Step{
            Name:             "Test Backend",
            Run:              "go test ./...",
            WorkingDirectory: "backend",
        },
    },
}

var FrontendJob = Job{
    Name:   "frontend",
    RunsOn: "ubuntu-latest",
    Steps: []any{
        checkout.Checkout{},
        setup_node.SetupNode{NodeVersion: "20"},
        Step{
            Name:             "Install",
            Run:              "npm ci",
            WorkingDirectory: "frontend",
        },
        Step{
            Name:             "Test Frontend",
            Run:              "npm test",
            WorkingDirectory: "frontend",
        },
    },
}
```

### Generated YAML

```yaml
name: Monorepo CI
on:
  push:
    branches: [main]
    paths: ["backend/**", "frontend/**"]
  pull_request:
    paths: ["backend/**", "frontend/**"]

jobs:
  backend:
    name: backend
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.23"
      - name: Test Backend
        run: go test ./...
        working-directory: backend

  frontend:
    name: frontend
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: "20"
      - name: Install
        run: npm ci
        working-directory: frontend
      - name: Test Frontend
        run: npm test
        working-directory: frontend
```

---

## Scheduled Maintenance

Run maintenance tasks on a schedule.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
)

var Maintenance = Workflow{
    Name: "Maintenance",
    On:   MaintenanceTriggers,
    Jobs: Jobs{"cleanup": CleanupJob},
}

var MaintenanceTriggers = Triggers{
    Schedule: List(
        ScheduleTrigger{Cron: "0 0 * * 0"},  // Weekly on Sunday
    ),
    WorkflowDispatch: &WorkflowDispatchTrigger{},  // Allow manual runs
}

var CleanupJob = Job{
    Name:   "cleanup",
    RunsOn: "ubuntu-latest",
    Steps: []any{
        checkout.Checkout{},
        Step{
            Name: "Clean old artifacts",
            Run: `
                gh api repos/${{ github.repository }}/actions/artifacts \
                  --jq '.artifacts[] | select(.created_at < (now - 30*24*60*60 | todate)) | .id' \
                  | xargs -I {} gh api -X DELETE repos/${{ github.repository }}/actions/artifacts/{}
            `,
            Env: Env{"GH_TOKEN": Secrets.Get("GITHUB_TOKEN")},
        },
    },
}
```

### Generated YAML

```yaml
name: Maintenance
on:
  schedule:
    - cron: "0 0 * * 0"
  workflow_dispatch: {}

jobs:
  cleanup:
    name: cleanup
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Clean old artifacts
        run: |
          gh api repos/${{ github.repository }}/actions/artifacts \
            --jq '.artifacts[] | select(.created_at < (now - 30*24*60*60 | todate)) | .id' \
            | xargs -I {} gh api -X DELETE repos/${{ github.repository }}/actions/artifacts/{}
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## PR Labeling

Automatically label PRs based on changed files.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/github_script"
)

var PRLabeler = Workflow{
    Name: "PR Labeler",
    On:   LabelerTriggers,
    Jobs: Jobs{"label": LabelJob},
}

var LabelerTriggers = Triggers{
    PullRequest: PullRequestTrigger{
        Types: List("opened", "synchronize"),
    },
}

var LabelJob = Job{
    Name:   "label",
    RunsOn: "ubuntu-latest",
    Permissions: &Permissions{
        Contents:     "read",
        PullRequests: "write",
    },
    Steps: []any{
        github_script.GithubScript{
            Script: `
                const { data: files } = await github.rest.pulls.listFiles({
                    owner: context.repo.owner,
                    repo: context.repo.repo,
                    pull_number: context.issue.number,
                });

                const labels = new Set();
                for (const file of files) {
                    if (file.filename.startsWith('docs/')) labels.add('documentation');
                    if (file.filename.endsWith('.go')) labels.add('go');
                    if (file.filename.endsWith('.ts')) labels.add('typescript');
                    if (file.filename.includes('test')) labels.add('tests');
                }

                if (labels.size > 0) {
                    await github.rest.issues.addLabels({
                        owner: context.repo.owner,
                        repo: context.repo.repo,
                        issue_number: context.issue.number,
                        labels: [...labels],
                    });
                }
            `,
        },
    },
}
```

### Generated YAML

```yaml
name: PR Labeler
on:
  pull_request:
    types: [opened, synchronize]

jobs:
  label:
    name: label
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: write
    steps:
      - uses: actions/github-script@v7
        with:
          script: |
            const { data: files } = await github.rest.pulls.listFiles({
                owner: context.repo.owner,
                repo: context.repo.repo,
                pull_number: context.issue.number,
            });

            const labels = new Set();
            for (const file of files) {
                if (file.filename.startsWith('docs/')) labels.add('documentation');
                if (file.filename.endsWith('.go')) labels.add('go');
                if (file.filename.endsWith('.ts')) labels.add('typescript');
                if (file.filename.includes('test')) labels.add('tests');
            }

            if (labels.size > 0) {
                await github.rest.issues.addLabels({
                    owner: context.repo.owner,
                    repo: context.repo.repo,
                    issue_number: context.issue.number,
                    labels: [...labels],
                });
            }
```

---

## Deploy to Multiple Environments

Deploy to staging on PR merge, production on tag.

### Go Source

```go
package workflows

import (
    . "github.com/lex00/wetwire-github-go/workflow"
    "github.com/lex00/wetwire-github-go/actions/checkout"
)

var Deploy = Workflow{
    Name: "Deploy",
    On:   DeployTriggers,
    Jobs: Jobs{
        "staging":    StagingDeploy,
        "production": ProductionDeploy,
    },
}

var DeployTriggers = Triggers{
    Push: PushTrigger{
        Branches: List("main"),
        Tags:     List("v*"),
    },
}

var StagingDeploy = Job{
    Name:   "staging",
    RunsOn: "ubuntu-latest",
    If:     Branch("main"),
    Environment: &JobEnvironment{
        Name: "staging",
        URL:  "https://staging.example.com",
    },
    Steps: []any{
        checkout.Checkout{},
        Step{
            Name: "Deploy to Staging",
            Run:  "echo 'Deploying to staging...'",
            Env: Env{
                "DEPLOY_ENV": "staging",
                "API_KEY":    Secrets.Get("STAGING_API_KEY"),
            },
        },
    },
}

var ProductionDeploy = Job{
    Name:   "production",
    RunsOn: "ubuntu-latest",
    If:     "startsWith(github.ref, 'refs/tags/v')",
    Environment: &JobEnvironment{
        Name: "production",
        URL:  "https://example.com",
    },
    Steps: []any{
        checkout.Checkout{},
        Step{
            Name: "Deploy to Production",
            Run:  "echo 'Deploying to production...'",
            Env: Env{
                "DEPLOY_ENV": "production",
                "API_KEY":    Secrets.Get("PRODUCTION_API_KEY"),
            },
        },
    },
}
```

### Generated YAML

```yaml
name: Deploy
on:
  push:
    branches: [main]
    tags: ["v*"]

jobs:
  staging:
    name: staging
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment:
      name: staging
      url: https://staging.example.com
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to Staging
        run: echo 'Deploying to staging...'
        env:
          DEPLOY_ENV: staging
          API_KEY: ${{ secrets.STAGING_API_KEY }}

  production:
    name: production
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/v')
    environment:
      name: production
      url: https://example.com
    steps:
      - uses: actions/checkout@v4
      - name: Deploy to Production
        run: echo 'Deploying to production...'
        env:
          DEPLOY_ENV: production
          API_KEY: ${{ secrets.PRODUCTION_API_KEY }}
```

---

## See Also

- [Quick Start](QUICK_START.md) - Getting started
- [CLI Reference](CLI.md) - CLI commands
- [Developer Guide](DEVELOPERS.md) - Contributing
