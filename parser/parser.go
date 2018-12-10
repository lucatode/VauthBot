package parser

type Message struct {
	Text   string
	ChatId string
}

type Parser interface {
	ParseMessage(Message) (bool, string)
}
