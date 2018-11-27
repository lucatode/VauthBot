package subscriber

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Subscription struct {
	Channel string
	ChatId  string
}

func httpPost(url string, subscription Subscription) {
	jsonStr, _ := json.Marshal(subscription)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "subscription")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func AddSubscription(url string, message []string, chatId string) string {

	if len(message) == 3 {
		channel := message[2]
		s := Subscription{channel, chatId}
		httpPost(url+"/"+channel+".json", s)
		return "Subscribed this chat to channel " + channel
	}
	return "invalid message for subscribe channel"
}

func GetChatIdForChannel(url string, channel string) []string {
	chatIds := []string{}
	client := http.Client{}
	resp, err := client.Get(url + "/" + channel + ".json")

	if err != nil {
		//repo.Logger.Err("FireBaseRepository", "First err: "+err.Error())
	}
	defer resp.Body.Close()

	var bytesArray []byte
	if resp.StatusCode == http.StatusOK {
		bytesArray, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			//repo.Logger.Err("FireBaseRepository", "First err: "+err.Error())
		}
	}

	if bytesArray != nil {
		var subscriptions map[string]Subscription
		json.Unmarshal(bytesArray, &subscriptions)

		for _, sub := range subscriptions {
			chatIds = append(chatIds, sub.ChatId)
		}
	}
	return chatIds
}
