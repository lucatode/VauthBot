package parser

import "strings"

type ExactIgnoreCaseMatcher struct {
	dictionary map[string]string
}

func (em ExactIgnoreCaseMatcher) ParseMessage(message Message) (bool, string) {
	lower := strings.ToLower(message.Text)
	val, ok := em.dictionary[lower]
	return ok, val
}

func NewExactIgnoreCaseMatcher (dict map[string]string)Parser {
	return ExactIgnoreCaseMatcher{dict}
}
