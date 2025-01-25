package parser

import (
	"barba/compiler"
	"barba/compiler/tokenizer"
	"fmt"
)

var curt *tokenizer.Token
var nodes *compiler.Node

func next() *tokenizer.Token {
	return curt.GetNext()
}

func eof() bool {
	return next().GetKind() == tokenizer.TK_EOF
}

func startWith(kind tokenizer.TokenKind) bool {
	return next().GetKind() == kind
}

func startWithIdent(id string) bool {
	return startWith(tokenizer.TK_IDENT) && next().GetText() == id
}

func advance() *tokenizer.Token {
	t := next()
	curt = curt.GetNext()
	return t
}

func consume(kind tokenizer.TokenKind) *tokenizer.Token {
	if startWith(kind) {
		return advance()
	}
	return nil
}

func expect(kind tokenizer.TokenKind) error {
	v := consume(kind)
	if v == nil {
		return fmt.Errorf("expect %v, but got %v", kind.String(), next())
	}
	return nil
}

func consumeLiteralLv() error {
	switch {
	case startWith(tokenizer.TK_INT):
		nodes.SetNext(compiler.NewLeafNode(compiler.ST_INTEGER, consume(tokenizer.TK_INT).ShallowClone()))
		nodes = nodes.GetNext()
	}
	return nil
}

func consumeAccessLv() error {
	return consumeLiteralLv()
}

func consumePrimaryLv() error {
	return consumeAccessLv()
}

func consumeUnaryLv() error {
	return consumePrimaryLv()
}

func consumeMulLv() error {
	return consumeUnaryLv()
}

func consumeAddLv() error {
	return consumeMulLv()
}

func consumeRelationalLv() error {
	return consumeAddLv()
}

func consumeEqualityLv() error {
	return consumeRelationalLv()
}

func consumeAndorLv() error {
	return consumeEqualityLv()
}

func consumeAssignLv() error {
	return consumeAndorLv()
}

func consumeExprLv() error {
	return consumeAssignLv()
}

func consumeReturn() error {
	// return
	_ = advance()

	// 戻り値は複数記述される可能性がありreturnとしてまとめたいので
	backup := nodes
	dummyForReturn := compiler.NewDummyNode()
	nodes = dummyForReturn

	for !eof() {
		// returnだけで戻り値がない可能性あり, exprがないならstmtまたは}なのでどちらにしろエラー起こる.
		// エラーが出なかったらそれが戻り値になるというだけ
		if err := consumeExprLv(); err != nil {
			break
		}
		if consume(tokenizer.TK_COMMA) == nil {
			break
		}
	}
	// 復元
	nodes = backup
	//
	nodes.SetNext(
		compiler.NewLRNode(compiler.ST_RETURN, dummyForReturn.GetNext(), nil),
	)
	nodes = nodes.GetNext()

	return nil
}

func consumeStmtLv() error {
	switch {
	case startWithIdent("return"):
		return consumeReturn()
	default:
		return consumeExprLv()
	}
}

func consumeTopLv() error {
	switch {
	case startWithIdent("func"):
		return consumeDefineFunction()
	default:
		return fmt.Errorf("unsupported toplevel: %v", next())
	}
}

func consumeBlock() error {
	if err := expect(tokenizer.TK_LCB); err != nil { // {
		return err
	}

	// stmtを直接つなげるのではなくblockでまとめたいので
	backup := nodes
	dummyForBlock := compiler.NewDummyNode()
	nodes = dummyForBlock
	for consume(tokenizer.TK_RCB) == nil { // }を見つけるまで
		if err := consumeStmtLv(); err != nil {
			return err
		}
	}

	// 復元
	nodes = backup
	//
	nodes.SetNext(
		compiler.NewBlockNode(dummyForBlock.GetNext()),
	)
	nodes = nodes.GetNext()

	return nil
}

func consumeFuncArg() error {
	arg := consume(tokenizer.TK_IDENT)
	if arg == nil {
		return fmt.Errorf("argument ident expect ident, but got %v", next())
	}

	typ := consume(tokenizer.TK_IDENT)
	if typ == nil {
		return fmt.Errorf("argument type expect ident, but got %v", next())
	}

	// 現在に直接接続
	nodes.SetNext(
		compiler.NewFunctionArgumentNode(
			compiler.NewLeafNode(compiler.ST_IDENT, arg.ShallowClone()),
			compiler.NewLeafNode(compiler.ST_IDENT, typ.ShallowClone())),
	)
	// 前進
	nodes = nodes.GetNext()

	return nil
}

func consumeFuncArgs() error {
	// func foo(arg...) { stmt... }
	//         ^
	if err := expect(tokenizer.TK_LRB); err != nil {
		return err
	}
	// 直接argをつけるのではなくargsとしてまとめてつけたいので
	backup := nodes
	dummyForArg := compiler.NewDummyNode()
	nodes = dummyForArg

	for consume(tokenizer.TK_RRB) == nil { // )を見つけたら終わる
		if err := consumeFuncArg(); err != nil {
			return err
		}
		if consume(tokenizer.TK_COMMA) == nil { // ,がなかったら終わる
			// )で終わっていることを確認
			if err := expect(tokenizer.TK_RRB); err != nil {
				return err
			}
			break
		}
	}
	// )は消費済み

	// 復元
	nodes = backup

	// 現在のnodeにつける, 回収は呼び出し元
	nodes.SetNext(compiler.NewFunctionArgumentsNode(dummyForArg.GetNext()))
	// 前進
	nodes = nodes.GetNext()
	return nil
}

