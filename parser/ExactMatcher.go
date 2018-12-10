package parser

type ExactMatcher struct {
	exactMatchDict map[string]string
}

func (em ExactMatcher) ParseMessage(message Message) (bool, string) {
	val, ok := em.exactMatchDict[message.Text]
	return ok, val
}

func NewExactMatcher(dict map[string]string) Parser {
	return ExactMatcher{dict}
}
