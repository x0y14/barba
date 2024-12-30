package runtime

import (
	"fmt"
	"os"
)

type Runtime struct {
	program Program
	sym     SymbolTable
	reg     []Object
	stack   []Object
	mem     Memory
}

func NewRuntime(stackSize, memSize int) *Runtime {
	r := &Runtime{
		program: nil,
		sym:     *NewSymbolTable(),
		reg:     *NewRegisterSet(),
		stack:   make([]Object, stackSize),
		mem:     *NewMemory(memSize),
	}
	r.setSp(stackSize - 1)
	r.setPc(0)
	r.setBp(0)
	return r
}

func (r *Runtime) Load(program Program) {
	// 擬似的なプロセス呼び出し用コード
	// これがないとmainでretを使えなくなる
	startup := Program{
		// root(l_-1):
		//   call main(l_0)
		//   exit
		DefLabel(-1),
		Call, Label(0),
		Exit,
	}
	program = append(startup, program...)
	r.program = program
	return
}

func (r *Runtime) CollectLabels() error {
	for pc, code := range r.program {
		switch code.(type) {
		case DefLabel:
			if err := r.sym.Set(Label(code.Value()), ProgramAbsoluteOffset(pc)); err != nil {
				return err
			}
		default:
			continue
		}
	}
	return nil
}