func consumeFuncReturnDetail() error {
	typ := consume(tokenizer.TK_IDENT)
	if typ == nil {
		return fmt.Errorf("return detail expect ident, but got %v", next())
	}
	nodes.SetNext(
		compiler.NewFunctionReturnDetailNode(compiler.NewLeafNode(compiler.ST_IDENT, typ.ShallowClone())),
	)
	nodes = nodes.GetNext()
	return nil
}

func consumeFuncReturnDetails() error {
	// detailではなくdetailsでまとめてつけたいので
	backup := nodes
	dummyForRetDetails := compiler.NewDummyNode()
	nodes = dummyForRetDetails
	if startWith(tokenizer.TK_LCB) {
		// 復元
		nodes = backup
		nodes.SetNext(
			compiler.NewFunctionReturnDetailsNode(nil),
		)
		nodes = nodes.GetNext()
	}

	// (がない場合
	if consume(tokenizer.TK_LRB) == nil {
		if err := consumeFuncReturnDetail(); err != nil {
			return err
		}
		// 復元
		nodes = backup
		//
		nodes.SetNext(
			compiler.NewFunctionReturnDetailsNode(dummyForRetDetails.GetNext()),
		)
		nodes = nodes.GetNext()
		return nil
	}
	// (があった場合, (はすでに消費されている
	for consume(tokenizer.TK_RRB) == nil { // )が見つかるまで
		if err := consumeFuncReturnDetail(); err != nil {
			return err
		}
		if consume(tokenizer.TK_COMMA) == nil { // ,がなかったら
			if err := expect(tokenizer.TK_RRB); err != nil { // )で終わっていることを確認
				return err
			}
			break
		}
	}
	// )はすでに消費されている

	// 復元
	nodes = backup
	//
	nodes.SetNext(
		compiler.NewFunctionReturnDetailsNode(dummyForRetDetails.GetNext()),
	)
	nodes = nodes.GetNext()

	return nil
}

func consumeFuncHeader() error {
	// func foo(arg...) { stmt... }
	//      ^
	name := consume(tokenizer.TK_IDENT)
	funcId := compiler.NewLeafNode(compiler.ST_IDENT, name.ShallowClone())
	// func foo(arg...) { stmt... }
	//         ^

	// 現在のノードに直接つけたらダメなので
	backup := nodes
	// ダミーを用意
	dummyForArgs := compiler.NewDummyNode()
	nodes = dummyForArgs
	if err := consumeFuncArgs(); err != nil {
		return err
	}
	// 復元
	nodes = backup

	// 現在にHeaderをつける
	// 回収は親がする
	nodes.SetNext(compiler.NewFunctionHeaderNode(
		funcId,
		dummyForArgs.GetNext()))
	// 前進
	nodes = nodes.GetNext()

	return nil
}

func consumeFuncDecl() error {
	// 現在のノードに直接つけたらダメなので
	backup := nodes
	// ダミーを用意
	dummyForHeader := compiler.NewDummyNode()
	nodes = dummyForHeader
	if err := consumeFuncHeader(); err != nil {
		return err
	}
	// ダミーを用意
	dummyForRetDetails := compiler.NewDummyNode()
	nodes = dummyForRetDetails
	if err := consumeFuncReturnDetails(); err != nil {
		return err
	}
	// 復元
	nodes = backup

	// 現在にDeclをつける
	// 回収は親がする
	nodes.SetNext(compiler.NewFunctionDeclarationNode(dummyForHeader.GetNext(), dummyForRetDetails.GetNext()))
	// 前進
	nodes = nodes.GetNext()

	return nil
}

func consumeDefineFunction() error {
	// func foo(arg...) { stmt... }
	// ^
	_ = advance()

	// func foo(arg...) { stmt... }
	//      ^
	// 現在に直接つけるのはダメなので
	backup := nodes
	// ダミーを用意
	dummyForDecl := compiler.NewDummyNode()
	nodes = dummyForDecl
	if err := consumeFuncDecl(); err != nil {
		return err
	}

	// func foo(arg...) { stmt... }
	//                  ^
	// 現在に直接つけるのはダメなので
	// ダミーを用意
	dummyForBlock := compiler.NewDummyNode()
	nodes = dummyForBlock
	if err := consumeBlock(); err != nil {
		return err
	}
	// 復元
	nodes = backup

	// 現在のノードの次にFUNCをつけて
	nodes.SetNext(compiler.NewDefineFunctionNode(dummyForDecl.GetNext(), dummyForBlock.GetNext()))
	// 前進
	nodes = nodes.GetNext()

	return nil
}

func Parse(token *tokenizer.Token) (*compiler.Node, error) {
	head := tokenizer.NewToken(tokenizer.TK_INVALID, "")
	head.SetNext(token)
	curt = head

	nodes = compiler.NewDummyNode()

	for !eof() {
		if err := consumeTopLv(); err != nil {
			return nil, err
		}
	}

	return nodes, nil
}
