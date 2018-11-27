package initializer

import (
	"os"
	"strconv"
)

func NewInitializer(ps InitializerStorage) *ParameterInitializer {
	return &ParameterInitializer{
		storage: ps,
	}
}

type Initializer interface {
	GetApiToken() string
	GetServerUrl() string
	GetTimerSeconds() int
	GetFireBaseResponsesUrl() string
	GetFireBaseLogsUrl() string
	GetFireBaseSubscriptionsUrl() string
}

type ParameterInitializer struct {
	storage InitializerStorage
}

func (init ParameterInitializer) GetFireBaseLogsUrl() string {
	return init.storage.GetData()["FireBaseLogsUrl"]
}

func (init ParameterInitializer) GetApiToken() string {
	return init.storage.GetData()["ApiToken"]
}

func (init ParameterInitializer) GetServerUrl() string {
	return init.storage.GetData()["ServerUrl"]
}

func (init ParameterInitializer) GetTimerSeconds() int {
	ret, _ := strconv.Atoi(init.storage.GetData()["TimerSeconds"])
	return ret
}

func (init ParameterInitializer) GetFireBaseResponsesUrl() string {
	return init.storage.GetData()["FireBaseResponsesUrl"]
}

func (init ParameterInitializer) GetFireBaseSubscriptionsUrl() string {
	return init.storage.GetData()["FireBaseSubscriptionsUrl"]
}

//--- Initializer Storage
type InitializerStorage interface {
	GetData() map[string]string
}

//--- Env Storage
func NewEnvReader() *EnvStorage {
	return &EnvStorage{""}
}

type EnvStorage struct {
	storageType string
}

func (storage EnvStorage) GetData() map[string]string {
	return map[string]string{
		"ApiToken":                 os.Getenv("ApiToken"),
		"ServerUrl":                os.Getenv("ServerUrl"),
		"TimerSeconds":             os.Getenv("TimerSeconds"),
		"FireBaseResponsesUrl":     os.Getenv("FireBaseResponsesUrl"),
		"FireBaseLogsUrl":          os.Getenv("FireBaseLogsUrl"),
		"FireBaseSubscriptionsUrl": os.Getenv("FireBaseSubscriptionsUrl"),
	}
}
