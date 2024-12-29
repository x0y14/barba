package runtime

type StandardIO int

func (s StandardIO) Value() int {
	return int(s)
}

const (
	StdIn StandardIO = iota
	StdOut
	StdErr
)
