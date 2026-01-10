package workflows

import (
	"github.com/lex00/wetwire-github-go/actions/checkout"
	"github.com/lex00/wetwire-github-go/actions/docker_build_push"
	"github.com/lex00/wetwire-github-go/actions/docker_login"
	"github.com/lex00/wetwire-github-go/actions/docker_setup_buildx"
	"github.com/lex00/wetwire-github-go/workflow"
)

// GHCRLogin logs into GitHub Container Registry.
var GHCRLogin = docker_login.DockerLogin{
	Registry: "ghcr.io",
	Username: "${{ github.actor }}",
	Password: "${{ secrets.GITHUB_TOKEN }}",
}

// DockerBuild builds the Docker image without pushing.
var DockerBuild = docker_build_push.DockerBuildPush{
	Context:   ".",
	Push:      false,
	Load:      true,
	Tags:      "ghcr.io/${{ github.repository }}:${{ github.sha }}",
	CacheFrom: "type=gha",
	CacheTo:   "type=gha,mode=max",
}

// DockerPushStep pushes the image to GHCR (only on main branch).
var DockerPushStep = workflow.Step{
	Name: "Push to GHCR",
	If:   "github.ref == 'refs/heads/main'",
	Uses: "docker/build-push-action@v6",
	With: map[string]any{
		"context":    ".",
		"push":       true,
		"tags":       "ghcr.io/${{ github.repository }}:latest\nghcr.io/${{ github.repository }}:${{ github.sha }}",
		"cache-from": "type=gha",
		"cache-to":   "type=gha,mode=max",
	},
}

// LoginStep logs into GHCR (only on main branch).
var LoginStep = workflow.Step{
	Name: "Login to GHCR",
	If:   "github.ref == 'refs/heads/main'",
	Uses: GHCRLogin.Action(),
	With: GHCRLogin.Inputs(),
}

// BuildSteps are the steps for the Docker build job.
var BuildSteps = []any{
	checkout.Checkout{},
	docker_setup_buildx.DockerSetupBuildx{},
	LoginStep,
	DockerBuild,
	DockerPushStep,
}
