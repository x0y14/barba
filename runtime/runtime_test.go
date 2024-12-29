package runtime

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRuntime(t *testing.T) {
	// スタックは自動で設定される
	reg := *NewRegisterSet()
	reg[StackPointer] = Integer(1 - 1)
	reg[ProgramCounter] = Integer(0)
	reg[BasePointer] = Integer(0)
	assert.Equal(t, &Runtime{
		program: nil,
		sym:     *NewSymbolTable(),
		reg:     reg,
		stack:   make([]Object, 1),
		mem:     make(Memory, 1),
	}, NewRuntime(1, 1))
	//
	reg = *NewRegisterSet()
	reg[StackPointer] = Integer(15 - 1)
	reg[ProgramCounter] = Integer(0)
	reg[BasePointer] = Integer(0)
	assert.Equal(t, &Runtime{
		program: nil,
		sym:     *NewSymbolTable(),
		reg:     reg,
		stack:   make([]Object, 15),
		mem:     make(Memory, 15),
	}, NewRuntime(15, 15))
}

func TestRuntime_CollectLabels(t *testing.T) {
	rt := NewRuntime(1, 1)
	rt.Load(Program{
		// main:
		//   ret
		DefLabel(0),
		Ret,
	})
	assert.Nil(t, rt.CollectLabels())
	// 擬似プロセスコードが挿入されるのでコードが追加される
	assert.Equal(t, Program{
		DefLabel(-1),   // ここから->
		Call, Label(0), //
		Exit,        // <- ここまで追加
		DefLabel(0), // ここ以降はユーザーが読み込ませたプログラム
		Ret,
	}, rt.program)
	assert.Equal(t, ProgramAbsoluteOffset(0), rt.sym[-1]) // root label, DefLabel(-1)の配列での位置
	assert.Equal(t, ProgramAbsoluteOffset(4), rt.sym[0])  // main label, DefLabel(0)の配列での位置
}

func TestRuntime_Run_Exit(t *testing.T) {
	// Mainなし
	rt := NewRuntime(1, 1)
	rt.program = Program{
		Exit,
	}
	rt.sym[Label(-1)] = ProgramAbsoluteOffset(0) // RunでcallされるものをCollectLabelsの代わりにセットしてあげる
	err := rt.Run()
	assert.Nil(t, err)
	// Mainあり
	rt = NewRuntime(2, 1)
	rt.Load(Program{
		DefLabel(0),
		Ret,
	})
	err = rt.CollectLabels()
	assert.Nil(t, err)
	err = rt.Run()
	assert.Nil(t, err)
}

func TestRuntime_Run_Call(t *testing.T) {
	rt := NewRuntime(10, 10)
	rt.Load(Program{
		DefLabel(5),
		Ret,

		DefLabel(4),
		Call, Label(5),
		Ret,

		DefLabel(3),
		Call, Label(4),
		Ret,

		DefLabel(2),
		Call, Label(3),
		Ret,

		DefLabel(1),
		Call, Label(2),
		Ret,

		DefLabel(0),
		Call, Label(1),
		Ret,
	})
	err := rt.CollectLabels()
	assert.Nil(t, err)
	err = rt.Run()
	assert.Nil(t, err)
}
