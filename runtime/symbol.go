package runtime

import "fmt"

type SymbolTable map[string]Offset

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{}
}

func (s *SymbolTable) Get(name string) (Offset, error) {
	v, ok := (*s)[name]
	if !ok {
		return v, fmt.Errorf("failed to get symbol: not registered: %s", name)
	}
	return v, nil
}

func (s *SymbolTable) Set(name string, offset Offset) error {
	_, err := s.Get(name)
	if err == nil {
		return fmt.Errorf("failed to set symbol: already registered: %s", name)
	}
	(*s)[name] = offset
	return nil
}

func (s *SymbolTable) Delete(name string) {
	delete(*s, name)
}
