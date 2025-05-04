package typemetadata

type Underlying int

const (
	Undef Underlying = iota
	Int
	Float
	String
	Bool
	Array
	Map
)

type Constant struct {
	Underlying
	Next *Constant
}

func NewConstant(primary Underlying, next *Constant) *Constant {
	return &Constant{
		Underlying: primary,
		Next:       next,
	}
}

func NewUndef() *Constant {
	return NewConstant(Undef, nil)
}

func NewInt() *Constant {
	return NewConstant(Int, nil)
}

func NewFloat() *Constant {
	return NewConstant(Float, nil)
}

func NewString() *Constant {
	return NewConstant(String, nil)
}

func NewBool() *Constant {
	return NewConstant(Bool, nil)
}

// NewArray []stringはArray->Stringで表現します．
func NewArray(typ *Constant) *Constant {
	return NewConstant(Array, typ)
}

// NewMap map[string]intはMap->string->Intで表現します．
// keyの要素を最後まで展開した後，valueを展開します．
func NewMap(ktyp *Constant, vtyp *Constant) *Constant {
	c := NewConstant(Map, ktyp)
	tail := c
	for tail.Next != nil {
		tail = tail.Next
	}
	tail.Next = vtyp
	return c
}

// NewConstantFromStr "string"からString, "[]string"からArray->Stringのように文字列から与えられた型データを
// 解析してConstantにして返却します．
func NewConstantFromStr(rawConstant string) (*Constant, error) {
	runes := []rune(rawConstant)
	tokens, err := Tokenize(runes)
	if err != nil {
		return nil, err
	}
	return Parse(tokens)
}
