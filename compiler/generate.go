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
	c := &Node{next: nd.lhs}
	prog := runtime.Program{}
	for {
		log.Println("block")
		// go next
		if c.next == nil {
			break
		} else {
			c = c.next
		}

		switch c.kind {
		case ST_RETURN:
			stmt, err := genReturn(c)
			if err != nil {
				return nil, err
			}
			prog = append(prog, stmt...)
		default:
			return nil, fmt.Errorf("unsupported statement: %v", c)
		}
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
	block, err := genBlock(nd.rhs)
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

		switch curt.kind {
		case ST_DEFINE_FUNCTION:
			prog, err := genDefineFunction(curt)
			if err != nil {
				return nil, err
			}
			program = append(program, prog...)
		default:
			return nil, fmt.Errorf("unsupported syntax: %v", curt.kind.String())
		}
	}

	return program, nil
}
