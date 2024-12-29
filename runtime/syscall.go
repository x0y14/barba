package runtime

type SystemCall int

func (s SystemCall) Value() int {
	return int(s)
}

const (
	Write SystemCall = iota
)
