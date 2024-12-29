package runtime

type Offset struct {
	target     Register
	distance   int
	isAbsolute bool
}

// [bp-1] => Offset(BasePointer.Value, -1)
// [sp+2] => Offset(StackPointer.Value, 2)

func NewOffset(target Register, distance int, isAbs bool) *Offset {
	return &Offset{target: target, distance: distance, isAbsolute: isAbs}
}
func NewRelativeOffset(target Register, distance int) *Offset {
	return NewOffset(target, distance, false)
}
func NewAbsoluteOffset(target Register, distance int) *Offset {
	return NewOffset(target, distance, true)
}

func (o Offset) Value() int {
	return o.distance
}

func (o Offset) IsPc() bool {
	return o.target == ProgramCounter
}
func (o Offset) IsSp() bool {
	return o.target == StackPointer
}
func (o Offset) IsBp() bool {
	return o.target == BasePointer
}

//
//func AbsoluteOffset(target Register) Offset {
//	return Offset{target: target, distance: 0}
//}
