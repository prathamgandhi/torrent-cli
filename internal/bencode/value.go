package bencode

import (
	"strconv"
	"strings"
)

type Value interface {
	bencodedValue()
}

type Integer struct {
	Value int64
}

type List struct {
	Value []Value
}

type String struct {
	Value string
}

type Dict struct {
	Value map[string]Value
}

func (String) bencodedValue()  {}
func (Integer) bencodedValue() {}
func (List) bencodedValue()    {}
func (Dict) bencodedValue()    {}


func Describe(v Value) string {
	switch val := v.(type) {
	case String:
		return val.Value

	case Integer:
		return strconv.FormatInt(val.Value, 10)

	case List:
		parts := make([]string, len(val.Value))
		for i, item := range val.Value {
			parts[i] = Describe(item)
		}
		return "[" + strings.Join(parts, ", ") + "]"

	case Dict:
		parts := make([]string, 0, len(val.Value))
		for key, item := range val.Value {
			parts = append(parts, key+": "+Describe(item))
		}
		return "{" + strings.Join(parts, ", ") + "}"
	}
	return ""
}
