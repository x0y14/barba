package typemetadata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeMetadata(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect *Constant
	}{
		{
			"string",
			"string",
			NewString(),
		},
		{
			"int",
			"int",
			NewInt(),
		},
		{
			"float",
			"float",
			NewFloat(),
		},
		{
			"bool",
			"bool",
			NewBool(),
		},
		{
			"array of strings",
			"[]string",
			NewArray(NewString()),
		},
		{
			"map of string to int",
			"map[string]int",
			NewMap(NewString(), NewInt()),
		},
		{
			"nested array",
			"[][]int",
			NewArray(NewArray(NewInt())),
		},
		{
			"map of string to array of bool",
			"map[string][]bool",
			NewMap(NewString(), NewArray(NewBool())),
		},
		{
			"map of int to map of string to float",
			"map[int]map[string]float",
			NewMap(NewInt(), NewMap(NewString(), NewFloat())),
		},
		{
			"complex nested type",
			"map[string][]map[int]string",
			NewMap(
				NewString(),
				NewArray(
					NewMap(NewInt(), NewString()),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ, err := NewConstantFromStr(tt.input)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, typ)
		})
	}
}
