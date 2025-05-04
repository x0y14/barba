package typemetadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect *Constant
	}{
		{
			"int",
			"int",
			NewInt(),
		},
		{
			"map",
			"map[string]int",
			NewMap(NewString(), NewInt()),
		},
		{
			"map map",
			"map[string]map[string]int",
			NewMap(NewString(), NewMap(NewString(), NewInt())),
		},
		{
			"int array",
			"[]int",
			NewArray(NewInt()),
		},
		{
			"map array",
			"[]map[string]string",
			NewArray(NewMap(NewString(), NewString())),
		},
		{
			"array array int",
			"[][]int",
			NewArray(NewArray(NewInt())),
		},
		{
			"map array v2 err",
			"map[map[string]string][]int",
			NewMap(NewMap(NewString(), NewString()), NewArray(NewInt())),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok, err := Tokenize([]rune(tt.input))
			assert.Nil(t, err)
			typ, err := Parse(tok)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, typ)
		})
	}
}
