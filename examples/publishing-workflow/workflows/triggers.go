// Package workflows defines GitHub Actions workflow declarations.
package workflows

import "github.com/lex00/wetwire-github-go/workflow"

// PublishPush triggers on version tags (v*).
var PublishPush = workflow.PushTrigger{
	Tags: []string{"v*"},
}

// PublishTriggers activates on version tag pushes.
var PublishTriggers = workflow.Triggers{
	Push: &PublishPush,
}

// ReleasePublished triggers when a release is published.
var ReleasePublished = workflow.ReleaseTrigger{
	Types: []string{"published"},
}

// ReleaseTriggers activates when a GitHub release is published.
var ReleaseTriggers = workflow.Triggers{
	Release: &ReleasePublished,
}
