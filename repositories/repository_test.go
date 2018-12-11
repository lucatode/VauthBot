package repositories

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"vauthbot/logger"
)

func TestExactMatch(t *testing.T) {
	client := http.Client{}
	const logsJson = "https://xxxxxxxxx.firebaseio.com/logs.json"
	const repoJson = "https://xxxxxxxxx.firebaseio.com/responses.json"
	const request = "test"
	const response = "TEST"

	logger := logger.FirebaseLogger{logsJson}
	repo := FireBaseRepository{client.Get, logger, func(repo FireBaseRepository, url string) []MatchCase {
		return []MatchCase{{true, request, response, false}}
	}}

	var dict = repo.FilterMatchCases(repoJson, ExactMatchCasesToMap)

	l := len(dict)
	assert.Equal(t, 1, l)
	m := dict[request]
	assert.Equal(t, response, m)
}

func TestExactMatchIgnoreCase(t *testing.T) {
	client := http.Client{}
	const repoJson = "https://xxxxxxxxx.firebaseio.com/responses.json"
	const logsJson = "https://xxxxxxxxx.firebaseio.com/logs.json"
	const request = "teST"
	const response = "TEST"

	logger := logger.FirebaseLogger{logsJson}
	repo := FireBaseRepository{client.Get, logger, func(repo FireBaseRepository, url string) []MatchCase {
		return []MatchCase{{true, request, response, true}}
	}}

	var dict = repo.FilterMatchCases(repoJson, ExactMatchCasesIgnoreCaseToMap)

	l := len(dict)
	assert.Equal(t, 1, l)
	m := dict["test"]
	assert.Equal(t, response, m)
}

func TestDictionariesBuilder(t *testing.T) {
	const repoJson = "https://xxxxxxxxx.firebaseio.com/responses.json"

	const logsJson = "https://xxxxxxxxx.firebaseio.com/logs.json"

	dicts := BuildDictionaries(repoJson, nil)

	assert.Equal(t, 1, len(dicts.ExactMatch))
	assert.Equal(t, 1, len(dicts.ExactMatchIgnoreCase))
	assert.Equal(t, 1, len(dicts.StringMatch))
	assert.Equal(t, 1, len(dicts.StringMatchIgnoreCase))
}


