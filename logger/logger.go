package logger

import (
	"net/http"
	"bytes"
	"time"
	"encoding/json"
)

type Log struct{
	Source string
	Message string
	Level string
	Time string
}

type Logger interface{
	Log(source string, message string)
	Warn(source string, message string)
	Err(source string, message string)
}

type FirebaseLogger struct{
	Url string
}

func HttpPost(url string, Log Log){
	jsonStr, _ := json.Marshal(Log)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "log")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (logger FirebaseLogger) Log(source string, message string) {
	l := Log { Source: source, Message: message, Level:"Info", Time:time.Now().String() }
	HttpPost(logger.Url, l)
}

func (logger FirebaseLogger) Warn(source string, message string) {
	l := Log { Source: source, Message: message, Level:"Warning", Time:time.Now().String() }
	HttpPost(logger.Url, l)
}

func (logger FirebaseLogger) Err(source string, message string) {
	l := Log { Source: source, Message: message, Level:"Error", Time:time.Now().String() }
	HttpPost(logger.Url, l)
}

