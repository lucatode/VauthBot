package dispatcher

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewMockedCommandDispatcher() Dispatcher {
	return CommandDispatcher{map[string]func([]string, string) string{
		"#subscribe": func([]string, string) string { return "subscribed" },
		"#subscribe_2": func(params []string, chatId string) string {
			msg := "subscribed"
			for i, p := range params {
				if i > 0 {
					msg = msg + " " + p
				}
			}
			return msg
		},
	}}
}

func TestCommandMatched(t *testing.T) {
	d := NewMockedCommandDispatcher()
	ok, _ := d.GetActionFunc("#subscribe_2")
	assert.Equal(t, true, ok, "Getting specific match bool")
}

func TestGetActionFunc(t *testing.T) {
	d := NewMockedCommandDispatcher()
	_, f := d.GetActionFunc("#subscribe_2")
	input := "#subscribe_2 channel xxxx"
	splittedInput := strings.Split(input, " ")
	m := f(splittedInput, "chatId")
	assert.Equal(t, "subscribed channel xxxx", m, "Getting specific match bool")
}
