package utils

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

func TestBenStringNonEmpty(t *testing.T) {
	byteArray := []byte("12:hello worlds")
	got, _, _ := ParseString(byteArray)
	want := "hello worlds"

	if got.Value != want {
		t.Errorf("Required result = %s, received %s", want, got)
	}
}

func TestBenStringEmpty(t *testing.T) {
	byteArray := []byte("")
	_, _, err := ParseString(byteArray)

	if err == nil {
		t.Errorf("EOF error absent")
	}
}

func TestBenStringOneDigit(t *testing.T) {
	byteArray := []byte("6:coding")
	got, _, _ := ParseString(byteArray)
	want := "coding"

	if got.Value != want {
		t.Errorf("Required result = %s, received %s", want, got)
	}
}

func TestBenStringAbsentColon(t *testing.T) {
	byteArray := []byte("6coding")
	_, index, err := ParseString(byteArray)

	if err == nil {
		t.Errorf("Should fail with a missing colon")
	}
	if index != -1 {
		t.Errorf("Expecting -1 as next index")
	}
}

func TestBenStringInsufficientLength(t *testing.T) {
	byteArray := []byte("6:codin")
	_, index, err := ParseString(byteArray)

	if err == nil {
		t.Errorf("Should fail with a missing colon")
	}
	if index != -1 {
		t.Errorf("Expecting -1 as next index")
	}
}

func TestBenIntegerValidInput(t *testing.T) {
	byteArray := []byte("i32e")
	got, _, _ := ParseInteger(byteArray)

	want := int64(32)

	if got.Value != want {
		t.Errorf("Received: %d, Expected: %d", got, want)
	}
}

func TestBenIntegerMissingI(t *testing.T) {
	byteArray := []byte("32e")
	_, index, err := ParseInteger(byteArray)

	if err == nil {
		t.Errorf("Should fail with error")
	}
	if index != -1 {
		t.Errorf("Expecting -1 as next index")
	}
}

func TestBenIntegerMissingE(t *testing.T) {
	byteArray := []byte("i32")
	_, index, err := ParseInteger(byteArray)

	if err == nil {
		t.Errorf("Should fail with error")
	}
	if index != -1 {
		t.Errorf("Expecting -1 as next index")
	}
}

func TestBenIntegerNegative(t *testing.T) {
	byteArray := []byte("i-1e")
	got, index, err := ParseInteger(byteArray)

	want := int64(-1)
	if got.Value != want {
		t.Errorf("Failed to parse. Expected %d, got %d", want, got.Value)
	}

	if index != 4 {
		t.Errorf("Invalid next pointer. Expected %d, got %d", 4, index)
	}

	if err != nil {
		t.Errorf("Parsing failed with error: %s", err)
	}
}

func TestBenListEmpty(t *testing.T) {
	byteArray := []byte("le")
	got, i, _ := ParseList(byteArray)
	want := []BencodedValue{}

	if !slices.Equal(got.Value, want) {
		t.Errorf("Output slices don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := 2
	if i != 2 {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenListSingleInteger(t *testing.T) {
	byteArray := []byte("li1ee")
	got, i, _ := ParseList(byteArray)
	want := []BencodedValue{
		BencodedInteger{ Value: 1 },
	}

	if !slices.Equal(got.Value, want) {
		t.Errorf("Output slices don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := 5
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenListNegativeInteger(t *testing.T) {
	byteArray := []byte("li-1ee")
	got, i, _ := ParseList(byteArray)
	want := []BencodedValue{
		BencodedInteger{ Value: -1 },
	}

	if !slices.Equal(got.Value, want) {
		t.Errorf("Output slices don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := 6
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenListNegativeIntegerAndString(t *testing.T) {
	byteArray := []byte("li-1e6:codinge")
	got, i, _ := ParseList(byteArray)
	want := []BencodedValue{
		BencodedInteger{ Value: -1 },
		BencodedString{ Value: "coding" },
	}

	if !slices.Equal(got.Value, want) {
		t.Errorf("Output slices don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	fmt.Println("Expecting index: ", expectedIndex)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenListMultipleStringsInvalid(t *testing.T) {
	byteArray := []byte("l4:star4:wars")
	_, _, err := ParseList(byteArray)
	if err == nil {
		t.Errorf("Expected error because of invalid string, received %s", err)
	}	
}

func TestBenListMultipleStringsValid(t *testing.T) {
	byteArray := []byte("l4:star4:warse")
	got, i, _ := ParseList(byteArray)
	want := []BencodedValue{
		BencodedString{ Value: "star" },
		BencodedString{ Value: "wars" },
	}

	if !slices.Equal(got.Value, want) {
		t.Errorf("Output slices don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenListEmptyStringValid(t *testing.T) {
	byteArray := []byte("l0:e")
	got, i, _ := ParseList(byteArray)
	want := []BencodedValue{
		BencodedString{ Value: "" },
	}

	if !slices.Equal(got.Value, want) {
		t.Errorf("Output slices don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenDictEmpty(t *testing.T) {
	byteArray := []byte("de")
	got, i, _ := ParseDict(byteArray)

	want := map[BencodedString]BencodedValue{}

	if !reflect.DeepEqual(got.Value, want) {
		t.Errorf("Output maps don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenDictSingleInteger(t *testing.T) {
	byteArray := []byte("d1:ai1ee")
	got, i, _ := ParseDict(byteArray)

	want := map[BencodedString]BencodedValue{
		{Value: "a"}: BencodedInteger{Value: 1},
	}

	if !reflect.DeepEqual(got.Value, want) {
		t.Errorf("Output maps don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenDictSingleString(t *testing.T) {
	byteArray := []byte("d3:key5:valuee")
	got, i, _ := ParseDict(byteArray)

	want := map[BencodedString]BencodedValue{
		{Value: "key"}: BencodedString{Value: "value"},
	}

	if !reflect.DeepEqual(got.Value, want) {
		t.Errorf("Output maps don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenDictMultipleEntries(t *testing.T) {
	byteArray := []byte("d1:ai1e1:b3:foo1:ci-5ee")
	got, i, _ := ParseDict(byteArray)

	want := map[BencodedString]BencodedValue{
		{Value: "a"}: BencodedInteger{Value: 1},
		{Value: "b"}: BencodedString{Value: "foo"},
		{Value: "c"}: BencodedInteger{Value: -5},
	}

	if !reflect.DeepEqual(got.Value, want) {
		t.Errorf("Output maps don't match. Expected %v, got %v", want, got.Value)
	}

	expectedIndex := len(byteArray)
	if i != expectedIndex {
		t.Errorf("Next indices don't match. Expected %d, got %d", expectedIndex, i)
	}
}

func TestBenDictMissingLeadingD(t *testing.T) {
	byteArray := []byte("1:ai1ee")

	_, _, err := ParseDict(byteArray)

	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestBenDictNonStringKey(t *testing.T) {
	byteArray := []byte("di1e1:ae")

	_, _, err := ParseDict(byteArray)

	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestBenDictMissingValue(t *testing.T) {
	byteArray := []byte("d1:ae")

	_, _, err := ParseDict(byteArray)

	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

func TestBenDictMissingTrailingE(t *testing.T) {
	byteArray := []byte("d1:ai1e")

	_, _, err := ParseDict(byteArray)

	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
}

