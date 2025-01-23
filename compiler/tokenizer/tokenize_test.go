package tokenizer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewTokenChain(tokens []*Token) *Token {
	head := NewToken(TK_INVALID, "")
	curt := head
	for _, tok := range tokens {
		curt.next = tok
	}
	return head.next
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect *Token
	}{
		{
			"only ident with under-bar",
			"this_is_ident",
			NewToken(TK_IDENT, "this_is_ident"),
		},
		{
			"escaped string",
			`"hello world\""`,
			NewToken(TK_STRING, "hello world\""),
		},
		{
			"int",
			"123",
			NewToken(TK_INT, "123"),
		},
		{
			"float",
			"123.45",
			NewToken(TK_FLOAT, "123.45"),
		},
		{
			"-int",
			"-123",
			NewTokenChain([]*Token{NewToken(TK_SUB, ""), NewToken(TK_INT, "123")}),
		},
		{
			"-float",
			"-123.45",
			NewTokenChain([]*Token{NewToken(TK_SUB, ""), NewToken(TK_FLOAT, "123.45")}),
		},
		{
			"whitespace",
			" \n\t\r",
			NewToken(TK_WHITESPACE, " \n\t\r"),
		},
		{
			"escape",
			`"hel\"lo"`,
			NewToken(TK_STRING, `hel"lo`),
		},
		{
			"=+",
			"=+",
			NewTokenChain([]*Token{NewToken(TK_ASSIGN, ""), NewToken(TK_ADD, "")}),
		},
		{
			"==",
			"==",
			NewToken(TK_EQ, ""),
		},
		{
			"+",
			"+",
			NewToken(TK_ADD, ""),
		},
		{
			"hello world",
			`print("hello world!")`,
			NewTokenChain([]*Token{
				NewToken(TK_IDENT, "print"),
				NewToken(TK_LRB, ""),
				NewToken(TK_STRING, "hello world!"),
				NewToken(TK_RRB, ""),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := Tokenize(tt.src)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, nodes)
		})
	}
}
