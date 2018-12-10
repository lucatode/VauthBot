package parser

import (
	"strings"

	"vauthbot/dispatcher"
)

type CommandsMatcher struct {
	delegate  func(Message) (bool, string)
	dispatcher dispatcher.Dispatcher
}

func (cm CommandsMatcher) ParseMessage(message Message) (bool, string) {
	inputString := message.Text
	if !strings.Contains(inputString, "#") {
		return cm.delegate(message)
	}

	splittedMessage := strings.Split(inputString, " ")
	ok, f := cm.dispatcher.GetActionFunc(splittedMessage[0])
	if ok {

		return ok, f(splittedMessage, message.ChatId)
	}

	return cm.delegate(message)
}

func NewCommandsMatcher(dispatcher dispatcher.Dispatcher) Parser {
	delegate := func(input Message) (bool, string) { return false, "" }
	return CommandsMatcher{delegate, dispatcher}
}

func CommandsDecorated(dispatcher dispatcher.Dispatcher, matcher Parser) Parser {
	return CommandsMatcher{matcher.ParseMessage, dispatcher}
}

