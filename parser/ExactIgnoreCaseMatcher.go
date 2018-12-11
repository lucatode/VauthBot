package parser

import "strings"

type ExactIgnoreCaseMatcher struct {
	delegate   func(Message) (bool, string)
	dictionary map[string]string
}

func (em ExactIgnoreCaseMatcher) ParseMessage(message Message) (bool, string) {
	lower := strings.ToLower(message.Text)
	val, ok := em.dictionary[lower]
	return ok, val
}

func NewExactIgnoreCaseMatcher(dict map[string]string) Parser {
	return ExactIgnoreCaseMatcher{dict, nil}
}

func ExactIgnoreCaseMatcherDecorated(dict map[string]string, matcher Parser) Parser {
	return ExactIgnoreCaseMatcher{matcher.ParseMessage, dict}
}
