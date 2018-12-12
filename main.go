package main

import (
	"vauthbot/dispatcher"
	"vauthbot/initializer"
	"vauthbot/subscriber"
	"log"
	"net/http"
	"os"
	"strconv"

	"vauthbot/logger"
	"vauthbot/parser"
	"vauthbot/repositories"
	"gopkg.in/telegram-bot-api.v4"
	"strings"
	"io/ioutil"
	"encoding/json"
	"vauthbot/replacer"
)

func main() {
	//INIT
	init := Init()
	logger := CreateLogger(init)

	m := repositories.BuildDictionaries(init.GetFireBaseResponsesUrl(), logger)
	p := BuildParser(init, m)

	// SETUP BOT
	bot, err := tgbotapi.NewBotAPI(init.GetApiToken())
	if err != nil {
		log.Fatal(err)
	}

	// BOT CONFIG
	res, errWebhook := bot.SetWebhook(tgbotapi.NewWebhook(init.GetServerUrl() + bot.Token))
	if errWebhook != nil {
		log.Fatal(errWebhook)
	}
	logger.Log("MAIN", res.Description)

	// SETUP INPUT ROUTES
	port := os.Getenv("PORT")
	logger.Log("MAIN", "port: "+port)
	go http.ListenAndServe(":"+port, nil)
	http.HandleFunc("/notify/", NotifyHandler(init, bot))

	// FETCH MESSAGES
	updates := bot.ListenForWebhook("/" + bot.Token)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		ok, text := p.ParseMessage(BuildMessage(update.Message))

		placeholder := "%randomNumber"
		if strings.Contains(text, placeholder) {
			rnd := replacer.GetRandomRangeNumberReplacer(1000, placeholder, replacer.GenerateRandomNumeber)
			text = rnd.ReplaceIn(text)
		}

		if ok {
			logger.Log("Match hit", update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
			m := repositories.BuildDictionaries(init.GetFireBaseResponsesUrl(), logger)
			p = BuildParser(init, m)
		}
	}
}

func Init() initializer.Initializer {
	return initializer.NewInitializer(initializer.NewEnvReader())
}
func CreateLogger(init initializer.Initializer) logger.FirebaseLogger {
	logger := logger.FirebaseLogger{init.GetFireBaseLogsUrl()}
	logger.Log("MAIN", "Starting")
	return logger
}
func CreateRepository(logger logger.FirebaseLogger) repositories.FireBaseRepository {
	client := http.Client{}
	return repositories.FireBaseRepository{client.Get, logger, repositories.GetMatchCases }
}

func BuildCommandDispatcher(url string) dispatcher.Dispatcher {
	return dispatcher.CommandDispatcher{map[string]func([]string, string) string{
		"#subscribe": func(split []string, chatId string) string { return subscriber.AddSubscription(url, split, chatId) },
	}}
}
func BuildParser(init initializer.Initializer, m repositories.MatchDictionaries) parser.Parser {
	return parser.CommandsDecorated(
		BuildCommandDispatcher(init.GetFireBaseSubscriptionsUrl()),
		parser.ContainsStringDecorated(m.StringMatch, parser.ContainsIgnoreCaseMatcherDecorated(m.StringMatchIgnoreCase,
			parser.ExactIgnoreCaseMatcherDecorated(m.ExactMatchIgnoreCase,
			parser.NewExactMatcher(m.ExactMatch)))))
}

func BuildMessage(message *tgbotapi.Message) parser.Message {
	return parser.Message{message.Text, strconv.FormatInt(message.Chat.ID, 10)}
}

func NotifyHandler(init initializer.Initializer, bot *tgbotapi.BotAPI) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		channel := strings.TrimPrefix(r.URL.Path, "/notify/")
		channelsToNotify := subscriber.GetChatIdForChannel(init.GetFireBaseSubscriptionsUrl(), channel)

		type NotificationMessage struct {
			Key string
			Source  string
			Message string
		}

		var mex NotificationMessage

		if r.Method == "POST" {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body",
					http.StatusInternalServerError)
			}

			if body != nil {
				json.Unmarshal(body, &mex)
			}
		}

		if mex.Key == os.Getenv("SecuriyKey") {
			for _, c := range channelsToNotify {
				i, _ := strconv.ParseInt(c, 10, 64)
				msg := tgbotapi.NewMessage(i, "["+mex.Source+"]: "+mex.Message)
				bot.Send(msg)
			}
		}
	}
}


