package transformer

type Default struct {
	processors map[string]RuleProcessor
}

func NewDefault() *Default {
	return &Default{
		processors: make(map[string]RuleProcessor),
	}
}

func (s *Default) ApplyRules(content string, rules []Rule) (string, error) {
	modified := content

	for _, rule := range rules {
		processor, err := s.processorFor(rule)
		if err != nil {
			return content, err
		}

		modified, err = processor.Process(modified, rule)
		if err != nil {
			return content, err
		}
	}

	return modified, nil
}

func (s *Default) RegisterProcessor(processor RuleProcessor) {
	s.processors[processor.RuleName()] = processor
}

func (s *Default) processorFor(rule Rule) (RuleProcessor, error) {
	processor, ok := s.processors[rule.RuleName()]
	if !ok {
		return nil, ProcessorNotFound()
	}

	return processor, nil
}
