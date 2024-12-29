package runtime

type Label int

func (l Label) Value() int {
	return int(l)
}

type DefLabel int

func (d DefLabel) Value() int {
	return int(d)
}
