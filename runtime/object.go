package runtime

type Object interface {
	Value() int
}

type Integer int

func (i Integer) Value() int {
	return int(i)
}

type Character int

func (c Character) Value() int {
	return int(c)
}

type Bool bool

func (b Bool) Value() int {
	if b == true {
		return 1
	}
	return 0
}

const (
	True  Bool = true
	False Bool = false
)

type Null struct{}

func (n Null) Value() int {
	return 0
}

type List int

func (l List) Value() int {
	return int(l)
}
