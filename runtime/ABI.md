# 呼び出し規則
intel記法っぽいものを使います．  
`MOV DEST SRC` // DEST = SRC  
`ADD DEST SRC` // DEST += SRC  

## 関数の引数
呼び出し前にスタックに逆順に入れてください．(要検証)

## 計算結果
スタックに入れます．
## 関数の戻り値
(レジスタのみを用いる場合)最大2つです．  
1つの場合はR10, 2つの場合は加えてR11を使います．  
```go
// 戻り値が1つの場合
func f() int { return 1 } // R10に1が入ります．
// 戻り値が2つの場合
func f() (int, bool) { return 1, true } // R10に1が，R11にtrueが入ります．
```

以下フィボナッチ数列の10番目までの合計を計算する関数のコードと解説です．  
arrttyで動作したコードを書き起こしたものなのでbarbaの動作と完全に一致しているわけではありません．  
Barbaでコンパイラが作成されたタイミング更新します．
```text
// 変換前
func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
func main() int {
	return fib(10)
}
```
```text
// 変換後
// 関数の定義、関数ラベルの作成
fib:
	// メイン関数でなければ、関数終了時の戻り場所を記録
	push bp
	mov bp sp
	// 関数内で使用される変数の数だけSPを下げる(変数領域の確保)
	push 1
	pop r1
	sub sp r1
	// 引数と変数を結びつける(代入によって)
	mov [bp-1] [bp+2]
	// [[ IFの条件 ]]
	// << lt >>
	// 左辺をプッシュ
	push [bp-1]
	// 右辺をプッシュ
	push 2
	// 右辺の取り出し
	pop r2
	// 左辺の取り出し
	pop r1
	// 比較
	lt r1 r2
	// 結果が真ならIFブロックへ
    jz if_if_block_jGcBdPTUNWrbPQrSxTuS
	// そうでないならRETURNブロックへ
    jmp if_end_gNEXgcrQiVuTXurmOGFW

	// [[ IFブロック ]]
	if_if_block_jGcBdPTUNWrbPQrSxTuS:
		push [bp-1]
		// 計算結果をR1に移動
		pop r1
		mov r10 r1
		// リターン本文
		mov sp bp
		pop bp
		ret
	
	// [[ RETURNブロック ]]
	if_end_gNEXgcrQiVuTXurmOGFW:
		// << sub >>
		// 左辺をプッシュ
		push [bp-1]
		// 右辺をプッシュ
		push 1
		// 右辺の取り出し
		pop r2
		// 左辺の取り出し
		pop r1
		// 引き算
		sub r2 r1
		// 結果をプッシュ
		push r1

		// << call >>
		// 呼び出し
		call fib
		// 引数分spを加算
		push 1
		pop r1
		add r1 sp
		// 戻り値を結果としてプッシュ
		push r10 // 答え1

		<< sub >>
		// 左辺をプッシュ
		push [bp-1]
		// 右辺をプッシュ
		push 2
		// 右辺を取り出し
		pop r2
		// 左辺を取り出し
		pop r1
		// 引き算
		sub r1 r2
		// 結果をプッシュ
		push r1

		// << call >>
		// 呼び出し
		call fib
		// 引数分spを加算
		push 1
		pop r1
		add sp r1
		// 戻り値を結果としてプッシュ
		push r10 // 答え2

		// 直前の計算結果つまり'答え2'を右辺として
		pop r2
		// 一個前の結果つまり'答え1'を左辺として取り出す
		pop r1
		// 足し算
		add r1 r2
		// 結果をスタックへ
		push r1

		// 計算結果をR1に移動
		pop r1
		mov r10 r1
		// リターン本文
		mov sp bp
		pop bp
		ret

// 関数の定義、関数ラベルの作成
main:
	// メイン関数でなければ、関数終了時の戻り場所を記録
	// 戻り場所に関する処理
	// 

	// 関数内で使用される変数の数だけSPを下げる(変数領域の確保)
	push 0
	pop r1
	sub sp r1
	
	// 引数をプッシュ
	push 10
	// << call >>
	// 呼び出し
	call fib
	// 引数分spを加算
	push 1
	pop r1
	add sp r1
	// 戻り値を結果としてプッシュ
	push r10
	// 計算結果をR1に移動
	pop r1
	mov r10 r1
```