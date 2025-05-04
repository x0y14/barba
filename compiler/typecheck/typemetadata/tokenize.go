package typemetadata

import "fmt"

var pos int

func startWith(target, v []rune) bool {
	for i, r := range v {
		if pos+i >= len(target) {
			return false
		}
		if target[pos+i] != r {
			return false
		}
	}
	return true
}

func advance(n int) {
	pos += n
}

func Tokenize(runes []rune) ([]string, error) {
	pos = 0
	var result []string
	for len(runes) > pos {
		switch {
		case startWith(runes, []rune("int")):
			result = append(result, "int")
			advance(len("int"))
		case startWith(runes, []rune("float")):
			result = append(result, "float")
			advance(len("float"))
		case startWith(runes, []rune("string")):
			result = append(result, "string")
			advance(len("string"))
		case startWith(runes, []rune("bool")):
			result = append(result, "bool")
			advance(len("bool"))
		case startWith(runes, []rune("array")):
			result = append(result, "array")
			advance(len("array"))
		case startWith(runes, []rune("map")):
			result = append(result, "map")
			advance(len("map"))
		case startWith(runes, []rune("[")):
			result = append(result, "[")
			advance(len("["))
		case startWith(runes, []rune("]")):
			result = append(result, "]")
			advance(len("]"))
		default:
			return nil, fmt.Errorf("unsupported rune: pos=%v", pos)
		}
	}
	return result, nil
}
