package utils

import "testing"

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
