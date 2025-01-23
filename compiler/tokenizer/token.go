package tokenizer

import (
	"fmt"
	"strconv"
)

type TokenKind int

const (
	TK_INVALID TokenKind = iota
	TK_EOF

	TK_NULL   // null
	TK_INT    // 12
	TK_FLOAT  // 12.3
	TK_STRING // "string"

	TK_IDENT      // name
	TK_KEYWORD    // var, fn, ...
	TK_COMMENT    // // this is comment, start with double slash
	TK_WHITESPACE // " ", "\n", "\t"

	TK_LRB // (
	TK_RRB // )
	TK_LCB // {
	TK_RCB // }

	TK_SEMI  // ;
	TK_COLON // :
	TK_COMMA // ,
	TK_DOT   // .

	TK_EQ  // ==
	TK_NE  // !=
	TK_LT  // <
	TK_LE  // <=
	TK_GT  // >
	TK_GE  // >=
	TK_NOT // !

	TK_ASSIGN // =
	TK_ADD    // +
	TK_SUB    // -
	TK_MUL    // *
	TK_DIV    // /
)

var tokKinds = [...]string{
	TK_INVALID: "INVALID",
	TK_EOF:     "EOF",

	TK_NULL:   "NULL",
	TK_INT:    "INT",
	TK_FLOAT:  "FLOAT",
	TK_STRING: "STRING",

	TK_IDENT:      "IDENT",
	TK_KEYWORD:    "KEYWORD",
	TK_WHITESPACE: "WHITESPACE",
	TK_COMMENT:    "COMMENT",

	TK_LRB: "(",
	TK_RRB: ")",
	TK_LCB: "{",
	TK_RCB: "}",

	TK_SEMI:  ";",
	TK_COLON: ":",
	TK_COMMA: ",",
	TK_DOT:   ".",

	TK_EQ:  "==",
	TK_NE:  "!=",
	TK_LT:  "<",
	TK_LE:  "<=",
	TK_GT:  ">",
	TK_GE:  ">=",
	TK_NOT: "!",

	TK_ASSIGN: "=",
	TK_ADD:    "+",
	TK_SUB:    "-",
	TK_MUL:    "*",
	TK_DIV:    "/",
}

func (tk TokenKind) String() string {
	return tokKinds[tk]
}

func NewToken(kind TokenKind, text string) *Token {
	return &Token{
		kind: kind,
		text: text,
		next: nil,
	}
}

type Token struct {
	kind TokenKind
	text string
	next *Token
}

func (t *Token) GetInt() (int, error) {
	if t.kind != TK_INT {
		return 0, fmt.Errorf("this token is not int: %v", t.kind.String())
	}
	return strconv.Atoi(t.text)
}

func (t *Token) GetFloat() (float64, error) {
	if t.kind != TK_FLOAT {
		return 0, fmt.Errorf("this token is not float: %v", t.kind.String())
	}
	return strconv.ParseFloat(t.text, 64)
}

func (t *Token) GetString() (string, error) {
	if t.kind != TK_STRING {
		return "", fmt.Errorf("this token is not string: %v", t.kind.String())
	}
	return t.text, nil
}

func (t *Token) GetIdent() (string, error) {
	if t.kind != TK_IDENT {
		return "", fmt.Errorf("this token is not ident: %v", t.kind.String())
	}
	return t.text, nil
}

func (t *Token) GetKind() TokenKind {
	return t.kind
}

func (t *Token) GetText() string {
	return t.text
}

func (t *Token) SetNext(token *Token) {
	t.next = token
}

func (t *Token) GetNext() *Token {
	return t.next
}

func NewEofToken() *Token {
	return NewToken(TK_EOF, "")
}
