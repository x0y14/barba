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

	_reg_end
)

func NewRegisterSet() *[]Object {
	rSet := make([]Object, _reg_end-1)
	return &rSet
}
