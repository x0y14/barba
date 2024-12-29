package runtime

type Offset struct {
	baseAddress      int
	relativeDistance int
}

func (o Offset) Addr() int {
	return o.baseAddress + o.relativeDistance
}

// [bp-1] => Offset(BasePointer.Value, -1)
// [sp+2] => Offset(StackPointer.Value, 2)
