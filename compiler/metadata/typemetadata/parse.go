package typemetadata

import "fmt"

var curt *Constant
var tokens []string

func setAndGoNext(next *Constant) {
	curt.Next = next
	curt = curt.Next
}

func isEof() bool {
	return len(tokens) == 0
}

func expect(s string) error {
	if peek() != s {
		return fmt.Errorf("expected %v, but got %v", s, peek())
	}
	consume()
	return nil
}

func peek() string {
	return tokens[0]
}

func consume() string {
	s := tokens[0]
	tokens = tokens[1:]
	return s
}

func parseArray() error {
	if err := expect("["); err != nil {
		return err
	}
	if err := expect("]"); err != nil {
		return err
	}
	setAndGoNext(NewConstant(Array, nil))
	if err := parse(); err != nil {
		return err
	}
	return nil
}

func parseMap() error {
	if err := expect("map"); err != nil {
		return err
	}
	if err := expect("["); err != nil {
		return err
	}

	setAndGoNext(NewConstant(Map, nil))

	if err := parse(); err != nil {
		return err
	}

	if err := expect("]"); err != nil {
		return err
	}

	return nil
}

func parse() error {
	switch tokens[0] {
	case "[": // array
		if err := parseArray(); err != nil {
			return err
		}
	case "map": // map
		if err := parseMap(); err != nil {
			return err
		}
	case "int":
		setAndGoNext(NewInt())
		consume()
	case "float":
		setAndGoNext(NewFloat())
		consume()
	case "string":
		setAndGoNext(NewString())
		consume()
	case "bool":
		setAndGoNext(NewBool())
		consume()
	default:
		return fmt.Errorf("unsupported token: %v", tokens[0])
	}
	return nil
}

func Parse(tok []string) (*Constant, error) {
	head := NewUndef()
	curt = head

	tokens = tok
	for !isEof() {
		if err := parse(); err != nil {
			return nil, err
		}
	}

	return head.Next, nil
}
