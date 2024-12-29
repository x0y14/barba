package runtime

import "fmt"

type Memory []Object

func NewMemory(size int) *Memory {
	mem := make(Memory, size)
	return &mem
}

func (m *Memory) Set(offset Offset, obj Object) error {
	if 0 <= offset.Value() && offset.Value() < len(*m) {
		(*m)[offset.Value()] = obj
		return nil
	}
	return fmt.Errorf("offset must be 0< = x < %d", len(*m))
}

func (m *Memory) Get(offset Offset) Object {
	return (*m)[offset.Value()]
}

func (m *Memory) Delete(offset Offset) {
	(*m)[offset.Value()] = nil
}

func (m *Memory) IsEmpty(offset Offset) bool {
	return (*m)[offset.Value()] == nil
}
