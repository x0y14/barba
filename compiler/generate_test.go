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
					lhs:  nil, // mainは引数なし
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

func TestGenerate_IF_ELSE(t *testing.T) {

	tests := []struct {
		name         string
		ir           *Node
		expectAsm    runtime.Program
		expectStatus int
	}{
		{
			//	func main() int {
			//		if 1 == 1 {
			//			return 1
			//		} else {
			//			return 2
			//		}
			//		return 3
			//	}
			"if_true",
			&Node{ // `func main() int {...}`
				kind: ST_DEFINE_FUNCTION,
				lhs: &Node{ // `func main() int`
					kind: ST_FUNCTION_DECLARATION,
					lhs: &Node{ // `main()`
						kind: ST_FUNCTION_HEADER,
						lhs: &Node{ // `main`
							kind: ST_IDENT,
							leaf: NewToken(TK_IDENT, "main"),
						},
						rhs: &Node{ // `()`
							kind: ST_FUNCTION_ARGUMENTS,
							lhs:  nil, // mainは引数なし
						},
					},
					rhs: &Node{ // `int`
						kind: ST_FUNCTION_RETURN_DETAIL,
						// 本来であれば戻り値の情報がここにある
					},
				},
				rhs: &Node{ // `{...}`
					kind: ST_BLOCK,
					lhs: &Node{ // `if 1 == 1 {...} else {...}`
						kind: ST_IF_ELSE,
						lhs: &Node{ // `if 1 == 1 {...}`
							kind: ST_IF,
							lhs: &Node{ // `1 == 1`
								kind: ST_EQ,
								lhs: &Node{ // `1` (左辺)
									kind: ST_PRIMITIVE,
									lhs: &Node{
										kind: ST_INTEGER,
										leaf: NewToken(TK_INT, "1"),
									},
								},
								rhs: &Node{ // `1` (右辺)
									kind: ST_PRIMITIVE,
									lhs: &Node{
										kind: ST_INTEGER,
										leaf: NewToken(TK_INT, "1"),
									},
								},
							},
							rhs: &Node{ // `{...}`
								kind: ST_BLOCK,
								lhs: &Node{ // return 1
									kind: ST_RETURN,
									lhs: &Node{
										kind: ST_PRIMITIVE,
										lhs: &Node{
											kind: ST_INTEGER,
											leaf: NewToken(TK_INT, "1"),
										},
									},
								},
							},
						},
						rhs: &Node{ // `else {...}`
							kind: ST_BLOCK,
							lhs: &Node{ // return 2
								kind: ST_RETURN,
								lhs: &Node{
									kind: ST_PRIMITIVE,
									lhs: &Node{
										kind: ST_INTEGER,
										leaf: NewToken(TK_INT, "2"),
									},
								},
							},
						},

						next: &Node{ // return 3
							kind: ST_RETURN,
							lhs: &Node{
								kind: ST_PRIMITIVE,
								lhs: &Node{
									kind: ST_INTEGER,
									leaf: NewToken(TK_INT, "3"),
								},
							},
						},
					},
				},
			},
			runtime.Program{
				// main:
				runtime.DefLabel(0),
				//
				// # 関数の初期設定 #
				// ## 現状復帰のための保存 ##
				runtime.Push, runtime.BasePointer,
				runtime.Mov, runtime.BasePointer, runtime.StackPointer,
				// ## 引数を含む変数領域の確保 ##
				runtime.Push, runtime.Integer(0),
				runtime.Pop, runtime.R1,
				runtime.Sub, runtime.StackPointer, runtime.R1,
				// ## 変数と引数の結び付け ##
				//
				// # 本体 #
				// < if-beg >
				// ## 条件式 ##
				// ### 左辺の計算 ###
				runtime.Push, runtime.Integer(1), // 即値の評価
				// ### 右辺の計算 ###
				runtime.Push, runtime.Integer(1), // 即値の評価
				// ## 条件の評価 ##
				// < eq-beg >
				runtime.Pop, runtime.R2, // 右辺の取り出し
				runtime.Pop, runtime.R1, // 左辺の取り出し
				runtime.Eq, runtime.R1, runtime.R2,
				runtime.Push, runtime.ZeroFlag,
				// < eq-end >
				// ## 条件分岐 ##
				runtime.Pop, runtime.Temporal1, // ZFを使うので結果は捨てる
				runtime.Je, runtime.Label(GETA_LABEL + 1), // if
				runtime.Jmp, runtime.Label(GETA_LABEL + 2), // else
				// if_1_if:
				runtime.DefLabel(GETA_LABEL + 1),
				// return 1
				// ## 戻り値の準備 ##
				// 即値の評価
				runtime.Push, runtime.Integer(1),
				// < return-beg >
				runtime.Pop, runtime.R1,
				runtime.Mov, runtime.ACM1, runtime.R1,
				// < return-end >
				// # 関数の終了処理 #
				runtime.Mov, runtime.StackPointer, runtime.BasePointer,
				runtime.Pop, runtime.BasePointer,
				runtime.Ret,
				runtime.Jmp, runtime.Label(GETA_LABEL + 3), // ifの終了へ
				// if_1_else:
				runtime.DefLabel(GETA_LABEL + 2),
				// return 2
				// ## 戻り値の準備 ##
				// 即値の評価
				runtime.Push, runtime.Integer(2),
				// < return-beg >
				runtime.Pop, runtime.R1,
				runtime.Mov, runtime.ACM1, runtime.R1,
				// < return-end >
				// # 関数の終了処理 #
				runtime.Mov, runtime.StackPointer, runtime.BasePointer,
				runtime.Pop, runtime.BasePointer,
				runtime.Ret,
				runtime.Jmp, runtime.Label(GETA_LABEL + 3), // ifの終了へ
				// if_1_end:
				runtime.DefLabel(GETA_LABEL + 3),
				// < if-end >
				//
				// ## 戻り値の準備 ##
				// 即値の評価
				runtime.Push, runtime.Integer(3),
				// < return-beg >
				runtime.Pop, runtime.R1,
				runtime.Mov, runtime.ACM1, runtime.R1,
				// < return-end >
				//
				// # 関数の終了処理 #
				runtime.Mov, runtime.StackPointer, runtime.BasePointer,
				runtime.Pop, runtime.BasePointer,
				runtime.Ret,
			},
			1,
		},
		{
			//	func main() int {
			//		if 1 == 2 {
			//			return 1
			//		} else {
			//			return 2
			//		}
			//		return 3
			//	}
			"if_false",
			&Node{ // `func main() int {...}`
				kind: ST_DEFINE_FUNCTION,
				lhs: &Node{ // `func main() int`
					kind: ST_FUNCTION_DECLARATION,
					lhs: &Node{ // `main()`
						kind: ST_FUNCTION_HEADER,
						lhs: &Node{ // `main`
							kind: ST_IDENT,
							leaf: NewToken(TK_IDENT, "main"),
						},
						rhs: &Node{ // `()`
							kind: ST_FUNCTION_ARGUMENTS,
							lhs:  nil, // mainは引数なし
						},
					},
					rhs: &Node{ // `int`
						kind: ST_FUNCTION_RETURN_DETAIL,
						// 本来であれば戻り値の情報がここにある
					},
				},
				rhs: &Node{ // `{...}`
					kind: ST_BLOCK,
					lhs: &Node{ // `if 1 == 1 {...} else {...}`
						kind: ST_IF_ELSE,
						lhs: &Node{ // `if 1 == 1 {...}`
							kind: ST_IF,
							lhs: &Node{ // `1 == 1`
								kind: ST_EQ,
								lhs: &Node{ // `1` (左辺)
									kind: ST_PRIMITIVE,
									lhs: &Node{
										kind: ST_INTEGER,
										leaf: NewToken(TK_INT, "1"),
									},
								},
								rhs: &Node{ // `2` (右辺)
									kind: ST_PRIMITIVE,
									lhs: &Node{
										kind: ST_INTEGER,
										leaf: NewToken(TK_INT, "2"),
									},
								},
							},
							rhs: &Node{ // `{...}`
								kind: ST_BLOCK,
								lhs: &Node{ // return 1
									kind: ST_RETURN,
									lhs: &Node{
										kind: ST_PRIMITIVE,
										lhs: &Node{
											kind: ST_INTEGER,
											leaf: NewToken(TK_INT, "1"),
										},
									},
								},
							},
						},
						rhs: &Node{ // `else {...}`
							kind: ST_BLOCK,
							lhs: &Node{ // return 2
								kind: ST_RETURN,
								lhs: &Node{
									kind: ST_PRIMITIVE,
									lhs: &Node{
										kind: ST_INTEGER,
										leaf: NewToken(TK_INT, "2"),
									},
								},
							},
						},

						next: &Node{ // return 3
							kind: ST_RETURN,
							lhs: &Node{
								kind: ST_PRIMITIVE,
								lhs: &Node{
									kind: ST_INTEGER,
									leaf: NewToken(TK_INT, "3"),
								},
							},
						},
					},
				},
			},
			runtime.Program{
				// main:
				runtime.DefLabel(0),
				//
				// # 関数の初期設定 #
				// ## 現状復帰のための保存 ##
				runtime.Push, runtime.BasePointer,
				runtime.Mov, runtime.BasePointer, runtime.StackPointer,
				// ## 引数を含む変数領域の確保 ##
				runtime.Push, runtime.Integer(0),
				runtime.Pop, runtime.R1,
				runtime.Sub, runtime.StackPointer, runtime.R1,
				// ## 変数と引数の結び付け ##
				//
				// # 本体 #
				// < if-beg >
				// ## 条件式 ##
				// ### 左辺の計算 ###
				runtime.Push, runtime.Integer(1), // 即値の評価
				// ### 右辺の計算 ###
				runtime.Push, runtime.Integer(2), // 即値の評価
				// ## 条件の評価 ##
				// < eq-beg >
				runtime.Pop, runtime.R2, // 右辺の取り出し
				runtime.Pop, runtime.R1, // 左辺の取り出し
				runtime.Eq, runtime.R1, runtime.R2,
				runtime.Push, runtime.ZeroFlag,
				// < eq-end >
				// ## 条件分岐 ##
				runtime.Pop, runtime.Temporal1, // ZFを使うので結果は捨てる
				runtime.Je, runtime.Label(GETA_LABEL + 1), // if
				runtime.Jmp, runtime.Label(GETA_LABEL + 2), // else
				// if_1_if:
				runtime.DefLabel(GETA_LABEL + 1),
				// return 1
				// ## 戻り値の準備 ##
				// 即値の評価
				runtime.Push, runtime.Integer(1),
				// < return-beg >
				runtime.Pop, runtime.R1,
				runtime.Mov, runtime.ACM1, runtime.R1,
				// < return-end >
				// # 関数の終了処理 #
				runtime.Mov, runtime.StackPointer, runtime.BasePointer,
				runtime.Pop, runtime.BasePointer,
				runtime.Ret,
				runtime.Jmp, runtime.Label(GETA_LABEL + 3), // ifの終了へ
				// if_1_else:
				runtime.DefLabel(GETA_LABEL + 2),
				// return 2
				// ## 戻り値の準備 ##
				// 即値の評価
				runtime.Push, runtime.Integer(2),
				// < return-beg >
				runtime.Pop, runtime.R1,
				runtime.Mov, runtime.ACM1, runtime.R1,
				// < return-end >
				// # 関数の終了処理 #
				runtime.Mov, runtime.StackPointer, runtime.BasePointer,
				runtime.Pop, runtime.BasePointer,
				runtime.Ret,
				runtime.Jmp, runtime.Label(GETA_LABEL + 3), // ifの終了へ
				// if_1_end:
				runtime.DefLabel(GETA_LABEL + 3),
				// < if-end >
				//
				// ## 戻り値の準備 ##
				// 即値の評価
				runtime.Push, runtime.Integer(3),
				// < return-beg >
				runtime.Pop, runtime.R1,
				runtime.Mov, runtime.ACM1, runtime.R1,
				// < return-end >
				//
				// # 関数の終了処理 #
				runtime.Mov, runtime.StackPointer, runtime.BasePointer,
				runtime.Pop, runtime.BasePointer,
				runtime.Ret,
			},
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := Generate(tt.ir)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectAsm, prog)
			rt := runtime.NewRuntime(10, 10)
			rt.Load(prog)
			assert.Nil(t, rt.CollectLabels())
			assert.Nil(t, rt.Run())
			assert.Equal(t, tt.expectStatus, rt.Status())
		})
	}

}
