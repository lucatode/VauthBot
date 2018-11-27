package repositories

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"net/http"
	"condorbot/logger"
)



func TestCheckStringMatch(t *testing.T) {
	client := http.Client{}
	logger := logger.FirebaseLogger{"https://xxxxxxxxx.firebaseio.com/logs.json"}
	repo := FireBaseRepository{client.Get, logger}

	var dict = repo.GetExactMatchMap("https://xxxxxxxxx.firebaseio.com/responses.json")

	l := len(dict)
	assert.Equal(t, 1, l, "Getting specific match bool")
}
