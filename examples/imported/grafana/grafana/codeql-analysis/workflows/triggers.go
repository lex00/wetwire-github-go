package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CodeQLChecksPush = workflow.PushTrigger{
	Branches: []string{"main", "release-*.*.*"},
}

var CodeQLChecksTriggers = workflow.Triggers{
	Push: &CodeQLChecksPush,
	Schedule: []workflow.ScheduleTrigger{
		{Cron: "0 4 * * 6"},
	},
}
