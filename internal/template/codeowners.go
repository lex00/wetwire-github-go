package template

import (
	"github.com/lex00/wetwire-github-go/codeowners"
	"github.com/lex00/wetwire-github-go/internal/discover"
	"github.com/lex00/wetwire-github-go/internal/runner"
	"github.com/lex00/wetwire-github-go/internal/serialize"
)

// CodeownersBuildResult contains the result of building Codeowners configs.
type CodeownersBuildResult struct {
	// Configs contains the assembled configs with text output
	Configs []BuiltCodeowners

	// Errors contains any non-fatal errors encountered
	Errors []string
}

// BuiltCodeowners represents a Codeowners config ready for output.
type BuiltCodeowners struct {
	// Name is the config variable name
	Name string

	// Config is the assembled Owners config
	Config *codeowners.Owners

	// Content is the serialized CODEOWNERS text
	Content []byte
}

// BuildCodeowners assembles Codeowners configs from discovery and extraction results.
func (b *Builder) BuildCodeowners(discovered *discover.CodeownersDiscoveryResult, extracted *runner.CodeownersExtractionResult) (*CodeownersBuildResult, error) {
	result := &CodeownersBuildResult{
		Configs: []BuiltCodeowners{},
		Errors:  []string{},
	}

	// Process each config
	for _, dc := range discovered.Configs {
		// Find the extracted config data
		var rules []serialize.ExtractedCodeownersRule
		var found bool
		for _, ec := range extracted.Configs {
			if ec.Name == dc.Name {
				// Convert runner rules to serialize rules
				rules = make([]serialize.ExtractedCodeownersRule, len(ec.Rules))
				for i, r := range ec.Rules {
					rules[i] = serialize.ExtractedCodeownersRule{
						Pattern: r.Pattern,
						Owners:  r.Owners,
						Comment: r.Comment,
					}
				}
				found = true
				break
			}
		}

		if !found {
			result.Errors = append(result.Errors, "config "+dc.Name+": extraction data not found")
			continue
		}

		// Create the Owners config
		cfg := &codeowners.Owners{
			Rules: make([]codeowners.Rule, len(rules)),
		}
		for i, r := range rules {
			cfg.Rules[i] = codeowners.Rule{
				Pattern: r.Pattern,
				Owners:  r.Owners,
				Comment: r.Comment,
			}
		}

		// Serialize to CODEOWNERS format
		content, err := serialize.CodeownersRulesToText(rules)
		if err != nil {
			result.Errors = append(result.Errors, "config "+dc.Name+": "+err.Error())
			continue
		}

		result.Configs = append(result.Configs, BuiltCodeowners{
			Name:    dc.Name,
			Config:  cfg,
			Content: content,
		})
	}

	return result, nil
}
