package runtime

type Label int

func (l Label) Value() int {
	return int(l)
}
