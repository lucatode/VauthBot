package dispatcher

type Dispatcher interface {
	GetActionFunc(input string) (bool, func([]string, string) string)
}

type CommandDispatcher struct {
	ActionDictionary map[string]func([]string, string) string
}

func (cd CommandDispatcher) GetActionFunc(input string) (bool, func([]string, string) string) {
	f, ok := cd.ActionDictionary[input]
	return ok, f
}

func NewCommandDispatcher(dict map[string]func([]string, string) string) Dispatcher {
	return CommandDispatcher{dict}
}
