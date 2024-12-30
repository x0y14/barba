package runtime

type Register int

func (r Register) Value() int {
	return int(r)
}
func (r Register) String() string {
	var kinds = [...]string{
		ProgramCounter: "pc",
		BasePointer:    "bp",
		StackPointer:   "sp",
		ZeroFlag:       "zf",
		ExitFlag:       "ef",
		General1:       "g1",
		General2:       "g2",
		Temporal1:      "t1",
		_reg_end:       "",
	}
	return kinds[r]
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
