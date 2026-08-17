package bencode

import (
	"bytes"
	"testing"
)

func TestEncodeInteger(t *testing.T) {
    tests := []struct {
        in   Value
        want string
    }{
        {Integer{0}, "i0e"},
        {Integer{1}, "i1e"},
        {Integer{-1}, "i-1e"},
        {Integer{42}, "i42e"},
        {Integer{-42}, "i-42e"},
        {Integer{9223372036854775807}, "i9223372036854775807e"},
        {Integer{-9223372036854775808}, "i-9223372036854775808e"},
    }

    for _, tt := range tests {
        got := string(EncodeValue(tt.in))
        if got != tt.want {
            t.Fatalf("got %q want %q", got, tt.want)
        }
    }
}

func TestEncodeString(t *testing.T) {
    tests := []struct {
        in   Value
        want string
    }{
        {String{""}, "0:"},
        {String{"a"}, "1:a"},
        {String{"spam"}, "4:spam"},
        {String{"hello world"}, "11:hello world"},
        {String{"é"}, "2:é"},      // UTF-8 = 2 bytes
        {String{"日本"}, "6:日本"}, // 6 bytes
    }

    for _, tt := range tests {
        got := string(EncodeValue(tt.in))
        if got != tt.want {
            t.Fatalf("got %q want %q", got, tt.want)
        }
    }
}

func TestEncodeBinaryString(t *testing.T) {
    s := string([]byte{'a', 0, 'b'})

    got := EncodeValue(String{s})

    want := append([]byte("3:"), []byte{'a', 0, 'b'}...)

    if !bytes.Equal(got, want) {
        t.Fatal()
    }
}

func TestEncodeEmptyContainers(t *testing.T) {
    if got := string(EncodeValue(List{})); got != "le" {
        t.Fatal(got)
    }

    if got := string(EncodeValue(Dict{})); got != "de" {
        t.Fatal(got)
    }
}

func TestEncodeList(t *testing.T) {
    in := List{
        []Value{
            String{"spam"},
            Integer{42},
            String{"eggs"},
        },
    }

    want := "l4:spami42e4:eggse"

    if got := string(EncodeValue(in)); got != want {
        t.Fatal(got)
    }
}

func TestEncodeNestedLists(t *testing.T) {
    in := List{
        []Value{
            Integer{1},
            List{
                []Value{
                    Integer{2},
                    List{
                        []Value{
                            Integer{3},
                        },
                    },
                },
            },
        },
    }

    want := "li1eli2eli3eeee"

    if got := string(EncodeValue(in)); got != want {
        t.Fatal(got)
    }
}

func TestEncodeDict(t *testing.T) {
    in := Dict{
        map[string]Value{
            "cow":  String{"moo"},
            "spam": String{"eggs"},
        },
    }

    want := "d3:cow3:moo4:spam4:eggse"

    if got := string(EncodeValue(in)); got != want {
        t.Fatal(got)
    }
}

func TestDictOrdering(t *testing.T) {
    in := Dict{
        map[string]Value{
            "b":  String{"1"},
            "aa": String{"2"},
            "a":  String{"0"},
        },
    }

    want := "d1:a1:02:aa1:21:b1:1e"

    if got := string(EncodeValue(in)); got != want {
        t.Fatalf("got %q", got)
    }
}

func TestDictionaryByteOrdering(t *testing.T) {
    in := Dict{
        map[string]Value{
            "ä": String{"1"},
            "z": String{"2"},
        },
    }

    got := string(EncodeValue(in))

    // UTF-8("z") = 0x7A
    // UTF-8("ä") = 0xC3 0xA4
    // therefore "z" sorts first.
    want := "d1:z1:22:ä1:1e"

    if got != want {
        t.Fatal(got)
    }
}

func TestNestedDict(t *testing.T) {
    in := Dict{
        map[string]Value{
            "a": Dict{
                map[string]Value{
                    "b": Integer{1},
                },
            },
        },
    }

    want := "d1:ad1:bi1eee"

    if got := string(EncodeValue(in)); got != want {
        t.Fatal(got)
    }
}

func TestMixedStructure(t *testing.T) {
    in := Dict{
        map[string]Value{
            "list": List{
                []Value{
                    Integer{1},
                    Dict{
                        map[string]Value{
                            "x": String{"abc"},
                        },
                    },
                },
            },
        },
    }

    want := "d4:listli1ed1:x3:abceee"

    if got := string(EncodeValue(in)); got != want {
        t.Fatal(got)
    }
}

func TestDeterministicEncoding(t *testing.T) {
    in := Dict{
        map[string]Value{
            "b": Integer{2},
            "a": Integer{1},
            "c": Integer{3},
        },
    }

    first := string(EncodeValue(in))

    for i := 0; i < 1000; i++ {
        got := string(EncodeValue(in))
        if got != first {
            t.Fatal("encoding not deterministic")
        }
    }
}

func TestDeepNesting(t *testing.T) {
    var v Value = Integer{1}

    for i := 0; i < 1000; i++ {
        v = List{[]Value{v}}
    }

    _ = EncodeValue(v)
}

func TestAllByteValues(t *testing.T) {
    data := make([]byte, 256)
    for i := 0; i < 256; i++ {
        data[i] = byte(i)
    }

    s := string(data)

    got := EncodeValue(String{s})

    expected := append([]byte("256:"), data...)

    if !bytes.Equal(got, expected) {
        t.Fatal()
    }
}
