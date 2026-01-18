package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var BuildPush = workflow.PushTrigger{
	Branches: []string{"main", "v[0-9]+.[0-9]+", "releng/**", "tsccr-auto-pinning/**", "dependabot/**"},
	Tags:     []string{"v[0-9]+.[0-9]+.[0-9]+*"},
}

var BuildTriggers = workflow.Triggers{
	Push: &BuildPush,
}
