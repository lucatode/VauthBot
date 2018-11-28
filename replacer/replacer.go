package replacer

import (
	"strings"
	"time"
	"strconv"
	"math/rand"
)

type Replacer interface {
	ReplaceIn(string) (string)
}

type RandomRangeNumberReplacer struct {
	Max int
	Placeholder string
	GenerateRandomNumber func(int) string
}

func (rnr RandomRangeNumberReplacer) ReplaceIn(s string) string{
	new := rnr.GenerateRandomNumber(rnr.Max)
	return strings.Replace(s, rnr.Placeholder, new, -1)
}

func GetRandomRangeNumberReplacer(max int, placeholder string, randomFunc func(int) string) RandomRangeNumberReplacer{
	return RandomRangeNumberReplacer{max,placeholder,randomFunc}
}

func GenerateRandomNumeber(max int) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	new := strconv.Itoa(r1.Intn(max))
	return new
}