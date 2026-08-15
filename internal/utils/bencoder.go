package utils

import (
	"fmt"
	"io"
	"strconv"
)

type BencodedValue interface {
	PrintValue()
}

type BencodedInteger struct {
	Value int64
}

func (i BencodedInteger) PrintValue() {fmt.Println(i.Value)}


type BencodedList struct {
	Value []BencodedValue
}

func (BencodedList) PrintValue() {}

type BencodedString struct {
	Value string
}

func (s BencodedString) PrintValue() {fmt.Println(s.Value)}

type BencodedDict struct {
	Value map[BencodedString]BencodedValue
}

func (BencodedDict) PrintValue() {}


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

func
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