package bencode


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
	want := []Value{}

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
	want := []Value{
		Integer{ Value: 1 },
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
	want := []Value{
		Integer{ Value: -1 },
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
	want := []Value{
		Integer{ Value: -1 },
		String{ Value: "coding" },
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
	want := []Value{
		String{ Value: "star" },
		String{ Value: "wars" },
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
	want := []Value{
		String{ Value: "" },
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

	want := map[string]Value{}

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

	want := map[string]Value{
		"a": Integer{Value: 1},
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

	want := map[string]Value{
		"key": String{Value: "value"},
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

	want := map[string]Value{
		"a": Integer{Value: 1},
		"b": String{Value: "foo"},
		"c": Integer{Value: -5},
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

func TestParseValueNestedLists(t *testing.T) {
	input := []byte("lli1ei2eei3ee")

	got, i, err := ParseValue(input)
	if err != nil {
		t.Fatal(err)
	}

	want := List{
		Value: []Value{
			List{
				Value: []Value{
					Integer{1},
					Integer{2},
				},
			},
			Integer{3},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %#v, got %#v", want, got)
	}

	if i != len(input) {
		t.Errorf("expected index %d got %d", len(input), i)
	}
}

func TestParseValueDictWithList(t *testing.T) {
	input := []byte("d4:listli1e3:abcee")

	got, _, err := ParseValue(input)
	if err != nil {
		t.Fatal(err)
	}

	want := Dict{
		Value: map[string]Value{
			"list": List{
				Value: []Value{
					Integer{1},
					String{"abc"},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %#v got %#v", want, got)
	}
}

func TestParseValueListOfDicts(t *testing.T) {
	input := []byte("ld1:ai1eed1:bi2eee")

	got, _, err := ParseValue(input)
	if err != nil {
		t.Fatal(err)
	}

	want := List{
		Value: []Value{
			Dict{
				Value: map[string]Value{
					"a": Integer{1},
				},
			},
			Dict{
				Value: map[string]Value{
					"b": Integer{2},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fail()
	}
}

func TestParseValueNestedDicts(t *testing.T) {
	input := []byte("d5:innerd1:ai42eee")

	got, _, err := ParseValue(input)
	if err != nil {
		t.Fatal(err)
	}

	want := Dict{
		Value: map[string]Value{
			"inner": Dict{
				Value: map[string]Value{
					"a": Integer{42},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fail()
	}
}

func TestParseValueDeepMixed(t *testing.T) {
	input := []byte("d1:ali1ed1:b3:fooeli2ei3eeee")

	got, _, err := ParseValue(input)
	if err != nil {
		t.Fatal(err)
	}

	want := Dict{
		Value: map[string]Value{
			"a": List{
				Value: []Value{
					Integer{1},
					Dict{
						Value: map[string]Value{
							"b": String{"foo"},
						},
					},
					List{
						Value: []Value{
							Integer{2},
							Integer{3},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fail()
	}
}

func TestParseValueEmptyContainers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Value
	}{
		{
			name:  "empty list",
			input: "le",
			want: List{
				Value: []Value{},
			},
		},
		{
			name:  "empty dict",
			input: "de",
			want: Dict{
				Value: map[string]Value{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, i, err := ParseValue([]byte(tt.input))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseValue(%q)\nwant: %#v\ngot:  %#v", tt.input, tt.want, got)
			}

			if i != len(tt.input) {
				t.Errorf("next index: want %d, got %d", len(tt.input), i)
			}
		})
	}
}

func TestParseValueUnterminatedNestedList(t *testing.T) {
	input := []byte("lli1ei2e")

	_, _, err := ParseValue(input)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseValueUnterminatedDict(t *testing.T) {
	input := []byte("d1:ad1:bi2ee")

	_, _, err := ParseValue(input)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseValueInvalidNestedValue(t *testing.T) {
	input := []byte("li1exe")

	_, _, err := ParseValue(input)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseValueVeryDeep(t *testing.T) {
	input := []byte("d1:ad1:bd1:cli1ei2e2:hieeee")

	got, _, err := ParseValue(input)
	if err != nil {
		t.Fatal(err)
	}

	want := Dict{
		Value: map[string]Value{
			"a": Dict{
				Value: map[string]Value{
					"b": Dict{
						Value: map[string]Value{
							"c": List{
								Value: []Value{
									Integer{1},
									Integer{2},
									String{"hi"},
								},
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fail()
	}
}
