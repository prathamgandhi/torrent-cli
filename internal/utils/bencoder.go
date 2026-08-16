package utils

import (
	"fmt"
	"io"
	"strconv"
)

type BencodedValue interface {
	PrintValue()
	GetValue() any
}

type BencodedInteger struct {
	Value int64
}

func (i BencodedInteger) PrintValue() {fmt.Println(i.Value)}

func (i BencodedInteger) GetValue()(any) { return i.Value }

type BencodedList struct {
	Value []BencodedValue
}

func (b BencodedList) PrintValue() {
	for _, v := range b.Value {
		fmt.Print(v, " ")
	}
	fmt.Println()
}

func (b BencodedList) GetValue()(any) {
	lissy := []any{}
	for _, v := range b.Value {
		lissy = append(lissy, v)
	}
	return lissy
}

type BencodedString struct {
	Value string
}

func (s BencodedString) PrintValue() {fmt.Println(s.Value)}

func (s BencodedString) GetValue()(any) {
	return s.Value
}

type BencodedDict struct {
	Value map[BencodedString]BencodedValue
}

func (d BencodedDict) PrintValue() {
	for key, value := range d.Value {
		fmt.Println(key.GetValue(), ":", value.GetValue())
	}
}

func (d BencodedDict) GetValue() {
	
}

// ParseString parses a bencoded string from the beginning of in.
// It returns the parsed string, the index immediately after it,
// and an error if the input is invalid or incomplete.
func ParseString(in []byte) (BencodedString, int, error)  {
	if len(in) == 0 {
		return BencodedString{}, -1, io.EOF
	}
	index := 0
	for (index < len(in) && in[index] >= '0' && in[index] <= '9') {
		index++;
	}
	intLength, err := strconv.Atoi(string(in[0:index]))
	if err != nil {
		return BencodedString{}, -1, err
	}
	if index >= len(in) || in[index] != ':' {
		return BencodedString{}, int(-1), fmt.Errorf("Expected : at index: %d", index)
	}
	index++
	if index + intLength > len(in) {
		return BencodedString{}, int(-1), fmt.Errorf("Error while parsing string. Byte slice has input size smaller than the specified length %d", index+intLength)
	}
	// 12:hello worlds
	return BencodedString{ Value: string(in[index:index+intLength]) }, int(index+intLength), nil
}

func ParseInteger(in []byte) (BencodedInteger, int, error) {
	if len(in) == 0 {
		return BencodedInteger{}, -1, io.EOF
	}
	index := 0
	if in[index] != 'i' {
		return BencodedInteger{}, -1, fmt.Errorf("Expected byte 'i', found %d", in[index])
	}

	// startIndex is assumed to be 0
	for index < len(in) && in[index] != 'e' {
		index++
	}
	if index == len(in) {
		return BencodedInteger{}, -1, fmt.Errorf("Reached unexpected end of bytes. Byte 'e' not found")
	}
	convertedInteger, err := strconv.ParseInt(string(in[1:index]), 10, 64)
	if err != nil {
		return BencodedInteger{}, -1, err
	}
	return BencodedInteger{ Value: convertedInteger}, index+1, nil
}

func ParseList(in []byte) (BencodedList, int, error) {
	if len(in) == 0 {
		return BencodedList{}, -1, io.EOF
	}
	index := 0
	if in[index] != 'l' {
		return BencodedList{}, -1, fmt.Errorf("Expected byte 'l', found %d", in[index])
	}
	index++
	list := BencodedList{
		Value: []BencodedValue{},
	}
	for index < len(in) && in[index] != 'e' {
		var (
			value BencodedValue
			nextInd int
			err error
		)
		if in[index] == 'i' {
			value, nextInd, err = ParseInteger(in[index:])
		} else if in[index] >= '0' && in[index] <= '9' {
			value, nextInd, err = ParseString(in[index:])
		} else {
			return BencodedList{}, -1, fmt.Errorf("Reached early EOF")
		}
		if err != nil {
			return BencodedList{}, -1, fmt.Errorf("Error while parsing the list, %s", err)
		}
		list.Value = append(list.Value, value)
		index += nextInd
	}
	if index+1 != len(in) {
		return BencodedList{}, -1, fmt.Errorf("End of list 'e' not found while parsing")
	}
	return list, index+1, nil
}

