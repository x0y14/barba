package compiler

import (
	"barba/runtime"
	"fmt"
	"log"
)

const (
	MAX_RETURN_VALUE = 2
)

var curt *Node
var st *SymbolTable

func nextNode() error {
	if curt.next == nil {
		return fmt.Errorf("end of node")
	}
	curt = curt.next
	return nil
}

func genToplevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	case ST_DEFINE_FUNCTION:
		return genDefineFunction(nd)
	default:
		return nil, fmt.Errorf("unsupported toplevel syntax: %v", nd.kind.String())
	}
}

func genStatementLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	case ST_RETURN:
		return genReturn(nd)
	case ST_IF_ELSE:
		return genIfElse(nd)
	case ST_BLOCK:
		return genBlock(nd)
	default:
		return genExprLevel(nd)
	}
}

func genIfElse(nd *Node) (runtime.Program, error) {
	// ラベルの作成
	curtIfId := RandomString(10)
	lIf, err := st.RegisterLabel(fmt.Sprintf("%s_if_%s_if", st.curtFn, curtIfId))
	lElse, err := st.RegisterLabel(fmt.Sprintf("%s_if_%s_else", st.curtFn, curtIfId))
	lEnd, err := st.RegisterLabel(fmt.Sprintf("%s_if_%s_end", st.curtFn, curtIfId))

	// lhs: if
	//	lhs: condition
	//	rhs: block
	// rhs: else(if)
	//	lhs: nil
	//	rhs: block
	prog := runtime.Program{}
	ifCond, ifBlock, err := genIf(nd.lhs)
	if err != nil {
		return nil, err
	}
	elseBlock, err := genBlock(nd.rhs)
	if err != nil {
		return nil, err
	}

	// # if cond {}をする #
	// ## 条件式 ##
	prog = append(prog, ifCond...)
	// zfで判断するので条件式の結果は捨てる
	prog = append(prog, runtime.Program{
		runtime.Pop, runtime.Temporal1,
	}...)
	// ## 条件分岐 ##
	prog = append(prog, runtime.Program{
		runtime.Je, runtime.Label(lIf),
		runtime.Jmp, runtime.Label(lElse),
	}...)

	// IF BLOCK
	prog = append(prog, runtime.Program{
		runtime.DefLabel(lIf),
	}...)
	// 中身
	prog = append(prog, ifBlock...)
	// IFの終了へ
	prog = append(prog, runtime.Program{
		runtime.Jmp, runtime.Label(lEnd),
	}...)

	// ELSE BLOCK
	prog = append(prog, runtime.Program{
		runtime.DefLabel(lElse),
	}...)
	// 中身
	prog = append(prog, elseBlock...)
	// IFの終了へ
	prog = append(prog, runtime.Program{
		runtime.Jmp, runtime.Label(lEnd),
	}...)

	// IFの終了
	prog = append(prog, runtime.Program{
		runtime.DefLabel(lEnd),
	}...)
	return prog, nil
}

func genIf(nd *Node) (runtime.Program, runtime.Program, error) {
	// lhs: condition
	// rhs: block
	condition, err := genExprLevel(nd.lhs)
	if err != nil {
		return nil, nil, err
	}
	block, err := genBlock(nd.rhs)
	if err != nil {
		return nil, nil, err
	}
	return condition, block, nil
}

func genExprLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genAssignLevel(nd)
	}
}

func genAssignLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genAndorLevel(nd)
	}
}

func genAndorLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genEqualityLevel(nd)
	}
}

func genEqualityLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	case ST_EQ:
		prog := runtime.Program{}
		// 左辺の評価
		lhs, err := genRelationalLevel(nd.lhs)
		if err != nil {
			return nil, err
		}
		prog = append(prog, lhs...)
		// 右辺の評価
		rhs, err := genRelationalLevel(nd.rhs)
		if err != nil {
			return nil, err
		}
		prog = append(prog, rhs...)
		// 比較
		prog = append(prog, runtime.Program{
			runtime.Pop, runtime.R2, // 右辺の取り出し
			runtime.Pop, runtime.R1, // 左辺の取り出し
			runtime.Eq, runtime.R1, runtime.R2, // r1 == r2
			runtime.Push, runtime.ZeroFlag, // 結果を投げる
		}...)
		return prog, nil
	default:
		return genRelationalLevel(nd)
	}
}

func genRelationalLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genAddLevel(nd)
	}
}

func genAddLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genMulLevel(nd)
	}
}

func genMulLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genUnaryLevel(nd)
	}
}

func genUnaryLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genPrimaryLevel(nd)
	}
}

func genPrimaryLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genAccessLevel(nd)
	}
}

func genAccessLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	default:
		return genLiteralLevel(nd)
	}
}

func genLiteralLevel(nd *Node) (runtime.Program, error) {
	switch nd.kind {
	case ST_PRIMITIVE:
		return genPrimitive(nd)
	default:
		return nil, fmt.Errorf("unsupported literal syntax: %v", nd.kind.String())
	}
}

func genPrimitive(nd *Node) (runtime.Program, error) {
	switch primValue := nd.lhs; primValue.kind {
	case ST_INTEGER:
		i, err := primValue.leaf.GetInt()
		if err != nil {
			return nil, err
		}
		return runtime.Program{runtime.Push, runtime.Integer(i)}, nil
	default:
		return nil, fmt.Errorf("genPrimitive: unsupported value: %s", primValue.String())
	}
}

