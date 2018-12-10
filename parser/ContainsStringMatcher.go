package parser

import "strings"

type ContainsStringMatcher struct {
	delegate         func(Message) (bool, string)
	containsWordDict map[string]string
}

func (cwm ContainsStringMatcher) ParseMessage(message Message) (bool, string) {
	inputString := message.Text
	for k, v := range cwm.containsWordDict {
		if strings.Contains(inputString, k) {
			return true, v
		}
	}
	return cwm.delegate(message)
}

func NewContainsStringMatcher(dict map[string]string) Parser {
	delegate := func(input Message) (bool, string) { return false, "" }
	return ContainsStringMatcher{delegate, dict}
}

func ContainsStringDecorated(dict map[string]string, matcher Parser) Parser {
	return ContainsStringMatcher{matcher.ParseMessage, dict}
}