func (r *Runtime) Run() error {
	// 擬似的なプロセスからの呼び出し
	entryPoint, err := r.sym.Get(Label(-1))
	if err != nil {
		return err
	}
	r.setPc(entryPoint.Value())
	//
	for {
		if r.mustExit() {
			return nil
		}
		switch code := r.program[r.pc()]; code.(type) {
		case DefLabel:
			r.incPc() // ラベル定義を読み飛ばす
		case Opcode:
			if err := r.do(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported code: %v", code)
		}
	}
}

// ############
// #レジスタ管理#
// ############
// プログラムカウンター
func (r *Runtime) pc() int {
	return r.reg[ProgramCounter].Value()
}
func (r *Runtime) setPc(pc int) {
	r.reg[ProgramCounter] = Integer(pc)
}
func (r *Runtime) incPc() {
	r.setPc(r.pc() + 1)
}

// ベースポインター
func (r *Runtime) bp() int {
	return r.reg[BasePointer].Value()
}
func (r *Runtime) setBp(bp int) {
	r.reg[BasePointer] = Integer(bp)
}

// スタックポインター
func (r *Runtime) sp() int {
	return r.reg[StackPointer].Value()
}
func (r *Runtime) setSp(sp int) {
	r.reg[StackPointer] = Integer(sp)
}

// offsetの計算
func (r *Runtime) calcOffset(offset StackRelativeOffset) int {
	switch offset.target {
	case StackPointer:
		return r.sp() + offset.relativeDistance
	case BasePointer:
		return r.bp() + offset.relativeDistance
	default:
		panic(fmt.Sprintf("unsupported offset: %v", offset))
	}
}

// 終了フラグ
func (r *Runtime) mustExit() bool {
	return r.reg[ExitFlag] == True
}

// ############
// #スタック管理#
// ############
func (r *Runtime) push(obj Object) {
	r.setSp(r.sp() - 1)
	if r.sp() < 0 {
		panic(fmt.Sprintf("stack overflow: stack_size=%d, access=%d", len(r.stack), r.sp()))
	}
	r.stack[r.sp()] = obj
}
func (r *Runtime) pop() Object {
	v := r.stack[r.sp()]
	r.stack[r.sp()] = nil
	r.setSp(r.sp() + 1)
	return v
}

// #####
// #命令#
// #####
func (r *Runtime) do() error {
	switch code := r.program[r.pc()]; code.(Opcode) {
	case Exit:
		r.reg[ExitFlag] = True
		return nil
	case Call: // CALL LABEL
		fnLabel := r.program[r.pc()+1] // call fnLabel <- koko
		switch fnLabel.(type) {
		case Label:
			dest, err := r.sym.Get(fnLabel.(Label))
			if err != nil {
				return err
			}
			r.push(ProgramAbsoluteOffset(r.pc() + 1 + Operand(Call))) // 戻る場所はoffsetとして与える
			r.setPc(dest.Value())
			return nil
		default:
			return fmt.Errorf("unsupported call dest: %v", fnLabel)
		}
	case Ret: // RET
		dest := r.pop()
		switch dest.(type) {
		case ProgramAbsoluteOffset: // offsetが入っているはず
			r.setPc(dest.Value())
			return nil
		default:
			return fmt.Errorf("unsupported ret dest: %v", dest)
		}
	case Jmp: // Jmp Label
		destLabel := r.program[r.pc()+1]
		switch destLabel.(type) {
		case Label:
			dest, err := r.sym.Get(destLabel.(Label))
			if err != nil {
				return err
			}
			r.setPc(dest.Value())
			return nil
		default:
			return fmt.Errorf("unsupported jmp dest: want label, but got: %v", destLabel)
		}
	case Je:
		destLabel := r.program[r.pc()+1]
		switch destLabel.(type) {
		case Label:
			dest, err := r.sym.Get(destLabel.(Label))
			if err != nil {
				return err
			}
			if r.reg[ZeroFlag] == True {
				r.setPc(dest.Value())
				return nil
			}
			r.setPc(r.pc() + 1 + Operand(code.(Opcode)))
			return nil
		default:
			return fmt.Errorf("unsupported je dest: want label, but got: %v", destLabel)
		}
	case Jne:
		destLabel := r.program[r.pc()+1]
		switch destLabel.(type) {
		case Label:
			dest, err := r.sym.Get(destLabel.(Label))
			if err != nil {
				return err
			}
			if r.reg[ZeroFlag] == False {
				r.setPc(dest.Value())
				return nil
			}
			r.setPc(r.pc() + 1 + Operand(code.(Opcode)))
			return nil
		default:
			return fmt.Errorf("unsupported jne dest: want label, but got: %v", destLabel)
		}
	case Mov: // MOV DEST SRC
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		dest := r.program[r.pc()+1]
		src := r.program[r.pc()+2]
		switch dest.(type) {
		case Register:
			switch src.(type) {
			case Register: // reg <- reg
				r.reg[dest.(Register)] = r.reg[src.(Register)]
				return nil
			case StackRelativeOffset: // reg <- offset
				return fmt.Errorf("unsupported mov src: %v", src)
			case Integer, Character, Bool, Null:
				r.reg[dest.(Register)] = src
				return nil
			default:
				return fmt.Errorf("unsupported mov src: %v", src)
			}
		case StackRelativeOffset:
			return fmt.Errorf("unsupported mov dest: %v", dest)
		default:
			return fmt.Errorf("unsupported mov dest: %v", dest)
		}
	case Push:
		defer func() { r.setPc(r.pc() + 1 + Operand(Push)) }()
		switch src := r.program[r.pc()+1]; src.(type) {
		case Register:
			r.push(r.reg[src.(Register)])
			return nil
		case Integer, Character, Bool, Null:
			r.push(src)
			return nil
		default:
			return fmt.Errorf("unsupported push src: %v", src)
		}
	case Pop:
		defer func() { r.setPc(r.pc() + 1 + Operand(Pop)) }()
		switch dest := r.program[r.pc()+1]; dest.(type) {
		case Register:
			r.reg[dest.(Register)] = r.pop()
			return nil
		default:
			return fmt.Errorf("unsupported pop dest: %v", dest)
		}
	case Add:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		dest := r.program[r.pc()+1]
		src := r.program[r.pc()+2]
		switch dest.(type) {
		case Register:
			switch src.(type) {
			case Integer: // reg += int
				switch r.reg[dest.(Register)].(type) { // destがなんなのか確かめる
				case Integer:
					r.reg[dest.(Register)] = Integer(r.reg[dest.(Register)].Value() + src.Value())
					return nil
				default:
					return fmt.Errorf("unsupported add match: %v+=%v", dest, src)
				}
			default:
				return fmt.Errorf("unsupported add src: %v", src)
			}
		default:
			return fmt.Errorf("unsupported add dest: %v", dest)
		}
	case Sub:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		dest := r.program[r.pc()+1]
		src := r.program[r.pc()+2]
		switch dest.(type) {
		case Register:
			switch src.(type) {
			case Integer: // reg -= int
				switch r.reg[dest.(Register)].(type) { // destがなんなのか確かめる
				case Integer:
					r.reg[dest.(Register)] = Integer(r.reg[dest.(Register)].Value() - src.Value())
					return nil
				default:
					return fmt.Errorf("unsupported sub match: %v+=%v", dest, src)
				}
			default:
				return fmt.Errorf("unsupported sub src: %v", src)
			}
		default:
			return fmt.Errorf("unsupported sub dest: %v", dest)
		}
	case Eq:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		lhs := r.program[r.pc()+1]
		rhs := r.program[r.pc()+2]
		switch lhs.(type) {
		case Register:
			switch rhs.(type) {
			case Register: // reg vs reg
				if r.reg[lhs.(Register)] == r.reg[rhs.(Register)] {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null: // reg vs int
				if r.reg[lhs.(Register)] == rhs {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported eq value: %v == %v", lhs, rhs)
			}
		case Integer, Character, Bool, Null:
			switch rhs.(type) {
			case Register:
				if lhs == r.reg[rhs.(Register)] {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null:
				if lhs == rhs {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported eq value: %v == %v", lhs, rhs)
			}
		default:
			return fmt.Errorf("unsupported eq value: %v == %v", lhs, rhs)
		}
	case Ne:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		lhs := r.program[r.pc()+1]
		rhs := r.program[r.pc()+2]
		switch lhs.(type) {
		case Register:
			switch rhs.(type) {
			case Register: // reg vs reg
				if r.reg[lhs.(Register)] != r.reg[rhs.(Register)] {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null: // reg vs int
				if r.reg[lhs.(Register)] != rhs {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported ne value: %v == %v", lhs, rhs)
			}
		case Integer, Character, Bool, Null:
			switch rhs.(type) {
			case Register:
				if lhs != r.reg[rhs.(Register)] {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null:
				if lhs != rhs {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported ne value: %v == %v", lhs, rhs)
			}
		default:
			return fmt.Errorf("unsupported ne value: %v == %v", lhs, rhs)
		}
	case Lt:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		lhs := r.program[r.pc()+1]
		rhs := r.program[r.pc()+2]
		switch lhs.(type) {
		case Register:
			switch rhs.(type) {
			case Register: // reg vs reg
				if r.reg[lhs.(Register)].Value() < r.reg[rhs.(Register)].Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null: // reg vs int
				if r.reg[lhs.(Register)].Value() < rhs.Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported lt value: %v == %v", lhs, rhs)
			}
		case Integer, Character, Bool, Null:
			switch rhs.(type) {
			case Register:
				if lhs.Value() < r.reg[rhs.(Register)].Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null:
				if lhs.Value() < rhs.Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported lt value: %v == %v", lhs, rhs)
			}
		default:
			return fmt.Errorf("unsupported lt value: %v == %v", lhs, rhs)
		}
	case Le:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		lhs := r.program[r.pc()+1]
		rhs := r.program[r.pc()+2]
		switch lhs.(type) {
		case Register:
			switch rhs.(type) {
			case Register: // reg vs reg
				if r.reg[lhs.(Register)].Value() <= r.reg[rhs.(Register)].Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null: // reg vs int
				if r.reg[lhs.(Register)].Value() <= rhs.Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported le value: %v == %v", lhs, rhs)
			}
		case Integer, Character, Bool, Null:
			switch rhs.(type) {
			case Register:
				if lhs.Value() <= r.reg[rhs.(Register)].Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			case Integer, Character, Bool, Null:
				if lhs.Value() <= rhs.Value() {
					r.reg[ZeroFlag] = True
				} else {
					r.reg[ZeroFlag] = False
				}
				return nil
			default:
				return fmt.Errorf("unsupported le value: %v == %v", lhs, rhs)
			}
		default:
			return fmt.Errorf("unsupported le value: %v == %v", lhs, rhs)
		}
	case Syscall:
		defer func() { r.setPc(r.pc() + 1 + Operand(code.(Opcode))) }()
		syscallNo := r.program[r.pc()+1]   // Write, ...
		syscallArg1 := r.program[r.pc()+2] // Stdout, ...
		syscallArg2 := r.program[r.pc()+3] // "Hello", ...
		switch syscallNo.(type) {
		case SystemCall:
			switch syscallNo.(SystemCall) {
			case Write:
				var f *os.File
				switch syscallArg1.(StandardIO) {
				case StdIn:
					panic(fmt.Errorf("unsupported syscall write dest: %v", StdIn))
				case StdOut:
					f = os.Stdout
				case StdErr:
					f = os.Stderr
				}
				_, err := fmt.Fprintf(f, syscallArg2.String())
				return err
			default:
				return fmt.Errorf("unsupported syscall number: %v", syscallNo)
			}
		default:
			return fmt.Errorf("unsupported syscall want type(syscall), but got: %v", syscallNo)
		}
	default:
		return fmt.Errorf("unsupported opcode: %v", code)
	}
}
