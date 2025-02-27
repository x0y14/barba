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

func NewConstant(primary Underlying, sub *Constant) *Constant {
	return &Constant{
		Underlying: primary,
		Next:       sub,
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
func NewArray(sub *Constant) *Constant {
	return NewConstant(Array, sub)
}

// NewMap map[string]intはMap->string->Intで表現します．
func NewMap(sub *Constant) *Constant {
	return NewConstant(Map, sub)
}

// NewConstantFromStr "string"からString, "[]string"からArray->Stringのように文字列から与えられた型データを
// 解析してConstantにして返却します．
//func NewConstantFromStr(rawConstant string) (*Constant, error) {
//	runes := []rune(rawConstant)
//	tokens := tokenize(runes)
//
//	return nil, nil
//}
