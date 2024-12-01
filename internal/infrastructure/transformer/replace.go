package transformer

import (
	"strings"

	"github.com/kolah/github-batch-updater/internal/domain/transformer"
)

type Replace struct {
	Match   string
	Replace string
}

type ReplaceRule struct {
	Replaces []Replace
}

func (r ReplaceRule) RuleName() string {
	return "replace"
}

type ReplaceProcessor struct{}

func NewReplaceProcessor() transformer.RuleProcessor {
	return transformer.NewProcessor(&ReplaceProcessor{})
}

func (r *ReplaceProcessor) Process(content string, rule ReplaceRule) (string, error) {
	for _, replace := range rule.Replaces {
		content = strings.ReplaceAll(content, replace.Match, replace.Replace)
	}

	return content, nil
}

func (r *ReplaceProcessor) RuleName() string {
	return ReplaceRule{}.RuleName()
}
