package transformer

type Rule interface {
	RuleName() string
}

type Service interface {
	ApplyRules(content string, rules []Rule) (string, error)
	RegisterProcessor(processor RuleProcessor)
}
