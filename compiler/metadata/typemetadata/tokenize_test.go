package typemetadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []string
	}{
		{
			"int",
			"int",
			[]string{"int"},
		},
		{
			"int array",
			"[]int",
			[]string{"[", "]", "int"},
		},
		{
			"map string int",
			"map[string]int",
			[]string{"map", "[", "string", "]", "int"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := Tokenize([]rune(tt.input))
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, tokens)
		})
	}
}
