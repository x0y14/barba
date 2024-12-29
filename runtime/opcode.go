package runtime

type Opcode int

func (o Opcode) Value() int {
	return int(o)
}

const (
	Nop Opcode = iota

	Exit

	Mov
	Push
	Pop

	Call
	Ret

	Add
	Sub

	Jmp
	Je
	Jne

	Eq
	Ne
	Lt
	Le

	Syscall
)

func Operand(op Opcode) int {
	switch op {
	case Nop, Exit, Ret:
		return 0
	case Push, Pop, Call, Jmp, Je, Jne:
		return 1
	case Mov, Add, Sub, Eq, Ne, Lt, Le, Syscall:
		return 2
	default:
		return 0
	}
}
