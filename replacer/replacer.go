package replacer

import (
	"strings"
	"time"
	"math/rand"
	"strconv"
)

type Replacer interface {
	ReplaceIn(string) (string)
}

type RandomRangeNumberReplacer struct {
	Placeholder string
}

func (rnr RandomRangeNumberReplacer) ReplaceIn(s string) string{
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	new := strconv.Itoa(r1.Intn(100000))
	return strings.Replace(s, rnr.Placeholder, new, -1)
}