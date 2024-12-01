package transformer

type ruleProcessor[T Rule] interface {
	Process(content string, rule T) (string, error)
	RuleName() string
}

type RuleProcessor interface {
	Process(content string, rule Rule) (string, error)
	RuleName() string
}

type genericProcessor[r Rule] struct {
	processor ruleProcessor[r]
}

func (g genericProcessor[r]) Process(content string, rule Rule) (string, error) {
	rle := rule.(r)

	return g.processor.Process(content, rle)
}

func (g genericProcessor[r]) RuleName() string {
	return g.processor.RuleName()
}

func NewProcessor[r Rule](processor ruleProcessor[r]) RuleProcessor {
	return &genericProcessor[r]{
		processor: processor,
	}
}
