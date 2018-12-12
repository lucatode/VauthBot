package parser

import "strings"

type ContainsIgnoreCaseMatcher struct {
	delegate   func(Message) (bool, string)
	dictionary map[string]string
}

func (cwm ContainsIgnoreCaseMatcher) ParseMessage(message Message) (bool, string) {
	inputString := strings.ToLower(message.Text)
	for k, v := range cwm.dictionary {
		if strings.Contains(inputString, k) {
			return true, v
		}
	}
	return cwm.delegate(message)
}

func NewContainsIgnoreCaseMatcher(dict map[string]string) Parser {
	return ContainsIgnoreCaseMatcher{ nil, dict}
}

func ContainsIgnoreCaseMatcherDecorated(dict map[string]string, matcher Parser) Parser {
	return ContainsIgnoreCaseMatcher{matcher.ParseMessage, dict}
}
