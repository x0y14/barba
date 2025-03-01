package compiler

import "barba/compiler/tokenizer"

type Syntax int

const (
	ST_ILLEGAL Syntax = iota
	ST_EOF

	ST_DEFINE_FUNCTION
	ST_FUNCTION_DECLARATION
	ST_FUNCTION_HEADER
	ST_FUNCTION_ARGUMENTS
	ST_FUNCTION_ARGUMENT
	ST_FUNCTION_RETURN_DETAILS
	ST_FUNCTION_RETURN_DETAIL

	ST_IDENT

	ST_PRIMITIVE
	ST_INTEGER

	ST_BLOCK
	ST_RETURN
	ST_IF_ELSE
	ST_IF

	ST_EQ
)

var stKinds = [...]string{
	ST_ILLEGAL: "ILLEGAL",
	ST_EOF:     "EOF",

	ST_DEFINE_FUNCTION:         "DEFINE_FUNCTION",
	ST_FUNCTION_DECLARATION:    "FUNCTION_DECLARATION",
	ST_FUNCTION_HEADER:         "FUNCTION_HEADER",
	ST_FUNCTION_ARGUMENTS:      "FUNCTION_ARGUMENTS",
	ST_FUNCTION_ARGUMENT:       "FUNCTION_ARGUMENT",
	ST_FUNCTION_RETURN_DETAILS: "FUNCTION_RETURN_DETAILS",
	ST_FUNCTION_RETURN_DETAIL:  "FUNCTION_RETURN_DETAIL",

	ST_IDENT:     "IDENT",
	ST_PRIMITIVE: "PRIMITIVE",
	ST_INTEGER:   "INTEGER",

	// Statement Level
	ST_BLOCK:   "BLOCK",
	ST_RETURN:  "RETURN",
	ST_IF_ELSE: "IF_ELSE",

	// EQ LEVEL
	ST_EQ: "EQ",
}

func (st Syntax) String() string {
	return stKinds[st]
}

func NewNode(kind Syntax, lhs, rhs, next *Node, leaf *tokenizer.Token) *Node {
	return &Node{
		kind: kind,
		leaf: leaf,
		lhs:  lhs,
		rhs:  rhs,
		next: next,
	}
}

func NewLeafNode(syntax Syntax, leaf *tokenizer.Token) *Node {
	return NewNode(syntax, nil, nil, nil, leaf)
}

func NewBlockNode(children *Node) *Node {
	return NewNode(ST_BLOCK, children, nil, nil, nil)
}

func NewLRNode(syntax Syntax, lhs, rhs *Node) *Node {
	return NewNode(syntax, lhs, rhs, nil, nil)
}

func NewFunctionArgumentNode(nameLeaf, typeLeaf *Node) *Node {
	return NewNode(ST_FUNCTION_ARGUMENT, nameLeaf, typeLeaf, nil, nil)
}

func NewFunctionArgumentsNode(children *Node) *Node {
	return NewNode(ST_FUNCTION_ARGUMENTS, children, nil, nil, nil)
}

func NewFunctionReturnDetailNode(typLeaf *Node) *Node {
	return NewNode(ST_FUNCTION_RETURN_DETAIL, typLeaf, nil, nil, nil)
}

func NewFunctionReturnDetailsNode(children *Node) *Node {
	return NewNode(ST_FUNCTION_RETURN_DETAILS, children, nil, nil, nil)
}

func NewFunctionHeaderNode(id, args *Node) *Node {
	return NewNode(ST_FUNCTION_HEADER, id, args, nil, nil)
}

func NewFunctionDeclarationNode(header, retDetails *Node) *Node {
	return NewNode(ST_FUNCTION_DECLARATION, header, retDetails, nil, nil)
}

func NewDefineFunctionNode(decl, block *Node) *Node {
	return NewNode(ST_DEFINE_FUNCTION, decl, block, nil, nil)
}

func NewDummyNode() *Node {
	return NewNode(ST_ILLEGAL, nil, nil, nil, nil)
}

func NewEofNode() *Node {
	return NewNode(ST_EOF, nil, nil, nil, nil)
}

type Node struct {
	kind Syntax
	leaf *tokenizer.Token
	lhs  *Node // 1個しか要素がないならLHSを使う
	rhs  *Node
	next *Node
}

func (n *Node) GetKind() Syntax {
	return n.kind
}

func (n *Node) String() string {
	return ""
}

func (n *Node) SetNext(nd *Node) {
	n.next = nd
}

func (n *Node) GetNext() *Node {
	return n.next
}

func (n *Node) SetLhs(nd *Node) {
	n.lhs = nd
}

func (n *Node) GetLhs() *Node {
	return n.lhs
}

func (n *Node) SetRhs(nd *Node) {
	n.rhs = nd
}

func (n *Node) GetRhs() *Node {
	return n.rhs
}

func (n *Node) SetLeaf(tk *tokenizer.Token) {
	n.leaf = tk
}

func (n *Node) GetLeaf() *tokenizer.Token {
	return n.leaf
}
