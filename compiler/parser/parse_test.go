package parser_test

import (
	"barba/compiler"
	"barba/compiler/parser"
	"barba/compiler/tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func NewNodeChain(nodes []*compiler.Node) *compiler.Node {
	head := compiler.NewDummyNode()
	curt := head
	for _, nd := range nodes {
		curt.SetNext(nd)
		curt = curt.GetNext()
	}
	return head.GetNext()
}

func TestNewNodeChain(t *testing.T) {
	nd := compiler.NewLRNode(compiler.ST_IDENT, nil, nil)
	nd.SetNext(compiler.NewLRNode(compiler.ST_RETURN, nil, nil))
	assert.Equal(t, nd, NewNodeChain([]*compiler.Node{
		compiler.NewLRNode(compiler.ST_IDENT, nil, nil),
		compiler.NewLRNode(compiler.ST_RETURN, nil, nil),
	}))
}

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect *compiler.Node
	}{
		{
			"return 100",
			`
func main() int {
	return 100
}
`,
			compiler.NewDefineFunctionNode(
				compiler.NewFunctionDeclarationNode(
					compiler.NewFunctionHeaderNode(
						compiler.NewLeafNode(compiler.ST_IDENT, tokenizer.NewToken(tokenizer.TK_IDENT, "main")),
						compiler.NewFunctionArgumentsNode(nil),
					),
					compiler.NewFunctionReturnDetailsNode(
						compiler.NewFunctionReturnDetailNode(
							compiler.NewLeafNode(compiler.ST_IDENT, tokenizer.NewToken(tokenizer.TK_IDENT, "int")),
						),
					),
				),
				compiler.NewBlockNode(
					compiler.NewLRNode(compiler.ST_RETURN,
						compiler.NewLeafNode(compiler.ST_INTEGER, tokenizer.NewToken(tokenizer.TK_INT, "100")),
						nil),
				),
			),
		},
		{
			"return 200",
			`
func main() int {
	return
	return 200
}
`,
			compiler.NewDefineFunctionNode(
				compiler.NewFunctionDeclarationNode(
					compiler.NewFunctionHeaderNode(
						compiler.NewLeafNode(compiler.ST_IDENT, tokenizer.NewToken(tokenizer.TK_IDENT, "main")),
						compiler.NewFunctionArgumentsNode(nil),
					),
					compiler.NewFunctionReturnDetailsNode(
						compiler.NewFunctionReturnDetailNode(
							compiler.NewLeafNode(compiler.ST_IDENT, tokenizer.NewToken(tokenizer.TK_IDENT, "int")),
						),
					),
				),
				compiler.NewBlockNode(
					NewNodeChain([]*compiler.Node{
						compiler.NewLRNode(compiler.ST_RETURN,
							nil,
							nil),
						compiler.NewLRNode(compiler.ST_RETURN,
							compiler.NewLeafNode(compiler.ST_INTEGER, tokenizer.NewToken(tokenizer.TK_INT, "200")),
							nil),
					}),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := tokenizer.Tokenize(tt.src)
			assert.Nil(t, err)
			nodes, err := parser.Parse(tokens)
			assert.Nil(t, err)
			assert.Equal(t, tt.expect, nodes)
		})
	}
}
