package tokenizer

import (
	"fmt"
	"strings"
)

var curtToken *Token
var code []rune
var pos int

func eof() bool {
	return pos+1 >= len(code) // == でぴったり
}

func advance(n int) {
	pos += n
}

func next() rune {
	return code[pos+1]
}

func startWith(r rune) bool {
	return next() == r
}

func startWithWhitespace() bool {
	for _, sr := range []rune{' ', '\n', '\t', '\r'} {
		if startWith(sr) {
			return true
		}
	}
	return false
}

func startWithString() bool {
	return startWith('"')
}

func startWithNumber() bool {
	for _, nr := range []rune("0123456789") {
		if startWith(nr) {
			return true
		}
	}
	return false
}

func startWithSymbol() bool {
	var symbolOps string
	// bracket
	symbolOps += "(){}"
	// symbol
	symbolOps += ";:,."
	// logical
	symbolOps += "><"
	// op
	symbolOps += "=+-*/"

	for _, sor := range []rune(symbolOps) {
		if startWith(sor) {
			return true
		}
	}
	return false
}

func startWithIdent() bool {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := strings.ToUpper(lower)
	under := "_"
	for _, ilr := range []rune(lower + upper + under) {
		if startWith(ilr) {
			return true
		}
	}
	return false
}

func consume() rune {
	r := next()
	advance(1)
	return r
}

func expect(r rune) error {
	if !startWith(r) {
		return fmt.Errorf("expect %v, but got %v", string(r), string(next()))
	}
	consume()
	return nil
}

func consumeWhitespace() error {
	var ws []rune
loop:
	for !eof() {
		switch {
		case startWithWhitespace():
			ws = append(ws, consume())
		default:
			break loop
		}
	}
	curtToken.next = NewToken(TK_WHITESPACE, string(ws))
	return nil
}

func consumeString() error {
	var str []rune
	escapeMode := false
	_ = expect('"')

loop:
	for !eof() {
		if escapeMode { // エスケープモード
			str = append(str, consume()) // newlineとかでエラー出るかも
			escapeMode = false
		} else { // 通常モード
			switch {
			case startWith('\\'): // turn on escape mode
				_ = consume()
				escapeMode = true
			case startWith('"'): // エスケープされておらず，"を発見
				_ = expect('"')
				break loop
			default:
				str = append(str, consume())
			}
		}
	}
	curtToken.next = NewToken(TK_STRING, string(str))
	return nil
}

func consumeNumber() error {
	var num []rune
	dotCount := 0

loop:
	for !eof() {
		switch {
		case startWithNumber():
			num = append(num, consume())
		case startWith('.'):
			num = append(num, consume())
			dotCount++
		default:
			break loop
		}
	}

	if dotCount == 0 {
		curtToken.next = NewToken(TK_INT, string(num))
		return nil
	} else if dotCount == 1 {
		if num[0] == '.' || num[len(num)-1] == '.' { // .012とか123.とかはエラー
			return fmt.Errorf("number has invalid-position dot")
		}
		curtToken.next = NewToken(TK_FLOAT, string(num))
		return nil
	} else {
		return fmt.Errorf("number has too many dot")
	}
}

func consumeSymbol() error {
	singles := map[rune]TokenKind{
		'(': TK_LRB,
		')': TK_RRB,
		'{': TK_LCB,
		'}': TK_RCB,

		';': TK_SEMI,
		':': TK_COLON,
		',': TK_COMMA,
		'.': TK_DOT,

		'<': TK_LT,
		'>': TK_GT,

		'=': TK_ASSIGN,
		'+': TK_ADD,
		'-': TK_SUB,
		'*': TK_MUL,
		'/': TK_DIV,
	}
	compounds := map[[2]TokenKind]TokenKind{
		{TK_ASSIGN, TK_ASSIGN}: TK_EQ, // ==
		{TK_NOT, TK_ASSIGN}:    TK_NE, // !=
		{TK_LT, TK_ASSIGN}:     TK_LE, // <=
		{TK_GT, TK_ASSIGN}:     TK_GE, // >=
	}

	var first TokenKind = TK_INVALID
	for !eof() {
		if stk, ok := singles[next()]; ok { // 複合でなかった場合を考慮して, consumeじゃなくてnextを使って確認だけする
			if first == TK_INVALID {
				first = stk
				consume()
			} else { // 二つ目の記号を発見したら最初のものと組み合わせて複合記号ができるかチェック
				ctk, ok := compounds[[2]TokenKind{first, stk}]
				if ok { // できた
					first = ctk
					consume() // 確認できたので消費
				} else { // できなかったのでなにもせず
					break
				}
			}
		} else {
			break
		}
	}
	curtToken.next = NewToken(first, "")
	return nil
}

func consumeIdent() error {
	var id []rune

loop:
	for !eof() {
		switch {
		case startWithIdent():
			id = append(id, consume())
		case startWithNumber():
			id = append(id, consume())
		default:
			break loop
		}
	}

	curtToken.next = NewToken(TK_IDENT, string(id))
	return nil
}

func Tokenize(sourceCode string) (*Token, error) {
	headToken := NewToken(TK_INVALID, "")
	curtToken = headToken

	code = []rune(sourceCode)
	pos = -1 // ダミーを参照

	for !eof() {
		switch {
		case startWithWhitespace(): // whitespace
			if err := consumeWhitespace(); err != nil {
				return nil, err
			}
		case startWithString(): // string
			if err := consumeString(); err != nil {
				return nil, err
			}
		case startWithNumber(): // number
			if err := consumeNumber(); err != nil {
				return nil, err
			}
		case startWithSymbol(): // symbol or logical or op
			if err := consumeSymbol(); err != nil {
				return nil, err
			}
		case startWithIdent(): // ident
			if err := consumeIdent(); err != nil {
				return nil, err
			}
		}
	}

	return headToken.next, nil
}
