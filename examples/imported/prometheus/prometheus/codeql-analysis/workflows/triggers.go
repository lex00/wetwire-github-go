package workflows

import (
	"github.com/lex00/wetwire-github-go/workflow"
)

var CodeQLTriggers = workflow.Triggers{
	Schedule: []workflow.ScheduleTrigger{
		{Cron: "26 14 * * 1"},
	},
}