func ParseDict(in []byte) (BencodedDict, int, error) {
	if len(in) == 0 {
		return BencodedDict{}, -1, io.EOF
	}
	index := 0
	if in[index] != 'd' {
		return BencodedDict{}, -1, fmt.Errorf("Expected byte 'd', found %d", in[index])
	}
	index++
	dict := BencodedDict{
		Value: make(map[BencodedString]BencodedValue),
	}
	for index < len(in) && in[index] != 'e' {
		var (
			value BencodedValue
			key BencodedString
			nextInd int
			err error
		)
		// Parse the key
		if in[index] >= '0' && in[index] <= '9' {
			key, nextInd, err = ParseString(in[index:])
		} else {
			return BencodedDict{}, int(-1), fmt.Errorf("Error parsing the dict. All keys must be string.")
		}
		if err != nil {
			return BencodedDict{}, -1, fmt.Errorf("Error while parsing the dict, %s", err)
		}

		index += nextInd

		// Parse the value
		if index >= len(in) {
			return BencodedDict{}, -1, fmt.Errorf("Reached early EOF")
		}
		if in[index] == 'i' {
			value, nextInd, err = ParseInteger(in[index:])
		} else if in[index] >= '0' && in[index] <= '9' {
			value, nextInd, err = ParseString(in[index:])
		} else {
			return BencodedDict{}, -1, fmt.Errorf("Reached early EOF")
		}
		if err != nil {
			return BencodedDict{}, -1, fmt.Errorf("Error while parsing the dict, %s", err)
		}
		dict.Value[key] = value

		index += nextInd
	}
	if index+1 != len(in) {
		return BencodedDict{}, -1, fmt.Errorf("End of dict 'e' not found while parsing")
	}
	return dict, index+1, nil
}

// func parseInt(in []byte) BencodedInteger {

// }

// func parseDict(in []byte) BencodedDict {

// }

// func parseList(in []byte) BencodedList {

// }

// func parseValue(in []byte) (BencodedValue, error){
// 	if (in[0] >= '0' && in[0] <= '9') {
// 		return parseString(in), nil
// 	}

// 	if (in[0] == 'i') {
// 		return parseInt(in), nil
// 	}

// 	if (in[0] == 'd') {
// 		return parseDict(in), nil
// 	}

// 	if (in[0] == 'l') {
// 		return parseList(in), nil;
// 	}
// 	return nil, errors.New("Failed to parse the byte slice. Invalid input")
// }


// func Decoder(in[] byte) BencodedValue {
// 	if len(in) == 0 {
// 		return BencodedString{ Value: "" }
// 	}

// 	switch (in[0]) {
// 	case 'i':

// 	case 'l':

// 	case 'd':

// 	default:

// 	}
// }

// func Decoder(in[] byte) {
// 	i := 0
// 	for i < len(in) {
// 		switch (in[i]) {
// 			case 'i':
// 				fmt.Println("Integer")

// 			case 'd':
// 				fmt.Println("Dictionary")
// 			case 'l':
// 				fmt.Println("List")
// 			default:
// 				fmt.Println("String")
// 				// 6:coding
// 				// i
// 				// 01234567
// 				// 123:
// 				intStart := i
// 				for in[i] != ':' {
// 					i++;
// 				}
// 				intEnd := i
// 				strLength, err := strconv.Atoi(string(in[intStart:intEnd]))
// 				if err != nil {
// 					panic(err) 
// 				}
// 				fmt.Println(strLength)
// 				// 12:123456789012
// 				acStr := string(in[intEnd+1: intEnd+1+12])
// 				fmt.Println("String: ", acStr)
// 				i := 
// 		}
// 	}
// }