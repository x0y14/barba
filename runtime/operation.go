package runtime

type Operation int

func (o Operation) Value() int {
	return int(o)
}

const (
	Illegal Operation = iota
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
