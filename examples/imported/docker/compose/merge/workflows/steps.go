package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BinImageSteps = []any{
	workflow.Step{
		Name: "Free disk space",
		Uses: "jlumbroso/free-disk-space@54081f138730dfa15788a46383842cd2f914a1be",
		With: map[string]any{
			"android":        true,
			"dotnet":         true,
			"haskell":        true,
			"large-packages": true,
			"swap-storage":   true,
		},
	},
	workflow.Step{
		Name: "Checkout",
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Name: "Login to DockerHub",
		Uses: "docker/login-action@v3",
		If:   "github.event_name != 'pull_request'",
		With: map[string]any{
			"password": "${{ secrets.DOCKERPUBLICBOT_WRITE_PAT }}",
			"username": "${{ secrets.DOCKERPUBLICBOT_USERNAME }}",
		},
	},
	workflow.Step{
		Name: "Set up QEMU",
		Uses: "docker/setup-qemu-action@v3",
	},
	workflow.Step{
		Name: "Set up Docker Buildx",
		Uses: "docker/setup-buildx-action@v3",
	},
	workflow.Step{
		ID:   "meta",
		Name: "Docker meta",
		Uses: "docker/metadata-action@v5",
		With: map[string]any{
			"bake-target": "meta-helper",
			"images":      "${{ env.REPO_SLUG }}\n",
			"tags":        "type=ref,event=tag\ntype=edge\n",
		},
	},
	workflow.Step{
		ID:   "bake",
		Name: "Build and push image",
		Uses: "docker/bake-action@v6",
		With: map[string]any{
			"files":      "./docker-bake.hcl\n${{ steps.meta.outputs.bake-file }}\n",
			"provenance": "mode=max",
			"push":       "${{ github.event_name != 'pull_request' }}",
			"sbom":       true,
			"set":        "*.cache-from=type=gha,scope=bin-image\n*.cache-to=type=gha,scope=bin-image,mode=max\n",
			"source":     ".",
			"targets":    "image-cross",
		},
	},
}

var DesktopEdgeTestSteps = []any{
	workflow.Step{
		ID:   "generate_token",
		Name: "Generate Token",
		Uses: "actions/create-github-app-token@v1",
		With: map[string]any{
			"app-id":       "${{ vars.DOCKERDESKTOP_APP_ID }}",
			"owner":        "docker",
			"private-key":  "${{ secrets.DOCKERDESKTOP_APP_PRIVATEKEY }}",
			"repositories": "${{ secrets.DOCKERDESKTOP_REPO }}\n",
		},
	},
	workflow.Step{
		Name: "Trigger Docker Desktop e2e with edge version",
		Uses: "actions/github-script@v7",
		With: map[string]any{
			"github-token": "${{ steps.generate_token.outputs.token }}",
			"script":       "await github.rest.actions.createWorkflowDispatch({\n  owner: 'docker',\n  repo: '${{ secrets.DOCKERDESKTOP_REPO }}',\n  workflow_id: 'compose-edge-integration.yml',\n  ref: 'main',\n  inputs: {\n    \"image-tag\": \"${{ needs.bin-image.outputs.digest }}\"\n  }\n})\n",
		},
	},
}

var E2eSteps = []any{
	workflow.Step{
		Uses: "actions/checkout@v4",
	},
	workflow.Step{
		Uses: "actions/setup-go@v6",
		With: map[string]any{
			"cache":           true,
			"check-latest":    true,
			"go-version-file": ".go-version",
		},
	},
	workflow.Step{
		Name: "List Docker resources on machine",
		Run: `docker ps --all
docker volume ls
docker network ls
docker image ls
`,
	},
	workflow.Step{
		Name: "Remove Docker resources on machine",
		Run: `docker kill $(docker ps -q)
docker rm -f $(docker ps -aq)
docker volume rm -f $(docker volume ls -q)
docker ps --all
`,
	},
	workflow.Step{
		Name: "Unit tests",
		Run:  "make test",
	},
	workflow.Step{
		Name: "Build binaries",
		Run: `make
`,
	},
	workflow.Step{
		Name: "Check arch of go compose binary",
		Run: `file ./bin/build/docker-compose
`,
		If: "${{ !contains(matrix.os, 'desktop-windows') }}",
	},
	workflow.Step{
		Name: "Test plugin mode",
		Run: `make e2e-compose
`,
		If: "${{ matrix.mode == 'plugin' }}",
	},
	workflow.Step{
		Name: "Test standalone mode",
		Run: `make e2e-compose-standalone
`,
		If: "${{ matrix.mode == 'standalone' }}",
	},
}
