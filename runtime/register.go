package runtime

type Register int

func (r Register) Value() int {
	return int(r)
}

const (
	ProgramCounter Register = iota
	BasePointer
	StackPointer
	ZeroFlag
	ExitFlag
	General1
	General2
	Temporal1
)
