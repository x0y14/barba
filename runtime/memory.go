package runtime

import "fmt"

type Memory []Object

func (m *Memory) Set(offset Offset, obj Object) error {
	if 0 <= offset.Addr() && offset.Addr() < len(*m) {
		(*m)[offset.Addr()] = obj
		return nil
	}
	return fmt.Errorf("offset must be 0< = x < %d", len(*m))
}

func (m *Memory) Get(offset Offset) Object {
	return (*m)[offset.Addr()]
}

func (m *Memory) Delete(offset Offset) {
	(*m)[offset.Addr()] = nil
}

func (m *Memory) IsEmpty(offset Offset) bool {
	return (*m)[offset.Addr()] == nil
}
