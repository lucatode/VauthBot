package repositories

import (
	"vauthbot/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type MatchCase struct {
	MatchExact bool
	Request    string
	Response   string
	IgnoreCase bool
}

type MatchDictionaries struct {
	ExactMatch            map[string]string
	ExactMatchIgnoreCase  map[string]string
	StringMatch           map[string]string
	StringMatchIgnoreCase map[string]string
}

type Repository interface {
	FilterMatchCases(url string, f func([]MatchCase) map[string]string) map[string]string
}

type FireBaseRepository struct {
	Delegate      func(string) (*http.Response, error)
	Logger        logger.Logger
	GetMatchCases func(repo FireBaseRepository, url string) []MatchCase
}

func (repo FireBaseRepository) FilterMatchCases(url string, f func([]MatchCase) map[string]string) map[string]string {
	cases := repo.GetMatchCases(repo, url)
	return f(cases)
}

func GetMatchCases(repo FireBaseRepository, url string) []MatchCase {
	resp, err := repo.Delegate(url)
	if err != nil {
		repo.Logger.Err("FireBaseRepository", "First err: "+err.Error())
	}
	defer resp.Body.Close()

	var byteArray []byte
	if resp.StatusCode == http.StatusOK {
		byteArray, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			repo.Logger.Err("FireBaseRepository", "First err: "+err.Error())
		}
	}

	if byteArray != nil {
		var cases []MatchCase
		s := string(byteArray)
		repo.Logger.Log("REPO", "received "+s)
		json.Unmarshal(byteArray, &cases)
		return cases
	}

	return nil
}

func ExactMatchCasesToMap(matchCases []MatchCase) map[string]string {
	dict := make(map[string]string)
	for _, matchCase := range matchCases {
		if matchCase.MatchExact && !matchCase.IgnoreCase{
			dict[matchCase.Request] = matchCase.Response
		}
	}
	return dict
}

func ExactMatchCasesIgnoreCaseToMap(matchCases []MatchCase) map[string]string {
	dict := make(map[string]string)
	for _, matchCase := range matchCases {
		if matchCase.MatchExact && matchCase.IgnoreCase {
			low := strings.ToLower(matchCase.Request)
			dict[low] = matchCase.Response
		}
	}
	return dict
}

func WordMatchCasesToMap(matchCases []MatchCase) map[string]string {
	dict := make(map[string]string)
	for _, matchCase := range matchCases {
		if !matchCase.MatchExact && !matchCase.IgnoreCase{
			dict[matchCase.Request] = matchCase.Response
		}
	}
	return dict
}

func WordMatchCasesIgnoreCaseToMap(matchCases []MatchCase) map[string]string {
	dict := make(map[string]string)
	for _, matchCase := range matchCases {
		if !matchCase.MatchExact && matchCase.IgnoreCase{
			dict[matchCase.Request] = matchCase.Response
		}
	}
	return dict
}

func BuildDictionaries(url string, logger logger.Logger) MatchDictionaries {
	client := http.Client{}
	repo := FireBaseRepository{client.Get, logger, func(repo FireBaseRepository, url string) []MatchCase {
		return []MatchCase{
			{true, "abc", "def", false},
			{true, "ABC", "DEF", true},
			{false, "A B C", "OoOoOoO", false},
			{false, "j k l", "zzzzZZZZzzz", true},
		}
	}}
	return MatchDictionaries{
		repo.FilterMatchCases(url, ExactMatchCasesToMap),
		repo.FilterMatchCases(url, ExactMatchCasesIgnoreCaseToMap),
		repo.FilterMatchCases(url, WordMatchCasesToMap),
		repo.FilterMatchCases(url, WordMatchCasesIgnoreCaseToMap),
	}
}