func genReturn(nd *Node) (runtime.Program, error) {
	prog := runtime.Program{}
	count := 0
	c := nd.lhs
retLoop:
	for {
		log.Println("ret")
		if count >= MAX_RETURN_VALUE {
			break retLoop
		}
		switch {
		case c == nil:
			break retLoop
		case c.kind == ST_PRIMITIVE:
			primProg, err := genPrimitive(c)
			if err != nil {
				return nil, err
			}
			// # 戻り値の返却 #
			// スタックに値が入っている
			prog = append(prog, primProg...)
			switch count {
			case 0: // 1つめの戻り値, ACM1に.
				prog = append(prog, runtime.Program{
					// スタックから取り出してアキュムレータ1に
					runtime.Pop, runtime.R1,
					runtime.Mov, runtime.ACM1, runtime.R1,
				}...)
			case 1: // 2つめの戻り値, ACM2に.
				prog = append(prog, runtime.Program{
					// スタックから取り出してアキュムレータ2に
					runtime.Pop, runtime.R1,
					runtime.Mov, runtime.ACM2, runtime.R1,
				}...)
			}
			c = c.next
			count++
		default:
			return nil, fmt.Errorf("genReturn: unsupported value: %s", c.String())
		}
	}

	prog = append(prog, runtime.Program{
		// # 関数の終了処理 #
		runtime.Mov, runtime.StackPointer, runtime.BasePointer,
		runtime.Pop, runtime.BasePointer,
		runtime.Ret,
	}...)

	return prog, nil
}

func genBlock(nd *Node) (runtime.Program, error) {
	prog := runtime.Program{}
	if nd == nil { // elseとかでnilが渡される場合がある
		return prog, nil
	}
	c := &Node{next: nd.lhs}
	for {
		log.Println("block")
		// go next
		if c.next == nil {
			break
		} else {
			c = c.next
		}

		stmt, err := genStatementLevel(c)
		if err != nil {
			return nil, err
		}
		prog = append(prog, stmt...)
	}
	return prog, nil
}

func analyzeFunctionIdent(nd *Node) (int, error) {
	id, err := nd.leaf.GetIdent()
	if err != nil {
		return 0, err
	}

	labelNo, ok := st.FindFn(id)
	if ok {
		return labelNo, nil
	}

	labelNo, err = st.RegisterFn(id)
	if err != nil {
		return 0, err
	}
	return labelNo, nil
}

func genFunctionArguments(nd *Node) (runtime.Program, error) {
	prog := runtime.Program{}
	c := nd.lhs
	argCount := 0
ArgLoop:
	for {
		log.Println("arg")
		switch {
		case c == nil:
			break ArgLoop
		case c.kind == ST_IDENT:
			argName, err := c.leaf.GetIdent()
			if err != nil {
				return nil, err
			}
			sym, err := st.RegisterVar(argName)
			if err != nil {
				return nil, err
			}
			prog = append(prog, runtime.Program{
				// ## 引数と変数の結び付け ##
				runtime.Mov, runtime.NewBPOffset(-sym), runtime.NewBPOffset(2 + argCount),
			}...)
			c = c.next
			argCount++
		default:
			return nil, fmt.Errorf("function argument is not ident: %v", c)
		}
	}

	prog = append(runtime.Program{
		// ## 引数を含む変数領域の確保 ##
		runtime.Push, runtime.Integer(st.TotalVariables()),
		runtime.Pop, runtime.R1,
		runtime.Sub, runtime.StackPointer, runtime.R1,
	}, prog...)

	return prog, nil
}

func genFunctionHeader(nd *Node) (runtime.Program, error) {
	prog := runtime.Program{}
	label, err := analyzeFunctionIdent(nd.lhs)
	if err != nil {
		return nil, err
	}
	args, err := genFunctionArguments(nd.rhs)
	if err != nil {
		return nil, err
	}
	prog = append(prog, runtime.Program{
		// l_func:
		runtime.DefLabel(label),
		// # 関数の初期設定 #
		// ## 現状復帰のための保存 ##
		runtime.Push, runtime.BasePointer,
		runtime.Mov, runtime.BasePointer, runtime.StackPointer,
	}...)
	// ## 引数を含む変数領域の確保 ##
	prog = append(prog, args...)
	return prog, nil
}

func genFunctionDeclaration(nd *Node) (runtime.Program, error) {
	return genFunctionHeader(nd.lhs)
}

func genDefineFunction(nd *Node) (runtime.Program, error) {
	prog := runtime.Program{}
	decl, err := genFunctionDeclaration(nd.lhs)
	if err != nil {
		return nil, err
	}
	block, err := genStatementLevel(nd.rhs)
	if err != nil {
		return nil, err
	}
	prog = append(prog, decl...)
	prog = append(prog, block...)
	return prog, nil
}

func Generate(nd *Node) (runtime.Program, error) {
	curt = &Node{next: nd} // dummy
	st = NewSymbolTable()

	program := runtime.Program{}
	for {
		// go next
		if err := nextNode(); err != nil { // end of nd
			break
		}
		// check toplevel
		prog, err := genToplevel(curt)
		if err != nil {
			return nil, err
		}
		program = append(program, prog...)
	}

	return program, nil
}
