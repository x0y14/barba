package compiler

import (
	"fmt"
	"maps"
)

type SymbolTable struct {
	fns map[string]int // fnName: labelNo

	curtFn   string
	curtNest int
	vars     map[string]map[int]map[string]int // fnName: nest: varName
	labels   map[string]map[string]int         // fnName: labelName: labelNo
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		fns:      make(map[string]int),
		curtFn:   "",
		curtNest: 0,
		vars:     make(map[string]map[int]map[string]int),
		labels:   make(map[string]map[string]int),
	}
}

// Functions

func (st *SymbolTable) RegisterFn(fnName string) (int, error) {
	_, ok := st.FindFn(fnName)
	if ok {
		return 0, fmt.Errorf("func alredy exists: %s", fnName)
	}
	var no int
	if fnName == "main" {
		no = 0
	} else {
		no = len(st.fns) + 1 // 連番
	}
	st.fns[fnName] = no
	return no, nil
}

func (st *SymbolTable) FindFn(fnName string) (int, bool) {
	no, ok := st.fns[fnName]
	return no, ok
}

// Variables

func (st *SymbolTable) RegisterVar(varName string) (int, error) {
	// すでに定義されていないかチェック, 重複も上書きもだめなのでエラー.
	_, ok := st.FindVar(varName)
	if ok {
		return 0, fmt.Errorf("var alredy exists: %s.%s", st.curtFn, varName)
	}
	// なかった場合
	// 関数用の領域があるか
	_, ok = st.vars[st.curtFn]
	if !ok {
		st.vars[st.curtFn] = make(map[int]map[string]int)
	}
	// 関数用領域の中にネストされた領域があるか
	_, ok = st.vars[st.curtFn][st.curtNest]
	if !ok {
		st.vars[st.curtFn][st.curtNest] = make(map[string]int)
	}
	// 連番してあげる
	st.vars[st.curtFn][st.curtNest][varName] = len(st.vars[st.curtFn][st.curtNest]) + 1
	return st.vars[st.curtFn][st.curtNest][varName], nil
}

func (st *SymbolTable) FindVar(varName string) (int, bool) {
	// そもそも関数用の領域がない
	_, ok := st.vars[st.curtFn]
	if !ok {
		return 0, false
	}
	//_, ok = st.vars[st.curtFn][st.curtNest]
	//if !ok {
	//	return 0, false
	//}
	//distance, ok := st.vars[st.curtFn][st.curtNest][varName]
	//return distance, ok

	for nest := range maps.Keys(st.vars[st.curtFn]) {
		if nest <= st.curtNest { // 同じかそれより浅いものしか参照できない
			for registeredVarName, distance := range st.vars[st.curtFn][nest] {
				if registeredVarName == varName {
					return distance, true
				}
			}
		}
	}

	return 0, false
}

func (st *SymbolTable) TotalVariables() int {
	// そもそも関数用の領域がない
	_, ok := st.vars[st.curtFn]
	if !ok {
		return 0
	}
	count := 0
	for nest := range maps.Keys(st.vars[st.curtFn]) {
		for range st.vars[st.curtFn][nest] {
			count++
		}
	}
	return count
}

// Labels

func (st *SymbolTable) RegisterLabel(label string) (int, error) {
	_, ok := st.FindLabel(label)
	if ok {
		return 0, fmt.Errorf("label alredy exists: %s.l_%s", st.curtFn, label)
	}
	// なかったら作ってあげる
	_, ok = st.labels[st.curtFn]
	if !ok {
		st.labels[st.curtFn] = make(map[string]int)
	}
	no := len(st.labels[st.curtFn]) + 1 // 連番
	st.labels[st.curtFn][label] = no
	return no, nil
}

func (st *SymbolTable) FindLabel(label string) (int, bool) {
	// 関数領域があるか
	_, ok := st.labels[st.curtFn]
	if !ok {
		return 0, false
	}
	no, ok := st.labels[st.curtFn][label]
	return no, ok
}
