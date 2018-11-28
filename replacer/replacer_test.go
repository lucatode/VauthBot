package replacer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExactMatchDecoratedWithCommand(t *testing.T) {

	replacer := RandomRangeNumberReplacer{"%n"}

	stringOutput := replacer.ReplaceIn("This is a text and this (%n) should be replaced with a random number")
	assert.Equal(t, "This is a text and this 23 should be replaced with a random number", stringOutput, "")

}

