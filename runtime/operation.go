package runtime

type Operation int

func (o Operation) Value() int {
	return int(o)
}

const (
	Exit Operation = iota

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

func Operand(op Operation) int {
	switch op {
	case Exit, Ret:
		return 0
	case Push, Pop, Call, Jmp, Je, Jne:
		return 1
	case Mov, Add, Sub, Eq, Ne, Lt, Le, Syscall:
		return 2
	default:
		return 0
	}
}
