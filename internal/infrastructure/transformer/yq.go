package transformer

import (
	"os"

	"github.com/kolah/github-batch-updater/internal/domain/transformer"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

type YQRule struct {
	Expressions []string
}

func (r YQRule) RuleName() string {
	return "yq"
}

type YQProcessor struct {
	evaluator yqlib.StringEvaluator
}

func NewYQProcessor() transformer.RuleProcessor {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	levelBackend := logging.AddModuleLevel(backend)
	levelBackend.SetLevel(logging.CRITICAL, "")
	yqlib.GetLogger().SetBackend(levelBackend)

	evaluator := yqlib.NewStringEvaluator()
	processor := &YQProcessor{
		evaluator: evaluator,
	}
	return transformer.NewProcessor(processor)
}

func (p *YQProcessor) Process(content string, rule YQRule) (string, error) {
	format, err := yqlib.FormatFromString("yaml")
	if err != nil {
		return content, err
	}
	encoder := format.EncoderFactory()
	decoder := format.DecoderFactory()
	result := content
	for _, expression := range rule.Expressions {
		result, err = p.evaluator.Evaluate(expression, result, encoder, decoder)
		if err != nil {
			return content, transformer.ProcessorError().WrapError(err)
		}
	}

	return result, nil
}

func (p *YQProcessor) RuleName() string {
	return YQRule{}.RuleName()
}
