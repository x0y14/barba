package runtime

type ProgramAbsoluteOffset int

func (p ProgramAbsoluteOffset) Value() int {
	return int(p)
}

type StackRelativeOffset struct {
	target           Register
	relativeDistance int
}

func (s StackRelativeOffset) Value() int {
	return s.relativeDistance
}
