package bencode

import (
	"fmt"
	"io"
	"strconv"
)


func DecodeValue(in []byte) (Value, int, error) {
	if len(in) == 0 {
		return nil, -1, io.EOF
	}
	switch in[0] {
	case 'd':
		return decodeDict(in)
	case 'l':
		return decodeList(in)
	case 'i':
		return decodeInteger(in)
	default:
		if in[0] >= '0' && in[0] <= '9' {
			return decodeString(in)
		}
	}
	return nil, -1, fmt.Errorf("Invalid bencoded value")
}

// decodeString Decodes a bencoded string from the beginning of in.
// It returns the Decoded string, the index immediately after it,
// and an error if the input is invalid or incomplete.
func decodeString(in []byte) (String, int, error)  {
	if len(in) == 0 {
		return String{}, -1, io.EOF
	}
	index := 0
	for (index < len(in) && in[index] >= '0' && in[index] <= '9') {
		index++;
	}
	intLength, err := strconv.Atoi(string(in[0:index]))
	if err != nil {
		return String{}, -1, err
	}
	if index >= len(in) || in[index] != ':' {
		return String{}, int(-1), fmt.Errorf("Expected : at index: %d", index)
	}
	index++
	if index + intLength > len(in) {
		return String{}, int(-1), fmt.Errorf("Error while parsing string. Byte slice has input size smaller than the specified length %d", index+intLength)
	}
	// 12:hello worlds
	return String{ Value: string(in[index:index+intLength]) }, int(index+intLength), nil
}

func decodeInteger(in []byte) (Integer, int, error) {
	if len(in) == 0 {
		return Integer{}, -1, io.EOF
	}
	index := 0
	if in[index] != 'i' {
		return Integer{}, -1, fmt.Errorf("Expected byte 'i', found %d", in[index])
	}

	// startIndex is assumed to be 0
	for index < len(in) && in[index] != 'e' {
		index++
	}
	if index == len(in) {
		return Integer{}, -1, fmt.Errorf("Reached unexpected end of bytes. Byte 'e' not found")
	}
	convertedInteger, err := strconv.ParseInt(string(in[1:index]), 10, 64)
	if err != nil {
		return Integer{}, -1, err
	}
	return Integer{ Value: convertedInteger}, index+1, nil
}

func decodeList(in []byte) (List, int, error) {
	if len(in) == 0 {
		return List{}, -1, io.EOF
	}
	index := 0
	if in[index] != 'l' {
		return List{}, -1, fmt.Errorf("Expected byte 'l', found %d", in[index])
	}
	index++
	list := List{
		Value: []Value{},
	}
	for index < len(in) && in[index] != 'e' {
		var (
			value Value
			nextInd int
			err error
		)
		value, nextInd, err = DecodeValue(in[index:])
		
		if err != nil {
			return List{}, -1, fmt.Errorf("Error while parsing the list, %s", err)
		}
		list.Value = append(list.Value, value)
		index += nextInd
	}
	if index >= len(in) || in[index] != 'e' {
		return List{}, -1, fmt.Errorf("End of dict 'e' not found while parsing")
	}
	return list, index+1, nil
}

func decodeDict(in []byte) (Dict, int, error) {
	if len(in) == 0 {
		return Dict{}, -1, io.EOF
	}
	index := 0
	if in[index] != 'd' {
		return Dict{}, -1, fmt.Errorf("Expected byte 'd', found %d", in[index])
	}
	index++
	dict := Dict{
		Value: make(map[string]Value),
	}
	for index < len(in) && in[index] != 'e' {
		var (
			value Value
			key String
			nextInd int
			err error
		)
		// Decode the key
		if in[index] >= '0' && in[index] <= '9' {
			key, nextInd, err = decodeString(in[index:])
		} else {
			return Dict{}, int(-1), fmt.Errorf("Error parsing the dict. All keys must be string.")
		}
		if err != nil {
			return Dict{}, -1, fmt.Errorf("Error while parsing the dict, %s", err)
		}

		index += nextInd

		// Decode the value
		if index >= len(in) {
			return Dict{}, -1, fmt.Errorf("Reached early EOF")
		}
		value, nextInd, err = DecodeValue(in[index:])
		
		if err != nil {
			return Dict{}, -1, fmt.Errorf("Error while parsing the dict, %s", err)
		}
		dict.Value[key.Value] = value

		index += nextInd
	}
	if index >= len(in) || in[index] != 'e' {
		return Dict{}, -1, fmt.Errorf("End of dict 'e' not found while parsing")
	}
	return dict, index+1, nil
}
