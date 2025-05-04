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
		curt = curt.next
	}
	return head.next
}

func TestNewTokenChain(t *testing.T) {
	tok := NewToken(TK_IDENT, "id")
	tok.next = NewToken(TK_EOF, "")
	assert.Equal(t, tok, NewTokenChain([]*Token{NewToken(TK_IDENT, "id"), NewEofToken()}))
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
			NewTokenChain([]*Token{NewToken(TK_IDENT, "this_is_ident"), NewEofToken()}),
		},
		{
			"escaped string",
			`"hello world\""`,
			NewTokenChain([]*Token{NewToken(TK_STRING, "hello world\""), NewEofToken()}),
		},
		{
			"int",
			"123",
			NewTokenChain([]*Token{NewToken(TK_INT, "123"), NewEofToken()}),
		},
		{
			"float",
			"123.45",
			NewTokenChain([]*Token{NewToken(TK_FLOAT, "123.45"), NewEofToken()}),
		},
		{
			"-int",
			"-123",
			NewTokenChain([]*Token{NewToken(TK_SUB, ""), NewToken(TK_INT, "123"), NewEofToken()}),
		},
		{
			"-float",
			"-123.45",
			NewTokenChain([]*Token{NewToken(TK_SUB, ""), NewToken(TK_FLOAT, "123.45"), NewEofToken()}),
		},
		{
			// 空白は全部読み飛ばされるのでEOFしか残らない
			"whitespace",
			" \n\t\r",
			NewTokenChain([]*Token{NewEofToken()}),
		},
		{
			"escape",
			`"hel\"lo"`,
			NewTokenChain([]*Token{NewToken(TK_STRING, `hel"lo`), NewEofToken()}),
		},
		{
			"=+",
			"=+",
			NewTokenChain([]*Token{NewToken(TK_ASSIGN, ""), NewToken(TK_ADD, ""), NewEofToken()}),
		},
		{
			"==",
			"==",
			NewTokenChain([]*Token{NewToken(TK_EQ, ""), NewEofToken()}),
		},
		{
			"+",
			"+",
			NewTokenChain([]*Token{NewToken(TK_ADD, ""), NewEofToken()}),
		},
		{
			"hello world",
			`print("hello world!")`,
			NewTokenChain([]*Token{
				NewToken(TK_IDENT, "print"),
				NewToken(TK_LRB, ""),
				NewToken(TK_STRING, "hello world!"),
				NewToken(TK_RRB, ""),
				NewEofToken(),
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
