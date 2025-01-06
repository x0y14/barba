package compiler

import (
	"barba/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate_Return100(t *testing.T) {
	n := &Node{ // `func main() int { return 100 }`
		kind: ST_DEFINE_FUNCTION,
		lhs: &Node{ // `func main() int`
			kind: ST_FUNCTION_DECLARATION,
			lhs: &Node{ // `func main()`
				kind: ST_FUNCTION_HEADER,
				lhs: &Node{ // `main`
					kind: ST_IDENT,
					leaf: NewToken(TK_IDENT, "main"),
				},
				rhs: &Node{ // `()`
					kind: ST_FUNCTION_ARGUMENTS,
					// lhs: nil, // mainは引数なし
				},
			},
			rhs: &Node{ // `int`
				kind: ST_FUNCTION_RETURN_DETAIL,
				// 本来であれば戻り値の情報がここにある
			},
		},
		rhs: &Node{ // `{ return 100 }`
			kind: ST_BLOCK,
			lhs: &Node{ // `return 100`
				kind: ST_RETURN,
				lhs: &Node{ //`100`
					kind: ST_PRIMITIVE,
					lhs: &Node{
						kind: ST_INTEGER,
						leaf: NewToken(TK_INT, "100"),
					},
				},
			},
		},
	}
	prog, err := Generate(n)
	assert.Nil(t, err)
	assert.Equal(t, runtime.Program{
		// func main() int {
		//     return 100
		// }

		// main:
		runtime.DefLabel(0),
		// # 関数の初期設定 #
		// ## 現状復帰のための保存 ##
		runtime.Push, runtime.BasePointer,
		runtime.Mov, runtime.BasePointer, runtime.StackPointer,
		// ## 引数を含む変数領域の確保 ##
		runtime.Push, runtime.Integer(0), // 今テストのmainは引数も変数も使用しない
		runtime.Pop, runtime.R1,
		runtime.Sub, runtime.StackPointer, runtime.R1,
		// # 戻り値の返却 #
		// 即値の評価
		runtime.Push, runtime.Integer(100),
		runtime.Pop, runtime.R1,
		runtime.Mov, runtime.ACM1, runtime.R1,
		// # 関数の終了処理 #
		runtime.Mov, runtime.StackPointer, runtime.BasePointer,
		runtime.Pop, runtime.BasePointer,
		runtime.Ret,
	}, prog)
	rt := runtime.NewRuntime(10, 10)
	rt.Load(prog)
	assert.Nil(t, rt.CollectLabels())
	assert.Nil(t, rt.Run())
	assert.Equal(t, 100, rt.Status())
}
