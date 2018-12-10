package parser

import (
	"vauthbot/dispatcher"
	"testing"

	"github.com/stretchr/testify/assert"
)

func MockExactMatchDictionary() map[string]string {
	return map[string]string{
		"notify": "notified",
	}
}

func MockWordMatcherDictionary() map[string]string {
	return map[string]string{
		"ab ba": "Found ab ba",
		"abba": "Found abba",
		"bccb": "Found bccb",
		"cddc": "Found cddc",
	}
}

func MockDispatcher() dispatcher.Dispatcher {
	return dispatcher.NewCommandDispatcher(map[string]func([]string, string) string{
		"#command": func([]string, string) string { return "command received" },
		"#subscribe_2": func(params []string, chatId string) string {
			msg := "subscribed"
			for i, p := range params {
				if i > 0 {
					msg = msg + " " + p
				}
			}
			return msg
		}})
}

var B = NewMessageBuilder()

////////////////////////
func TestExactMatchDecoratedWithCommand(t *testing.T) {

	matcher := CommandsDecorated(
		MockDispatcher(),
		NewExactMatcher(MockExactMatchDictionary()),
	)

	matchOutput, stringOutput := matcher.ParseMessage(B.WithText("notify").Build())
	assert.Equal(t, "notified", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = matcher.ParseMessage(B.WithText("#command").Build())
	assert.Equal(t, "command received", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")
}

func TestCheckCommandMatch(t *testing.T) {

	matchOutput, stringOutput := NewCommandsMatcher(MockDispatcher()).ParseMessage(B.WithText("#command").Build())
	assert.Equal(t, stringOutput, "command received", "")
	assert.Equal(t, matchOutput, true, "")
}

func TestCommandWithParameters(t *testing.T) {

	matchOutput, stringOutput := NewCommandsMatcher(MockDispatcher()).ParseMessage(B.WithText("#subscribe_2 channel xxxx").Build())
	assert.Equal(t, stringOutput, "subscribed channel xxxx", "")
	assert.Equal(t, matchOutput, true, "")
}

func TestCheckStringMatch(t *testing.T) {
	matchOutput, stringOutput := NewExactMatcher(MockExactMatchDictionary()).ParseMessage(B.WithText("notify").Build())
	assert.Equal(t, stringOutput, "notified", "")
	assert.Equal(t, matchOutput, true, "")
}

func TestContainsWordMatch(t *testing.T) {
	matchOutput, stringOutput := NewContainsStringMatcher(MockWordMatcherDictionary()).ParseMessage(B.WithText("ab ba abab").Build())
	assert.Equal(t, "Found ab ba", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = NewContainsStringMatcher(MockWordMatcherDictionary()).ParseMessage(B.WithText("abba abab ababa").Build())
	assert.Equal(t, "Found abba", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = NewContainsStringMatcher(MockWordMatcherDictionary()).ParseMessage(B.WithText("aaaa bccb").Build())
	assert.Equal(t, "Found bccb", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = NewContainsStringMatcher(MockWordMatcherDictionary()).ParseMessage(B.WithText("cddc").Build())
	assert.Equal(t, "Found cddc", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = NewContainsStringMatcher(MockWordMatcherDictionary()).ParseMessage(B.WithText("cccc").Build())
	assert.Equal(t, "", stringOutput, "")
	assert.Equal(t, matchOutput, false, "")
}

func TestExactIgnoreCase(t *testing.T){

	parser := NewExactIgnoreCaseMatcher(MockExactMatchDictionary())
	ok, val := parser.ParseMessage(B.WithText("Notify").Build())

	assert.Equal(t, true, ok, "")
	assert.Equal(t, "notified", val, "")
}

func TestExactMatchDecorated(t *testing.T) {

	matcher := ContainsStringDecorated(
		MockWordMatcherDictionary(),
		NewExactMatcher(MockExactMatchDictionary()),
	)
	matchOutput, stringOutput := matcher.ParseMessage(B.WithText("aaaa bccb").Build())
	assert.Equal(t, "Found bccb", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = ContainsStringDecorated(
		MockWordMatcherDictionary(),
		NewExactMatcher(MockExactMatchDictionary())).ParseMessage(B.WithText("notify").Build())
	assert.Equal(t, "notified", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")
}

func TestWithAllDecorations(t *testing.T) {

	matcher := CommandsDecorated(MockDispatcher(), ContainsStringDecorated(
		MockWordMatcherDictionary(),
		NewExactMatcher(MockExactMatchDictionary()),
	))
	matchOutput, stringOutput := matcher.ParseMessage(B.WithText("aaaa bccb").Build())
	assert.Equal(t, "Found bccb", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = matcher.ParseMessage(B.WithText("notify").Build())
	assert.Equal(t, "notified", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

	matchOutput, stringOutput = matcher.ParseMessage(B.WithText("#command").Build())
	assert.Equal(t, "command received", stringOutput, "")
	assert.Equal(t, matchOutput, true, "")

}

//////////////////////////////////

type Builder interface {
	Build() Message
}

type MessageBuilder struct {
	text   string
	chatId string
}

func (b MessageBuilder) WithText(t string) MessageBuilder   { b.text = t; return b }
func (b MessageBuilder) WithChatId(i string) MessageBuilder { b.chatId = i; return b }
func (b MessageBuilder) Build() Message                     { return Message{b.text, b.chatId} }

func NewMessageBuilder() MessageBuilder {
	return MessageBuilder{}
}
