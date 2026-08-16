package main

import "simpletorrent/internal/bencode"
import "fmt"

func main() {
	// byteArray := []byte("li32e6:codinge")
	byteArray := []byte("d1:ali32eee")
	str, nextInd, err := bencode.ParseValue(byteArray)
	fmt.Println("Printing answers now: ")
	fmt.Println(bencode.Describe(str))
	fmt.Println("NextIndex: ", nextInd)
	fmt.Println(err)
}