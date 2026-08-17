package bencode

import (
	"sort"
	"strconv"
)


func EncodeValue(in Value) ([]byte) {
	switch v := in.(type) {
	case String:
		return encodeString(v)
	case Integer:
		return encodeInteger(v)
	case List:
		return encodeList(v)
	case Dict:
		return encodeDict(v)
	default:
    	return nil
	}
}

func encodeString(in String) ([]byte) {
	return []byte(strconv.Itoa(len(in.Value)) + ":" + in.Value)
}

func encodeInteger(in Integer) ([]byte) {
	return []byte("i" + strconv.FormatInt(in.Value, 10) + "e")
}

func encodeList(in List) ([]byte) {
	encodedList := []byte{}
	for _, v := range(in.Value) {
		encodedValue := EncodeValue(v)
		encodedList = append(encodedList, encodedValue...)
	}
	result := make([]byte, 0, len(encodedList)+2)
	result = append(result, 'l')
	result = append(result, encodedList...)
	result = append(result, 'e')
	return result
}

func encodeDict(in Dict) ([]byte) {
	encodedDict := []byte{}
	keys := make([]string, 0, len(in.Value))
	for k := range(in.Value) {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		byteKey := encodeString(String{ Value: k })
		byteValue := EncodeValue(in.Value[k])
		encodedDict = append(encodedDict, byteKey...)
		encodedDict = append(encodedDict, byteValue...)
	}
	result := make([]byte, 0, len(encodedDict)+2)
	result = append(result, 'd')
	result = append(result, encodedDict...)
	result = append(result, 'e')
	return result
}