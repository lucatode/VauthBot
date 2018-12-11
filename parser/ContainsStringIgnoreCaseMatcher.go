package parser

import "strings"

type ContainsIgnoreCaseMatcher struct {
	delegate   func(Message) (bool, string)
	dictionary map[string]string
}

func (em ContainsIgnoreCaseMatcher) ParseMessage(message Message) (bool, string) {
	lower := strings.ToLower(message.Text)
	val, ok := em.dictionary[lower]
	return ok, val
}

func NewContainsIgnoreCaseMatcher(dict map[string]string) Parser {
	return ContainsIgnoreCaseMatcher{ nil, dict}
}

func ContainsIgnoreCaseMatcherDecorated(dict map[string]string, matcher Parser) Parser {
	return ContainsIgnoreCaseMatcher{matcher.ParseMessage, dict}
}
