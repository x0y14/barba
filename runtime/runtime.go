package runtime

type Runtime struct {
	program []Object
	sym     SymbolTable
	reg     []Object
	stack   []Object
	mem     Memory
}

func NewRuntime(stackSize, memSize int) *Runtime {
	return &Runtime{
		program: nil,
		sym:     *NewSymbolTable(),
		reg:     *NewRegisterSet(),
		stack:   make([]Object, stackSize),
		mem:     *NewMemory(memSize),
	}
}
